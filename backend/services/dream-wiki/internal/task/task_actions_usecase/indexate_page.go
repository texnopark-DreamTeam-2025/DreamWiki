package task_actions_usecase

import (
	"fmt"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/indexing"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

func (u *taskActionUsecaseImpl) indexatePageInTransaction(repo repository.AppRepository, pageID api.PageID) error {
	err := repo.RemovePageIndexation(pageID)
	if err != nil {
		return err
	}

	page, _, err := repo.GetPageByID(pageID)
	if err != nil {
		return err
	}

	paragraphs := indexing.SplitPageToParagraphs(pageID, page.Content)

	contentStrings := make([]string, len(paragraphs))
	for i, paragraph := range paragraphs {
		contentStrings[i] = paragraph.Content
	}

	embeddings, err := u.deps.InferenceClient.GenerateEmbeddings(u.ctx, contentStrings)
	if err != nil {
		return err
	}

	stems, err := u.deps.InferenceClient.GenerateStems(u.ctx, contentStrings)
	if err != nil {
		return err
	}

	for i, paragraph := range paragraphs {
		paragraph.Embedding = embeddings[i]

		err = repo.AddIndexedParagraph(paragraph)
		if err != nil {
			return err
		}

		termCount := make(map[string]int64)
		for _, stem := range stems[i] {
			termCount[stem]++
		}

		terms := make([]internals.Term, 0, len(termCount))
		for term, count := range termCount {
			if count == 0 {
				continue
			}
			terms = append(terms, internals.Term{
				Term:           term,
				PageId:         pageID,
				ParagraphIndex: int64(paragraph.ParagraphIndex),
				TimesIn:        count,
			})
		}

		if len(terms) > 0 {
			err = repo.AddTerms(terms)
			if err != nil {
				return fmt.Errorf("failed to add terms for page %s: %w", pageID, err)
			}
		}
	}

	return nil
}

func (u *taskActionUsecaseImpl) executeIndexatePageAction(repo repository.AppRepository, actionID internals.TaskActionID, taskAction *internals.TaskAction) error {
	indexatePageAction, err := taskAction.AsTaskActionIndexatePage()
	if err != nil {
		return fmt.Errorf("failed to parse task action as TaskActionIndexatePage: %w", err)
	}

	err = u.indexatePageInTransaction(repo, indexatePageAction.PageId)
	if err != nil {
		return fmt.Errorf("failed to indexate page: %w", err)
	}

	err = repo.SetTaskActionStatus(actionID, internals.Finished)
	if err != nil {
		return fmt.Errorf("failed to set task action status to finished: %w", err)
	}

	result := internals.TaskActionResult{}
	indexatePageResult := internals.TaskActionResultIndexatePage{
		TaskActionType: internals.IndexatePage,
		PageId:         indexatePageAction.PageId,
	}
	err = result.FromTaskActionResultIndexatePage(indexatePageResult)
	if err != nil {
		return fmt.Errorf("failed to create task action result: %w", err)
	}

	err = repo.CreateTaskActionResult(actionID, result)
	if err != nil {
		return fmt.Errorf("failed to create task action result: %w", err)
	}

	err = repo.EnqueueTaskActionResult(actionID)
	if err != nil {
		return fmt.Errorf("failed to enqueue task action result: %w", err)
	}

	return nil
}
