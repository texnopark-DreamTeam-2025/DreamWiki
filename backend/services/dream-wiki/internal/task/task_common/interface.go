package task_common

import (
	"context"
	"fmt"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
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
		repo      repository.AppRepository
	}

	TaskLogic interface {
		CalculateSubtasks() ([]api.Subtask, error)
		OnActionResult(result internals.TaskActionResult) error
	}

	// TaskLogicCreator is a function that creates a TaskLogic based on the task state.
	// Essential to break import cycle.
	TaskLogicCreator = func(ctx context.Context,
		deps *TaskDeps) (TaskLogic, error)

	TaskDeps struct {
		Deps   *deps.Deps
		Digest api.TaskDigest
		State  *internals.TaskState
		Repo   repository.AppRepository
	}
)

var (
	_ Task = (*taskImpl)(nil)
)

func NewTask(ctx context.Context, deps *TaskDeps, taskLogicCreator TaskLogicCreator) Task {
	taskLogic, err := taskLogicCreator(ctx, deps)
	if err != nil {
		panic(err.Error())
	}

	return &taskImpl{
		state:     deps.State,
		status:    deps.Digest.Status,
		taskLogic: taskLogic,
		repo:      deps.Repo,
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
	fmt.Println("ON ACTION RESULT")
	return t.taskLogic.OnActionResult(result)
}
