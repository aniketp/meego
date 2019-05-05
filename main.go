package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aniketp/go-compiler/lexer"
	"github.com/aniketp/go-compiler/parser"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Parse(input string) {
	l := lexer.NewLexer([]byte(input))
	p := parser.NewParser()

	_, err := p.Parse(l)
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

	Parse(string(input))
}
