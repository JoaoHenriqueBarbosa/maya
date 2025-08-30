package widgets

import (
	"context"
	"sync"
	"sync/atomic"
	
	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/reactive"
)

// BaseWidget provides default implementation for WidgetImpl interface
type BaseWidget struct {
	id       string
	widgetType string
	
	// Tree structure
	parent   WidgetImpl
	children []WidgetImpl
	mu       sync.RWMutex
	
	// Props with reactive signals
	props    *reactive.Signal[Props]
	
	// State
	needsRepaint *reactive.Signal[bool]
	needsLayout  *reactive.Signal[bool]
	
	// Lifecycle
	initialized atomic.Bool
	disposed    atomic.Bool
	
	// Layout cache
	cachedSize   Size
	cachedConstraints *core.Constraints
	
	// Effects for reactive updates
	effects []func()
}

// NewBaseWidget creates a new base widget
func NewBaseWidget(id string, widgetType string) *BaseWidget {
	w := &BaseWidget{
		id:         id,
		widgetType: widgetType,
		children:   make([]WidgetImpl, 0),
		props:      reactive.NewSignal(make(Props)),
		needsRepaint: reactive.NewSignal(false),
		needsLayout:  reactive.NewSignal(false),
		effects:    make([]func(), 0),
	}
	
	// Setup reactive effects
	w.setupEffects()
	
	return w
}

// setupEffects creates reactive effects for the widget
func (w *BaseWidget) setupEffects() {
	// Watch for prop changes that require repaint
	dispose := reactive.Watch(func() {
		props := w.props.Get()
		if _, ok := props["visible"]; ok {
			w.MarkNeedsRepaint()
		}
	})
	w.effects = append(w.effects, dispose)
}

// ID returns the widget's unique identifier
func (w *BaseWidget) ID() string {
	return w.id
}

// Type returns the widget type
func (w *BaseWidget) Type() string {
	return w.widgetType
}

// Parent returns the parent widget
func (w *BaseWidget) Parent() WidgetImpl {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.parent
}

// SetParent sets the parent widget
func (w *BaseWidget) SetParent(parent WidgetImpl) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.parent = parent
}

// Children returns child widgets
func (w *BaseWidget) Children() []WidgetImpl {
	w.mu.RLock()
	defer w.mu.RUnlock()
	
	// Return a copy to prevent external modification
	children := make([]WidgetImpl, len(w.children))
	copy(children, w.children)
	return children
}

// AddChild adds a child widget
func (w *BaseWidget) AddChild(child WidgetImpl) {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	w.children = append(w.children, child)
	// Set parent for any widget that has SetParent method
	if setter, ok := child.(interface{ SetParent(WidgetImpl) }); ok {
		setter.SetParent(w)
	}
	
	w.needsLayout.Set(true)
}

// RemoveChild removes a child widget
func (w *BaseWidget) RemoveChild(child WidgetImpl) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	for i, c := range w.children {
		if c == child {
			// Remove from slice
			w.children = append(w.children[:i], w.children[i+1:]...)
			
			// Clear parent reference
			if setter, ok := child.(interface{ SetParent(WidgetImpl) }); ok {
				setter.SetParent(nil)
			}
			
			w.needsLayout.Set(true)
			return true
		}
	}
	return false
}

// Build creates the render object for this widget
func (w *BaseWidget) Build(ctx context.Context) RenderObject {
	// Default implementation - subclasses should override
	return &RenderBox{
		Size: w.cachedSize,
	}
}

// Init initializes the widget
func (w *BaseWidget) Init() {
	if w.initialized.CompareAndSwap(false, true) {
		// Initialization logic
		w.needsLayout.Set(true)
		w.needsRepaint.Set(true)
	}
}

// Dispose cleans up the widget
func (w *BaseWidget) Dispose() {
	if w.disposed.CompareAndSwap(false, true) {
		// Cleanup effects
		for _, dispose := range w.effects {
			dispose()
		}
		w.effects = nil
		
		// Dispose children
		w.mu.Lock()
		for _, child := range w.children {
			child.Dispose()
		}
		w.children = nil
		w.mu.Unlock()
		
		// Dispose signals
		w.props.Dispose()
		w.needsRepaint.Dispose()
		w.needsLayout.Dispose()
	}
}

// GetProps returns widget properties
func (w *BaseWidget) GetProps() Props {
	return w.props.Get()
}

