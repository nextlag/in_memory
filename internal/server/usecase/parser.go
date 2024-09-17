package usecase

import (
	"errors"
	"fmt"
	"strings"
)

const (
	startState = iota
	symbolCheckState
	whiteSpaceState
)

type Parser struct {
	state int
	sb    strings.Builder
}

func NewParser() *Parser {
	return &Parser{state: startState}
}

func (p *Parser) Parse(query string) (tokens []string, err error) {
	tokens = make([]string, 0, len(query))

	p.sb.Reset()

	for i := range query {
		err = errors.New(fmt.Sprintf("invalid symbol: '%s'", string(query[i])))

		switch p.state {
		case startState:
			if !symbolCheck(query[i]) {
				return
			}
			p.sb.WriteByte(query[i])
			p.state = symbolCheckState
		case symbolCheckState:
			if isSpaceSymbol(query[i]) {
				tokens = append(tokens, p.sb.String())
				p.sb.Reset()
				p.state = whiteSpaceState
				break
			}
			if !symbolCheck(query[i]) {
				return
			}
			p.sb.WriteByte(query[i])
		case whiteSpaceState:
			if isSpaceSymbol(query[i]) {
				continue
			}
			if !symbolCheck(query[i]) {
				return
			}

			p.sb.WriteByte(query[i])
			p.state = symbolCheckState
		}
	}

	if p.state == symbolCheckState {
		tokens = append(tokens, p.sb.String())
	}

	return
}

func isSpaceSymbol(s byte) bool {
	return s == '\t' || s == '\n' || s == ' '
}

func symbolCheck(s byte) bool {
	return s >= 'a' && s <= 'z' ||
		s >= 'A' && s <= 'Z' ||
		s >= '0' && s <= '9' ||
		s == '*' || s == '/' ||
		s == '_' || s == '.'
}
