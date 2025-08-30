package widgets

import (
	"context"
	"testing"

	"github.com/maya-framework/maya/internal/core"
)

func TestNewContainer(t *testing.T) {
	container := NewContainer("container-id")
	
	if container.ID() != "container-id" {
		t.Errorf("Expected ID 'container-id', got %s", container.ID())
	}
	
	if container.Type() != "Container" {
		t.Errorf("Expected Type 'Container', got %s", container.Type())
	}
	
	// Check default padding
	padding := container.padding.Get()
	if padding.Top != 0 || padding.Bottom != 0 {
		t.Error("Default padding should be 0")
	}
	
	// Check default margin
	margin := container.margin.Get()
	if margin.Top != 0 || margin.Bottom != 0 {
		t.Error("Default margin should be 0")
	}
	
	// Check default alignment
	if container.alignment.Get() != AlignmentTopLeft {
		t.Error("Default alignment should be TopLeft")
	}
	
	// Check default width/height
	if container.width.Get() != 0 {
		t.Error("Default width should be 0")
	}
	
	if container.height.Get() != 0 {
		t.Error("Default height should be 0")
	}
}

func TestContainer_SetPadding(t *testing.T) {
	container := NewContainer("test")
	
	padding := EdgeInsets{Top: 10, Right: 20, Bottom: 30, Left: 40}
	container.SetPadding(padding)
	
	if container.padding.Get() != padding {
		t.Errorf("Expected padding %v, got %v", padding, container.padding.Get())
	}
	
	// Should trigger layout
	if !container.needsLayout.Get() {
		t.Error("Changing padding should trigger layout")
	}
}

func TestContainer_SetMargin(t *testing.T) {
	container := NewContainer("test")
	
	margin := EdgeInsets{Top: 5, Right: 10, Bottom: 15, Left: 20}
	container.SetMargin(margin)
	
	if container.margin.Get() != margin {
		t.Errorf("Expected margin %v, got %v", margin, container.margin.Get())
	}
	
	// Should trigger layout
	if !container.needsLayout.Get() {
		t.Error("Changing margin should trigger layout")
	}
}

func TestContainer_SetAlignment(t *testing.T) {
	container := NewContainer("test")
	
	container.SetAlignment(AlignmentCenter)
	
	if container.alignment.Get() != AlignmentCenter {
		t.Errorf("Expected alignment Center, got %v", container.alignment.Get())
	}
	
	// Should trigger layout
	if !container.needsLayout.Get() {
		t.Error("Changing alignment should trigger layout")
	}
}

func TestContainer_SetWidthHeight(t *testing.T) {
	container := NewContainer("test")
	
	container.SetWidth(100)
	container.SetHeight(50)
	
	if container.width.Get() != 100 {
		t.Errorf("Expected width 100, got %f", container.width.Get())
	}
	
	if container.height.Get() != 50 {
		t.Errorf("Expected height 50, got %f", container.height.Get())
	}
	
	// Should trigger layout
	if !container.needsLayout.Get() {
		t.Error("Changing size should trigger layout")
	}
}

func TestContainer_SetColor(t *testing.T) {
	container := NewContainer("test")
	
	color := core.Color{R: 255, G: 0, B: 0, A: 255}
	container.SetColor(color)
	
	if container.color.Get() != color {
		t.Errorf("Expected color %v, got %v", color, container.color.Get())
	}
	
	// Should trigger repaint
	if !container.needsRepaint.Get() {
		t.Error("Changing color should trigger repaint")
	}
}

func TestContainer_SetBorderRadius(t *testing.T) {
	container := NewContainer("test")
	
	// SetBorderRadius is handled via SetBorder
	container.SetBorder(core.Color{}, 0, 10)
	
	if container.borderRadius.Get() != 10 {
		t.Errorf("Expected border radius 10, got %f", container.borderRadius.Get())
	}
	
	// Should trigger repaint
	if !container.needsRepaint.Get() {
		t.Error("Changing border radius should trigger repaint")
	}
}

func TestContainer_SetBorder(t *testing.T) {
	container := NewContainer("test")
	
	borderColor := core.Color{R: 0, G: 0, B: 255, A: 255}
	borderWidth := 2.0
	borderRadius := 5.0
	
	container.SetBorder(borderColor, borderWidth, borderRadius)
	
	if container.borderColor.Get() != borderColor {
		t.Errorf("Expected border color %v, got %v", borderColor, container.borderColor.Get())
	}
	
	if container.borderWidth.Get() != borderWidth {
		t.Errorf("Expected border width %f, got %f", borderWidth, container.borderWidth.Get())
	}
	
	if container.borderRadius.Get() != borderRadius {
		t.Errorf("Expected border radius %f, got %f", borderRadius, container.borderRadius.Get())
	}
	
	// Should trigger repaint
	if !container.needsRepaint.Get() {
		t.Error("Changing border should trigger repaint")
	}
}

