package task_common

import "github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"

func IsTerminalTaskStatus(status api.TaskStatus) bool {
	return status != api.Executing
}
