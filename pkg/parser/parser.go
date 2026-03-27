package parser

import (
	"fmt"
	"strconv"

	"github.com/johnivanpuayap/lexo/pkg/lexer"
)

type ParseError struct {
	Message string
	Line    int
	Column  int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Line %d, Column %d: %s", e.Line, e.Column, e.Message)
}

type parser struct {
	tokens  []lexer.Token
	pos     int
	current lexer.Token
}

func Parse(tokens []lexer.Token) (*Program, error) {
	p := &parser{
		tokens:  tokens,
		pos:     0,
		current: tokens[0],
	}
	return p.parseProgram()
}

func (p *parser) parseProgram() (*Program, error) {
	prog := &Program{}
	for p.current.Type != lexer.EOF {
		// Skip bare newlines at top level
		if p.current.Type == lexer.NEWLINE {
			p.advance()
			continue
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		prog.Body = append(prog.Body, stmt)
	}
	return prog, nil
}

func (p *parser) parseStatement() (Statement, error) {
	switch p.current.Type {
	case lexer.DEF:
		return p.parseFuncDecl()
	case lexer.IF:
		return p.parseIfStmt()
	case lexer.WHILE:
		return p.parseWhileStmt()
	case lexer.FOR:
		return p.parseForStmt()
	case lexer.RETURN:
		return p.parseReturnStmt()
	case lexer.BREAK:
		return p.parseBreakStmt()
	case lexer.CONTINUE:
		return p.parseContinueStmt()
	case lexer.IDENTIFIER:
		if p.peekType() == lexer.COLON {
			return p.parseVarDecl()
		}
		return p.parseExprStmtOrAssignment()
	default:
		return p.parseExprStmtOrAssignment()
	}
}

func (p *parser) parseVarDecl() (*VarDecl, error) {
	line := p.current.Line
	name := p.current.Value
	p.advance() // skip identifier

	if err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}

	typeName, err := p.parseTypeAnnotation()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.ASSIGN); err != nil {
		return nil, err
	}

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.NEWLINE); err != nil {
		return nil, err
	}

	return &VarDecl{Name: name, TypeName: typeName, Value: value, Line: line}, nil
}

func (p *parser) parseExprStmtOrAssignment() (Statement, error) {
	line := p.current.Line
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.current.Type == lexer.ASSIGN {
		p.advance() // skip =
		value, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(lexer.NEWLINE); err != nil {
			return nil, err
		}
		return &Assignment{Target: expr, Value: value, Line: line}, nil
	}

	if err := p.expect(lexer.NEWLINE); err != nil {
		return nil, err
	}
	return &ExpressionStmt{Expr: expr, Line: line}, nil
}

func (p *parser) parseFuncDecl() (*FuncDecl, error) {
	line := p.current.Line
	p.advance() // skip 'def'

	if p.current.Type != lexer.IDENTIFIER {
		return nil, p.error("expected function name")
	}
	name := p.current.Value
	p.advance()

	if err := p.expect(lexer.LPAREN); err != nil {
		return nil, err
	}

	var params []Param
	if p.current.Type != lexer.RPAREN {
		for {
			if p.current.Type != lexer.IDENTIFIER {
				return nil, p.error("expected parameter name")
			}
			paramName := p.current.Value
			p.advance()

			if err := p.expect(lexer.COLON); err != nil {
				return nil, err
			}

			paramType, err := p.parseTypeAnnotation()
			if err != nil {
				return nil, err
			}

			params = append(params, Param{Name: paramName, TypeName: paramType})

			if p.current.Type != lexer.COMMA {
				break
			}
			p.advance() // skip comma
		}
	}

	if err := p.expect(lexer.RPAREN); err != nil {
		return nil, err
	}
	if err := p.expect(lexer.ARROW); err != nil {
		return nil, err
	}

	returnType, err := p.parseTypeAnnotation()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}
	if err := p.expect(lexer.NEWLINE); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &FuncDecl{
		Name: name, Params: params, ReturnType: returnType, Body: body, Line: line,
	}, nil
}

