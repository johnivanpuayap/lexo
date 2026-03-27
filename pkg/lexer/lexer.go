package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type LexerError struct {
	Message string
	Line    int
	Column  int
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("Line %d, Column %d: %s", e.Line, e.Column, e.Message)
}

type lexer struct {
	source      string
	pos         int
	line        int
	column      int
	tokens      []Token
	indentStack []int
	atLineStart bool
}

func Tokenize(source string) ([]Token, error) {
	// Ensure source ends with newline for consistent handling
	if len(source) > 0 && source[len(source)-1] != '\n' {
		source += "\n"
	}

	l := &lexer{
		source:      source,
		pos:         0,
		line:        1,
		column:      1,
		tokens:      []Token{},
		indentStack: []int{0},
		atLineStart: true,
	}

	if err := l.tokenize(); err != nil {
		return nil, err
	}

	return l.tokens, nil
}

func (l *lexer) tokenize() error {
	for l.pos < len(l.source) {
		if l.atLineStart {
			if err := l.handleIndentation(); err != nil {
				return err
			}
			if l.pos >= len(l.source) {
				break
			}
			// If still atLineStart (blank line was consumed), re-enter loop
			if l.atLineStart {
				continue
			}
		}

		ch := l.source[l.pos]

		switch {
		case ch == '\n':
			l.emit(NEWLINE, "\\n")
			l.advance()
			l.atLineStart = true

		case ch == '#':
			l.skipComment()

		case ch == ' ' || ch == '\t':
			l.advance() // skip whitespace within a line

		case ch == '"':
			if err := l.readString(); err != nil {
				return err
			}

		case isDigit(ch):
			l.readNumber()

		case isAlpha(ch) || ch == '_':
			l.readIdentifier()

		case ch == '+':
			l.emit(PLUS, "+")
			l.advance()
		case ch == '*':
			l.emit(STAR, "*")
			l.advance()
		case ch == '/':
			l.emit(SLASH, "/")
			l.advance()
		case ch == '%':
			l.emit(PERCENT, "%")
			l.advance()
		case ch == '(':
			l.emit(LPAREN, "(")
			l.advance()
		case ch == ')':
			l.emit(RPAREN, ")")
			l.advance()
		case ch == '[':
			l.emit(LBRACKET, "[")
			l.advance()
		case ch == ']':
			l.emit(RBRACKET, "]")
			l.advance()
		case ch == ',':
			l.emit(COMMA, ",")
			l.advance()
		case ch == '.':
			l.emit(DOT, ".")
			l.advance()
		case ch == ':':
			l.emit(COLON, ":")
			l.advance()

		case ch == '-':
			if l.peek() == '>' {
				l.emit(ARROW, "->")
				l.advance()
				l.advance()
			} else {
				l.emit(MINUS, "-")
				l.advance()
			}

		case ch == '=':
			if l.peek() == '=' {
				l.emit(EQ, "==")
				l.advance()
				l.advance()
			} else {
				l.emit(ASSIGN, "=")
				l.advance()
			}

		case ch == '!':
			if l.peek() == '=' {
				l.emit(NEQ, "!=")
				l.advance()
				l.advance()
			} else {
				return l.error("unexpected character '!'")
			}

		case ch == '<':
			if l.peek() == '=' {
				l.emit(LTE, "<=")
				l.advance()
				l.advance()
			} else {
				l.emit(LT, "<")
				l.advance()
			}

		case ch == '>':
			if l.peek() == '=' {
				l.emit(GTE, ">=")
				l.advance()
				l.advance()
			} else {
				l.emit(GT, ">")
				l.advance()
			}

		default:
			return l.error(fmt.Sprintf("unexpected character '%c'", ch))
		}
	}

	// Emit remaining DEDENTs at EOF
	for len(l.indentStack) > 1 {
		l.indentStack = l.indentStack[:len(l.indentStack)-1]
		l.tokens = append(l.tokens, Token{Type: DEDENT, Value: "", Line: l.line, Column: l.column})
	}

	l.tokens = append(l.tokens, Token{Type: EOF, Value: "", Line: l.line, Column: l.column})
	return nil
}

