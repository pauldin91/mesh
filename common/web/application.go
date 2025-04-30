package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type HttpApplication struct {
	cert    string
	certKey string
	server  *http.Server
}

func (app *HttpApplication) SetServer(serverAddress string, routes Routes) {
	app.server = &http.Server{
		Addr:    serverAddress,
		Handler: app.setRouter(routes),
	}
}

func (app *HttpApplication) setRouter(routes Routes) *chi.Mux {
	router := chi.NewMux()
	for key, route := range routes.Hanlders {
		if route.Method == http.MethodGet {
			router.Get(key, route.Handler)
		} else if route.Method == http.MethodPost {
			router.Post(key, route.Handler)
		}
	}
	return router

}

func (app *HttpApplication) WithCertificateAndKey(cert, certKey string) {
	app.cert = cert
	app.certKey = certKey

}

func (app *HttpApplication) Start() {

	go func() {
		log.Info().Msgf("INFO: HTTP server started on %s\n", app.server.Addr)
		if app.cert != "" {
			if _, err := os.Stat(app.cert); err != nil {
				if err := app.server.ListenAndServeTLS(app.cert, app.certKey); err != nil && err != http.ErrServerClosed {
					log.Fatal().Msgf("Could not start HTTP server: %s", err)
				}
			}
		} else {
			if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Msgf("Could not start HTTP server: %s", err)
			}
		}
	}()

}

func (app *HttpApplication) WaitForShutdown(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signalChan
	log.Info().Msgf("Received signal: %s. Shutting down gracefully...", sig)

	if err := app.server.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("HTTP server Shutdown: %v", err)
	}

	log.Info().Msg("Server gracefully stopped")
}
