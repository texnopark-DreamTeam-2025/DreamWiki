package usecase

import (
	"errors"
	"fmt"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

func (u *appUsecaseImpl) FetchPageFromYWiki(pageURL string) error {
	slug := extractYWikiSlugFromURL(pageURL)

	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	// Write integration log
	err := repo.WriteIntegrationLogField("ywiki", fmt.Sprintf("Fetching page with slug: %s", slug))
	if err != nil {
		u.log.Errorf("Failed to write integration log: %v", err)
	}

	// 2. Go to YWiki client from deps, fetch page
	pageResponse, err := u.deps.YWikiClient.GetPage(u.ctx, slug)
	if err != nil {
		// TODO: If fetch return 404, delete page
		// Will be done later
		u.log.Errorf("Failed to fetch page from YWiki: %w", err)
		return err
	}

	// 3. Upsert page to repository
	pageFromYWIki := api.Page{
		Content: *pageResponse.Content,
		Title:   pageResponse.Title,
	}
	pageFromRepository, err := repo.GetPageBySlug(slug)
	if err != nil && !errors.Is(err, models.ErrNoRows) {
		u.log.Errorf("Failed to get page by slug: %w", err)
		return err
	}
	var pageID api.PageID

	if errors.Is(err, models.ErrNoRows) {
		newPageID, err := repo.CreatePage(slug, pageFromYWIki.Title, pageFromYWIki.Content)
		if err != nil {
			return err
		}
		pageID = *newPageID
	} else {
		pageID = pageFromRepository.PageId
		if pageFromRepository.Content != pageFromYWIki.Content {
			_, err := repo.AppendPageRevision(pageID, pageFromYWIki.Content)
			if err != nil {
				return err
			}
		}
		if pageFromRepository.Title != pageFromYWIki.Title {
			err := repo.SetPageTitle(pageID, pageFromYWIki.Title)
			if err != nil {
				return err
			}
		}
	}

	// 4. Indexate page using usecase function
	indexateReq := api.V1IndexatePageRequest{
		PageId: pageID,
	}
	_, err = u.indexatePageInTransaction(repo, indexateReq.PageId)
	if err != nil {
		u.log.Errorf("Failed to indexate page: %w", err)
		return err
	}

	// 5. Commit transaction
	return repo.Commit()
}

func (u *appUsecaseImpl) GetIntegrationLogs(integrationID api.IntegrationID, cursor *string) (fields []api.IntegrationLogField, nextInfo *api.NextInfo, err error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	return repo.GetIntegrationLogFields(integrationID, cursor, 50)
}

func (u *appUsecaseImpl) GithubAccountPRAsync(prURL string) (*api.TaskID, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	taskState := internals.TaskStateGitHubAccountPR{
		PrUrl:    prURL,
		TaskType: internals.GithubAccountPr,
	}

	var taskStateUnion internals.TaskState
	err := taskStateUnion.FromTaskStateGitHubAccountPR(taskState)
	if err != nil {
		return nil, err
	}

	taskID, err := repo.CreateTask(taskStateUnion)
	if err != nil {
		return nil, err
	}

	taskAction := internals.TaskAction{}
	taskAction.FromTaskActionNewTask(internals.TaskActionNewTask{TaskActionType: internals.NewTask})
	taskActionID, err := repo.CreateTaskAction(*taskID, taskAction)
	if err != nil {
		return nil, err
	}

	err = repo.EnqueueTaskAction(*taskActionID)
	if err != nil {
		return nil, err
	}

	err = repo.Commit()
	if err != nil {
		return nil, err
	}

	return taskID, nil
}

func (u *appUsecaseImpl) YwikiFetchAllAsync() (*api.TaskID, error) {
	panic("unimplemented")
}