// SetProps updates widget properties
func (w *BaseWidget) SetProps(props Props) {
	reactive.Batch(func() {
		w.props.Set(props)
		w.needsLayout.Set(true)
		w.needsRepaint.Set(true)
	})
}

// UpdateProp updates a single property
func (w *BaseWidget) UpdateProp(key string, value interface{}) {
	props := w.props.Get()
	newProps := make(Props)
	for k, v := range props {
		newProps[k] = v
	}
	newProps[key] = value
	w.SetProps(newProps)
}

// Layout calculates widget size given constraints
func (w *BaseWidget) Layout(constraints core.Constraints) (width, height float64) {
	// Check cache
	if w.cachedConstraints != nil && 
	   w.cachedConstraints.MinWidth == constraints.MinWidth &&
	   w.cachedConstraints.MaxWidth == constraints.MaxWidth &&
	   w.cachedConstraints.MinHeight == constraints.MinHeight &&
	   w.cachedConstraints.MaxHeight == constraints.MaxHeight {
		return w.cachedSize.Width, w.cachedSize.Height
	}
	
	// Default implementation - take maximum available space
	width = constraints.MaxWidth
	height = constraints.MaxHeight
	
	// Cache results
	w.cachedConstraints = &constraints
	w.cachedSize = Size{Width: width, Height: height}
	w.needsLayout.Set(false)
	
	return width, height
}

// GetIntrinsicWidth returns the intrinsic width for a given height
func (w *BaseWidget) GetIntrinsicWidth(height float64) float64 {
	// Default implementation
	return 0
}

// GetIntrinsicHeight returns the intrinsic height for a given width
func (w *BaseWidget) GetIntrinsicHeight(width float64) float64 {
	// Default implementation
	return 0
}

// Paint renders the widget
func (w *BaseWidget) Paint(context core.PaintContext) {
	// Default implementation - subclasses should override
	// Paint background if specified
	props := w.GetProps()
	if bg, ok := props["backgroundColor"].(core.Color); ok {
		bounds := core.Bounds{
			X:      0,
			Y:      0,
			Width:  w.cachedSize.Width,
			Height: w.cachedSize.Height,
		}
		context.DrawRect(bounds, core.Paint{
			Color: bg,
			Alpha: 1.0,
		})
	}
	
	// Paint children
	w.mu.RLock()
	children := w.children
	w.mu.RUnlock()
	
	for _, child := range children {
		if coreWidget, ok := child.(core.Widget); ok {
			coreWidget.Paint(context)
		}
	}
	
	w.needsRepaint.Set(false)
}

// NeedsRepaint returns true if widget needs repainting
func (w *BaseWidget) NeedsRepaint() bool {
	return w.needsRepaint.Get()
}

// MarkNeedsRepaint marks widget as needing repaint
func (w *BaseWidget) MarkNeedsRepaint() {
	w.needsRepaint.Set(true)
	
	// Propagate to parent
	if parent := w.Parent(); parent != nil {
		parent.MarkNeedsRepaint()
	}
}

// MarkNeedsLayout marks widget as needing layout
func (w *BaseWidget) MarkNeedsLayout() {
	w.needsLayout.Set(true)
	
	// Propagate to parent
	if parent := w.Parent(); parent != nil {
		if base, ok := parent.(*BaseWidget); ok {
			base.MarkNeedsLayout()
		}
	}
}

// HandleEvent processes input events
func (w *BaseWidget) HandleEvent(event core.Event) bool {
	// Default implementation - propagate to children
	w.mu.RLock()
	children := w.children
	w.mu.RUnlock()
	
	// Process in reverse order (top to bottom)
	for i := len(children) - 1; i >= 0; i-- {
		if coreWidget, ok := children[i].(core.Widget); ok {
			if coreWidget.HandleEvent(event) {
				return true
			}
		}
	}
	
	return false
}

// IsVisible returns true if widget is visible
func (w *BaseWidget) IsVisible() bool {
	props := w.GetProps()
	if visible, ok := props["visible"].(bool); ok {
		return visible
	}
	return true // Default to visible
}

// GetOpacity returns widget opacity
func (w *BaseWidget) GetOpacity() float64 {
	props := w.GetProps()
	if opacity, ok := props["opacity"].(float64); ok {
		return opacity
	}
	return 1.0 // Default to fully opaque
}