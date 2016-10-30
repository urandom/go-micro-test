package main

import (
	"time"

	"golang.org/x/net/context"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
)

type generatorHandler struct {
	dbProfiler db.ProfilerClient
}

func (h *generatorHandler) Generate(ctx context.Context, req *auth.GenerateRequest, resp *auth.GenerateResponse) error {
	p, err := h.dbProfiler.UserProfile(ctx, &db.UserRequest{req.User})
	if err != nil {
		return errors.Wrap(err, "getting the user from db")
	}

	if p.Exists && p.Profile.Pass == req.Pass {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Subject:   req.User,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		})

		tokenStr, err := token.SignedString(auth.SecretKey)
		if err != nil {
			return errors.Wrap(err, "generating jwt token string")
		}

		resp.Exists = true
		resp.Auth = tokenStr
	}

	return nil
}
