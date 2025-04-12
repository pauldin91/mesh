package api

import (
	"net/http"
)

// uploadHandler handles file uploads
// @Summary      Uploads a file
// @Description  Upload a file via multipart/form-data
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Param        file        formData  file   true  "File to upload"
// @Param        description formData  string false "Optional description"
// @Success      200 {string} string
// @Router       /upload [post]
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`"File uploaded"`)) // Proper JSON-encoded string
}
