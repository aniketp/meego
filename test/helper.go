package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aniketp/meego/src/ast"
	"github.com/aniketp/meego/src/checker"
	"github.com/aniketp/meego/src/codegen"
	"github.com/aniketp/meego/src/lexer"
	"github.com/aniketp/meego/src/parser"
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

func Compile(code bytes.Buffer) string {
	f, err := os.Create("../template/main.cpp")
	check(err)
	defer f.Close()

	f.Write(code.Bytes())

	var out bytes.Buffer
	cmd1 := exec.Command("g++", "-std=c++11", "-o", "apple", "../template/main.cpp",
		"../template/Builtins.cpp")
	cmd1.Stderr = &out
	err = cmd1.Run()

	// Check if output was a valid one
	if len(out.String()) != 0 {
		panic(fmt.Sprintf("error: %s", out.String()))
	}

	// Now, execute the resulting 'apple' binary
	cmd := exec.Command("./apple")
	var outb bytes.Buffer
	cmd.Stdout = &outb
	err = cmd.Run()
	check(err)

	// This is the generated output of our transpiled program
	return outb.String()

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
	// Generate vanilla C++
	code := codegen.GenWrapper(program)
	result := Compile(code)
	fmt.Println(result)
}
