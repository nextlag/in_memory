package network

import (
	"errors"
	"github.com/nextlag/in_memory/configuration"
	"github.com/nextlag/in_memory/pkg/logger/l"
)

const srvTypeCMD = "cmd"

func CreateNetwork(cfg *configuration.Config, log *l.Logger) (srv Server, err error) {
	switch cfg.Server.SrvType {
	case srvTypeCMD:
		srv, err = NewServer(log)
		if err != nil {
			return nil, errors.New("failed new server")
		}
	default:
		return nil, errors.New("unknown server type")
	}

	return
}
