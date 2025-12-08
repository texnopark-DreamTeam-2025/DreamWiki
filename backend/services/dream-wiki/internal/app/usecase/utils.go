package usecase

import (
	"strings"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/db_adapter"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
)

func extractYWikiSlugFromURL(pageURL string) string {
	const prefix = "https://wiki.yandex.ru/"

	return strings.TrimSuffix(strings.TrimPrefix(pageURL, prefix), "/")
}

func (u *appUsecaseImpl) createReadOnlyRepository() repository.AppRepository {
	return repository.NewAppRepository(u.ctx, &deps.RepositoryDeps{
		TX:   u.deps.YDBDriver.NewTransaction(u.ctx, db_adapter.SnapshotReadOnly),
		Deps: u.deps,
	})
}

func (u *appUsecaseImpl) createReadWriteRepository() repository.AppRepository {
	return repository.NewAppRepository(u.ctx, &deps.RepositoryDeps{
		TX:   u.deps.YDBDriver.NewTransaction(u.ctx, db_adapter.SerializableReadWrite),
		Deps: u.deps,
	})
}
