package ast

import (
	"bytes"
	"fmt"
	"github.com/muter3000/monkeparser/pkg/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

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

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Identifier struct {
	Token token.Token
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
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) String() string {
	if ls.Value != nil {
		return fmt.Sprintf("%s %s = %s;", ls.TokenLiteral(), ls.Name.String(), ls.Value.String())
	}
	return fmt.Sprintf("%s %s;", ls.TokenLiteral(), ls.Name.String())
}

func (ls *LetStatement) expressionNode() {}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) String() string {
	if rs.ReturnValue != nil {
		return fmt.Sprintf("%s %s;", rs.TokenLiteral(), rs.ReturnValue.String())
	}
	return fmt.Sprintf("%s;", rs.TokenLiteral())
}

func (rs *ReturnStatement) expressionNode() {}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

func (b *BooleanLiteral) statementNode() {}

func (b *BooleanLiteral) expressionNode() {}

func (b *BooleanLiteral) TokenLiteral() string {
	return b.Token.Literal
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) statementNode() {}

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Right.String())
}

func (p *PrefixExpression) statementNode() {}

func (p *PrefixExpression) expressionNode() {}

func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Operator, i.Right.String())
}

func (i *InfixExpression) statementNode() {}

func (i *InfixExpression) expressionNode() {}

func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) String() string {
	buf := strings.Builder{}
	buf.WriteString("{ ")
	for _, statement := range b.Statements {
		buf.WriteString(statement.String())
		if buf.String()[buf.Len()-1] != ';' {
			buf.WriteString(";")
		}
	}
	buf.WriteString(" }")
	return buf.String()
}

func (b *BlockStatement) statementNode() {}

func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token
	Predicate   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) String() string {
	buf := strings.Builder{}
	buf.WriteString("if(")
	buf.WriteString(i.Predicate.String())
	buf.WriteString(")")
	buf.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		buf.WriteString("else")
		buf.WriteString(i.Alternative.String())
	}

	buf.WriteString(";")

	return buf.String()
}

func (i *IfExpression) statementNode() {}

func (i *IfExpression) expressionNode() {}

func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())
	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
