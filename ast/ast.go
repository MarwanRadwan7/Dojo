package ast

import "dojo/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node

	// dummy method helps by guiding the Go compiler and possibly causing it to throw errors
	// when we use a Statement instead of an Expression and vice versa.
	statementNode()
}

type Expression interface {
	Node
	// dummy method helps by guiding the Go compiler and possibly causing it to throw errors
	// when we use a Statement instead of an Expression and vice versa.
	expressionNode()
}

// Program node is going to be the root node of every AST the parser produces.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
