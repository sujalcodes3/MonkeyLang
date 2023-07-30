package ast

import (
	"bytes"
	"monkeylang/token"
)

type Node interface {
	TokenLiteral() string
	String() string // this will print the AST nodes for debugging purposes and to compare them with other AST notes
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program is our first implementation of the Statement interface
// Program is a slice of Statement
type Program struct { // it is also the entry node of our Program
	Statements []Statement
}
func(p * Program) String() string{
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}
func (i *Identifier) String() string {
	return i.Value
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name *Identifier
	Value Expression
}
func (ls *LetStatement) String() string{
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}
func (ls *LetStatement) statementNode(){}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type ReturnStatement struct {
	Token token.Token // the token.RETURN token
	ReturnValue Expression
}
func (rs * ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + "")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
func (rs * ReturnStatement) statementNode(){}
func (rs * ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

type ExpressionStatement struct { // the "statement" keyword only indicates that this is just a wrapper for a standalone expression 'cause our language supports expression statements
	Token token.Token
	Expression Expression
}
func (es * ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}
func (es *ExpressionStatement) statementNode(){}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

