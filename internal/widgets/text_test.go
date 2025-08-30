package widgets

import (
	"context"
	"testing"
	"time"

	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/reactive"
)

func TestNewText(t *testing.T) {
	text := NewText("text-id", "Hello World")
	
	if text.ID() != "text-id" {
		t.Errorf("Expected ID 'text-id', got %s", text.ID())
	}
	
	if text.Type() != "Text" {
		t.Errorf("Expected Type 'Text', got %s", text.Type())
	}
	
	if text.GetText() != "Hello World" {
		t.Errorf("Expected text 'Hello World', got %s", text.GetText())
	}
	
	// Check default style
	style := text.GetStyle()
	if style.FontFamily != "system-ui" {
		t.Errorf("Expected default font family 'system-ui', got %s", style.FontFamily)
	}
	
	if style.FontSize != 14 {
		t.Errorf("Expected default font size 14, got %f", style.FontSize)
	}
	
	if style.FontWeight != FontWeightNormal {
		t.Errorf("Expected normal font weight, got %d", style.FontWeight)
	}
	
	if style.Color != ColorBlack {
		t.Errorf("Expected black color, got %v", style.Color)
	}
	
	if style.LineHeight != 1.2 {
		t.Errorf("Expected line height 1.2, got %f", style.LineHeight)
	}
}

func TestText_SetText(t *testing.T) {
	text := NewText("test", "Initial")
	
	text.SetText("Updated")
	
	if text.GetText() != "Updated" {
		t.Errorf("Expected text 'Updated', got %s", text.GetText())
	}
	
	// Should trigger layout and repaint
	if !text.NeedsRepaint() {
		t.Error("Changing text should trigger repaint")
	}
}

func TestText_SetStyle(t *testing.T) {
	text := NewText("test", "Test")
	
	newStyle := TextStyle{
		FontFamily: "monospace",
		FontSize:   16,
		FontWeight: FontWeightBold,
		Color:      ColorRed,
		LineHeight: 1.5,
	}
	
	text.SetStyle(newStyle)
	
	style := text.GetStyle()
	if style.FontFamily != "monospace" {
		t.Errorf("Expected font family 'monospace', got %s", style.FontFamily)
	}
	
	if style.FontSize != 16 {
		t.Errorf("Expected font size 16, got %f", style.FontSize)
	}
	
	if style.FontWeight != FontWeightBold {
		t.Errorf("Expected bold font weight, got %d", style.FontWeight)
	}
	
	if style.Color != ColorRed {
		t.Errorf("Expected red color, got %v", style.Color)
	}
	
	// Should trigger layout and repaint
	if !text.NeedsRepaint() {
		t.Error("Changing style should trigger repaint")
	}
}

func TestText_Build(t *testing.T) {
	text := NewText("test", "Hello")
	
	ctx := context.Background()
	renderObj := text.Build(ctx)
	
	if renderObj == nil {
		t.Fatal("Build should return a RenderObject")
	}
	
	paragraph, ok := renderObj.(*RenderParagraph)
	if !ok {
		t.Fatal("Text.Build should return RenderParagraph")
	}
	
	if paragraph.Text != "Hello" {
		t.Errorf("Expected text 'Hello', got %s", paragraph.Text)
	}
	
	if paragraph.Style.FontFamily != "system-ui" {
		t.Errorf("Expected font family 'system-ui', got %s", paragraph.Style.FontFamily)
	}
}

func TestText_Layout(t *testing.T) {
	text := NewText("test", "Hello")
	
	constraints := core.Constraints{
		MinWidth:  0,
		MaxWidth:  1000,
		MinHeight: 0,
		MaxHeight: 1000,
	}
	
	width, height := text.Layout(constraints)
	
	// Text layout should calculate based on character count and font size
	// "Hello" = 5 chars, default font size = 14
	expectedWidth := 5 * 14 * 0.6  // ~42
	expectedHeight := 14 * 1.2       // ~16.8
	
	if width != expectedWidth {
		t.Errorf("Expected width ~%f, got %f", expectedWidth, width)
	}
	
	if height != expectedHeight {
		t.Errorf("Expected height ~%f, got %f", expectedHeight, height)
	}
	
	// Test constraints
	narrowConstraints := core.Constraints{
		MinWidth:  10,
		MaxWidth:  20,
		MinHeight: 5,
		MaxHeight: 10,
	}
	
	width, height = text.Layout(narrowConstraints)
	
	if width != 20 {
		t.Errorf("Width should be constrained to max, got %f", width)
	}
	
	if height != 10 {
		t.Errorf("Height should be constrained to max, got %f", height)
	}
}

func TestText_GetIntrinsicWidth(t *testing.T) {
	text := NewText("test", "Hello World")
	
	// Layout first to calculate measured dimensions
	text.Layout(core.Constraints{
		MinWidth: 0, MaxWidth: 1000,
		MinHeight: 0, MaxHeight: 1000,
	})
	
	intrinsicWidth := text.GetIntrinsicWidth(100)
	
	if intrinsicWidth != text.measuredWidth {
		t.Errorf("Intrinsic width should match measured width, got %f", intrinsicWidth)
	}
}

