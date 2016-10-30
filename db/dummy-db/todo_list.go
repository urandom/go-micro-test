package main

import (
	"golang.org/x/net/context"

	"github.com/urandom/go-micro-test/db"
)

type todoHandler struct{}

func (h *todoHandler) TodoList(ctx context.Context, req *db.TodoListRequest, resp *db.TodoListResponse) error {
	if req.Limit <= 0 {
		resp.Entries, resp.Exists = todoData[req.User]
		return nil
	}

	if all, ok := todoData[req.User]; ok {
		start := req.Offset
		if start <= 0 {
			start = 0
		}

		total := int64(len(all))
		if start+req.Limit < total {
			total = start + req.Limit
		}

		resp.Entries = make([]*db.TodoEntry, 0, req.Limit)
		for i := start; i < total; i++ {
			resp.Entries = append(resp.Entries, all[i])
		}

		resp.Exists = true
	}

	return nil
}

var (
	todoData = map[string][]*db.TodoEntry{
		"foo": []*db.TodoEntry{
			{"title1", "body1"},
			{"title2", "body2"},
			{"title3", "body3"},
			{"title4", "body4"},
		},
		"bar": []*db.TodoEntry{
			{"title1", "body1"},
			{"title2", "body2"},
		},
	}
)
