package eval

type Expression interface{}

type DataType int

type Subr func(arg Expression, env *Env) Expression
type FSubr func(arg Expression, env *Env) Expression

type Cell struct {
	Car Expression
	Cdr Expression
}

func NewCell(car Expression, cdr Expression) *Cell {
	return &Cell{Car: car, Cdr: cdr}
}

type SymbolAtom struct {
	Name string
	Bind Expression
}

type SubrSymbolAtom struct {
	SymbolAtom
	Subr Subr
}

type FSubrSymbolAtom struct {
	SymbolAtom
	FSubr FSubr
}

func NewSymbolAtom(name string) *SymbolAtom {
	atom := &SymbolAtom{Name: name}
	atom.Bind = atom
	return atom
}

func NewSubrSymbolAtom(name string, subr Subr) *SubrSymbolAtom {
	return &SubrSymbolAtom{
		SymbolAtom{name, nil},
		subr,
	}
}

func NewFSubrSymbolAtom(name string, fsubr FSubr) *FSubrSymbolAtom {
	return &FSubrSymbolAtom{
		SymbolAtom{name, nil},
		fsubr,
	}
}

type NumberAtom struct {
	Value int
}

func NewNumber(num int) *NumberAtom {
	return &NumberAtom{Value: num}
}

type ErrorStatement struct {
	Message string
}

func NewError(detail string) *ErrorStatement {
	return &ErrorStatement{Message: detail}
}
