package app

import (
	"github.com/google/uuid"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type AppRepository interface {
	Commit()
	Rollback()

	SearchByEmbedding(query string, queryEmbedding models.Embedding) ([]api.SearchResultItem, error)
	RetrievePageByID(pageID uuid.UUID) (*api.Page, error)
	RemovePageIndexation(pageID uuid.UUID) error
	AddIndexedParagraph(paragraph models.ParagraphWithEmbedding) error
	DeleteAllPages() error
	GetUserByLogin(login string) (*models.User, error)
	WriteIntegrationLogField(integrationID api.IntegrationID, logText string) error
	GetPageBySlug(yWikiSlug string) (*api.Page, error)
	UpsertPage(page api.Page, ywikiSlug string) error
	DeletePageBySlug(yWikiSlug string) error
}

type AppUsecase interface {
	Search(req api.V1SearchRequest) (*api.V1SearchResponse, error)
	GetDiagnosticInfo(req api.V1DiagnosticInfoGetRequest) (*api.V1DiagnosticInfoGetResponse, error)
	IndexatePage(req api.V1IndexatePageRequest) (*api.V1IndexatePageResponse, error)
	Login(req api.V1LoginRequest) (*api.V1LoginResponse, error)
	FetchPageFromYWiki(pageURL string) error
	// AccountGitHubPullRequest(pullRequestURL string)error
}