func (l *lexer) handleIndentation() error {
	spaces := 0
	startPos := l.pos
	for l.pos < len(l.source) && l.source[l.pos] == ' ' {
		spaces++
		l.pos++
		l.column++
	}

	if l.pos < len(l.source) && l.source[l.pos] == '\t' {
		return &LexerError{
			Message: "tabs are not allowed for indentation, use spaces",
			Line:    l.line,
			Column:  l.column,
		}
	}

	// Skip blank lines and comment-only lines — don't change indentation
	if l.pos >= len(l.source) {
		_ = startPos
		l.atLineStart = false
		return nil
	}
	if l.source[l.pos] == '\n' {
		// Blank line: consume the newline and stay in atLineStart mode
		l.advance()
		l.atLineStart = true
		return nil
	}
	if l.source[l.pos] == '#' {
		// Comment-only line: skip comment, then let the main loop handle the newline
		_ = startPos
		l.atLineStart = false
		return nil
	}

	l.atLineStart = false
	currentIndent := l.indentStack[len(l.indentStack)-1]

	if spaces > currentIndent {
		l.indentStack = append(l.indentStack, spaces)
		l.tokens = append(l.tokens, Token{Type: INDENT, Value: "", Line: l.line, Column: 1})
	} else if spaces < currentIndent {
		for len(l.indentStack) > 1 && l.indentStack[len(l.indentStack)-1] > spaces {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			l.tokens = append(l.tokens, Token{Type: DEDENT, Value: "", Line: l.line, Column: 1})
		}
		if l.indentStack[len(l.indentStack)-1] != spaces {
			return &LexerError{
				Message: fmt.Sprintf("inconsistent indentation: expected %d spaces, got %d", l.indentStack[len(l.indentStack)-1], spaces),
				Line:    l.line,
				Column:  1,
			}
		}
	}

	return nil
}

func (l *lexer) readString() error {
	startLine := l.line
	startCol := l.column
	l.advance() // skip opening "

	var sb strings.Builder
	for l.pos < len(l.source) {
		ch := l.source[l.pos]
		if ch == '"' {
			l.tokens = append(l.tokens, Token{Type: STRING_LIT, Value: sb.String(), Line: startLine, Column: startCol})
			l.advance() // skip closing "
			return nil
		}
		if ch == '\n' {
			return &LexerError{Message: "unterminated string literal", Line: startLine, Column: startCol}
		}
		if ch == '\\' && l.pos+1 < len(l.source) {
			next := l.source[l.pos+1]
			switch next {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			default:
				sb.WriteByte(ch)
				sb.WriteByte(next)
			}
			l.advance()
			l.advance()
			continue
		}
		sb.WriteByte(ch)
		l.advance()
	}
	return &LexerError{Message: "unterminated string literal", Line: startLine, Column: startCol}
}

func (l *lexer) readNumber() {
	startCol := l.column
	start := l.pos
	isFloat := false

	for l.pos < len(l.source) && isDigit(l.source[l.pos]) {
		l.advance()
	}

	if l.pos < len(l.source) && l.source[l.pos] == '.' && l.pos+1 < len(l.source) && isDigit(l.source[l.pos+1]) {
		isFloat = true
		l.advance() // skip .
		for l.pos < len(l.source) && isDigit(l.source[l.pos]) {
			l.advance()
		}
	}

	value := l.source[start:l.pos]
	tokType := INT_LIT
	if isFloat {
		tokType = FLOAT_LIT
	}
	l.tokens = append(l.tokens, Token{Type: tokType, Value: value, Line: l.line, Column: startCol})
}

func (l *lexer) readIdentifier() {
	startCol := l.column
	start := l.pos

	for l.pos < len(l.source) && (isAlpha(l.source[l.pos]) || isDigit(l.source[l.pos]) || l.source[l.pos] == '_') {
		l.advance()
	}

	value := l.source[start:l.pos]

	if tokType, ok := LookupKeyword(value); ok {
		l.tokens = append(l.tokens, Token{Type: tokType, Value: value, Line: l.line, Column: startCol})
	} else {
		l.tokens = append(l.tokens, Token{Type: IDENTIFIER, Value: value, Line: l.line, Column: startCol})
	}
}

func (l *lexer) skipComment() {
	for l.pos < len(l.source) && l.source[l.pos] != '\n' {
		l.advance()
	}
}

func (l *lexer) emit(tokenType TokenType, value string) {
	l.tokens = append(l.tokens, Token{Type: tokenType, Value: value, Line: l.line, Column: l.column})
}

func (l *lexer) advance() {
	if l.pos < len(l.source) {
		if l.source[l.pos] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.pos++
	}
}

func (l *lexer) peek() byte {
	if l.pos+1 < len(l.source) {
		return l.source[l.pos+1]
	}
	return 0
}

func (l *lexer) error(message string) *LexerError {
	return &LexerError{Message: message, Line: l.line, Column: l.column}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlpha(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}