func (p *parser) parseIfStmt() (*IfStmt, error) {
	line := p.current.Line
	p.advance() // skip 'if'

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}
	if err := p.expect(lexer.NEWLINE); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	var elifs []ElifClause
	for p.current.Type == lexer.ELIF {
		p.advance() // skip 'elif'
		elifCond, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(lexer.COLON); err != nil {
			return nil, err
		}
		if err := p.expect(lexer.NEWLINE); err != nil {
			return nil, err
		}
		elifBody, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		elifs = append(elifs, ElifClause{Condition: elifCond, Body: elifBody})
	}

	var elseBody []Statement
	if p.current.Type == lexer.ELSE {
		p.advance() // skip 'else'
		if err := p.expect(lexer.COLON); err != nil {
			return nil, err
		}
		if err := p.expect(lexer.NEWLINE); err != nil {
			return nil, err
		}
		elseBody, err = p.parseBlock()
		if err != nil {
			return nil, err
		}
	}

	return &IfStmt{
		Condition: condition, Body: body, Elifs: elifs, ElseBody: elseBody, Line: line,
	}, nil
}

func (p *parser) parseWhileStmt() (*WhileStmt, error) {
	line := p.current.Line
	p.advance() // skip 'while'

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}
	if err := p.expect(lexer.NEWLINE); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &WhileStmt{Condition: condition, Body: body, Line: line}, nil
}

func (p *parser) parseForStmt() (*ForStmt, error) {
	line := p.current.Line
	p.advance() // skip 'for'

	if p.current.Type != lexer.IDENTIFIER {
		return nil, p.error("expected loop variable name")
	}
	varName := p.current.Value
	p.advance()

	if err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}

	varType, err := p.parseTypeAnnotation()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.IN); err != nil {
		return nil, err
	}

	iterable, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}
	if err := p.expect(lexer.NEWLINE); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ForStmt{
		VarName: varName, VarType: varType, Iterable: iterable, Body: body, Line: line,
	}, nil
}

func (p *parser) parseReturnStmt() (*ReturnStmt, error) {
	line := p.current.Line
	p.advance() // skip 'return'

	var value Expression
	if p.current.Type != lexer.NEWLINE {
		var err error
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	if err := p.expect(lexer.NEWLINE); err != nil {
		return nil, err
	}

	return &ReturnStmt{Value: value, Line: line}, nil
}

func (p *parser) parseBreakStmt() (*BreakStmt, error) {
	line := p.current.Line
	p.advance() // skip 'break'
	if err := p.expect(lexer.NEWLINE); err != nil {
		return nil, err
	}
	return &BreakStmt{Line: line}, nil
}

func (p *parser) parseContinueStmt() (*ContinueStmt, error) {
	line := p.current.Line
	p.advance() // skip 'continue'
	if err := p.expect(lexer.NEWLINE); err != nil {
		return nil, err
	}
	return &ContinueStmt{Line: line}, nil
}

func (p *parser) parseBlock() ([]Statement, error) {
	if err := p.expect(lexer.INDENT); err != nil {
		return nil, err
	}

	var stmts []Statement
	for p.current.Type != lexer.DEDENT && p.current.Type != lexer.EOF {
		if p.current.Type == lexer.NEWLINE {
			p.advance()
			continue
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	if p.current.Type == lexer.DEDENT {
		p.advance() // consume DEDENT
	}

	return stmts, nil
}

// --- Expression Parsing (precedence climbing) ---

func (p *parser) parseExpression() (Expression, error) {
	return p.parseOr()
}

func (p *parser) parseOr() (Expression, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.current.Type == lexer.OR {
		line := p.current.Line
		p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: "or", Right: right, Line: line}
	}
	return left, nil
}

func (p *parser) parseAnd() (Expression, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for p.current.Type == lexer.AND {
		line := p.current.Line
		p.advance()
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: "and", Right: right, Line: line}
	}
	return left, nil
}

func (p *parser) parseNot() (Expression, error) {
	if p.current.Type == lexer.NOT {
		line := p.current.Line
		p.advance()
		operand, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{Operator: "not", Operand: operand, Line: line}, nil
	}
	return p.parseComparison()
}

func (p *parser) parseComparison() (Expression, error) {
	left, err := p.parseAddition()
	if err != nil {
		return nil, err
	}
	if isComparisonOp(p.current.Type) {
		line := p.current.Line
		op := tokenToOp(p.current.Type)
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: op, Right: right, Line: line}
	}
	return left, nil
}

