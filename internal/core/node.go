package core

import (
	"runtime"
	"sync/atomic"
	"weak"
)

// NodeID is a unique identifier for nodes
type NodeID string

// DirtyFlags represents different types of changes that require updates
type DirtyFlags uint32

const (
	CleanFlag DirtyFlags = 0
	LayoutDirty DirtyFlags = 1 << iota
	PaintDirty
	ChildrenDirty
	PropertiesDirty
	TransformDirty
)

// Transform represents 2D transformation matrix
type Transform struct {
	Matrix [6]float64 // 2x3 affine transformation matrix
}

// Bounds represents the bounding box of a node
type Bounds struct {
	X, Y, Width, Height float64
}

// Offset represents a 2D position
type Offset struct {
	X, Y float64
}

// Constraints represents layout constraints
type Constraints struct {
	MinWidth, MaxWidth   float64
	MinHeight, MaxHeight float64
}

// Node represents a single element in the UI tree
type Node struct {
	// Identity
	ID       NodeID
	Type     string
	ZIndex   int
	
	// Tree structure with weak parent reference (Go 1.24)
	Parent   *weak.Pointer[Node]
	Children []*Node
	
	// Widget reference
	Widget   Widget
	
	// Layout properties
	Bounds           Bounds
	Transform        Transform
	CachedPosition   Offset
	ResolvedConstraints Constraints
	
	// State management
	isDirty     atomic.Bool
	dirtyFlags  atomic.Uint32
	version     atomic.Uint64
	
	// Performance optimizations
	intrinsicWidth  float64
	intrinsicHeight float64
	
	// Weak cache for computed values (Go 1.24)
	weakCache *weak.Pointer[ComputedValues]
}

// ComputedValues holds cached computed values for a node
type ComputedValues struct {
	Layout     LayoutData
	Paint      PaintData
	HitRegion  *HitRegion
}

// LayoutData contains computed layout information
type LayoutData struct {
	Position Offset
	Size     struct{ Width, Height float64 }
	Padding  struct{ Top, Right, Bottom, Left float64 }
	Margin   struct{ Top, Right, Bottom, Left float64 }
}

// PaintData contains rendering information
type PaintData struct {
	Opacity    float64
	Visibility bool
	ClipPath   []Offset
}

// HitRegion defines the area that responds to input
type HitRegion struct {
	Path   []Offset
	Bounds Bounds
}

// cleanupData holds cleanup data for a node
type cleanupData struct {
	widget Widget
}

// cleanup performs cleanup when node is garbage collected
func cleanup(data *cleanupData) {
	if data != nil && data.widget != nil {
		// Notify widget of disposal if it implements Disposable
		if d, ok := data.widget.(Disposable); ok {
			d.Dispose()
		}
	}
}

// NewNode creates a new node with automatic cleanup (Go 1.24)
func NewNode(id string, widget Widget) *Node {
	node := &Node{
		ID:       NodeID(id),
		Widget:   widget,
		Children: make([]*Node, 0, 4), // Pre-allocate for common case
	}
	
	// Go 1.24: Use runtime.AddCleanup for automatic cleanup
	// The cleanup data must be different from the node itself
	if widget != nil {
		cd := &cleanupData{widget: widget}
		runtime.AddCleanup(node, cleanup, cd)
	}
	
	return node
}

// GetParent safely retrieves the parent node using weak pointer
func (n *Node) GetParent() *Node {
	if n.Parent == nil {
		return nil
	}
	if ptr := n.Parent.Value(); ptr != nil {
		return ptr
	}
	return nil
}

// SetParent sets the parent with a weak reference
func (n *Node) SetParent(parent *Node) {
	if parent != nil {
		wp := weak.Make(parent)
		n.Parent = &wp
	} else {
		n.Parent = nil
	}
}

// AddChild adds a child node
func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
	child.SetParent(n)
	n.markDirty(ChildrenDirty)
}

// RemoveChild removes a child node
func (n *Node) RemoveChild(child *Node) bool {
	for i, c := range n.Children {
		if c == child {
			// Remove without preserving order for performance
			n.Children[i] = n.Children[len(n.Children)-1]
			n.Children = n.Children[:len(n.Children)-1]
			child.SetParent(nil)
			n.markDirty(ChildrenDirty)
			return true
		}
	}
	return false
}

// MarkDirty marks the node as needing update
func (n *Node) MarkDirty(flags DirtyFlags) {
	n.markDirty(flags)
}

// markDirty internal implementation
func (n *Node) markDirty(flags DirtyFlags) {
	// Update dirty flags atomically
	oldFlags := DirtyFlags(n.dirtyFlags.Load())
	newFlags := oldFlags | flags
	
	if oldFlags != newFlags {
		n.dirtyFlags.Store(uint32(newFlags))
		n.isDirty.Store(true)
		n.version.Add(1)
		
		// Propagate to parent if needed
		if parent := n.GetParent(); parent != nil {
			parent.markDirty(ChildrenDirty)
		}
	}
}

// IsDirty checks if the node needs update
func (n *Node) IsDirty() bool {
	return n.isDirty.Load()
}

// GetDirtyFlags returns current dirty flags
func (n *Node) GetDirtyFlags() DirtyFlags {
	return DirtyFlags(n.dirtyFlags.Load())
}

// ClearDirty clears dirty flags
func (n *Node) ClearDirty() {
	n.dirtyFlags.Store(0)
	n.isDirty.Store(false)
}

// GetVersion returns the current version number
func (n *Node) GetVersion() uint64 {
	return n.version.Load()
}

// GetCachedValues retrieves cached computed values if available
func (n *Node) GetCachedValues() *ComputedValues {
	if n.weakCache == nil {
		return nil
	}
	if ptr := n.weakCache.Value(); ptr != nil {
		return ptr
	}
	return nil
}

// SetCachedValues stores computed values in weak cache
func (n *Node) SetCachedValues(values *ComputedValues) {
	wc := weak.Make(values)
	n.weakCache = &wc
}

// Widget interface that all UI widgets must implement
type Widget interface {
	// Layout calculates the widget's size given constraints
	Layout(constraints Constraints) (width, height float64)
	
	// Paint renders the widget
	Paint(context PaintContext)
	
	// HandleEvent processes input events
	HandleEvent(event Event) bool
	
	// GetIntrinsicWidth returns the widget's natural width
	GetIntrinsicWidth(height float64) float64
	
	// GetIntrinsicHeight returns the widget's natural height
	GetIntrinsicHeight(width float64) float64
}

// Disposable interface for widgets that need cleanup
type Disposable interface {
	Dispose()
}

// PaintContext provides rendering context
type PaintContext interface {
	DrawRect(bounds Bounds, paint Paint)
	DrawText(text string, offset Offset, paint Paint)
	DrawPath(path []Offset, paint Paint)
	PushTransform(transform Transform)
	PopTransform()
	PushClip(bounds Bounds)
	PopClip()
}

// Paint defines rendering properties
type Paint struct {
	Color     Color
	Alpha     float64
	BlendMode string
}

// Color represents RGBA color
type Color struct {
	R, G, B, A uint8
}

// Event represents an input event
type Event interface {
	Type() string
	Timestamp() int64
}