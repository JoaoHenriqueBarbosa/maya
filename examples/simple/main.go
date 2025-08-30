//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"github.com/maya-framework/maya"
)

func main() {
	// Create reactive state
	counter := maya.Signal(0)
	message := maya.Signal("Click the buttons!")
	
	// Create the app
	app := maya.New(func() maya.Widget {
		return maya.Container(
			// Title
			maya.Title("Maya Counter App"),
			
			// Counter display
			maya.Row(
				maya.Text("Count: "),
				maya.TextSignal(counter, func(v int) string {
					return fmt.Sprintf("%d", v)
				}),
			),
			
			// Buttons
			maya.Row(
				maya.Button("Increment", func() {
					counter.Set(counter.Get() + 1)
					message.Set(fmt.Sprintf("Incremented to %d", counter.Get()))
				}),
				
				maya.Button("Decrement", func() {
					counter.Set(counter.Get() - 1)
					message.Set(fmt.Sprintf("Decremented to %d", counter.Get()))
				}),
				
				maya.Button("Reset", func() {
					counter.Set(0)
					message.Set("Counter reset!")
				}),
			),
			
			// Message display
			maya.TextSignal(message, func(v string) string {
				return v
			}),
		)
	})
	
	// Run the app
	app.Run()
}