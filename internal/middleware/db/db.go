package db

import (
	"net/http"
)

var (
	userCtxKey = &contextKey{"user"}
)

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return next
	}
}
