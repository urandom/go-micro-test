package main

import (
	"context"
	"log"
	"net/http"

	"github.com/urandom/go-micro-test/auth"
)

type loginHandler struct {
	token auth.TokenClient
}

func (h loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing post form: %v", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	user := r.Form.Get("user")
	pass := r.Form.Get("password")

	resp, err := h.token.Generate(context.TODO(), &auth.TokenGenerateRequest{user, pass})
	if err != nil {
		log.Printf("Error generating jwt token for %s:%s: %v", user, pass, err)
		http.Error(w, "invalid auth token", http.StatusUnauthorized)
		return
	}

	if !resp.Exists {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	w.Write([]byte(resp.Auth))

	return
}
