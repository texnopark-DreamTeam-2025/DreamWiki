package task_actions_usecase

import (
	"context"
	"fmt"
	"runtime"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

type (
	TaskActionUsecase interface {
		ExecuteAction(actionID internals.TaskActionID) error
	}

	taskActionUsecaseImpl struct {
		ctx  context.Context
		deps *deps.Deps
		log  logger.Logger
	}
)

func NewTaskActionUsecase(ctx context.Context, deps *deps.Deps) TaskActionUsecase {
	return &taskActionUsecaseImpl{
		ctx:  ctx,
		deps: deps,
		log:  deps.Logger,
	}
}

func (u *taskActionUsecaseImpl) failTaskActionAndTask(repo repository.AppRepository, actionID internals.TaskActionID, taskID api.TaskID) {
	setErr := repo.SetTaskActionStatus(actionID, internals.Failed)
	if setErr != nil {
		u.log.Error("failed to set task action status to failed", "action_id", actionID, "error", setErr)
	}

	setTaskErr := repo.SetTaskStatus(taskID, api.FailedByError)
	if setTaskErr != nil {
		u.log.Error("failed to set task status to failed_by_error", "task_id", taskID, "error", setTaskErr)
	}

	commitErr := repo.Commit()
	if commitErr != nil {
		u.log.Error("failed to commit transaction", "action_id", actionID, "error", commitErr)
	}
}

func (u *taskActionUsecaseImpl) ExecuteAction(actionID internals.TaskActionID) (err error) {
	repo := repository.NewAppRepository(u.ctx, &deps.RepositoryDeps{})
	defer repo.Rollback()

	defer func() {
		if r := recover(); r != nil {
			stackBuf := make([]byte, 4096)
			stackBuf = stackBuf[:runtime.Stack(stackBuf, false)]
			u.log.Error("panic recovered in task action execution",
				"action_id", actionID,
				"panic", r,
				"stack", string(stackBuf))

			_, taskActionAdditionalInfo, getErr := repo.GetTaskActionByID(actionID)
			if getErr != nil {
				u.log.Error("failed to get task action info for panic handling", "error", getErr)
				err = fmt.Errorf("panic occurred: %v, and failed to get task info: %w", r, getErr)
				return
			}

			u.failTaskActionAndTask(repo, actionID, taskActionAdditionalInfo.TaskId)
			err = fmt.Errorf("panic occurred in task action execution: %v", r)
		}
	}()

	taskAction, taskActionAdditionalInfo, err := repo.GetTaskActionByID(actionID)
	if err != nil {
		return fmt.Errorf("failed to get task action by ID: %w", err)
	}

	taskDigest, _, err := repo.GetTaskByID(taskActionAdditionalInfo.TaskId)
	if err != nil {
		return fmt.Errorf("failed to get task by ID: %w", err)
	}

	if task_common.IsTerminalTaskStatus(taskDigest.Status) {
		u.log.Info("skipping task action execution because task is already in terminal status",
			"action_id", actionID,
			"task_id", taskActionAdditionalInfo.TaskId,
			"task_status", taskDigest.Status)
		return nil
	}

	actionType, err := taskAction.Discriminator()
	if err != nil {
		u.failTaskActionAndTask(repo, actionID, taskActionAdditionalInfo.TaskId)
		return fmt.Errorf("failed to get task action type: %w", err)
	}

	switch internals.TaskActionType(actionType) {
	case internals.NewTask:
		err = u.executeNewTaskAction(repo, actionID, taskAction)
	case internals.AskLlm:
		err = u.executeAskLLMAction(repo, actionID, taskAction)
	case internals.IndexatePage:
		err = u.executeIndexatePageAction(repo, actionID, taskAction)
	default:
		err = fmt.Errorf("unsupported task action type: %s", actionType)
	}

	if err != nil {
		u.failTaskActionAndTask(repo, actionID, taskActionAdditionalInfo.TaskId)
		return err
	}

	return repo.Commit()
}

func (u *taskActionUsecaseImpl) executeNewTaskAction(repo repository.AppRepository, actionID internals.TaskActionID, _ *internals.TaskAction) error {
	err := repo.SetTaskActionStatus(actionID, internals.Finished)
	if err != nil {
		return fmt.Errorf("failed to set task action status to finished: %w", err)
	}

	result := internals.TaskActionResult{}
	newTaskResult := internals.TaskActionResultNewTask{
		TaskActionType: internals.NewTask,
	}
	err = result.FromTaskActionResultNewTask(newTaskResult)
	if err != nil {
		return fmt.Errorf("failed to create task action result: %w", err)
	}

	err = repo.CreateTaskActionResult(actionID, result)
	if err != nil {
		return fmt.Errorf("failed to create task action result: %w", err)
	}

	err = repo.EnqueueTaskActionResult(actionID)
	if err != nil {
		return fmt.Errorf("failed to enqueue task action result: %w", err)
	}

	return nil
}
