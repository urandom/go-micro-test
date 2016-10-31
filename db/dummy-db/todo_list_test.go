// +build go1.7

package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/urandom/go-micro-test/db"
)

func TestTodoList(t *testing.T) {
	cases := []struct {
		user   string
		limit  int64
		offset int64
		exists bool
		total  int
		titles []string
		bodies []string
	}{
		{
			user:   "foo",
			exists: true,
			total:  4,
			titles: []string{"title1", "title2", "title3", "title4"},
			bodies: []string{"body1", "body2", "body3", "body4"},
		},
		{
			user:   "foo",
			limit:  1,
			exists: true,
			total:  1,
			titles: []string{"title1"},
			bodies: []string{"body1"},
		},
		{
			user:   "foo",
			limit:  2,
			exists: true,
			total:  2,
			titles: []string{"title1", "title2"},
			bodies: []string{"body1", "body2"},
		},
		{
			user:   "foo",
			limit:  2,
			offset: 1,
			exists: true,
			total:  2,
			titles: []string{"title2", "title3"},
			bodies: []string{"body2", "body3"},
		},
		{
			user:   "foo",
			limit:  2,
			offset: 3,
			exists: true,
			total:  1,
			titles: []string{"title4"},
			bodies: []string{"body4"},
		},
		{
			user:   "foo",
			limit:  2,
			offset: 4,
			exists: true,
			total:  0,
			titles: []string{},
			bodies: []string{},
		},
		{user: "none", exists: false},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			h := todoHandler{}

			req := &db.TodoListRequest{tc.user, tc.limit, tc.offset}
			resp := &db.TodoListResponse{}

			if err := h.List(context.TODO(), req, resp); err != nil {
				t.Fatal(err)
			}

			if resp.Exists != tc.exists {
				t.Fatalf("expected exists %v, got %v", tc.exists, resp.Exists)
			}

			if !tc.exists {
				return
			}

			if len(resp.Entries) != tc.total {
				t.Fatalf("expected %d entries, got %d", tc.total, len(resp.Entries))
			}

			for j, e := range resp.Entries {
				if e.Title != tc.titles[j] {
					t.Fatalf("expected %d title %s, got %s", j, tc.titles[j], e.Title)
				}
				if e.Body != tc.bodies[j] {
					t.Fatalf("expected %d body %s, got %s", j, tc.bodies[j], e.Body)
				}
			}
		})
	}
}
