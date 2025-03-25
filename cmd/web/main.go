//go:build wasm

package main

import (
	"mbpe-dyn/internal/web"
	"syscall/js"
)

func main() {
	js.Global().Set("tokenizeWeb", web.WrapTokenizeWeb())

	// Keep the Go runtime alive
	select {}
}
