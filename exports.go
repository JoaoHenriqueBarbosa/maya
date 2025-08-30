//go:build wasm
// +build wasm

package maya

import (
	"fmt"
	"syscall/js"
	"github.com/maya-framework/maya/internal/render"
)

// domReadyCallbacks stores callbacks for DOM ready events
var domReadyCallbacks []func()
var domReady = false

// onDOMReady is exported to JavaScript for DOM ready notification
//
//go:wasmexport onDOMReady
func onDOMReady() {
	domReady = true
	for _, callback := range domReadyCallbacks {
		callback()
	}
	domReadyCallbacks = nil
}

// handleEvent is exported to JavaScript to handle events
//
//go:wasmexport handleEvent  
func handleEvent(callbackID int32) {
	fmt.Printf("handleEvent called with ID: %d\n", callbackID)
	if callback := render.GetCallback(callbackID); callback != nil {
		callback()
	}
}

// handleButtonClick is exported for button clicks
//
//go:wasmexport handleButtonClick
func handleButtonClick(buttonID string) {
	fmt.Printf("handleButtonClick called with ID: %s\n", buttonID)
}

// waitForDOM waits for the DOM to be ready (without js.FuncOf)
func waitForDOM(callback func()) {
	if domReady {
		callback()
		return
	}
	
	doc := js.Global().Get("document")
	readyState := doc.Get("readyState").String()
	
	if readyState == "complete" || readyState == "interactive" {
		domReady = true
		callback()
	} else {
		// Store callback to be called when DOM is ready
		domReadyCallbacks = append(domReadyCallbacks, callback)
		
		// Setup listener in JavaScript
		// The JavaScript side will call wasmExports.onDOMReady when ready
		// This is set up in wasm_init_correct.js
	}
}