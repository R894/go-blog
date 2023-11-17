package handlers

import "github.com/R894/go-blog/api/interfaces"

type Handlers struct {
	server interfaces.ServerInterface
}

func NewHandlers(server interfaces.ServerInterface) *Handlers {
	return &Handlers{
		server: server,
	}
}
