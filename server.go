package gap

import (
	"fmt"
	"net/http"
)

type Server struct {
	Options *Options
}

func NewServer(options *Options) *Server {
	return &Server{
		Options: options,
	}
}

func (server *Server) Run() error {
	http.HandleFunc("/_ping", HandlePing)
	ah, err := NewAuthHandler(server.Options)

	if err != nil {
		return err
	}

	http.Handle("/", ah)
	return http.ListenAndServe(fmt.Sprintf(":%d", server.Options.Port), nil)
}
