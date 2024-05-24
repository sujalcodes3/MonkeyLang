package lexer

import "monkeylang/token"

// our lexer has the input string(which the source code writter by the user).
type Lexer struct {
	input        string
	position     int  // current position
	readPosition int  // current reading position (after current char)
	ch           byte // current character
}

// basically a constructor of our lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// returns the next token in the input source code.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// skips all types of whitespaces in our source code before giving the next token.
	l.skipWhiteSpace()

	// checks the type of the current character and returns the corresponding token.
	switch l.ch {
	case '=':
		// checks if the next character is '=' and if so, returns the token for '=='.
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)} // this is the equality operator.
		} else { // otherwise, returns the token for '='.
			tok = newToken(token.ASSIGN, l.ch) // this is the assignment operator.
		}
	case '+':
		tok = newToken(token.PLUS, l.ch) // normal stuff
	case '-':
		tok = newToken(token.MINUS, l.ch) // normal stuff
	case '!':
		if l.peekChar() == '=' { // checks for the not equal operator.
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch) // checks for the "bang" operator.
		}
	case '/':
		tok = newToken(token.SLASH, l.ch) // normal stuff
	case '*':
		tok = newToken(token.ASTERISK, l.ch) // normal stuff
	case '<':
		tok = newToken(token.LT, l.ch) // normal stuff
	case '>':
		tok = newToken(token.GT, l.ch) // normal stuff
	case ';':
		tok = newToken(token.SEMICOLON, l.ch) // normal stuff
	case ',':
		tok = newToken(token.COMMA, l.ch) // normal stuff
	case '(':
		tok = newToken(token.LPAREN, l.ch) // normal stuff
	case ')':
		tok = newToken(token.RPAREN, l.ch) // normal stuff
	case '{':
		tok = newToken(token.LBRACE, l.ch) // normal stuff
	case '}':
		tok = newToken(token.RBRACE, l.ch) // normal stuff
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case 0: // end of file
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default: // checks if the character is a letter or a digit.
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch) // illegal stuff.
		}
	}
	l.readChar() // moves to the next character.
	return tok   // returns the token.
}

func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()

		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

// creates a new token with the given type and literal.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// returns the identifier as a string from the input source code.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) { // checks that the character is a letter.
		l.readChar() // advances our position in the input string.
	}
	return l.input[position:l.position] // returns the identifier as a string.
}

// returns the number as string.
func (l *Lexer) readNumber() string {
	position := l.position // save the position of the first digit.
	for isDigit(l.ch) {
		l.readChar() // advance our position in the input string.
	}
	return l.input[position:l.position] // return the number as a string.
}

// readChar advances our position in the input source code string
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // end of file, 0 is ASCII for NULL
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

// used to peek at the next character without advancing our position in the input.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// the following two functions are helper functions to check if a character is a letter or a digit.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// skips all types of whitespaces in our source code.
func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
