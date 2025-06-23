// Package collector provides functionality for collecting metrics from Gitea Actions
package collector

import (
	"gitea_actions_prometheus_exporter/model"
)

const (
	// UnknownRepoName is used when a repository name is not available.
	UnknownRepoName = "unknown"
	// UnknownWorkflowID is used when a workflow ID is not available.
	UnknownWorkflowID = "unknown"
)

// UpdateMetrics updates all Prometheus metrics with the latest data.
func UpdateMetrics(env *model.Environment, actionRuns []ActionRun) {
	updateActionRunsFailureTotal(env, actionRuns)
	updateActionRunsNotSuccessTotal(env, actionRuns)
	updateActionRunsFailureOrCancelledTotal(env, actionRuns)
}

func updateActionRunsFailureTotal(env *model.Environment, actionRuns []ActionRun) {
	env.PreviousActionRunsFailureTotal = env.CurrentActionRunsFailureTotal

	// Group failures by repository name and workflow ID
	failuresByRepoAndWorkflow := mapMatchingActionRunsByRepositoryNameAndWorkflowID(
		actionRuns,
		func(ar ActionRun) bool {
			return *ar.Status == StatusFailure
		})

	env.CurrentActionRunsFailureTotal = failuresByRepoAndWorkflow

	// Update the counter with repository_name and workflow_id labels
	for repoName, workflows := range failuresByRepoAndWorkflow {
		for workflowID := range workflows {
			delta := env.CurrentActionRunsFailureTotal[repoName][workflowID] -
				env.PreviousActionRunsFailureTotal[repoName][workflowID]
			env.ActionRunsFailureTotal.WithLabelValues(repoName, workflowID).Add(float64(delta))
		}
	}
}

func updateActionRunsNotSuccessTotal(env *model.Environment, actionRuns []ActionRun) {
	env.PreviousActionRunsNotSuccessTotal = env.CurrentActionRunsNotSuccessTotal

	// Group not success by repository name and workflow ID
	notSuccessByRepoAndWorkflow := mapMatchingActionRunsByRepositoryNameAndWorkflowID(
		actionRuns,
		func(ar ActionRun) bool {
			return *ar.Status != StatusSuccess && *ar.Stopped > 0
		})

	env.CurrentActionRunsNotSuccessTotal = notSuccessByRepoAndWorkflow

	// Update the counter with repository_name and workflow_id labels
	for repoName, workflows := range notSuccessByRepoAndWorkflow {
		for workflowID := range workflows {
			delta := env.CurrentActionRunsNotSuccessTotal[repoName][workflowID] -
				env.PreviousActionRunsNotSuccessTotal[repoName][workflowID]
			env.ActionRunsNotSuccessTotal.WithLabelValues(repoName, workflowID).Add(float64(delta))
		}
	}
}

func updateActionRunsFailureOrCancelledTotal(env *model.Environment, actionRuns []ActionRun) {
	env.PreviousActionRunsFailureOrCancelledTotal = env.CurrentActionRunsFailureOrCancelledTotal

	// Group failures or cancelled by repository name and workflow ID
	failureOrCancelledByRepoAndWorkflow := mapMatchingActionRunsByRepositoryNameAndWorkflowID(
		actionRuns,
		func(ar ActionRun) bool {
			return *ar.Status == StatusFailure || *ar.Status == StatusCancelled
		})

	env.CurrentActionRunsFailureOrCancelledTotal = failureOrCancelledByRepoAndWorkflow

	// Update the counter with repository_name and workflow_id labels
	for repoName, workflows := range failureOrCancelledByRepoAndWorkflow {
		for workflowID := range workflows {
			delta := env.CurrentActionRunsFailureOrCancelledTotal[repoName][workflowID] -
				env.PreviousActionRunsFailureOrCancelledTotal[repoName][workflowID]
			env.ActionRunsFailureOrCancelledTotal.WithLabelValues(repoName, workflowID).
				Add(float64(delta))
		}
	}
}
