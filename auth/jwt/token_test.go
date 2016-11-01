// +build go1.7

package main

import (
	"fmt"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/micro/go-micro/client"
	"github.com/pkg/errors"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
	"golang.org/x/net/context"
)

func TestGenerate(t *testing.T) {
	cases := []struct {
		dbexists bool
		profile  db.UserProfile
		exists   bool
		err      error
		user     string
		pass     string
	}{
		{
			dbexists: true,
			profile:  db.UserProfile{"foo", "", "foopass"},
		},
		{
			dbexists: true,
			profile:  db.UserProfile{"foo", "", "foopass"},
			exists:   true,
			user:     "foo",
			pass:     "foopass",
		},
		{
			dbexists: true,
			profile:  db.UserProfile{"foo", "", "foopass"},
			exists:   true,
			user:     "foo",
			pass:     "foopass",
		},
		{
			dbexists: true,
			profile:  db.UserProfile{"foo", "", "foopass"},
			exists:   false,
			user:     "foo",
			pass:     "foopass2",
		},
		{
			err: errors.New("err"),
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			h := tokenHandler{userClient{tc.dbexists, tc.profile, tc.err}}

			req := &auth.TokenGenerateRequest{tc.user, tc.pass}
			resp := &auth.TokenGenerateResponse{}

			if err := h.Generate(context.TODO(), req, resp); err != nil {
				if tc.err == nil && errors.Cause(err) != tc.err {
					t.Fatal(err)
				}
			}

			if resp.Exists != tc.exists {
				t.Fatalf("expected exists %v, got %v", tc.exists, resp.Exists)
			}

			if !tc.exists {
				return
			}

			if u, err := parseToken(resp.Auth); err != nil {
				t.Fatal(err)
			} else if u != tc.user {
				t.Fatalf("expected user %s, got %s", tc.user, u)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	cases := []struct {
		dbexists bool
		profile  db.UserProfile
		err      error
		user     string
		expired  bool
		valid    bool
	}{
		{
			dbexists: true,
			profile:  db.UserProfile{"foo", "", "foopass"},
		},
		{
			dbexists: true,
			profile:  db.UserProfile{"foo", "", "foopass"},
			user:     "foo",
			valid:    true,
		},
		{
			dbexists: true,
			profile:  db.UserProfile{"foo", "", "foopass"},
			user:     "foo",
			expired:  true,
		},
		{
			err: errors.New("err"),
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			h := tokenHandler{userClient{tc.dbexists, tc.profile, tc.err}}

			expiresAt := time.Now().Add(time.Hour).Unix()
			if tc.expired {
				expiresAt = time.Now().Add(-1 * time.Hour).Unix()
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
				Subject:   tc.user,
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: expiresAt,
			})

			tokenStr, err := token.SignedString(auth.SecretKey)
			if err != nil {
				t.Fatal(err)
			}

			req := &auth.TokenCheckRequest{tokenStr}
			resp := &auth.TokenCheckResponse{}

			if err := h.Check(context.TODO(), req, resp); err != nil {
				if tc.err == nil && errors.Cause(err) != tc.err {
					t.Fatal(err)
				}
			}

			if resp.Expired != tc.expired {
				t.Fatalf("expected expired %v, got %v", tc.expired, resp.Expired)
			}

			if resp.Valid != tc.valid {
				t.Fatalf("expected valid %v, got %v", tc.valid, resp.Valid)
			}

			if tc.valid {
				if resp.User != tc.user {
					t.Fatalf("expected user %v, got %v", tc.user, resp.User)
				}
			}
		})
	}
}

type userClient struct {
	exists  bool
	profile db.UserProfile
	err     error
}

func (c userClient) Profile(ctx context.Context, in *db.UserProfileRequest, opts ...client.CallOption) (*db.UserProfileResponse, error) {
	if c.err != nil {
		return nil, c.err
	}

	if c.profile.User != in.User {
		return &db.UserProfileResponse{false, nil}, nil
	}

	resp := &db.UserProfileResponse{c.exists, &c.profile}

	return resp, nil
}
