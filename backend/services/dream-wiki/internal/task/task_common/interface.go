package task_common

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

type (
	Task interface {
		GetStatus() api.TaskStatus
		GetState() *internals.TaskState

		CalculateSubtasks() ([]api.Subtask, error)

		OnActionResult(result internals.TaskActionResult) error
	}

	taskImpl struct {
		state     *internals.TaskState
		status    api.TaskStatus
		taskLogic TaskLogic
	}

	TaskLogic interface {
		CalculateSubtasks() ([]api.Subtask, error)
		OnActionResult(result internals.TaskActionResult) error
	}

	// TaskLogicCreator is a function that creates a TaskLogic based on the task state
	TaskLogicCreator = func(state *internals.TaskState) TaskLogic
)

var (
	_ Task = (*taskImpl)(nil)
)

func NewTask(digest api.TaskDigest, state *internals.TaskState, taskLogicCreator TaskLogicCreator) Task {
	return &taskImpl{
		state:     state,
		status:    digest.Status,
		taskLogic: taskLogicCreator(state),
	}
}

func (t *taskImpl) CalculateSubtasks() ([]api.Subtask, error) {
	return t.taskLogic.CalculateSubtasks()
}

func (t *taskImpl) GetState() *internals.TaskState {
	return t.state
}

func (t *taskImpl) GetStatus() api.TaskStatus {
	return t.status
}

func (t *taskImpl) OnActionResult(result internals.TaskActionResult) error {
	return t.taskLogic.OnActionResult(result)
}
