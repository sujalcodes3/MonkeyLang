package parser

import (
	"monkeylang/ast"
	"monkeylang/lexer"
	"monkeylang/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l : l}

	p.nextToken()
	p.nextToken()

	return p
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
	default:
		return nil
	}
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

func (p * Parser) curTokenIs(t token.TokenType) bool{
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p * Parser) expectPeek( t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}