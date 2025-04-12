package api

import (
	"net/http"

	"github.com/pauldin91/common/web"
	_ "github.com/pauldin91/mesh/document_service/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	uploadEndpoint  string = "/upload"
	swaggerEndpoint string = "/swagger/*"
	healthEndpoint  string = "/health"
)

func GetRoutes() web.Routes {
	return web.Routes{
		Hanlders: map[string]web.Route{
			healthEndpoint:  {Method: http.MethodGet, Handler: healthHandler},
			swaggerEndpoint: {Method: http.MethodGet, Handler: httpSwagger.WrapHandler},
			uploadEndpoint:  {Method: http.MethodPost, Handler: uploadHandler},
		},
	}
}
