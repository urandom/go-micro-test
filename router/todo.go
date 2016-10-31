package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/urandom/go-micro-test/todo"
)

type todoHandler struct {
	todo todo.TodoClient
}

func (h todoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	token := checkAuth(r)
	if token == "" {
		http.Error(w, "wrong auth", http.StatusBadRequest)
		return
	}

	resp, err := h.todo.List(context.TODO(), &todo.ListRequest{Auth: token})
	if err != nil {
		log.Printf("Error getting todo list: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if !resp.Valid {
		http.NotFound(w, r)
		return
	}

	b, err := json.Marshal(resp.Items)
	if err != nil {
		log.Printf("Error marshaling todo items to json: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
