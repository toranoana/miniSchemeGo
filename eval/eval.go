package eval

const (
	DEFINE = "define"
	LAMBDA = "lambda"
	QUOTE  = "quote"
	TRUE   = "#t"
	FALSE  = "#f"
	NIL    = "nil"
	EOF    = "EOF"
)

type Stack struct {
	mem []map[string]Expression
}

func NewStack() *Stack {
	stack := &Stack{}
	stack.mem = []map[string]Expression{}
	return stack
}

func (s *Stack) Push(exp map[string]Expression) {
	s.mem = append(s.mem, exp)
}

func (s *Stack) Pop() map[string]Expression {
	exp := s.mem[len(s.mem)-1]
	s.mem = s.mem[:len(s.mem)-1]
	return exp
}

func (s *Stack) Top() map[string]Expression {
	if len(s.mem) > 0 {
		return s.mem[len(s.mem)-1]
	} else {
		return nil
	}
}

type Env struct {
	SymbolAtom map[string]Expression
	Stack      *Stack
}

func NewEnv() *Env {
	env := &Env{}
	env.SymbolAtom = map[string]Expression{}

	env.SetSystemSymbolAtom()
	env.SetSubrAtom()
	env.SetFSubrAtom()

	env.Stack = NewStack()
	return env
}

func (e *Env) GetSysSymbolAtom(name string) (Expression, bool) {
	if sym, ok := e.SymbolAtom[name]; ok {
		return sym, ok
	} else {
		return e.SymbolAtom[NIL], ok
	}
}

func (e *Env) GetStackSymbolAtom(name string) (Expression, bool) {
	stack := e.Stack.Top()
	if stack == nil {
		return e.SymbolAtom[NIL], false
	}
	if sym, ok := stack[name]; ok {
		return sym, ok
	}
	return e.SymbolAtom[NIL], false
}

func (e *Env) SetSystemSymbolAtom() {
	e.SetSymbolAtom(TRUE)
	e.SetSymbolAtom(FALSE)
	e.SetSymbolAtom(EOF)
	e.SetSymbolAtom(NIL)
}

func (e *Env) SetSymbolAtom(name string) *SymbolAtom {
	atom := NewSymbolAtom(name)
	e.SymbolAtom[name] = atom
	return atom
}

func (e *Env) SetSubrAtom() {
	e.SymbolAtom["+"] = NewSubrSymbolAtom("+", fplus)
	e.SymbolAtom["-"] = NewSubrSymbolAtom("-", fminus)
	e.SymbolAtom["*"] = NewSubrSymbolAtom("*", fmult)
	e.SymbolAtom[">"] = NewSubrSymbolAtom(">", fgreater)
	e.SymbolAtom[">="] = NewSubrSymbolAtom("=>", fgreaterEq)
	e.SymbolAtom["<"] = NewSubrSymbolAtom("<", fless)
	e.SymbolAtom["<="] = NewSubrSymbolAtom("<=", flessEq)
	e.SymbolAtom["car"] = NewSubrSymbolAtom("car", fcar)
	e.SymbolAtom["cdr"] = NewSubrSymbolAtom("cdr", fcdr)
	e.SymbolAtom["cons"] = NewSubrSymbolAtom("cons", fcons)
}

func (e *Env) SetFSubrAtom() {
	e.SymbolAtom[DEFINE] = NewFSubrSymbolAtom(DEFINE, fdefine)
	e.SymbolAtom[QUOTE] = NewFSubrSymbolAtom(QUOTE, fquote)
	e.SymbolAtom[LAMBDA] = NewFSubrSymbolAtom(LAMBDA, flambda)
	e.SymbolAtom["if"] = NewFSubrSymbolAtom("if", fif)
}

func (e *Env) isNill(atom Expression) bool {
	if atom, ok := atom.(*SymbolAtom); ok {
		return atom.Name == NIL
	}
	return false
}

func (e *Env) Eval(exp Expression) Expression {
	switch exp.(type) {
	case *NumberAtom:
		return exp
	case *SymbolAtom:
		return e.atomEval(exp.(*SymbolAtom))
	case *Cell:
		cell := exp.(*Cell)
		fun := e.Eval(cell.Car)
		if e.isEvalArgs(fun) {
			exp = e.evalList(cell.Cdr)
		} else {
			exp = cell.Cdr
		}
		exp = e.apply(fun, exp)
	}
	return exp
}

func (e *Env) isEvalArgs(exp Expression) bool {
	switch exp.(type) {
	case *SubrSymbolAtom:
		return true
	case *FSubrSymbolAtom:
		return false
	case *Cell:
		if symbol, ok := exp.(*Cell).Car.(*FSubrSymbolAtom); ok {
			if symbol.Name == LAMBDA {
				return true
			}
		}
	}
	return false
}

