package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func SendApiMessage(w http.ResponseWriter, status int, v string) error {
	return WriteJSON(w, status, Message{Message: v})
}

func SendApiError(w http.ResponseWriter, status int, v string) error {
	return WriteJSON(w, status, ErrorResponse{Error: v})
}

// GetBearerHeader extracts the bearer token from the header then returns it as a string
func GetBearerHeader(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		return "", errors.New("Unauthorized")
	}
	return strings.TrimPrefix(tokenString, "Bearer "), nil
}

func GenerateRandomFileName() string {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	return timestamp
}
