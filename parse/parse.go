package parse

import (
	. "miniSchemeGo/eval"
	"miniSchemeGo/types"
	"strconv"
)

type Parser struct {
	tokens []types.Token
	index  int
	Error  string
}

func NewParser(tokens []types.Token) *Parser {
	return &Parser{tokens, 0, ""}
}

func (p *Parser) NextToken() {
	if p.index < len(p.tokens) || p.tokens[p.index].Type != types.EOF {
		p.index++
	}
}

func (p *Parser) Token() types.TokenType {
	return p.tokens[p.index].Type
}

func (p *Parser) Literal() string {
	return p.tokens[p.index].Literal
}

func (p *Parser) SetNumber() Expression {
	val, err := strconv.Atoi(p.Literal())
	if err == nil {
		return NewNumber(val)
	}
	return nil
}

func (p *Parser) SetSymbol(env *Env) Expression {
	name := p.Literal()
	if atom, ok := env.GetSysSymbolAtom(name); ok {
		return atom
	} else {
		return env.SetSymbolAtom(name)
	}
}

func (p *Parser) Parse(env *Env) Expression {
	for p.Token() != types.EOF {
		switch p.Token() {
		case types.NUMBER:
			return p.SetNumber()
		case types.SYMBOL:
			return p.SetSymbol(env)
		case types.LPARAM:
			p.NextToken()
			return p.MakeList(env)
		case types.QUOTE:
			quoteSymbol, _ := env.GetSysSymbolAtom(QUOTE)
			nilSymbol, _ := env.GetSysSymbolAtom(NIL)
			p.NextToken()
			return NewCell(quoteSymbol, NewCell(p.Parse(env), nilSymbol))
		}
		p.NextToken()
	}
	nilSymbol, _ := env.GetSysSymbolAtom(NIL)
	return nilSymbol
}

func (p *Parser) MakeList(env *Env) Expression {
	if p.Token() == types.RPARAM {
		return env.SymbolAtom[NIL]
	} else if p.Token() != types.EOF {
		car := p.Parse(env)
		p.NextToken()
		if p.Token() == types.DOT {
			p.NextToken()
			if p.Token() == types.LPARAM {
				p.NextToken()
				cdr := p.MakeList(env)
				return NewCell(car, cdr)
			} else {
				p.NextToken()
				cdr := p.Parse(env)
				return NewCell(car, cdr)
			}
		} else {
			cdr := p.MakeList(env)
			return NewCell(car, cdr)
		}
	} else {
		p.Error = "list Unfinished."
		return NewSymbolAtom(NIL)
	}
}
