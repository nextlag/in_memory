package network

import (
	"context"
)

type Server interface {
	LaunchServer(ctx context.Context, handler func(context.Context, []byte, int) string) error
	Close(ctx context.Context) error
}

type ServerOption func(Server)
