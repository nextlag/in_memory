package initialize

import (
	"errors"

	config "github.com/nextlag/in_memory/config/server"
	"github.com/nextlag/in_memory/internal/server/usecase/repository/engine/in_memory"
	"github.com/nextlag/in_memory/pkg/logger/l"
)

const engineInMemory = "in_memory"

func CreateEngine(cfg *config.Engine, log *l.Logger) (*in_memory.Engine, error) {
	if cfg == nil {
		return nil, errors.New("config is invalid")
	}
	if log == nil {
		return nil, errors.New("logger is invalid")
	}

	if cfg.Type != "" {
		supportedTypes := map[string]struct{}{
			engineInMemory: {},
		}

		if _, found := supportedTypes[cfg.Type]; !found {
			return nil, errors.New("engine type is incorrect")
		}
	}

	var opt []in_memory.InMemoryOption
	if cfg.PartitionsNumber != 0 {
		opt = append(opt, in_memory.WithPartitions(cfg.PartitionsNumber))
	}

	return in_memory.New(log, opt...)
}
