package checker

import (
	"errors"
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
			return "", errors.New("invalid type assignment")
		}
	} else {
		return "", errors.New("ident not exist")
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
		return "", errors.New("incorrect return type")
	}

	SetFunctionSignature(node.Name, Signature{node.Return, params})
	return "", nil
}
