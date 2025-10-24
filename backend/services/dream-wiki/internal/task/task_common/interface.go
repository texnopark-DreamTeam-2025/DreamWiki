package task_common

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

type (
	TaskWrapper interface {
		FailByActionTimeout() error
		FailByActionError() error
		Cancel() error

		GetStatus() api.TaskStatus
		CalculateSubtasks() ([]api.Subtask, error)

		Logic() TaskLogic
	}

	TaskLogic interface {
		OnActionResult(result internals.TaskActionResult) error
	}
)