func (e *Env) evalList(arg Expression) Expression {
	if _, ok := arg.(*Cell); !ok {
		return NewSymbolAtom(NIL)
	}

	cell := NewCell(e.Eval(arg.(*Cell).Car), nil)
	cell.Cdr = e.evalList(arg.(*Cell).Cdr)
	return cell
}

func (e *Env) atomEval(atom *SymbolAtom) Expression {
	if exp, ok := e.GetStackSymbolAtom(atom.Name); ok {
		if atom, ok := exp.(*SymbolAtom); ok {
			return atom.Bind
		}
	}
	if exp, ok := e.GetSysSymbolAtom(atom.Name); ok {
		if atom, ok := exp.(*SymbolAtom); ok {
			return atom.Bind
		}
	}
	return atom.Bind
}

func (e *Env) applyLambda(exp *SymbolAtom) (*Cell, bool) {
	if exp.Name == LAMBDA {
		if lambda, ok := exp.Bind.(*Cell); ok {
			return lambda, ok
		}
	}
	return nil, false
}

func (e *Env) apply(fun Expression, exp Expression) Expression {
	switch fun.(type) {
	case *SubrSymbolAtom:
		return fun.(*SubrSymbolAtom).Subr(exp, e)
	case *FSubrSymbolAtom:
		return fun.(*FSubrSymbolAtom).FSubr(exp, e)
	case *SymbolAtom:
		if lambda, ok := e.applyLambda(fun.(*SymbolAtom)); ok {
			e.Stack.Push(e.bind(lambda.Car.(*Cell), exp))
			if body, ok := lambda.Cdr.(*Cell); ok {
				exp = e.Eval(body.Car)
			}
			e.Stack.Pop()
			return exp
		}
	}
	return NewError("Unbound Atom : " + fun.(*SymbolAtom).Name)
}

func fplus(arg Expression, env *Env) Expression {
	if _, ok := arg.(*Cell); !ok {
		return NewNumber(0)
	}
	cell := arg.(*Cell)
	res := 0
	for {
		if car, ok := cell.Car.(*NumberAtom); ok {
			res += car.Value
		}
		if cdr, ok := cell.Cdr.(*Cell); ok {
			cell = cdr
		} else {
			break
		}
	}
	return NewNumber(res)
}

func fminus(arg Expression, env *Env) Expression {
	if _, ok := arg.(*Cell); !ok {
		return NewNumber(0)
	}
	cell := arg.(*Cell)
	res := 0
	if car, ok := cell.Car.(*NumberAtom); ok {
		res = car.Value
	}
	for {
		if cdr, ok := cell.Cdr.(*Cell); ok {
			cell = cdr
		} else {
			break
		}
		if car, ok := cell.Car.(*NumberAtom); ok {
			res -= car.Value
		}
	}
	return NewNumber(res)
}

func fmult(arg Expression, env *Env) Expression {
	if _, ok := arg.(*Cell); !ok {
		return NewNumber(0)
	}
	cell := arg.(*Cell)
	res := 0
	if car, ok := cell.Car.(*NumberAtom); ok {
		res = car.Value
	}
	for {
		if cdr, ok := cell.Cdr.(*Cell); ok {
			cell = cdr
		} else {
			break
		}
		if car, ok := cell.Car.(*NumberAtom); ok {
			res *= car.Value
		}
	}
	return NewNumber(res)
}

func fgreater(arg Expression, env *Env) Expression {
	return listLoop(arg, greater, env)
}

func fgreaterEq(arg Expression, env *Env) Expression {
	return listLoop(arg, greaterEq, env)
}

func fless(arg Expression, env *Env) Expression {
	return listLoop(arg, less, env)
}

func flessEq(arg Expression, env *Env) Expression {
	return listLoop(arg, lessEq, env)
}

func greater(arg1 *NumberAtom, arg2 *NumberAtom, env *Env) Expression {
	if arg1.Value > arg2.Value {
		trueAtom, _ := env.GetSysSymbolAtom(TRUE)
		return trueAtom
	} else {
		falseAtom, _ := env.GetSysSymbolAtom(FALSE)
		return falseAtom
	}
}

func less(arg1 *NumberAtom, arg2 *NumberAtom, env *Env) Expression {
	if arg1.Value < arg2.Value {
		trueAtom, _ := env.GetSysSymbolAtom(TRUE)
		return trueAtom
	} else {
		falseAtom, _ := env.GetSysSymbolAtom(FALSE)
		return falseAtom
	}
}

