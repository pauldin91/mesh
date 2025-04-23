package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/pauldin91/common/utils"
)

type HttpApplication struct {
	cfg    utils.Config
	server *http.Server
}

func (app *HttpApplication) SetCfg(cfg utils.Config) {
	app.cfg = cfg
}

func (app *HttpApplication) SetServer(routes Routes) {
	app.server = &http.Server{
		Addr:    app.cfg.HttpServerAddress,
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

func (app *HttpApplication) Start() {

	certFile := filepath.Join(app.cfg.CertPath, app.cfg.CertFile)
	certKey := filepath.Join(app.cfg.CertPath, app.cfg.CertKey)

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Fatal().Msg("unable to load certs")
	}

	go func() {
		log.Info().Msgf("INFO: HTTP server started on %s\n", app.cfg.HttpServerAddress)
		if app.cfg.Tls {

			if err := app.server.ListenAndServeTLS(certFile, certKey); err != nil && err != http.ErrServerClosed {
				log.Fatal().Msgf("Could not start HTTP server: %s", err)
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
