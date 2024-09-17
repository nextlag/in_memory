package storage

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

type Storage struct {
	engine    Engine
	generator *IDGenerator
	log       *l.Logger
}

func New(engine Engine, log *l.Logger, options ...Option) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("engine is invalid")
	}

	storage := &Storage{
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

func (s *Storage) Set(ctx context.Context, key, value string) error {
	txID := s.generator.Generate()
	ctx = ContextWithTxID(ctx, txID)

	s.engine.Set(ctx, key, value)
	return nil
}

func (s *Storage) Del(ctx context.Context, key string) error {

	txID := s.generator.Generate()
	ctx = ContextWithTxID(ctx, txID)

	s.engine.Del(ctx, key)
	return nil
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = ContextWithTxID(ctx, txID)

	value, found := s.engine.Get(ctx, key)
	if !found {
		return "", ErrNotFound
	}

	return value, nil
}
