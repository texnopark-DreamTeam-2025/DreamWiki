package usecase

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type (
	AppUsecase interface {
		// domain_auth.go
		Login(req api.V1LoginRequest) (*api.V1LoginResponse, error)

		// domain_drafts.go
		CreateDraft(pageURL string) (api.DraftDigest, error)
		DeleteDraft(draftID api.DraftID) error
		GetDraft(draftID api.DraftID) (api.Draft, error)
		UpdateDraft(draftID api.DraftID, newContent *string, newTitle *string) error
		ApplyDraft(draftID api.DraftID) error
		ListDrafts(cursor *string) ([]api.DraftDigest, error)

		// domain_integrations.go
		FetchPageFromYWiki(pageURL string) error
		YwikiFetchAllAsync() (*api.TaskID, error)
		GetIntegrationLogs(integrationID api.IntegrationID, cursor *string) (fields []api.IntegrationLogField, newCursor string, err error)
		GithubAccountPRAsync(prURL string) (*api.TaskID, error)

		// domain_page_indexation.go
		IndexatePage(pageID api.PageID) (*api.V1IndexatePageResponse, error)

		// domain_pages.go
		GetDiagnosticInfo(req api.V1DiagnosticInfoGetRequest) (*api.V1DiagnosticInfoGetResponse, error)
		GetPagesTree(activePagesIDs []api.PageID) ([]api.TreeItem, error)

		// domain_search.go
		Search(req api.V1SearchRequest) (*api.V1SearchResponse, error)

		// domain_tasks.go
		CancelTask(taskID api.TaskID) error
		GetTaskDetails(taskID api.TaskID) (api.Task, error)
		ListTasks(cursor *string) (tasks []api.TaskDigest, newCursor string, err error)
		RetryTask(taskID api.TaskID) error
		GetTaskInternalState(taskID api.TaskID) (*api.V1TasksInternalStateGetResponse, error)
		RecreateTask(taskID api.TaskID) (*api.TaskID, error)
	}

	appUsecaseImpl struct {
		ctx  context.Context
		deps *deps.Deps
		log  logger.Logger
	}
)

func NewAppUsecaseImpl(ctx context.Context, deps *deps.Deps) AppUsecase {
	return &appUsecaseImpl{ctx: ctx, deps: deps, log: deps.Logger}
}
