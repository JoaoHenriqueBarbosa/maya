//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"github.com/maya-framework/maya"
	"github.com/maya-framework/maya/internal/widgets"
)

func main() {
	// Create reactive state
	counter := maya.Signal(0)
	message := maya.Signal("Click the buttons!")
	
	// Create the app
	app := maya.New(func() widgets.WidgetImpl {
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
					fmt.Println("Increment button clicked!")
					counter.Set(counter.Get() + 1)
					message.Set(fmt.Sprintf("Incremented to %d", counter.Get()))
					fmt.Printf("Counter is now: %d\n", counter.Get())
				}),
				
				maya.Button("Decrement", func() {
					fmt.Println("Decrement button clicked!")
					counter.Set(counter.Get() - 1)
					message.Set(fmt.Sprintf("Decremented to %d", counter.Get()))
					fmt.Printf("Counter is now: %d\n", counter.Get())
				}),
				
				maya.Button("Reset", func() {
					fmt.Println("Reset button clicked!")
					counter.Set(0)
					message.Set("Counter reset!")
					fmt.Printf("Counter is now: %d\n", counter.Get())
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