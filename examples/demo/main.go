//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"github.com/maya-framework/maya"
	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/logger"
	"github.com/maya-framework/maya/internal/widgets"
)

// Custom colors
var (
	ColorPrimary = &core.Color{R: 79, G: 70, B: 229, A: 255}    // Indigo-600
	ColorSuccess = &core.Color{R: 34, G: 197, B: 94, A: 255}    // Green-500
	ColorDanger  = &core.Color{R: 239, G: 68, B: 68, A: 255}    // Red-500
	ColorWarning = &core.Color{R: 245, G: 158, B: 11, A: 255}   // Amber-500
	ColorDark    = &core.Color{R: 31, G: 41, B: 55, A: 255}     // Gray-800
	ColorMuted   = &core.Color{R: 107, G: 114, B: 128, A: 255}  // Gray-500
	ColorLight   = &core.Color{R: 243, G: 244, B: 246, A: 255}  // Gray-100
)

func main() {
	logger.Info("APP", "Maya Demo Starting")
	
	// Reactive state
	counter := maya.Signal(0)
	
	// Create the app
	app := maya.New(func() widgets.WidgetImpl {
		return maya.StyledContainer(
			maya.ContainerStyle{
				Background: ColorLight,
				Padding:    &widgets.EdgeInsets{Top: 40, Right: 40, Bottom: 40, Left: 40},
			},
			maya.Column(
				// Title card
				maya.StyledContainer(
					maya.ContainerStyle{
						Background:   maya.ColorWhite,
						BorderRadius: 12,
						Padding:      &widgets.EdgeInsets{Top: 30, Right: 40, Bottom: 30, Left: 40},
						Shadow: &widgets.BoxShadow{
							Color:      *ColorMuted,
							Offset:     maya.Offset{X: 0, Y: 4},
							BlurRadius: 6,
						},
					},
					maya.Column(
						maya.Title("🚀 Maya Framework Demo"),
						maya.Text("A modern reactive UI framework for Go"),
					),
				),
				
				// Counter card
				maya.StyledContainer(
					maya.ContainerStyle{
						Background:   maya.ColorWhite,
						BorderRadius: 12,
						Padding:      &widgets.EdgeInsets{Top: 30, Right: 40, Bottom: 30, Left: 40},
						Shadow: &widgets.BoxShadow{
							Color:      *ColorMuted,
							Offset:     maya.Offset{X: 0, Y: 2},
							BlurRadius: 4,
						},
					},
					maya.Column(
						maya.Title("Counter Example"),
						
						// Counter display
						maya.StyledContainer(
							maya.ContainerStyle{
								Background:   ColorPrimary,
								BorderRadius: 8,
								Padding:      &widgets.EdgeInsets{Top: 20, Right: 30, Bottom: 20, Left: 30},
							},
							maya.TextSignal(counter, func(v int) string {
								return fmt.Sprintf("🔢 Count: %d", v)
							}),
						),
						
						// Buttons
						maya.Row(
							maya.Button("➕ Increment", func() {
								counter.Set(counter.Get() + 1)
							}),
							maya.Button("➖ Decrement", func() {
								counter.Set(counter.Get() - 1)
							}),
							maya.Button("🔄 Reset", func() {
								counter.Set(0)
							}),
						),
					),
				),
				
				// Info cards row
				maya.Row(
					// Card 1
					createInfoCard("Users", "1,234", ColorSuccess, "👥"),
					// Card 2
					createInfoCard("Revenue", "$12,345", ColorWarning, "💰"),
					// Card 3
					createInfoCard("Growth", "+23%", ColorPrimary, "📈"),
				),
			),
		)
	})
	
	// Run the app
	app.Run()
}

func createInfoCard(title, value string, color *core.Color, icon string) widgets.WidgetImpl {
	return maya.StyledContainer(
		maya.ContainerStyle{
			Background:   maya.ColorWhite,
			BorderRadius: 12,
			Padding:      &widgets.EdgeInsets{Top: 25, Right: 30, Bottom: 25, Left: 30},
			Shadow: &widgets.BoxShadow{
				Color:      *ColorMuted,
				Offset:     maya.Offset{X: 0, Y: 2},
				BlurRadius: 4,
			},
		},
		maya.Column(
			maya.Row(
				maya.StyledContainer(
					maya.ContainerStyle{
						Background:   color,
						BorderRadius: 8,
						Padding:      &widgets.EdgeInsets{Top: 10, Right: 10, Bottom: 10, Left: 10},
					},
					maya.Text(icon),
				),
				maya.Text(title),
			),
			maya.Title(value),
		),
	)
}