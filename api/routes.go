package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (s *Server) routes() http.Handler {
	mux := httprouter.New()
	mux.HandlerFunc(http.MethodGet, "/", s.home)
	protected := alice.New(s.withJWTAuth)

	mux.HandlerFunc(http.MethodGet, "/posts", s.viewPosts)
	mux.Handler(http.MethodPost, "/posts", protected.ThenFunc(s.createPost))
	mux.HandlerFunc(http.MethodGet, "/posts/:id", s.viewPostById)

	mux.HandlerFunc(http.MethodGet, "/comments/:id", s.viewCommentsByPostId)
	mux.Handler(http.MethodPost, "/comments/:id", protected.ThenFunc(s.createComment))

	mux.HandlerFunc(http.MethodPost, "/login", s.login)
	mux.HandlerFunc(http.MethodPost, "/register", s.register)

	standard := alice.New(s.logRequest)

	return standard.Then(mux)
}
