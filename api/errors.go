package api

import (
	"net/http"

	"github.com/R894/go-blog/utils"
)

func (s *Server) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	s.logger.Error(err.Error(), "method", method, "uri", uri)
	//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	utils.SendApiError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (s *Server) ClientError(w http.ResponseWriter, status int) {
	//http.Error(w, http.StatusText(status), status)
	utils.SendApiError(w, status, http.StatusText(status))
}

func (s *Server) NotFound(w http.ResponseWriter) {
	s.ClientError(w, http.StatusNotFound)
}
