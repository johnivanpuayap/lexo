package interpreter

import (
	"strings"
	"testing"

	"github.com/johnivanpuayap/lexo/pkg/checker"
	"github.com/johnivanpuayap/lexo/pkg/lexer"
	"github.com/johnivanpuayap/lexo/pkg/parser"
)

type testIO struct {
	output []string
	inputs []string
	idx    int
}

func (t *testIO) Print(text string) { t.output = append(t.output, text) }
func (t *testIO) Input(prompt string) string {
	if t.idx >= len(t.inputs) {
		return ""
	}
	val := t.inputs[t.idx]
	t.idx++
	return val
}

func run(t *testing.T, source string) *testIO {
	t.Helper()
	return runWithInput(t, source, nil)
}

func runWithInput(t *testing.T, source string, inputs []string) *testIO {
	t.Helper()
	tokens, err := lexer.Tokenize(source)
	if err != nil {
		t.Fatalf("lexer error: %v", err)
	}
	prog, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
	if err := checker.Check(prog); err != nil {
		t.Fatalf("checker error: %v", err)
	}
	io := &testIO{inputs: inputs}
	interp := New(io)
	if err := interp.Execute(prog); err != nil {
		t.Fatalf("runtime error: %v", err)
	}
	return io
}

func mustFail(t *testing.T, source string, expectedMsg string) {
	t.Helper()
	tokens, _ := lexer.Tokenize(source)
	prog, _ := parser.Parse(tokens)
	_ = checker.Check(prog)
	io := &testIO{}
	interp := New(io)
	err := interp.Execute(prog)
	if err == nil {
		t.Fatal("expected runtime error, got none")
	}
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Fatalf("expected error containing %q, got: %v", expectedMsg, err)
	}
}

func TestPrintInt(t *testing.T) {
	io := run(t, `print(42)`)
	if len(io.output) != 1 || io.output[0] != "42" {
		t.Errorf("expected [42], got %v", io.output)
	}
}

func TestPrintString(t *testing.T) {
	io := run(t, `print("hello")`)
	if len(io.output) != 1 || io.output[0] != "hello" {
		t.Errorf("expected [hello], got %v", io.output)
	}
}

func TestVariableDeclaration(t *testing.T) {
	io := run(t, "x: int = 10\nprint(x)\n")
	if io.output[0] != "10" {
		t.Errorf("expected 10, got %s", io.output[0])
	}
}

func TestVariableAssignment(t *testing.T) {
	io := run(t, "x: int = 10\nx = 20\nprint(x)\n")
	if io.output[0] != "20" {
		t.Errorf("expected 20, got %s", io.output[0])
	}
}

func TestArithmetic(t *testing.T) {
	io := run(t, "x: int = 2 + 3 * 4\nprint(x)\n")
	if io.output[0] != "14" {
		t.Errorf("expected 14, got %s", io.output[0])
	}
}

func TestStringConcat(t *testing.T) {
	io := run(t, "x: string = \"hello\" + \" \" + \"world\"\nprint(x)\n")
	if io.output[0] != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", io.output[0])
	}
}

func TestIfTrue(t *testing.T) {
	io := run(t, "if true:\n    print(\"yes\")\n")
	if io.output[0] != "yes" {
		t.Errorf("expected 'yes', got '%s'", io.output[0])
	}
}

func TestIfFalseElse(t *testing.T) {
	io := run(t, "if false:\n    print(\"yes\")\nelse:\n    print(\"no\")\n")
	if io.output[0] != "no" {
		t.Errorf("expected 'no', got '%s'", io.output[0])
	}
}

func TestWhileLoop(t *testing.T) {
	src := "i: int = 0\nwhile i < 3:\n    print(i)\n    i = i + 1\n"
	io := run(t, src)
	expected := []string{"0", "1", "2"}
	if len(io.output) != 3 {
		t.Fatalf("expected 3 outputs, got %d: %v", len(io.output), io.output)
	}
	for i, exp := range expected {
		if io.output[i] != exp {
			t.Errorf("output %d: expected %s, got %s", i, exp, io.output[i])
		}
	}
}

func TestForLoop(t *testing.T) {
	src := "nums: int[] = [10, 20, 30]\nfor n: int in nums:\n    print(n)\n"
	io := run(t, src)
	expected := []string{"10", "20", "30"}
	for i, exp := range expected {
		if io.output[i] != exp {
			t.Errorf("output %d: expected %s, got %s", i, exp, io.output[i])
		}
	}
}

func TestBreak(t *testing.T) {
	src := "i: int = 0\nwhile true:\n    if i == 3:\n        break\n    print(i)\n    i = i + 1\n"
	io := run(t, src)
	if len(io.output) != 3 {
		t.Fatalf("expected 3 outputs, got %d: %v", len(io.output), io.output)
	}
}

func TestContinue(t *testing.T) {
	src := "i: int = 0\nwhile i < 5:\n    i = i + 1\n    if i == 3:\n        continue\n    print(i)\n"
	io := run(t, src)
	// Should print 1, 2, 4, 5 (skip 3)
	expected := []string{"1", "2", "4", "5"}
	if len(io.output) != len(expected) {
		t.Fatalf("expected %d outputs, got %d: %v", len(expected), len(io.output), io.output)
	}
	for i, exp := range expected {
		if io.output[i] != exp {
			t.Errorf("output %d: expected %s, got %s", i, exp, io.output[i])
		}
	}
}