func (p *parser) parseAddition() (Expression, error) {
	left, err := p.parseMultiplication()
	if err != nil {
		return nil, err
	}
	for p.current.Type == lexer.PLUS || p.current.Type == lexer.MINUS {
		line := p.current.Line
		op := tokenToOp(p.current.Type)
		p.advance()
		right, err := p.parseMultiplication()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: op, Right: right, Line: line}
	}
	return left, nil
}

func (p *parser) parseMultiplication() (Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.current.Type == lexer.STAR || p.current.Type == lexer.SLASH || p.current.Type == lexer.PERCENT {
		line := p.current.Line
		op := tokenToOp(p.current.Type)
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: op, Right: right, Line: line}
	}
	return left, nil
}

func (p *parser) parseUnary() (Expression, error) {
	if p.current.Type == lexer.MINUS {
		line := p.current.Line
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{Operator: "-", Operand: operand, Line: line}, nil
	}
	return p.parsePostfix()
}

func (p *parser) parsePostfix() (Expression, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		switch p.current.Type {
		case lexer.LPAREN:
			// Function call: expr(args)
			// expr must be an Identifier
			ident, ok := expr.(*Identifier)
			if !ok {
				return nil, &ParseError{
					Message: "only named functions can be called",
					Line:    p.current.Line,
					Column:  p.current.Column,
				}
			}
			line := p.current.Line
			p.advance() // skip (
			args, err := p.parseArgList()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.RPAREN); err != nil {
				return nil, err
			}
			expr = &FuncCall{Name: ident.Name, Args: args, Line: line}

		case lexer.LBRACKET:
			// Index access: expr[index]
			line := p.current.Line
			p.advance() // skip [
			index, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.RBRACKET); err != nil {
				return nil, err
			}
			expr = &IndexAccess{Object: expr, Index: index, Line: line}

		case lexer.DOT:
			// Method call: expr.method(args)
			p.advance() // skip .
			if p.current.Type != lexer.IDENTIFIER {
				return nil, p.error("expected method name after '.'")
			}
			method := p.current.Value
			line := p.current.Line
			p.advance()

			if err := p.expect(lexer.LPAREN); err != nil {
				return nil, err
			}
			args, err := p.parseArgList()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.RPAREN); err != nil {
				return nil, err
			}
			expr = &MethodCall{Object: expr, Method: method, Args: args, Line: line}

		default:
			return expr, nil
		}
	}
}

func (p *parser) parsePrimary() (Expression, error) {
	switch p.current.Type {
	case lexer.INT_LIT:
		val, _ := strconv.ParseInt(p.current.Value, 10, 64)
		lit := &Literal{Value: val, LitType: "int", Line: p.current.Line}
		p.advance()
		return lit, nil

	case lexer.FLOAT_LIT:
		val, _ := strconv.ParseFloat(p.current.Value, 64)
		lit := &Literal{Value: val, LitType: "float", Line: p.current.Line}
		p.advance()
		return lit, nil

	case lexer.STRING_LIT:
		lit := &Literal{Value: p.current.Value, LitType: "string", Line: p.current.Line}
		p.advance()
		return lit, nil

	case lexer.TRUE:
		lit := &Literal{Value: true, LitType: "bool", Line: p.current.Line}
		p.advance()
		return lit, nil

	case lexer.FALSE:
		lit := &Literal{Value: false, LitType: "bool", Line: p.current.Line}
		p.advance()
		return lit, nil

	case lexer.IDENTIFIER:
		ident := &Identifier{Name: p.current.Value, Line: p.current.Line}
		p.advance()
		return ident, nil

	case lexer.LPAREN:
		p.advance() // skip (
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(lexer.RPAREN); err != nil {
			return nil, err
		}
		return expr, nil

	case lexer.LBRACKET:
		return p.parseArrayLiteral()

	default:
		return nil, p.error(fmt.Sprintf("unexpected token %s", p.current.Type))
	}
}

