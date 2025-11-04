package task_actions_usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/ycloud_client_gen"
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

func (u *taskActionUsecaseImpl) ExecuteAction(actionID internals.TaskActionID) error {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	taskAction, taskActionAdditionalInfo, err := repo.GetTaskActionByID(actionID)
	if err != nil {
		return fmt.Errorf("failed to get task action by ID: %w", err)
	}

	actionType, err := taskAction.Discriminator()
	if err != nil {
		u.failTaskActionAndTask(repo, actionID, taskActionAdditionalInfo.TaskId)
		return fmt.Errorf("failed to get task action type: %w", err)
	}

	switch actionType {
	case string(internals.NewTask):
		err = u.executeNewTaskAction(repo, actionID, taskAction)
	case string(internals.AskLlm):
		err = u.executeAskLLMAction(repo, actionID, taskAction)
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

func (u *taskActionUsecaseImpl) executeAskLLMAction(repo repository.AppRepository, actionID internals.TaskActionID, taskAction *internals.TaskAction) error {
	// Parse the task action as TaskActionAskLLM
	askLLMAction, err := taskAction.AsTaskActionAskLLM()
	if err != nil {
		return fmt.Errorf("failed to parse task action as TaskActionAskLLM: %w", err)
	}

	// Convert LLMMessage to ycloud_client_gen.Message
	var messages []ycloud_client_gen.Message
	for _, msg := range askLLMAction.Messages {
		messages = append(messages, ycloud_client_gen.Message{
			Role: ycloud_client_gen.MessageRole(msg.Role),
			Text: msg.Content,
		})
	}

	// Send request to ycloud
	operationID, err := u.deps.YCloudClient.StartAsyncLLMRequest(u.ctx, messages)
	if err != nil {
		return fmt.Errorf("failed to start async LLM request: %w", err)
	}

	// Poll ycloud each 5 seconds until operation status is not done
	// Timeout 3 minutes
	timeout := time.After(3 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var operation *ycloud_client_gen.Operation
pollingLoop:
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout while waiting for LLM response")
		case <-ticker.C:
			operation, err = u.deps.YCloudClient.GetLLMResponse(u.ctx, *operationID)
			u.log.Info(operation)
			if err != nil {
				return fmt.Errorf("failed to get LLM response: %w", err)
			}
			if operation == nil {
				return fmt.Errorf("operation is nil")
			}

			// Check if operation is done
			if operation.Done {
				break pollingLoop
			}
		}
	}

	// Save result
	err = repo.SetTaskActionStatus(actionID, internals.Finished)
	if err != nil {
		return fmt.Errorf("failed to set task action status to finished: %w", err)
	}

	responseBytes, err := json.Marshal(operation)
	if err != nil {
		return err
	}
	rawResponse := make(map[string]any)
	err = json.Unmarshal(responseBytes, &rawResponse)
	if err != nil {
		return err
	}

	if operation.Response == nil {
		return fmt.Errorf("response is nil")
	}

	if len(operation.Response.Alternatives) == 0 {
		return fmt.Errorf("no alternatives")
	}

	result := internals.TaskActionResult{}
	askLLMResult := internals.TaskActionResultAskLLM{
		TaskActionType:  internals.AskLlm,
		ResponseMessage: operation.Response.Alternatives[0].Message.Text,
		ServerResponse:  rawResponse,
	}
	err = result.FromTaskActionResultAskLLM(askLLMResult)
	if err != nil {
		return fmt.Errorf("failed to create task action result: %w", err)
	}

	err = repo.CreateTaskActionResult(actionID, result)
	if err != nil {
		return fmt.Errorf("failed to create task action result: %w", err)
	}

	// Enqueue result
	err = repo.EnqueueTaskActionResult(actionID)
	if err != nil {
		return fmt.Errorf("failed to enqueue task action result: %w", err)
	}

	return nil
}
