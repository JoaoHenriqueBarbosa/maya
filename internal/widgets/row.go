package widgets

import (
	"context"
	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/reactive"
)

// Row arranges children horizontally
type Row struct {
	*BaseWidget
	
	// Layout properties
	gap         *reactive.Signal[float64]
	alignment   *reactive.Signal[Alignment]
}

// NewRow creates a new row widget
func NewRow(id string, children ...WidgetImpl) *Row {
	r := &Row{
		BaseWidget: NewBaseWidget(id, "Row"),
		gap:        reactive.NewSignal(10.0),
		alignment:  reactive.NewSignal(AlignmentCenterLeft),
	}
	
	// Add children
	for _, child := range children {
		r.AddChild(child)
	}
	
	return r
}

// SetGap sets the gap between children
func (r *Row) SetGap(gap float64) {
	r.gap.Set(gap)
	r.MarkNeedsLayout()
}

// Build creates the render object
func (r *Row) Build(ctx context.Context) RenderObject {
	return &RenderBox{
		Size: r.cachedSize,
	}
}

// Layout arranges children horizontally
func (r *Row) Layout(constraints core.Constraints) (width, height float64) {
	gap := r.gap.Get()
	children := r.Children()
	
	if len(children) == 0 {
		r.cachedSize = Size{Width: constraints.MinWidth, Height: constraints.MinHeight}
		return constraints.MinWidth, constraints.MinHeight
	}
	
	// Calculate total width and max height
	totalWidth := 0.0
	maxHeight := 0.0
	
	for i, child := range children {
		if coreWidget, ok := child.(core.Widget); ok {
			// Give each child flexible width constraint
			childConstraints := core.Constraints{
				MinWidth:  0,
				MaxWidth:  constraints.MaxWidth,
				MinHeight: 0,
				MaxHeight: constraints.MaxHeight,
			}
			
			childWidth, childHeight := coreWidget.Layout(childConstraints)
			
			totalWidth += childWidth
			if i > 0 {
				totalWidth += gap
			}
			
			if childHeight > maxHeight {
				maxHeight = childHeight
			}
		}
	}
	
	// Apply constraints
	width = totalWidth
	if width < constraints.MinWidth {
		width = constraints.MinWidth
	}
	if width > constraints.MaxWidth {
		width = constraints.MaxWidth
	}
	
	height = maxHeight
	if height < constraints.MinHeight {
		height = constraints.MinHeight
	}
	if height > constraints.MaxHeight {
		height = constraints.MaxHeight
	}
	
	r.cachedSize = Size{Width: width, Height: height}
	r.needsLayout.Set(false)
	
	return width, height
}

// Paint renders the row (delegates to children)
func (r *Row) Paint(context core.PaintContext) {
	// Row itself doesn't paint, only its children do
	children := r.Children()
	for _, child := range children {
		if coreWidget, ok := child.(core.Widget); ok {
			coreWidget.Paint(context)
		}
	}
	r.needsRepaint.Set(false)
}

// HandleEvent processes events
func (r *Row) HandleEvent(event core.Event) bool {
	// Propagate to children
	children := r.Children()
	for i := len(children) - 1; i >= 0; i-- {
		if coreWidget, ok := children[i].(core.Widget); ok {
			if coreWidget.HandleEvent(event) {
				return true
			}
		}
	}
	return false
}