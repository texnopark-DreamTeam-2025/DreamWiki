package topic_reader

import (
	"context"
	"fmt"
	"strconv"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/db_adapter"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_actions_usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_factory"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicreader"
)

const (
	taskActionsTopicName       = "TaskActionToExecute"
	taskActionResultsTopicName = "TaskActionResultReady"
)

type (
	TopicReaders struct {
		ctx                          context.Context
		cancel                       context.CancelFunc
		TaskActionsTopicReader       db_adapter.TopicReader
		TaskActionResultsTopicReader db_adapter.TopicReader
		log                          logger.Logger
		deps                         *deps.Deps
		onReaderClose                chan struct{}
	}
)

func NewTopicReader(ctx context.Context, d *deps.Deps) (*TopicReaders, error) {
	ctx, cancel := context.WithCancel(ctx)
	t := TopicReaders{
		ctx:           ctx,
		cancel:        cancel,
		log:           d.Logger,
		deps:          d,
		onReaderClose: make(chan struct{}, 2),
	}

	processTaskActionMessage := func(message []byte) {
		taskActionID, _ := strconv.ParseInt(string(message), 10, 64)
		t.log.Info("processing task action", "action_id", taskActionID)

		taskActionUsecase := task_actions_usecase.NewTaskActionUsecase(context.Background(), t.deps)
		err := taskActionUsecase.ExecuteAction(internals.TaskActionID(taskActionID))
		if err != nil {
			t.log.Error("failed to execute task action", " action_id ", taskActionID, " error ", err)
		} else {
			t.log.Info("successfully executed task action", "action_id", taskActionID)
		}
	}

	processTaskActionResultMessage := func(message []byte) {
		taskActionID, _ := strconv.ParseInt(string(message), 10, 64)
		t.log.Info("processing task action result", "action_id", taskActionID)

		repoDeps := &deps.RepositoryDeps{
			TX:   d.YDBDriver.NewTransaction(ctx, db_adapter.SerializableReadWrite),
			Deps: d,
		}
		repo := repository.NewAppRepository(context.Background(), repoDeps)
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

	taskActionsTopicReader := d.YDBDriver.NewTopicReader(taskActionsTopicName, processTaskActionMessage)
	taskActionResultsTopicReader := d.YDBDriver.NewTopicReader(taskActionResultsTopicName, processTaskActionResultMessage)

	t.TaskActionsTopicReader = taskActionsTopicReader
	t.TaskActionResultsTopicReader = taskActionResultsTopicReader

	return &t, nil
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

	t.TaskActionsTopicReader.ReadMessages(context.Background())
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

	t.TaskActionResultsTopicReader.ReadMessages(context.Background())
}
