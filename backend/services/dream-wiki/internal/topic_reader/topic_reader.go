package topic_reader

import (
	"context"
	"strconv"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_actions_usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_factory"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicoptions"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicreader"
)

type (
	TopicReaders struct {
		ctx                          context.Context
		TaskActionsTopicReader       *topicreader.Reader
		TaskActionResultsTopicReader *topicreader.Reader
		log                          logger.Logger
		deps                         *deps.Deps
	}
)

func NewTopicReader(deps *deps.Deps) (*TopicReaders, error) {
	topic := deps.YDBDriver.Topic()
	taskActionsTopicReader, err := topic.StartReader("dream_wiki", topicoptions.ReadTopic("TaskActionToExecute"))
	if err != nil {
		return nil, err
	}

	taskActionResultsTopicReader, err := topic.StartReader("dream_wiki", topicoptions.ReadTopic("TaskActionResultReady"))
	if err != nil {
		taskActionsTopicReader.Close(context.Background())
		return nil, err
	}

	return &TopicReaders{
		TaskActionsTopicReader:       taskActionsTopicReader,
		TaskActionResultsTopicReader: taskActionResultsTopicReader,
		log:                          deps.Logger,
		deps:                         deps,
		ctx:                          context.Background(), // TODO: use global application context and implement graceful shutdown
	}, nil
}

func (t *TopicReaders) Close(ctx context.Context) error {
	err := t.TaskActionsTopicReader.Close(ctx)
	if err != nil {
		return err
	}

	return t.TaskActionResultsTopicReader.Close(ctx)
}

func (t *TopicReaders) ReadMessages() {
	go t.readTaskActionMessages()
	go t.readTaskActionResultMessages()
}

func readTopic(ctx context.Context, reader *topicreader.Reader, log logger.Logger, onMessage func(int64)) {
	for {
		mess, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Error("failed to read message", err)
			break
		}

		// We commit here. Reason is not all actions can be idempotent.
		// There is better to avoid multiple (maybe infinity looped) requests
		// than get minor reliability improvement
		commitError := reader.Commit(mess.Context(), mess)
		if commitError != nil {
			continue
		}

		messageContent := make([]byte, 1024)
		n, readError := mess.Read(messageContent)
		if readError != nil {
			log.Error("failed to get message data", "error", err)
			continue
		}
		if n == 0 {
			log.Error("zero message length")
			continue
		}
		messageContent = messageContent[:n]

		actionID, err := strconv.ParseInt(string(messageContent), 10, 64)
		if err != nil {
			log.Error("failed to parse action ID from message", " data ", string(messageContent), " error ", err)
			continue
		}

		log.Info("processing message with task action ID: ", actionID)
		onMessage(actionID)
	}
}

func (t *TopicReaders) readTaskActionMessages() {
	t.log.Info("START READING TaskActionsTopicReader")
	processTaskActionMessage := func(taskActionID int64) {
		t.log.Info("processing task action", "action_id", taskActionID)

		taskActionUsecase := task_actions_usecase.NewTaskActionUsecase(context.Background(), t.deps)
		err := taskActionUsecase.ExecuteAction(internals.TaskActionID(taskActionID))
		if err != nil {
			t.log.Error("failed to execute task action", "action_id", taskActionID, "error", err)
		} else {
			t.log.Info("successfully executed task action", "action_id", taskActionID)
		}

		if err != nil {
			t.log.Error("failed to commit task action message", "action_id", taskActionID, "error", err)
		}
	}

	readTopic(t.ctx, t.TaskActionsTopicReader, t.log, processTaskActionMessage)
}

func (t *TopicReaders) readTaskActionResultMessages() {
	t.log.Info("START READING TaskActionResultsTopicReader")
	processTaskActionResultMessage := func(taskActionID int64) {
		t.log.Info("processing task action result", "action_id", taskActionID)

		repo := repository.NewAppRepository(context.Background(), t.deps)
		defer repo.Rollback()

		_, taskActionAdditionalInfo, err := repo.GetTaskActionByID(internals.TaskActionID(taskActionID))
		if err != nil {
			t.log.Error("failed to get task action by ID", "action_id", taskActionID, "error", err)
			return
		}

		taskDigest, taskState, err := repo.GetTaskByID(taskActionAdditionalInfo.TaskId)
		if err != nil {
			t.log.Error("failed to get task by ID", "task_id", taskActionAdditionalInfo.TaskId, "error", err)
			return
		}

		taskLogicCreator := task_factory.CreateTaskLogicCreator(taskState)
		task := task_common.NewTask(*taskDigest, taskState, taskLogicCreator)

		taskActionResult, _, err := repo.GetTaskActionResultByID(internals.TaskActionID(taskActionID))
		if err != nil {
			t.log.Error("failed to get task action result by ID", "action_id", taskActionID, "error", err)
			return
		}

		err = task.OnActionResult(*taskActionResult)
		if err != nil {
			t.log.Error("failed to process task action result", "action_id", taskActionID, "error", err)
			return
		}

		err = repo.Commit()
		if err != nil {
			t.log.Error("failed to commit transaction", "error", err)
			return
		}

		t.log.Info("successfully processed task action result", "action_id", taskActionID)
	}

	readTopic(t.ctx, t.TaskActionResultsTopicReader, t.log, processTaskActionResultMessage)
}
