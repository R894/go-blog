package api

import (
	"net/http"

	"github.com/R894/go-blog/utils"
)

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		s.logger.Info("received request", "ip", ip, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) withJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := utils.GetBearerHeader(r)
		if err != nil {
			s.clientError(w, http.StatusUnauthorized)
			return
		}

		token, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.SendApiError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		if !token.Valid {
			utils.SendApiError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
	})
}
