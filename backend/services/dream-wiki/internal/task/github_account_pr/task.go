package github_account_pr

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/ycloud_client_gen"
)

type (
	gitHubAccountPRTask struct {
		taskID api.TaskID
		status api.TaskStatus
		state  internals.TaskStateGitHubAccountPR
		ctx    context.Context
		deps   *deps.Deps
		repo   repository.AppRepository
	}
)

var (
	_ task_common.TaskLogic = (*gitHubAccountPRTask)(nil)
)

func NewGitHubAccountPRTask(ctx context.Context, state internals.TaskStateGitHubAccountPR, deps *task_common.TaskDeps) *gitHubAccountPRTask {
	return &gitHubAccountPRTask{
		state:  state,
		status: deps.Digest.Status,
		ctx:    ctx,
		deps:   deps.Deps,
		taskID: deps.Digest.TaskId,
		repo:   deps.Repo,
	}
}

func (t *gitHubAccountPRTask) Close() {
	t.repo.Rollback()
}

func (t *gitHubAccountPRTask) CalculateSubtasks() ([]api.Subtask, error) {
	subtasks := []api.Subtask{}

	subtask1 := api.Subtask{
		Description: "Fetch PR data from GitHub",
		Status:      t.getSubtaskStatus(t.state.PrPatch != nil && t.state.PrDescription != nil, api.Done),
		Subsubtasks: []api.SubSubtask{},
	}
	subtasks = append(subtasks, subtask1)

	subtask2 := api.Subtask{
		Description: "Detect product changes with LLM",
		Status:      t.getSubtaskStatus(t.state.LlmDetectedProductChanges != nil, subtasks[len(subtasks)-1].Status),
		Subsubtasks: []api.SubSubtask{},
	}
	subtasks = append(subtasks, subtask2)

	subtask3 := api.Subtask{
		Description: "Generate search queries with LLM",
		Status:      t.getSubtaskStatus(t.state.LlmSuggestedSearchQueries != nil, subtasks[len(subtasks)-1].Status),
		Subsubtasks: []api.SubSubtask{},
	}
	subtasks = append(subtasks, subtask3)

	subtask4 := api.Subtask{
		Description: "Search for relevant documentation",
		Status:      t.getSubtaskStatus(t.state.HotParagraphs != nil, subtasks[len(subtasks)-1].Status),
		Subsubtasks: []api.SubSubtask{},
	}
	subtasks = append(subtasks, subtask4)

	subtask5 := api.Subtask{
		Description: "Rephrase documentation with LLM",
		Status:      t.getSubtaskStatus(t.state.LlmRephrasedParagraphContents != nil, subtasks[len(subtasks)-1].Status),
		Subsubtasks: []api.SubSubtask{},
	}
	subtasks = append(subtasks, subtask5)

	subtask6 := api.Subtask{
		Description: "Create drafts with rephrased content",
		Status:      t.getSubtaskStatus(t.state.CreatedDraftIds != nil, subtasks[len(subtasks)-1].Status),
		Subsubtasks: []api.SubSubtask{},
	}
	subtasks = append(subtasks, subtask6)

	if t.state.LlmSuggestedSearchQueries != nil {
		for i, query := range *t.state.LlmSuggestedSearchQueries {
			if i >= 3 {
				break
			}
			subsubtask := api.SubSubtask{
				Description: fmt.Sprintf("Search query: %s", query),
				Status:      t.getSubtaskStatus(t.state.HotParagraphs != nil, subtasks[len(subtasks)-1].Status),
			}
			subtasks[3].Subsubtasks = append(subtasks[3].Subsubtasks, subsubtask)
		}
	}

	if t.state.HotParagraphs != nil {
		if t.state.LlmRephrasedParagraphContents == nil {
			temp := make([]string, 0)
			t.state.LlmRephrasedParagraphContents = &temp
		}
		for i, paragraph := range *t.state.HotParagraphs {
			if i >= 3 {
				break
			}
			subsubtask := api.SubSubtask{
				Description: fmt.Sprintf("Rephrasing paragraph from page %s", paragraph.PageId),
				Status:      t.getSubtaskStatus(t.state.LlmRephrasedParagraphContents != nil, subtasks[len(subtasks)-1].Status),
			}
			subtasks[4].Subsubtasks = append(subtasks[4].Subsubtasks, subsubtask)
		}
	}

	if t.state.CreatedDraftIds != nil {
		for i, draftID := range *t.state.CreatedDraftIds {
			if i >= 3 {
				break
			}
			subsubtask := api.SubSubtask{
				Description: fmt.Sprintf("Created draft %s", draftID),
				Status:      api.Done,
			}
			subtasks[5].Subsubtasks = append(subtasks[5].Subsubtasks, subsubtask)
		}
	}

	return subtasks, nil
}

