package repository

import (
	"fmt"
	"strings"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicwriter"
)

func (r *appRepositoryImpl) CreateTaskAction(taskID api.TaskID, actionState internals.TaskAction) (*internals.TaskActionID, error) {
	panic("unimplemented")
}

func (r *appRepositoryImpl) CreateTaskActionResult(actionID internals.TaskActionID, result internals.TaskActionResult) error {
	panic("unimplemented")
}

func (r *appRepositoryImpl) EnqueueTaskAction(actionID internals.TaskActionID) error {
	topicClient := r.ydbClient.TopicClient()
	writer, err := topicClient.StartTransactionalWriter(r.ydbClient.GetTX(), "TaskActionToExecute")
	if err != nil {
		return err
	}
	err = writer.Write(r.ctx, topicwriter.Message{Data: strings.NewReader(fmt.Sprintf("adsfadsasf %d", actionID))})
	if err != nil {
		return err
	}
	r.log.Info("Enqueued message")
	return nil
}

func (r *appRepositoryImpl) EnqueueTaskActionResult(actionID internals.TaskActionID) error {
	panic("unimplemented")
}

func (r *appRepositoryImpl) GetTaskActionByID(actionID internals.TaskActionID) (*internals.TaskAction, *internals.TaskActionAdditionalInfo, error) {
	panic("unimplemented")
}

func (r *appRepositoryImpl) GetTaskActionResultByID(actionID internals.TaskActionID) (*internals.TaskActionResult, *internals.TaskActionResultAdditionalInfo, error) {
	panic("unimplemented")
}

func (r *appRepositoryImpl) SetTaskActionStatus(actionID internals.TaskActionID, newStatus internals.TaskActionStatus) error {
	panic("unimplemented")
}
