//go:build wasm
// +build wasm

package render

import (
	"fmt"
	"syscall/js"
	"github.com/maya-framework/maya/internal/core"
)

// DOMRenderer renders to HTML DOM
type DOMRenderer struct {
	container js.Value
	elements  map[string]js.Value // Track created elements by ID
	document  js.Value
	frameStarted bool
}

// NewDOMRenderer creates a new DOM renderer
func NewDOMRenderer() *DOMRenderer {
	return &DOMRenderer{
		elements:  make(map[string]js.Value),
		document:  js.Global().Get("document"),
	}
}

func (r *DOMRenderer) Init(container interface{}) error {
	if cont, ok := container.(js.Value); ok {
		r.container = cont
		return nil
	}
	return fmt.Errorf("DOMRenderer requires js.Value container")
}

func (r *DOMRenderer) Clear() {
	r.container.Set("innerHTML", "")
	r.elements = make(map[string]js.Value)
}

func (r *DOMRenderer) BeginFrame() {
	// Clear container for full re-render
	r.container.Set("innerHTML", "")
	r.elements = make(map[string]js.Value)
}

func (r *DOMRenderer) Paint(cmd PaintCommand) {
	var elem js.Value
	
	switch cmd.Type {
	case PaintText:
		elem = r.document.Call("createElement", "span")
		elem.Set("textContent", cmd.Text)
		
	case PaintButton:
		elem = r.document.Call("createElement", "button")
		elem.Set("textContent", cmd.Text)
		
		if cmd.OnClick != nil {
			// Create direct event handler
			elem.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				cmd.OnClick()
				return nil
			}))
		}
		
	case PaintContainer, PaintRect:
		elem = r.document.Call("createElement", "div")
		
	default:
		elem = r.document.Call("createElement", "div")
	}
	
	// Store element for future updates
	if cmd.ID != "" {
		r.elements[cmd.ID] = elem
	}
	
	// Apply positioning
	style := elem.Get("style")
	style.Set("position", "absolute")
	style.Set("left", fmt.Sprintf("%fpx", cmd.Bounds.X))
	style.Set("top", fmt.Sprintf("%fpx", cmd.Bounds.Y))
	
	if cmd.Bounds.Width > 0 {
		style.Set("width", fmt.Sprintf("%fpx", cmd.Bounds.Width))
	}
	if cmd.Bounds.Height > 0 {
		style.Set("height", fmt.Sprintf("%fpx", cmd.Bounds.Height))
	}
	
	// Apply styling
	if cmd.Background.A > 0 {
		style.Set("backgroundColor", formatColor(cmd.Background))
	}
	
	if cmd.Color.A > 0 {
		style.Set("color", formatColor(cmd.Color))
	}
	
	if cmd.FontSize > 0 {
		style.Set("fontSize", fmt.Sprintf("%fpx", cmd.FontSize))
	}
	
	// Apply border
	if cmd.Border != nil {
		style.Set("borderWidth", fmt.Sprintf("%fpx", cmd.Border.Width))
		style.Set("borderStyle", "solid")
		style.Set("borderColor", formatColor(cmd.Border.Color))
		if cmd.Border.Radius > 0 {
			style.Set("borderRadius", fmt.Sprintf("%fpx", cmd.Border.Radius))
		}
	}
	
	// Apply shadow
	if cmd.Shadow != nil {
		style.Set("boxShadow", fmt.Sprintf("%fpx %fpx %fpx %s",
			cmd.Shadow.OffsetX,
			cmd.Shadow.OffsetY,
			cmd.Shadow.BlurRadius,
			formatColor(cmd.Shadow.Color)))
	}
	
	// Default button styling
	if cmd.Type == PaintButton {
		style.Set("backgroundColor", "#007acc")
		style.Set("color", "white")
		style.Set("border", "none")
		style.Set("borderRadius", "5px")
		style.Set("cursor", "pointer")
		style.Set("fontSize", "14px")
		style.Set("fontFamily", "system-ui, -apple-system, sans-serif")
		style.Set("textAlign", "center")
	}
	
	r.container.Call("appendChild", elem)
}

func (r *DOMRenderer) EndFrame() {
	// Nothing to do for DOM
}

func (r *DOMRenderer) ApplyUpdates(updates []PaintCommand, allCommands []PaintCommand) bool {
	// DOM can handle selective updates
	for _, cmd := range updates {
		if elem, exists := r.elements[cmd.ID]; exists {
			if cmd.Type == UpdateText {
				println("[DOM-UPDATE] Updating text for", cmd.ID, "to", cmd.Text)
				elem.Set("textContent", cmd.Text)
			}
			// Add more update types as needed
		} else {
			println("[DOM-UPDATE] Element not found for ID:", cmd.ID)
		}
	}
	return true // Handled successfully
}

func (r *DOMRenderer) Resize(width, height float64) {
	// DOM handles this automatically
}

func (r *DOMRenderer) Name() string {
	return "DOM"
}


func formatColor(c core.Color) string {
	return fmt.Sprintf("rgba(%d,%d,%d,%f)", 
		c.R, c.G, c.B, float64(c.A)/255.0)
}