//go:build wasm
// +build wasm

package render

import (
	"fmt"
	"syscall/js"
	"github.com/maya-framework/maya/internal/core"
)

// CanvasRenderer renders to HTML5 Canvas
type CanvasRenderer struct {
	canvas    js.Value
	ctx       js.Value
	width     float64
	height    float64
	clickHandlers []ClickHandler
	lastCommands []PaintCommand // Store commands for redraw
}

type ClickHandler struct {
	Bounds  core.Bounds
	Handler func()
}

// NewCanvasRenderer creates a new canvas renderer
func NewCanvasRenderer() *CanvasRenderer {
	return &CanvasRenderer{
		clickHandlers: make([]ClickHandler, 0),
	}
}

func (r *CanvasRenderer) Init(container interface{}) error {
	doc := js.Global().Get("document")
	
	// Create canvas element
	r.canvas = doc.Call("createElement", "canvas")
	r.canvas.Set("width", 800)
	r.canvas.Set("height", 600)
	r.width = 800
	r.height = 600
	
	// Get 2D context
	r.ctx = r.canvas.Call("getContext", "2d")
	
	// Add to container
	if cont, ok := container.(js.Value); ok {
		cont.Call("appendChild", r.canvas)
	} else {
		return fmt.Errorf("CanvasRenderer requires js.Value container")
	}
	
	// Setup click handler
	r.canvas.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			event := args[0]
			rect := r.canvas.Call("getBoundingClientRect")
			x := event.Get("clientX").Float() - rect.Get("left").Float()
			y := event.Get("clientY").Float() - rect.Get("top").Float()
			
			// Check all click handlers
			for _, handler := range r.clickHandlers {
				if x >= handler.Bounds.X && x <= handler.Bounds.X+handler.Bounds.Width &&
				   y >= handler.Bounds.Y && y <= handler.Bounds.Y+handler.Bounds.Height {
					if handler.Handler != nil {
						handler.Handler()
					}
					break
				}
			}
		}
		return nil
	}))
	
	return nil
}

func (r *CanvasRenderer) Clear() {
	r.ctx.Call("clearRect", 0, 0, r.width, r.height)
	r.clickHandlers = r.clickHandlers[:0]
}

func (r *CanvasRenderer) BeginFrame() {
	r.Clear()
	
	// Set default font
	r.ctx.Set("font", "16px system-ui, -apple-system, sans-serif")
	r.ctx.Set("textBaseline", "top")
}

func (r *CanvasRenderer) Paint(cmd PaintCommand) {
	// Save context state
	r.ctx.Call("save")
	
	// Draw background if specified
	if cmd.Background.A > 0 || cmd.Type == PaintContainer {
		if cmd.Background.A > 0 {
			r.ctx.Set("fillStyle", formatCanvasColor(cmd.Background))
		} else {
			r.ctx.Set("fillStyle", "transparent")
		}
		
		if cmd.Border != nil && cmd.Border.Radius > 0 {
			r.drawRoundedRect(cmd.Bounds, cmd.Border.Radius)
			r.ctx.Call("fill")
		} else {
			r.ctx.Call("fillRect", cmd.Bounds.X, cmd.Bounds.Y, cmd.Bounds.Width, cmd.Bounds.Height)
		}
	}
	
	// Draw border
	if cmd.Border != nil && cmd.Border.Width > 0 {
		r.ctx.Set("strokeStyle", formatCanvasColor(cmd.Border.Color))
		r.ctx.Set("lineWidth", cmd.Border.Width)
		
		if cmd.Border.Radius > 0 {
			r.drawRoundedRect(cmd.Bounds, cmd.Border.Radius)
			r.ctx.Call("stroke")
		} else {
			r.ctx.Call("strokeRect", cmd.Bounds.X, cmd.Bounds.Y, cmd.Bounds.Width, cmd.Bounds.Height)
		}
	}
	
	// Draw shadow
	if cmd.Shadow != nil {
		r.ctx.Set("shadowColor", formatCanvasColor(cmd.Shadow.Color))
		r.ctx.Set("shadowOffsetX", cmd.Shadow.OffsetX)
		r.ctx.Set("shadowOffsetY", cmd.Shadow.OffsetY)
		r.ctx.Set("shadowBlur", cmd.Shadow.BlurRadius)
	}
	
	// Draw content based on type
	switch cmd.Type {
	case PaintText:
		if cmd.FontSize > 0 {
			r.ctx.Set("font", fmt.Sprintf("%fpx system-ui, -apple-system, sans-serif", cmd.FontSize))
		}
		r.ctx.Set("fillStyle", formatCanvasColor(cmd.Color))
		r.ctx.Call("fillText", cmd.Text, cmd.Bounds.X, cmd.Bounds.Y)
		
	case PaintButton:
		// Draw button background
		r.ctx.Set("fillStyle", "#007acc")
		r.drawRoundedRect(cmd.Bounds, 5)
		r.ctx.Call("fill")
		
		// Draw button text
		r.ctx.Set("fillStyle", "white")
		r.ctx.Set("font", "14px system-ui, -apple-system, sans-serif")
		r.ctx.Set("textAlign", "center")
		r.ctx.Set("textBaseline", "middle")
		r.ctx.Call("fillText", cmd.Text, 
			cmd.Bounds.X + cmd.Bounds.Width/2,
			cmd.Bounds.Y + cmd.Bounds.Height/2)
		r.ctx.Set("textAlign", "start")
		r.ctx.Set("textBaseline", "top")
		
		// Register click handler
		if cmd.OnClick != nil {
			r.clickHandlers = append(r.clickHandlers, ClickHandler{
				Bounds:  cmd.Bounds,
				Handler: cmd.OnClick,
			})
		}
	}
	
	// Restore context state
	r.ctx.Call("restore")
}

func (r *CanvasRenderer) drawRoundedRect(bounds core.Bounds, radius float64) {
	x := bounds.X
	y := bounds.Y
	width := bounds.Width
	height := bounds.Height
	
	r.ctx.Call("beginPath")
	r.ctx.Call("moveTo", x + radius, y)
	r.ctx.Call("lineTo", x + width - radius, y)
	r.ctx.Call("quadraticCurveTo", x + width, y, x + width, y + radius)
	r.ctx.Call("lineTo", x + width, y + height - radius)
	r.ctx.Call("quadraticCurveTo", x + width, y + height, x + width - radius, y + height)
	r.ctx.Call("lineTo", x + radius, y + height)
	r.ctx.Call("quadraticCurveTo", x, y + height, x, y + height - radius)
	r.ctx.Call("lineTo", x, y + radius)
	r.ctx.Call("quadraticCurveTo", x, y, x + radius, y)
	r.ctx.Call("closePath")
}

func (r *CanvasRenderer) EndFrame() {
	// Nothing to do for canvas
}

func (r *CanvasRenderer) ApplyUpdates(updates []PaintCommand, allCommands []PaintCommand) bool {
	// Canvas can't do selective updates efficiently, request full redraw
	println("[CANVAS-UPDATE] Updates detected, requesting full redraw")
	return false // Tell pipeline we need full redraw
}

func (r *CanvasRenderer) Resize(width, height float64) {
	r.width = width
	r.height = height
	r.canvas.Set("width", width)
	r.canvas.Set("height", height)
}

func (r *CanvasRenderer) Name() string {
	return "Canvas"
}

func formatCanvasColor(c core.Color) string {
	return fmt.Sprintf("rgba(%d,%d,%d,%f)", 
		c.R, c.G, c.B, float64(c.A)/255.0)
}