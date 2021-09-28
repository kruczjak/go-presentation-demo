package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"webserver/pkg/accesslogger"
)

const Addr = "0.0.0.0:5000"

type Server struct {
	server *http.Server
}

func New() *Server {
	return &Server{server: &http.Server{Addr: Addr}}
}

func (s *Server) ListenAndServe() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", s.handleHello())
	r.HandleFunc("/articles/{id}", s.handleArticles())

	http.Handle("/", loggerMiddleware(r))

	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	fmt.Println("Server has started")
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		accesslogger.Push(request)

		next.ServeHTTP(writer, request)
	})
}

func (s *Server) Shutdown(ctx context.Context) {
	s.server.Shutdown(ctx)
}

