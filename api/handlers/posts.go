package handlers

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

func (s *Handlers) ViewPostById(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		s.server.NotFound(w)
		return
	}
	post, err := s.server.GetDB().GetPostById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.server.NotFound(w)
			return
		}
		s.server.ClientError(w, http.StatusInternalServerError)
		return
	}

	if post == nil {
		s.server.NotFound(w)
		return
	}
	utils.WriteJSON(w, http.StatusOK, post)
}

func (s *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	var newPost models.NewPostRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newPost); err != nil {
		s.server.ClientError(w, http.StatusBadRequest)
		return
	}

	token, err := utils.GetBearerHeader(r)
	if err != nil {
		s.server.ServerError(w, r, err)
		return
	}
	userId, err := utils.GetUserIdFromJWT(token)
	if err != nil {
		s.server.ServerError(w, r, err)
		return
	}
	newPost.UserId = userId

	id, err := s.server.GetDB().CreatePost(newPost)

	if err != nil {
		s.server.ClientError(w, http.StatusBadRequest)
		return
	}
	utils.SendApiMessage(w, http.StatusCreated, strconv.Itoa(id))
}

func (s *Handlers) ViewPosts(w http.ResponseWriter, r *http.Request) {
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

	posts, totalPages, err := s.server.GetDB().GetPosts(6, offset)
	if err != nil {
		s.server.ServerError(w, r, err)
		return
	}

	// Return a response containing the total pages
	response := map[string]any{
		"posts":      posts,
		"totalPages": totalPages,
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Handlers) DeletePost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		s.server.NotFound(w)
		return
	}

	token, err := utils.GetBearerHeader(r)
	if err != nil {
		s.server.ClientError(w, http.StatusForbidden)
		return
	}
	userId, _ := utils.GetUserIdFromJWT(token)
	postUserId, err := s.server.GetDB().GetPostById(id)
	if err != nil {
		s.server.ClientError(w, http.StatusNotFound)
		return
	}

	if postUserId.UserId != userId {
		s.server.ClientError(w, http.StatusForbidden)
		return
	}

	err = s.server.GetDB().DeletePostById(id)
	if err != nil {
		s.server.ServerError(w, r, err)
		return
	}

	utils.SendApiMessage(w, http.StatusOK, "Post deleted successfully")
}

func (s *Handlers) UpdatePost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		s.server.NotFound(w)
		return
	}

	token, err := utils.GetBearerHeader(r)
	if err != nil {
		s.server.ClientError(w, http.StatusForbidden)
		return
	}
	userId, _ := utils.GetUserIdFromJWT(token)
	postUserId, _ := s.server.GetDB().GetPostById(id)

	if postUserId.UserId != userId {
		s.server.ClientError(w, http.StatusForbidden)
		return
	}

	var updatePost models.UpdatePostRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatePost); err != nil {
		s.server.ClientError(w, http.StatusBadRequest)
		return
	}

	err = s.server.GetDB().UpdatePostById(id, updatePost)
	if err != nil {
		s.server.ServerError(w, r, err)
		return
	}
	utils.SendApiMessage(w, http.StatusOK, "Post updated successfully")
}
