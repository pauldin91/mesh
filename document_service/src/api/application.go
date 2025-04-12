package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	_ "github.com/pauldin91/mesh/document_service/docs"
	"github.com/pauldin91/mesh/document_service/src/utils"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Application interface {
	Start()
	SetServer()
	SetCfg(cfg utils.Config)
	WaitForShutdown(context.Context)
}

type HttpApplication struct {
	cfg        utils.Config
	server     *http.Server
	signalChan chan os.Signal
}

func (app *HttpApplication) setRouter() *chi.Mux {
	var router *chi.Mux = chi.NewRouter()
	router.Get(swaggerEndpoint, httpSwagger.WrapHandler)
	router.Get(healthEndpoint, app.healthHandler)
	router.Post(uploadEndpoint, app.uploadHandler)
	return router

}

func (app *HttpApplication) SetCfg(cfg utils.Config) {
	app.cfg = cfg
}

func (app *HttpApplication) SetServer() {
	app.server = &http.Server{
		Addr:    app.cfg.HttpServerAddress,
		Handler: app.setRouter(),
	}
}

func (app *HttpApplication) Start() {

	certFile := filepath.Join(app.cfg.CertPath, app.cfg.CertFile)
	certKey := filepath.Join(app.cfg.CertPath, app.cfg.CertKey)

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Fatal().Msg("unable to load certs")
	}

	go func() {
		log.Printf("INFO: HTTP server started on %s\n", app.cfg.HttpServerAddress)
		if err := app.server.ListenAndServeTLS(certFile, certKey); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Could not start HTTP server: %s", err)
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
