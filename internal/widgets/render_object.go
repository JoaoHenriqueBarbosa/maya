package widgets

import "github.com/maya-framework/maya/internal/core"

// RenderObject is the base interface for all render objects
type RenderObject interface {
	Layout(constraints core.Constraints) Size
	Paint(context core.PaintContext)
}

// RenderBox is a basic render object
type RenderBox struct {
	Size Size
}

func (r *RenderBox) Layout(constraints core.Constraints) Size {
	return r.Size
}

func (r *RenderBox) Paint(context core.PaintContext) {
	// Basic implementation
}

// RenderParagraph renders text
type RenderParagraph struct {
	Text  string
	Style TextStyle
}

func (r *RenderParagraph) Layout(constraints core.Constraints) Size {
	// Calculate text dimensions
	charWidth := r.Style.FontSize * 0.6
	lineHeight := r.Style.FontSize * r.Style.LineHeight
	
	width := float64(len(r.Text)) * charWidth
	height := lineHeight
	
	// Apply constraints
	if width < constraints.MinWidth {
		width = constraints.MinWidth
	}
	if width > constraints.MaxWidth {
		width = constraints.MaxWidth
	}
	
	if height < constraints.MinHeight {
		height = constraints.MinHeight
	}
	if height > constraints.MaxHeight {
		height = constraints.MaxHeight
	}
	
	return Size{Width: width, Height: height}
}

func (r *RenderParagraph) Paint(context core.PaintContext) {
	context.DrawText(r.Text, core.Offset{X: 0, Y: 0}, core.Paint{
		Color: r.Style.Color,
		Alpha: 1.0,
	})
}

// RenderButton represents a button render object
type RenderButton struct {
	Label     string
	OnPressed func()
	Disabled  bool
}

func (r *RenderButton) Layout(constraints core.Constraints) Size {
	// Simple button layout
	textWidth := float64(len(r.Label)) * 8.4
	textHeight := 14.0 * 1.2
	
	width := textWidth + 32  // Add padding
	height := textHeight + 16
	
	// Apply constraints
	if width < constraints.MinWidth {
		width = constraints.MinWidth
	}
	if width > constraints.MaxWidth {
		width = constraints.MaxWidth
	}
	
	if height < constraints.MinHeight {
		height = constraints.MinHeight
	}
	if height > constraints.MaxHeight {
		height = constraints.MaxHeight
	}
	
	return Size{Width: width, Height: height}
}

func (r *RenderButton) Paint(context core.PaintContext) {
	// Draw button background
	bgColor := ColorBlue
	if r.Disabled {
		bgColor = core.Color{R: 128, G: 128, B: 128, A: 255}
	}
	
	bounds := core.Bounds{
		X:      0,
		Y:      0,
		Width:  100, // Default width
		Height: 40,  // Default height
	}
	context.DrawRect(bounds, core.Paint{
		Color: bgColor,
		Alpha: 1.0,
	})
	
	// Draw label
	textOffset := core.Offset{
		X: 16,
		Y: 8,
	}
	
	textColor := ColorWhite
	if r.Disabled {
		textColor = core.Color{R: 200, G: 200, B: 200, A: 255}
	}
	
	context.DrawText(r.Label, textOffset, core.Paint{
		Color: textColor,
		Alpha: 1.0,
	})
}

// RenderDecoratedBox renders a decorated box
type RenderDecoratedBox struct {
	Decoration BoxDecoration
	Child      RenderObject
}

// BoxDecoration describes box decoration
type BoxDecoration struct {
	Color        core.Color
	BorderRadius float64
	BorderColor  core.Color
	BorderWidth  float64
	BoxShadow    *BoxShadow  // Optional shadow
}

func (r *RenderDecoratedBox) Layout(constraints core.Constraints) Size {
	if r.Child != nil {
		return r.Child.Layout(constraints)
	}
	return Size{
		Width:  constraints.MaxWidth,
		Height: constraints.MaxHeight,
	}
}

func (r *RenderDecoratedBox) Paint(context core.PaintContext) {
	// Draw decoration
	if r.Decoration.Color.A > 0 {
		bounds := core.Bounds{
			X:      0,
			Y:      0,
			Width:  100, // Default width
			Height: 100, // Default height
		}
		context.DrawRect(bounds, core.Paint{
			Color: r.Decoration.Color,
			Alpha: 1.0,
		})
	}
	
	// Draw border
	if r.Decoration.BorderWidth > 0 {
		bounds := core.Bounds{
			X:      0,
			Y:      0,
			Width:  100, // Default width
			Height: 100, // Default height
		}
		context.DrawRect(bounds, core.Paint{
			Color: r.Decoration.BorderColor,
			Alpha: 1.0,
		})
	}
	
	// Paint child
	if r.Child != nil {
		r.Child.Paint(context)
	}
}