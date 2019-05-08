package codegen

import (
	"bytes"
	"strings"

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

func genExpressionStatement(node *ast.ExpressionStatement, b *bytes.Buffer) string {
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
func genFunctionStatement(node *ast.FunctionStatement, b *bytes.Buffer) string {
	if checker.IsBuiltin(node.Name) {
		panic("Already a builtin function")
	}

	write(b, "%s %s(", node.Return, node.Name)

	// Generate all the arguments
	for i, arg := range node.Parameters {
		write(b, "%s %s", arg.Type, arg.Arg)
		// Check if it's the last argument
		if i != len(node.Parameters)-1 {
			write(b, ",") // Write a comma
		}
	}

	// Generate function body
	write(b, ") {\n")
	codeGen(node.Body, b)
	write(b, "}\n\n")
	return ""
}

func genIfStatement(node *ast.IfStatement, b *bytes.Buffer) string {
	cond := codeGen(node.Condition, b)
	write(b, "if (\"true\" == %s.val) {\n", cond)
	codeGen(node.Block, b)
	write(b, "} else {\n")
	codeGen(node.Alternative, b)
	write(b, "}\n\n")
	return ""
}

// Generate integers, strings and booleans
func genInteger(node *ast.IntegerLiteral, b *bytes.Buffer) string {
	tmpVar := freshTemp()
	write(b, "Int %s = Int(%s);\n", tmpVar, string(node.Token.Lit))
	return tmpVar
}

func genString(node *ast.IntegerLiteral, b *bytes.Buffer) string {
	tmpVar := freshTemp()
	str := string(node.Token.Lit)
	str = strings.Replace(str, `\`, "\\", -1)

	// Write the sanitized buffer (Type casted to string class)
	write(b, "String %s = String(%s);\n", tmpVar, str)
	return tmpVar
}

func genBoolean(node *ast.Boolean, b *bytes.Buffer) string {
	// Node can be either of the boolean values
	if node.Value {
		return "Bool(\"true\")"
	} else {
		return "Bool(\"false\")"
	}
	return ""
}

func genIdentifier(node *ast.Identifier, b *bytes.Buffer) string {
	return node.Value
}

// Evaluate function call and infix expression
func genInfixExpression(node *ast.InfixExpression, b *bytes.Buffer) string {
	left := codeGen(node.Left, b)
	right := codeGen(node.Right, b)
	kind := node.Type

	tmp := freshTemp()
	methods := map[string]string{"+": checker.PLUS, "-": checker.MINUS, "==": checker.EQUAL,
		"<": checker.LT, ">": checker.GT, "*": checker.TIMES, "/": checker.DIVIDE, "or": checker.OR, "and": checker.AND}

	method, _ := checker.GetMethod(kind, methods[node.Operator])
	write(b, "%s %s = %s.%s(%s);\n", method.Return, tmp, left, methods[node.Operator], right)
	return tmp
}

func genFunctionCall(node *ast.FunctionCall, b *bytes.Buffer) string {
	var sig checker.Signature
	args := make([]string, len(node.Args))
	// store expression tmp vars
	for i, arg := range node.Args {
		res := codeGen(arg, b)
		args[i] = res
	}

	tmp := freshTemp()
	if checker.IsBuiltin(node.Name) {
		sig, ok := checker.GetMethod(node.Type, node.Name)
		if !ok {
			panic("No builtin function")
		}

		write(b, "%s %s = %s.%s(", sig.Return, tmp, args[0], node.Name)
	} else {
		sig, _ = checker.GetFunctionSignature(node.Name)
		write(b, "%s %s = %s(", sig.Return, tmp, node.Name)
		for i, arg := range args {
			write(b, arg)
			if i != len(args)-1 {
				write(b, ",")
			}
		}
	}

	write(b, ");")
	return tmp
}
