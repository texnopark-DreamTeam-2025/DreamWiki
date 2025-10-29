package reindexate_all_pages

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

type (
	reindexateAllPagesTask struct {
		state internals.TaskStateReindexateAllPages
	}
)

var (
	_ task_common.TaskLogic = (*reindexateAllPagesTask)(nil)
)

func NewReindexateAllPagesTask(state internals.TaskStateReindexateAllPages) *reindexateAllPagesTask {
	return &reindexateAllPagesTask{
		state: state,
	}
}

func (t *reindexateAllPagesTask) CalculateSubtasks() ([]api.Subtask, error) {
	// Create subtasks based on the pages that need to be indexed
	var subtasks []api.Subtask

	// Add a subtask for each page that needs to be indexed
	if t.state.PagesToIndexateIds != nil {
		for _, pageID := range *t.state.PagesToIndexateIds {
			subtasks = append(subtasks, api.Subtask{
				Description: "Indexating page " + pageID.String(),
				Status:      api.Executing, // Default status
				Subsubtasks: []api.SubSubtask{},
			})
		}
	}

	// Add a subtask for each page that has been indexed
	for _, pageID := range t.state.IndexatedPageIds {
		subtasks = append(subtasks, api.Subtask{
			Description: "Indexating page " + pageID.String(),
			Status:      api.Done, // These pages are already indexed
			Subsubtasks: []api.SubSubtask{},
		})
	}

	return subtasks, nil
}

func (t *reindexateAllPagesTask) OnActionResult(result internals.TaskActionResult) error {
	// Handle action results if needed
	return nil
}
