//go:build wasm
// +build wasm

package maya

import (
	"context"
	"fmt"
	"syscall/js"
	"time"

	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/graph"
	"github.com/maya-framework/maya/internal/reactive"
	"github.com/maya-framework/maya/internal/widgets"
	"github.com/maya-framework/maya/internal/workflow"
)

// App represents a Maya application following the planned architecture
type App struct {
	// UI Tree (core.Node based)
	rootWidget widgets.WidgetImpl
	uiTree     *UITree

	// Rendering Pipeline (workflow based)
	renderPipeline *RenderPipeline

	// Reactive System
	batcher *reactive.UpdateBatcher

	// Context
	ctx    context.Context
	cancel context.CancelFunc

	// Theme
	theme *Theme

	// DOM
	container js.Value
}

// UITree wraps the core.Node tree structure
type UITree struct {
	root       *core.Node
	nodeMap    map[core.NodeID]*core.Node
	dirtyNodes []*core.Node
}

// RenderPipeline implements the multipass rendering system
type RenderPipeline struct {
	// Workflow engine for stage orchestration
	engine *workflow.WorkflowEngine

	// Individual render stages
	markDirtyStage       *workflow.Stage
	intrinsicStage       *workflow.Stage
	constraintStage      *workflow.Stage
	sizeCalculationStage *workflow.Stage
	positionStage        *workflow.Stage
	paintStage           *workflow.Stage
	commitStage          *workflow.Stage

	// Stage dependency graph (NOT the UI tree!)
	stageDependencies *graph.Graph
}

// New creates a new Maya application with proper architecture
func New(root Component) *App {
	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		ctx:    ctx,
		cancel: cancel,
		theme:  DefaultTheme(),
		uiTree: &UITree{
			nodeMap: make(map[core.NodeID]*core.Node),
		},
	}

	// Build widget tree
	app.rootWidget = root()

	// Create UI tree from widgets
	app.uiTree.root = app.buildUITree(app.rootWidget)

	// Setup render pipeline
	app.renderPipeline = app.createRenderPipeline()

	// Setup reactive batching
	app.batcher = reactive.NewUpdateBatcher()

	return app
}

// buildUITree converts widgets to core.Node tree
func (app *App) buildUITree(widget widgets.WidgetImpl) *core.Node {
	node := core.NewNode(widget.ID(), widget)
	app.uiTree.nodeMap[core.NodeID(widget.ID())] = node

	// Recursively build children
	for _, child := range widget.Children() {
		childNode := app.buildUITree(child)
		node.AddChild(childNode)
	}

	return node
}

