package parse

import (
	"errors"
	"fmt"
	"strings"
)

const (
	startState = iota
	letterOrPunctuationState
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

	p.state = startState
	p.sb.Reset()

	for i := range query {
		err = errors.New(fmt.Sprintf("invalid symbol: '%s'", string(query[i])))
		switch p.state {
		case startState:
			if !isLetterOrPunctuationSymbol(query[i]) {
				return
			}
			p.sb.WriteByte(query[i])
			p.state = letterOrPunctuationState
		case letterOrPunctuationState:
			if isSpaceSymbol(query[i]) {
				tokens = append(tokens, p.sb.String())
				p.sb.Reset()
				p.state = whiteSpaceState
				break
			}
			if !isLetterOrPunctuationSymbol(query[i]) {
				return
			}
			p.sb.WriteByte(query[i])
		case whiteSpaceState:
			if isSpaceSymbol(query[i]) {
				continue
			}
			if !isLetterOrPunctuationSymbol(query[i]) {
				return
			}

			p.sb.WriteByte(query[i])
			p.state = letterOrPunctuationState
		}
	}

	if p.state == letterOrPunctuationState {
		tokens = append(tokens, p.sb.String())
	}

	return tokens, nil
}

func isSpaceSymbol(ch byte) bool {
	return ch == '\t' || ch == '\n' || ch == ' '
}

func isLetterOrPunctuationSymbol(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '*' || ch == '/' ||
		ch == '_' || ch == '.'
}
