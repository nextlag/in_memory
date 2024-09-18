package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/nextlag/in_memory/internal"
	"github.com/nextlag/in_memory/internal/server/usecase/repository"
	"github.com/nextlag/in_memory/pkg/logger/l"
)

type computeLayer interface {
	Process(queryStr string) (Query, error)
}

type repositoryLayer interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
}

type UseCase struct {
	computeLayer computeLayer
	storageLayer repositoryLayer
	log          *l.Logger
}

func New(compute computeLayer, repo repositoryLayer, log *l.Logger) (*UseCase, error) {
	if compute == nil {
		return nil, errors.New("compute is invalid")
	}
	if repo == nil {
		return nil, errors.New("repository is invalid")
	}

	return &UseCase{
		computeLayer: compute,
		storageLayer: repo,
		log:          log,
	}, nil
}

func (uc *UseCase) HandleQuery(ctx context.Context, queryStr string) string {
	uc.log.Debug("handling query", l.StringAttr("query", queryStr))
	query, err := uc.computeLayer.Process(queryStr)
	if err != nil {
		return fmt.Sprintf("%s %s", internal.ResponseErr, err.Error())
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

	return fmt.Sprintf("%s internal error", internal.ResponseErr)
}

func (uc *UseCase) setQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	if err := uc.storageLayer.Set(ctx, arguments[0], arguments[1]); err != nil {
		return fmt.Sprintf("%s %s", internal.ResponseErr, err.Error())
	}

	return internal.ResponseOk
}

func (uc *UseCase) getQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	value, err := uc.storageLayer.Get(ctx, arguments[0])
	if errors.Is(err, repository.ErrNotFound) {
		return internal.ResponseNotFound
	} else if err != nil {
		return fmt.Sprintf("%s %s", internal.ResponseErr, err.Error())
	}

	return fmt.Sprintf("%s %s", internal.ResponseOk, value)
}

func (uc *UseCase) delQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	if err := uc.storageLayer.Del(ctx, arguments[0]); err != nil {
		return fmt.Sprintf("%s %s", internal.ResponseErr, err.Error())
	}

	return internal.ResponseOk
}
