package usecase

import (
	"encoding/json"
	"math"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_factory"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

func makeTaskDescription(taskDigest api.TaskDigest, state internals.TaskState) string {
	if discriminator, _ := state.Discriminator(); internals.TaskType(discriminator) == internals.GithubAccountPr {
		return "Применить PR GitHub"
	}
	if discriminator, _ := state.Discriminator(); internals.TaskType(discriminator) == internals.ReindexatePages {
		return "Проиндексировать страницы"
	}
	return "Какая-то задача"
}

func (u *appUsecaseImpl) ListTasks(cursor *api.Cursor) (tasks []api.TaskDigest, newCursor *api.NextInfo, err error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	taskDigests, taskStates, newCursor, err := repo.ListTasks(cursor, 20)
	if err != nil {
		return nil, nil, err
	}
	for i := range taskDigests {
		taskDigests[i].Description = makeTaskDescription(taskDigests[i], taskStates[i])
		task := task_common.NewTask(u.ctx, &task_common.TaskDeps{
			Deps:   u.deps,
			Digest: taskDigests[i],
			State:  &taskStates[i],
		}, task_factory.CreateTaskLogicCreator())

		subtasks, err := task.CalculateSubtasks()
		if err != nil {
			return nil, nil, err
		}
		doneSubtasks := 0
		for _, subtask := range subtasks {
			if subtask.Status == api.Done {
				doneSubtasks++
			}
		}
		taskDigests[i].ProgressPercentage = int(math.Round(float64(100) * float64(doneSubtasks) / float64(len(subtasks))))
	}

	return taskDigests, newCursor, nil
}

func (u *appUsecaseImpl) RetryTask(taskID api.TaskID) error {
	panic("unimplemented")
}

func (u *appUsecaseImpl) CancelTask(taskID api.TaskID) error {
	panic("unimplemented")
}

func (u *appUsecaseImpl) GetTaskDetails(taskID api.TaskID) (api.Task, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	taskDigest, taskState, err := repo.GetTaskByID(taskID)
	if err != nil {
		return api.Task{}, err
	}

	taskLogicCreator := task_factory.CreateTaskLogicCreator()

	task := task_common.NewTask(u.ctx, &task_common.TaskDeps{
		Deps:   u.deps,
		Digest: *taskDigest,
		State:  taskState,
	}, taskLogicCreator)

	subtasks, err := task.CalculateSubtasks()
	if err != nil {
		return api.Task{}, err
	}

	taskDetails := api.Task{
		CreatedAt:  time.Now(), // TODO: Get actual created time from task state
		Subtasks:   subtasks,
		TaskDigest: *taskDigest,
		UpdatedAt:  time.Now(), // TODO: Get actual updated time from task state
	}

	return taskDetails, nil
}

func (u *appUsecaseImpl) GetTaskInternalState(taskID api.TaskID) (*api.V1TasksInternalStateGetResponse, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	// Get task digest and state from repository
	taskDigest, taskState, err := repo.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	// Convert taskState to RawJSON
	taskStateBytes, err := json.Marshal(taskState)
	if err != nil {
		return nil, err
	}

	taskStateRaw := make(map[string]any)
	if err := json.Unmarshal(taskStateBytes, &taskStateRaw); err != nil {
		return nil, err
	}

	// Create and return the internal state response
	response := &api.V1TasksInternalStateGetResponse{
		Actions:   []api.RawJSON{}, // TODO: include actions in future
		TaskId:    taskDigest.TaskId,
		TaskState: taskStateRaw,
	}

	return response, nil
}

func (u *appUsecaseImpl) RecreateTask(taskID api.TaskID) (*api.TaskID, error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) CreatePageReindexationTask(pageIDs []api.PageID) (*api.TaskID, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	pageTitles := make(map[string]string)
	for _, pageID := range pageIDs {
		page, _, err := repo.GetPageByID(pageID)
		if err != nil {
			return nil, err
		}
		pageTitles[pageID.String()] = page.Title
	}

	taskState := internals.TaskStateReindexatePages{
		PagesToIndexateIds: pageIDs,
		IndexatedPageIds:   []api.PageID{},
		PageTitles:         pageTitles,
		TaskType:           internals.ReindexatePages,
	}

	var taskStateUnion internals.TaskState
	err := taskStateUnion.FromTaskStateReindexatePages(taskState)
	if err != nil {
		return nil, err
	}

	taskID, err := repo.CreateTask(taskStateUnion)
	if err != nil {
		return nil, err
	}

	taskAction := internals.TaskAction{}
	taskAction.FromTaskActionNewTask(internals.TaskActionNewTask{TaskActionType: internals.NewTask})
	taskActionID, err := repo.CreateTaskAction(*taskID, taskAction)
	if err != nil {
		return nil, err
	}

	err = repo.EnqueueTaskAction(*taskActionID)
	if err != nil {
		return nil, err
	}

	err = repo.Commit()
	if err != nil {
		return nil, err
	}

	return taskID, nil
}
