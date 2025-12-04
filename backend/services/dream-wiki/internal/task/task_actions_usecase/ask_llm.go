package task_actions_usecase

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/ycloud_client_gen"
)

func (u *taskActionUsecaseImpl) executeAskLLMAction(repo repository.AppRepository, actionID internals.TaskActionID, taskAction *internals.TaskAction) error {
	askLLMAction, err := taskAction.AsTaskActionAskLLM()
	if err != nil {
		return fmt.Errorf("failed to parse task action as TaskActionAskLLM: %w", err)
	}

	var messages []ycloud_client_gen.Message
	for _, msg := range askLLMAction.Messages {
		messages = append(messages, ycloud_client_gen.Message{
			Role: ycloud_client_gen.MessageRole(msg.Role),
			Text: msg.Content,
		})
	}

	operationID, err := u.deps.YCloudClient.StartAsyncLLMRequest(u.ctx, messages)
	if err != nil {
		return fmt.Errorf("failed to start async LLM request: %w", err)
	}

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

			if operation.Done {
				break pollingLoop
			}
		}
	}

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

	err = repo.EnqueueTaskActionResult(actionID)
	if err != nil {
		return fmt.Errorf("failed to enqueue task action result: %w", err)
	}

	return nil
}
