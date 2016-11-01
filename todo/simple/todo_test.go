package main

import (
	"fmt"
	"testing"

	"github.com/micro/go-micro/client"
	"github.com/pkg/errors"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
	"github.com/urandom/go-micro-test/todo"
	"golang.org/x/net/context"
)

func TestTodo(t *testing.T) {
	cases := []struct {
		user     string
		entries  []*db.TodoEntry
		dbErr    error
		tokenErr error
		valid    bool
	}{
		{
			user:    "foo",
			entries: []*db.TodoEntry{{"title1", "body1"}, {"title2", "body2"}},
			valid:   true,
		},
		{
			user:  "dbNotExist",
			valid: true,
		},
		{
			user:  "dbError",
			dbErr: errors.New("dbErr"),
		},
		{
			user:     "error",
			tokenErr: errors.New("dbErr"),
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			h := todoHandler{
				dbTodoClient{tc.user, tc.entries, tc.dbErr},
				tokenClient{tc.user, tc.tokenErr},
			}

			req := &todo.ListRequest{Auth: tc.user}
			resp := &todo.ListResponse{}

			if err := h.List(context.TODO(), req, resp); err != nil {
				cause := errors.Cause(err)
				if tc.dbErr != nil && tc.dbErr == cause {
					return
				} else if tc.tokenErr != nil && tc.tokenErr == cause {
					return
				}

				t.Fatal(err)
			}

			if resp.Valid != tc.valid {
				t.Fatalf("expected valid %v, got %v", tc.valid, resp.Valid)
			}

			if !tc.valid {
				return
			}

			if len(resp.Items) != len(tc.entries) {
				t.Fatalf("expected %d items, got %d", len(tc.entries), len(resp.Items))
			}

			for i, item := range resp.Items {
				if item.Title != tc.entries[i].Title {
					t.Fatalf("expected title %s, got %s", tc.entries[i].Title, item.Title)
				}

				if item.Body != tc.entries[i].Body {
					t.Fatalf("expected body %s, got %s", tc.entries[i].Body, item.Body)
				}
			}
		})
	}
}

type dbTodoClient struct {
	user    string
	entries []*db.TodoEntry
	err     error
}
type tokenClient struct {
	user string
	err  error
}

func (c dbTodoClient) List(ctx context.Context, in *db.TodoListRequest, opts ...client.CallOption) (*db.TodoListResponse, error) {
	if in.User == "dbError" {
		return nil, c.err
	} else if in.User == "dbNotExist" {
		return &db.TodoListResponse{}, nil
	} else {
		return &db.TodoListResponse{true, c.entries}, nil
	}
}

func (c tokenClient) Generate(ctx context.Context, in *auth.TokenGenerateRequest, opts ...client.CallOption) (*auth.TokenGenerateResponse, error) {
	panic("not implemented")
}

func (c tokenClient) Check(ctx context.Context, in *auth.TokenCheckRequest, opts ...client.CallOption) (*auth.TokenCheckResponse, error) {
	if in.Auth == "error" {
		return nil, c.err
	} else if in.Auth == "expired" {
		return &auth.TokenCheckResponse{Expired: true}, nil
	} else {
		return &auth.TokenCheckResponse{Valid: true, User: c.user}, nil
	}
}