func (t *gitHubAccountPRTask) getSubtaskStatus(completed bool, previousStatus api.TaskStatus) api.TaskStatus {
	if completed {
		return api.Done
	}
	if slices.Contains([]api.TaskStatus{api.FailedByError, api.Cancelled, api.FailedByTimeout}, previousStatus) {
		return api.Cancelled
	}
	if t.status == api.FailedByError || t.status == api.FailedByTimeout {
		return t.status
	}
	return api.Executing
}

func (t *gitHubAccountPRTask) saveChanges() error {
	taskState := internals.TaskState{}
	taskState.FromTaskStateGitHubAccountPR(t.state)
	err := t.repo.SetTaskState(t.taskID, taskState)
	if err != nil {
		return err
	}

	return t.repo.Commit()
}

func (t *gitHubAccountPRTask) accountResult(result internals.TaskActionResult) error {
	discriminator, err := result.Discriminator()
	if err != nil {
		return err
	}
	if internals.TaskActionType(discriminator) != internals.AskLlm {
		return nil
	}

	resCast, err := result.AsTaskActionResultAskLLM()
	if err != nil {
		return err
	}
	if t.state.LlmDetectedProductChanges == nil {
		res := strings.Split(resCast.ResponseMessage, "\n")
		t.deps.Logger.Info("{{{{{{{}}}}}}}", resCast.ResponseMessage)
		t.state.LlmDetectedProductChanges = &res
		return nil
	}
	if t.state.LlmSuggestedSearchQueries == nil {
		res := strings.Split(resCast.ResponseMessage, "\n")
		t.deps.Logger.Info("((((((((((((((()))))))))))))))", resCast.ResponseMessage)
		t.state.LlmSuggestedSearchQueries = &res
		return nil
	}

	rephrased := t.state.LlmRephrasedParagraphContents
	*rephrased = append(*rephrased, resCast.ResponseMessage)

	return nil
}

func (t *gitHubAccountPRTask) OnActionResult(result internals.TaskActionResult) error {
	t.accountResult(result)

	if t.state.PrPatch == nil || t.state.PrDescription == nil {
		if err := t.fetchPRData(); err != nil {
			return fmt.Errorf("failed to fetch PR data: %w", err)
		}
	}

	if t.state.LlmDetectedProductChanges == nil {
		if err := t.askLLMForProductChanges(); err != nil {
			return fmt.Errorf("failed to ask LLM for product changes: %w", err)
		}
		return t.saveChanges()
	}

	if t.state.LlmSuggestedSearchQueries == nil || len(*t.state.LlmSuggestedSearchQueries) == 0 {
		if err := t.askLLMForSearchQueries(); err != nil {
			return fmt.Errorf("failed to ask LLM for search queries: %w", err)
		}
		return t.saveChanges()
	}

	maybeHotParagraphs := t.state.HotParagraphs
	maybeRephrased := t.state.LlmRephrasedParagraphContents

	if maybeHotParagraphs == nil {
		fmt.Println("################################## here")
		err := t.searchHotParagraphs()
		if err != nil {
			return fmt.Errorf("failed to search hot parageaphs: %w", err)
		}

		if err := t.startLLMSearchAndRephrase(); err != nil {
			return fmt.Errorf("failed to perform search and rephrase: %w", err)
		}
		return t.saveChanges()
	}

	if maybeRephrased != nil && len(*maybeRephrased) != len(*maybeHotParagraphs) {
		if err := t.startLLMSearchAndRephrase(); err != nil {
			return fmt.Errorf("failed to perform search and rephrase: %w", err)
		}
	}

	if t.state.LlmRephrasedParagraphContents != nil && len(*t.state.LlmRephrasedParagraphContents) > 0 {
		if err := t.createDraftsWithRephrasedParagraphs(); err != nil {
			return fmt.Errorf("failed to create drafts with rephrased paragraphs: %w", err)
		}
	}

	return t.saveChanges()
}

