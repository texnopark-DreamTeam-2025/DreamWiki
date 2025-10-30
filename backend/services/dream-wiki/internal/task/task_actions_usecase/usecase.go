package task_actions_usecase

import (
	"context"
	"fmt"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
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

func (u *taskActionUsecaseImpl) ExecuteAction(actionID internals.TaskActionID) error {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	taskAction, _, err := repo.GetTaskActionByID(actionID)
	if err != nil {
		return fmt.Errorf("failed to get task action by ID: %w", err)
	}

	actionType, err := taskAction.Discriminator()
	if err != nil {
		return fmt.Errorf("failed to get task action type: %w", err)
	}

	switch actionType {
	case string(internals.NewTask):
		err = u.executeNewTaskAction(repo, actionID, taskAction)
	default:
		err = fmt.Errorf("unsupported task action type: %s", actionType)
	}

	if err != nil {
		return err
	}

	return repo.Commit()
}

func (u *taskActionUsecaseImpl) executeNewTaskAction(repo repository.AppRepository, actionID internals.TaskActionID, taskAction *internals.TaskAction) error {
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
