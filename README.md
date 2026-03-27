# Lexo

A statically-typed, Python-like programming language with a web-based IDE, designed to teach beginners the fundamentals of static typing.

## What is Lexo?

Lexo bridges the gap between dynamically-typed languages like Python and statically-typed languages like Go, Java, and TypeScript. It uses familiar Python-like syntax with explicit type annotations, friendly error messages, and an interactive web IDE — all running entirely in the browser.

```
# Every variable has a type
name: string = "World"
count: int = 5

def greet(name: string, times: int) -> void:
    i: int = 0
    while i < times:
        print("Hello, " + name + "!")
        i = i + 1

greet(name, count)
```

## Features

**Language**
- Python-like syntax with explicit static types
- Types: `int`, `float`, `string`, `bool`, and arrays (`int[]`, `string[]`, etc.)
- Functions with typed parameters and return types
- Control flow: `if`/`elif`/`else`, `while`, `for...in`, `break`, `continue`
- String methods: `.length()`, `.upper()`, `.lower()`, `.substring()`
- No implicit type coercion — teaches you to think about types

**IDE**
- Monaco Editor (VS Code engine) with Lexo syntax highlighting
- Catppuccin Mocha dark theme
- Output console with `print()` and `input()` support
- Inline error highlighting
- Built-in language tutorial

**Architecture**
- Go interpreter compiled to WebAssembly — runs entirely in the browser
- Pipeline: Lexer → Parser → Type Checker → Tree-Walking Interpreter
- React/TypeScript frontend with Vite
- No backend required — deploys as a static site

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.21+
- [Node.js](https://nodejs.org/) 20+

### Build & Run

```bash
# Clone the repo
git clone https://github.com/johnivanpuayap/lexo.git
cd lexo

# Build the WASM binary
make build

# Install frontend dependencies and start dev server
cd web
npm install
npm run dev
```

Open **http://localhost:5173** in your browser.

### Run Tests

```bash
# Go interpreter tests
go test ./pkg/... -v

# Build frontend for production
cd web && npm run build
```

## Project Structure

```
lexo/
├── cmd/wasm/           # WASM entry point
├── pkg/
│   ├── lexer/          # Tokenizer with INDENT/DEDENT
│   ├── parser/         # Recursive descent parser
│   ├── checker/        # Static type checker
│   └── interpreter/    # Tree-walking interpreter
├── wasm/               # JS ↔ Go bridge (syscall/js)
├── programs/           # Example .lexo programs
└── Makefile

web/
├── src/
│   ├── components/     # React UI components
│   ├── hooks/          # WASM loader
│   └── App.tsx         # Main app
├── public/
│   ├── tutorial.html   # Language documentation
│   └── wasm_exec.js    # Go WASM runtime
└── package.json
```

## Example Programs

See the [`programs/`](programs/) directory for examples:

- [`hello.lexo`](programs/hello.lexo) — Hello World
- [`variables.lexo`](programs/variables.lexo) — All type declarations
- [`functions.lexo`](programs/functions.lexo) — Functions and recursion
- [`loops.lexo`](programs/loops.lexo) — While and for loops
- [`arrays.lexo`](programs/arrays.lexo) — Array operations

## License

MIT
