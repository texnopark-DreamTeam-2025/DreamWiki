package delivery

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type AppDelivery struct {
	deps *deps.Deps
	log  logger.Logger
}

var (
	_ api.StrictServerInterface = &AppDelivery{}
)

const (
	internalErrorMessage string = "internal error"
)

func NewAppDelivery(deps *deps.Deps) *AppDelivery {
	return &AppDelivery{deps: deps, log: deps.Logger}
}
