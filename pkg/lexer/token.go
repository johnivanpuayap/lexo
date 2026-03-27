package lexer

import "fmt"

type TokenType int

const (
	// Literals
	INT_LIT TokenType = iota
	FLOAT_LIT
	STRING_LIT

	// Identifier
	IDENTIFIER

	// Keywords
	DEF
	RETURN
	IF
	ELIF
	ELSE
	WHILE
	FOR
	IN
	BREAK
	CONTINUE
	AND
	OR
	NOT
	TRUE
	FALSE
	VOID

	// Type keywords
	INT_TYPE
	FLOAT_TYPE
	STRING_TYPE
	BOOL_TYPE

	// Operators
	PLUS
	MINUS
	STAR
	SLASH
	PERCENT
	ASSIGN
	EQ
	NEQ
	LT
	GT
	LTE
	GTE

	// Delimiters
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	COLON
	COMMA
	ARROW
	DOT

	// Whitespace / structure
	NEWLINE
	INDENT
	DEDENT

	// Special
	EOF
)

var tokenNames = map[TokenType]string{
	INT_LIT: "INT_LIT", FLOAT_LIT: "FLOAT_LIT", STRING_LIT: "STRING_LIT",
	IDENTIFIER: "IDENTIFIER",
	DEF: "DEF", RETURN: "RETURN", IF: "IF", ELIF: "ELIF", ELSE: "ELSE",
	WHILE: "WHILE", FOR: "FOR", IN: "IN", BREAK: "BREAK", CONTINUE: "CONTINUE",
	AND: "AND", OR: "OR", NOT: "NOT", TRUE: "TRUE", FALSE: "FALSE", VOID: "VOID",
	INT_TYPE: "INT_TYPE", FLOAT_TYPE: "FLOAT_TYPE", STRING_TYPE: "STRING_TYPE", BOOL_TYPE: "BOOL_TYPE",
	PLUS: "PLUS", MINUS: "MINUS", STAR: "STAR", SLASH: "SLASH", PERCENT: "PERCENT",
	ASSIGN: "ASSIGN", EQ: "EQ", NEQ: "NEQ", LT: "LT", GT: "GT", LTE: "LTE", GTE: "GTE",
	LPAREN: "LPAREN", RPAREN: "RPAREN", LBRACKET: "LBRACKET", RBRACKET: "RBRACKET",
	COLON: "COLON", COMMA: "COMMA", ARROW: "ARROW", DOT: "DOT",
	NEWLINE: "NEWLINE", INDENT: "INDENT", DEDENT: "DEDENT", EOF: "EOF",
}

func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return fmt.Sprintf("UNKNOWN(%d)", int(t))
}

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q, %d:%d)", t.Type, t.Value, t.Line, t.Column)
}

var keywords = map[string]TokenType{
	"def":      DEF,
	"return":   RETURN,
	"if":       IF,
	"elif":     ELIF,
	"else":     ELSE,
	"while":    WHILE,
	"for":      FOR,
	"in":       IN,
	"break":    BREAK,
	"continue": CONTINUE,
	"and":      AND,
	"or":       OR,
	"not":      NOT,
	"true":     TRUE,
	"false":    FALSE,
	"void":     VOID,
	"int":      INT_TYPE,
	"float":    FLOAT_TYPE,
	"string":   STRING_TYPE,
	"bool":     BOOL_TYPE,
}

func LookupKeyword(ident string) (TokenType, bool) {
	tok, ok := keywords[ident]
	return tok, ok
}
