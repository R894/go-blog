package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/R894/go-blog/models"
	"github.com/R894/go-blog/utils"
	"github.com/julienschmidt/httprouter"
)

func (s *Handlers) ViewCommentsByPostId(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		s.server.ServerError(w, r, err)
		return
	}
	comment, err := s.server.GetDB().GetCommentsByPostId(id)
	if err != nil {
		s.server.ClientError(w, http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, comment)
}

func (s *Handlers) CreateComment(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		s.server.NotFound(w)
		return
	}

	var newComment models.CreateCommentRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newComment); err != nil {
		s.server.ClientError(w, http.StatusBadRequest)
		return
	}
	newComment.PostId = id

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
	newComment.UserId = userId

	err = s.server.GetDB().CreateComment(newComment)
	if err != nil {
		s.server.ServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, "New comment created")
}
