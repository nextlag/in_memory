package usecase

import (
	"errors"
)

var (
	errInvalidQuery     = errors.New("empty query")
	errInvalidCommand   = errors.New("invalid command")
	errInvalidArguments = errors.New("invalid arguments")
)

type Analyzer struct{}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(tokens []string) (Query, error) {
	if len(tokens) == 0 {
		return Query{}, errInvalidQuery
	}

	command := tokens[0]
	commandID := commandNameToCommandID(command)
	if commandID == UnknownCommandID {
		return Query{}, errInvalidCommand
	}

	query := NewQuery(commandID, tokens[1:])
	argumentsNum := commandArgumentsNumber(commandID)
	if len(query.Arguments()) != argumentsNum {
		return Query{}, errInvalidArguments
	}

	return query, nil
}
