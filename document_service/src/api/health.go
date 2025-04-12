package api

import "net/http"

// healthHandler retrieves health for document service
// @Summary      Get health check
// @Description  Retrieves health for document service
// @Tags         health
// @Produce      json
// @Success      200 {string}  string
// @Router       /health [get]
func (app *HttpApplication) healthHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}
