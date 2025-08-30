//go:build wasm
// +build wasm

package maya

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/maya-framework/maya/internal/graph"
	"github.com/maya-framework/maya/internal/reactive"
)

// App represents a Maya application
type App struct {
	root     Widget
	renderer *renderer
	ctx      context.Context
	cancel   context.CancelFunc
	theme    *Theme
}

// Widget is the public interface for UI components
type Widget interface {
	Build(context.Context) js.Value
	Style() Style
}

// Component is a functional component
type Component func() Widget

// Style represents widget styling
type Style struct {
	Width           float64
	Height          float64
	Padding         float64
	Margin          float64
	BackgroundColor string
	Color           string
	FontSize        float64
	FontWeight      string
	Border          string
	BorderRadius    float64
	Display         string
	FlexDirection   string
	Gap             float64
	AlignItems      string
	JustifyContent  string
	Cursor          string
	Transition      string
	BoxShadow       string
}

// Theme defines the application theme
type Theme struct {
	Primary      string
	Secondary    string
	Background   string
	Surface      string
	Text         string
	TextLight    string
	FontFamily   string
	BorderRadius float64
	Spacing      float64
}

// DefaultTheme returns the default Maya theme
func DefaultTheme() *Theme {
	return &Theme{
		Primary:      "#667eea",
		Secondary:    "#48bb78",
		Background:   "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
		Surface:      "#ffffff",
		Text:         "#333333",
		TextLight:    "#666666",
		FontFamily:   "system-ui, -apple-system, sans-serif",
		BorderRadius: 8,
		Spacing:      10,
	}
}

// New creates a new Maya application
func New(root Component) *App {
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		root:     root(),
		renderer: newRenderer(),
		ctx:      ctx,
		cancel:   cancel,
		theme:    DefaultTheme(),
	}
}

// Run starts the application
func (app *App) Run() {
	// Set the app context
	app.setAppContext()
	
	// Initialize when DOM is ready
	waitForDOM(func() {
		// Setup viewport
		app.setupViewport()

		// Get or create container
		container := js.Global().Get("document").Call("getElementById", "app")
		if container.IsNull() {
			body := js.Global().Get("document").Get("body")
			container = js.Global().Get("document").Call("createElement", "div")
			container.Set("id", "app")
			body.Call("appendChild", container)
		}

		// Clear loading state
		container.Set("innerHTML", "")

		// Apply container styles
		app.applyStyle(container, Style{
			Width:          0, // Full width
			Height:         0, // Full height
			Padding:        app.theme.Spacing * 2,
			Display:        "flex",
			FlexDirection:  "column",
			AlignItems:     "center",
			JustifyContent: "center",
		})

		// Create wrapper for content
		wrapper := js.Global().Get("document").Call("createElement", "div")
		app.applyStyle(wrapper, Style{
			BackgroundColor: app.theme.Surface,
			Padding:         app.theme.Spacing * 3,
			BorderRadius:    app.theme.BorderRadius,
			BoxShadow:       "0 10px 40px rgba(0,0,0,0.2)",
		})

		// Render the app
		dom := app.root.Build(app.ctx)
		wrapper.Call("appendChild", dom)
		container.Call("appendChild", wrapper)
	})

	// Keep running
	select {}
}

// SetTheme sets a custom theme
func (app *App) SetTheme(theme *Theme) {
	app.theme = theme
}

// Signal creates a reactive signal
func Signal[T comparable](initial T) *reactive.Signal[T] {
	return reactive.NewSignal(initial)
}

// Container creates a container widget
func Container(children ...Widget) Widget {
	return &containerWidget{
		children: children,
		style: Style{
			Display:       "flex",
			FlexDirection: "column",
			Gap:           10,
		},
	}
}

// Text creates a text widget
func Text(content string) Widget {
	return &textWidget{
		content: content,
		style: Style{
			FontSize: 16,
			Color:    "#333333",
		},
	}
}

// Title creates a title text widget
func Title(content string) Widget {
	return &textWidget{
		content: content,
		style: Style{
			FontSize:   24,
			FontWeight: "bold",
			Color:      "#333333",
			Margin:     10,
		},
	}
}

// TextSignal creates a reactive text widget
func TextSignal[T any](signal *reactive.Signal[T], format func(T) string) Widget {
	return &reactiveTextWidget[T]{
		signal: signal,
		format: format,
		style: Style{
			FontSize: 16,
			Color:    "#333333",
		},
	}
}

