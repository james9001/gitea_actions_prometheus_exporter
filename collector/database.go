// Package collector provides functionality for collecting metrics from Gitea Actions
package collector

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	errOpeningDB      = errors.New("error opening database")
	errConnectingDB   = errors.New("error connecting to database")
	errExecutingQuery = errors.New("error executing query")
	errScanningRow    = errors.New("error scanning row")
	errIteratingRows  = errors.New("error iterating rows")
)

// ActionRun represents a single action run record from the database.
type ActionRun struct {
	ID             int64   `json:"id"`
	Title          *string `json:"title,omitempty"`
	RepoID         *int64  `json:"repoId,omitempty"`
	OwnerID        *int64  `json:"ownerId,omitempty"`
	WorkflowID     *string `json:"workflowId,omitempty"`
	Index          *int64  `json:"index,omitempty"`
	TriggerUserID  *int64  `json:"triggerUserId,omitempty"`
	ScheduleID     *int64  `json:"scheduleId,omitempty"`
	Ref            *string `json:"ref,omitempty"`
	Event          *string `json:"event,omitempty"`
	TriggerEvent   *string `json:"triggerEvent,omitempty"`
	Status         *Status `json:"status,omitempty"`
	Version        *int32  `json:"version,omitempty"`
	Started        *int64  `json:"started,omitempty"`
	Stopped        *int64  `json:"stopped,omitempty"`
	Created        *int64  `json:"created,omitempty"`
	Updated        *int64  `json:"updated,omitempty"`
	RepositoryName *string `json:"repositoryName,omitempty"`
}

// DBConfig holds the database connection configuration.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// GetDBConfig returns a DBConfig with values from environment variables or defaults.
func getDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", ""),
		DBName:   getEnvOrDefault("DB_NAME", "postgres"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

// InitialiseDatabaseConnection initialises the database connection, based on environment variables.
func InitialiseDatabaseConnection() (*pgxpool.Pool, error) {
	config := getDBConfig()

	dbPool, dbErr := connectDB(config)

	if dbErr != nil {
		return nil, dbErr
	}

	return dbPool, nil
}

func connectDB(config DBConfig) (*pgxpool.Pool, error) {
	hostPort := net.JoinHostPort(config.Host, config.Port)
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.User, config.Password, hostPort, config.DBName)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, errors.Join(errOpeningDB, err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, errors.Join(errConnectingDB, err)
	}

	return pool, nil
}

// scanActionRun scans a single row into an ActionRun struct.
func scanActionRun(rows pgx.Rows) (ActionRun, error) {
	var actionRun ActionRun

	var statusInt *int32
	err := rows.Scan(
		&actionRun.ID,
		&actionRun.Title,
		&actionRun.RepoID,
		&actionRun.OwnerID,
		&actionRun.WorkflowID,
		&actionRun.Index,
		&actionRun.TriggerUserID,
		&actionRun.ScheduleID,
		&actionRun.Ref,
		&actionRun.Event,
		&actionRun.TriggerEvent,
		&statusInt,
		&actionRun.Version,
		&actionRun.Started,
		&actionRun.Stopped,
		&actionRun.Created,
		&actionRun.Updated,
		&actionRun.RepositoryName,
	)

	if err != nil {
		return ActionRun{}, errors.Join(errScanningRow, err)
	}

	// Convert int32 to Status
	if statusInt != nil {
		status := Status(*statusInt)
		actionRun.Status = &status
	}

	return actionRun, nil
}

// processActionRunRows processes the rows and returns a slice of ActionRun.
func processActionRunRows(rows pgx.Rows) ([]ActionRun, error) {
	var actionRuns []ActionRun

	for rows.Next() {
		actionRun, err := scanActionRun(rows)

		if err != nil {
			return nil, err
		}

		actionRuns = append(actionRuns, actionRun)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Join(errIteratingRows, err)
	}

	return actionRuns, nil
}

// GetActionRuns retrieves all action runs from the database, ordered by ID in descending order.
func GetActionRuns(pool *pgxpool.Pool) ([]ActionRun, error) {
	actionRunsQuery := `
SELECT
    ar.id,
    ar.title,
    ar.repo_id,
    ar.owner_id,
    ar.workflow_id,
    ar.index,
    ar.trigger_user_id,
    ar.schedule_id,
    ar.ref,
    ar.event,
    ar.trigger_event,
    ar.status,
    ar.version,
    ar.started,
    ar.stopped,
    ar.created,
    ar.updated,
    r.name AS repository_name
FROM public.action_run ar
INNER JOIN public.repository r ON ar.repo_id = r.id
ORDER BY ar.id DESC;
	`

	rows, err := pool.Query(context.Background(), actionRunsQuery)
	if err != nil {
		return nil, errors.Join(errExecutingQuery, err)
	}

	defer rows.Close()

	return processActionRunRows(rows)
}
