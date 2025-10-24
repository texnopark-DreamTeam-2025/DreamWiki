package usecase

import "github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"

func (u *appUsecaseImpl) ListTasks(cursor *string) (tasks []api.TaskDigest, newCursor string, err error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) RetryTask(taskID api.TaskID) error {
	panic("unimplemented")
}

func (u *appUsecaseImpl) CancelTask(taskID api.TaskID) error {
	panic("unimplemented")
}

func (u *appUsecaseImpl) GetTaskDetails(taskID api.TaskID) (api.Task, error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) GetTaskInternalState(taskID api.TaskID) (*api.V1TasksInternalStateGetResponse, error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) RecreateTask(taskID api.TaskID) (*api.TaskID, error) {
	panic("unimplemented")
}
