package staletaskfailer

import (
	"context"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/components/component"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/db_adapter"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type StaleTaskFailer struct {
	deps *deps.Deps
}

func NewStaleTaskFailer(deps *deps.Deps) *StaleTaskFailer {
	return &StaleTaskFailer{
		deps: deps,
	}
}

var _ component.Component = &StaleTaskFailer{}

func (s *StaleTaskFailer) Name() string {
	return "StaleTaskFailer"
}

func (s *StaleTaskFailer) Run(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := s.failStaleTasks(ctx); err != nil {
				s.deps.Logger.Error("failed to fail stale tasks", err)
			}
		}
	}
}

func (s *StaleTaskFailer) failStaleTasks(ctx context.Context) error {
	readOnlyTx := s.deps.YDBDriver.NewTransaction(ctx, db_adapter.SnapshotReadOnly)
	defer readOnlyTx.Rollback()

	repo := repository.NewAppRepository(ctx, &deps.RepositoryDeps{
		Deps: s.deps,
		TX:   readOnlyTx,
	})

	taskIDs, err := repo.GetStaleTaskIDs()
	if err != nil {
		return err
	}

	if len(taskIDs) == 0 {
		return nil
	}

	tx := s.deps.YDBDriver.NewTransaction(ctx, db_adapter.SerializableReadWrite)
	defer tx.Rollback()

	repo = repository.NewAppRepository(ctx, &deps.RepositoryDeps{
		Deps: s.deps,
		TX:   tx,
	})

	for _, taskID := range taskIDs {
		if err := repo.SetTaskStatus(taskID, api.FailedByTimeout); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
