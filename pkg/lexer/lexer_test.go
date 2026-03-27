package lexer

import (
	"testing"
)

func tokenTypes(tokens []Token) []TokenType {
	types := make([]TokenType, len(tokens))
	for i, t := range tokens {
		types[i] = t.Type
	}
	return types
}

func TestSimpleVarDecl(t *testing.T) {
	tokens, err := Tokenize(`name: string = "hello"`)
	if err != nil {
		t.Fatal(err)
	}
	expected := []TokenType{IDENTIFIER, COLON, STRING_TYPE, ASSIGN, STRING_LIT, NEWLINE, EOF}
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("expected %d tokens, got %d: %v", len(expected), len(types), tokens)
	}
	for i, exp := range expected {
		if types[i] != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, types[i])
		}
	}
}

func TestNumberLiterals(t *testing.T) {
	tokens, err := Tokenize("42 3.14 -7")
	if err != nil {
		t.Fatal(err)
	}
	// -7 is tokenized as MINUS + INT_LIT (unary negation handled by parser)
	expected := []TokenType{INT_LIT, FLOAT_LIT, MINUS, INT_LIT, NEWLINE, EOF}
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("expected %d tokens, got %d: %v", len(expected), len(types), tokens)
	}
	for i, exp := range expected {
		if types[i] != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, types[i])
		}
	}
	if tokens[0].Value != "42" {
		t.Errorf("expected 42, got %s", tokens[0].Value)
	}
	if tokens[1].Value != "3.14" {
		t.Errorf("expected 3.14, got %s", tokens[1].Value)
	}
}

func TestOperators(t *testing.T) {
	tokens, err := Tokenize("== != <= >= -> + - * / %")
	if err != nil {
		t.Fatal(err)
	}
	expected := []TokenType{EQ, NEQ, LTE, GTE, ARROW, PLUS, MINUS, STAR, SLASH, PERCENT, NEWLINE, EOF}
	types := tokenTypes(tokens)
	for i, exp := range expected {
		if i >= len(types) {
			t.Fatalf("missing token at position %d", i)
		}
		if types[i] != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, types[i])
		}
	}
}

func TestKeywords(t *testing.T) {
	tokens, err := Tokenize("def return if elif else while for in break continue and or not true false void int float string bool")
	if err != nil {
		t.Fatal(err)
	}
	expected := []TokenType{
		DEF, RETURN, IF, ELIF, ELSE, WHILE, FOR, IN, BREAK, CONTINUE,
		AND, OR, NOT, TRUE, FALSE, VOID, INT_TYPE, FLOAT_TYPE, STRING_TYPE, BOOL_TYPE,
		NEWLINE, EOF,
	}
	types := tokenTypes(tokens)
	for i, exp := range expected {
		if types[i] != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, types[i])
		}
	}
}

func TestIndentDedent(t *testing.T) {
	src := "if true:\n    x: int = 1\n    y: int = 2\nz: int = 3\n"
	tokens, err := Tokenize(src)
	if err != nil {
		t.Fatal(err)
	}
	expected := []TokenType{
		IF, TRUE, COLON, NEWLINE,
		INDENT,
		IDENTIFIER, COLON, INT_TYPE, ASSIGN, INT_LIT, NEWLINE,
		IDENTIFIER, COLON, INT_TYPE, ASSIGN, INT_LIT, NEWLINE,
		DEDENT,
		IDENTIFIER, COLON, INT_TYPE, ASSIGN, INT_LIT, NEWLINE,
		EOF,
	}
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("expected %d tokens, got %d:\n%v", len(expected), len(types), tokens)
	}
	for i, exp := range expected {
		if types[i] != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, types[i])
		}
	}
}

func TestNestedIndent(t *testing.T) {
	src := "def foo() -> void:\n    if true:\n        x: int = 1\n    y: int = 2\n"
	tokens, err := Tokenize(src)
	if err != nil {
		t.Fatal(err)
	}
	expected := []TokenType{
		DEF, IDENTIFIER, LPAREN, RPAREN, ARROW, VOID, COLON, NEWLINE,
		INDENT,
		IF, TRUE, COLON, NEWLINE,
		INDENT,
		IDENTIFIER, COLON, INT_TYPE, ASSIGN, INT_LIT, NEWLINE,
		DEDENT,
		IDENTIFIER, COLON, INT_TYPE, ASSIGN, INT_LIT, NEWLINE,
		DEDENT,
		EOF,
	}
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("expected %d tokens, got %d:\n%v", len(expected), len(types), tokens)
	}
	for i, exp := range expected {
		if types[i] != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, types[i])
		}
	}
}

