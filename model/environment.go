// Package model provides data structures and types used throughout the application
package model

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

// Environment holds the application's runtime configuration and state.
type Environment struct {
	DbPool                                    *pgxpool.Pool
	PreviousActionRunsFailureTotal            map[string]map[string]int
	CurrentActionRunsFailureTotal             map[string]map[string]int
	ActionRunsFailureTotal                    *prometheus.CounterVec
	PreviousActionRunsNotSuccessTotal         map[string]map[string]int
	CurrentActionRunsNotSuccessTotal          map[string]map[string]int
	ActionRunsNotSuccessTotal                 *prometheus.CounterVec
	PreviousActionRunsFailureOrCancelledTotal map[string]map[string]int
	CurrentActionRunsFailureOrCancelledTotal  map[string]map[string]int
	ActionRunsFailureOrCancelledTotal         *prometheus.CounterVec
	Registry                                  *prometheus.Registry
}
