package parser

import (
	"testing"

	"github.com/johnivanpuayap/lexo/pkg/lexer"
)

func mustParse(t *testing.T, source string) *Program {
	t.Helper()
	tokens, err := lexer.Tokenize(source)
	if err != nil {
		t.Fatalf("lexer error: %v", err)
	}
	prog, err := Parse(tokens)
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
	return prog
}

func mustFail(t *testing.T, source string) {
	t.Helper()
	tokens, err := lexer.Tokenize(source)
	if err != nil {
		return // lexer error is fine too
	}
	_, err = Parse(tokens)
	if err == nil {
		t.Fatal("expected parser error, got none")
	}
}

func TestParseVarDecl(t *testing.T) {
	prog := mustParse(t, `name: string = "hello"`)
	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Body))
	}
	decl, ok := prog.Body[0].(*VarDecl)
	if !ok {
		t.Fatalf("expected VarDecl, got %T", prog.Body[0])
	}
	if decl.Name != "name" {
		t.Errorf("expected name 'name', got '%s'", decl.Name)
	}
	if decl.TypeName != "string" {
		t.Errorf("expected type 'string', got '%s'", decl.TypeName)
	}
	lit, ok := decl.Value.(*Literal)
	if !ok {
		t.Fatalf("expected Literal, got %T", decl.Value)
	}
	if lit.Value != "hello" {
		t.Errorf("expected value 'hello', got '%v'", lit.Value)
	}
}

func TestParseIntVarDecl(t *testing.T) {
	prog := mustParse(t, "age: int = 25")
	decl := prog.Body[0].(*VarDecl)
	if decl.TypeName != "int" {
		t.Errorf("expected type 'int', got '%s'", decl.TypeName)
	}
	lit := decl.Value.(*Literal)
	if lit.Value != int64(25) {
		t.Errorf("expected 25, got %v", lit.Value)
	}
}

func TestParseArrayDecl(t *testing.T) {
	prog := mustParse(t, "nums: int[] = [1, 2, 3]")
	decl := prog.Body[0].(*VarDecl)
	if decl.TypeName != "int[]" {
		t.Errorf("expected type 'int[]', got '%s'", decl.TypeName)
	}
	arr, ok := decl.Value.(*ArrayLiteral)
	if !ok {
		t.Fatalf("expected ArrayLiteral, got %T", decl.Value)
	}
	if len(arr.Elements) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arr.Elements))
	}
}

func TestParseAssignment(t *testing.T) {
	prog := mustParse(t, "x: int = 1\nx = 2\n")
	if len(prog.Body) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(prog.Body))
	}
	assign, ok := prog.Body[1].(*Assignment)
	if !ok {
		t.Fatalf("expected Assignment, got %T", prog.Body[1])
	}
	ident := assign.Target.(*Identifier)
	if ident.Name != "x" {
		t.Errorf("expected target 'x', got '%s'", ident.Name)
	}
}

func TestParseFuncDecl(t *testing.T) {
	src := "def add(a: int, b: int) -> int:\n    return a + b\n"
	prog := mustParse(t, src)
	fn, ok := prog.Body[0].(*FuncDecl)
	if !ok {
		t.Fatalf("expected FuncDecl, got %T", prog.Body[0])
	}
	if fn.Name != "add" {
		t.Errorf("expected name 'add', got '%s'", fn.Name)
	}
	if len(fn.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(fn.Params))
	}
	if fn.Params[0].Name != "a" || fn.Params[0].TypeName != "int" {
		t.Errorf("param 0: expected a: int, got %s: %s", fn.Params[0].Name, fn.Params[0].TypeName)
	}
	if fn.ReturnType != "int" {
		t.Errorf("expected return type 'int', got '%s'", fn.ReturnType)
	}
	if len(fn.Body) != 1 {
		t.Fatalf("expected 1 statement in body, got %d", len(fn.Body))
	}
}

func TestParseVoidFunc(t *testing.T) {
	src := "def sayHi() -> void:\n    print(\"Hi\")\n"
	prog := mustParse(t, src)
	fn := prog.Body[0].(*FuncDecl)
	if fn.ReturnType != "void" {
		t.Errorf("expected return type 'void', got '%s'", fn.ReturnType)
	}
}

func TestParseIfStmt(t *testing.T) {
	src := "if x > 0:\n    print(\"positive\")\nelif x == 0:\n    print(\"zero\")\nelse:\n    print(\"negative\")\n"
	prog := mustParse(t, src)
	ifStmt, ok := prog.Body[0].(*IfStmt)
	if !ok {
		t.Fatalf("expected IfStmt, got %T", prog.Body[0])
	}
	if len(ifStmt.Elifs) != 1 {
		t.Errorf("expected 1 elif, got %d", len(ifStmt.Elifs))
	}
	if ifStmt.ElseBody == nil {
		t.Error("expected else body")
	}
}

func TestParseWhileStmt(t *testing.T) {
	src := "while i < 10:\n    i = i + 1\n"
	prog := mustParse(t, src)
	ws, ok := prog.Body[0].(*WhileStmt)
	if !ok {
		t.Fatalf("expected WhileStmt, got %T", prog.Body[0])
	}
	if len(ws.Body) != 1 {
		t.Errorf("expected 1 body statement, got %d", len(ws.Body))
	}
}

