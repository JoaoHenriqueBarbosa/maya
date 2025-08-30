package widgets

import (
	"context"
	"github.com/maya-framework/maya/internal/core"
)

// WidgetImpl extends the core.Widget interface with additional functionality
// specific to the widgets package implementation
type WidgetImpl interface {
	core.Widget  // Embed the core Widget interface
	
	// Core identification
	ID() string
	Type() string
	
	// Tree structure
	Parent() WidgetImpl
	Children() []WidgetImpl
	AddChild(child WidgetImpl)
	RemoveChild(child WidgetImpl) bool
	
	// Lifecycle
	Build(ctx context.Context) RenderObject
	Init()
	Dispose()
	
	// State & Props
	GetProps() Props
	SetProps(props Props)
	
	// Additional methods
	NeedsRepaint() bool
	MarkNeedsRepaint()
}

// Props represents widget properties
type Props map[string]interface{}

// Size represents widget dimensions
type Size struct {
	Width  float64
	Height float64
}

// Rect represents a rectangle
type Rect struct {
	core.Offset
	Size
}

// Canvas provides low-level drawing operations
type Canvas interface {
	DrawRect(rect Rect, paint core.Paint)
	DrawText(text string, offset core.Offset, style TextStyle)
	DrawImage(image Image, offset core.Offset)
	ClipRect(rect Rect)
	Save()
	Restore()
}

// Common colors - using core.Color type
var (
	ColorBlack       = core.Color{R: 0, G: 0, B: 0, A: 255}
	ColorWhite       = core.Color{R: 255, G: 255, B: 255, A: 255}
	ColorRed         = core.Color{R: 255, G: 0, B: 0, A: 255}
	ColorGreen       = core.Color{R: 0, G: 255, B: 0, A: 255}
	ColorBlue        = core.Color{R: 0, G: 0, B: 255, A: 255}
	ColorTransparent = core.Color{R: 0, G: 0, B: 0, A: 0}
)

// TextStyle describes text rendering
type TextStyle struct {
	FontFamily  string
	FontSize    float64
	FontWeight  FontWeight
	Color       core.Color
	LineHeight  float64
}

// FontWeight enumeration
type FontWeight int

const (
	FontWeightNormal FontWeight = 400
	FontWeightBold   FontWeight = 700
)

// Image represents a bitmap image
type Image interface {
	Width() int
	Height() int
	Data() []byte
}

// EventType enumeration
type EventType int

const (
	EventTypePointerDown EventType = iota
	EventTypePointerMove
	EventTypePointerUp
	EventTypePointerCancel
	EventTypeKeyDown
	EventTypeKeyUp
	EventTypeScroll
)

// PointerEvent represents mouse/touch events
type PointerEvent struct {
	EventType EventType
	Position  core.Offset
	Button    MouseButton
	Pressure  float64
	TimeStamp int64
}

func (e *PointerEvent) Type() EventType {
	return e.EventType
}

func (e *PointerEvent) Timestamp() int64 {
	return e.TimeStamp
}

// MouseButton enumeration
type MouseButton int

const (
	MouseButtonPrimary MouseButton = iota
	MouseButtonSecondary
	MouseButtonMiddle
)

// KeyEvent represents keyboard events
type KeyEvent struct {
	EventType   EventType
	Key         string
	Code        string
	CtrlKey     bool
	ShiftKey    bool
	AltKey      bool
	MetaKey     bool
	TimeStamp   int64
}

func (e *KeyEvent) Type() EventType {
	return e.EventType
}

func (e *KeyEvent) Timestamp() int64 {
	return e.TimeStamp
}