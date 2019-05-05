package ast

import "github.com/aniketp/go-compiler/token"

type Attrib interface{}

type Program struct {
	Statements []Statement `json:"statements"`
	Functions  []Statement `json:"functions"`
}

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Statement structures
type AssignStatement struct {
	Token *token.Token `json:"-"`
	Left  Identifier   `json:"left"`
	Right Expression   `json:"right"`
}

type FunctionStatement struct {
	Token      *token.Token    `json:"-"`
	Name       string          `json:"name"`
	Parameters []FormalArg     `json:"params"`
	Body       *BlockStatement `json:"body"`
	Return     string          `json:"return"`
}

type FormalArg struct {
	Arg  string `json:"arg"`
	Type string `json:"type"`
}

type ForStatement struct {
	Token          *token.Token    `json:"-"`
	Condition      Expression      `json:"condition"`
	BlockStatement *BlockStatement `json:"block"`
}

type ReturnStatement struct {
	Token       *token.Token `json:"-"`
	ReturnValue Expression   `json:"return"`
}

type BlockStatement struct {
	Token      *token.Token `json:"-"`
	Statements []Statement  `json:"statements"`
}

type IfStatement struct {
	Token       *token.Token    `json:"-"`
	Condition   Expression      `json:"condition"`
	Block       *BlockStatement `json:"block"`
	Alternative *BlockStatement `json:"alternative"`
}

type ExpressionStatement struct {
	Token      *token.Token `json:"-"`
	Expression Expression   `json:"statement"`
}

type InitStatement struct {
	Token    *token.Token `json:"-"`
	Expr     Expression   `json:"expression"`
	Location string       `json:"location"`
}

// Expression structures
type Identifier struct {
	Token *token.Token `json:"-"`
	Value string       `json:"value"`
}

type Boolean struct {
	Token *token.Token `json:"-"`
	Value bool         `json:"value"`
}

type IntegerLiteral struct {
	Token *token.Token `json:"-"`
	Value string       `json:"value"`
}

type StringLiteral struct {
	Token *token.Token `json:"-"`
	Value string       `json:"value"`
}

type InfixExpression struct {
	Token    *token.Token `json:"-"`
	Type     string       `json:"-"`
	Left     Expression   `json:"left"`
	Right    Expression   `json:"right"`
	Operator string       `json:"operator"`
}

type FunctionCall struct {
	Token *token.Token `json:"-"`
	Name  string       `json:"name"`
	Args  []Expression `json:"args"`
	Type  string       `json:"type"`
}
