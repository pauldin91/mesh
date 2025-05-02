package main

import (
	"context"
	"log"
	"time"

	"github.com/pauldin91/common/utils"
	"github.com/pauldin91/common/web"
	"github.com/pauldin91/mesh/registry/src/api"
)

func main() {
	s := web.NewBuilder(func() *web.HttpApplication {
		return &web.HttpApplication{}
	})

	cfg, err := utils.LoadConfig("../../registry")
	if err != nil {
		log.Fatalln("could not load cfg")
	}
	httpServer := s.
		WithServer(cfg.HttpServerAddress, api.GetRoutes()).
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	httpServer.Start()
	httpServer.WaitForShutdown(ctx)
}
