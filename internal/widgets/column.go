package widgets

import (
	"context"
	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/reactive"
)

// Column arranges children vertically
type Column struct {
	*BaseWidget
	
	// Layout properties
	gap         *reactive.Signal[float64]
	alignment   *reactive.Signal[Alignment]
}

// NewColumn creates a new column widget
func NewColumn(id string, children ...WidgetImpl) *Column {
	c := &Column{
		BaseWidget: NewBaseWidget(id, "Column"),
		gap:        reactive.NewSignal(10.0),
		alignment:  reactive.NewSignal(AlignmentTopLeft),
	}
	
	// Add children
	for _, child := range children {
		c.AddChild(child)
	}
	
	return c
}

// SetGap sets the gap between children
func (c *Column) SetGap(gap float64) {
	c.gap.Set(gap)
	c.MarkNeedsLayout()
}

// Build creates the render object
func (c *Column) Build(ctx context.Context) RenderObject {
	return &RenderBox{
		Size: c.cachedSize,
	}
}

// Layout arranges children vertically
func (c *Column) Layout(constraints core.Constraints) (width, height float64) {
	gap := c.gap.Get()
	children := c.Children()
	
	if len(children) == 0 {
		c.cachedSize = Size{Width: constraints.MinWidth, Height: constraints.MinHeight}
		return constraints.MinWidth, constraints.MinHeight
	}
	
	// Calculate total height and max width
	totalHeight := 0.0
	maxWidth := 0.0
	
	for i, child := range children {
		if coreWidget, ok := child.(core.Widget); ok {
			// Give each child the full width constraint
			childConstraints := core.Constraints{
				MinWidth:  0,
				MaxWidth:  constraints.MaxWidth,
				MinHeight: 0,
				MaxHeight: constraints.MaxHeight,
			}
			
			childWidth, childHeight := coreWidget.Layout(childConstraints)
			
			totalHeight += childHeight
			if i > 0 {
				totalHeight += gap
			}
			
			if childWidth > maxWidth {
				maxWidth = childWidth
			}
		}
	}
	
	// Apply constraints
	width = maxWidth
	if width < constraints.MinWidth {
		width = constraints.MinWidth
	}
	if width > constraints.MaxWidth {
		width = constraints.MaxWidth
	}
	
	height = totalHeight
	if height < constraints.MinHeight {
		height = constraints.MinHeight
	}
	if height > constraints.MaxHeight {
		height = constraints.MaxHeight
	}
	
	c.cachedSize = Size{Width: width, Height: height}
	c.needsLayout.Set(false)
	
	return width, height
}

// Paint renders the column (delegates to children)
func (c *Column) Paint(context core.PaintContext) {
	// Column itself doesn't paint, only its children do
	children := c.Children()
	for _, child := range children {
		if coreWidget, ok := child.(core.Widget); ok {
			coreWidget.Paint(context)
		}
	}
	c.needsRepaint.Set(false)
}

// HandleEvent processes events
func (c *Column) HandleEvent(event core.Event) bool {
	// Propagate to children
	children := c.Children()
	for i := len(children) - 1; i >= 0; i-- {
		if coreWidget, ok := children[i].(core.Widget); ok {
			if coreWidget.HandleEvent(event) {
				return true
			}
		}
	}
	return false
}