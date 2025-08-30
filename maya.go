//go:build wasm
// +build wasm

package maya

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/logger"
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
	tree         *core.Tree              // Use REAL tree
	pipeline     *render.Pipeline        // Use REAL pipeline
	batcher      *reactive.UpdateBatcher // Use REAL batcher
	container    js.Value
	ctx          context.Context
	cancel       context.CancelFunc
	root         Component        // Root component for re-rendering
	renderEffect *reactive.Effect // Single effect for re-rendering
	widgetElements map[widgets.WidgetImpl]js.Value // Direct widget->DOM mapping for fine-grained updates
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
		widgetElements: make(map[widgets.WidgetImpl]js.Value),
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
		
		// Determine renderer type from window.MAYA_RENDERER
		rendererType := "dom" // default
		if rt := js.Global().Get("window").Get("MAYA_RENDERER"); !rt.IsUndefined() {
			rendererType = rt.String()
		}
		
		// Create appropriate renderer
		var renderer render.Renderer
		switch rendererType {
		case "canvas":
			renderer = render.NewCanvasRenderer()
		default:
			renderer = render.NewDOMRenderer()
		}
		
		// Initialize renderer with container
		if err := renderer.Init(container); err != nil {
			panic(fmt.Sprintf("Failed to initialize renderer: %v", err))
		}
		
		logger.Info(logger.TagMaya, "Using renderer: %s", renderer.Name())

		// Create pipeline with our REAL components and selected renderer
		app.pipeline = render.NewPipeline(app.tree, renderer, &render.Theme{
			Primary:    "#007acc",
			Text:       "#333333",
			Background: "#ffffff",
			FontFamily: "system-ui, -apple-system, sans-serif",
		})

		// Create single root effect for reactive updates
		app.setupReactiveEffect()

		// Start reactive batching
		app.batcher.Start()
	})

	// Keep running
	select {}
}

// updateWidget updates a specific widget in the DOM without recreating everything
func (app *App) updateWidget(widget widgets.WidgetImpl) {
	logger.Trace(logger.TagWidget, "Updating specific widget: %s", widget.ID())
	
	// Find the node in the tree
	var targetNode *core.Node
	for node := range app.tree.PreOrderDFS() {
		if node.Widget == widget {
			targetNode = node
			break
		}
	}
	
	if targetNode != nil {
		// Mark as dirty and run selective update
		targetNode.MarkDirty(core.LayoutDirty)
		
		// For now, re-run the full pipeline (can optimize later)
		if app.pipeline != nil {
			app.pipeline.Execute(context.Background())
		}
	}
}

// render executes the pipeline with current tree
func (app *App) render() {
	logger.Trace(logger.TagRender, "Starting pipeline execution...")
	if err := app.pipeline.Execute(app.ctx); err != nil {
		logger.Error(logger.TagRender, "Error: %s", err.Error())
	} else {
		logger.Trace(logger.TagRender, "Pipeline execution complete")
	}
}


