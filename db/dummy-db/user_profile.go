package main

import (
	"golang.org/x/net/context"

	"github.com/urandom/go-micro-test/db"
)

type userHandler struct{}

func (h *userHandler) Profile(ctx context.Context, req *db.UserProfileRequest, resp *db.UserProfileResponse) error {
	resp.Profile, resp.Exists = profileData[req.User]

	return nil
}

var (
	profileData = map[string]*db.UserProfile{
		"foo": {"foo", "full name 1", "foopass"},
		"bar": {"bar", "full name 2", "barpass"},
	}
)