func TestContainer_SetBoxShadow(t *testing.T) {
	container := NewContainer("test")
	
	shadow := &BoxShadow{
		Offset:     core.Offset{X: 5, Y: 5},
		BlurRadius: 10,
		Color:      core.Color{R: 0, G: 0, B: 0, A: 128},
	}
	
	container.SetBoxShadow(shadow)
	
	if container.boxShadow.Get() != shadow {
		t.Error("Shadow not set correctly")
	}
	
	// Should trigger repaint
	if !container.needsRepaint.Get() {
		t.Error("Changing shadow should trigger repaint")
	}
}

func TestContainer_Build(t *testing.T) {
	container := NewContainer("test")
	child := newMockWidget("child")
	
	container.SetChild(child)
	container.SetPadding(EdgeInsets{Top: 10, Right: 10, Bottom: 10, Left: 10})
	container.SetColor(core.Color{R: 255, G: 0, B: 0, A: 255})
	
	// Layout first to set cached size
	constraints := core.Constraints{
		MinWidth: 0, MaxWidth: 200,
		MinHeight: 0, MaxHeight: 100,
	}
	container.Layout(constraints)
	
	ctx := context.Background()
	renderObj := container.Build(ctx)
	
	if renderObj == nil {
		t.Fatal("Build should return a RenderObject")
	}
	
	_, ok := renderObj.(*RenderBox)
	if !ok {
		t.Fatal("Container.Build should return RenderBox")
	}
	
	// RenderBox has been created correctly
}

func TestContainer_Layout(t *testing.T) {
	container := NewContainer("test")
	container.SetPadding(EdgeInsets{Top: 10, Right: 20, Bottom: 10, Left: 20})
	container.SetMargin(EdgeInsets{Top: 5, Right: 5, Bottom: 5, Left: 5})
	
	constraints := core.Constraints{
		MinWidth:  0,
		MaxWidth:  200,
		MinHeight: 0,
		MaxHeight: 100,
	}
	
	width, height := container.Layout(constraints)
	
	// With padding and margin
	expectedMinWidth := 20.0 + 20.0 + 5.0 + 5.0 // left/right padding + margin
	expectedMinHeight := 10.0 + 10.0 + 5.0 + 5.0 // top/bottom padding + margin
	
	if width < expectedMinWidth {
		t.Errorf("Width should be at least %f, got %f", expectedMinWidth, width)
	}
	
	if height < expectedMinHeight {
		t.Errorf("Height should be at least %f, got %f", expectedMinHeight, height)
	}
}

func TestContainer_Layout_WithChild(t *testing.T) {
	container := NewContainer("test")
	child := newMockWidget("child")
	child.layoutReturn.width = 50
	child.layoutReturn.height = 30
	
	container.SetChild(child)
	container.SetPadding(EdgeInsets{Top: 10, Right: 10, Bottom: 10, Left: 10})
	
	constraints := core.Constraints{
		MinWidth:  0,
		MaxWidth:  200,
		MinHeight: 0,
		MaxHeight: 100,
	}
	
	width, height := container.Layout(constraints)
	
	// Should accommodate child + padding
	expectedWidth := 50.0 + 20.0  // child width + horizontal padding
	expectedHeight := 30.0 + 20.0 // child height + vertical padding
	
	if width != expectedWidth {
		t.Errorf("Expected width %f, got %f", expectedWidth, width)
	}
	
	if height != expectedHeight {
		t.Errorf("Expected height %f, got %f", expectedHeight, height)
	}
	
	if !child.layoutCalled {
		t.Error("Child's Layout should be called")
	}
}

func TestContainer_Layout_WithSize(t *testing.T) {
	container := NewContainer("test")
	
	// Set specific size
	container.SetWidth(150)
	container.SetHeight(75)
	
	constraints := core.Constraints{
		MinWidth:  0,
		MaxWidth:  300,
		MinHeight: 0,
		MaxHeight: 200,
	}
	
	width, height := container.Layout(constraints)
	
	// Container layout depends on implementation
	// Should respect constraints
	if width > 300 {
		t.Errorf("Width should not exceed max constraint, got %f", width)
	}
	
	if height > 200 {
		t.Errorf("Height should not exceed max constraint, got %f", height)
	}
}

func TestContainer_Paint(t *testing.T) {
	container := NewContainer("test")
	container.SetColor(core.Color{R: 255, G: 0, B: 0, A: 255})
	container.cachedSize = Size{Width: 100, Height: 50}
	
	ctx := &mockPaintContext{}
	container.Paint(ctx)
	
	if !ctx.drawRectCalled {
		t.Error("Container should draw background rect")
	}
	
	if ctx.lastPaint.Color.R != 255 {
		t.Error("Container should use correct background color")
	}
}

