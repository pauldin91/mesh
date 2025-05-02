package api

import (
	"net/http"

	"github.com/pauldin91/common/web"
)

const (
	serviceEndpoint string = "/services"
)

func GetRoutes() web.Routes {
	return web.Routes{
		Hanlders: map[string]web.Route{
			serviceEndpoint: {Method: http.MethodGet, Handler: serviceHandler},
		},
	}
}
