package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/R894/go-blog/models"
	"github.com/R894/go-blog/utils"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		s.notFound(w)
		return
	}
	utils.SendApiMessage(w, http.StatusOK, "Hi")
}

func (s *Server) viewPostById(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		s.notFound(w)
		return
	}
	post, err := s.db.GetPostById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.notFound(w)
			return
		}
		s.clientError(w, http.StatusInternalServerError)
		return
	}

	if post == nil {
		s.notFound(w)
		return
	}
	utils.WriteJSON(w, http.StatusOK, post)
}

func (s *Server) createPost(w http.ResponseWriter, r *http.Request) {
	var newPost models.NewPostRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newPost); err != nil {
		s.clientError(w, http.StatusBadRequest)
		return
	}

	token, err := utils.GetBearerHeader(r)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	userId, err := utils.GetUserIdFromJWT(token)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	newPost.UserId = userId

	id, err := s.db.CreatePost(newPost)

	if err != nil {
		s.clientError(w, http.StatusBadRequest)
		return
	}
	utils.SendApiMessage(w, http.StatusCreated, strconv.Itoa(id))
}

func (s *Server) viewCommentsByPostId(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	comment, err := s.db.GetCommentsByPostId(id)
	if err != nil {
		s.clientError(w, http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, comment)
}

func (s *Server) createComment(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		s.notFound(w)
		return
	}

	var newComment models.CreateCommentRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newComment); err != nil {
		s.clientError(w, http.StatusBadRequest)
		return
	}
	newComment.PostId = id

	token, err := utils.GetBearerHeader(r)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	userId, err := utils.GetUserIdFromJWT(token)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	newComment.UserId = userId

	err = s.db.CreateComment(newComment)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, "New comment created")
}

func (s *Server) viewPosts(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	limit := 6
	page, err := strconv.Atoi(params.ByName("page"))
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	fmt.Println(offset)

	posts, totalPages, err := s.db.GetPosts(6, offset)
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	// Return a response containing the total pages
	response := map[string]any{
		"posts":      posts,
		"totalPages": totalPages,
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Server) deletePost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		s.notFound(w)
		return
	}

	token, err := utils.GetBearerHeader(r)
	if err != nil {
		s.clientError(w, http.StatusForbidden)
		return
	}
	userId, _ := utils.GetUserIdFromJWT(token)
	postUserId, err := s.db.GetPostById(id)
	if err != nil {
		s.clientError(w, http.StatusNotFound)
		return
	}

	if postUserId.UserId != userId {
		s.clientError(w, http.StatusForbidden)
		return
	}

	err = s.db.DeletePostById(id)
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	utils.SendApiMessage(w, http.StatusOK, "Post deleted successfully")
}

func (s *Server) updatePost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		s.notFound(w)
		return
	}

	token, err := utils.GetBearerHeader(r)
	if err != nil {
		s.clientError(w, http.StatusForbidden)
		return
	}
	userId, _ := utils.GetUserIdFromJWT(token)
	postUserId, _ := s.db.GetPostById(id)

	if postUserId.UserId != userId {
		s.clientError(w, http.StatusForbidden)
		return
	}

	var updatePost models.UpdatePostRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatePost); err != nil {
		s.clientError(w, http.StatusBadRequest)
		return
	}

	err = s.db.UpdatePostById(id, updatePost)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	utils.SendApiMessage(w, http.StatusOK, "Post updated successfully")
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		s.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := s.db.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		s.clientError(w, http.StatusUnauthorized)
		return
	}
	token, err := utils.CreateJWT(id)
	if err != nil {
		s.clientError(w, http.StatusUnauthorized)
	}
	utils.SendApiMessage(w, http.StatusOK, token)
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var newUser models.CreateUserRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		s.serverError(w, r, err)
		return
	}
	fmt.Println(newUser)
	id, err := s.db.CreateUser(newUser)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	jwt, _ := utils.CreateJWT(id)
	fmt.Println(jwt)
	utils.SendApiMessage(w, http.StatusOK, "User created successfully")
}