func TestFunctionCall(t *testing.T) {
	src := "def double(x: int) -> int:\n    return x * 2\nprint(double(5))\n"
	io := run(t, src)
	if io.output[0] != "10" {
		t.Errorf("expected 10, got %s", io.output[0])
	}
}

func TestRecursion(t *testing.T) {
	src := `def factorial(n: int) -> int:
    if n <= 1:
        return 1
    return n * factorial(n - 1)
print(factorial(5))
`
	io := run(t, src)
	if io.output[0] != "120" {
		t.Errorf("expected 120, got %s", io.output[0])
	}
}

func TestStringMethods(t *testing.T) {
	io := run(t, "print(\"hello\".length())")
	if io.output[0] != "5" {
		t.Errorf("expected 5, got %s", io.output[0])
	}

	io = run(t, "print(\"hello\".upper())")
	if io.output[0] != "HELLO" {
		t.Errorf("expected HELLO, got %s", io.output[0])
	}

	io = run(t, "print(\"HELLO\".lower())")
	if io.output[0] != "hello" {
		t.Errorf("expected hello, got %s", io.output[0])
	}

	io = run(t, "print(\"hello\".substring(1, 3))")
	if io.output[0] != "el" {
		t.Errorf("expected 'el', got '%s'", io.output[0])
	}
}

func TestArrayLength(t *testing.T) {
	io := run(t, "nums: int[] = [1, 2, 3]\nprint(nums.length())\n")
	if io.output[0] != "3" {
		t.Errorf("expected 3, got %s", io.output[0])
	}
}

func TestArrayIndex(t *testing.T) {
	io := run(t, "nums: int[] = [10, 20, 30]\nprint(nums[1])\n")
	if io.output[0] != "20" {
		t.Errorf("expected 20, got %s", io.output[0])
	}
}

func TestDivisionByZero(t *testing.T) {
	mustFail(t, "x: int = 1 / 0\nprint(x)\n", "division by zero")
}

func TestIndexOutOfBounds(t *testing.T) {
	mustFail(t, "nums: int[] = [1, 2]\nprint(nums[5])\n", "index out of bounds")
}

func TestInput(t *testing.T) {
	io := runWithInput(t, "name: string = input(\"Name: \")\nprint(name)\n", []string{"Alice"})
	if io.output[0] != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", io.output[0])
	}
}

func TestNestedFunction(t *testing.T) {
	src := `def add(a: int, b: int) -> int:
    return a + b
def multiply(a: int, b: int) -> int:
    result: int = 0
    i: int = 0
    while i < b:
        result = add(result, a)
        i = i + 1
    return result
print(multiply(3, 4))
`
	io := run(t, src)
	if io.output[0] != "12" {
		t.Errorf("expected 12, got %s", io.output[0])
	}
}

func TestElifChain(t *testing.T) {
	src := `x: int = 2
if x == 1:
    print("one")
elif x == 2:
    print("two")
elif x == 3:
    print("three")
else:
    print("other")
`
	io := run(t, src)
	if io.output[0] != "two" {
		t.Errorf("expected 'two', got '%s'", io.output[0])
	}
}

func TestBooleanLogic(t *testing.T) {
	io := run(t, "print(true and false)")
	if io.output[0] != "false" {
		t.Errorf("expected 'false', got '%s'", io.output[0])
	}

	io = run(t, "print(true or false)")
	if io.output[0] != "true" {
		t.Errorf("expected 'true', got '%s'", io.output[0])
	}

	io = run(t, "print(not true)")
	if io.output[0] != "false" {
		t.Errorf("expected 'false', got '%s'", io.output[0])
	}
}

func TestFloatArithmetic(t *testing.T) {
	io := run(t, "x: float = 1.5 + 2.5\nprint(x)\n")
	if io.output[0] != "4" {
		t.Errorf("expected '4', got '%s'", io.output[0])
	}
}

func TestModulo(t *testing.T) {
	io := run(t, "print(10 % 3)")
	if io.output[0] != "1" {
		t.Errorf("expected '1', got '%s'", io.output[0])
	}
}

func TestShortCircuitAnd(t *testing.T) {
	// This should NOT cause a runtime error because "and" short-circuits
	// The division by zero should never be evaluated
	io := run(t, "x: bool = false and (1 / 0 == 0)\nprint(x)\n")
	if io.output[0] != "false" {
		t.Errorf("expected 'false', got '%s'", io.output[0])
	}
}

func TestShortCircuitOr(t *testing.T) {
	io := run(t, "x: bool = true or (1 / 0 == 0)\nprint(x)\n")
	if io.output[0] != "true" {
		t.Errorf("expected 'true', got '%s'", io.output[0])
	}
}

func TestComparison(t *testing.T) {
	io := run(t, "print(5 > 3)")
	if io.output[0] != "true" {
		t.Errorf("expected 'true', got '%s'", io.output[0])
	}
	io = run(t, "print(5 == 5)")
	if io.output[0] != "true" {
		t.Errorf("expected 'true', got '%s'", io.output[0])
	}
	io = run(t, "print(5 != 3)")
	if io.output[0] != "true" {
		t.Errorf("expected 'true', got '%s'", io.output[0])
	}
}
