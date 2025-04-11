package main

import (
	"sync"

	"github.com/pauldin91/mesh/document_service/src/api"
)

func main() {

	s := api.NewBuilder(func() *api.HttpApplication {
		return &api.HttpApplication{}
	})

	httpServer := s.
		WithConfig(".").
		WithRouter().
		Build()

	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Start(done)
	}()

}
