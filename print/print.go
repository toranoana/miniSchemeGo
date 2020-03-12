package print

import (
	"fmt"
	. "miniSchemeGo/eval"
)

func Print(exp Expression) {
	switch exp.(type) {
	case *NumberAtom:
		fmt.Printf("%d", exp.(*NumberAtom).Value)
	case *SymbolAtom:
		fmt.Printf("%s", exp.(*SymbolAtom).Name)
	case *Cell:
		listPrint(exp.(*Cell))
	case *ErrorStatement:
		fmt.Printf("\n[Error] %s\n", exp.(*ErrorStatement).Message)
	}
}

func listPrint(cell *Cell) {
	fmt.Print("(")
	for {
		Print(cell.Car)
		if cdr, ok := cell.Cdr.(*Cell); ok {
			cell = cdr
			fmt.Print(" ")
		} else {
			if sym, ok := cell.Cdr.(*SymbolAtom); ok {
				if sym.Name == NIL {
					break
				} else {
					fmt.Print(" . ")
					Print(cell.Cdr)
					break
				}
			}
		}
	}
	fmt.Print(")")
}
