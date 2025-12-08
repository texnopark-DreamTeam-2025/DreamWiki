package repository

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicwriter"
)

func (r *appRepositoryImpl) CreateTaskAction(taskID api.TaskID, actionState internals.TaskAction) (*internals.TaskActionID, error) {
	yql := `
	INSERT INTO TaskAction(task_id, status, action, created_at, updated_at)
	VALUES (
		$taskID,
		'new',
		$action,
		CurrentUtcDatetime(),
		CurrentUtcDatetime()
	)
	RETURNING task_action_id;`

	actionBytes, err := json.Marshal(actionState)
	if err != nil {
		return nil, err
	}

	parameters := []table.ParameterOption{
		table.ValueParam("$taskID", types.Int64Value(taskID)),
		table.ValueParam("$action", types.JSONValueFromBytes(actionBytes)),
	}

	result, err := r.tx.InTX().Execute(yql, parameters...)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var taskActionID internals.TaskActionID
	err = result.FetchExactlyOne(&taskActionID)
	if err != nil {
		return nil, err
	}

	r.log.Debug("Inserted task action with id ", taskActionID)

	return &taskActionID, nil
}

func (r *appRepositoryImpl) CreateTaskActionResult(actionID internals.TaskActionID, actionResult internals.TaskActionResult) error {
	yql := `
	INSERT INTO TaskActionResult(task_action_id, result, created_at)
	VALUES (
		$actionID,
		$result,
		CurrentUtcDatetime()
	);`

	resultBytes, err := json.Marshal(actionResult)
	if err != nil {
		return err
	}

	parameters := []table.ParameterOption{
		table.ValueParam("$actionID", types.Int64Value(actionID)),
		table.ValueParam("$result", types.JSONValueFromBytes(resultBytes)),
	}

	result, err := r.tx.InTX().Execute(yql, parameters...)
	if err != nil {
		return err
	}
	defer result.Close()

	r.log.Debug("Inserted task action result for action id ", actionID)

	return nil
}

func (r *appRepositoryImpl) EnqueueTaskAction(actionID internals.TaskActionID) error {
	topicClient := r.tx.TopicClient()
	writer, err := topicClient.StartTransactionalWriter(r.tx.GetTX(), "TaskActionToExecute")
	if err != nil {
		return err
	}
	reader := strings.NewReader(fmt.Sprintf("%d", actionID))
	err = writer.Write(r.ctx, topicwriter.Message{Data: reader})
	if err != nil {
		return err
	}
	r.log.Info("Enqueued message for actionID ", actionID)
	return nil
}

func (r *appRepositoryImpl) EnqueueTaskActionResult(actionID internals.TaskActionID) error {
	topicClient := r.tx.TopicClient()
	writer, err := topicClient.StartTransactionalWriter(r.tx.GetTX(), "TaskActionResultReady")
	if err != nil {
		return err
	}
	err = writer.Write(r.ctx, topicwriter.Message{Data: strings.NewReader(fmt.Sprintf("%d", actionID))})
	if err != nil {
		return err
	}
	r.log.Info("Enqueued result ready message for actionID ", actionID)
	return nil
}

func (r *appRepositoryImpl) GetTaskActionByID(actionID internals.TaskActionID) (*internals.TaskAction, *internals.TaskActionAdditionalInfo, error) {
	yql := `
	SELECT
		task_action_id,
		task_id,
		status,
		action,
		created_at,
		updated_at
	FROM TaskAction
	WHERE task_action_id = $actionID;
	`

	result, err := r.tx.InTX().Execute(yql, table.ValueParam("$actionID", types.Int64Value(actionID)))
	if err != nil {
		return nil, nil, err
	}
	defer result.Close()

	var taskActionID internals.TaskActionID
	var taskID api.TaskID
	var status string
	var actionBytes []byte
	var createdAt time.Time
	var updatedAt time.Time

	if err = result.FetchExactlyOne(&taskActionID, &taskID, &status, &actionBytes, &createdAt, &updatedAt); err != nil {
		return nil, nil, err
	}

	var taskAction internals.TaskAction
	if err = json.Unmarshal(actionBytes, &taskAction); err != nil {
		return nil, nil, err
	}

	taskActionAdditionalInfo := &internals.TaskActionAdditionalInfo{
		CreatedAt: createdAt,
		Status:    internals.TaskActionStatus(status),
		TaskId:    taskID,
		UpdatedAt: updatedAt,
	}

	return &taskAction, taskActionAdditionalInfo, nil
}

