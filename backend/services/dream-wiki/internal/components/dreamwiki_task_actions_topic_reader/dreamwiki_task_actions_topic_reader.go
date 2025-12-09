package dreamwikitaskactionstopicreader

import (
	"context"
	"errors"
	"strconv"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/components/component"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/db_adapter"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_actions_usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

const taskActionsTopicName = "TaskActionToExecute"

type DreamWikiTaskActionsTopicReader struct {
	deps   *deps.Deps
	reader db_adapter.TopicReader
}

func NewDreamWikiTaskActionsTopicReader(deps *deps.Deps) *DreamWikiTaskActionsTopicReader {
	processTaskActionMessage := func(message []byte) {
		taskActionID, _ := strconv.ParseInt(string(message), 10, 64)
		deps.Logger.Info("processing task action", "action_id", taskActionID)

		taskActionUsecase := task_actions_usecase.NewTaskActionUsecase(context.Background(), deps)
		err := taskActionUsecase.ExecuteAction(internals.TaskActionID(taskActionID))
		if err != nil {
			deps.Logger.Error("failed to execute task action", " action_id ", taskActionID, " error ", err)
		} else {
			deps.Logger.Info("successfully executed task action", "action_id", taskActionID)
		}
	}

	reader := deps.YDBDriver.NewTopicReader(taskActionsTopicName, processTaskActionMessage)

	return &DreamWikiTaskActionsTopicReader{
		deps:   deps,
		reader: reader,
	}
}

var _ component.Component = &DreamWikiTaskActionsTopicReader{}

func (d *DreamWikiTaskActionsTopicReader) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- d.reader.ReadMessages(ctx)
	}()

	err := <-errCh
	if err != nil && !errors.Is(err, context.Canceled) {
		d.deps.Logger.Error("task actions topic reader error", "error", err)
		return err
	}
	return nil
}

func (d *DreamWikiTaskActionsTopicReader) Name() string {
	return "DreamWikiTaskActionsTopicReader"
}
