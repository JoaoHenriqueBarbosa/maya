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
	
	// Memo example - computed values that cache
	doubled := maya.Memo(func() int {
		fmt.Println("[MEMO] Computing doubled value...")
		return counter.Get() * 2
	})
	
	squared := maya.Memo(func() int {
		fmt.Println("[MEMO] Computing squared value...")
		return counter.Get() * counter.Get()
	})
	
	// Computed example - derived state
	analysis := maya.Computed(func() string {
		fmt.Println("[COMPUTED] Analyzing counter...")
		value := counter.Get()
		if value < 0 {
			return "Negative"
		} else if value == 0 {
			return "Zero"
		} else if value < 10 {
			return "Small"
		} else if value < 100 {
			return "Medium"
		}
		return "Large"
	})
	
	// Create the app
	app := maya.New(func() widgets.WidgetImpl {
		return maya.Container(
			// Title
			maya.Title("Maya Reactive Framework Demo"),
			
			// Main content in two columns
			maya.Row(
				// Left column - Counter section
				maya.StyledContainer(
					maya.ContainerStyle{
						Background:   maya.ColorLightGray,
						BorderColor:  maya.ColorGray,
						BorderWidth:  2,
						BorderRadius: 8,
						Padding:      &widgets.EdgeInsets{Top: 10, Right: 10, Bottom: 10, Left: 10},
					},
					maya.Column(
						maya.Title("Counter (Signal)"),
					
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
					),
				),
				
				// Right column - Memo and Computed sections
				maya.StyledContainer(
					maya.ContainerStyle{
						Background:   maya.ColorWhite,
						BorderColor:  maya.ColorBlue,
						BorderWidth:  2,
						BorderRadius: 8,
						Padding:      &widgets.EdgeInsets{Top: 10, Right: 10, Bottom: 10, Left: 10},
						Shadow: &widgets.BoxShadow{
							Color:     *maya.ColorGray,
							Offset:    maya.Offset{X: 2, Y: 2},
							BlurRadius: 4,
						},
					},
					maya.Column(
						// Memo section
						maya.Title("Memoized Values"),
					maya.Row(
						maya.Text("Doubled: "),
						maya.TextMemo(doubled, func(v int) string {
							return fmt.Sprintf("%d", v)
						}),
					),
					maya.Row(
						maya.Text("Squared: "),
						maya.TextMemo(squared, func(v int) string {
							return fmt.Sprintf("%d", v)
						}),
					),
					
					// Computed section
					maya.Title("Computed Analysis"),
					maya.Row(
						maya.Text("Status: "),
						maya.TextMemo(analysis.Memo, func(v string) string {
							return v
						}),
						),
					),
				),
			),
		)
	})
	
	// Run the app
	app.Run()
}