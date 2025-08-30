//go:build wasm
// +build wasm

package render

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/graph"
	"github.com/maya-framework/maya/internal/widgets"
	"github.com/maya-framework/maya/internal/workflow"
)

// Pipeline implements the multipass rendering system using the actual components
type Pipeline struct {
	// Use the REAL workflow engine
	engine *workflow.WorkflowEngine

	// Use the REAL graph for dependencies
	dependencies *graph.Graph

	// Reference to the REAL tree
	tree *core.Tree

	// DOM container
	container js.Value

	// Theme
	theme *Theme
	
	// Node to DOM element mapping
	nodeElements map[*core.Node]js.Value
}

// Theme for styling
type Theme struct {
	Primary    string
	Text       string
	Background string
	FontFamily string
}

// NewPipeline creates a new rendering pipeline using existing components
func NewPipeline(tree *core.Tree, container js.Value, theme *Theme) *Pipeline {
	p := &Pipeline{
		engine:       workflow.NewWorkflowEngine("render-pipeline"),
		dependencies: graph.NewGraph(),
		tree:         tree,
		container:    container,
		theme:        theme,
		nodeElements: make(map[*core.Node]js.Value),
	}

	p.setupStages()
	return p
}

// setupStages configures all rendering stages
func (p *Pipeline) setupStages() {
	// Create stages using workflow.Stage
	stages := []*workflow.Stage{
		{
			ID:   "mark-dirty",
			Name: "Mark Dirty Nodes",
			Execute: func(ctx context.Context, stageCtx *workflow.StageContext) error {
				// Use tree's DirtyNodes iterator
				for node := range p.tree.DirtyNodes() {
					p.propagateDirty(node)
				}
				stageCtx.Output = p.tree
				return nil
			},
		},
		{
			ID:   "calculate-sizes",
			Name: "Calculate Widget Sizes",
			Execute: func(ctx context.Context, stageCtx *workflow.StageContext) error {
				// Use PostOrderDFS for bottom-up calculation
				for node := range p.tree.PostOrderDFS() {
					p.calculateNodeSize(node)
				}
				stageCtx.Output = p.tree
				return nil
			},
		},
		{
			ID:   "assign-positions",
			Name: "Assign Positions",
			Execute: func(ctx context.Context, stageCtx *workflow.StageContext) error {
				// Use PreOrderDFS for top-down positioning
				for node := range p.tree.PreOrderDFS() {
					p.assignNodePosition(node)
				}
				stageCtx.Output = p.tree
				return nil
			},
		},
		{
			ID:   "commit-dom",
			Name: "Commit to DOM",
			Execute: func(ctx context.Context, stageCtx *workflow.StageContext) error {
				p.commitToDOM()
				stageCtx.Output = p.tree
				return nil
			},
		},
	}

	// Add stages to workflow engine
	for _, stage := range stages {
		p.engine.AddStage(stage)
	}

	// Setup dependencies in graph
	p.dependencies.AddNode(graph.NodeID("mark-dirty"), nil)
	p.dependencies.AddNode(graph.NodeID("calculate-sizes"), nil)
	p.dependencies.AddNode(graph.NodeID("assign-positions"), nil)
	p.dependencies.AddNode(graph.NodeID("commit-dom"), nil)

	p.dependencies.AddEdge(graph.NodeID("mark-dirty"), graph.NodeID("calculate-sizes"), 1.0)
	p.dependencies.AddEdge(graph.NodeID("calculate-sizes"), graph.NodeID("assign-positions"), 1.0)
	p.dependencies.AddEdge(graph.NodeID("assign-positions"), graph.NodeID("commit-dom"), 1.0)
}

// Execute runs the rendering pipeline
func (p *Pipeline) Execute(ctx context.Context) error {
	// Get topological order from dependency graph
	order, err := p.dependencies.TopologicalSort()
	if err != nil {
		return fmt.Errorf("failed to sort stages: %w", err)
	}

	// Execute stages in order
	stageCtx := &workflow.StageContext{
		Input:    p.tree,
		Metadata: make(map[string]interface{}),
	}

	for _, nodeID := range order {
		stageID := string(nodeID)
		if stage, exists := p.engine.GetStage(stageID); exists {
			stageCtx.Stage = stage
			if err := stage.Execute(ctx, stageCtx); err != nil {
				return fmt.Errorf("stage %s failed: %w", stageID, err)
			}
			// Use output as input for next stage
			stageCtx.Input = stageCtx.Output
		}
	}

	return nil
}

// propagateDirty marks ancestors as dirty
func (p *Pipeline) propagateDirty(node *core.Node) {
	for ancestor := range node.Ancestors() {
		ancestor.MarkDirty(core.LayoutDirty)
	}
}

