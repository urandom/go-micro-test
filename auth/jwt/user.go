package main

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
)

type userHandler struct {
	dbUser db.UserClient
}

func (h *userHandler) Profile(ctx context.Context, req *auth.UserProfileRequest, resp *auth.UserProfileResponse) error {
	if user, err := parseToken(req.Auth); err == nil {
		if ur, err := h.dbUser.Profile(ctx, &db.UserProfileRequest{user}); err == nil {
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
