package web

import (
	"log"

	"github.com/pauldin91/common/utils"
)

type AppBuilder[T Application] struct {
	app T
}

func NewBuilder[T Application](new func() T) AppBuilder[T] {
	return AppBuilder[T]{
		app: new(),
	}
}

func (builder *AppBuilder[T]) WithServer(routes Routes) *AppBuilder[T] {
	builder.app.SetServer(routes)
	return builder
}

func (builder *AppBuilder[T]) WithConfig(path string) *AppBuilder[T] {
	cfg, err := utils.LoadConfig(path)
	if err != nil {
		log.Fatal("unable to load config")
	}
	builder.app.SetCfg(cfg)
	return builder
}

func (builder *AppBuilder[T]) Build() T {
	return builder.app
}
