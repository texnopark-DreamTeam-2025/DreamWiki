package dreamwikitaskactionresultstopicreader

import (
	"context"
	"strconv"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/components/component"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/db_adapter"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_factory"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

const taskActionResultsTopicName = "TaskActionResultReady"

type DreamWikiTaskActionResultsTopicReader struct {
	deps   *deps.Deps
	reader db_adapter.TopicReader
}

func NewDreamWikiTaskActionResultsTopicReader(d *deps.Deps) *DreamWikiTaskActionResultsTopicReader {
	processTaskActionResultMessage := func(message []byte) {
		taskActionID, _ := strconv.ParseInt(string(message), 10, 64)
		d.Logger.Info("processing task action result", "action_id", taskActionID)

		repoDeps := &deps.RepositoryDeps{
			TX:   d.YDBDriver.NewTransaction(context.Background(), db_adapter.SerializableReadWrite),
			Deps: d,
		}
		repo := repository.NewAppRepository(context.Background(), repoDeps)
		defer repo.Rollback()

		_, taskActionAdditionalInfo, err := repo.GetTaskActionByID(internals.TaskActionID(taskActionID))
		if err != nil {
			d.Logger.Error("failed to get task action by ID. action_id=", taskActionID, " error ", err)
			return
		}

		taskDigest, taskState, err := repo.GetTaskByID(taskActionAdditionalInfo.TaskId)
		if err != nil {
			d.Logger.Error("failed to get task by ID. task_id=", taskActionAdditionalInfo.TaskId, " error ", err)
			return
		}

		taskLogicCreator := task_factory.CreateTaskLogicCreator()
		task := task_common.NewTask(
			context.Background(),
			&task_common.TaskDeps{
				Deps:   d,
				Digest: *taskDigest,
				State:  taskState,
				Repo:   repo,
			},
			taskLogicCreator,
		)

		taskActionResult, _, err := repo.GetTaskActionResultByID(internals.TaskActionID(taskActionID))
		if err != nil {
			d.Logger.Error("failed to get task action result by ID", "action_id", taskActionID, "error", err)
			return
		}

		err = task.OnActionResult(*taskActionResult)
		if err != nil {
			d.Logger.Error("failed to process task action result", "action_id", taskActionID, "error", err)
			return
		}

		d.Logger.Info("successfully processed task action result", "action_id", taskActionID)
	}

	reader := d.YDBDriver.NewTopicReader(taskActionResultsTopicName, processTaskActionResultMessage)

	return &DreamWikiTaskActionResultsTopicReader{
		deps:   d,
		reader: reader,
	}
}

var _ component.Component = &DreamWikiTaskActionResultsTopicReader{}

func (d *DreamWikiTaskActionResultsTopicReader) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- d.reader.ReadMessages(ctx)
	}()

	select {
	case <-ctx.Done():
		d.deps.Logger.Info("task action results topic reader is shutting down")
		return nil
	case err := <-errCh:
		if err != nil {
			d.deps.Logger.Error("task action results topic reader error", "error", err)
			// TODO in Q5: fail task and commit
			return err
		}
		return nil
	}
}

func (d *DreamWikiTaskActionResultsTopicReader) Name() string {
	return "DreamWikiTaskActionResultsTopicReader"
}
