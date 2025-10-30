package repository

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func decodeTasksCursor(cursor *api.Cursor) int64 {
	if cursor == nil {
		return 0
	}

	idUpperLimit, err := strconv.ParseInt(string(*cursor), 10, 64)
	if err != nil {
		return 0
	}

	return idUpperLimit
}

func encodeTasksCursor(idUpperLimit int64) api.Cursor {
	return api.Cursor(strconv.FormatInt(idUpperLimit, 10))
}

func (r *appRepositoryImpl) CreateTask(taskState internals.TaskState) (*api.TaskID, error) {
	yql := `
	INSERT INTO Task (status, state, created_at, updated_at)
	VALUES ('executing', $state, CurrentUtcDatetime(), CurrentUtcDatetime())
	RETURNING task_id;
	`

	stateBytes, err := json.Marshal(taskState)
	if err != nil {
		return nil, err
	}

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$state", types.JSONValueFromBytes(stateBytes)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var taskID api.TaskID
	err = result.FetchExactlyOne(&taskID)
	if err != nil {
		return nil, err
	}

	return &taskID, nil
}

func (r *appRepositoryImpl) GetTaskByID(taskID api.TaskID) (*api.TaskDigest, *internals.TaskState, error) {
	yql := `
	SELECT
		task_id,
		status,
		state,
		created_at,
		updated_at
	FROM Task
	WHERE task_id = $taskID;
	`

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$taskID", types.Int64Value(taskID)))
	if err != nil {
		return nil, nil, err
	}
	defer result.Close()

	var retrievedTaskID api.TaskID
	var status string
	var stateBytes []byte
	var createdAt time.Time
	var updatedAt time.Time

	err = result.FetchExactlyOne(&retrievedTaskID, &status, &stateBytes, &createdAt, &updatedAt)
	if err != nil {
		return nil, nil, err
	}

	var taskState internals.TaskState
	if err = json.Unmarshal(stateBytes, &taskState); err != nil {
		return nil, nil, err
	}

	// For now, we'll use a placeholder for triggered_by and description
	// In a real implementation, these would be derived from the task state
	taskDigest := &api.TaskDigest{
		TaskId:      retrievedTaskID,
		Status:      api.TaskStatus(status),
		TriggeredBy: "system",           // Placeholder
		Description: "Task description", // Placeholder
	}

	return taskDigest, &taskState, nil
}

func (r *appRepositoryImpl) ListTasks(cursor *api.Cursor, limit int64) ([]api.TaskDigest, []internals.TaskState, *api.Cursor, error) {
	yql := `
	SELECT
		task_id,
		status,
		state,
		created_at,
		updated_at
	FROM Task
	WHERE task_id < $idUpperLimit
	ORDER BY task_id
	LIMIT $limit;
	`

	idUpperLimit := decodeTasksCursor(cursor)

	result, err := r.ydbClient.InTX().Execute(yql,
		table.ValueParam("$idUpperLimit", types.Int64Value(idUpperLimit)),
		table.ValueParam("$limit", types.Uint64Value(uint64(limit))),
	)
	if err != nil {
		return nil, nil, nil, err
	}
	defer result.Close()

	taskDigests := make([]api.TaskDigest, 0, result.RowCount())
	taskStates := make([]internals.TaskState, 0, result.RowCount())

	if result.RowCount() == 0 {
		if cursor == nil {
			return taskDigests, taskStates, nil, nil
		}
		return taskDigests, taskStates, cursor, nil
	}

	newIDFrom := int64(0)
	for result.NextRow() {
		var taskID api.TaskID
		var status string
		var stateBytes []byte
		var createdAt time.Time
		var updatedAt time.Time

		err := result.FetchRow(&taskID, &status, &stateBytes, &createdAt, &updatedAt)
		if err != nil {
			return nil, nil, nil, err
		}

		var taskState internals.TaskState
		if err = json.Unmarshal(stateBytes, &taskState); err != nil {
			return nil, nil, nil, err
		}

		// For now, we'll use a placeholder for triggered_by and description
		// In a real implementation, these would be derived from the task state
		taskDigest := api.TaskDigest{
			TaskId:      taskID,
			Status:      api.TaskStatus(status),
			TriggeredBy: "system",           // Placeholder
			Description: "Task description", // Placeholder
		}

		taskDigests = append(taskDigests, taskDigest)
		taskStates = append(taskStates, taskState)
		newIDFrom = int64(taskID)
	}

	newCursor := encodeTasksCursor(newIDFrom)
	return taskDigests, taskStates, &newCursor, nil
}

func (r *appRepositoryImpl) SetTaskState(taskID api.TaskID, newState internals.TaskState) error {
	yql := `
	UPDATE Task
	SET
		state = $newState,
		updated_at = CurrentUtcDatetime()
	WHERE task_id = $taskID;
	`

	stateBytes, err := json.Marshal(newState)
	if err != nil {
		return err
	}

	parameters := []table.ParameterOption{
		table.ValueParam("$taskID", types.Int64Value(taskID)),
		table.ValueParam("$newState", types.JSONValueFromBytes(stateBytes)),
	}

	result, err := r.ydbClient.InTX().Execute(yql, parameters...)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (r *appRepositoryImpl) SetTaskStatus(taskID api.TaskID, newStatus api.TaskStatus) error {
	yql := `
	UPDATE Task
	SET
		status = $newStatus,
		updated_at = CurrentUtcDatetime()
	WHERE task_id = $taskID;
	`

	parameters := []table.ParameterOption{
		table.ValueParam("$taskID", types.Int64Value(taskID)),
		table.ValueParam("$newStatus", types.TextValue(string(newStatus))),
	}

	result, err := r.ydbClient.InTX().Execute(yql, parameters...)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}
