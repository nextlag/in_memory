package initialize

import (
	"errors"

	"github.com/nextlag/in_memory/configuration"
	"github.com/nextlag/in_memory/internal/server/storage/engine/in_memory"
	"github.com/nextlag/in_memory/pkg/logger/l"
)

func CreateEngine(cfg *configuration.Engine, log *l.Logger) (*in_memory.Engine, error) {
	if log == nil {
		return nil, errors.New("logger is invalid")
	}

	if cfg.Type != "" {
		supportedTypes := map[string]struct{}{
			"in_memory": {},
		}

		if _, found := supportedTypes[cfg.Type]; !found {
			return nil, errors.New("engine type is incorrect")
		}
	}

	var opt []in_memory.EngineOption
	if cfg.PartitionsNumber != 0 {
		opt = append(opt, in_memory.WithPartitions(cfg.PartitionsNumber))
	}

	return in_memory.NewEngine(log, opt...)
}
