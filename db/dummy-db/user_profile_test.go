// +build go1.7

package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/urandom/go-micro-test/db"
)

func TestUserProfile(t *testing.T) {
	cases := []struct {
		user   string
		exists bool
		name   string
		pass   string
	}{
		{"foo", true, "full name 1", "foopass"},
		{"bar", true, "full name 2", "barpass"},
		{"asdasd", false, "", ""},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			h := userHandler{}
			req := &db.UserProfileRequest{tc.user}
			resp := &db.UserProfileResponse{}

			if err := h.Profile(context.TODO(), req, resp); err != nil {
				t.Fatal(err)
			}

			if resp.Exists != tc.exists {
				t.Fatalf("expected exists %v, got %v", tc.exists, resp.Exists)
			}

			if !tc.exists {
				return
			}

			if resp.Profile.Name != tc.name {
				t.Fatalf("expected name %s, got %s", tc.name, resp.Profile.Name)
			}

			if resp.Profile.User != tc.user {
				t.Fatalf("expected user %s, got %s", tc.user, resp.Profile.User)
			}

			if resp.Profile.Pass != tc.pass {
				t.Fatalf("expected pass %s, got %s", tc.pass, resp.Profile.Pass)
			}
		})
	}
}
