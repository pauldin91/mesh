package api

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/pauldin91/mesh/document_service/src/utils"
)

type ApiBuilder[T Application] struct {
	app T
}

func NewBuilder[T Application](new func() T) ApiBuilder[T] {
	return ApiBuilder[T]{
		app: new(),
	}
}

func (builder *ApiBuilder[T]) WithConfig(path string) *ApiBuilder[T] {
	cfg, err := utils.LoadConfig(path)
	if err != nil {
		log.Fatal("unable to load config")
	}
	builder.app.SetCfg(cfg)
	return builder
}

func (builder *ApiBuilder[T]) WithRouter() *ApiBuilder[T] {
	builder.app.SetRouter(chi.NewRouter())
	return builder
}

func (builder *ApiBuilder[T]) Build() T {
	return builder.app
}
