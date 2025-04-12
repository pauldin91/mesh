package web

import (
	"context"
	"net/http"

	"github.com/pauldin91/common/utils"
)

type Application interface {
	Start()
	SetServer(routes Routes)
	SetCfg(cfg utils.Config)
	WaitForShutdown(context.Context)
}

type Routes struct {
	Hanlders map[string]Route
}

type Route struct {
	Method  string
	Handler func(w http.ResponseWriter, r *http.Request)
}