func TestParseForStmt(t *testing.T) {
	src := "for x: int in nums:\n    print(x)\n"
	prog := mustParse(t, src)
	fs, ok := prog.Body[0].(*ForStmt)
	if !ok {
		t.Fatalf("expected ForStmt, got %T", prog.Body[0])
	}
	if fs.VarName != "x" {
		t.Errorf("expected var 'x', got '%s'", fs.VarName)
	}
	if fs.VarType != "int" {
		t.Errorf("expected var type 'int', got '%s'", fs.VarType)
	}
}

func TestParseFuncCall(t *testing.T) {
	src := `print("hello")` + "\n"
	prog := mustParse(t, src)
	exprStmt, ok := prog.Body[0].(*ExpressionStmt)
	if !ok {
		t.Fatalf("expected ExpressionStmt, got %T", prog.Body[0])
	}
	call, ok := exprStmt.Expr.(*FuncCall)
	if !ok {
		t.Fatalf("expected FuncCall, got %T", exprStmt.Expr)
	}
	if call.Name != "print" {
		t.Errorf("expected 'print', got '%s'", call.Name)
	}
}

func TestParseBinaryExpr(t *testing.T) {
	prog := mustParse(t, "x: int = 1 + 2 * 3")
	decl := prog.Body[0].(*VarDecl)
	// Should parse as 1 + (2 * 3) due to precedence
	binExpr, ok := decl.Value.(*BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", decl.Value)
	}
	if binExpr.Operator != "+" {
		t.Errorf("expected '+', got '%s'", binExpr.Operator)
	}
	right, ok := binExpr.Right.(*BinaryExpr)
	if !ok {
		t.Fatalf("expected right to be BinaryExpr, got %T", binExpr.Right)
	}
	if right.Operator != "*" {
		t.Errorf("expected '*', got '%s'", right.Operator)
	}
}

func TestParseUnaryExpr(t *testing.T) {
	prog := mustParse(t, "x: int = -5")
	decl := prog.Body[0].(*VarDecl)
	unary, ok := decl.Value.(*UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr, got %T", decl.Value)
	}
	if unary.Operator != "-" {
		t.Errorf("expected '-', got '%s'", unary.Operator)
	}
}

func TestParseMethodCall(t *testing.T) {
	prog := mustParse(t, `x: int = msg.length()`)
	decl := prog.Body[0].(*VarDecl)
	mc, ok := decl.Value.(*MethodCall)
	if !ok {
		t.Fatalf("expected MethodCall, got %T", decl.Value)
	}
	if mc.Method != "length" {
		t.Errorf("expected method 'length', got '%s'", mc.Method)
	}
}

func TestParseIndexAccess(t *testing.T) {
	prog := mustParse(t, "x: int = nums[0]")
	decl := prog.Body[0].(*VarDecl)
	idx, ok := decl.Value.(*IndexAccess)
	if !ok {
		t.Fatalf("expected IndexAccess, got %T", decl.Value)
	}
	obj := idx.Object.(*Identifier)
	if obj.Name != "nums" {
		t.Errorf("expected 'nums', got '%s'", obj.Name)
	}
}

func TestParseBreakContinue(t *testing.T) {
	src := "while true:\n    break\n    continue\n"
	prog := mustParse(t, src)
	ws := prog.Body[0].(*WhileStmt)
	if _, ok := ws.Body[0].(*BreakStmt); !ok {
		t.Errorf("expected BreakStmt, got %T", ws.Body[0])
	}
	if _, ok := ws.Body[1].(*ContinueStmt); !ok {
		t.Errorf("expected ContinueStmt, got %T", ws.Body[1])
	}
}

func TestParseLogicalExpr(t *testing.T) {
	prog := mustParse(t, "x: bool = a and b or c")
	decl := prog.Body[0].(*VarDecl)
	// Should parse as (a and b) or c
	bin, ok := decl.Value.(*BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", decl.Value)
	}
	if bin.Operator != "or" {
		t.Errorf("expected 'or', got '%s'", bin.Operator)
	}
}

func TestParseNotExpr(t *testing.T) {
	prog := mustParse(t, "x: bool = not true")
	decl := prog.Body[0].(*VarDecl)
	unary, ok := decl.Value.(*UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr, got %T", decl.Value)
	}
	if unary.Operator != "not" {
		t.Errorf("expected 'not', got '%s'", unary.Operator)
	}
}

func TestParseGroupedExpr(t *testing.T) {
	prog := mustParse(t, "x: int = (1 + 2) * 3")
	decl := prog.Body[0].(*VarDecl)
	bin := decl.Value.(*BinaryExpr)
	if bin.Operator != "*" {
		t.Errorf("expected '*' at top level, got '%s'", bin.Operator)
	}
	left := bin.Left.(*BinaryExpr)
	if left.Operator != "+" {
		t.Errorf("expected '+' in group, got '%s'", left.Operator)
	}
}
