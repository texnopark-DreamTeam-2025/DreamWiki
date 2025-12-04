package topic_reader

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_actions_usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_factory"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicoptions"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicreader"
)

type (
	TopicReaders struct {
		ctx                          context.Context
		cancel                       context.CancelFunc
		TaskActionsTopicReader       *topicreader.Reader
		TaskActionResultsTopicReader *topicreader.Reader
		log                          logger.Logger
		deps                         *deps.Deps
		onReaderClose                chan struct{}
	}
)

func NewTopicReader(ctx context.Context, deps *deps.Deps) (*TopicReaders, error) {
	ctx, cancel := context.WithCancel(ctx)

	topic1 := deps.YDBDriver.Topic()
	taskActionsTopicReader, err := topic1.StartReader(
		"dream_wiki",
		topicoptions.ReadTopic("TaskActionToExecute"),
		topicoptions.WithReaderOperationTimeout(time.Second),
	)
	if err != nil {
		cancel()
		return nil, err
	}

	topic2 := deps.YDBDriver.Topic()
	taskActionResultsTopicReader, err := topic2.StartReader(
		"dream_wiki",
		topicoptions.ReadTopic("TaskActionResultReady"),
		topicoptions.WithReaderOperationTimeout(time.Second),
	)
	if err != nil {
		taskActionsTopicReader.Close(ctx)
		cancel()
		return nil, err
	}

	return &TopicReaders{
		ctx:                          ctx,
		cancel:                       cancel,
		TaskActionsTopicReader:       taskActionsTopicReader,
		TaskActionResultsTopicReader: taskActionResultsTopicReader,
		log:                          deps.Logger,
		deps:                         deps,
		onReaderClose:                make(chan struct{}, 2),
	}, nil
}

func (t *TopicReaders) Close(ctx context.Context) error {
	t.cancel()

	for range 2 {
		select {
		case <-t.onReaderClose:
			continue
		case <-ctx.Done():
			return fmt.Errorf("cant close readers, by timeout")
		}
	}

	return nil
}

func (t *TopicReaders) ReadMessages() {
	go t.readTaskActionMessages()
	go t.readTaskActionResultMessages()
}

func readTopic(ctx context.Context, reader *topicreader.Reader, onReaderClose *chan struct{}, log logger.Logger, onMessage func(int64)) {
	closeReader := func() {
		log.Info("context cancelled, stopping topic reader")
		reader.Close(context.Background())
		*onReaderClose <- struct{}{}
	}
	for {
		select {
		case <-ctx.Done():
			closeReader()
			return
		default:
			log.Info("Ready to get message")
			mess, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					closeReader()
					return
				}
				log.Error("failed to read message", err)
				continue
			}
			log.Info("Got message")

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
				log.Error("failed to get message data", "error", readError)
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
}

func (t *TopicReaders) readTaskActionMessages() {
	t.log.Info("START READING TaskActionsTopicReader")
	processTaskActionMessage := func(taskActionID int64) {
		t.log.Info("processing task action", "action_id", taskActionID)

		taskActionUsecase := task_actions_usecase.NewTaskActionUsecase(context.Background(), t.deps)
		err := taskActionUsecase.ExecuteAction(internals.TaskActionID(taskActionID))
		if err != nil {
			t.log.Error("failed to execute task action", " action_id ", taskActionID, " error ", err)
		} else {
			t.log.Info("successfully executed task action", "action_id", taskActionID)
		}
	}

	readTopic(t.ctx, t.TaskActionsTopicReader, &t.onReaderClose, t.log, processTaskActionMessage)
}

func (t *TopicReaders) failTaskAndCommit(repo repository.AppRepository, taskID api.TaskID, taskActionID int64) {
	setTaskErr := repo.SetTaskStatus(taskID, api.FailedByError)
	if setTaskErr != nil {
		t.log.Error("failed to set task status to failed", "task_id", taskID, "error", setTaskErr)
	} else {
		commitErr := repo.Commit()
		if commitErr != nil {
			t.log.Error("failed to commit transaction when setting task status to failed", "error", commitErr)
		} else {
			t.log.Info("successfully set task status to failed", "task_id", taskID)
		}
	}
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

		taskLogicCreator := task_factory.CreateTaskLogicCreator()
		task := task_common.NewTask(
			context.Background(),
			&task_common.TaskDeps{
				Deps:   t.deps,
				Digest: *taskDigest,
				State:  taskState,
				Repo:   repo,
			},
			taskLogicCreator,
		)

		taskActionResult, _, err := repo.GetTaskActionResultByID(internals.TaskActionID(taskActionID))
		if err != nil {
			t.log.Error("failed to get task action result by ID", "action_id", taskActionID, "error", err)
			return
		}

		err = task.OnActionResult(*taskActionResult)
		if err != nil {
			t.log.Error("failed to process task action result", "action_id", taskActionID, "error", err)
			t.failTaskAndCommit(repo, taskActionAdditionalInfo.TaskId, taskActionID)
			return
		}

		t.log.Info("successfully processed task action result", "action_id", taskActionID)
	}

	readTopic(t.ctx, t.TaskActionResultsTopicReader, &t.onReaderClose, t.log, processTaskActionResultMessage)
}
