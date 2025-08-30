package widgets

import (
	"context"
	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/reactive"
)

// Button widget represents a clickable button
type Button struct {
	*BaseWidget
	
	// Button-specific properties
	label      *reactive.Signal[string]
	onPressed  func()
	isPressed  *reactive.Signal[bool]
	isHovered  *reactive.Signal[bool]
	isDisabled *reactive.Signal[bool]
	
	// Style properties
	backgroundColor *reactive.Signal[core.Color]
	pressedColor    *reactive.Signal[core.Color]
	hoverColor      *reactive.Signal[core.Color]
	textColor       *reactive.Signal[core.Color]
	borderRadius    float64
	padding         EdgeInsets
}

// EdgeInsets represents padding/margin values
type EdgeInsets struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

// NewButton creates a new button widget
func NewButton(id string, label string, onPressed func()) *Button {
	b := &Button{
		BaseWidget:      NewBaseWidget(id, "Button"),
		label:           reactive.NewSignal(label),
		onPressed:       onPressed,
		isPressed:       reactive.NewSignal(false),
		isHovered:       reactive.NewSignal(false),
		isDisabled:      reactive.NewSignal(false),
		backgroundColor: reactive.NewSignal(ColorBlue),
		pressedColor:    reactive.NewSignal(core.Color{R: 0, G: 0, B: 200, A: 255}),
		hoverColor:      reactive.NewSignal(core.Color{R: 100, G: 100, B: 255, A: 255}),
		textColor:       reactive.NewSignal(ColorWhite),
		borderRadius:    4,
		padding: EdgeInsets{
			Top:    8,
			Right:  16,
			Bottom: 8,
			Left:   16,
		},
	}
	
	// Setup reactive effects for visual feedback
	b.setupButtonEffects()
	
	return b
}

// setupButtonEffects creates reactive effects for button states
func (b *Button) setupButtonEffects() {
	// Watch for state changes that require repaint
	dispose := reactive.Watch(func() {
		_ = b.isPressed.Get()
		_ = b.isHovered.Get()
		_ = b.isDisabled.Get()
		b.MarkNeedsRepaint()
	})
	b.effects = append(b.effects, dispose)
}

// SetLabel updates the button label
func (b *Button) SetLabel(label string) {
	b.label.Set(label)
	b.MarkNeedsLayout()
	b.MarkNeedsRepaint()
}

// GetLabel returns the button label
func (b *Button) GetLabel() string {
	return b.label.Get()
}

// SetDisabled sets the disabled state
func (b *Button) SetDisabled(disabled bool) {
	b.isDisabled.Set(disabled)
}

// IsDisabled returns the disabled state
func (b *Button) IsDisabled() bool {
	return b.isDisabled.Get()
}

// Build creates the render object
func (b *Button) Build(ctx context.Context) RenderObject {
	// Create a text widget for the label
	textWidget := NewText(b.id+"-label", b.label.Get())
	textWidget.SetStyle(TextStyle{
		FontFamily: "system-ui",
		FontSize:   14,
		FontWeight: FontWeightNormal,
		Color:      b.textColor.Get(),
		LineHeight: 1.2,
	})
	
	return &RenderButton{
		Label:     b.label.Get(),
		OnPressed: b.onPressed,
		Disabled:  b.isDisabled.Get(),
	}
}

// Layout calculates button dimensions
func (b *Button) Layout(constraints core.Constraints) (width, height float64) {
	label := b.label.Get()
	
	// Calculate text size
	textWidth := float64(len(label)) * 8.4  // Approximate
	textHeight := 14.0 * 1.2                // FontSize * LineHeight
	
	// Add padding
	width = textWidth + b.padding.Left + b.padding.Right
	height = textHeight + b.padding.Top + b.padding.Bottom
	
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
	
	b.cachedSize = Size{Width: width, Height: height}
	b.needsLayout.Set(false)
	
	return width, height
}

// Paint renders the button
func (b *Button) Paint(context core.PaintContext) {
	if !b.IsVisible() {
		return
	}
	
	// Determine background color based on state
	var bgColor core.Color
	if b.isDisabled.Get() {
		bgColor = core.Color{R: 128, G: 128, B: 128, A: 255} // Gray for disabled
	} else if b.isPressed.Get() {
		bgColor = b.pressedColor.Get()
	} else if b.isHovered.Get() {
		bgColor = b.hoverColor.Get()
	} else {
		bgColor = b.backgroundColor.Get()
	}
	
	// Draw button background
	bounds := core.Bounds{
		X:      0,
		Y:      0,
		Width:  b.cachedSize.Width,
		Height: b.cachedSize.Height,
	}
	
	context.DrawRect(bounds, core.Paint{
		Color: bgColor,
		Alpha: 1.0,
	})
	
	// Draw border if needed
	if b.borderRadius > 0 {
		// In production, implement rounded rectangle
		context.DrawRect(bounds, core.Paint{
			Color: core.Color{R: 0, G: 0, B: 0, A: 100},
			Alpha: 1.0,
		})
	}
	
	// Draw text centered
	label := b.label.Get()
	textOffset := core.Offset{
		X: b.padding.Left,
		Y: b.padding.Top,
	}
	
	textColor := b.textColor.Get()
	if b.isDisabled.Get() {
		textColor = core.Color{R: 200, G: 200, B: 200, A: 255}
	}
	
	context.DrawText(label, textOffset, core.Paint{
		Color: textColor,
		Alpha: 1.0,
	})
	
	b.needsRepaint.Set(false)
}

// HandleEvent processes button events
func (b *Button) HandleEvent(event core.Event) bool {
	if b.isDisabled.Get() {
		return false
	}
	
	// For now, buttons don't handle events directly from core
	// This would be handled by a higher-level event system
	// that converts core.Event to widget-specific events
	return false
}

// containsPoint checks if a point is within the button bounds
func (b *Button) containsPoint(point core.Offset) bool {
	size := b.cachedSize
	return point.X >= 0 && point.X <= size.Width &&
	       point.Y >= 0 && point.Y <= size.Height
}

// Dispose cleans up the button widget
func (b *Button) Dispose() {
	b.label.Dispose()
	b.isPressed.Dispose()
	b.isHovered.Dispose()
	b.isDisabled.Dispose()
	b.backgroundColor.Dispose()
	b.pressedColor.Dispose()
	b.hoverColor.Dispose()
	b.textColor.Dispose()
	b.BaseWidget.Dispose()
}