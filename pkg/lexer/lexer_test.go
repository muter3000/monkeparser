// Description: This file contains tests for the lexer package.
package lexer_test

import (
	"monkeparser/pkg/lexer"
	"monkeparser/pkg/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	code := `
let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
`
	expected := []token.Token{
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},

		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "ten"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},

		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "add"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.FUNCTION, Literal: "fn"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},

		{Type: token.IDENT, Literal: "x"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.SEMICOLON, Literal: ";"},

		{Type: token.RBRACE, Literal: "}"},
		{Type: token.SEMICOLON, Literal: ";"},

		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "result"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.IDENT, Literal: "add"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENT, Literal: "ten"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},

		{Type: token.EOF, Literal: "\x00"},
	}

	l := lexer.New(code)
	for _, e := range expected {
		tok := l.NextToken()
		assert.Equal(t, e, tok)
	}
}

func TestBasicNextToken(t *testing.T) {
	code := "=+(){},;"
	expected := []token.Token{
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.SEMICOLON, Literal: ";"},
	}

	l := lexer.New(code)
	for _, e := range expected {
		tok := l.NextToken()
		assert.Equal(t, e, tok)
	}
}

func TestNextTokenOperators(t *testing.T) {
	code := `
		1==1
		1!=1
		1<1
		1>1
		1<=1
		1>=1
		1/1
		1*1
		1-1
		1+1
		!
	`
	expected := []token.Token{
		{Type: token.INT, Literal: "1"},
		{Type: token.EQ, Literal: "=="},
		{Type: token.INT, Literal: "1"},

		{Type: token.INT, Literal: "1"},
		{Type: token.NOT_EQ, Literal: "!="},
		{Type: token.INT, Literal: "1"},

		{Type: token.INT, Literal: "1"},
		{Type: token.LT, Literal: "<"},
		{Type: token.INT, Literal: "1"},

		{Type: token.INT, Literal: "1"},
		{Type: token.GT, Literal: ">"},
		{Type: token.INT, Literal: "1"},

		{Type: token.INT, Literal: "1"},
		{Type: token.LTE, Literal: "<="},
		{Type: token.INT, Literal: "1"},

		{Type: token.INT, Literal: "1"},
		{Type: token.GTE, Literal: ">="},
		{Type: token.INT, Literal: "1"},

		{Type: token.INT, Literal: "1"},
		{Type: token.DIV, Literal: "/"},
		{Type: token.INT, Literal: "1"},

		{Type: token.INT, Literal: "1"},
		{Type: token.MUL, Literal: "*"},
		{Type: token.INT, Literal: "1"},

		{Type: token.INT, Literal: "1"},
		{Type: token.SUB, Literal: "-"},
		{Type: token.INT, Literal: "1"},

		{Type: token.INT, Literal: "1"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.INT, Literal: "1"},

		{Type: token.BANG, Literal: "!"},

		{Type: token.EOF, Literal: "\x00"},
	}

	l := lexer.New(code)
	for _, e := range expected {
		tok := l.NextToken()
		assert.Equal(t, e, tok)
	}
}

func TestNextTokenKeywords(t *testing.T) {
	code := `
		fn
		let
		true
		false
		if
		else
		return
	`
	expected := []token.Token{
		{Type: token.FUNCTION, Literal: "fn"},
		{Type: token.LET, Literal: "let"},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.IF, Literal: "if"},
		{Type: token.ELSE, Literal: "else"},
		{Type: token.RETURN, Literal: "return"},
		{Type: token.EOF, Literal: "\x00"},
	}

	l := lexer.New(code)
	for _, e := range expected {
		tok := l.NextToken()
		assert.Equal(t, tok, e)
	}
}
