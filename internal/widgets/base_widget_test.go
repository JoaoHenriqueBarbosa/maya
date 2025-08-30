package widgets

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/reactive"
)

// mockWidget implements WidgetImpl for testing
type mockWidget struct {
	*BaseWidget
	layoutCalled   bool
	paintCalled    bool
	handleCalled   bool
	buildCalled    bool
	lastEvent      core.Event
	layoutReturn   struct{ width, height float64 }
}

func newMockWidget(id string) *mockWidget {
	m := &mockWidget{
		BaseWidget: NewBaseWidget(id, "Mock"),
	}
	m.layoutReturn.width = 100
	m.layoutReturn.height = 50
	return m
}

func (m *mockWidget) Layout(constraints core.Constraints) (width, height float64) {
	m.layoutCalled = true
	return m.layoutReturn.width, m.layoutReturn.height
}

func (m *mockWidget) Paint(context core.PaintContext) {
	m.paintCalled = true
}

func (m *mockWidget) HandleEvent(event core.Event) bool {
	m.handleCalled = true
	m.lastEvent = event
	return true
}

func (m *mockWidget) Build(ctx context.Context) RenderObject {
	m.buildCalled = true
	return &RenderBox{Size: Size{Width: 100, Height: 50}}
}

// mockPaintContext implements core.PaintContext for testing
type mockPaintContext struct {
	drawRectCalled bool
	drawTextCalled bool
	lastBounds     core.Bounds
	lastPaint      core.Paint
}

func (m *mockPaintContext) DrawRect(bounds core.Bounds, paint core.Paint) {
	m.drawRectCalled = true
	m.lastBounds = bounds
	m.lastPaint = paint
}

func (m *mockPaintContext) DrawText(text string, offset core.Offset, paint core.Paint) {
	m.drawTextCalled = true
}

func (m *mockPaintContext) DrawPath(path []core.Offset, paint core.Paint) {}
func (m *mockPaintContext) PushTransform(transform core.Transform)         {}
func (m *mockPaintContext) PopTransform()                                  {}
func (m *mockPaintContext) PushClip(bounds core.Bounds)                    {}
func (m *mockPaintContext) PopClip()                                       {}

func TestNewBaseWidget(t *testing.T) {
	widget := NewBaseWidget("test-id", "TestWidget")
	
	if widget.ID() != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", widget.ID())
	}
	
	if widget.Type() != "TestWidget" {
		t.Errorf("Expected Type 'TestWidget', got %s", widget.Type())
	}
	
	if widget.Parent() != nil {
		t.Error("Expected Parent to be nil for new widget")
	}
	
	if len(widget.Children()) != 0 {
		t.Error("Expected no children for new widget")
	}
	
	if widget.NeedsRepaint() {
		t.Error("New widget should not need repaint")
	}
}

func TestBaseWidget_ParentChild(t *testing.T) {
	parent := NewBaseWidget("parent", "Parent")
	child1 := newMockWidget("child1")
	child2 := newMockWidget("child2")
	
	// Test AddChild
	parent.AddChild(child1)
	if len(parent.Children()) != 1 {
		t.Errorf("Expected 1 child, got %d", len(parent.Children()))
	}
	
	if child1.Parent() != parent {
		t.Error("Child's parent should be set")
	}
	
	// Add second child
	parent.AddChild(child2)
	if len(parent.Children()) != 2 {
		t.Errorf("Expected 2 children, got %d", len(parent.Children()))
	}
	
	// Test RemoveChild
	removed := parent.RemoveChild(child1)
	if !removed {
		t.Error("RemoveChild should return true for existing child")
	}
	
	if len(parent.Children()) != 1 {
		t.Errorf("Expected 1 child after removal, got %d", len(parent.Children()))
	}
	
	if child1.Parent() != nil {
		t.Error("Removed child's parent should be nil")
	}
	
	// Try removing non-existent child
	removed = parent.RemoveChild(child1)
	if removed {
		t.Error("RemoveChild should return false for non-existent child")
	}
}

func TestBaseWidget_Props(t *testing.T) {
	widget := NewBaseWidget("test", "Test")
	
	props := Props{
		"color": "blue",
		"size":  42,
	}
	
	widget.SetProps(props)
	
	gotProps := widget.GetProps()
	if gotProps["color"] != "blue" {
		t.Errorf("Expected color 'blue', got %v", gotProps["color"])
	}
	
	if gotProps["size"] != 42 {
		t.Errorf("Expected size 42, got %v", gotProps["size"])
	}
	
	// Test UpdateProp
	widget.UpdateProp("color", "red")
	gotProps = widget.GetProps()
	if gotProps["color"] != "red" {
		t.Errorf("Expected updated color 'red', got %v", gotProps["color"])
	}
}

