package main

import (
	"golang.org/x/net/context"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
)

type checkerHandler struct {
	dbProfiler db.ProfilerClient
}

func (h *checkerHandler) Check(ctx context.Context, req *auth.CheckRequest, resp *auth.CheckResponse) error {
	user, err := parseToken(req.Auth)
	if err == nil {
		resp.Valid = true
		resp.User = user
	} else if vErr, ok := errors.Cause(err).(*jwt.ValidationError); ok {
		if vErr.Errors&jwt.ValidationErrorExpired > 0 {
			resp.Expired = true
			return nil
		}
	}
	return err
}

func parseToken(tokenString string) (user string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return auth.SecretKey, nil
	})

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid && claims.Valid() == nil {
		return claims.Subject, nil
	} else {
		return "", errors.Wrap(err, "parsing jwt token string")
	}
}
