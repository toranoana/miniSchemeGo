package lexer

import (
	"miniSchemeGo/types"
)

type Lexer struct {
	input        string
	position     int
	nextPosition int
	char         byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) ReadToken() []types.Token {
	var tokens []types.Token
	for {
		tok := l.nextToken()
		tokens = append(tokens, tok)
		if tok.Type == types.EOF {
			break
		}
	}
	return tokens
}

func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.nextPosition]
	}
	l.position = l.nextPosition
	l.nextPosition++
}

func (l *Lexer) nextToken() types.Token {
	var token types.Token

	l.skipWhitespace()
	switch l.char {
	case '(':
		token = types.NewToken(types.LPARAM, string(l.char))
	case ')':
		token = types.NewToken(types.RPARAM, string(l.char))
	case '.':
		token = types.NewToken(types.DOT, string(l.char))
	case '\'':
		token = types.NewToken(types.QUOTE, string(l.char))
	default:
		if isDigit(l.char) {
			return types.NewToken(types.NUMBER, l.readDigit())
		} else if isCharacter(l.char) {
			return types.NewToken(types.SYMBOL, l.readSymbol())
		} else if isOperator(l.char) {
			return types.NewToken(types.SYMBOL, l.readOperator())
		} else {
			return types.NewToken(types.EOF, "")
		}
	}
	l.readChar()
	return token
}

func (l *Lexer) readDigit() string {
	p := l.position
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[p:l.position]
}

func (l *Lexer) readSymbol() string {
	p := l.position
	for isCharacter(l.char) {
		l.readChar()
	}
	return l.input[p:l.position]
}

func (l *Lexer) readOperator() string {
	p := l.position
	for isOperator(l.char) {
		l.readChar()
	}
	return l.input[p:l.position]
}

func isCharacter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_' || isDigit(char) || char == '#'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func isOperator(char byte) bool {
	switch char {
	case '+', '-', '*', '/', '=', '<', '>':
		return true
	}
	return false
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}
