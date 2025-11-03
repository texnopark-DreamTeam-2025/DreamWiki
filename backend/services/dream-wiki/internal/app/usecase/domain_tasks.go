package usecase

import (
	"encoding/json"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_factory"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (u *appUsecaseImpl) ListTasks(cursor *api.Cursor) (tasks []api.TaskDigest, newCursor *api.NextInfo, err error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	taskDigests, _, newCursor, err := repo.ListTasks(cursor, 20)
	if err != nil {
		return nil, nil, err
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

	// Get task digest and state from repository
	taskDigest, taskState, err := repo.GetTaskByID(taskID)
	if err != nil {
		return api.Task{}, err
	}

	// Create task logic creator
	taskLogicCreator := task_factory.CreateTaskLogicCreator()

	// Create task instance
	task := task_common.NewTask(u.ctx, &task_common.TaskDeps{
		Deps:   u.deps,
		Digest: *taskDigest,
		State:  taskState,
	}, taskLogicCreator)

	// Calculate subtasks
	subtasks, err := task.CalculateSubtasks()
	if err != nil {
		return api.Task{}, err
	}

	// Create and return the task details
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