// Button creates a button widget
func Button(text string, onClick func()) Widget {
	return &buttonWidget{
		text:    text,
		onClick: onClick,
		style: Style{
			Padding:         12,
			BackgroundColor: "#667eea",
			Color:           "#ffffff",
			FontSize:        16,
			FontWeight:      "500",
			Border:          "none",
			BorderRadius:    5,
			Cursor:          "pointer",
			Transition:      "all 0.3s ease",
		},
	}
}

// PrimaryButton creates a primary button
func PrimaryButton(text string, onClick func()) Widget {
	return &buttonWidget{
		text:    text,
		onClick: onClick,
		style: Style{
			Padding:         12,
			BackgroundColor: "#667eea",
			Color:           "#ffffff",
			FontSize:        16,
			FontWeight:      "500",
			Border:          "none",
			BorderRadius:    5,
			Cursor:          "pointer",
			Transition:      "all 0.3s ease",
		},
	}
}

// SecondaryButton creates a secondary button
func SecondaryButton(text string, onClick func()) Widget {
	return &buttonWidget{
		text:    text,
		onClick: onClick,
		style: Style{
			Padding:         12,
			BackgroundColor: "#48bb78",
			Color:           "#ffffff",
			FontSize:        16,
			FontWeight:      "500",
			Border:          "none",
			BorderRadius:    5,
			Cursor:          "pointer",
			Transition:      "all 0.3s ease",
		},
	}
}

// Column creates a vertical layout
func Column(children ...Widget) Widget {
	return &columnWidget{
		children: children,
		style: Style{
			Display:       "flex",
			FlexDirection: "column",
			Gap:           10,
		},
	}
}

// Row creates a horizontal layout
func Row(children ...Widget) Widget {
	return &rowWidget{
		children: children,
		style: Style{
			Display:       "flex",
			FlexDirection: "row",
			Gap:           10,
			AlignItems:    "center",
		},
	}
}

// Spacer creates empty space
func Spacer(height float64) Widget {
	return &spacerWidget{
		height: height,
	}
}

// --- Internal widget implementations ---

type containerWidget struct {
	children []Widget
	style    Style
}

func (c *containerWidget) Build(ctx context.Context) js.Value {
	doc := js.Global().Get("document")
	div := doc.Call("createElement", "div")

	app := getApp(ctx)
	app.applyStyle(div, c.style)

	for _, child := range c.children {
		div.Call("appendChild", child.Build(ctx))
	}

	return div
}

func (c *containerWidget) Style() Style {
	return c.style
}

type textWidget struct {
	content string
	style   Style
}

func (t *textWidget) Build(ctx context.Context) js.Value {
	doc := js.Global().Get("document")
	span := doc.Call("createElement", "span")
	span.Set("textContent", t.content)

	app := getApp(ctx)
	app.applyStyle(span, t.style)

	return span
}

func (t *textWidget) Style() Style {
	return t.style
}

type reactiveTextWidget[T any] struct {
	signal *reactive.Signal[T]
	format func(T) string
	style  Style
}

func (t *reactiveTextWidget[T]) Build(ctx context.Context) js.Value {
	doc := js.Global().Get("document")
	span := doc.Call("createElement", "span")

	app := getApp(ctx)
	app.applyStyle(span, t.style)

	// Initial value
	span.Set("textContent", t.format(t.signal.Get()))

	// Subscribe to changes
	t.signal.Subscribe(func(value T) {
		span.Set("textContent", t.format(value))
	})

	return span
}

func (t *reactiveTextWidget[T]) Style() Style {
	return t.style
}

type buttonWidget struct {
	text    string
	onClick func()
	style   Style
}

func (b *buttonWidget) Build(ctx context.Context) js.Value {
	doc := js.Global().Get("document")
	button := doc.Call("createElement", "button")
	button.Set("textContent", b.text)

	app := getApp(ctx)
	app.applyStyle(button, b.style)

	// Add hover effect
	button.Call("addEventListener", "mouseenter", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		this.Get("style").Set("transform", "translateY(-2px)")
		this.Get("style").Set("boxShadow", "0 5px 15px rgba(0,0,0,0.2)")
		return nil
	}))

	button.Call("addEventListener", "mouseleave", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		this.Get("style").Set("transform", "translateY(0)")
		this.Get("style").Set("boxShadow", "none")
		return nil
	}))

	button.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if b.onClick != nil {
			b.onClick()
		}
		return nil
	}))

	return button
}

func (b *buttonWidget) Style() Style {
	return b.style
}

type columnWidget struct {
	children []Widget
	style    Style
}

