package api

import (
	"net/http"

	"github.com/R894/go-blog/api/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (s *Server) routes() http.Handler {
	mux := httprouter.New()
	handler := handlers.NewHandlers(s)
	protected := alice.New(s.withJWTAuth)

	// Define file server folder
	dir := "./static"

	// Posts
	mux.HandlerFunc(http.MethodGet, "/posts", handler.ViewPosts)
	mux.HandlerFunc(http.MethodGet, "/posts/page/:page", handler.ViewPosts)
	mux.HandlerFunc(http.MethodGet, "/posts/view/:id", handler.ViewPostById)
	mux.Handler(http.MethodDelete, "/posts/delete/:id", protected.ThenFunc(handler.DeletePost))
	mux.Handler(http.MethodPut, "/posts/update/:id", protected.ThenFunc(handler.UpdatePost))
	mux.Handler(http.MethodPost, "/posts", protected.ThenFunc(handler.CreatePost))

	// Comments
	mux.HandlerFunc(http.MethodGet, "/comments/post/:id", handler.ViewCommentsByPostId)
	mux.Handler(http.MethodPost, "/comments/post/:id", protected.ThenFunc(handler.CreateComment))

	// Authentication
	mux.HandlerFunc(http.MethodPost, "/login", handler.Login)
	mux.HandlerFunc(http.MethodPost, "/register", handler.Register)

	// File Server
	mux.ServeFiles("/static/*filepath", http.Dir(dir))
	mux.HandlerFunc(http.MethodPost, "/upload", handler.Upload)

	standard := alice.New(s.logRequest, s.secureHeaders)

	return standard.Then(mux)
}
