package ast

import (
	"fmt"

	"github.com/aniketp/go-compiler/token"
)

// Interface methods
func (p Program) TokenLiteral() string {
	return "Program"
}

// Statements (and TokenLiterals)
func (ls AssignStatement) statementNode() {}
func (ls AssignStatement) TokenLiteral() string {
	return "AssignStatement"
}

func (rs ReturnStatement) statementNode() {}
func (rs ReturnStatement) TokenLiteral() string {
	return "ReturnStatement"
}

func (es ExpressionStatement) statementNode() {}
func (es ExpressionStatement) TokenLiteral() string {
	return "ExpressionStatement"
}

func (is IfStatement) statementNode() {}
func (is IfStatement) TokenLiteral() string {
	return "IfStatement"
}

func (bs BlockStatement) statementNode() {}
func (bs BlockStatement) TokenLiteral() string {
	return "BlockStatement"
}

func (is InitStatement) statementNode() {}
func (is InitStatement) TokenLiteral() string {
	return "InitStatement"
}

func (fs FunctionStatement) statementNode() {}
func (fs FunctionStatement) TokenLiteral() string {
	return "FunctionStatement"
}

// AST Constructors
func NewProgram(funcs, stmts Attrib) (*Program, error) {
	s, ok := stmts.([]Statement)
	if !ok {
		return nil, Error("NewProgram", "[]Statement", "stmts", stmts)
	}

	f, ok := funcs.([]Statement)
	if !ok {
		return nil, Error("NewProgram", "[]Statement", "funcs", funcs)
	}

	// Combine the functions and statements
	return &Program{Functions: f, Statements: s}, nil
}

func NewStatementList() ([]Statement, error) {
	return []Statement{}, nil
}

func AppendStatement(stmtList, stmt Attrib) ([]Statement, error) {
	s, ok := stmt.(Statement)
	if !ok {
		return nil, Error("AppendStatement", "Statement", "stmt", stmt)
	}
	return append(stmtList.([]Statement), s), nil
}

func NewAssignStatement(left, right Attrib) (Statement, error) {
	l, ok := left.(*token.Token)
	if !ok {
		return nil, Error("NewAssignStatement", "Identifier", "left", left)
	}

	r, ok := right.(*token.Token)
	if !ok {
		return nil, Error("NewAssignStatement", "Expression", "right", right)
	}

	// Return the modified assign statement
	return &AssignStatement{Left: Identifier{Value: string(l.Lit)}, Right: r}, nil
}

func NewExpressionStatement(expr Attrib) (Statement, error) {
	e, ok := expr.(Expression)
	if !ok {
		return nil, Error("NewExpressionStatement", "Expression", "expr", expr)
	}

	return &ExpressionStatement{Expression: e}, nil

}

func NewBlockStatement(stmts Attrib) (*BlockStatement, error) {
	s, ok := stmts.([]Statement)
	if !ok {
		return nil, Error("NewBlockStatement", "[]Statement", "stmts", stmts)
	}

	return &BlockStatement{Statements: s}, nil
}

func NewFunctionStatement(name, args, ret, block Attrib) (Statement, error) {
	n, ok := name.(*token.Token)
	if !ok {
		return nil, Error("NewFunctionStatement", "*token.Token", "name", name)
	}

	b, ok := block.(*BlockStatement)
	if !ok {
		return nil, Error("NewFunctionStatement", "BlockStatement", "block", block)
	}

	a := []FormalArg{}
	if args != nil {
		a, ok = args.([]FormalArg)
		if !ok {
			return nil, Error("NewFunctionStatement", "[]FormalArg", "args", args)
		}
	}

	r, ok := ret.(*token.Token)
	if !ok {
		panic("No return value")
	}

	return &FunctionStatement{Name: string(n.Lit), Body: b, Parameters: a, Return: string(r.Lit)}, nil
}

func NewIfStatement(cond, cons, alt Attrib) (Statement, error) {
	c, ok := cond.(Expression)
	if !ok {
		return nil, fmt.Errorf("invalid type of cond. got=%T", cond)
	}

	cs, ok := cons.(*BlockStatement)
	if !ok {
		return nil, fmt.Errorf("invalid type of cons. got=%T", cons)
	}

	a, ok := alt.(*BlockStatement)
	if !ok {
		return nil, fmt.Errorf("invalid type of alt. got=%T", alt)
	}

	return &IfStatement{Condition: c, Block: cs, Alternative: a}, nil
}
