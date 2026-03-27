package wasm

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/johnivanpuayap/lexo/pkg/checker"
	"github.com/johnivanpuayap/lexo/pkg/interpreter"
	"github.com/johnivanpuayap/lexo/pkg/lexer"
	"github.com/johnivanpuayap/lexo/pkg/parser"
)

type WasmIO struct {
	onOutput       js.Value
	onRequestInput js.Value
	inputChan      chan string
}

func (w *WasmIO) Print(text string) {
	w.onOutput.Invoke(text)
}

func (w *WasmIO) Input(prompt string) string {
	if w.onRequestInput.IsUndefined() || w.onRequestInput.IsNull() || w.inputChan == nil {
		return ""
	}
	w.onRequestInput.Invoke(prompt)
	return <-w.inputChan
}

type ErrorResult struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
	Line    int    `json:"line"`
}

func RegisterCallbacks() {
	lexoObj := js.Global().Get("Object").New()

	lexoObj.Set("run", js.FuncOf(runProgram))
	lexoObj.Set("check", js.FuncOf(checkProgram))
	lexoObj.Set("provideInput", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return nil // will be set per-execution
	}))

	js.Global().Set("lexo", lexoObj)
}

func runProgram(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return errorResult("bridge", "run() requires (source, callbacks)", 0)
	}

	source := args[0].String()
	callbacks := args[1]

	onOutput := callbacks.Get("onOutput")
	onError := callbacks.Get("onError")
	onComplete := callbacks.Get("onComplete")

	// Check for input callback (optional)
	var onRequestInput js.Value
	var inputChan chan string
	if ri := callbacks.Get("onRequestInput"); !ri.IsUndefined() && !ri.IsNull() {
		onRequestInput = ri
		inputChan = make(chan string, 1)

		// Register the input provider
		js.Global().Get("lexo").Set("provideInput", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 {
				inputChan <- args[0].String()
			}
			return nil
		}))
	}

	// Run in a goroutine so we don't block the JS event loop
	go func() {
		// Lexer
		tokens, err := lexer.Tokenize(source)
		if err != nil {
			if !onError.IsUndefined() {
				errJSON := errorResult("lexer", err.Error(), extractLine(err))
				onError.Invoke(errJSON)
			}
			if !onComplete.IsUndefined() {
				onComplete.Invoke(false)
			}
			return
		}

		// Parser
		prog, err := parser.Parse(tokens)
		if err != nil {
			if !onError.IsUndefined() {
				errJSON := errorResult("parser", err.Error(), extractLine(err))
				onError.Invoke(errJSON)
			}
			if !onComplete.IsUndefined() {
				onComplete.Invoke(false)
			}
			return
		}

		// Type Checker
		if err := checker.Check(prog); err != nil {
			if !onError.IsUndefined() {
				errJSON := errorResult("type", err.Error(), extractLine(err))
				onError.Invoke(errJSON)
			}
			if !onComplete.IsUndefined() {
				onComplete.Invoke(false)
			}
			return
		}

		// Interpreter
		io := &WasmIO{
			onOutput:       onOutput,
			onRequestInput: onRequestInput,
			inputChan:      inputChan,
		}
		interp := interpreter.New(io)
		if err := interp.Execute(prog); err != nil {
			if !onError.IsUndefined() {
				errJSON := errorResult("runtime", err.Error(), extractLine(err))
				onError.Invoke(errJSON)
			}
			if !onComplete.IsUndefined() {
				onComplete.Invoke(false)
			}
			return
		}

		if !onComplete.IsUndefined() {
			onComplete.Invoke(true)
		}
	}()

	return nil
}

func checkProgram(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return errorResult("bridge", "check() requires (source)", 0)
	}

	source := args[0].String()

	tokens, err := lexer.Tokenize(source)
	if err != nil {
		return errorResult("lexer", err.Error(), extractLine(err))
	}

	prog, err := parser.Parse(tokens)
	if err != nil {
		return errorResult("parser", err.Error(), extractLine(err))
	}

	if err := checker.Check(prog); err != nil {
		return errorResult("type", err.Error(), extractLine(err))
	}

	return nil // no errors
}

func extractLine(err error) int {
	msg := err.Error()
	prefix := "Line "
	idx := strings.Index(msg, prefix)
	if idx < 0 {
		return 0
	}
	numStr := ""
	for i := idx + len(prefix); i < len(msg); i++ {
		if msg[i] >= '0' && msg[i] <= '9' {
			numStr += string(msg[i])
		} else {
			break
		}
	}
	if numStr == "" {
		return 0
	}
	var n int
	fmt.Sscanf(numStr, "%d", &n)
	return n
}

func errorResult(kind, message string, line int) string {
	result := ErrorResult{Kind: kind, Message: message, Line: line}
	b, _ := json.Marshal(result)
	return string(b)
}