func TestBaseWidget_Layout(t *testing.T) {
	widget := NewBaseWidget("test", "Test")
	
	constraints := core.Constraints{
		MinWidth:  50,
		MaxWidth:  200,
		MinHeight: 30,
		MaxHeight: 150,
	}
	
	width, height := widget.Layout(constraints)
	
	// Default implementation should take maximum available space
	if width != 200 {
		t.Errorf("Expected width 200, got %f", width)
	}
	
	if height != 150 {
		t.Errorf("Expected height 150, got %f", height)
	}
	
	// Test caching
	width2, height2 := widget.Layout(constraints)
	if width != width2 || height != height2 {
		t.Error("Layout should return cached results for same constraints")
	}
}

func TestBaseWidget_Paint(t *testing.T) {
	widget := NewBaseWidget("test", "Test")
	ctx := &mockPaintContext{}
	
	// Paint without background color
	widget.Paint(ctx)
	if ctx.drawRectCalled {
		t.Error("Should not draw rect without backgroundColor prop")
	}
	
	// Set background color
	widget.UpdateProp("backgroundColor", core.Color{R: 255, G: 0, B: 0, A: 255})
	widget.cachedSize = Size{Width: 100, Height: 50}
	
	widget.Paint(ctx)
	if !ctx.drawRectCalled {
		t.Error("Should draw rect with backgroundColor prop")
	}
	
	if ctx.lastPaint.Color.R != 255 {
		t.Errorf("Expected red color, got %v", ctx.lastPaint.Color)
	}
}

func TestBaseWidget_Paint_WithChildren(t *testing.T) {
	parent := NewBaseWidget("parent", "Parent")
	child := newMockWidget("child")
	
	parent.AddChild(child)
	
	ctx := &mockPaintContext{}
	parent.Paint(ctx)
	
	if !child.paintCalled {
		t.Error("Child's Paint should be called")
	}
}

func TestBaseWidget_HandleEvent(t *testing.T) {
	parent := NewBaseWidget("parent", "Parent")
	child1 := newMockWidget("child1")
	child2 := newMockWidget("child2")
	
	parent.AddChild(child1)
	parent.AddChild(child2)
	
	event := &core.MockEvent{EventType: "test"}
	
	// Event should be handled by children in reverse order
	handled := parent.HandleEvent(event)
	if !handled {
		t.Error("Event should be handled by child")
	}
	
	if !child2.handleCalled {
		t.Error("Child2 should handle event first (reverse order)")
	}
}

func TestBaseWidget_Visibility(t *testing.T) {
	widget := NewBaseWidget("test", "Test")
	
	// Default should be visible
	if !widget.IsVisible() {
		t.Error("Widget should be visible by default")
	}
	
	// Set visible to false
	widget.UpdateProp("visible", false)
	if widget.IsVisible() {
		t.Error("Widget should not be visible after setting prop")
	}
}

func TestBaseWidget_Opacity(t *testing.T) {
	widget := NewBaseWidget("test", "Test")
	
	// Default opacity should be 1.0
	if widget.GetOpacity() != 1.0 {
		t.Errorf("Expected default opacity 1.0, got %f", widget.GetOpacity())
	}
	
	// Set custom opacity
	widget.UpdateProp("opacity", 0.5)
	if widget.GetOpacity() != 0.5 {
		t.Errorf("Expected opacity 0.5, got %f", widget.GetOpacity())
	}
}

func TestBaseWidget_IntrinsicDimensions(t *testing.T) {
	widget := NewBaseWidget("test", "Test")
	
	// Default implementation returns 0
	if widget.GetIntrinsicWidth(100) != 0 {
		t.Error("Default intrinsic width should be 0")
	}
	
	if widget.GetIntrinsicHeight(100) != 0 {
		t.Error("Default intrinsic height should be 0")
	}
}

func TestBaseWidget_NeedsRepaint(t *testing.T) {
	parent := NewBaseWidget("parent", "Parent")
	child := NewBaseWidget("child", "Child")
	
	parent.AddChild(child)
	
	// Mark child as needing repaint
	child.MarkNeedsRepaint()
	
	if !child.NeedsRepaint() {
		t.Error("Child should need repaint")
	}
	
	// Parent should also need repaint (propagation)
	if !parent.NeedsRepaint() {
		t.Error("Parent should need repaint when child needs repaint")
	}
}

func TestBaseWidget_MarkNeedsLayout(t *testing.T) {
	parent := NewBaseWidget("parent", "Parent")
	child := NewBaseWidget("child", "Child")
	
	parent.AddChild(child)
	
	// Mark child as needing layout
	child.MarkNeedsLayout()
	
	// Parent should also need layout (propagation)
	if !parent.needsLayout.Get() {
		t.Error("Parent should need layout when child needs layout")
	}
}

