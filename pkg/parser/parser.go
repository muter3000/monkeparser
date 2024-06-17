package parser

import (
	"fmt"
	"github.com/muter3000/monkeparser/pkg/ast"
	"github.com/muter3000/monkeparser/pkg/lexer"
	"github.com/muter3000/monkeparser/pkg/token"
	"strconv"
)

const (
	LOWEST int = iota + 1
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) NextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = map[token.TokenType]prefixParseFn{
		token.IDENT: p.parseIdentifier,
		token.INT:   p.parseIntegerLiteral,
		token.FALSE: p.parseBooleanLiteral,
		token.TRUE:  p.parseBooleanLiteral,
		token.SUB:   p.parsePrefixModifier,
		token.BANG:  p.parsePrefixModifier,

		token.LPAREN: p.parseGroupedExpression,

		token.IF:       p.parseIfExpression,
		token.FUNCTION: p.parseFuncExpression,
	}

	p.infixParseFns = map[token.TokenType]infixParseFn{
		token.EQ:     p.parseInfixExpression,
		token.NOT_EQ: p.parseInfixExpression,
		token.GT:     p.parseInfixExpression,
		token.LT:     p.parseInfixExpression,
		token.GTE:    p.parseInfixExpression,
		token.LTE:    p.parseInfixExpression,
		token.SUB:    p.parseInfixExpression,
		token.PLUS:   p.parseInfixExpression,
		token.MUL:    p.parseInfixExpression,
		token.DIV:    p.parseInfixExpression,

		token.LPAREN: p.parseCallExpression,
	}

	p.NextToken()
	p.NextToken()
	return p
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be '%s', got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	ls := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	ls.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.NextToken()

	ls.Value = p.parseExpression(LOWEST)
	p.NextToken()

	return ls
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	rs := &ast.ReturnStatement{Token: p.curToken}
	p.NextToken()

	if !p.curTokenIs(token.SEMICOLON) {
		rs.ReturnValue = p.parseExpression(LOWEST)
		p.NextToken()
	}

	return rs
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		st := p.parseStatement()
		if st != nil {
			program.Statements = append(program.Statements, st)
		}
		p.NextToken()
	}

	return program
}

var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,
	token.GT:     LESSGREATER,
	token.LT:     LESSGREATER,
	token.GTE:    LESSGREATER,
	token.LTE:    LESSGREATER,
	token.SUB:    SUM,
	token.PLUS:   SUM,
	token.MUL:    PRODUCT,
	token.DIV:    PRODUCT,
	token.LPAREN: CALL,
}

func (p *Parser) peekPrecedence() int {
	val, ok := precedences[p.peekToken.Type]
	if ok {
		return val
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	val, ok := precedences[p.curToken.Type]
	if ok {
		return val
	}
	return LOWEST
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.NextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: val}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parsePrefixModifier() ast.Expression {
	pe := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.NextToken()

	pe.Right = p.parseExpression(PREFIX)

	return pe
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	ie := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.NextToken()

	ie.Right = p.parseExpression(precedence)

	return ie
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.NextToken()

	left := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return left
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.NextToken()

	exp.Predicate = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	exp.Consequence = p.parseBlockStatement()
	if !p.curTokenIs(token.RBRACE) {
		return nil
	}

	if p.peekTokenIs(token.ELSE) {
		p.NextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()

		if !p.curTokenIs(token.RBRACE) {
			return nil
		}
	}

	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	bExp := &ast.BlockStatement{Statements: []ast.Statement{}}
	for !p.peekTokenIs(token.RBRACE) {
		p.NextToken()
		bExp.Statements = append(bExp.Statements, p.parseStatement())
	}
	p.NextToken()
	return bExp
}

func (p *Parser) parseFuncExpression() ast.Expression {
	fExpr := &ast.FunctionLiteral{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fExpr.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fExpr.Body = p.parseBlockStatement()

	return fExpr
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var pars []*ast.Identifier
	p.NextToken()
	for p.peekTokenIs(token.COMMA) {
		pars = append(pars, p.parseIdentifier().(*ast.Identifier))
		p.NextToken()
		p.NextToken()
	}

	if !p.curTokenIs(token.RPAREN) {
		pars = append(pars, p.parseIdentifier().(*ast.Identifier))
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
	}

	return pars
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	cExp := &ast.CallExpression{
		Token:     p.curToken,
		Function:  function,
		Arguments: p.parseCallArguments(),
	}

	return cExp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression
	if p.peekTokenIs(token.RPAREN) {
		p.NextToken()
		return args
	}
	p.NextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.NextToken()
		p.NextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}
