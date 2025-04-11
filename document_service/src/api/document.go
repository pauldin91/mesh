package api

import "net/http"

// UploadFile godoc
// @Summary Uploads a file
// @Description Upload a file via multipart/form-data
// @Tags files
// @Accept mpfd
// @Produce json
// @Param file formData file true "File to upload"
// @Param description formData string false "Optional description"
// @Success 200 {object} string
// @Router /upload [post]
func (app *HttpApplication) uploadHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}
