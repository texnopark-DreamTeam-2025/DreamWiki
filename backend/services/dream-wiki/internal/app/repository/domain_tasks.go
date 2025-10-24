package repository

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

func (r *appRepositoryImpl) CreateTask(taskState internals.TaskState) (*api.TaskID, error) {
	panic("unimplemented")
}

func (r *appRepositoryImpl) GetTaskByID(taskID api.TaskID) (*api.TaskDigest, *internals.TaskState, error) {
	panic("unimplemented")
}

func (r *appRepositoryImpl) ListTasks(cursor *api.Cursor, limit int64) ([]api.TaskDigest, []internals.TaskState, *api.Cursor, error) {
	panic("unimplemented")
}

func (r *appRepositoryImpl) SetTaskState(taskID api.TaskID, newState internals.TaskState) error {
	panic("unimplemented")
}

func (r *appRepositoryImpl) SetTaskStatus(taskID api.TaskID, newStatus api.TaskStatus) error {
	panic("unimplemented")
}
