//go:build wasm
// +build wasm

package maya

import (
	"context"
	"syscall/js"

	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/reactive"
	"github.com/maya-framework/maya/internal/render"
	"github.com/maya-framework/maya/internal/widgets"
)

// Global app instance for reactive updates
var globalApp *App

// Component is a function that returns a widget
type Component func() widgets.WidgetImpl

// App represents a Maya application - SIMPLIFIED API
type App struct {
	tree      *core.Tree              // Use REAL tree
	pipeline  *render.Pipeline        // Use REAL pipeline
	batcher   *reactive.UpdateBatcher // Use REAL batcher
	container js.Value
	ctx       context.Context
	cancel    context.CancelFunc
	root      Component               // Root component for re-rendering
}

// New creates a new Maya application - SIMPLE API
func New(root Component) *App {
	ctx, cancel := context.WithCancel(context.Background())
	
	app := &App{
		tree:    core.NewTree(),
		batcher: reactive.NewUpdateBatcher(),
		ctx:     ctx,
		cancel:  cancel,
		root:    root,
	}
	
	// Build widget and convert to tree
	rootWidget := root()
	rootNode := app.widgetToNode(rootWidget)
	app.tree.SetRoot(rootNode)
	
	return app
}

// Run starts the application
func (app *App) Run() {
	// Set global app for reactive updates
	globalApp = app
	
	// Wait for DOM
	waitForDOM(func() {
		app.setupViewport()
		
		// Get container
		container := js.Global().Get("document").Call("getElementById", "app")
		if container.IsNull() {
			body := js.Global().Get("document").Get("body")
			container = js.Global().Get("document").Call("createElement", "div")
			container.Set("id", "app")
			body.Call("appendChild", container)
		}
		
		app.container = container
		container.Set("innerHTML", "")
		
		// Create pipeline with our REAL components
		app.pipeline = render.NewPipeline(app.tree, container, &render.Theme{
			Primary:    "#007acc",
			Text:       "#333333",
			Background: "#ffffff",
			FontFamily: "system-ui, -apple-system, sans-serif",
		})
		
		// Initial render
		app.render()
		
		// Start reactive batching
		app.batcher.Start()
		app.setupReactiveLoop()
	})
	
	// Keep running
	select {}
}

// render executes the pipeline
func (app *App) render() {
	println("Re-rendering application...")
	
	// Rebuild the widget tree
	rootWidget := app.root()
	rootNode := app.widgetToNode(rootWidget)
	app.tree.SetRoot(rootNode)
	
	// Execute the render pipeline
	if err := app.pipeline.Execute(app.ctx); err != nil {
		println("Render error:", err.Error())
	}
}

// scheduleRender schedules a render using the batcher
func (app *App) scheduleRender() {
	println("Scheduling render via batcher...")
	app.batcher.Add(func() {
		app.render()
	})
}

// widgetToNode converts a widget to a core.Node recursively
func (app *App) widgetToNode(widget widgets.WidgetImpl) *core.Node {
	node := core.NewNode(widget.ID(), widget)
	
	// Handle different widget types
	switch w := widget.(type) {
	case *widgets.Container:
		if w.GetChild() != nil {
			child := app.widgetToNode(w.GetChild())
			node.AddChild(child)
		}
		
	case *widgets.Column:
		for _, child := range w.Children() {
			childNode := app.widgetToNode(child)
			node.AddChild(childNode)
		}
		
	case *widgets.Row:
		for _, child := range w.Children() {
			childNode := app.widgetToNode(child)
			node.AddChild(childNode)
		}
	}
	
	return node
}

// setupReactiveLoop sets up the reactive update loop
func (app *App) setupReactiveLoop() {
	// Use the REAL batcher for updates
	app.batcher.Add(func() {
		app.render()
	})
}

// setupViewport configures the viewport
func (app *App) setupViewport() {
	doc := js.Global().Get("document")
	body := doc.Get("body")
	
	bodyStyle := body.Get("style")
	bodyStyle.Set("margin", "0")
	bodyStyle.Set("padding", "0")
	bodyStyle.Set("fontFamily", "system-ui, -apple-system, sans-serif")
	bodyStyle.Set("background", "#ffffff")
	bodyStyle.Set("minHeight", "100vh")
	bodyStyle.Set("display", "flex")
	bodyStyle.Set("alignItems", "center")
	bodyStyle.Set("justifyContent", "center")
}

// waitForDOM is now in exports.go without js.FuncOf

// ============================================================
// PUBLIC API - Simple functions for users
// ============================================================

// Container creates a container widget
func Container(children ...widgets.WidgetImpl) widgets.WidgetImpl {
	if len(children) == 0 {
		return widgets.NewContainer("container")
	}
	
	// Use Column for multiple children
	if len(children) > 1 {
		return Column(children...)
	}
	
	c := widgets.NewContainer("container")
	c.SetChild(children[0])
	return c
}

// Column creates a vertical layout
func Column(children ...widgets.WidgetImpl) widgets.WidgetImpl {
	return widgets.NewColumn("column", children...)
}

// Row creates a horizontal layout
func Row(children ...widgets.WidgetImpl) widgets.WidgetImpl {
	return widgets.NewRow("row", children...)
}

// Text creates a text widget
func Text(text string) widgets.WidgetImpl {
	t := widgets.NewText("text-"+text, text)
	t.SetStyle(widgets.TextStyle{
		FontSize:   16,
		FontWeight: widgets.FontWeightNormal,
		Color:      widgets.ColorBlack,
		LineHeight: 1.5,
	})
	return t
}

// Title creates a title text widget
func Title(text string) widgets.WidgetImpl {
	t := widgets.NewText("title-"+text, text)
	t.SetStyle(widgets.TextStyle{
		FontSize:   24,
		FontWeight: widgets.FontWeightBold,
		Color:      widgets.ColorBlack,
		LineHeight: 1.2,
	})
	return t
}

// Button creates a button widget
func Button(text string, onClick func()) widgets.WidgetImpl {
	return widgets.NewButton("button-"+text, text, onClick)
}

// Signal creates a reactive signal - using REAL reactive system
func Signal[T comparable](initial T) *reactive.Signal[T] {
	return reactive.NewSignal(initial)
}

// TextSignal creates a reactive text widget
func TextSignal[T any](signal *reactive.Signal[T], format func(T) string) widgets.WidgetImpl {
	text := widgets.NewText("reactive-text", format(signal.Get()))
	
	// Create effect to update text when signal changes
	reactive.CreateEffect(func() {
		text.SetText(format(signal.Get()))
		// Trigger re-render if app is running
		if globalApp != nil {
			globalApp.scheduleRender()
		}
	})
	
	return text
}