// calculateNodeSize calculates size based on widget type and children
func (p *Pipeline) calculateNodeSize(node *core.Node) {
	if node.Widget == nil {
		return
	}

	switch w := node.Widget.(type) {
	case *widgets.Column:
		// Vertical layout: sum heights, max width
		totalHeight := 0.0
		maxWidth := 0.0
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
		totalWidth := 0.0
		maxHeight := 0.0
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
		constraints := core.Constraints{
			MinWidth:  0,
			MaxWidth:  800,
			MinHeight: 0,
			MaxHeight: 600,
		}
		width, height := w.Layout(constraints)
		node.Bounds.Width = width
		node.Bounds.Height = height
	}
}

// assignNodePosition assigns position based on parent and layout type
func (p *Pipeline) assignNodePosition(node *core.Node) {
	if node.Parent == nil {
		// Root node
		node.Bounds.X = 0
		node.Bounds.Y = 0
		return
	}

	// Get actual parent node (dereference weak pointer)
	parentPtr := node.Parent.Value()
	if parentPtr == nil || parentPtr.Widget == nil {
		return
	}
	parent := parentPtr // Use pointer directly, not copy

	// Find child index
	childIndex := -1
	for i, child := range parent.Children {
		if child == node {
			childIndex = i
			break
		}
	}

	if childIndex == -1 {
		return
	}

	// Calculate position based on parent's layout type
	switch parent.Widget.(type) {
	case *widgets.Column:
		// Vertical layout
		node.Bounds.X = parent.Bounds.X
		node.Bounds.Y = parent.Bounds.Y

		// Add offset from previous siblings
		for i := 0; i < childIndex; i++ {
			node.Bounds.Y += parent.Children[i].Bounds.Height + 10 // gap
		}

	case *widgets.Row:
		// Horizontal layout
		node.Bounds.X = parent.Bounds.X
		node.Bounds.Y = parent.Bounds.Y

		// Add offset from previous siblings
		for i := 0; i < childIndex; i++ {
			node.Bounds.X += parent.Children[i].Bounds.Width + 10 // gap
		}

	default:
		// Default positioning
		node.Bounds.X = parent.Bounds.X
		node.Bounds.Y = parent.Bounds.Y
	}
}

// commitToDOM renders the tree to the DOM
func (p *Pipeline) commitToDOM() {
	// Clear container and mapping
	p.container.Set("innerHTML", "")
	p.nodeElements = make(map[*core.Node]js.Value)
	
	// Build DOM from tree recursively
	if root := p.tree.GetRoot(); root != nil {
		p.createDOMTree(root, p.container)
	}
}

// createDOMTree creates DOM elements recursively
func (p *Pipeline) createDOMTree(node *core.Node, parentElement js.Value) {
	if node.Widget == nil {
		return
	}

	doc := js.Global().Get("document")
	var elem js.Value

	// Create appropriate element based on widget type
	switch widget := node.Widget.(type) {
	case *widgets.Text:
		elem = doc.Call("createElement", "span")
		elem.Set("textContent", widget.GetText())

	case *widgets.Button:
		elem = doc.Call("createElement", "button")
		elem.Set("textContent", widget.GetLabel())
		
		// Register callback without js.FuncOf
		callbackID := RegisterCallback(widget.Click)
		
		// Set onclick to call our exported function
		elem.Set("onclick", js.Global().Get("Function").New(
			"return window.handleEvent("+fmt.Sprintf("%d", callbackID)+");",
		))

	case *widgets.Container, *widgets.Column, *widgets.Row:
		elem = doc.Call("createElement", "div")

	default:
		elem = doc.Call("createElement", "div")
	}

	// Apply calculated positioning
	style := elem.Get("style")
	style.Set("position", "absolute")
	style.Set("left", fmt.Sprintf("%fpx", node.Bounds.X))
	style.Set("top", fmt.Sprintf("%fpx", node.Bounds.Y))

	if node.Bounds.Width > 0 {
		style.Set("width", fmt.Sprintf("%fpx", node.Bounds.Width))
	}
	if node.Bounds.Height > 0 {
		style.Set("height", fmt.Sprintf("%fpx", node.Bounds.Height))
	}

	// Apply theme styles
	if widgetImpl, ok := node.Widget.(widgets.WidgetImpl); ok {
		p.applyThemeStyles(elem, widgetImpl)
	}

	// Append to parent element
	parentElement.Call("appendChild", elem)
	
	// Store mapping
	p.nodeElements[node] = elem
	
	// Recursively create children
	for _, child := range node.Children {
		p.createDOMTree(child, elem)
	}
}

// applyThemeStyles applies visual styles based on widget type
func (p *Pipeline) applyThemeStyles(elem js.Value, widget widgets.WidgetImpl) {
	style := elem.Get("style")

	switch widget.(type) {
	case *widgets.Text:
		style.Set("color", p.theme.Text)
		style.Set("fontSize", "16px")
		style.Set("fontFamily", p.theme.FontFamily)

	case *widgets.Button:
		style.Set("backgroundColor", p.theme.Primary)
		style.Set("color", "white")
		style.Set("border", "none")
		style.Set("borderRadius", "5px")
		style.Set("cursor", "pointer")
		style.Set("fontSize", "14px")
		style.Set("fontFamily", p.theme.FontFamily)
		style.Set("textAlign", "center")
	}
}
