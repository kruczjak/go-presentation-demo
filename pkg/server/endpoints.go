package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type ArticlesResponse struct {
	ID int `json:"id"`
	Time time.Time `json:"time"`
	RequestedURL string `json:"request_url"`
}


func (s *Server) handleHello() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		//fmt.Printf("incoming request %v\n", request)

		//writer.Write([]byte("world"))

		http.Error(writer, "sorry", http.StatusBadRequest)
	}
}

func (s *Server) handleArticles() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		fmt.Printf("vars %v\n", vars)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(writer, "ID is not integer", http.StatusUnprocessableEntity)
			return
		}

		response := &ArticlesResponse{
			ID:           id,
			Time:         time.Now(),
			RequestedURL: request.URL.String(),
		}

		respJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, "can't marshal", http.StatusUnprocessableEntity)
			return
		}

		writer.Write(respJSON)
	}
}