func greaterEq(arg1 *NumberAtom, arg2 *NumberAtom, env *Env) Expression {
	if arg1.Value >= arg2.Value {
		trueAtom, _ := env.GetSysSymbolAtom(TRUE)
		return trueAtom
	} else {
		falseAtom, _ := env.GetSysSymbolAtom(FALSE)
		return falseAtom
	}
}

func lessEq(arg1 *NumberAtom, arg2 *NumberAtom, env *Env) Expression {
	if arg1.Value <= arg2.Value {
		trueAtom, _ := env.GetSysSymbolAtom(TRUE)
		return trueAtom
	} else {
		falseAtom, _ := env.GetSysSymbolAtom(FALSE)
		return falseAtom
	}
}

func listLoop(arg Expression, comparison func(arg1 *NumberAtom, arg2 *NumberAtom, env *Env) Expression, env *Env) Expression {
	if _, ok := arg.(*Cell); !ok {
		falseAtom, _ := env.GetSysSymbolAtom(FALSE)
		return falseAtom
	}
	cell := arg.(*Cell)
	var res Expression
	for {
		if car, ok := cell.Car.(*NumberAtom); ok {
			if cdr, ok := cell.Cdr.(*Cell); ok {
				if atom, ok := cdr.Car.(*NumberAtom); ok {
					res = comparison(car, atom, env)
				}
				cell = cdr
			} else {
				break
			}
		}
	}
	return res
}

func fcar(arg Expression, env *Env) Expression {
	if cell, ok := arg.(*Cell); ok {
		if car, ok := cell.Car.(*Cell); ok {
			return car.Car
		}
	}
	return env.SymbolAtom[NIL]
}

func fcdr(arg Expression, env *Env) Expression {
	if cell, ok := arg.(*Cell); ok {
		if car, ok := cell.Car.(*Cell); ok {
			return car.Cdr
		}
	}
	return env.SymbolAtom[NIL]
}

func fcons(arg Expression, env *Env) Expression {
	if cell, ok := arg.(*Cell); ok {
		if cdr, ok := cell.Cdr.(*Cell); ok {
			return cons(cell.Car, cdr.Car)
		}
	}
	return env.SymbolAtom[NIL]
}

func cons(car Expression, cdr Expression) *Cell {
	return NewCell(car, cdr)
}

func fif(arg Expression, env *Env) Expression {
	if cell, ok := arg.(*Cell); ok {
		comp := env.Eval(cell.Car)
		then := cell.Cdr
		if sym, ok := comp.(*SymbolAtom); ok {
			if sym.Name == TRUE {
				if thenCell, ok := then.(*Cell); ok {
					return env.Eval(thenCell.Car)
				}
			} else {
				if thenCell, ok := then.(*Cell).Cdr.(*Cell); ok {
					return env.Eval(thenCell.Car)
				}
			}
		}
	}
	return NewError("Syntax Error. if argument")
}

func fquote(arg Expression, env *Env) Expression {
	if cell, ok := arg.(*Cell); ok {
		return cell.Car
	}
	return env.SymbolAtom[NIL]
}

func flambda(arg Expression, env *Env) Expression {
	if cell, ok := arg.(*Cell); ok {
		if _, ok := cell.Car.(*Cell); !ok {
			return NewError("Syntax Error. Lambda argument")
		}
		sym := NewSymbolAtom(LAMBDA)
		sym.Bind = cell
		return sym
	}
	return NewError("Syntax Error. Lambda argument")
}

func (e *Env) bind(larg *Cell, arg Expression) map[string]Expression {
	lambdaArg := map[string]Expression{}
	for {
		if symbol, ok := larg.Car.(*SymbolAtom); ok {
			if cell, ok := arg.(*Cell); ok {
				bindSymbol := NewSymbolAtom(symbol.Name)
				bindSymbol.Bind = e.Eval(cell.Car)
				lambdaArg[symbol.Name] = bindSymbol

				arg = cell.Cdr
				if newCell, ok := larg.Cdr.(*Cell); ok {
					larg = newCell
				} else {
					break
				}
			}
		} else {
			break
		}
	}
	return lambdaArg
}

func fdefine(arg Expression, env *Env) Expression {
	if cell, ok := arg.(*Cell); ok {
		if atom, ok := cell.Car.(*SymbolAtom); ok {
			if bind, ok := cell.Cdr.(*Cell); ok {
				atom.Bind = env.Eval(bind.Car)
			}
			return atom
		}
	}
	return NewError("duplicate variable")
}
