package parser

import (
	"fmt"
	"monkeylang/ast"
	"monkeylang/lexer"
	"monkeylang/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

// a pratt parser will create an associations between token types and functions that will parse the token
type Parser struct {
	l *lexer.Lexer
	errors []string

	curToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression // this function takes an expression which is on the left side of the operator
)

func (p * Parser) registerPrefix (tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p * Parser) registerInfix (tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l : l, errors: []string{}}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	p.registerPrefix(token.IDENT, p.parseIdentifier)

	return p
}

func (p * Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType){
	msg := fmt.Sprintf("expected next token error: expected = {%s} | got = {%s}", t, p.peekToken.Type)

	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // root node of AST
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil { //if there exists a statement
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement { // this is a helper method for the ParseProgram method
	switch p.curToken.Type {
	case token.LET :
		return p.parseLetStatement()
	case token.RETURN :
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
	p.nextToken()
	}
	return stmt
	}

func (p * Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	return leftExp
}
func (p *Parser) parseLetStatement() *ast.LetStatement { // this is a helper method for the parseStatement method
	stmt := &ast.LetStatement{Token: p.curToken} // create a new let statement

	if !p.expectPeek(token.IDENT) { // if the next token is not an identifier
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal} // set the name of the let statement

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO : We're skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p * Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO : We're skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}


func (p * Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p * Parser) expectPeek(t token.TokenType) bool {

	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t) // we add error if we do not have the expected token.
		return false
	}
}