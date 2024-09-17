package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/nextlag/in_memory/internal/server/storage"
	"github.com/nextlag/in_memory/pkg/logger/l"
)

const (
	responseOk       = "[ok]"
	responseErr      = "[error]"
	responseNotFound = "[not found]"
)

type computeLayer interface {
	Process(queryStr string) (Query, error)
}

type storageLayer interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
}

type UseCase struct {
	computeLayer computeLayer
	storageLayer storageLayer
	response     string
	log          *l.Logger
}

func New(computeLayer computeLayer, storageLayer storageLayer, log *l.Logger) (*UseCase, error) {
	if computeLayer == nil {
		return nil, errors.New("compute is invalid")
	}

	if storageLayer == nil {
		return nil, errors.New("storage is invalid")
	}

	return &UseCase{
		computeLayer: computeLayer,
		storageLayer: storageLayer,
		log:          log,
	}, nil
}

func (uc *UseCase) HandleQuery(ctx context.Context, queryStr string) string {
	uc.log.Debug("handling query", l.StringAttr("query", queryStr))
	query, err := uc.computeLayer.Process(queryStr)
	if err != nil {
		return fmt.Sprintf("%s %s", responseErr, err.Error())
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

	return fmt.Sprintf("%s internal error", responseErr)
}

func (uc *UseCase) setQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	if err := uc.storageLayer.Set(ctx, arguments[0], arguments[1]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return responseOk
}

func (uc *UseCase) getQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	value, err := uc.storageLayer.Get(ctx, arguments[0])
	if errors.Is(err, storage.ErrNotFound) {
		return responseNotFound
	} else if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return fmt.Sprintf("%s %s", responseOk, value)
}

func (uc *UseCase) delQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	if err := uc.storageLayer.Del(ctx, arguments[0]); err != nil {
		return fmt.Sprintf("%s %s", responseErr, err.Error())
	}

	return responseOk
}
