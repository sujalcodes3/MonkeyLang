package parser

import (
	"fmt"
	"monkeylang/ast"
	"monkeylang/lexer"
	"monkeylang/token"
	"strconv"
)

const (
	// increases one by one
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precendences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
    token.LPAREN:   CALL,
}

// a pratt parser will create an associations between token types and functions that will parse the token
type Parser struct {
	l      *lexer.Lexer // takes lexer as an input
	errors []string

    // reading through the file.
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression // this function takes an expression which is on the left side of the operator
)

func (p *Parser) peekPrecedence() int {
	if p, ok := precendences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// we get the prefix function for the current token
	prefix := p.prefixParseFns[p.curToken.Type]

	// if the prefix function is nil, then we have an error
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// we call the prefix function
	leftExp := prefix()

	// we loop until we encounter a semicolon or the precedence of the next token is lower than the precedence of the current token
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) curPrecedence() int {
	if p, ok := precendences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	// = something
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	// 5
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// !something
	p.registerPrefix(token.BANG, p.parsePrefixExpression)

	// -something
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

    // TRUE
    p.registerPrefix(token.TRUE, p.parseBoolean)

    // FALSE
    p.registerPrefix(token.FALSE, p.parseBoolean)

    // grouped expression
    p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

    // if statement
    p.registerPrefix(token.IF, p.parseIfExpression)

    // functions
    p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	// something + something
	p.registerInfix(token.PLUS, p.parseInfixExpression)

	// something - something
	p.registerInfix(token.MINUS, p.parseInfixExpression)

	// something / something
	p.registerInfix(token.SLASH, p.parseInfixExpression)

	// something * something
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)

	// something == something
	p.registerInfix(token.EQ, p.parseInfixExpression)

	// something != something
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)

	// something < something
	p.registerInfix(token.LT, p.parseInfixExpression)

	// something > something
	p.registerInfix(token.GT, p.parseInfixExpression)

    // call expressions
    p.registerInfix(token.LPAREN, p.parseCallExpression)

	return p
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// The function of the Parser.
func (p *Parser) ParseProgram() *ast.Program {

	// creating an instance of the program type.
	program := &ast.Program{} // root node of AST
	program.Statements = []ast.Statement{}

	// we are parsing the program until we reach the end of the file.
	for p.curToken.Type != token.EOF {

		// we are parsing the statements and appending them to the program.
		stmt := p.parseStatement()
		// NOTE : we removed an if check here which checked if stmt is nil
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement { // this is a helper method for the ParseProgram method
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
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
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value

	return lit
}

// parseBoolean is a helper function for the parseExpression method
func (p *Parser) parseBoolean() ast.Expression {
    return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)} 
}


// parseGroupedExpression is for grouped Expression
func (p * Parser) parseGroupedExpression () ast.Expression {
    p.nextToken()

    exp := p.parseExpression(LOWEST)

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return exp
}



func (p *Parser) parseLetStatement() *ast.LetStatement { // this is a helper method for the parseStatement method

	stmt := &ast.LetStatement{Token: p.curToken} // create a new let statement

	// ? The expectPeek method also moves the pointer ahead - keep in mind
	if !p.expectPeek(token.IDENT) { // if the next token is not an identifier, then something is wrong in the program
		return nil
	}

	// the name of a statement is an IDENTIFIER.
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal} // set the name of the let statement

	// the next token should be an ASSIGN token, it moves the pointer ahead and checks too, the idea of the expectPeek method is very good.
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO : We're skipping the expressions until we encounter a semicolon

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO : We're skipping the expression after the 'return' keyword until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {

	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t) // we add error if we do not have the expected token.
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token error: expected = {%s} | got = {%s}", t, p.peekToken.Type)

	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}



// parsing of if-else statements
func (p *Parser) parseIfExpression() ast.Expression {
    expression := &ast.IfExpression{Token: p.curToken}
    if !p.expectPeek(token.LPAREN) {
        return nil
    }
    p.nextToken()
    expression.Condition = p.parseExpression(LOWEST)
    if !p.expectPeek(token.RPAREN) {
        return nil
    }
    if !p.expectPeek(token.LBRACE) {
        return nil
    }
    expression.Consequence = p.parseBlockStatement()

    // checks for an existing else statement but does not throw an error if there is none
    if p.peekTokenIs(token.ELSE) {
        p.nextToken()
        
        if !p.expectPeek(token.LBRACE) {
            return nil
        }

        expression.Alternative = p.parseBlockStatement()
    }

    return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
    block := &ast.BlockStatement{Token: p.curToken}
    block.Statements = []ast.Statement{}
    p.nextToken()
    for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
        stmt := p.parseStatement()
        if stmt != nil {
            block.Statements = append(block.Statements, stmt)
        }
        p.nextToken()
    }
    return block
}


// parsing of functions

func (p * Parser) parseFunctionLiteral () ast.Expression {
    lit := &ast.FunctionLiteral{Token: p.curToken}

    if !p.expectPeek(token.LPAREN) {
        return nil
    }

    lit.Parameters = p.parseFunctionParameters()

    if !p.expectPeek(token.LBRACE) {
        return nil
    }

    lit.Body = p.parseBlockStatement()

    return lit
}

func (p * Parser) parseFunctionParameters () []*ast.Identifier {
    identifiers := []*ast.Identifier{}

    if p.peekTokenIs(token.RPAREN) {
        p.nextToken()
        return identifiers
    }

    p.nextToken()

    ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
    
    identifiers = append(identifiers, ident)
    
    for p.peekTokenIs(token.COMMA) {
        p.nextToken()
        p.nextToken()

        ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

        identifiers = append(identifiers, ident)
    }
    
    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return identifiers
}


// call expressions

func (p * Parser) parseCallExpression(function ast.Expression) ast.Expression {
    exp := &ast.CallExpression{Token: p.curToken, Function: function}

    exp.Arguments = p.parseCallArguments()
    return exp;
}

func (p * Parser) parseCallArguments() []ast.Expression {
    args := []ast.Expression{}

    if p.peekTokenIs(token.RPAREN) {
        p.nextToken()
        return args
    }

    p.nextToken()

    args = append(args, p.parseExpression(LOWEST))

    for p.peekTokenIs(token.COMMA) {
        p.nextToken()
        p.nextToken()

        args = append(args, p.parseExpression(LOWEST))
    } 

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return args
}
