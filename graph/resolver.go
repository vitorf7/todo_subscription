package graph

import "log/slog"

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	Logger *slog.Logger
}