func (c *columnWidget) Build(ctx context.Context) js.Value {
	doc := js.Global().Get("document")
	div := doc.Call("createElement", "div")

	app := getApp(ctx)
	app.applyStyle(div, c.style)

	for _, child := range c.children {
		div.Call("appendChild", child.Build(ctx))
	}

	return div
}

func (c *columnWidget) Style() Style {
	return c.style
}

type rowWidget struct {
	children []Widget
	style    Style
}

func (r *rowWidget) Build(ctx context.Context) js.Value {
	doc := js.Global().Get("document")
	div := doc.Call("createElement", "div")

	app := getApp(ctx)
	app.applyStyle(div, r.style)

	for _, child := range r.children {
		div.Call("appendChild", child.Build(ctx))
	}

	return div
}

func (r *rowWidget) Style() Style {
	return r.style
}

type spacerWidget struct {
	height float64
}

func (s *spacerWidget) Build(ctx context.Context) js.Value {
	doc := js.Global().Get("document")
	div := doc.Call("createElement", "div")
	div.Get("style").Set("height", fmt.Sprintf("%fpx", s.height))
	return div
}

func (s *spacerWidget) Style() Style {
	return Style{Height: s.height}
}

// --- Internal renderer ---

type renderer struct {
	multipass *multipassRenderer
}

func newRenderer() *renderer {
	return &renderer{
		multipass: &multipassRenderer{
			graph: graph.NewGraph(),
		},
	}
}

type multipassRenderer struct {
	graph *graph.Graph
}

// --- Utilities ---

func (app *App) setupViewport() {
	doc := js.Global().Get("document")
	body := doc.Get("body")

	// Reset body styles
	bodyStyle := body.Get("style")
	bodyStyle.Set("margin", "0")
	bodyStyle.Set("padding", "0")
	bodyStyle.Set("fontFamily", app.theme.FontFamily)
	bodyStyle.Set("background", app.theme.Background)
	bodyStyle.Set("minHeight", "100vh")
	bodyStyle.Set("display", "flex")
	bodyStyle.Set("alignItems", "center")
	bodyStyle.Set("justifyContent", "center")
}

func (app *App) applyStyle(element js.Value, style Style) {
	s := element.Get("style")

	if style.Width > 0 {
		s.Set("width", fmt.Sprintf("%fpx", style.Width))
	}
	if style.Height > 0 {
		s.Set("height", fmt.Sprintf("%fpx", style.Height))
	}
	if style.Padding > 0 {
		s.Set("padding", fmt.Sprintf("%fpx", style.Padding))
	}
	if style.Margin > 0 {
		s.Set("margin", fmt.Sprintf("%fpx", style.Margin))
	}
	if style.BackgroundColor != "" {
		s.Set("backgroundColor", style.BackgroundColor)
	}
	if style.Color != "" {
		s.Set("color", style.Color)
	}
	if style.FontSize > 0 {
		s.Set("fontSize", fmt.Sprintf("%fpx", style.FontSize))
	}
	if style.FontWeight != "" {
		s.Set("fontWeight", style.FontWeight)
	}
	if style.Border != "" {
		s.Set("border", style.Border)
	}
	if style.BorderRadius > 0 {
		s.Set("borderRadius", fmt.Sprintf("%fpx", style.BorderRadius))
	}
	if style.Display != "" {
		s.Set("display", style.Display)
	}
	if style.FlexDirection != "" {
		s.Set("flexDirection", style.FlexDirection)
	}
	if style.Gap > 0 {
		s.Set("gap", fmt.Sprintf("%fpx", style.Gap))
	}
	if style.AlignItems != "" {
		s.Set("alignItems", style.AlignItems)
	}
	if style.JustifyContent != "" {
		s.Set("justifyContent", style.JustifyContent)
	}
	if style.Cursor != "" {
		s.Set("cursor", style.Cursor)
	}
	if style.Transition != "" {
		s.Set("transition", style.Transition)
	}
	if style.BoxShadow != "" {
		s.Set("boxShadow", style.BoxShadow)
	}
}

func waitForDOM(callback func()) {
	doc := js.Global().Get("document")
	readyState := doc.Get("readyState").String()

	if readyState == "complete" || readyState == "interactive" {
		callback()
	} else {
		js.Global().Call("addEventListener", "DOMContentLoaded",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				callback()
				return nil
			}))
	}
}

var appContext *App

func getApp(ctx context.Context) *App {
	// In a real implementation, would get from context
	return appContext
}

// SetAppContext sets the global app context (internal use)
func (app *App) setAppContext() {
	appContext = app
}
