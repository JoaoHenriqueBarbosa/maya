package widgets

import (
	"context"
	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/reactive"
)

// Container is a widget that contains a single child with optional styling
type Container struct {
	*BaseWidget
	
	// Container properties
	width           *reactive.Signal[float64]
	height          *reactive.Signal[float64]
	color           *reactive.Signal[core.Color]
	padding         *reactive.Signal[EdgeInsets]
	margin          *reactive.Signal[EdgeInsets]
	alignment       *reactive.Signal[Alignment]
	borderRadius    *reactive.Signal[float64]
	borderColor     *reactive.Signal[core.Color]
	borderWidth     *reactive.Signal[float64]
	boxShadow       *reactive.Signal[*BoxShadow]
	
	// Single child
	child WidgetImpl
}

// Alignment specifies how to align a child within its parent
type Alignment struct {
	X float64 // -1.0 (left) to 1.0 (right)
	Y float64 // -1.0 (top) to 1.0 (bottom)
}

// Common alignments
var (
	AlignmentTopLeft      = Alignment{-1, -1}
	AlignmentTopCenter    = Alignment{0, -1}
	AlignmentTopRight     = Alignment{1, -1}
	AlignmentCenterLeft   = Alignment{-1, 0}
	AlignmentCenter       = Alignment{0, 0}
	AlignmentCenterRight  = Alignment{1, 0}
	AlignmentBottomLeft   = Alignment{-1, 1}
	AlignmentBottomCenter = Alignment{0, 1}
	AlignmentBottomRight  = Alignment{1, 1}
)

// BoxShadow represents a box shadow effect
type BoxShadow struct {
	Color      core.Color
	Offset     core.Offset
	BlurRadius float64
	SpreadRadius float64
}

// NewContainer creates a new container widget
func NewContainer(id string) *Container {
	c := &Container{
		BaseWidget:   NewBaseWidget(id, "Container"),
		width:        reactive.NewSignal(0.0),
		height:       reactive.NewSignal(0.0),
		color:        reactive.NewSignal(ColorTransparent),
		padding:      reactive.NewSignal(EdgeInsets{}),
		margin:       reactive.NewSignal(EdgeInsets{}),
		alignment:    reactive.NewSignal(AlignmentTopLeft),
		borderRadius: reactive.NewSignal(0.0),
		borderColor:  reactive.NewSignal(ColorTransparent),
		borderWidth:  reactive.NewSignal(0.0),
		boxShadow:    reactive.NewSignal[*BoxShadow](nil),
	}
	
	// Setup reactive effects
	c.setupContainerEffects()
	
	return c
}

// setupContainerEffects creates reactive effects for container
func (c *Container) setupContainerEffects() {
	// Watch for visual changes
	dispose := reactive.Watch(func() {
		_ = c.color.Get()
		_ = c.borderColor.Get()
		_ = c.borderWidth.Get()
		_ = c.boxShadow.Get()
		c.MarkNeedsRepaint()
	})
	c.effects = append(c.effects, dispose)
	
	// Watch for layout changes
	dispose2 := reactive.Watch(func() {
		_ = c.width.Get()
		_ = c.height.Get()
		_ = c.padding.Get()
		_ = c.margin.Get()
		_ = c.alignment.Get()
		c.MarkNeedsLayout()
	})
	c.effects = append(c.effects, dispose2)
}

// SetChild sets the container's child widget
func (c *Container) SetChild(child WidgetImpl) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Remove old child
	if c.child != nil {
		if setter, ok := c.child.(interface{ SetParent(WidgetImpl) }); ok {
			setter.SetParent(nil)
		}
	}
	
	// Set new child
	c.child = child
	if child != nil {
		if setter, ok := child.(interface{ SetParent(WidgetImpl) }); ok {
			setter.SetParent(c)
		}
		// Update children slice for base widget compatibility
		c.children = []WidgetImpl{child}
	} else {
		c.children = []WidgetImpl{}
	}
	
	c.MarkNeedsLayout()
	c.MarkNeedsRepaint()
}

// GetChild returns the container's child
func (c *Container) GetChild() WidgetImpl {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.child
}

// SetWidth sets the container width
func (c *Container) SetWidth(width float64) {
	c.width.Set(width)
}

// SetHeight sets the container height
func (c *Container) SetHeight(height float64) {
	c.height.Set(height)
}

// SetColor sets the background color
func (c *Container) SetColor(color core.Color) {
	c.color.Set(color)
}

// SetPadding sets the padding
func (c *Container) SetPadding(padding EdgeInsets) {
	c.padding.Set(padding)
}

// SetMargin sets the margin
func (c *Container) SetMargin(margin EdgeInsets) {
	c.margin.Set(margin)
}

// SetAlignment sets the child alignment
func (c *Container) SetAlignment(alignment Alignment) {
	c.alignment.Set(alignment)
}

// SetBorder sets border properties
func (c *Container) SetBorder(color core.Color, width float64, radius float64) {
	reactive.Batch(func() {
		c.borderColor.Set(color)
		c.borderWidth.Set(width)
		c.borderRadius.Set(radius)
	})
}

// SetBoxShadow sets the box shadow
func (c *Container) SetBoxShadow(shadow *BoxShadow) {
	c.boxShadow.Set(shadow)
}

// Build creates the render object
func (c *Container) Build(ctx context.Context) RenderObject {
	return &RenderDecoratedBox{
		Decoration: BoxDecoration{
			Color:        c.color.Get(),
			BorderRadius: c.borderRadius.Get(),
			BorderColor:  c.borderColor.Get(),
			BorderWidth:  c.borderWidth.Get(),
		},
	}
}

