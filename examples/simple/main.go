//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"github.com/maya-framework/maya"
	"github.com/maya-framework/maya/internal/logger"
	"github.com/maya-framework/maya/internal/widgets"
)

func main() {
	logger.Info("APP", "Maya App Starting")
	
	// Create reactive state
	counter := maya.Signal(0)
	message := maya.Signal("Click the buttons!")
	logger.Debug("APP", "Signals created")
	
	// Derived signals - these update automatically when counter changes
	doubled := maya.Signal(0)
	squared := maya.Signal(0)
	analysis := maya.Signal("Zero")
	
	// Create effects to update derived signals when counter changes
	maya.CreateEffect(func() {
		value := counter.Get()
		doubled.Set(value * 2)
		squared.Set(value * value)
		
		// Update analysis
		if value < 0 {
			analysis.Set("Negative")
		} else if value == 0 {
			analysis.Set("Zero")
		} else if value < 10 {
			analysis.Set("Small")
		} else if value < 100 {
			analysis.Set("Medium")
		} else {
			analysis.Set("Large")
		}
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
							logger.Debug("UI", "Increment button clicked")
							counter.Set(counter.Get() + 1)
							message.Set(fmt.Sprintf("Incremented to %d", counter.Get()))
							logger.Debug("UI", "Counter is now: %d", counter.Get())
						}),
						
						maya.Button("Decrement", func() {
							logger.Debug("UI", "Decrement button clicked")
							counter.Set(counter.Get() - 1)
							message.Set(fmt.Sprintf("Decremented to %d", counter.Get()))
							logger.Debug("UI", "Counter is now: %d", counter.Get())
						}),
						
						maya.Button("Reset", func() {
							logger.Debug("UI", "Reset button clicked")
							counter.Set(0)
							message.Set("Counter reset!")
							logger.Debug("UI", "Counter is now: %d", counter.Get())
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
						// Derived values section
						maya.Title("Derived Values"),
					maya.Row(
						maya.Text("Doubled: "),
						maya.TextSignal(doubled, func(v int) string {
							return fmt.Sprintf("%d", v)
						}),
					),
					maya.Row(
						maya.Text("Squared: "),
						maya.TextSignal(squared, func(v int) string {
							return fmt.Sprintf("%d", v)
						}),
					),
					
					// Analysis section
					maya.Title("Analysis"),
					maya.Row(
						maya.Text("Status: "),
						maya.TextSignal(analysis, func(v string) string {
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