func TestBaseWidget_Init(t *testing.T) {
	widget := NewBaseWidget("test", "Test")
	
	// First init
	widget.Init()
	if !widget.initialized.Load() {
		t.Error("Widget should be initialized")
	}
	
	if !widget.needsLayout.Get() {
		t.Error("Widget should need layout after init")
	}
	
	if !widget.needsRepaint.Get() {
		t.Error("Widget should need repaint after init")
	}
	
	// Second init should not re-initialize
	widget.needsLayout.Set(false)
	widget.needsRepaint.Set(false)
	widget.Init()
	
	if widget.needsLayout.Get() || widget.needsRepaint.Get() {
		t.Error("Second init should not change state")
	}
}

func TestBaseWidget_Dispose(t *testing.T) {
	parent := NewBaseWidget("parent", "Parent")
	child := NewBaseWidget("child", "Child")
	
	parent.AddChild(child)
	
	// Setup an effect to test cleanup
	var effectCleaned bool
	parent.effects = append(parent.effects, func() {
		effectCleaned = true
	})
	
	parent.Dispose()
	
	if !parent.disposed.Load() {
		t.Error("Widget should be marked as disposed")
	}
	
	if !effectCleaned {
		t.Error("Effects should be cleaned up")
	}
	
	if len(parent.Children()) != 0 {
		t.Error("Children should be cleared after dispose")
	}
	
	// Second dispose should not panic
	parent.Dispose()
}

func TestBaseWidget_Build(t *testing.T) {
	widget := NewBaseWidget("test", "Test")
	widget.cachedSize = Size{Width: 100, Height: 50}
	
	ctx := context.Background()
	renderObj := widget.Build(ctx)
	
	if renderObj == nil {
		t.Fatal("Build should return a RenderObject")
	}
	
	box, ok := renderObj.(*RenderBox)
	if !ok {
		t.Fatal("Default Build should return RenderBox")
	}
	
	if box.Size.Width != 100 || box.Size.Height != 50 {
		t.Errorf("RenderBox should have cached size, got %v", box.Size)
	}
}

func TestBaseWidget_ConcurrentAccess(t *testing.T) {
	widget := NewBaseWidget("test", "Test")
	
	var wg sync.WaitGroup
	iterations := 100
	
	// Concurrent adds
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			child := NewBaseWidget(string(rune(idx)), "Child")
			widget.AddChild(child)
		}(i)
	}
	
	// Concurrent reads
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = widget.Children()
		}()
	}
	
	wg.Wait()
	
	// Should have all children
	if len(widget.Children()) != iterations {
		t.Errorf("Expected %d children, got %d", iterations, len(widget.Children()))
	}
}

func TestBaseWidget_ReactiveEffects(t *testing.T) {
	widget := NewBaseWidget("test", "Test")
	
	// Set visible property and verify effect triggers repaint
	var repaintTriggered bool
	dispose := reactive.Watch(func() {
		if widget.NeedsRepaint() {
			repaintTriggered = true
		}
	})
	defer dispose()
	
	widget.UpdateProp("visible", false)
	
	// Give reactive system time to process
	time.Sleep(10 * time.Millisecond)
	
	if !repaintTriggered {
		t.Error("Changing visible prop should trigger repaint")
	}
}

func TestBaseWidget_SetParent(t *testing.T) {
	child := NewBaseWidget("child", "Child")
	parent := NewBaseWidget("parent", "Parent")
	
	child.SetParent(parent)
	
	if child.Parent() != parent {
		t.Error("Parent should be set")
	}
	
	// Set to nil
	child.SetParent(nil)
	if child.Parent() != nil {
		t.Error("Parent should be nil")
	}
}

func TestBaseWidget_LayoutCaching(t *testing.T) {
	widget := NewBaseWidget("test", "Test")
	
	constraints1 := core.Constraints{
		MinWidth: 0, MaxWidth: 100,
		MinHeight: 0, MaxHeight: 100,
	}
	
	constraints2 := core.Constraints{
		MinWidth: 0, MaxWidth: 200,
		MinHeight: 0, MaxHeight: 200,
	}
	
	// First layout
	w1, h1 := widget.Layout(constraints1)
	
	// Same constraints should return cached
	w2, h2 := widget.Layout(constraints1)
	if w1 != w2 || h1 != h2 {
		t.Error("Should return cached results for same constraints")
	}
	
	// Different constraints should recalculate
	w3, h3 := widget.Layout(constraints2)
	if w3 == w1 || h3 == h1 {
		t.Error("Should recalculate for different constraints")
	}
	
	if w3 != 200 || h3 != 200 {
		t.Errorf("Expected 200x200, got %fx%f", w3, h3)
	}
}