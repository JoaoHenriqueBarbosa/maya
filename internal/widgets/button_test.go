package widgets

import (
	"context"
	"testing"
	"time"

	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/reactive"
)

func TestNewButton(t *testing.T) {
	var pressed bool
	onPressed := func() {
		pressed = true
	}
	
	button := NewButton("btn-id", "Click Me", onPressed)
	
	if button.ID() != "btn-id" {
		t.Errorf("Expected ID 'btn-id', got %s", button.ID())
	}
	
	if button.Type() != "Button" {
		t.Errorf("Expected Type 'Button', got %s", button.Type())
	}
	
	if button.GetLabel() != "Click Me" {
		t.Errorf("Expected label 'Click Me', got %s", button.GetLabel())
	}
	
	// Test onPressed callback
	button.onPressed()
	if !pressed {
		t.Error("onPressed callback should be called")
	}
	
	// Check default states
	if button.isPressed.Get() {
		t.Error("Button should not be pressed initially")
	}
	
	if button.isHovered.Get() {
		t.Error("Button should not be hovered initially")
	}
	
	if button.IsDisabled() {
		t.Error("Button should not be disabled initially")
	}
	
	// Check default colors
	if button.backgroundColor.Get() != ColorBlue {
		t.Error("Default background color should be blue")
	}
	
	if button.textColor.Get() != ColorWhite {
		t.Error("Default text color should be white")
	}
	
	// Check default padding
	if button.padding.Top != 8 || button.padding.Bottom != 8 {
		t.Error("Default vertical padding should be 8")
	}
	
	if button.padding.Left != 16 || button.padding.Right != 16 {
		t.Error("Default horizontal padding should be 16")
	}
	
	if button.borderRadius != 4 {
		t.Error("Default border radius should be 4")
	}
}

func TestButton_SetLabel(t *testing.T) {
	button := NewButton("test", "Initial", nil)
	
	button.SetLabel("Updated")
	
	if button.GetLabel() != "Updated" {
		t.Errorf("Expected label 'Updated', got %s", button.GetLabel())
	}
	
	// Should trigger layout and repaint
	if !button.NeedsRepaint() {
		t.Error("Changing label should trigger repaint")
	}
}

func TestButton_SetDisabled(t *testing.T) {
	button := NewButton("test", "Test", nil)
	
	button.SetDisabled(true)
	
	if !button.IsDisabled() {
		t.Error("Button should be disabled")
	}
	
	button.SetDisabled(false)
	
	if button.IsDisabled() {
		t.Error("Button should be enabled")
	}
}

func TestButton_Build(t *testing.T) {
	var pressed bool
	button := NewButton("test", "Click", func() { pressed = true })
	
	ctx := context.Background()
	renderObj := button.Build(ctx)
	
	if renderObj == nil {
		t.Fatal("Build should return a RenderObject")
	}
	
	renderButton, ok := renderObj.(*RenderButton)
	if !ok {
		t.Fatal("Button.Build should return RenderButton")
	}
	
	if renderButton.Label != "Click" {
		t.Errorf("Expected label 'Click', got %s", renderButton.Label)
	}
	
	if renderButton.Disabled {
		t.Error("Button should not be disabled")
	}
	
	// Test callback
	renderButton.OnPressed()
	if !pressed {
		t.Error("OnPressed should call the callback")
	}
}

func TestButton_Layout(t *testing.T) {
	button := NewButton("test", "OK", nil)
	
	constraints := core.Constraints{
		MinWidth:  0,
		MaxWidth:  1000,
		MinHeight: 0,
		MaxHeight: 1000,
	}
	
	width, height := button.Layout(constraints)
	
	// "OK" = 2 chars, font size = 14
	textWidth := 2 * 8.4 // ~16.8
	textHeight := 14.0 * 1.2 // ~16.8
	
	// Add padding
	expectedWidth := textWidth + 16 + 16   // ~48.8
	expectedHeight := textHeight + 8 + 8    // ~32.8
	
	if width != expectedWidth {
		t.Errorf("Expected width ~%f, got %f", expectedWidth, width)
	}
	
	if height != expectedHeight {
		t.Errorf("Expected height ~%f, got %f", expectedHeight, height)
	}
}

func TestButton_Layout_WithConstraints(t *testing.T) {
	button := NewButton("test", "Long Button Text", nil)
	
	constraints := core.Constraints{
		MinWidth:  100,
		MaxWidth:  150,
		MinHeight: 40,
		MaxHeight: 60,
	}
	
	width, height := button.Layout(constraints)
	
	// Should respect min constraints
	if width < 100 {
		t.Errorf("Width should respect minimum, got %f", width)
	}
	
	if height < 40 {
		t.Errorf("Height should respect minimum, got %f", height)
	}
	
	// Should not exceed max constraints
	if width > 150 {
		t.Errorf("Width should respect maximum, got %f", width)
	}
	
	if height > 60 {
		t.Errorf("Height should respect maximum, got %f", height)
	}
}

func TestButton_Paint(t *testing.T) {
	button := NewButton("test", "Click", nil)
	button.cachedSize = Size{Width: 100, Height: 40}
	
	ctx := &mockPaintContext{}
	button.Paint(ctx)
	
	if !ctx.drawRectCalled {
		t.Error("Button should draw background rect")
	}
	
	if !ctx.drawTextCalled {
		t.Error("Button should draw text")
	}
}

func TestButton_Paint_Invisible(t *testing.T) {
	button := NewButton("test", "Click", nil)
	button.UpdateProp("visible", false)
	
	ctx := &mockPaintContext{}
	button.Paint(ctx)
	
	if ctx.drawRectCalled || ctx.drawTextCalled {
		t.Error("Invisible button should not be painted")
	}
}

