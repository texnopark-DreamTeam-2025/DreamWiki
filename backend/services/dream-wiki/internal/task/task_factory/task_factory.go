package task_factory

import (
	"context"
	"fmt"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/github_account_pr"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/reindexate_all_pages"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

func CreateTaskLogicCreator() task_common.TaskLogicCreator {
	return func(ctx context.Context, deps *task_common.TaskDeps) (task_common.TaskLogic, error) {
		taskType, err := deps.State.Discriminator()
		if err != nil {
			return nil, err
		}

		switch internals.TaskType(taskType) {
		case internals.ReindexateAllPages:
			taskState, err := deps.State.AsTaskStateReindexateAllPages()
			if err != nil {
				return nil, err
			}
			task := reindexate_all_pages.NewReindexateAllPagesTask(taskState)
			if task == nil {
				return nil, fmt.Errorf("task is nil")
			}
			return task, nil

		case internals.GithubAccountPr:
			taskState, err := deps.State.AsTaskStateGitHubAccountPR()
			if err != nil {
				return nil, err
			}
			task := github_account_pr.NewGitHubAccountPRTask(ctx, taskState, deps)
			if task == nil {
				return nil, fmt.Errorf("task is nil")
			}
			return task, nil
		}
		return nil, fmt.Errorf("unknown task type")
	}
}
