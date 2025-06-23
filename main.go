// Package main provides a Prometheus exporter for Gitea Actions
package main

import (
	"encoding/json"
	"fmt"
	"gitea_actions_prometheus_exporter/collector"
	"gitea_actions_prometheus_exporter/model"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/joho/godotenv"
)

const readHeaderTimeoutSeconds = 3

func handleActionRuns(env *model.Environment) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, _ *http.Request) {
		actionRuns, err := collector.GetActionRuns(env.DbPool)
		if err != nil {
			slog.Error("Failed to get action runs", "error", err)
			http.Error(responseWriter, "Internal server error", http.StatusInternalServerError)

			return
		}

		responseWriter.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(responseWriter).Encode(actionRuns)
		if err != nil {
			slog.Error("Failed to encode action runs", "error", err)
			http.Error(responseWriter, "Internal server error", http.StatusInternalServerError)

			return
		}
	}
}

// startMetricsUpdater starts a goroutine that updates metrics at the specified interval.
func startMetricsUpdater(env *model.Environment, intervalSeconds int) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	updateMetrics := func() {
		for range ticker.C {
			actionRuns, err := collector.GetActionRuns(env.DbPool)
			if err != nil {
				slog.Error("Failed to get action runs for metrics update", "error", err)

				continue
			}

			collector.UpdateMetrics(env, actionRuns)
			slog.Info("Updated Prometheus metrics", "interval_seconds", intervalSeconds)
		}
	}

	go updateMetrics()
}

func initialiseEnvironment() (*model.Environment, error) {
	dbPool, dbErr := collector.InitialiseDatabaseConnection()
	if dbErr != nil {
		return nil, fmt.Errorf("err initialising database connection %w", dbErr)
	}

	// Create a custom registry
	registry := prometheus.NewRegistry()

	env := &model.Environment{
		DbPool:                         dbPool,
		PreviousActionRunsFailureTotal: make(map[string]map[string]int),
		CurrentActionRunsFailureTotal:  make(map[string]map[string]int),
		ActionRunsFailureTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "",
			Subsystem: "",
			Name:      "action_runs_failure_total",
			Help:      "Total number of all action runs with status 'failure'",
		}, []string{"repository_name", "workflow_id"}),
		PreviousActionRunsNotSuccessTotal: make(map[string]map[string]int),
		CurrentActionRunsNotSuccessTotal:  make(map[string]map[string]int),
		ActionRunsNotSuccessTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "",
			Subsystem: "",
			Name:      "action_runs_not_success_total",
			Help:      "Total number of stopped action runs with status that isnt success",
		}, []string{"repository_name", "workflow_id"}),
		PreviousActionRunsFailureOrCancelledTotal: make(map[string]map[string]int),
		CurrentActionRunsFailureOrCancelledTotal:  make(map[string]map[string]int),
		ActionRunsFailureOrCancelledTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "",
			Subsystem: "",
			Name:      "action_runs_failure_or_cancelled_total",
			Help:      "Total number of all action runs with status 'failure' or 'cancelled'",
		}, []string{"repository_name", "workflow_id"}),
		Registry: registry,
	}

	// Register the counter with our custom registry
	registry.MustRegister(env.ActionRunsFailureTotal)
	registry.MustRegister(env.ActionRunsNotSuccessTotal)
	registry.MustRegister(env.ActionRunsFailureOrCancelledTotal)

	return env, nil
}

func configureLogging() {
	if os.Getenv("SLOG_HANDLER") == "JSON" {
		logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})
		slog.SetDefault(slog.New(logHandler))
	}
}

// https://github.com/go-critic/go-critic/issues/1022#issuecomment-1443876315
func mainImpl() error {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("error loading .env file", "error", err)

		return fmt.Errorf("error loading .env file %w", err)
	}

	configureLogging()

	env, envErr := initialiseEnvironment()
	if envErr != nil {
		slog.Error("Failed to prepare environment", "error", envErr)

		return envErr
	}
	defer env.DbPool.Close()

	// Get update interval from environment variable, default to 60 seconds if not set
	updateIntervalStr := os.Getenv("UPDATE_INTERVAL")
	updateInterval := 60 // Default value

	if updateIntervalStr != "" {
		updateInterval, err = strconv.Atoi(updateIntervalStr)
		if err != nil {
			slog.Warn(
				"Invalid UPDATE_INTERVAL value, using default of 60 seconds",
				"value",
				updateIntervalStr,
			)
		}
	}

	// Start the metrics updater
	startMetricsUpdater(env, updateInterval)

	router := mux.NewRouter()
	router.HandleFunc("/action-runs", handleActionRuns(env)).Methods("GET")

	// Add Prometheus metrics endpoint with custom registry
	router.Handle("/metrics", promhttp.HandlerFor(env.Registry, promhttp.HandlerOpts{
		Registry: env.Registry,
	})).Methods("GET")

	slog.Info("Starting up...")
	slog.Info("Listening on port " + os.Getenv("SERVER_PORT"))
	slog.Info("Metrics update interval: " + strconv.Itoa(updateInterval) + " seconds")

	server := &http.Server{
		Addr:              ":" + os.Getenv("SERVER_PORT"),
		ReadHeaderTimeout: readHeaderTimeoutSeconds * time.Second,
		Handler:           router,
	}

	listenErr := server.ListenAndServe()
	if listenErr != nil {
		slog.Error("error when starting http server %w", "error", listenErr)

		return fmt.Errorf("err starting http server %w", listenErr)
	}

	slog.Info("Shutting down exporter, bye bye...")

	return nil
}

func main() {
	err := mainImpl()
	if err != nil {
		os.Exit(1)
	}
}
