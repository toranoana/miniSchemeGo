package main

import (
	"bufio"
	"fmt"
	"miniSchemeGo/eval"
	"miniSchemeGo/lexer"
	"miniSchemeGo/parse"
	"miniSchemeGo/print"
	"os"
)

func main() {
	env := eval.NewEnv()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("miniSchemeGo> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		input := scanner.Text()
		l := lexer.New(input)
		p := parse.NewParser(l.ReadToken())
		exp := p.Parse(env)
		if len(p.Error) == 0 {
			print.Print(env.Eval(exp))
		} else {
			fmt.Println(p.Error)
		}
		fmt.Printf("\n")
	}
}