// Layout calculates container dimensions
func (c *Container) Layout(constraints core.Constraints) (width, height float64) {
	margin := c.margin.Get()
	padding := c.padding.Get()
	
	// Reduce available space by margin
	innerConstraints := core.Constraints{
		MinWidth:  constraints.MinWidth - margin.Left - margin.Right,
		MaxWidth:  constraints.MaxWidth - margin.Left - margin.Right,
		MinHeight: constraints.MinHeight - margin.Top - margin.Bottom,
		MaxHeight: constraints.MaxHeight - margin.Top - margin.Bottom,
	}
	
	// Check for explicit dimensions
	width = c.width.Get()
	height = c.height.Get()
	
	if width > 0 {
		innerConstraints.MinWidth = width
		innerConstraints.MaxWidth = width
	}
	if height > 0 {
		innerConstraints.MinHeight = height
		innerConstraints.MaxHeight = height
	}
	
	// Layout child if present
	var childWidth, childHeight float64
	if c.child != nil {
		// Reduce space for padding
		childConstraints := core.Constraints{
			MinWidth:  innerConstraints.MinWidth - padding.Left - padding.Right,
			MaxWidth:  innerConstraints.MaxWidth - padding.Left - padding.Right,
			MinHeight: innerConstraints.MinHeight - padding.Top - padding.Bottom,
			MaxHeight: innerConstraints.MaxHeight - padding.Top - padding.Bottom,
		}
		
		// Ensure constraints are valid
		if childConstraints.MinWidth < 0 {
			childConstraints.MinWidth = 0
		}
		if childConstraints.MinHeight < 0 {
			childConstraints.MinHeight = 0
		}
		
		if coreWidget, ok := c.child.(core.Widget); ok {
			childWidth, childHeight = coreWidget.Layout(childConstraints)
		}
	}
	
	// Calculate final size
	finalWidth := childWidth + padding.Left + padding.Right
	finalHeight := childHeight + padding.Top + padding.Bottom
	
	// Apply explicit dimensions if set
	if width > 0 {
		finalWidth = width
	}
	if height > 0 {
		finalHeight = height
	}
	
	// Apply constraints
	if finalWidth < innerConstraints.MinWidth {
		finalWidth = innerConstraints.MinWidth
	}
	if finalWidth > innerConstraints.MaxWidth {
		finalWidth = innerConstraints.MaxWidth
	}
	
	if finalHeight < innerConstraints.MinHeight {
		finalHeight = innerConstraints.MinHeight
	}
	if finalHeight > innerConstraints.MaxHeight {
		finalHeight = innerConstraints.MaxHeight
	}
	
	width = finalWidth + margin.Left + margin.Right
	height = finalHeight + margin.Top + margin.Bottom
	
	c.cachedSize = Size{Width: width, Height: height}
	c.needsLayout.Set(false)
	
	return width, height
}

// Paint renders the container
func (c *Container) Paint(context core.PaintContext) {
	if !c.IsVisible() {
		return
	}
	
	margin := c.margin.Get()
	
	// Calculate content bounds (excluding margin)
	contentBounds := core.Bounds{
		X:      margin.Left,
		Y:      margin.Top,
		Width:  c.cachedSize.Width - margin.Left - margin.Right,
		Height: c.cachedSize.Height - margin.Top - margin.Bottom,
	}
	
	// Draw box shadow if present
	if shadow := c.boxShadow.Get(); shadow != nil {
		shadowBounds := core.Bounds{
			X:      contentBounds.X + shadow.Offset.X,
			Y:      contentBounds.Y + shadow.Offset.Y,
			Width:  contentBounds.Width,
			Height: contentBounds.Height,
		}
		
		alpha := 0.0
		if shadow.Color.A != 0 {
			alpha = float64(shadow.Color.A) / 255.0
		}
		context.DrawRect(shadowBounds, core.Paint{
			Color: shadow.Color,
			Alpha: alpha,
		})
	}
	
	// Draw background
	if bgColor := c.color.Get(); bgColor.A > 0 {
		context.DrawRect(contentBounds, core.Paint{
			Color: bgColor,
			Alpha: float64(bgColor.A)/255.0,
		})
	}
	
	// Draw border
	if borderWidth := c.borderWidth.Get(); borderWidth > 0 {
		context.DrawRect(contentBounds, core.Paint{
			Color: c.borderColor.Get(),
			Alpha: 1.0,
		})
	}
	
	// Paint child
	if c.child != nil {
		if coreWidget, ok := c.child.(core.Widget); ok {
			// For simplicity, just paint the child at the current context
			// In a full implementation, we'd need to transform the context
			coreWidget.Paint(context)
		}
	}
	
	c.needsRepaint.Set(false)
}

// HandleEvent processes input events
func (c *Container) HandleEvent(event core.Event) bool {
	// Propagate to child if exists
	if c.child != nil {
		if coreWidget, ok := c.child.(core.Widget); ok {
			return coreWidget.HandleEvent(event)
		}
	}
	return false
}

// Dispose cleans up the container
func (c *Container) Dispose() {
	c.width.Dispose()
	c.height.Dispose()
	c.color.Dispose()
	c.padding.Dispose()
	c.margin.Dispose()
	c.alignment.Dispose()
	c.borderRadius.Dispose()
	c.borderColor.Dispose()
	c.borderWidth.Dispose()
	c.boxShadow.Dispose()
	
	if c.child != nil {
		c.child.Dispose()
	}
	
	c.BaseWidget.Dispose()
}