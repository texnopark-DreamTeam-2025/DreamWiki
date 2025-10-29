package github_account_pr

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

type (
	gitHubAccountPRTask struct {
		state internals.TaskStateGitHubAccountPR
	}
)

var (
	_ task_common.TaskLogic
)

func NewGitHubAccountPRTask(state internals.TaskStateGitHubAccountPR) *gitHubAccountPRTask {
	return &gitHubAccountPRTask{
		state: state,
	}
}

func (t *gitHubAccountPRTask) CalculateSubtasks() ([]api.Subtask, error) {
	// For now, return an empty list of subtasks
	return []api.Subtask{}, nil
}

func (t *gitHubAccountPRTask) OnActionResult(result internals.TaskActionResult) error {
	// For now, do nothing
	return nil
}
