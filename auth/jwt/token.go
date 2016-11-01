package main

import (
	"time"

	"golang.org/x/net/context"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
)

type tokenHandler struct {
	dbUser db.UserClient
}

func (h *tokenHandler) Generate(ctx context.Context, req *auth.TokenGenerateRequest, resp *auth.TokenGenerateResponse) error {
	p, err := h.dbUser.Profile(ctx, &db.UserProfileRequest{req.User})
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

func (h *tokenHandler) Check(ctx context.Context, req *auth.TokenCheckRequest, resp *auth.TokenCheckResponse) error {
	user, err := parseToken(req.Auth)
	if err == nil {
		p, err := h.dbUser.Profile(ctx, &db.UserProfileRequest{user})
		if err != nil {
			return errors.Wrap(err, "getting the user from db")
		}

		if p.Exists {
			resp.Valid = true
			resp.User = user
		}
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
