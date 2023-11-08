package api

import (
	"net/http"

	"github.com/R894/go-blog/utils"
)

func (s *Server) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	s.logger.Error(err.Error(), "method", method, "uri", uri)
	//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	utils.SendApiError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (s *Server) clientError(w http.ResponseWriter, status int) {
	//http.Error(w, http.StatusText(status), status)
	utils.SendApiError(w, status, http.StatusText(status))
}

func (s *Server) notFound(w http.ResponseWriter) {
	s.clientError(w, http.StatusNotFound)
}
