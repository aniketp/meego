package checker

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/aniketp/meego/ast"
)

//  Driver Type-Checker function
func Checker(program *ast.Program) error {
	env = NewEnvironment()
	_, err := checker(program)
	return err
}

func checker(node ast.Node) (string, error) {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ReturnStatement:
		return evalReturnStatement(node)
	case *ast.IfStatement:
		return evalIfStatement(node)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(node)
	case *ast.AssignStatement:
		return evalAssignStatement(node)
	case *ast.InitStatement:
		return evalInitStatement(node)
	case *ast.FunctionStatement:
		return evalFunctionStatement(node)

	// Expressions
	case *ast.InfixExpression:
		return evalInfixExpression(node)
	case *ast.IntegerLiteral:
		return evalInteger(node)
	case *ast.StringLiteral:
		return evalString(node)
	case *ast.Boolean:
		return evalBoolean(node)
	case *ast.Identifier:
		return evalIdentifier(node)
	case *ast.FunctionCall:
		return evalFunctionCall(node)
	}

	return "", nil
}

// Target program
func evalProgram(p *ast.Program) (string, error) {
	for _, function := range p.Functions {
		_, err := checker(function)
		if err != nil {
			return "", err
		}
	}

	for _, statement := range p.Statements {
		_, err := checker(statement)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

// Statements
func evalBlockStatement(node *ast.BlockStatement) (string, error) {
	for _, statement := range node.Statements {
		result, err := checker(statement)
		if err != nil {
			return "", err
		}

		if reflect.TypeOf(statement) == reflect.TypeOf(&ast.ReturnStatement{}) {
			return result, nil
		}
	}

	return NOTHING_TYPE, nil
}

func evalReturnStatement(node *ast.ReturnStatement) (string, error) {
	res, err := checker(node.ReturnValue)
	return res, err
}

func evalIfStatement(node *ast.IfStatement) (string, error) {
	cond, _ := checker(node.Condition)
	if cond != BOOL_TYPE {
		return "", errors.New("Condition not of Boolean type")
	}

	checker(node.Block)
	checker(node.Alternative)
	return "", nil
}

func evalExpressionStatement(node *ast.ExpressionStatement) (string, error) {
	_, err := checker(node.Expression)
	if err != nil {
		return "", err
	}

	return "", nil
}

func evalInitStatement(node *ast.InitStatement) (string, error) {
	if env.IdentExist(node.Location) {
		return "", errors.New("Identifier already exists")
	}

	right, err := checker(node.Expr)
	if err != nil {
		return "", nil
	}

	env.Set(node.Location, right) // Set identifier type
	return "", nil
}

// TODO: Fix this
func evalAssignStatement(node *ast.AssignStatement) (string, error) {
	right, err := checker(node.Right)
	if err != nil {
		return "", nil
	}

	if kind, ok := env.Get(node.Left.Value); ok {
		if kind != right {
			return "", errors.New("Invalid type assignment")
		}
	} else {
		return "", errors.New("Identifier does not exist")
	}
	return "", nil
}

func evalFunctionStatement(node *ast.FunctionStatement) (string, error) {
	var params []string
	for _, param := range node.Parameters {
		env.Set(param.Arg, param.Type) // set params into scope
		params = append(params, param.Type)
	}

	res, err := checker(node.Body)
	if err != nil {
		return "", err
	}
	// check if correct return type
	if res != node.Return {
		return "", errors.New("Incorrect return type")
	}

	SetFunctionSignature(node.Name, Signature{node.Return, params})
	return "", nil
}

// Expressions
func evalFunctionCall(node *ast.FunctionCall) (string, error) {
	// print() function call
	if IsBuiltin(node.Name) {
		res, err := checker(node.Args[0])
		if err != nil {
			return "", err
		}

		node.Type = res
		return NOTHING_TYPE, nil
	}

	var sig Signature
	// Check if the function called is a valid one
	if sig, ok := GetFunctionSignature(node.Name); !ok {
		return "", errors.New("Function signature does not exist")
	}

	if len(node.Args) != len(sig.Params) {
		return "", errors.New("Incorrect number of function arguments")
	}

	// Validate parameters
	for i, arg := range node.Args {
		res, err := checker(arg)
		if err != nil {
			return "", errors.New(err.Error())
		}

		if res != sig.Params[i] {
			return "", errors.New("Invalid argument type")
		}
	}

	return sig.Return, nil
}

func evalIdentifier(node *ast.Identifier) (string, error) {
	kind, _ := env.Get(node.Value)
	return kind, nil
}

// Trivial
func evalBoolean(node *ast.Boolean) (string, error) {
	return BOOL_TYPE, nil
}

func evalInteger(node *ast.IntegerLiteral) (string, error) {
	return INT_TYPE, nil
}

func evalString(node *ast.StringLiteral) (string, error) {
	return STRING_TYPE, nil
}

// Evaluate the provided expression in infix form
func evalInfixExpression(node *ast.InfixExpression) (string, error) {
	left, err := checker(node.Left)
	if err != nil {
		return left, err
	}

	right, err := checker(node.Right)
	if err != nil {
		return right, err
	}

	if left != right {
		return "", errors.New("Incorrect types for operation")
	}

	// Type setting for code generation
	node.Type = left // (or right)

	// Construct a map for all allowed methods
	methods := map[string]string{
		"+":   PLUS,
		"-":   MINUS,
		"==":  EQUAL,
		"<":   LT,
		">":   GT,
		"*":   TIMES,
		"/":   DIVIDE,
		"or":  OR,
		"and": AND,
	}

	if !MethodExist(right, methods[node.Operator]) {
		return NOTHING_TYPE, errors.New(fmt.Sprintf("Method %s does not exist for type %s",
			methods[node.Operator], left))
	}

	// Locate the node's operator and if found, return BOOL_TYPE
	for _, opr := range []string{"<=", "<", ">=", ">", "or", "and"} {
		if node.Operator == opr {
			return BOOL_TYPE, nil
		}
	}

	return left, nil
}
