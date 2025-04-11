package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/pauldin91/mesh/document_service/src/utils"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Application interface {
	Start(chan bool)
	SetRouter(router *chi.Mux)
	SetCfg(cfg utils.Config)
}

type HttpApplication struct {
	cfg    utils.Config
	router *chi.Mux
}

func (app *HttpApplication) Start(done chan bool) {

}
func (app *HttpApplication) SetRouter(router *chi.Mux) {
	router.Get(swaggerEndpoint, httpSwagger.WrapHandler)
	router.Get(healthEndpoint, app.healthHandler)
	router.Post(uploadEndpoint, app.uploadHandler)
	app.router = router

}

func (app *HttpApplication) SetCfg(cfg utils.Config) {
	app.cfg = cfg
}
