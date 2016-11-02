package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/urandom/go-micro-test/auth"
)

type profileHandler struct {
	user auth.UserClient
}

func (h profileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	token := checkAuth(r)
	if token == "" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	resp, err := h.user.Profile(context.TODO(), &auth.UserProfileRequest{token})
	if err != nil {
		log.Printf("Error getting user profile: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if !resp.Exists {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	b, err := json.Marshal(resp.Profile)
	if err != nil {
		log.Printf("Error marshaling profile to json: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func checkAuth(r *http.Request) string {
	authH := strings.Fields(r.Header.Get("Authorization"))
	if len(authH) != 2 || authH[0] != "Bearer" {
		log.Printf("Unexpected authorization header: %s", strings.Join(authH, " "))
		return ""
	}

	return authH[1]
}
