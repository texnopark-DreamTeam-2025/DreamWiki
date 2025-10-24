package topic_reader

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicoptions"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicreader"
)

type (
	TopicReaders struct {
		TaskActionsTopicReader       *topicreader.Reader
		TaskActionResultsTopicReader *topicreader.Reader
		log                          logger.Logger
	}
)

func NewTopicReader(db *ydb.Driver, log logger.Logger) (*TopicReaders, error) {
	taskActionsTopicReader, err := db.Topic().StartReader("dream_wiki", topicoptions.ReadTopic("TaskActionToExecute"))
	if err != nil {
		return nil, err
	}

	taskActionResultsTopicReader, err := db.Topic().StartReader("dream_wiki", topicoptions.ReadTopic("my-topic"))
	if err != nil {
		taskActionsTopicReader.Close(context.Background())
		return nil, err
	}

	return &TopicReaders{
		TaskActionsTopicReader:       taskActionsTopicReader,
		TaskActionResultsTopicReader: taskActionResultsTopicReader,
		log:                          log,
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

func (t *TopicReaders) readTaskActionMessages() {
	t.log.Info("START READING TaskActionsTopicReader")
	for {
		mess, err := t.TaskActionsTopicReader.ReadMessage(context.Background())
		if err != nil {
			t.log.Error(err)
			break
		}
		t.log.Info("SKDJFDLSJFKLDFJLSJFSDLFJFLKSJFLKDJSFKLSJFKSDJFKLSDJFLKSJFLKSJFLSKFJKLJDFLKSDJFLKSJFLKSDJFSKDFJSKJFHJK")
		t.log.Info(mess)
		t.TaskActionsTopicReader.Commit(mess.Context(), mess)
	}
}

func (t *TopicReaders) readTaskActionResultMessages() {

}
