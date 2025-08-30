//go:build wasm
// +build wasm

package maya

import (
	"syscall/js"
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
		
		// Set up JavaScript side to call our exported function
		js.Global().Call("eval", `
			if (document.readyState === "complete" || document.readyState === "interactive") {
				window.onDOMReady();
			} else {
				document.addEventListener("DOMContentLoaded", function() {
					window.onDOMReady();
				});
			}
		`)
	}
}