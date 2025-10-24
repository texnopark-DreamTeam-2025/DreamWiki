package github_account_pr

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/task/task_common"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

type (
	gitHubAccountPRTask struct {
		state internals.TaskStateGitHubAccountPR
	}
)

var (
	_ task_common.TaskLogic
)

func NewGitHubAccountPRTask(prURL string) *gitHubAccountPRTask {

}
