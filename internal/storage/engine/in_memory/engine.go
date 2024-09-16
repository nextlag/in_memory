package in_memory

import (
	"context"
	"errors"
	"github.com/nextlag/in_memory/pkg/logger/l"
	"github.com/nextlag/in_memory/pkg/tools"
	"hash/fnv"
)

type Engine struct {
	partitions []*HashTable
	log        *l.Logger
}

func NewEngine(log *l.Logger, options ...EngineOption) (*Engine, error) {
	if log == nil {
		return nil, errors.New("logger is invalid")
	}

	engine := &Engine{
		log: log,
	}

	for _, opt := range options {
		opt(engine)
	}

	if len(engine.partitions) == 0 {
		engine.partitions = make([]*HashTable, 1)
		engine.partitions[0] = NewHashTable()
	}

	return engine, nil
}

func (e *Engine) Set(ctx context.Context, key, value string) {
	partitionIdx := 0
	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	partition.Set(key, value)

	txID := tools.GetTxIDFromContext(ctx)
	e.log.Debug("successfull set query", "tx", txID)
}

func (e *Engine) Get(ctx context.Context, key string) (string, bool) {
	partitionIdx := 0
	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	value, found := partition.Get(key)

	txID := tools.GetTxIDFromContext(ctx)
	e.log.Debug("successfully get query", "tx", txID)
	return value, found
}

func (e *Engine) Del(ctx context.Context, key string) {
	partitionIdx := 0
	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	partition.Del(key)

	txID := tools.GetTxIDFromContext(ctx)
	e.log.Debug("successfull del query", "tx", txID)
}

func (e *Engine) partitionIdx(key string) int {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))
	return int(hash.Sum32()) % len(e.partitions)
}