// createRenderPipeline sets up the multipass rendering stages
func (app *App) createRenderPipeline() *RenderPipeline {
	pipeline := &RenderPipeline{
		engine:            workflow.NewWorkflowEngine("render-pipeline"),
		stageDependencies: graph.NewGraph(),
	}

	// Create render stages (each processes the ENTIRE UI tree)

	// Stage 1: Mark dirty nodes (BFS traversal)
	pipeline.markDirtyStage = &workflow.Stage{
		ID:   "mark-dirty",
		Name: "Mark Dirty Nodes",
		Execute: func(ctx context.Context, stageCtx *workflow.StageContext) error {
			tree := stageCtx.Input.(*UITree)

			// BFS to propagate dirty flags
			queue := []*core.Node{tree.root}
			for len(queue) > 0 {
				node := queue[0]
				queue = queue[1:]

				if node.IsDirty() {
					tree.dirtyNodes = append(tree.dirtyNodes, node)
					// Propagate to ancestors
					app.propagateDirtyUp(node)
				}

				queue = append(queue, node.Children...)
			}

			stageCtx.Output = tree
			return nil
		},
	}

	// Stage 2: Calculate intrinsic dimensions (Post-order DFS)
	pipeline.intrinsicStage = &workflow.Stage{
		ID:   "intrinsic",
		Name: "Calculate Intrinsic Dimensions",
		Execute: func(ctx context.Context, stageCtx *workflow.StageContext) error {
			tree := stageCtx.Input.(*UITree)

			// Post-order DFS (children first, then parent)
			var calculateIntrinsic func(*core.Node) error
			calculateIntrinsic = func(node *core.Node) error {
				// Process children first
				for _, child := range node.Children {
					if err := calculateIntrinsic(child); err != nil {
						return err
					}
				}

				// Then calculate this node's intrinsic size
				if widget := node.Widget; widget != nil {
					// Children have been processed, we can use their sizes
					childWidths := make([]float64, len(node.Children))
					childHeights := make([]float64, len(node.Children))

					for i, child := range node.Children {
						childWidths[i] = child.Widget.GetIntrinsicWidth(0)
						childHeights[i] = child.Widget.GetIntrinsicHeight(0)
					}

					// Store intrinsic dimensions (would be used by widget)
				}

				return nil
			}

			err := calculateIntrinsic(tree.root)
			stageCtx.Output = tree
			return err
		},
	}

	// Stage 3: Resolve constraints (Pre-order DFS)
	pipeline.constraintStage = &workflow.Stage{
		ID:   "constraints",
		Name: "Resolve Constraints",
		Execute: func(ctx context.Context, stageCtx *workflow.StageContext) error {
			tree := stageCtx.Input.(*UITree)

			// Pre-order DFS (parent first, then children)
			var resolveConstraints func(*core.Node, core.Constraints) error
			resolveConstraints = func(node *core.Node, parentConstraints core.Constraints) error {
				// Resolve this node's constraints first
				node.ResolvedConstraints = parentConstraints

				// Then propagate to children with adjusted constraints
				for _, child := range node.Children {
					childConstraints := parentConstraints // Adjust based on layout type
					if err := resolveConstraints(child, childConstraints); err != nil {
						return err
					}
				}

				return nil
			}

			// Start with viewport constraints
			rootConstraints := core.Constraints{
				MinWidth:  0,
				MaxWidth:  1920,
				MinHeight: 0,
				MaxHeight: 1080,
			}

			err := resolveConstraints(tree.root, rootConstraints)
			stageCtx.Output = tree
			return err
		},
	}

	// Stage 4: Calculate final sizes (Post-order DFS)
	pipeline.sizeCalculationStage = &workflow.Stage{
		ID:   "size-calculation",
		Name: "Calculate Final Sizes",
		Execute: func(ctx context.Context, stageCtx *workflow.StageContext) error {
			tree := stageCtx.Input.(*UITree)

			// Post-order DFS - children sizes needed for parent
			var calculateSizes func(*core.Node) error
			calculateSizes = func(node *core.Node) error {
				// Calculate children first
				for _, child := range node.Children {
					if err := calculateSizes(child); err != nil {
						return err
					}
				}

				// Now calculate this node's size based on children
				if widget := node.Widget; widget != nil {
					// For layout widgets, calculate based on children
					switch w := widget.(type) {
					case *widgets.Column:
						// Vertical layout: sum heights, max width
						totalHeight := float64(0)
						maxWidth := float64(0)
						for i, child := range node.Children {
							totalHeight += child.Bounds.Height
							if i > 0 {
								totalHeight += 10 // gap
							}
							if child.Bounds.Width > maxWidth {
								maxWidth = child.Bounds.Width
							}
						}
						node.Bounds.Width = maxWidth
						node.Bounds.Height = totalHeight
						
					case *widgets.Row:
						// Horizontal layout: sum widths, max height
						totalWidth := float64(0)
						maxHeight := float64(0)
						for i, child := range node.Children {
							totalWidth += child.Bounds.Width
							if i > 0 {
								totalWidth += 10 // gap
							}
							if child.Bounds.Height > maxHeight {
								maxHeight = child.Bounds.Height
							}
						}
						node.Bounds.Width = totalWidth
						node.Bounds.Height = maxHeight
						
					default:
						// Regular widget uses its own layout
						width, height := w.Layout(node.ResolvedConstraints)
						node.Bounds.Width = width
						node.Bounds.Height = height
					}
				}

				return nil
			}

			err := calculateSizes(tree.root)
			stageCtx.Output = tree
			return err
		},
	}

	// Stage 5: Assign positions (Pre-order DFS)
	pipeline.positionStage = &workflow.Stage{
		ID:   "position",
		Name: "Assign Positions",
		Execute: func(ctx context.Context, stageCtx *workflow.StageContext) error {
			tree := stageCtx.Input.(*UITree)

			// Pre-order DFS - parent position needed for children
			var assignPositions func(*core.Node, core.Offset) error
			assignPositions = func(node *core.Node, parentOffset core.Offset) error {
				// Set this node's position
				node.CachedPosition = parentOffset
				node.Bounds.X = parentOffset.X
				node.Bounds.Y = parentOffset.Y

				// Calculate child positions based on layout
				childOffset := parentOffset
				
				// Determine layout type
				switch node.Widget.(type) {
				case *widgets.Row:
					// Horizontal layout
					for _, child := range node.Children {
						if err := assignPositions(child, childOffset); err != nil {
							return err
						}
						childOffset.X += child.Bounds.Width + 10 // gap
					}
				case *widgets.Column:
					// Vertical layout
					for _, child := range node.Children {
						if err := assignPositions(child, childOffset); err != nil {
							return err
						}
						childOffset.Y += child.Bounds.Height + 10 // gap
					}
				default:
					// Default vertical layout
					for _, child := range node.Children {
						if err := assignPositions(child, childOffset); err != nil {
							return err
						}
						childOffset.Y += child.Bounds.Height
					}
				}

				return nil
			}

			err := assignPositions(tree.root, core.Offset{X: 0, Y: 0})
			stageCtx.Output = tree
			return err
		},
	}

	// Stage 6: Paint (DFS)
	pipeline.paintStage = &workflow.Stage{
		ID:   "paint",
		Name: "Paint Widgets",
		Execute: func(ctx context.Context, stageCtx *workflow.StageContext) error {
			tree := stageCtx.Input.(*UITree)
			container := stageCtx.Metadata["container"].(js.Value)

			// Create paint context
			paintCtx := &domPaintContext{
				element: container,
				theme:   app.theme,
			}

			// DFS paint traversal
			var paint func(*core.Node) error
			paint = func(node *core.Node) error {
				if widget := node.Widget; widget != nil {
					widget.Paint(paintCtx)
				}

				for _, child := range node.Children {
					if err := paint(child); err != nil {
						return err
					}
				}

				return nil
			}

			err := paint(tree.root)
			stageCtx.Output = tree
			return err
		},
	}

	// Stage 7: Commit to DOM
	pipeline.commitStage = &workflow.Stage{
		ID:   "commit",
		Name: "Commit to DOM",
		Execute: func(ctx context.Context, stageCtx *workflow.StageContext) error {
			tree := stageCtx.Input.(*UITree)
			container := stageCtx.Metadata["container"].(js.Value)

			// Clear and rebuild DOM (simplified for POC)
			container.Set("innerHTML", "")

			// Build DOM from tree
			var buildDOM func(*core.Node, js.Value)
			buildDOM = func(node *core.Node, parent js.Value) {
				if widgetImpl, ok := node.Widget.(widgets.WidgetImpl); ok {
					doc := js.Global().Get("document")
					var elem js.Value
					
					// Create appropriate element based on widget type
					switch widget := widgetImpl.(type) {
					case *widgets.Text:
						elem = doc.Call("createElement", "span")
						// Get text content from widget
						elem.Set("textContent", widget.GetText())
						
					case *widgets.Button:
						elem = doc.Call("createElement", "button")
						// Get button label
						elem.Set("textContent", widget.GetLabel())
						// Add click handler
						elem.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
							widget.Click()
							return nil
						}))
						
					case *widgets.Container:
						elem = doc.Call("createElement", "div")
					
					case *widgets.Column:
						elem = doc.Call("createElement", "div")
						
					case *widgets.Row:
						elem = doc.Call("createElement", "div")
						
					default:
						elem = doc.Call("createElement", "div")
					}

					// Maya calcula TUDO - não usa CSS layout
					style := elem.Get("style")
					style.Set("position", "absolute")
					style.Set("left", fmt.Sprintf("%fpx", node.Bounds.X))
					style.Set("top", fmt.Sprintf("%fpx", node.Bounds.Y))
					
					// Set calculated sizes
					if node.Bounds.Width > 0 {
						style.Set("width", fmt.Sprintf("%fpx", node.Bounds.Width))
					}
					if node.Bounds.Height > 0 {
						style.Set("height", fmt.Sprintf("%fpx", node.Bounds.Height))
					}
					
					// Visual styling apenas (não layout)
					switch widgetImpl.(type) {
					case *widgets.Text:
						style.Set("color", app.theme.Text)
						style.Set("fontSize", "16px")
						style.Set("fontFamily", app.theme.FontFamily)
					case *widgets.Button:
						style.Set("backgroundColor", app.theme.Primary)
						style.Set("color", "white")
						style.Set("border", "none")
						style.Set("borderRadius", "5px")
						style.Set("cursor", "pointer")
						style.Set("fontSize", "14px")
						style.Set("fontFamily", app.theme.FontFamily)
						style.Set("textAlign", "center")
						style.Set("lineHeight", fmt.Sprintf("%fpx", node.Bounds.Height))
					}

					parent.Call("appendChild", elem)

					// Recurse for children
					for _, child := range node.Children {
						buildDOM(child, elem)
					}
				}
			}

			buildDOM(tree.root, container)

			// Clear dirty flags
			for _, node := range tree.dirtyNodes {
				node.ClearDirty()
			}
			tree.dirtyNodes = nil

			return nil
		},
	}

	// Setup stage dependencies in the graph
	// The graph represents dependencies between STAGES, not widgets!
	pipeline.stageDependencies.AddNode(graph.NodeID("mark-dirty"), pipeline.markDirtyStage)
	pipeline.stageDependencies.AddNode(graph.NodeID("intrinsic"), pipeline.intrinsicStage)
	pipeline.stageDependencies.AddNode(graph.NodeID("constraints"), pipeline.constraintStage)
	pipeline.stageDependencies.AddNode(graph.NodeID("size-calculation"), pipeline.sizeCalculationStage)
	pipeline.stageDependencies.AddNode(graph.NodeID("position"), pipeline.positionStage)
	pipeline.stageDependencies.AddNode(graph.NodeID("paint"), pipeline.paintStage)
	pipeline.stageDependencies.AddNode(graph.NodeID("commit"), pipeline.commitStage)

	// Define dependencies
	pipeline.stageDependencies.AddEdge(graph.NodeID("mark-dirty"), graph.NodeID("intrinsic"), 1.0)
	pipeline.stageDependencies.AddEdge(graph.NodeID("intrinsic"), graph.NodeID("constraints"), 1.0)
	pipeline.stageDependencies.AddEdge(graph.NodeID("constraints"), graph.NodeID("size-calculation"), 1.0)
	pipeline.stageDependencies.AddEdge(graph.NodeID("size-calculation"), graph.NodeID("position"), 1.0)
	pipeline.stageDependencies.AddEdge(graph.NodeID("position"), graph.NodeID("paint"), 1.0)
	pipeline.stageDependencies.AddEdge(graph.NodeID("paint"), graph.NodeID("commit"), 1.0)

	// Add stages to workflow engine
	pipeline.engine.AddStage(pipeline.markDirtyStage)
	pipeline.engine.AddStage(pipeline.intrinsicStage)
	pipeline.engine.AddStage(pipeline.constraintStage)
	pipeline.engine.AddStage(pipeline.sizeCalculationStage)
	pipeline.engine.AddStage(pipeline.positionStage)
	pipeline.engine.AddStage(pipeline.paintStage)
	pipeline.engine.AddStage(pipeline.commitStage)

	return pipeline
}

