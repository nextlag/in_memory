package initialize

import (
	"errors"
	"fmt"
	config "github.com/nextlag/in_memory/configuration"
	"github.com/nextlag/in_memory/internal/server/initialize/network"
	"github.com/nextlag/in_memory/internal/storage"
	"github.com/nextlag/in_memory/internal/usecase"
	"github.com/nextlag/in_memory/pkg/logger/l"
	"sync"
)

type Initialize struct {
	srv network.Server
	uc  *usecase.UseCase
	wg  sync.WaitGroup
	cfg *config.Config
	log *l.Logger
}

func New(cfg *config.Config) (*Initialize, error) {
	if cfg == nil {
		return nil, errors.New("failed to initialize config")
	}
	log := l.NewLogger(cfg)

	compute, err := usecase.NewCompute(log)
	var opt []storage.Option
	engine, err := CreateEngine(cfg.Engine, log)
	if err != nil {
		return nil, err
	}
	store, err := storage.New(engine, log, opt...)

	uc, err := usecase.New(compute, store, log)

	srv, err := network.CreateNetwork(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize network: %w", err)
	}

	return &Initialize{
		srv: srv,
		uc:  uc,
		cfg: cfg,
		log: log,
	}, nil
}
