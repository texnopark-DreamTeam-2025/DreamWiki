package task_factory

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/github_account_pr"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/reindexate_all_pages"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

// CreateTaskLogicCreator creates a taskLogicCreator function based on the task state
func CreateTaskLogicCreator(ctx context.Context, deps *deps.Deps, state *internals.TaskState) task_common.TaskLogicCreator {
	return func(state *internals.TaskState) task_common.TaskLogic {
		// Get the task type from the state
		taskType, err := state.Discriminator()
		if err != nil {
			// Return a default task logic if we can't determine the type
			return &defaultTaskLogic{}
		}

		// Create the appropriate task logic based on the task type
		switch taskType {
		case "reindexate_all_pages":
			// We need to set the state in the task
			if taskState, err := state.AsTaskStateReindexateAllPages(); err == nil {
				task := reindexate_all_pages.NewReindexateAllPagesTask(taskState)
				return task
			}
			// If we can't get the state, return a default task
			return &defaultTaskLogic{}
		case "github_account_pr":
			// We need to set the state in the task
			if taskState, err := state.AsTaskStateGitHubAccountPR(); err == nil {
				task := github_account_pr.NewGitHubAccountPRTask(ctx, deps, taskState)
				return task
			}
			// If we can't get the state, return a default task
			return &defaultTaskLogic{}
		default:
			return &defaultTaskLogic{}
		}
	}
}

// defaultTaskLogic is a default implementation that does nothing
type defaultTaskLogic struct{}

func (d *defaultTaskLogic) CalculateSubtasks() ([]api.Subtask, error) {
	return []api.Subtask{}, nil
}

func (d *defaultTaskLogic) OnActionResult(result internals.TaskActionResult) error {
	return nil
}