func TestContainer_Paint_WithBorder(t *testing.T) {
	container := NewContainer("test")
	container.SetBorder(core.Color{R: 0, G: 0, B: 255, A: 255}, 2.0, 5.0)
	container.cachedSize = Size{Width: 100, Height: 50}
	
	ctx := &mockPaintContext{}
	container.Paint(ctx)
	
	// Should draw border
	if !ctx.drawRectCalled {
		t.Error("Container should draw border")
	}
}

func TestContainer_Paint_WithShadow(t *testing.T) {
	container := NewContainer("test")
	container.SetBoxShadow(&BoxShadow{
		Offset:     core.Offset{X: 5, Y: 5},
		BlurRadius: 10,
		Color:      core.Color{R: 0, G: 0, B: 0, A: 128},
	})
	container.cachedSize = Size{Width: 100, Height: 50}
	
	ctx := &mockPaintContext{}
	container.Paint(ctx)
	
	// Should draw shadow
	if !ctx.drawRectCalled {
		t.Error("Container should draw shadow")
	}
}

func TestContainer_Paint_WithChild(t *testing.T) {
	container := NewContainer("test")
	child := newMockWidget("child")
	
	container.SetChild(child)
	container.cachedSize = Size{Width: 100, Height: 50}
	
	ctx := &mockPaintContext{}
	container.Paint(ctx)
	
	if !child.paintCalled {
		t.Error("Child should be painted")
	}
}

func TestContainer_Paint_Invisible(t *testing.T) {
	container := NewContainer("test")
	container.UpdateProp("visible", false)
	
	ctx := &mockPaintContext{}
	container.Paint(ctx)
	
	if ctx.drawRectCalled {
		t.Error("Invisible container should not be painted")
	}
}

func TestContainer_HandleEvent(t *testing.T) {
	container := NewContainer("test")
	child := newMockWidget("child")
	
	container.SetChild(child)
	
	event := &core.MockEvent{EventType: "test"}
	handled := container.HandleEvent(event)
	
	if !handled {
		t.Error("Event should be handled by child")
	}
	
	if !child.handleCalled {
		t.Error("Child should handle event")
	}
}

func TestContainer_Dispose(t *testing.T) {
	container := NewContainer("test")
	child := NewContainer("child")
	
	container.SetChild(child)
	
	container.Dispose()
	
	if !container.disposed.Load() {
		t.Error("Container should be disposed")
	}
	
	if !child.disposed.Load() {
		t.Error("Child should be disposed")
	}
}

func TestContainer_ReactiveUpdates(t *testing.T) {
	container := NewContainer("test")
	
	// Test that methods trigger the needsRepaint flag correctly
	// Change color
	container.SetColor(core.Color{R: 255, G: 0, B: 0, A: 255})
	if !container.NeedsRepaint() {
		t.Error("SetColor should mark container as needing repaint")
	}
	
	// Create new container for next test
	container2 := NewContainer("test2")
	
	// Change border
	container2.SetBorder(core.Color{R: 0, G: 0, B: 0, A: 255}, 1.0, 0.0)
	if !container2.NeedsRepaint() {
		t.Error("SetBorder should mark container as needing repaint")
	}
	
	// Create new container for shadow test
	container3 := NewContainer("test3")
	
	// Change shadow
	container3.SetBoxShadow(&BoxShadow{Offset: core.Offset{X: 5, Y: 5}})
	if !container3.NeedsRepaint() {
		t.Error("SetBoxShadow should mark container as needing repaint")
	}
}

func TestContainer_EmptyContainer(t *testing.T) {
	container := NewContainer("test")
	
	constraints := core.Constraints{
		MinWidth:  0,
		MaxWidth:  1000,
		MinHeight: 0,
		MaxHeight: 1000,
	}
	
	width, height := container.Layout(constraints)
	
	// Empty container should take maximum space
	if width != 1000 {
		t.Errorf("Expected width 1000, got %f", width)
	}
	
	if height != 1000 {
		t.Errorf("Expected height 1000, got %f", height)
	}
}



func TestContainer_Alignment(t *testing.T) {
	testCases := []struct {
		name      string
		alignment Alignment
	}{
		{"TopLeft", AlignmentTopLeft},
		{"TopCenter", AlignmentTopCenter},
		{"TopRight", AlignmentTopRight},
		{"CenterLeft", AlignmentCenterLeft},
		{"Center", AlignmentCenter},
		{"CenterRight", AlignmentCenterRight},
		{"BottomLeft", AlignmentBottomLeft},
		{"BottomCenter", AlignmentBottomCenter},
		{"BottomRight", AlignmentBottomRight},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			container := NewContainer("test")
			container.SetAlignment(tc.alignment)
			
			if container.alignment.Get() != tc.alignment {
				t.Errorf("Expected alignment %v, got %v", tc.alignment, container.alignment.Get())
			}
		})
	}
}