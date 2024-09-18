package repository

import (
	"context"
	"errors"

	"github.com/nextlag/in_memory/pkg/logger/l"
)

var ErrNotFound = errors.New("not found")

type Engine interface {
	Set(context.Context, string, string)
	Get(context.Context, string) (string, bool)
	Del(context.Context, string)
}

type Repository struct {
	engine    Engine
	generator *IDGenerator
	log       *l.Logger
}

func New(engine Engine, log *l.Logger, options ...Option) (*Repository, error) {
	if engine == nil {
		return nil, errors.New("engine is invalid")
	}

	storage := &Repository{
		engine: engine,
		log:    log,
	}

	for _, opt := range options {
		opt(storage)
	}

	var lastLSN int64

	storage.generator = NewIDGenerator(lastLSN)
	return storage, nil
}

func (r *Repository) Set(ctx context.Context, key, value string) error {
	txID := r.generator.Generate()
	ctx = ContextWithTxID(ctx, txID)

	r.engine.Set(ctx, key, value)
	return nil
}

func (r *Repository) Del(ctx context.Context, key string) error {

	txID := r.generator.Generate()
	ctx = ContextWithTxID(ctx, txID)

	r.engine.Del(ctx, key)
	return nil
}

func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	txID := r.generator.Generate()
	ctx = ContextWithTxID(ctx, txID)

	value, found := r.engine.Get(ctx, key)
	if !found {
		return "", ErrNotFound
	}

	return value, nil
}
