package main

import (
	"context"
	"time"

	"github.com/pauldin91/common/web"
	"github.com/pauldin91/mesh/document_service/src/api"
)

func main() {

	s := web.NewBuilder(func() *web.HttpApplication {
		return &web.HttpApplication{}
	})
	cfgDir := "../../document_service"
	httpServer := s.
		WithConfig(cfgDir).
		WithServer(api.GetRoutes()).
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	httpServer.Start()
	httpServer.WaitForShutdown(ctx)
}
