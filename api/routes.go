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

	// Posts
	mux.HandlerFunc(http.MethodGet, "/posts/page/:page", s.viewPosts)
	mux.HandlerFunc(http.MethodGet, "/posts/view/:id", s.viewPostById)
	mux.Handler(http.MethodDelete, "/posts/delete/:id", protected.ThenFunc(s.deletePost))
	mux.Handler(http.MethodPut, "/posts/update/:id", protected.ThenFunc(s.updatePost))
	mux.Handler(http.MethodPost, "/posts", protected.ThenFunc(s.createPost))

	// Comments
	mux.HandlerFunc(http.MethodGet, "/comments/post/:id", s.viewCommentsByPostId)
	mux.Handler(http.MethodPost, "/comments/post/:id", protected.ThenFunc(s.createComment))

	// Authentication
	mux.HandlerFunc(http.MethodPost, "/login", s.login)
	mux.HandlerFunc(http.MethodPost, "/register", s.register)

	standard := alice.New(s.logRequest, s.secureHeaders)

	return standard.Then(mux)
}
