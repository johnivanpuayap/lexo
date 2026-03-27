package checker

import (
	"strings"
	"testing"

	"github.com/johnivanpuayap/lexo/pkg/lexer"
	"github.com/johnivanpuayap/lexo/pkg/parser"
)

func mustCheck(t *testing.T, source string) {
	t.Helper()
	tokens, err := lexer.Tokenize(source)
	if err != nil {
		t.Fatalf("lexer error: %v", err)
	}
	prog, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
	if err := Check(prog); err != nil {
		t.Fatalf("checker error: %v", err)
	}
}

func mustFailCheck(t *testing.T, source string, expectedMsg string) {
	t.Helper()
	tokens, err := lexer.Tokenize(source)
	if err != nil {
		t.Fatalf("lexer error: %v", err)
	}
	prog, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
	err = Check(prog)
	if err == nil {
		t.Fatal("expected type error, got none")
	}
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Fatalf("expected error containing %q, got: %v", expectedMsg, err)
	}
}

func TestValidVarDecl(t *testing.T) {
	mustCheck(t, `x: int = 42`)
	mustCheck(t, `name: string = "hello"`)
	mustCheck(t, `pi: float = 3.14`)
	mustCheck(t, `flag: bool = true`)
}

func TestTypeMismatchVarDecl(t *testing.T) {
	mustFailCheck(t, `x: int = "hello"`, "type mismatch")
	mustFailCheck(t, `x: string = 42`, "type mismatch")
	mustFailCheck(t, `x: bool = 5`, "type mismatch")
}

func TestTypeMismatchAssignment(t *testing.T) {
	mustFailCheck(t, "x: int = 1\nx = \"hello\"\n", "type mismatch")
}

func TestUndeclaredVariable(t *testing.T) {
	mustFailCheck(t, `print(x)`, "not declared")
}

func TestDuplicateDeclaration(t *testing.T) {
	mustFailCheck(t, "x: int = 1\nx: int = 2\n", "already declared")
}

func TestFunctionTypeChecking(t *testing.T) {
	mustCheck(t, "def add(a: int, b: int) -> int:\n    return a + b\n")
}

func TestFunctionWrongReturnType(t *testing.T) {
	mustFailCheck(t, "def foo() -> int:\n    return \"hello\"\n", "return type")
}

func TestFunctionWrongArgType(t *testing.T) {
	src := "def greet(name: string) -> void:\n    print(name)\ngreet(42)\n"
	mustFailCheck(t, src, "argument")
}

func TestFunctionWrongArgCount(t *testing.T) {
	src := "def add(a: int, b: int) -> int:\n    return a + b\nadd(1)\n"
	mustFailCheck(t, src, "argument")
}

func TestArrayTypeChecking(t *testing.T) {
	mustCheck(t, `nums: int[] = [1, 2, 3]`)
	mustCheck(t, `names: string[] = ["a", "b"]`)
}

func TestArrayTypeMismatch(t *testing.T) {
	mustFailCheck(t, `nums: int[] = ["a", "b"]`, "type mismatch")
}

func TestForLoopTypeChecking(t *testing.T) {
	mustCheck(t, "nums: int[] = [1, 2, 3]\nfor n: int in nums:\n    print(n)\n")
}

func TestForLoopWrongVarType(t *testing.T) {
	mustFailCheck(t, "nums: int[] = [1, 2, 3]\nfor n: string in nums:\n    print(n)\n", "type")
}

func TestConditionMustBeBool(t *testing.T) {
	mustFailCheck(t, "if 42:\n    print(\"hi\")\n", "bool")
}

func TestArithmeticTypes(t *testing.T) {
	mustCheck(t, "x: int = 1 + 2")
	mustCheck(t, "x: float = 1.0 + 2.0")
	mustCheck(t, `x: string = "a" + "b"`)
}

func TestArithmeticTypeMismatch(t *testing.T) {
	mustFailCheck(t, "x: int = 1 + 1.0", "type mismatch")
	mustFailCheck(t, `x: int = 1 + "a"`, "type mismatch")
}

func TestComparisonTypes(t *testing.T) {
	mustCheck(t, "x: bool = 1 < 2")
	mustCheck(t, "x: bool = 1 == 2")
}

func TestMethodCallTypes(t *testing.T) {
	mustCheck(t, `x: int = "hello".length()`)
	mustCheck(t, `x: string = "hello".upper()`)
	mustCheck(t, `x: string = "hello".substring(0, 3)`)
}

func TestIndexAccessType(t *testing.T) {
	mustCheck(t, "nums: int[] = [1, 2]\nx: int = nums[0]\n")
}

func TestIndexOnNonArray(t *testing.T) {
	mustFailCheck(t, "x: int = 5\ny: int = x[0]\n", "cannot index")
}

func TestPrintAcceptsAnyType(t *testing.T) {
	mustCheck(t, `print(42)`)
	mustCheck(t, `print("hello")`)
	mustCheck(t, `print(true)`)
}

func TestInputReturnsString(t *testing.T) {
	mustCheck(t, `name: string = input("Name: ")`)
	mustFailCheck(t, `x: int = input("X: ")`, "type mismatch")
}

func TestBreakOutsideLoop(t *testing.T) {
	mustFailCheck(t, "break\n", "outside")
}

func TestContinueOutsideLoop(t *testing.T) {
	mustFailCheck(t, "continue\n", "outside")
}

func TestBreakInsideLoopOK(t *testing.T) {
	mustCheck(t, "while true:\n    break\n")
}

func TestArrayLengthMethod(t *testing.T) {
	mustCheck(t, "nums: int[] = [1, 2, 3]\nx: int = nums.length()\n")
}