func TestText_GetIntrinsicHeight(t *testing.T) {
	text := NewText("test", "Hello")
	
	// Layout first to calculate measured dimensions
	text.Layout(core.Constraints{
		MinWidth: 0, MaxWidth: 1000,
		MinHeight: 0, MaxHeight: 1000,
	})
	
	intrinsicHeight := text.GetIntrinsicHeight(100)
	
	if intrinsicHeight != text.measuredHeight {
		t.Errorf("Intrinsic height should match measured height, got %f", intrinsicHeight)
	}
}

func TestText_Paint(t *testing.T) {
	text := NewText("test", "Hello")
	ctx := &mockPaintContext{}
	
	text.Paint(ctx)
	
	if !ctx.drawTextCalled {
		t.Error("Text should call DrawText on context")
	}
}

func TestText_Paint_Invisible(t *testing.T) {
	text := NewText("test", "Hello")
	text.UpdateProp("visible", false)
	
	ctx := &mockPaintContext{}
	text.Paint(ctx)
	
	if ctx.drawTextCalled {
		t.Error("Invisible text should not be painted")
	}
}

func TestText_Paint_WithOpacity(t *testing.T) {
	text := NewText("test", "Hello")
	text.UpdateProp("opacity", 0.5)
	
	style := text.GetStyle()
	style.Color = core.Color{R: 255, G: 0, B: 0, A: 255}
	text.SetStyle(style)
	
	ctx := &mockPaintContext{}
	text.Paint(ctx)
	
	if !ctx.drawTextCalled {
		t.Error("Text should be painted with opacity")
	}
	
	// The alpha should be adjusted based on opacity
	// Original alpha 255 * 0.5 = 127.5 -> 127
	// Note: This test assumes the paint context captures the modified color
}

func TestText_HandleEvent(t *testing.T) {
	text := NewText("test", "Hello")
	
	event := &core.MockEvent{EventType: "test"}
	
	// Text widgets don't handle events by default
	handled := text.HandleEvent(event)
	if handled {
		t.Error("Text should not handle events by default")
	}
}

func TestText_Dispose(t *testing.T) {
	text := NewText("test", "Hello")
	
	// Add a watcher to verify signals are disposed
	dispose := reactive.Watch(func() {
		_ = text.text.Get()
	})
	defer dispose()
	
	text.Dispose()
	
	// After dispose, the widget should be cleaned up
	// We can't directly test if signals are disposed, but we can verify
	// the Dispose method doesn't panic
	if !text.disposed.Load() {
		t.Error("Text should be marked as disposed")
	}
}

func TestText_ReactiveUpdates(t *testing.T) {
	text := NewText("test", "Initial")
	
	var layoutNeeded bool
	var repaintNeeded bool
	
	// Watch for changes
	dispose := reactive.Watch(func() {
		if text.needsLayout.Get() {
			layoutNeeded = true
		}
		if text.needsRepaint.Get() {
			repaintNeeded = true
		}
	})
	defer dispose()
	
	// Change text
	text.SetText("Updated")
	
	// Give reactive system time to process
	time.Sleep(10 * time.Millisecond)
	
	if !layoutNeeded {
		t.Error("Changing text should trigger layout")
	}
	
	if !repaintNeeded {
		t.Error("Changing text should trigger repaint")
	}
}

func TestText_EmptyText(t *testing.T) {
	text := NewText("test", "")
	
	constraints := core.Constraints{
		MinWidth:  0,
		MaxWidth:  1000,
		MinHeight: 0,
		MaxHeight: 1000,
	}
	
	width, height := text.Layout(constraints)
	
	// Empty text should still have height based on line height
	if width != 0 {
		t.Errorf("Empty text should have 0 width, got %f", width)
	}
	
	expectedHeight := 14 * 1.2 // font size * line height
	if height != expectedHeight {
		t.Errorf("Empty text should have line height, got %f", height)
	}
}

func TestText_LongText(t *testing.T) {
	longText := "This is a very long text that should be constrained by max width"
	text := NewText("test", longText)
	
	constraints := core.Constraints{
		MinWidth:  0,
		MaxWidth:  100, // Narrow max width
		MinHeight: 0,
		MaxHeight: 1000,
	}
	
	width, height := text.Layout(constraints)
	
	if width != 100 {
		t.Errorf("Long text should be constrained to max width, got %f", width)
	}
	
	// Height should be based on line height
	expectedHeight := 14 * 1.2
	if height != expectedHeight {
		t.Errorf("Expected height %f, got %f", expectedHeight, height)
	}
}

func TestText_MinConstraints(t *testing.T) {
	text := NewText("test", "Hi")
	
	constraints := core.Constraints{
		MinWidth:  100,
		MaxWidth:  200,
		MinHeight: 50,
		MaxHeight: 100,
	}
	
	width, height := text.Layout(constraints)
	
	// Text dimensions should respect minimum constraints
	if width < 100 {
		t.Errorf("Width should respect minimum constraint, got %f", width)
	}
	
	if height < 50 {
		t.Errorf("Height should respect minimum constraint, got %f", height)
	}
}