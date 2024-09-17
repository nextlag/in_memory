package initialize

import (
	"errors"
	"sync"

	"github.com/nextlag/in_memory/configuration"
	"github.com/nextlag/in_memory/internal/server/initialize/network"
	"github.com/nextlag/in_memory/internal/server/storage"
	"github.com/nextlag/in_memory/internal/server/usecase"
	"github.com/nextlag/in_memory/pkg/logger/l"
)

type Initialize struct {
	srv network.Server
	uc  *usecase.UseCase
	wg  sync.WaitGroup
	cfg *configuration.Config
	log *l.Logger
}

func New(cfg *configuration.Config) (i *Initialize, err error) {
	if cfg == nil {
		return nil, errors.New("failed to init config")
	}

	var opt []storage.Option

	log, err := l.NewLogger(cfg)
	if err != nil {
		return
	}

	parser := usecase.NewParser()
	analyzer := usecase.NewAnalyzer()

	compute, err := usecase.NewCompute(parser, analyzer)
	if err != nil {
		return
	}

	engine, err := CreateEngine(cfg.Engine, log)
	if err != nil {
		return
	}

	store, err := storage.New(engine, log, opt...)
	if err != nil {
		return
	}

	uc, err := usecase.New(compute, store, log)
	if err != nil {
		return
	}

	srv, err := network.New(cfg, log)
	if err != nil {
		return
	}

	i = &Initialize{
		srv: srv,
		uc:  uc,
		cfg: cfg,
		log: log,
	}

	return
}
