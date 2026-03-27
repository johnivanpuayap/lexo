# Lexo — Design Specification

A statically-typed, Python-like programming language with a web-based IDE, designed to teach beginners the fundamentals of static typing.

## 1. Overview

**Lexo** is an educational programming language and web-based IDE that bridges the gap between dynamically-typed languages (like Python) and statically-typed languages (like Java, Go, C#). It uses familiar Python-like syntax with explicit type annotations, friendly error messages, and a "type journey" visualization that helps learners understand why types matter.

### Goals

- Teach static typing concepts through a beginner-friendly syntax
- Provide an interactive web IDE with a step-through debugger
- Show learners *why* type errors happen, not just *that* they happened
- Run entirely in the browser with no backend required

### Target Audience

- Complete programming beginners learning typing concepts from day one
- Python/JS developers transitioning to statically-typed languages

### Scope (Moderate — Phase 1)

- Variables with explicit types
- Basic types: `int`, `float`, `string`, `bool`
- Arrays: `int[]`, `string[]`, etc.
- Control flow: `if`/`elif`/`else`, `while`, `for...in`
- Functions with typed parameters and return types
- String operations: `length()`, `upper()`, `lower()`, `substring()`
- I/O: `print()`, `input()`
- Comments: `#`

### Future Expansion (Phase 2)

- Classes/structs
- Imports/modules
- Enums
- Generics

## 2. Language Syntax

### Variable Declarations

```
name: string = "Lexo"
age: int = 25
pi: float = 3.14
active: bool = true
```

No `let` keyword — just `name: type = value`. The type annotation is always required.

### Functions

```
def greet(name: string) -> string:
    return "Hello, " + name

def add(a: int, b: int) -> int:
    return a + b

def sayHi() -> void:
    print("Hi!")
```

- `def` keyword (familiar from Python)
- Parameters are typed: `name: type`
- Return type after `->` arrow
- `void` for functions that return nothing
- Indentation-based blocks (Python-style, 4 spaces)

### Control Flow

```
if age >= 18:
    print("Adult")
elif age >= 13:
    print("Teen")
else:
    print("Child")
```

### Loops

```
# For-in loop (iterating arrays)
scores: int[] = [90, 85, 92]
for score: int in scores:
    print(score)

# While loop
i: int = 0
while i < 10:
    print(i)
    i = i + 1
```

The `for` loop variable must include a type annotation (`score: int`) to stay consistent with the "always explicit types" philosophy. The type checker verifies this matches the array's element type.

`break` and `continue` are supported inside loops.

### Arrays

```
numbers: int[] = [1, 2, 3, 4, 5]
names: string[] = ["Alice", "Bob"]

print(numbers[0])        # Indexing
print(numbers.length())  # Length
```

### String Operations

```
msg: string = "hello"
print(msg.length())       # 5
print(msg.upper())        # "HELLO"
print(msg.lower())        # "hello"
print(msg.substring(0, 3)) # "hel"
```

### I/O

```
print("What is your name?")
name: string = input("Name: ")
print("Hello, " + name)
```

`input()` always returns a `string`. The user must explicitly convert if they need another type (future: built-in `toInt()`, `toFloat()` functions).

### Operators

- **Arithmetic:** `+`, `-`, `*`, `/`, `%`
- **Comparison:** `==`, `!=`, `<`, `>`, `<=`, `>=`
- **Logical:** `and`, `or`, `not` (lowercase keywords)
- **String concatenation:** `+`

### Comments

```
# This is a comment
```

Single-line only. No multi-line comments in Phase 1.

### Keywords

`def`, `return`, `if`, `elif`, `else`, `while`, `for`, `in`, `break`, `continue`, `and`, `or`, `not`, `true`, `false`, `void`

### Type System

| Type | Description | Literal Examples |
|------|-------------|-----------------|
| `int` | Whole numbers | `0`, `42`, `-7` |
| `float` | Decimal numbers | `3.14`, `-0.5` |
| `string` | Text | `"hello"`, `"it's"` |
| `bool` | Boolean | `true`, `false` |
| `int[]` | Array of ints | `[1, 2, 3]` |
| `float[]` | Array of floats | `[1.0, 2.5]` |
| `string[]` | Array of strings | `["a", "b"]` |
| `bool[]` | Array of bools | `[true, false]` |
| `void` | No return value | (functions only) |

No implicit type coercion. `int` + `float` is a type error — the user must be explicit. This is intentional for teaching purposes.

## 3. Error Messages

All errors are written in plain English with line/column references and contextual explanations.

### Example: Type Mismatch

```
Line 3: Type mismatch in assignment

  1 | count: int = 10
  2 | label: string = "items"
> 3 | count = "five"

You declared `count` as an `int` on line 1, which means it can only hold
whole numbers like 10, -3, or 0. On line 3, you're trying to assign "five",
which is a `string` (text). An `int` variable can't hold text — they're
different types.
```

### Example: Undeclared Variable

```
Line 5: Variable not found

> 5 | print(naem)

There is no variable called `naem`. Did you mean `name` (declared on line 1)?
```

### Error Categories

- **Lexer errors:** Unexpected characters, unterminated strings, invalid numbers
- **Parse errors:** Syntax mistakes, unexpected tokens, indentation issues
- **Type errors:** Type mismatches, undeclared variables, wrong argument types, wrong return type (caught before execution)
- **Runtime errors:** Division by zero, array index out of bounds, stack overflow

## 4. Architecture

### Tech Stack

- **Interpreter:** Go, compiled to WebAssembly
- **IDE Frontend:** React + TypeScript + Vite
- **Code Editor:** Monaco Editor (VS Code engine)
- **Deployment:** Static site (Vercel, Netlify, GitHub Pages)

### Interpreter Pipeline

```
Source Code → Lexer → Parser → Type Checker → Interpreter → Output
                                    ↓
                            Type Journey Data
```

All four stages run in the browser via WASM.

#### Stage 1: Lexer (`pkg/lexer`)

- Character-by-character tokenization
- Produces tokens with type, value, line, and column
- Handles INDENT/DEDENT tokens for Python-style indentation
- Recognizes all keywords, operators, and literals

#### Stage 2: Parser (`pkg/parser`)

- Recursive descent parser
- Produces an Abstract Syntax Tree (AST)
- AST node types:
  - **Program** — root node, list of statements
  - **VarDecl, Assignment** — variable operations
  - **FuncDecl, ReturnStmt, FuncCall** — function operations
  - **If, While, For, Break, Continue** — control flow
  - **BinaryExpr, UnaryExpr, Literal, Identifier** — expressions
  - **Print, Input** — I/O
  - **ArrayLiteral, IndexAccess** — arrays

#### Stage 3: Type Checker (`pkg/checker`)

- Walks the AST and validates all types before execution
- Builds a symbol table with scope tracking
- Collects **type journey data** — records where each variable was declared and every assignment attempt with the type used
- Checks performed:
  - Variable declared before use
  - No duplicate declarations in the same scope
  - Assignment type matches declaration type
  - Function return type matches declared return type
  - Argument types match parameter types
  - Array element types are consistent
  - Condition expressions evaluate to `bool`

#### Stage 4: Interpreter (`pkg/interpreter`)

- Tree-walking interpreter that directly executes the AST
- Scope chain with lexical scoping (global scope + function scopes)
- Call stack for function calls
- **Debugger hooks:** Emits events to the IDE via WASM bridge callbacks
  - `onStep` — about to execute a statement
  - `onVariableUpdate` — a variable's value changed
  - `onScopeEnter` / `onScopeExit` — entering/leaving a function
- Async execution using Go channels — the debugger can pause the interpreter between statements by blocking on a channel, resumed when the user clicks "Step"
- Configurable execution speed for a "slow run" mode

### WASM Bridge (`wasm/bridge`)

Uses Go's `syscall/js` package to register functions on the global JS object.

**React → Go (function calls):**

| Function | Description |
|----------|-------------|
| `lexo.run(source)` | Execute program, return output |
| `lexo.debugStart(source)` | Start debug session |
| `lexo.debugStep()` | Step to next statement |
| `lexo.debugStepInto()` | Step into function call |
| `lexo.debugStepOut()` | Step out of current function |
| `lexo.debugStop()` | Stop debug session |
| `lexo.getVariables()` | Get current variable state |
| `lexo.getTypeJourney(name)` | Get type history for a variable |

**Go → React (callbacks):**

| Callback | Description |
|----------|-------------|
| `onOutput(text)` | `print()` produced output |
| `onError(error)` | Friendly error message |
| `onDebugPause(line, vars)` | Debugger paused at a line |
| `onVariableUpdate(name, value)` | Variable value changed |
| `onScopeEnter(name)` | Entered a function scope |
| `onScopeExit(name)` | Left a function scope |
| `onRequestInput(prompt)` | `input()` needs user input |
| `onComplete(result)` | Execution finished |

### Debugger Design

The debugger uses a channel-based approach in Go:

1. The interpreter runs in a goroutine
2. Before each statement, it sends the current state to a `pause` channel and blocks
3. When the user clicks "Step" in the IDE, React calls `lexo.debugStep()`
4. The bridge sends a signal to a `resume` channel
5. The interpreter unblocks and executes one statement
6. Repeat

This gives precise statement-level control without needing to restructure the interpreter.

## 5. IDE Design

### Layout

```
┌─────────────────────────────────────────────────────────┐
│ Lexo  │ untitled.lexo                    [▶ Run] [⏵ Debug] │
├────────────────────────────────┬────────────────────────┤
│                                │  [Variables] [Call Stack] [Type Journey] │
│   Monaco Code Editor           │                        │
│   - Syntax highlighting        │   Variable Inspector   │
│   - Line numbers               │   - Grouped by scope   │
│   - Debug line highlight       │   - Name, type, value  │
│   - Inline error markers       │   - Update indicators  │
│                                │                        │
├────────────────────────────────┤                        │
│ DEBUG: [⏸ Pause] [→ Step]     │   Type Journey Preview │
│        [↓ Into] [↑ Out] [■ Stop] │   - Declaration point  │
├────────────────────────────────┤   - Assignment history  │
│ Output Console                 │                        │
│ - print() output               │                        │
│ - input() prompts              │                        │
│ - Error messages               │                        │
└────────────────────────────────┴────────────────────────┘
```

### Components

| Component | Description |
|-----------|-------------|
| **TopBar** | File name, Run and Debug buttons |
| **Editor** | Monaco Editor with Lexo syntax highlighting, breakpoint gutters, debug line highlighting, inline error decorations |
| **DebugControls** | Pause, Step, Step Into, Step Out, Stop buttons — visible only during debug sessions |
| **OutputConsole** | Shows `print()` output, `input()` prompts (inline text field), error messages with line references, execution time |
| **VariableInspector** | Variables grouped by scope (Global, Function scopes). Each variable shows name, type badge, current value, and an update indicator when the value changes |
| **CallStack** | Shows the current function call chain during debugging |
| **TypeJourney** | For a selected variable, shows a timeline: where it was declared, every assignment with the type used, and whether each assignment was valid |

### Monaco Editor Configuration

- Custom Lexo language definition for syntax highlighting (keywords, types, operators, strings, comments)
- Custom theme matching the Catppuccin Mocha color scheme
- Error decorations: red underline + hover tooltip with the friendly error message
- Debug decorations: highlighted current line (blue background), breakpoint dots in gutter

## 6. Project Structure

```
lexo/
├── cmd/wasm/
│   └── main.go              # WASM entry point
├── pkg/lexer/
│   ├── lexer.go
│   ├── token.go
│   └── lexer_test.go
├── pkg/parser/
│   ├── parser.go
│   ├── ast.go                # AST node types
│   └── parser_test.go
├── pkg/checker/
│   ├── checker.go
│   ├── types.go
│   ├── journey.go            # Type journey data collection
│   └── checker_test.go
├── pkg/interpreter/
│   ├── interpreter.go
│   ├── environment.go        # Scope chain
│   ├── debugger.go           # Debug hooks + channel control
│   └── interpreter_test.go
├── wasm/
│   └── bridge.go             # JS ↔ Go interface via syscall/js
├── go.mod
└── Makefile

web/
├── src/
│   ├── components/
│   │   ├── Editor.tsx         # Monaco Editor wrapper
│   │   ├── OutputConsole.tsx
│   │   ├── VariableInspector.tsx
│   │   ├── TypeJourney.tsx
│   │   ├── CallStack.tsx
│   │   ├── DebugControls.tsx
│   │   └── TopBar.tsx
│   ├── hooks/
│   │   ├── useLexo.ts         # WASM loader + bridge wrapper
│   │   └── useDebugger.ts     # Debug state management
│   ├── wasm/
│   │   ├── lexo.wasm          # Built artifact (gitignored)
│   │   └── wasm_exec.js       # Go WASM runtime helper
│   ├── types/
│   │   └── lexo.d.ts          # TypeScript types for the bridge API
│   └── App.tsx
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 7. Build & Deploy

```bash
# Build Go interpreter to WASM
GOOS=js GOARCH=wasm go build -o web/src/wasm/lexo.wasm ./cmd/wasm

# Copy Go's WASM runtime helper
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" web/src/wasm/

# Build React frontend
cd web && npm run build

# Output: static files in web/dist/ — deploy to Vercel, Netlify, GitHub Pages
```

The final output is a static site (HTML + JS + WASM). No server required.

## 8. Testing Strategy

### Go (Interpreter)

- **Unit tests per package:** Each stage (lexer, parser, checker, interpreter) has its own test file
- **Lexer tests:** Verify token output for various inputs
- **Parser tests:** Verify AST structure for valid programs, verify error messages for invalid syntax
- **Type checker tests:** Verify type errors are caught, verify type journey data is collected correctly
- **Interpreter tests:** Verify correct output for sample programs, verify runtime errors are caught
- **Integration tests:** Feed complete `.lexo` programs through the full pipeline and verify output

### React (IDE)

- **Component tests:** Verify each panel renders correctly with mock data
- **WASM integration tests:** Verify the bridge API works end-to-end in a browser environment

### Sample Programs

A `programs/` directory with example `.lexo` files that serve as both tests and documentation:
- `hello.lexo` — basic print
- `variables.lexo` — all type declarations
- `functions.lexo` — function definitions and calls
- `loops.lexo` — for and while loops
- `arrays.lexo` — array operations
- `errors/` — programs that intentionally trigger each error type
