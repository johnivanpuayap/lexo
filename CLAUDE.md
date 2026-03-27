# Lexo

A statically-typed, Python-like programming language with a web-based IDE, designed to teach beginners the fundamentals of static typing.

## Project Status

- Design spec: `docs/specs/2026-03-26-lexo-design.md`
- Phase 1 implementation: **Complete** — interpreter + IDE working end-to-end
- Phase 2 (classes, imports, enums, generics): Not started

## Key Decisions

- **Language name:** Lexo (file extension: `.lexo`)
- **Syntax:** Python-like with explicit static types, `def` for functions, no `let` keyword for variables (`name: type = value`)
- **Tech stack:** Go interpreter compiled to WebAssembly + React/TypeScript IDE with Monaco Editor
- **Target audience:** Complete beginners learning static typing
- **Scope (Phase 1):** Variables, types (int/float/string/bool/arrays), control flow, functions, string operations, I/O
- **Scope (Phase 2):** Classes/structs, imports/modules, enums, generics
- **IDE features:** Syntax highlighting, variable inspector, friendly error messages, Catppuccin Mocha theme
- **Architecture:** Go WASM (lexer → parser → type checker → interpreter) ↔ WASM bridge (syscall/js) ↔ React frontend
- **Deployment:** Static site (no backend needed)

## Build Commands

```bash
# Run all Go tests
go test ./pkg/... -v

# Build WASM binary
make build

# Run frontend dev server
cd web && npm run dev

# Build frontend for production
cd web && npm run build
```

## Project Structure

```
lexo/                           # Go module (github.com/johnivanpuayap/lexo)
├── cmd/wasm/main.go            # WASM entry point
├── pkg/
│   ├── lexer/                  # Tokenizer with INDENT/DEDENT
│   │   ├── token.go            # Token types, Token struct, keywords
│   │   ├── lexer.go            # Lexer implementation
│   │   └── lexer_test.go
│   ├── parser/                 # Recursive descent parser
│   │   ├── ast.go              # AST node types
│   │   ├── parser.go           # Parser implementation
│   │   └── parser_test.go
│   ├── checker/                # Static type checker
│   │   ├── types.go            # LexoType enum, type utilities
│   │   ├── checker.go          # Type checker with scope tracking
│   │   └── checker_test.go
│   └── interpreter/            # Tree-walking interpreter
│       ├── value.go            # Runtime value types (IntVal, StringVal, etc.)
│       ├── environment.go      # Scope chain
│       ├── interpreter.go      # Interpreter with IOHandler interface
│       └── interpreter_test.go
├── wasm/bridge.go              # JS ↔ Go callback bridge (syscall/js)
├── programs/                   # Example .lexo files
├── go.mod
└── Makefile

web/                            # React IDE (Vite + TypeScript)
├── src/
│   ├── components/
│   │   ├── TopBar.tsx
│   │   ├── Editor.tsx          # Monaco Editor with Lexo highlighting
│   │   ├── OutputConsole.tsx
│   │   ├── VariableInspector.tsx
│   │   └── lexo-lang.ts        # Monarch tokenizer + Catppuccin theme
│   ├── hooks/useLexo.ts        # WASM loader + bridge wrapper
│   ├── types/lexo.d.ts         # Bridge TypeScript types
│   ├── App.tsx                 # Main app layout
│   └── App.css                 # Catppuccin Mocha styles
├── public/
│   ├── lexo.wasm               # Built WASM binary (gitignored)
│   └── wasm_exec.js            # Go WASM runtime
└── package.json
```

## Interpreter Pipeline

```
Source Code → Lexer (tokens) → Parser (AST) → Type Checker → Interpreter → Output
                                                    ↕
                                        WASM Bridge (syscall/js)
                                                    ↕
                                          React IDE (Monaco)
```

## Testing

The interpreter has comprehensive Go tests:
- **Lexer:** Token types, INDENT/DEDENT, comments, strings, operators, line tracking
- **Parser:** All statement/expression types, operator precedence, error cases
- **Checker:** Type mismatches, undeclared vars, function signatures, method calls, arrays
- **Interpreter:** Arithmetic, control flow, functions, recursion, built-ins, runtime errors
