package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/nextlag/in_memory/internal/storage"
	"github.com/nextlag/in_memory/pkg/logger/l"
)

type computeLayer interface {
	Parse(string) (Query, error)
}

type storageLayer interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
}

type UseCase struct {
	computeLayer computeLayer
	storageLayer storageLayer
	log          *l.Logger
}

func New(computeLayer computeLayer, storageLayer storageLayer, log *l.Logger) (*UseCase, error) {
	if computeLayer == nil {
		return nil, errors.New("compute is invalid")
	}

	if storageLayer == nil {
		return nil, errors.New("storage is invalid")
	}

	if log == nil {
		return nil, errors.New("logger is invalid")
	}

	return &UseCase{
		computeLayer: computeLayer,
		storageLayer: storageLayer,
		log:          log,
	}, nil
}

func (uc *UseCase) HandleQuery(ctx context.Context, queryStr string) string {
	uc.log.Debug("handling query", l.StringAttr("query", queryStr))
	query, err := uc.computeLayer.Parse(queryStr)
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	switch query.CommandID() {
	case SetCommandID:
		return uc.setQuery(ctx, query)
	case GetCommandID:
		return uc.getQuery(ctx, query)
	case DelCommandID:
		return uc.delQuery(ctx, query)
	default:
		uc.log.Error("Compute layer is incorrect", l.IntAttr("command_id", query.CommandID()))
	}

	return "[error] internal error"
}

func (uc *UseCase) setQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	if err := uc.storageLayer.Set(ctx, arguments[0], arguments[1]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return "[ok]"
}

func (uc *UseCase) getQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	value, err := uc.storageLayer.Get(ctx, arguments[0])
	if errors.Is(err, storage.ErrorNotFound) {
		return "[not found]"
	} else if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return fmt.Sprintf("[ok] %s", value)
}

func (uc *UseCase) delQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	if err := uc.storageLayer.Del(ctx, arguments[0]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return "[ok]"
}
