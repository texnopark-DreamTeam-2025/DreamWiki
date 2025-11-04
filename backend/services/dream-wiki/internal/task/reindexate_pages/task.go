package reindexate_pages

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

type (
	reindexatePagesTask struct {
		taskID api.TaskID
		status api.TaskStatus
		state  internals.TaskStateReindexatePages
		ctx    context.Context
		deps   *deps.Deps
		repo   repository.AppRepository
	}
)

var (
	_ task_common.TaskLogic = (*reindexatePagesTask)(nil)
)

func NewReindexatePagesTask(ctx context.Context, state internals.TaskStateReindexatePages, deps *task_common.TaskDeps) *reindexatePagesTask {
	return &reindexatePagesTask{
		state:  state,
		status: deps.Digest.Status,
		ctx:    ctx,
		deps:   deps.Deps,
		taskID: deps.Digest.TaskId,
		repo:   deps.Repo,
	}
}

func (t *reindexatePagesTask) updateState() error {
	taskState := internals.TaskState{}
	taskState.FromTaskStateReindexatePages(t.state)
	return t.repo.SetTaskState(t.taskID, taskState)
}

func (t *reindexatePagesTask) createIndexatePageTaskAction(pageID api.PageID) error {
	taskAction := internals.TaskAction{}
	taskAction.FromTaskActionIndexatePage(internals.TaskActionIndexatePage{
		TaskActionType: internals.IndexatePage,
		PageId:         pageID,
	})

	taskActionID, err := t.repo.CreateTaskAction(t.taskID, taskAction)
	if err != nil {
		return err
	}

	return t.repo.EnqueueTaskAction(*taskActionID)
}

func (t *reindexatePagesTask) CalculateSubtasks() ([]api.Subtask, error) {
	var subtasks []api.Subtask

	for _, pageID := range t.state.IndexatedPageIds {
		pageTitle := t.state.PageTitles[pageID.String()]
		subtasks = append(subtasks, api.Subtask{
			Description: "Indexating page " + pageTitle,
			Status:      api.Done,
			Subsubtasks: []api.SubSubtask{},
		})
	}

	if t.state.PagesToIndexateIds != nil {
		for _, pageID := range t.state.PagesToIndexateIds[len(t.state.IndexatedPageIds):] {
			pageTitle := t.state.PageTitles[pageID.String()]
			subtasks = append(subtasks, api.Subtask{
				Description: "Indexating page " + pageTitle,
				Status:      api.Executing,
				Subsubtasks: []api.SubSubtask{},
			})
		}
	}

	return subtasks, nil
}

func (t *reindexatePagesTask) OnActionResult(result internals.TaskActionResult) error {
	resultType, err := result.Discriminator()
	if err != nil {
		return err
	}

	switch internals.TaskActionType(resultType) {
	case internals.NewTask:
		if len(t.state.PagesToIndexateIds) > 0 {
			firstPageID := t.state.PagesToIndexateIds[0]
			err := t.createIndexatePageTaskAction(firstPageID)
			if err != nil {
				return err
			}
		}
		return t.repo.Commit()

	case internals.IndexatePage:
		indexatePageResult, err := result.AsTaskActionResultIndexatePage()
		if err != nil {
			return err
		}

		pageID := indexatePageResult.PageId
		t.state.IndexatedPageIds = append(t.state.IndexatedPageIds, pageID)

		err = t.updateState()
		if err != nil {
			return err
		}

		if len(t.state.PagesToIndexateIds) >= len(t.state.IndexatedPageIds) {
			err := t.repo.SetTaskStatus(t.taskID, api.Done)
			if err != nil {
				return err
			}
		} else {
			nextPageID := t.state.PagesToIndexateIds[len(t.state.IndexatedPageIds)]
			err = t.createIndexatePageTaskAction(nextPageID)
			if err != nil {
				return err
			}
		}

		return t.repo.Commit()
	}

	return nil
}
