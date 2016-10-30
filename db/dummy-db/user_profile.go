package main

import (
	"golang.org/x/net/context"

	"github.com/urandom/go-micro-test/db"
)

type profilerHandler struct{}

func (h *profilerHandler) UserProfile(ctx context.Context, req *db.UserRequest, resp *db.UserResponse) error {
	resp.Profile, resp.Exists = profileData[req.User]

	return nil
}

var (
	profileData = map[string]*db.UserProfile{
		"foo": {"foo", "full name 1", "foopass"},
		"bar": {"bar", "full name 2", "barpass"},
	}
)
