package checker

import "github.com/aniketp/meego/ast"

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
