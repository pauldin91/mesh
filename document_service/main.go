package main

import (
	"context"
	"time"

	"github.com/pauldin91/mesh/document_service/src/api"
)

func main() {

	s := api.NewBuilder(func() *api.HttpApplication {
		return &api.HttpApplication{}
	})

	httpServer := s.
		WithConfig(".").
		WithServer().
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpServer.Start()
	httpServer.WaitForShutdown(ctx)
}