// Run starts the application
func (app *App) Run() {
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

		// Initial render
		app.render()

		// Setup reactive updates
		app.setupReactiveLoop()
	})

	// Keep running
	select {}
}

// render executes the multipass rendering pipeline
func (app *App) render() {
	// Mark root as dirty for initial render
	app.uiTree.root.MarkDirty(core.LayoutDirty | core.PaintDirty)

	// Create stage context
	stageCtx := &workflow.StageContext{
		Input: app.uiTree,
		Metadata: map[string]interface{}{
			"container": app.container,
			"theme":     app.theme,
		},
	}

	// Execute the pipeline in topological order
	order, err := app.renderPipeline.stageDependencies.TopologicalSort()
	if err != nil {
		fmt.Println("Error in topological sort:", err)
		return
	}

	// Execute stages in order
	for _, nodeID := range order {
		node, _ := app.renderPipeline.stageDependencies.GetNode(nodeID)
		if stage, ok := node.Data.(*workflow.Stage); ok {
			if err := stage.Execute(app.ctx, stageCtx); err != nil {
				fmt.Println("Error in stage", stage.Name, ":", err)
				return
			}
			// Pass output to next stage
			stageCtx.Input = stageCtx.Output
		}
	}
}

// setupReactiveLoop sets up the reactive update loop
func (app *App) setupReactiveLoop() {
	// Create a ticker for frame updates (60 FPS)
	ticker := time.NewTicker(16 * time.Millisecond)

	go func() {
		for {
			select {
			case <-app.ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				// Check if any nodes are dirty
				if app.hasDirtyNodes() {
					app.render()
				}
			}
		}
	}()
}

