package interfaces

import (
	"net/http"

	"github.com/R894/go-blog/models"
)

type ServerInterface interface {
	GetDB() models.Database
	NotFound(w http.ResponseWriter)
	ClientError(w http.ResponseWriter, status int)
	ServerError(w http.ResponseWriter, r *http.Request, err error)
	Start() error
}
