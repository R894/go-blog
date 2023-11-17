package handlers

import (
	"io"
	"net/http"
	"os"

	"github.com/R894/go-blog/utils"
)

func (s *Handlers) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10mb limit
	if err != nil {
		s.server.ServerError(w, r, err)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file on the server to store the uploaded file
	f, err := os.Create("./static/" + handler.Filename)
	if err != nil {
		http.Error(w, "Error creating the file on the server", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the file data to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Error copying file data", http.StatusInternalServerError)
		return
	}
	utils.SendApiMessage(w, http.StatusCreated, "File "+handler.Filename+" created")

}
