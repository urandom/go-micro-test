package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
)

func TestProfile(t *testing.T) {
	cases := []struct {
		dbexists bool
		profile  db.UserProfile
		exists   bool
		err      error
		user     string
		name     string
		expired  bool
	}{
		{
			dbexists: true,
			profile:  db.UserProfile{"foo", "", "foopass"},
		},
		{
			dbexists: true,
			profile:  db.UserProfile{"foo", "full name", "foopass"},
			exists:   true,
			user:     "foo",
			name:     "full name",
		},
		{
			dbexists: true,
			profile:  db.UserProfile{"foo", "full name", "foopass"},
			exists:   false,
			user:     "foo2",
		},
		{
			err: errors.New("err"),
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			h := userHandler{userClient{tc.dbexists, tc.profile, tc.err}}

			req := &auth.UserProfileRequest{token(t, tc.user, tc.expired)}
			resp := &auth.UserProfileResponse{}

			if err := h.Profile(context.TODO(), req, resp); err != nil {
				if tc.expired {
					return
				}
				if tc.err == nil || errors.Cause(err) != tc.err {
					t.Fatal(err)
				}
			}

			if resp.Exists != tc.exists {
				t.Fatalf("expected exists %v, got %v", tc.exists, resp.Exists)
			}

			if !tc.exists {
				return
			}

			if resp.Profile.Name != tc.profile.Name {
				t.Fatalf("expected name %v, got %v", tc.name, resp.Profile.Name)
			}

			if resp.Profile.User != tc.profile.User {
				t.Fatalf("expected user %v, got %v", tc.user, resp.Profile.User)
			}
		})
	}
}
