package parser_test

import (
	"github.com/muter3000/monkeparser/pkg/ast"
	"github.com/muter3000/monkeparser/pkg/lexer"
	"github.com/muter3000/monkeparser/pkg/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestParser_ReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
		`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	//testStatements := []string{"5", "10", "993322"}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, st := range program.Statements {
		returnSt, ok := st.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("st not *ast.ReturnStatement. got=%T", st)
			continue
		}
		if returnSt.TokenLiteral() != "return" {
			t.Errorf("returnSt.TokenLiteral not 'return', got %q", returnSt.TokenLiteral())
		}
		//if returnSt.ReturnValue.TokenLiteral() != testStatements[i] {
		//	t.Errorf("returnSt.ReturnValue not %s, got %s", testStatements[i], returnSt.ReturnValue.TokenLiteral())
		//}
	}
}

func TestParser_Errors(t *testing.T) {
	input := `
		let x 5;
		let = 10;
		let 838383;
		`

	l := lexer.New(input)
	p := parser.New(l)

	_ = p.ParseProgram()
	errors := p.Errors()
	if len(errors) == 0 {
		t.Errorf("parser.Errors() returned no errors")
	}

	assert.Equal(t, 3, len(errors))
	assert.Equal(t, "expected next token to be =, got INT instead", errors[0])
	assert.Equal(t, "expected next token to be IDENT, got = instead", errors[1])
	assert.Equal(t, "expected next token to be IDENT, got INT instead", errors[2])
}

func TestParseProgram(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
		`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	testStatements := []string{
		"x", "y", "foobar",
	}

	for i, st := range program.Statements {
		if !testLetStatement(t, st, testStatements[i]) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestIntegerExpression(t *testing.T) {
	input := "2137;"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if ident.Value != 2137 {
		t.Errorf("ident.Value not %d. got=%d", 2137, ident.Value)
	}
	if ident.TokenLiteral() != "2137" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "2137",
			ident.TokenLiteral())
	}
}