func (p *parser) parseArrayLiteral() (*ArrayLiteral, error) {
	line := p.current.Line
	p.advance() // skip [

	var elements []Expression
	if p.current.Type != lexer.RBRACKET {
		for {
			elem, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			elements = append(elements, elem)
			if p.current.Type != lexer.COMMA {
				break
			}
			p.advance() // skip comma
		}
	}

	if err := p.expect(lexer.RBRACKET); err != nil {
		return nil, err
	}

	return &ArrayLiteral{Elements: elements, Line: line}, nil
}

func (p *parser) parseArgList() ([]Expression, error) {
	var args []Expression
	if p.current.Type == lexer.RPAREN {
		return args, nil
	}
	for {
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if p.current.Type != lexer.COMMA {
			break
		}
		p.advance() // skip comma
	}
	return args, nil
}

func (p *parser) parseTypeAnnotation() (TypeAnnotation, error) {
	var base string
	switch p.current.Type {
	case lexer.INT_TYPE:
		base = "int"
	case lexer.FLOAT_TYPE:
		base = "float"
	case lexer.STRING_TYPE:
		base = "string"
	case lexer.BOOL_TYPE:
		base = "bool"
	case lexer.VOID:
		p.advance()
		return TypeAnnotation("void"), nil
	default:
		return "", p.error(fmt.Sprintf("expected type, got %s", p.current.Type))
	}
	p.advance()

	// Check for array type: int[], string[], etc.
	if p.current.Type == lexer.LBRACKET {
		p.advance()
		if err := p.expect(lexer.RBRACKET); err != nil {
			return "", err
		}
		return TypeAnnotation(base + "[]"), nil
	}

	return TypeAnnotation(base), nil
}

// --- Helpers ---

func (p *parser) advance() {
	p.pos++
	if p.pos < len(p.tokens) {
		p.current = p.tokens[p.pos]
	}
}

func (p *parser) peekType() lexer.TokenType {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1].Type
	}
	return lexer.EOF
}

func (p *parser) expect(tt lexer.TokenType) error {
	if p.current.Type != tt {
		return &ParseError{
			Message: fmt.Sprintf("expected %s, got %s", tt, p.current.Type),
			Line:    p.current.Line,
			Column:  p.current.Column,
		}
	}
	p.advance()
	return nil
}

func (p *parser) error(msg string) *ParseError {
	return &ParseError{Message: msg, Line: p.current.Line, Column: p.current.Column}
}

func isComparisonOp(t lexer.TokenType) bool {
	return t == lexer.EQ || t == lexer.NEQ || t == lexer.LT || t == lexer.GT || t == lexer.LTE || t == lexer.GTE
}

func tokenToOp(t lexer.TokenType) string {
	switch t {
	case lexer.PLUS:
		return "+"
	case lexer.MINUS:
		return "-"
	case lexer.STAR:
		return "*"
	case lexer.SLASH:
		return "/"
	case lexer.PERCENT:
		return "%"
	case lexer.EQ:
		return "=="
	case lexer.NEQ:
		return "!="
	case lexer.LT:
		return "<"
	case lexer.GT:
		return ">"
	case lexer.LTE:
		return "<="
	case lexer.GTE:
		return ">="
	default:
		return ""
	}
}
