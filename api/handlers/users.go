package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/R894/go-blog/models"
	"github.com/R894/go-blog/utils"
)

func (s *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		s.server.ClientError(w, http.StatusBadRequest)
		return
	}

	id, err := s.server.GetDB().AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		s.server.ClientError(w, http.StatusUnauthorized)
		return
	}
	token, err := utils.CreateJWT(id)
	if err != nil {
		s.server.ClientError(w, http.StatusUnauthorized)
	}
	utils.SendApiMessage(w, http.StatusOK, token)
}

func (s *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var newUser models.CreateUserRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		s.server.ServerError(w, r, err)
		return
	}
	fmt.Println(newUser)
	id, err := s.server.GetDB().CreateUser(newUser)
	if err != nil {
		s.server.ServerError(w, r, err)
		return
	}
	jwt, _ := utils.CreateJWT(id)
	fmt.Println(jwt)
	utils.SendApiMessage(w, http.StatusOK, "User created successfully")
}