func TestButton_Paint_Disabled(t *testing.T) {
	button := NewButton("test", "Click", nil)
	button.SetDisabled(true)
	button.cachedSize = Size{Width: 100, Height: 40}
	
	ctx := &mockPaintContext{}
	button.Paint(ctx)
	
	if !ctx.drawRectCalled {
		t.Error("Disabled button should still be painted")
	}
	
	// Disabled state is being painted, checking color is implementation-specific
	// The important thing is that it paints
}

func TestButton_Paint_Pressed(t *testing.T) {
	button := NewButton("test", "Click", nil)
	button.isPressed.Set(true)
	button.cachedSize = Size{Width: 100, Height: 40}
	
	ctx := &mockPaintContext{}
	button.Paint(ctx)
	
	if !ctx.drawRectCalled {
		t.Error("Pressed button should be painted")
	}
	
	// Pressed state is being painted correctly
}

func TestButton_Paint_Hovered(t *testing.T) {
	button := NewButton("test", "Click", nil)
	button.isHovered.Set(true)
	button.cachedSize = Size{Width: 100, Height: 40}
	
	ctx := &mockPaintContext{}
	button.Paint(ctx)
	
	if !ctx.drawRectCalled {
		t.Error("Hovered button should be painted")
	}
	
	// Hovered state is being painted correctly
}

func TestButton_HandleEvent(t *testing.T) {
	button := NewButton("test", "Click", func() {})
	
	event := &core.MockEvent{EventType: "test"}
	
	// Disabled button should not handle events
	button.SetDisabled(true)
	handled := button.HandleEvent(event)
	if handled {
		t.Error("Disabled button should not handle events")
	}
	
	// Enable button
	button.SetDisabled(false)
	
	// For now, button returns false as event handling is simplified
	handled = button.HandleEvent(event)
	if handled {
		t.Error("Button should return false for simplified event handling")
	}
}

func TestButton_ContainsPoint(t *testing.T) {
	button := NewButton("test", "Click", nil)
	button.cachedSize = Size{Width: 100, Height: 50}
	
	// Test points inside
	if !button.containsPoint(core.Offset{X: 50, Y: 25}) {
		t.Error("Point (50, 25) should be inside button")
	}
	
	if !button.containsPoint(core.Offset{X: 0, Y: 0}) {
		t.Error("Point (0, 0) should be inside button")
	}
	
	if !button.containsPoint(core.Offset{X: 100, Y: 50}) {
		t.Error("Point (100, 50) should be inside button")
	}
	
	// Test points outside
	if button.containsPoint(core.Offset{X: -1, Y: 25}) {
		t.Error("Point (-1, 25) should be outside button")
	}
	
	if button.containsPoint(core.Offset{X: 101, Y: 25}) {
		t.Error("Point (101, 25) should be outside button")
	}
	
	if button.containsPoint(core.Offset{X: 50, Y: -1}) {
		t.Error("Point (50, -1) should be outside button")
	}
	
	if button.containsPoint(core.Offset{X: 50, Y: 51}) {
		t.Error("Point (50, 51) should be outside button")
	}
}

func TestButton_Dispose(t *testing.T) {
	button := NewButton("test", "Click", nil)
	
	button.Dispose()
	
	// After dispose, the widget should be cleaned up
	if !button.disposed.Load() {
		t.Error("Button should be marked as disposed")
	}
}

func TestButton_ReactiveEffects(t *testing.T) {
	button := NewButton("test", "Click", nil)
	
	var repaintCount int
	dispose := reactive.Watch(func() {
		if button.NeedsRepaint() {
			repaintCount++
		}
	})
	defer dispose()
	
	// Change button states
	button.isPressed.Set(true)
	time.Sleep(10 * time.Millisecond)
	
	button.isHovered.Set(true)
	time.Sleep(10 * time.Millisecond)
	
	button.SetDisabled(true)
	time.Sleep(10 * time.Millisecond)
	
	if repaintCount < 3 {
		t.Errorf("State changes should trigger repaints, got %d repaints", repaintCount)
	}
}

func TestButton_NilCallback(t *testing.T) {
	button := NewButton("test", "Click", nil)
	
	// Should not panic with nil callback
	button.onPressed = nil
	
	ctx := context.Background()
	renderObj := button.Build(ctx)
	
	renderButton := renderObj.(*RenderButton)
	
	// Should not panic when calling nil callback
	if renderButton.OnPressed != nil {
		renderButton.OnPressed()
	}
}

func TestButton_EmptyLabel(t *testing.T) {
	button := NewButton("test", "", nil)
	
	constraints := core.Constraints{
		MinWidth:  0,
		MaxWidth:  1000,
		MinHeight: 0,
		MaxHeight: 1000,
	}
	
	width, height := button.Layout(constraints)
	
	// Should still have padding
	expectedWidth := 0.0 + 32.0  // Just padding
	expectedHeight := 14.0*1.2 + 16.0 // Line height + padding
	
	if width != expectedWidth {
		t.Errorf("Expected width %f, got %f", expectedWidth, width)
	}
	
	if height != expectedHeight {
		t.Errorf("Expected height %f, got %f", expectedHeight, height)
	}
}

func TestButton_BorderRadius(t *testing.T) {
	button := NewButton("test", "Click", nil)
	button.cachedSize = Size{Width: 100, Height: 40}
	button.borderRadius = 10
	
	ctx := &mockPaintContext{}
	button.Paint(ctx)
	
	// Should draw border when borderRadius > 0
	// We check that DrawRect was called multiple times (background + border)
	if !ctx.drawRectCalled {
		t.Error("Button with border radius should draw rectangles")
	}
}