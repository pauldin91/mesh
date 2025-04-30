package web

import (
	"context"
	"net/http"
)

type Application interface {
	Start()
	SetServer(addr string, routes Routes)
	WaitForShutdown(context.Context)
}

type Routes struct {
	Hanlders map[string]Route
}

type Route struct {
	Method  string
	Handler func(w http.ResponseWriter, r *http.Request)
}
