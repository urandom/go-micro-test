package main

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
	"github.com/urandom/go-micro-test/todo"
)

type todoHandler struct {
	dbTodo db.TodoClient
	token  auth.TokenClient
}

func (h *todoHandler) List(ctx context.Context, req *todo.ListRequest, resp *todo.ListResponse) error {
	userResp, err := h.token.Check(ctx, &auth.TokenCheckRequest{req.Auth})
	if err != nil {
		return errors.Wrap(err, "getting user")
	}

	if !userResp.Valid {
		return nil
	}

	todoResp, err := h.dbTodo.List(ctx, &db.TodoListRequest{userResp.User, req.Limit, req.Offset})
	if err != nil {
		return errors.Wrap(err, "getting todo items from db")
	}
	resp.Valid = true

	if todoResp.Exists {
		resp.Items = make([]*todo.Item, 0, len(todoResp.Entries))
		for _, e := range todoResp.Entries {
			resp.Items = append(resp.Items, &todo.Item{Title: e.Title, Body: e.Body})
		}
	}

	return nil
}