func TestStringLiteral(t *testing.T) {
	tokens, err := Tokenize(`"hello world"`)
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Type != STRING_LIT {
		t.Errorf("expected STRING_LIT, got %s", tokens[0].Type)
	}
	if tokens[0].Value != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", tokens[0].Value)
	}
}

func TestComment(t *testing.T) {
	tokens, err := Tokenize("x: int = 1 # this is a comment\ny: int = 2\n")
	if err != nil {
		t.Fatal(err)
	}
	// Comments are skipped entirely
	expected := []TokenType{
		IDENTIFIER, COLON, INT_TYPE, ASSIGN, INT_LIT, NEWLINE,
		IDENTIFIER, COLON, INT_TYPE, ASSIGN, INT_LIT, NEWLINE,
		EOF,
	}
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("expected %d tokens, got %d:\n%v", len(expected), len(types), tokens)
	}
	for i, exp := range expected {
		if types[i] != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, types[i])
		}
	}
}

func TestBlankLinesIgnored(t *testing.T) {
	src := "if true:\n    x: int = 1\n\n    y: int = 2\n"
	tokens, err := Tokenize(src)
	if err != nil {
		t.Fatal(err)
	}
	// Blank line between indented lines should NOT produce DEDENT/INDENT
	expected := []TokenType{
		IF, TRUE, COLON, NEWLINE,
		INDENT,
		IDENTIFIER, COLON, INT_TYPE, ASSIGN, INT_LIT, NEWLINE,
		IDENTIFIER, COLON, INT_TYPE, ASSIGN, INT_LIT, NEWLINE,
		DEDENT,
		EOF,
	}
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("expected %d tokens, got %d:\n%v", len(expected), len(types), tokens)
	}
	for i, exp := range expected {
		if types[i] != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, types[i])
		}
	}
}

func TestArraySyntax(t *testing.T) {
	tokens, err := Tokenize("nums: int[] = [1, 2, 3]")
	if err != nil {
		t.Fatal(err)
	}
	expected := []TokenType{
		IDENTIFIER, COLON, INT_TYPE, LBRACKET, RBRACKET, ASSIGN,
		LBRACKET, INT_LIT, COMMA, INT_LIT, COMMA, INT_LIT, RBRACKET,
		NEWLINE, EOF,
	}
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("expected %d tokens, got %d:\n%v", len(expected), len(types), tokens)
	}
	for i, exp := range expected {
		if types[i] != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, types[i])
		}
	}
}

func TestFunctionDef(t *testing.T) {
	tokens, err := Tokenize("def add(a: int, b: int) -> int:")
	if err != nil {
		t.Fatal(err)
	}
	expected := []TokenType{
		DEF, IDENTIFIER, LPAREN, IDENTIFIER, COLON, INT_TYPE, COMMA,
		IDENTIFIER, COLON, INT_TYPE, RPAREN, ARROW, INT_TYPE, COLON,
		NEWLINE, EOF,
	}
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("expected %d tokens, got %d:\n%v", len(expected), len(types), tokens)
	}
	for i, exp := range expected {
		if types[i] != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, types[i])
		}
	}
}

func TestMethodCallSyntax(t *testing.T) {
	tokens, err := Tokenize(`msg.length()`)
	if err != nil {
		t.Fatal(err)
	}
	expected := []TokenType{IDENTIFIER, DOT, IDENTIFIER, LPAREN, RPAREN, NEWLINE, EOF}
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("expected %d tokens, got %d:\n%v", len(expected), len(types), tokens)
	}
	for i, exp := range expected {
		if types[i] != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, types[i])
		}
	}
}

func TestUnterminatedString(t *testing.T) {
	_, err := Tokenize(`"hello`)
	if err == nil {
		t.Fatal("expected error for unterminated string")
	}
}

func TestLineAndColumn(t *testing.T) {
	tokens, err := Tokenize("x: int = 1\ny: int = 2\n")
	if err != nil {
		t.Fatal(err)
	}
	// x should be at line 1, column 1
	if tokens[0].Line != 1 || tokens[0].Column != 1 {
		t.Errorf("x: expected 1:1, got %d:%d", tokens[0].Line, tokens[0].Column)
	}
	// y should be at line 2, column 1
	yIdx := 6 // After IDENT COLON INT_TYPE ASSIGN INT_LIT NEWLINE
	if tokens[yIdx].Line != 2 || tokens[yIdx].Column != 1 {
		t.Errorf("y: expected 2:1, got %d:%d", tokens[yIdx].Line, tokens[yIdx].Column)
	}
}