func (r *appRepositoryImpl) GetTaskActionResultByID(actionID internals.TaskActionID) (*internals.TaskActionResult, *internals.TaskActionResultAdditionalInfo, error) {
	yql := `
	SELECT
		task_action_id,
		result,
		created_at
	FROM TaskActionResult
	WHERE task_action_id = $actionID;
	`

	result, err := r.tx.InTX().Execute(yql, table.ValueParam("$actionID", types.Int64Value(actionID)))
	if err != nil {
		return nil, nil, err
	}
	defer result.Close()

	var taskActionID internals.TaskActionID
	var resultBytes []byte
	var createdAt time.Time

	if err = result.FetchExactlyOne(&taskActionID, &resultBytes, &createdAt); err != nil {
		return nil, nil, err
	}

	var taskActionResult internals.TaskActionResult
	if err = json.Unmarshal(resultBytes, &taskActionResult); err != nil {
		return nil, nil, err
	}

	taskActionResultAdditionalInfo := &internals.TaskActionResultAdditionalInfo{
		CreatedAt: createdAt,
		TaskId:    api.TaskID(taskActionID),
	}

	return &taskActionResult, taskActionResultAdditionalInfo, nil
}

func (r *appRepositoryImpl) SetTaskActionStatus(actionID internals.TaskActionID, newStatus internals.TaskActionStatus) error {
	yql := `
	UPDATE TaskAction
	SET
		status = $newStatus,
		updated_at = CurrentUtcDatetime()
	WHERE task_action_id = $actionID;
	`

	parameters := []table.ParameterOption{
		table.ValueParam("$actionID", types.Int64Value(actionID)),
		table.ValueParam("$newStatus", types.TextValue(string(newStatus))),
	}

	result, err := r.tx.InTX().Execute(yql, parameters...)
	if err != nil {
		return err
	}
	defer result.Close()

	r.log.Debug("Updated task action status for action id ", actionID)

	return nil
}

func (r *appRepositoryImpl) GetTaskActionsByTaskID(taskID api.TaskID) ([]api.TaskActionWithResult, error) {
	yql := `
	SELECT
		ta.task_action_id as task_action_id,
		ta.task_id as task_id,
		ta.status as status,
		ta.action as action,
		ta.created_at as created_at,
		ta.updated_at as updated_at,
		tar.result as result,
		tar.created_at as result_created_at
	FROM TaskAction as ta
	LEFT JOIN TaskActionResult as tar ON ta.task_action_id = tar.task_action_id
	WHERE ta.task_id = $taskID
	ORDER BY ta.task_action_id;
	`

	result, err := r.tx.InTX().Execute(yql, table.ValueParam("$taskID", types.Int64Value(taskID)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	taskActions := make([]api.TaskActionWithResult, 0)

	for result.NextRow() {
		var taskActionID internals.TaskActionID
		var taskID api.TaskID
		var status string
		var actionBytes []byte
		var createdAt time.Time
		var updatedAt time.Time
		var resultBytes []byte
		var resultCreatedAt *time.Time

		if err := result.FetchRow(&taskActionID, &taskID, &status, &actionBytes, &createdAt, &updatedAt, &resultBytes, &resultCreatedAt); err != nil {
			return nil, err
		}

		taskActionRaw := map[string]any{}
		if err := json.Unmarshal(actionBytes, &taskActionRaw); err != nil {
			return nil, err
		}

		var taskActionResultRaw *api.RawJSON
		if len(resultBytes) > 0 {
			taskActionResultRaw = &api.RawJSON{}
			if err := json.Unmarshal(resultBytes, taskActionResultRaw); err != nil {
				return nil, err
			}
		}

		taskActionWithResult := api.TaskActionWithResult{
			TaskActionId:     taskActionID,
			CreatedAt:        createdAt,
			UpdatedAt:        &updatedAt,
			TaskAction:       taskActionRaw,
			TaskActionResult: taskActionResultRaw,
		}

		taskActions = append(taskActions, taskActionWithResult)
	}

	return taskActions, nil
}