// hasDirtyNodes checks if any nodes need re-rendering
func (app *App) hasDirtyNodes() bool {
	// DFS check for dirty nodes
	var checkDirty func(*core.Node) bool
	checkDirty = func(node *core.Node) bool {
		if node.IsDirty() {
			return true
		}
		for _, child := range node.Children {
			if checkDirty(child) {
				return true
			}
		}
		return false
	}

	return checkDirty(app.uiTree.root)
}

// propagateDirtyUp propagates dirty flag to ancestors
func (app *App) propagateDirtyUp(node *core.Node) {
	parent := node.GetParent()
	for parent != nil {
		parent.MarkDirty(core.ChildrenDirty)
		parent = parent.GetParent()
	}
}

// The rest remains similar but uses the correct architecture...

// Component type
type Component func() widgets.WidgetImpl

// Helper functions remain the same
func Container(children ...widgets.WidgetImpl) widgets.WidgetImpl {
	// Use Column as default container for multiple children
	if len(children) > 1 {
		return widgets.NewColumn("container", children...)
	} else if len(children) == 1 {
		container := widgets.NewContainer("container")
		container.SetChild(children[0])
		return container
	}
	return widgets.NewContainer("container")
}

func Text(content string) widgets.WidgetImpl {
	return widgets.NewText("text-" + content, content)
}

