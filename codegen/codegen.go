package codegen

import (
	"bytes"

	"github.com/aniketp/meego/ast"
	"github.com/aniketp/meego/checker"
)

func GenWrapper(p *ast.Program) bytes.Buffer {
	TMP_COUNT = 0
	var b bytes.Buffer

	// Call the general code-generator method
	codeGen(p, &b)
	return b
}

// Core dispatcher function of code-generator
func codeGen(node ast.Node, b *bytes.Buffer) string {
	switch node := node.(type) {
	// Statement cases
	case *ast.Program:
		return genProgram(node, b)
	case *ast.BlockStatement:
		return genBlockStatement(node, b)
	case *ast.ReturnStatement:
		return genReturnStatement(node, b)
	case *ast.FunctionStatement:
		return genFunctionStatement(node, b)
	case *ast.IfStatement:
		return genIfStatement(node, b)
	case *ast.ExpressionStatement:
		return genExpressionStatement(node, b)
	case *ast.AssignStatement:
		return genAssignStatement(node, b)
	case *ast.InitStatement:
		return genInitStatement(node, b)
	// Expression cases
	case *ast.InfixExpression:
		return genInfixExpression(node, b)
	case *ast.IntegerLiteral:
		return genInteger(node, b)
	case *ast.StringLiteral:
		return genString(node, b)
	case *ast.Boolean:
		return genBoolean(node, b)
	case *ast.Identifier:
		return genIdentifier(node, b)
	case *ast.FunctionCall:
		return genFunctionCall(node, b)
	}

	return ""
}

// Generate main program
func genProgram(node *ast.Program, b *bytes.Buffer) string {
	write(b,
		"#include <iostream>\n#include <string>\ninclude \"Builtins.cpp\"\n\n")

	// We'll generate all functions before main, to ensure function
	// declaration before invocation
	for _, funcs := range node.Functions {
		codeGen(funcs, b)
	}

	write(b, "int main() {\n")
	for _, stmt := range node.Statements {
		codeGen(stmt, b)
	}

	// Here, we're done with the main function
	write(b, "return 0;\n")
	return ""
}

// Statements within main function
func genBlockStatement(node *ast.BlockStatement, b *bytes.Buffer) string {
	for _, stmt := range node.Statements {
		codeGen(stmt, b)
	}
	return ""
}

func genExpressionStatement(node *ast.BlockStatement, b *bytes.Buffer) string {
	expr := codeGen(node.Expression, b)
	write(b, "%s;\n", expr)
	return ""
}

func genAssignStatement(node *ast.AssignStatement, b *bytes.Buffer) string {
	right := codeGen(node.Right, b)
	write(b, "%s = %s;\n", node.Left.Value, right)
	return ""
}

func genInitStatement(node *ast.InitStatement, b *bytes.Buffer) string {
	right := codeGen(node.Expr, b)
	kind, _ := checker.GetIdentType(node.Location)
	write(b, "%s %s = %s;\n", kind, node.Location, right)
	return ""
}

func genReturnStatement(node *ast.ReturnStatement, b *bytes.Buffer) string {
	value := codeGen(node.ReturnValue, b)
	write(b, "return %s;\n", value)
	return ""
}

// Generate function body
