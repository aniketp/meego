package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aniketp/meego/ast"
	"github.com/aniketp/meego/checker"
	"github.com/aniketp/meego/lexer"
	"github.com/aniketp/meego/parser"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Parse(input string) *ast.Program {
	l := lexer.NewLexer([]byte(input))
	p := parser.NewParser()

	node, err := p.Parse(l)
	check(err)
	program, _ := node.(*ast.Program)
	return program
}

func TypeCheck(program *ast.Program) {
	err := checker.Checker(program)
	check(err)
}

func main() {
	if len(os.Args) < 2 {
		panic("Provide a valid file name")
	}

	path := os.Args[1]
	absPath, _ := filepath.Abs(path)
	input, err := ioutil.ReadFile(absPath)
	check(err)

	program := Parse(string(input))
	TypeCheck(program)
}
