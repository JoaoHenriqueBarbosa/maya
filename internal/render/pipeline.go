//go:build wasm
// +build wasm

package render

import (
	"context"
	"fmt"

	"github.com/maya-framework/maya/internal/core"
	"github.com/maya-framework/maya/internal/graph"
	"github.com/maya-framework/maya/internal/logger"
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

	// Renderer abstraction - THE ONLY WAY TO RENDER
	renderer Renderer

	// Theme
	theme *Theme
	
	// Track previous commands for selective updates
	previousCommands []PaintCommand
	firstRender bool
}

// Theme for styling
type Theme struct {
	Primary    string
	Text       string
	Background string
	FontFamily string
}

// NewPipeline creates a new rendering pipeline with a specific renderer
func NewPipeline(tree *core.Tree, renderer Renderer, theme *Theme) *Pipeline {
	p := &Pipeline{
		engine:       workflow.NewWorkflowEngine("render-pipeline"),
		dependencies: graph.NewGraph(),
		tree:         tree,
		renderer:     renderer,
		theme:        theme,
		firstRender:  true,
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
	logger.Trace("PIPELINE", "Starting execution...")
	
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
			logger.Trace("PIPELINE", "Executing stage: %s", stageID)
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
		// Vertical layout - children start at 0,0 relative to parent
		node.Bounds.X = 0
		node.Bounds.Y = 0

		// Add offset from previous siblings
		for i := 0; i < childIndex; i++ {
			node.Bounds.Y += parent.Children[i].Bounds.Height + 10 // gap
		}

	case *widgets.Row:
		// Horizontal layout - children start at 0,0 relative to parent
		node.Bounds.X = 0
		node.Bounds.Y = 0

		// Add offset from previous siblings
		for i := 0; i < childIndex; i++ {
			node.Bounds.X += parent.Children[i].Bounds.Width + 10 // gap
		}

	default:
		// Default positioning - relative to parent
		node.Bounds.X = 0
		node.Bounds.Y = 0
	}
}

// commitToDOM renders the tree using the abstract renderer
func (p *Pipeline) commitToDOM() {
	logger.Debug("RENDER", "Rendering with: %s", p.renderer.Name())
	
	if root := p.tree.GetRoot(); root != nil {
		// Convert tree to paint commands
		commands := ConvertNodeToCommands(root, 0, 0)
		
		if p.firstRender {
			// First render: full paint
			p.renderer.BeginFrame()
			for _, cmd := range commands {
				p.renderer.Paint(cmd)
			}
			p.renderer.EndFrame()
			
			p.previousCommands = commands
			p.firstRender = false
		} else {
			// Subsequent renders: let renderer decide how to handle updates
			updates := p.findUpdates(commands)
			
			if len(updates) > 0 {
				// Let renderer decide if it can handle selective updates
				if !p.renderer.ApplyUpdates(updates, commands) {
					// Renderer needs full redraw
					p.renderer.BeginFrame()
					for _, cmd := range commands {
						p.renderer.Paint(cmd)
					}
					p.renderer.EndFrame()
				}
			}
			
			p.previousCommands = commands
		}
	}
}

// findUpdates compares commands to find what changed
func (p *Pipeline) findUpdates(newCommands []PaintCommand) []PaintCommand {
	var updates []PaintCommand
	
	// Create map of previous commands for quick lookup
	prevMap := make(map[string]PaintCommand)
	for _, cmd := range p.previousCommands {
		prevMap[cmd.ID] = cmd
	}
	
	// Find changed commands
	for _, newCmd := range newCommands {
		if prevCmd, exists := prevMap[newCmd.ID]; exists {
			// Check if text changed
			if newCmd.Type == PaintText && newCmd.Text != prevCmd.Text {
				logger.Trace("UPDATE", "Text changed for %s from %s to %s", newCmd.ID, prevCmd.Text, newCmd.Text)
				newCmd.Type = UpdateText // Mark as update
				updates = append(updates, newCmd)
			}
			// Add more change detection as needed
		}
	}
	
	return updates
}