func (t *gitHubAccountPRTask) fetchPRData() error {
	parsedURL, err := url.Parse(t.state.PrUrl)
	if err != nil {
		return fmt.Errorf("failed to parse PR URL: %w", err)
	}

	// Expected format: https://github.com/owner/repo/pull/number
	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) < 4 || pathParts[2] != "pull" {
		return fmt.Errorf("invalid PR URL format")
	}

	owner := pathParts[0]
	repoName := pathParts[1]
	prNumber, err := strconv.Atoi(pathParts[3])
	if err != nil {
		return fmt.Errorf("failed to parse PR number: %w", err)
	}

	prResponse, err := t.deps.GitHubClient.GetPullRequest(t.ctx, owner, repoName, prNumber)
	if err != nil {
		return fmt.Errorf("failed to get PR from GitHub: %w", err)
	}

	t.deps.Logger.Info("fetching patch: ", prResponse.PatchUrl)
	patchResponse, err := http.Get(prResponse.PatchUrl)
	if err != nil {
		return fmt.Errorf("failed to get PR patch from GitHub: %w", err)
	}
	t.deps.Logger.Info("patch fetched: ", patchResponse.StatusCode, " ", patchResponse.ContentLength)
	defer patchResponse.Body.Close()

	patchBytes, err := io.ReadAll(patchResponse.Body)
	if err != nil {
		return err
	}
	patchStr := string(patchBytes)

	t.state.PrPatch = &patchStr
	t.state.PrDescription = &prResponse.Body

	return nil
}

func (t *gitHubAccountPRTask) askLLMForProductChanges() error {
	prompt := fmt.Sprintf(`Проанализируй этот GitHub PR и определи продуктовые изменения.
Сфокусируйся на продуктовых изменениях, а не на деталях реализации.
Обозначь максимум 2 изменения. Если значимых продуктовых изменений
не было, просто напиши NO_CHANGES и ничего больше не пиши.
Не учитывай исправление опечаток или что-то схожее.

Вот информация:
Описание PR: %s
PR Patch: %s`,
		*t.state.PrDescription, *t.state.PrPatch)

	messages := []internals.LLMMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return t.createAskLLMTaskAction(messages)
}

func (t *gitHubAccountPRTask) askLLMForSearchQueries() error {
	productChanges := ""
	if t.state.LlmDetectedProductChanges != nil {
		for _, change := range *t.state.LlmDetectedProductChanges {
			productChanges += fmt.Sprintf("- %s\n", change)
		}
	}

	prompt := fmt.Sprintf(`Я перечислю изменения в продукте. На основе их предложи поисковые запросы.
У компании есть база знаний, по которой есть семантический поиск. Надо сформировать поисковые запросы, которые могут
найти такие фрагменты базы знаний, которые теоретически надо изменить при внесении изменений в продукт.
Сформулируй как можно меньше запросов. Максимум - 3 запроса.

Вот продуктовые изменения:
%s
`, productChanges)

	messages := []internals.LLMMessage{
		{
			Role:    string(ycloud_client_gen.User),
			Content: prompt,
		},
	}

	return t.createAskLLMTaskAction(messages)
}

