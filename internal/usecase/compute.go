package usecase

import (
	"errors"
	"github.com/nextlag/in_memory/pkg/logger/l"
	"strings"
)

var (
	errInvalidQuery     = errors.New("empty query")
	errInvalidCommand   = errors.New("invalid command")
	errInvalidArguments = errors.New("invalid arguments")
)

type Compute struct {
	logger *l.Logger
}

func NewCompute(logger *l.Logger) (*Compute, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Compute{
		logger: logger,
	}, nil
}

func (d *Compute) Parse(queryStr string) (Query, error) {
	tokens := strings.Fields(queryStr)
	if len(tokens) == 0 {
		d.logger.Debug("empty tokens", l.StringAttr("query", queryStr))
		return Query{}, errInvalidQuery
	}

	command := tokens[0]
	commandID := commandNameToCommandID(command)
	if commandID == UnknownCommandID {
		d.logger.Debug("invalid command", l.StringAttr("query", queryStr))
		return Query{}, errInvalidCommand
	}

	query := NewQuery(commandID, tokens[1:])
	argumentsNumber := commandArgumentsNumber(commandID)
	if len(query.Arguments()) != argumentsNumber {
		d.logger.Debug("invalid arguments for query", l.StringAttr("query", queryStr))
		return Query{}, errInvalidArguments
	}

	return query, nil
}
