package network

import (
	"context"
	"errors"

	config "github.com/nextlag/in_memory/config/server"
	"github.com/nextlag/in_memory/pkg/logger/l"
	"github.com/nextlag/in_memory/pkg/parse"
)

const (
	srvTypeCMD = "cmd"
	srvTypeTCP = "tcp"
)

type Handler = func(context.Context, []byte) []byte

type Server interface {
	LaunchServer(ctx context.Context, handler Handler) error
	Close()
}

func New(cfg *config.Config, log *l.Logger) (srv Server, err error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}
	if log == nil {
		return nil, errors.New("logger is nil")
	}

	var (
		opt  []TCPServerOption
		size int
	)

	if cfg.Network.MaxConnections != 0 {
		opt = append(opt, WithServerMaxConnectionsNumber(uint(cfg.Network.MaxConnections)))
	}

	if cfg.Network.MaxMessageSize != "" {
		size, err = parse.Size(cfg.Network.MaxMessageSize)
		if err != nil {
			return nil, errors.New("incorrect max message size")
		}

		opt = append(opt, WithServerBufferSize(uint(size)))
	}

	if cfg.Network.IdleTimeout != 0 {
		opt = append(opt, WithServerIdleTimeout(cfg.Network.IdleTimeout))
	}

	switch cfg.Server.SrvType {
	case srvTypeCMD:
		srv, err = NewCMDServer(log)
		if err != nil {
			return nil, errors.New("failed to create new cmd server")
		}
	case srvTypeTCP:
		srv, err = NewTCPServer(cfg.Network.TCPSocket, log, opt...)
		if err != nil {
			return nil, errors.New("failed to create new tcp server")
		}

	default:
		return nil, errors.New("unknown server type")
	}

	return
}
