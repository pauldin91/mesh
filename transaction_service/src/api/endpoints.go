package api

import (
	"net/http"

	"github.com/pauldin91/common/web"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	transactionEndpoint string = "/transaction"
	swaggerEndpoint     string = "/swagger/*"
	healthEndpoint      string = "/health"
)

func GetRoutes() web.Routes {
	return web.Routes{
		Hanlders: map[string]web.Route{
			healthEndpoint:      web.Route{Method: http.MethodGet, Handler: healthHandler},
			swaggerEndpoint:     web.Route{Method: http.MethodGet, Handler: httpSwagger.WrapHandler},
			transactionEndpoint: web.Route{Method: http.MethodPost, Handler: transactionHandler},
		},
	}
}