// widgetToNode converts a widget to a core.Node recursively
func (app *App) widgetToNode(widget widgets.WidgetImpl) *core.Node {
	logger.Trace(logger.TagWidget, "Converting widget to node: %s", widget.ID())
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

// setupReactiveEffect creates the widget tree ONCE and lets widgets handle their own updates
func (app *App) setupReactiveEffect() {
	logger.Trace(logger.TagReactive, "Building widget tree ONCE...")
	
	// Build tree ONCE - widgets have their own reactive effects
	rootWidget := app.root()
	rootNode := app.widgetToNode(rootWidget)
	app.tree.SetRoot(rootNode)
	
	// Initial render
	app.render()
	
	// Widgets will update themselves via their own effects
	logger.Trace(logger.TagReactive, "Setup complete - using widget-level reactivity")
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

// StyledContainer creates a container with styling
func StyledContainer(style ContainerStyle, children ...widgets.WidgetImpl) widgets.WidgetImpl {
	var child widgets.WidgetImpl
	if len(children) == 0 {
		child = nil
	} else if len(children) == 1 {
		child = children[0]
	} else {
		child = Column(children...)
	}

	c := widgets.NewContainer("styled-container")
	if child != nil {
		c.SetChild(child)
	}
	
	// Apply styling
	if style.Background != nil {
		c.SetColor(*style.Background)
	}
	if style.BorderColor != nil && style.BorderWidth > 0 {
		c.SetBorder(*style.BorderColor, style.BorderWidth, style.BorderRadius)
	}
	if style.Padding != nil {
		c.SetPadding(*style.Padding)
	}
	if style.Shadow != nil {
		c.SetBoxShadow(style.Shadow)
	}
	
	return c
}

// ContainerStyle defines styling options for containers
type ContainerStyle struct {
	Background   *core.Color
	BorderColor  *core.Color
	BorderWidth  float64
	BorderRadius float64
	Padding      *widgets.EdgeInsets
	Shadow       *widgets.BoxShadow
}

// Predefined colors for convenience
var (
	ColorWhite   = &core.Color{R: 255, G: 255, B: 255, A: 255}
	ColorBlack   = &core.Color{R: 0, G: 0, B: 0, A: 255}
	ColorRed     = &core.Color{R: 255, G: 0, B: 0, A: 255}
	ColorGreen   = &core.Color{R: 0, G: 255, B: 0, A: 255}
	ColorBlue    = &core.Color{R: 0, G: 0, B: 255, A: 255}
	ColorGray    = &core.Color{R: 128, G: 128, B: 128, A: 255}
	ColorLightGray = &core.Color{R: 200, G: 200, B: 200, A: 255}
)

// Type aliases for convenience
type Offset = core.Offset

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

// Memo creates a memoized computed value
func Memo[T any](compute func() T) *reactive.Memo[T] {
	return reactive.NewMemo(compute)
}

// Computed creates a computed value with automatic dependency tracking
func Computed[T any](compute func() T) *reactive.Computed[T] {
	return reactive.NewComputed(compute)
}

// TextSignal creates a truly reactive text widget
func TextSignal[T any](signal *reactive.Signal[T], format func(T) string) widgets.WidgetImpl {
	id := fmt.Sprintf("reactive-text-%p", signal)
	
	// Create text widget with initial value (without tracking)
	initialValue := format(signal.Peek())
	text := widgets.NewText(id, initialValue)
	logger.Trace(logger.TagSignal, "Created TEXT widget with initial value: %s", initialValue)
	
	// Create effect that updates ONLY this widget's text
	// This effect will track the signal dependency on first run
	effect := reactive.CreateEffect(func() {
		// This Get() will register this effect as an observer
		newValue := format(signal.Get())
		logger.Trace(logger.TagSignal, "TEXT effect updating widget text to: %s", newValue)
		text.SetText(newValue)
		
		// Mark this specific widget for repaint
		text.MarkNeedsRepaint()
		
		// Schedule selective DOM update
		if globalApp != nil && globalApp.batcher != nil {
			globalApp.batcher.Add(func() {
				logger.Trace(logger.TagSignal, "Batched TEXT DOM update for: %s", id)
				globalApp.updateWidget(text)
			})
		}
	})
	
	logger.Trace(logger.TagSignal, "TEXT signal effect created with ID: %v", effect)
	
	return text
}

// TextMemo creates a reactive text widget from a computed value
func TextMemo[T any](memo *reactive.Memo[T], format func(T) string) widgets.WidgetImpl {
	id := fmt.Sprintf("reactive-memo-%p", memo)
	
	// Create text widget with initial value
	initialValue := format(memo.Peek())
	text := widgets.NewText(id, initialValue)
	logger.Trace(logger.TagMemo, "Created TEXT widget with initial value: %s", initialValue)
	
	// Create effect that updates when memo changes
	effect := reactive.CreateEffect(func() {
		newValue := format(memo.Get())
		logger.Trace(logger.TagMemo, "TEXT memo updated widget text to: %s", newValue)
		text.SetText(newValue)
		text.MarkNeedsRepaint()
		
		// Schedule selective DOM update
		if globalApp != nil && globalApp.batcher != nil {
			globalApp.batcher.Add(func() {
				logger.Trace(logger.TagMemo, "Batched TEXT DOM update for: %s", id)
				globalApp.updateWidget(text)
			})
		}
	})
	
	logger.Trace(logger.TagMemo, "TEXT memo effect created with ID: %v", effect)
	
	return text
}
