package types

type TokenType int

const (
	LPARAM TokenType = iota
	RPARAM
	QUOTE
	DOT
	NUMBER
	SYMBOL

	EOF // 入力の終了
)

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(tokenType TokenType, literal string) Token {
	return Token{tokenType, string(literal)}
}
