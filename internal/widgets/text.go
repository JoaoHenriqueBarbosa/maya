package widgets

import (
	"context"
	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/reactive"
)

// Text widget displays text content
type Text struct {
	*BaseWidget
	
	// Text-specific properties
	text  *reactive.Signal[string]
	style *reactive.Signal[TextStyle]
	
	// Layout measurements
	measuredWidth  float64
	measuredHeight float64
}

// NewText creates a new text widget
func NewText(id string, text string) *Text {
	t := &Text{
		BaseWidget: NewBaseWidget(id, "Text"),
		text:       reactive.NewSignal(text),
		style: reactive.NewSignal(TextStyle{
			FontFamily: "system-ui",
			FontSize:   14,
			FontWeight: FontWeightNormal,
			Color:      ColorBlack,
			LineHeight: 1.2,
		}),
	}
	
	// Watch for text changes
	dispose := reactive.Watch(func() {
		_ = t.text.Get()
		t.MarkNeedsLayout()
		t.MarkNeedsRepaint()
	})
	t.effects = append(t.effects, dispose)
	
	return t
}

// SetText updates the text content
func (t *Text) SetText(text string) {
	t.text.Set(text)
}

// GetText returns the text content
func (t *Text) GetText() string {
	return t.text.Get()
}

// SetStyle updates the text style
func (t *Text) SetStyle(style TextStyle) {
	t.style.Set(style)
	t.MarkNeedsLayout()
	t.MarkNeedsRepaint()
}

// GetStyle returns the text style
func (t *Text) GetStyle() TextStyle {
	return t.style.Get()
}

// Build creates the render object
func (t *Text) Build(ctx context.Context) RenderObject {
	return &RenderParagraph{
		Text:  t.text.Get(),
		Style: t.style.Get(),
	}
}

// Layout calculates text dimensions
func (t *Text) Layout(constraints core.Constraints) (width, height float64) {
	text := t.text.Get()
	style := t.style.Get()
	
	// Simple text measurement (in production, use proper text shaping)
	charWidth := style.FontSize * 0.6  // Approximate character width
	lineHeight := style.FontSize * style.LineHeight
	
	// Calculate dimensions
	t.measuredWidth = float64(len(text)) * charWidth
	t.measuredHeight = lineHeight
	
	// Apply constraints
	width = t.measuredWidth
	if width < constraints.MinWidth {
		width = constraints.MinWidth
	}
	if width > constraints.MaxWidth {
		width = constraints.MaxWidth
	}
	
	height = t.measuredHeight
	if height < constraints.MinHeight {
		height = constraints.MinHeight
	}
	if height > constraints.MaxHeight {
		height = constraints.MaxHeight
	}
	
	t.cachedSize = Size{Width: width, Height: height}
	t.needsLayout.Set(false)
	
	return width, height
}

// GetIntrinsicWidth returns natural width for text
func (t *Text) GetIntrinsicWidth(height float64) float64 {
	return t.measuredWidth
}

// GetIntrinsicHeight returns natural height for text
func (t *Text) GetIntrinsicHeight(width float64) float64 {
	// For multi-line text, calculate based on width
	// For now, return single line height
	return t.measuredHeight
}

// Paint renders the text
func (t *Text) Paint(context core.PaintContext) {
	if !t.IsVisible() {
		return
	}
	
	text := t.text.Get()
	style := t.style.Get()
	
	// Apply opacity if set
	opacity := t.GetOpacity()
	if opacity < 1.0 {
		textColor := style.Color
		textColor.A = uint8(float64(textColor.A) * opacity)
		style.Color = textColor
	}
	
	// Draw text
	context.DrawText(text, core.Offset{X: 0, Y: 0}, core.Paint{
		Color: style.Color,
		Alpha: 1.0,
	})
	
	t.needsRepaint.Set(false)
}

// HandleEvent processes input events - text widgets don't handle events by default
func (t *Text) HandleEvent(event core.Event) bool {
	return false
}

// Dispose cleans up the text widget
func (t *Text) Dispose() {
	t.text.Dispose()
	t.style.Dispose()
	t.BaseWidget.Dispose()
}