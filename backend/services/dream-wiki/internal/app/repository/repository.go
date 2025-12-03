package repository

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/ydb_wrapper"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

type (
	AppRepository interface {
		Commit() error
		Rollback()

		// domain_drafts.go
		CreateDraft(pageID api.PageID, draftTitle string, draftContent string) (*api.DraftID, error)
		GetDraftByID(draftID api.DraftID) (*api.Draft, error)
		ListDrafts(cursor *api.Cursor, limit int64) ([]api.DraftDigest, *api.NextInfo, error)
		RemoveDraft(draftID api.DraftID) error
		SetDraftStatus(draftID api.DraftID, newStatus api.DraftStatus) error
		SetDraftContent(draftID api.DraftID, newContent string) error
		SetDraftTitle(draftID api.DraftID, newTitle string) error
		SetDraftBaseRevision(draftID api.DraftID, newRevisionID internals.RevisionID) error

		// domain_integration_logs.go
		WriteIntegrationLogField(integrationID api.IntegrationID, logText string) error
		GetIntegrationLogFields(integrationID api.IntegrationID, cursor *api.Cursor, limit uint64) ([]api.IntegrationLogField, *api.NextInfo, error)

		// domain_page_indexation.go
		RemovePageIndexation(pageID api.PageID) error
		AddIndexedParagraph(paragraph internals.ParagraphWithEmbedding) error
		AddTerm(term string, pageID api.PageID, paragraphIndex int64, timesIn int64) error

		// domain_pages.go
		GetPageBySlug(yWikiSlug string) (*api.Page, error)
		CreatePage(yWikiSlug string, title string, content string) (*api.PageID, error)
		AppendPageRevision(pageID api.PageID, newContent string) (*internals.RevisionID, error)
		DeletePageBySlug(yWikiSlug string) error
		GetAllPageDigests() ([]api.PageDigest, error)
		GetPageByID(pageID api.PageID) (*api.Page, *internals.PageAdditionalInfo, error)
		SetPageTitle(pageID api.PageID, newTitle string) error
		DeleteAllPages() error

		// domain_search.go
		SearchByEmbedding(query string, queryEmbedding internals.Embedding, limit int) ([]internals.SearchResultItem, error)
		SearchByEmbeddingWithContext(query string, queryEmbedding internals.Embedding, contextSize int) ([]internals.ParagraphWithContext, error)
		SearchByTerms(terms []string, limit int) ([]internals.SearchResultItem, error)

		// domain_tasks.go
		GetTaskByID(taskID api.TaskID) (*api.TaskDigest, *internals.TaskState, error)
		ListTasks(cursor *api.Cursor, limit int64) ([]api.TaskDigest, []internals.TaskState, *api.NextInfo, error)
		CreateTask(taskState internals.TaskState) (*api.TaskID, error)
		SetTaskStatus(taskID api.TaskID, newStatus api.TaskStatus) error
		SetTaskState(taskID api.TaskID, newState internals.TaskState) error

		// domain_task_actions.go
		CreateTaskAction(taskID api.TaskID, actionState internals.TaskAction) (*internals.TaskActionID, error)
		GetTaskActionByID(actionID internals.TaskActionID) (*internals.TaskAction, *internals.TaskActionAdditionalInfo, error)
		EnqueueTaskAction(actionID internals.TaskActionID) error
		SetTaskActionStatus(actionID internals.TaskActionID, newStatus internals.TaskActionStatus) error
		CreateTaskActionResult(actionID internals.TaskActionID, result internals.TaskActionResult) error
		GetTaskActionResultByID(actionID internals.TaskActionID) (*internals.TaskActionResult, *internals.TaskActionResultAdditionalInfo, error)
		EnqueueTaskActionResult(actionID internals.TaskActionID) error

		// domain_users.go
		GetUserByLogin(username string) (*models.User, error)
	}

	appRepositoryImpl struct {
		ctx       context.Context
		ydbClient ydb_wrapper.YDBWrapper
		log       logger.Logger
	}
)

var (
	_ AppRepository = &appRepositoryImpl{}
)

func NewAppRepository(ctx context.Context, deps *deps.Deps) AppRepository {
	return &appRepositoryImpl{
		ctx:       ctx,
		log:       deps.Logger,
		ydbClient: ydb_wrapper.NewYDBWrapper(ctx, deps, true),
	}
}

func (r *appRepositoryImpl) Commit() error {
	return r.ydbClient.Commit()
}

func (r *appRepositoryImpl) Rollback() {
	r.ydbClient.Rollback()
}
