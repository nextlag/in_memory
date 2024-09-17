package usecase

import (
	"errors"
)

type Compute struct {
	parser   *Parser
	analyzer *Analyzer
}

func NewCompute(parser *Parser, analyzer *Analyzer) (*Compute, error) {
	if parser == nil {
		return nil, errors.New("query parser is invalid")
	}
	if analyzer == nil {
		return nil, errors.New("query analyzer is invalid")
	}

	return &Compute{
		parser:   parser,
		analyzer: analyzer,
	}, nil
}

func (d *Compute) Process(queryStr string) (Query, error) {
	tokens, err := d.parser.Parse(queryStr)
	if err != nil {
		return Query{}, err
	}

	query, err := d.analyzer.Analyze(tokens)
	if err != nil {
		return Query{}, errors.New("error Analyze")
	}

	return query, nil
}