package github_account_pr

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

type (
	gitHubAccountPRTask struct {
		state internals.TaskStateGitHubAccountPR
		ctx   context.Context
		deps  *deps.Deps
	}
)

var (
	_ task_common.TaskLogic = (*gitHubAccountPRTask)(nil)
)

func NewGitHubAccountPRTask(ctx context.Context, deps *deps.Deps, state internals.TaskStateGitHubAccountPR) *gitHubAccountPRTask {
	return &gitHubAccountPRTask{
		state: state,
		ctx:   ctx,
		deps:  deps,
	}
}

func (t *gitHubAccountPRTask) CalculateSubtasks() ([]api.Subtask, error) {
	// For now, return an empty list of subtasks
	return []api.Subtask{}, nil
}

func (t *gitHubAccountPRTask) OnActionResult(result internals.TaskActionResult) error {
	// TODO: Implement the full logic
	return nil
}