func (t *gitHubAccountPRTask) searchHotParagraphs() error {
	if t.state.LlmSuggestedSearchQueries == nil || len(*t.state.LlmSuggestedSearchQueries) == 0 {
		t.deps.Logger.Error("**************** no suggested search queries")
		return fmt.Errorf("no suggested search queries")
	}
	fmt.Println("YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY")

	hp := make([]internals.ParagraphWithContext, 0)

	for _, query := range *t.state.LlmSuggestedSearchQueries {
		embedding, err := t.deps.InferenceClient.GenerateEmbedding(t.ctx, query)
		if err != nil {
			t.deps.Logger.Warn("failed to generate embedding for query: %v", err)
			continue
		}

		hotParagraphs, err := t.repo.SearchByEmbeddingWithContext(query, internals.Embedding(embedding), 1)
		if err != nil {
			t.deps.Logger.Warn("failed to search by embedding: %v", err)
			continue
		}

		t.deps.Logger.Info("found %d hot paragraphs for query: %s", len(hotParagraphs), query)
		t.deps.Logger.Info(hotParagraphs)

		hp = append(hp, hotParagraphs...)
	}

	t.state.HotParagraphs = &hp
	temp := make([]string, 0)
	t.state.LlmRephrasedParagraphContents = &temp
	return nil
}

func (t *gitHubAccountPRTask) startLLMSearchAndRephrase() error {
	if len(*t.state.HotParagraphs) == 0 {
		return fmt.Errorf("empty hot paragraphs")
	}
	currentParagraphContent := (*t.state.HotParagraphs)[len(*t.state.LlmRephrasedParagraphContents)].Content

	prompt := fmt.Sprintf(`Измени это с учётом изменений в коде, чтобы текст соответствовал действительности.
Я тебе напишу какую-то информацию из корпоративной wiki. Твой ответ должен включать только отредактированную информацию.
Постарайся внести только те изменения, которые действительно произошли.

Вот информация:
Описание PR: %s
PR Patch: %s`,
		*t.state.PrDescription, *t.state.PrPatch)

	messages := []internals.LLMMessage{
		{
			Role:    string(ycloud_client_gen.System),
			Content: prompt,
		},
		{
			Role:    "user",
			Content: currentParagraphContent,
		},
	}

	return t.createAskLLMTaskAction(messages)
}

func (t *gitHubAccountPRTask) createDraftsWithRephrasedParagraphs() error {
	if t.state.HotParagraphs == nil || t.state.LlmRephrasedParagraphContents == nil {
		return nil
	}

	hotParagraphs := *t.state.HotParagraphs
	rephrasedParagraphs := *t.state.LlmRephrasedParagraphContents

	createdDraftIDs := make([]api.DraftID, 0)

	for i, paragraph := range hotParagraphs {
		if i >= len(rephrasedParagraphs) {
			break
		}

		page, _, err := t.repo.GetPageByID(paragraph.PageId)
		if err != nil {
			t.deps.Logger.Warn("failed to get page by ID: %v", err)
			continue
		}

		newContent := strings.Replace(page.Content, paragraph.Content, rephrasedParagraphs[i], 1)

		draftID, err := t.repo.CreateDraft(paragraph.PageId, page.Title, newContent)
		if err != nil {
			t.deps.Logger.Warn("failed to create draft: %v", err)
			continue
		}

		createdDraftIDs = append(createdDraftIDs, *draftID)
		break
	}

	t.state.CreatedDraftIds = &createdDraftIDs

	return nil
}

func (t *gitHubAccountPRTask) createAskLLMTaskAction(prompt []internals.LLMMessage) error {
	taskAction := internals.TaskAction{}
	taskAction.FromTaskActionAskLLM(internals.TaskActionAskLLM{
		TaskActionType: internals.AskLlm,
		Model:          internals.Yandexgpt5Lite,
		Messages:       prompt,
	})

	taskActionID, err := t.repo.CreateTaskAction(t.taskID, taskAction)
	if err != nil {
		return err
	}

	return t.repo.EnqueueTaskAction(*taskActionID)
}
