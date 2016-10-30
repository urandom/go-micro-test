package main

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
)

type profilerHandler struct {
	dbProfiler db.ProfilerClient
}

func (h *profilerHandler) UserProfile(ctx context.Context, req *auth.ProfileRequest, resp *auth.ProfileResponse) error {
	if user, err := parseToken(req.Auth); err == nil {
		if ur, err := h.dbProfiler.UserProfile(ctx, &db.UserRequest{user}); err == nil {
			resp.Exists = ur.Exists
			if ur.Exists {
				resp.Profile = &auth.AuthProfile{User: ur.Profile.User, Name: ur.Profile.Name}
			}
			return nil
		} else {
			return errors.Wrap(err, "fetching user profile from db")
		}
	} else {
		return errors.Wrap(err, "parsing jwt token")
	}
}
