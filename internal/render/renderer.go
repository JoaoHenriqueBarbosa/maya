//go:build wasm
// +build wasm

package render

import (
	"context"
	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/widgets"
)

// PaintCommand represents a drawing command
type PaintCommand struct {
	ID         string // Unique identifier for selective updates
	Type       PaintType
	Bounds     core.Bounds
	Text       string
	Color      core.Color
	Background core.Color
	Border     *BorderStyle
	Shadow     *ShadowStyle
	FontSize   float64
	OnClick    func()
}

type PaintType int

const (
	PaintRect PaintType = iota
	PaintText
	PaintButton
	PaintContainer
	UpdateText // Selective update for text content only
)

type BorderStyle struct {
	Color  core.Color
	Width  float64
	Radius float64
}

type ShadowStyle struct {
	Color      core.Color
	OffsetX    float64
	OffsetY    float64
	BlurRadius float64
}

// Renderer is the abstraction for different rendering backends
type Renderer interface {
	// Initialize the renderer with a container
	Init(container interface{}) error
	
	// Clear the rendering surface
	Clear()
	
	// Begin a new frame
	BeginFrame()
	
	// Paint a command
	Paint(cmd PaintCommand)
	
	// End the current frame
	EndFrame()
	
	// Apply updates - renderer decides if selective or full redraw
	// Returns true if handled, false if needs full redraw
	ApplyUpdates(updates []PaintCommand, allCommands []PaintCommand) bool
	
	// Handle resize
	Resize(width, height float64)
	
	// Get renderer name
	Name() string
}

// RenderContext holds rendering state
type RenderContext struct {
	Renderer    Renderer
	Commands    []PaintCommand
	Transform   Transform
	ClipBounds  *core.Bounds
}

type Transform struct {
	TranslateX float64
	TranslateY float64
	ScaleX     float64
	ScaleY     float64
}

// ConvertNodeToCommands converts a node tree to paint commands
func ConvertNodeToCommands(node *core.Node, offsetX, offsetY float64) []PaintCommand {
	if node.Widget == nil {
		return nil
	}
	
	var commands []PaintCommand
	
	// Calculate absolute position
	absX := node.Bounds.X + offsetX
	absY := node.Bounds.Y + offsetY
	
	// Create command based on widget type
	cmd := PaintCommand{
		ID:     string(node.ID), // Convert NodeID to string
		Bounds: core.Bounds{
			X:      absX,
			Y:      absY,
			Width:  node.Bounds.Width,
			Height: node.Bounds.Height,
		},
	}
	
	// Set command type and properties based on widget
	switch w := node.Widget.(type) {
	case *widgets.Text:
		cmd.Type = PaintText
		cmd.Text = w.GetText()
		style := w.GetStyle()
		cmd.FontSize = style.FontSize
		cmd.Color = style.Color
		
	case *widgets.Button:
		cmd.Type = PaintButton
		cmd.Text = w.GetLabel()
		cmd.OnClick = w.Click
		// Button gets default styling from renderer
		
	case *widgets.Container:
		cmd.Type = PaintContainer
		// Check for decoration via Build()
		if builder, ok := node.Widget.(interface{ Build(context.Context) widgets.RenderObject }); ok {
			if renderObj := builder.Build(context.Background()); renderObj != nil {
				if decorated, ok := renderObj.(*widgets.RenderDecoratedBox); ok {
					cmd.Background = decorated.Decoration.Color
					if decorated.Decoration.BorderWidth > 0 {
						cmd.Border = &BorderStyle{
							Color:  decorated.Decoration.BorderColor,
							Width:  decorated.Decoration.BorderWidth,
							Radius: decorated.Decoration.BorderRadius,
						}
					}
					if decorated.Decoration.BoxShadow != nil {
						shadow := decorated.Decoration.BoxShadow
						cmd.Shadow = &ShadowStyle{
							Color:      shadow.Color,
							OffsetX:    shadow.Offset.X,
							OffsetY:    shadow.Offset.Y,
							BlurRadius: shadow.BlurRadius,
						}
					}
				}
			}
		}
		
	default:
		cmd.Type = PaintRect // Generic rectangle
	}
	
	commands = append(commands, cmd)
	
	// Recursively add children
	for _, child := range node.Children {
		childCommands := ConvertNodeToCommands(child, absX, absY)
		commands = append(commands, childCommands...)
	}
	
	return commands
}