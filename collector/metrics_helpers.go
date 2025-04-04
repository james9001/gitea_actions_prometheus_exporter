// Package collector provides functionality for collecting metrics from Gitea Actions
package collector

func mapMatchingActionRunsByRepositoryNameAndWorkflowID(
	actionRuns []ActionRun,
	predicate func(ActionRun) bool,
) map[string]map[string]int {
	predicateMatchesByRepoAndWorkflow := make(map[string]map[string]int)

	for _, run := range actionRuns {
		if run.Status != nil && predicate(run) {
			// Track failures by repository name and workflow ID
			repoName := UnknownRepoName
			if run.RepositoryName != nil {
				repoName = *run.RepositoryName
			}

			workflowID := UnknownWorkflowID
			if run.WorkflowID != nil {
				workflowID = *run.WorkflowID
			}

			if _, exists := predicateMatchesByRepoAndWorkflow[repoName]; !exists {
				predicateMatchesByRepoAndWorkflow[repoName] = make(map[string]int)
			}

			predicateMatchesByRepoAndWorkflow[repoName][workflowID]++
		}
	}

	return predicateMatchesByRepoAndWorkflow
}
