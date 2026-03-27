package main

import (
	"github.com/johnivanpuayap/lexo/wasm"
)

func main() {
	wasm.RegisterCallbacks()

	// Keep the Go program alive
	select {}
}
