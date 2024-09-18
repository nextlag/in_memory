package initialize

import (
	"errors"
	"sync"

	config "github.com/nextlag/in_memory/config/server"
	"github.com/nextlag/in_memory/internal/server/initialize/network"
	"github.com/nextlag/in_memory/internal/server/usecase"
	"github.com/nextlag/in_memory/internal/server/usecase/repository"
	"github.com/nextlag/in_memory/pkg/logger/l"
	"github.com/nextlag/in_memory/pkg/parse"
)

type Initialize struct {
	srv network.Server
	uc  *usecase.UseCase
	wg  sync.WaitGroup
	cfg *config.Config
	log *l.Logger
}

func New(cfg *config.Config) (i *Initialize, err error) {
	if cfg == nil {
		return nil, errors.New("failed to init config")
	}

	var opt []repository.Option

	log, err := l.NewLogger(cfg)
	if err != nil {
		return
	}
	log.Info("Init config", "server socket", cfg.Network.TCPSocket)

	parser := parse.NewParser()
	analyzer := usecase.NewAnalyzer()

	compute, err := usecase.NewCompute(parser, analyzer)
	if err != nil {
		return
	}

	engine, err := CreateEngine(cfg.Engine, log)
	if err != nil {
		return
	}

	repo, err := repository.New(engine, log, opt...)
	if err != nil {
		return
	}

	uc, err := usecase.New(compute, repo, log)
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