func Title(content string) widgets.WidgetImpl {
	text := widgets.NewText("title-" + content, content)
	// Set title style (larger font)
	text.SetStyle(widgets.TextStyle{
		FontFamily: "system-ui",
		FontSize:   24,
		FontWeight: widgets.FontWeightBold,
		Color:      widgets.ColorBlack,
		LineHeight: 1.2,
	})
	return text
}

func Button(text string, onClick func()) widgets.WidgetImpl {
	return widgets.NewButton("button-" + text, text, onClick)
}

func Row(children ...widgets.WidgetImpl) widgets.WidgetImpl {
	return widgets.NewRow("row", children...)
}

func Column(children ...widgets.WidgetImpl) widgets.WidgetImpl {
	return widgets.NewColumn("column", children...)
}

func Signal[T comparable](initial T) *reactive.Signal[T] {
	return reactive.NewSignal(initial)
}

func TextSignal[T any](signal *reactive.Signal[T], format func(T) string) widgets.WidgetImpl {
	text := widgets.NewText("reactive-text", format(signal.Get()))

	// Subscribe to changes
	signal.Subscribe(func(value T) {
		// TODO: Update text widget
		// This needs to trigger a re-render
	})

	return text
}

// domPaintContext implements core.PaintContext for DOM
type domPaintContext struct {
	element js.Value
	theme   *Theme
}

func (d *domPaintContext) DrawRect(bounds core.Bounds, paint core.Paint) {
	// Implementation...
}

func (d *domPaintContext) DrawText(text string, offset core.Offset, paint core.Paint) {
	// Implementation...
}

func (d *domPaintContext) DrawPath(path []core.Offset, paint core.Paint) {
	// Implementation...
}

func (d *domPaintContext) PushTransform(transform core.Transform) {
	// Implementation...
}

func (d *domPaintContext) PopTransform() {
	// Implementation...
}

func (d *domPaintContext) PushClip(bounds core.Bounds) {
	// Implementation...
}

func (d *domPaintContext) PopClip() {
	// Implementation...
}

// Helper functions
func (app *App) setupViewport() {
	doc := js.Global().Get("document")
	body := doc.Get("body")

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

// Theme remains the same
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
