package graph

import (
	"context"
	"fmt"
	"iter"
	"sync"
	"sync/atomic"

)

// NodeID uniquely identifies a node in the graph
type NodeID string

// EdgeID uniquely identifies an edge in the graph
type EdgeID string

// Node represents a node in the dependency graph
type Node struct {
	ID       NodeID
	Data     interface{}
	Metadata map[string]interface{}
	
	// Dependencies
	InEdges  []EdgeID
	OutEdges []EdgeID
	
	// State tracking
	visited  atomic.Bool
	processing atomic.Bool
	processed atomic.Bool
	
	// Weak reference to graph
	graph *Graph
	
	mu sync.RWMutex
}

// Edge represents a directed edge between nodes
type Edge struct {
	ID       EdgeID
	From     NodeID
	To       NodeID
	Weight   float64
	Metadata map[string]interface{}
}

// Graph represents a directed acyclic graph for dependency management
type Graph struct {
	nodes map[NodeID]*Node
	edges map[EdgeID]*Edge
	
	// Index for fast lookups
	fromIndex map[NodeID][]EdgeID // Edges from node
	toIndex   map[NodeID][]EdgeID // Edges to node
	
	// Version for change tracking
	version atomic.Uint64
	
	mu sync.RWMutex
}

// NewGraph creates a new graph instance
func NewGraph() *Graph {
	return &Graph{
		nodes:     make(map[NodeID]*Node),
		edges:     make(map[EdgeID]*Edge),
		fromIndex: make(map[NodeID][]EdgeID),
		toIndex:   make(map[NodeID][]EdgeID),
	}
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(id NodeID, data interface{}) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if _, exists := g.nodes[id]; exists {
		return fmt.Errorf("node %s already exists", id)
	}
	
	node := &Node{
		ID:       id,
		Data:     data,
		Metadata: make(map[string]interface{}),
		InEdges:  []EdgeID{},
		OutEdges: []EdgeID{},
		graph:    g,
	}
	
	g.nodes[id] = node
	g.version.Add(1)
	
	return nil
}

// GetNode retrieves a node by ID
func (g *Graph) GetNode(id NodeID) (*Node, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	node, exists := g.nodes[id]
	return node, exists
}

// RemoveNode removes a node and all its edges
func (g *Graph) RemoveNode(id NodeID) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	node, exists := g.nodes[id]
	if !exists {
		return fmt.Errorf("node %s not found", id)
	}
	
	// Remove all edges connected to this node
	for _, edgeID := range node.InEdges {
		delete(g.edges, edgeID)
	}
	for _, edgeID := range node.OutEdges {
		delete(g.edges, edgeID)
	}
	
	// Update indices
	delete(g.fromIndex, id)
	delete(g.toIndex, id)
	
	// Remove node
	delete(g.nodes, id)
	g.version.Add(1)
	
	return nil
}

// AddEdge creates a directed edge from one node to another
func (g *Graph) AddEdge(from, to NodeID, weight float64) (EdgeID, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	// Verify nodes exist
	fromNode, fromExists := g.nodes[from]
	if !fromExists {
		return "", fmt.Errorf("source node %s not found", from)
	}
	
	toNode, toExists := g.nodes[to]
	if !toExists {
		return "", fmt.Errorf("target node %s not found", to)
	}
	
	// Generate edge ID
	edgeID := EdgeID(fmt.Sprintf("%s->%s", from, to))
	
	// Check if edge already exists
	if _, exists := g.edges[edgeID]; exists {
		return "", fmt.Errorf("edge %s already exists", edgeID)
	}
	
	// Create edge
	edge := &Edge{
		ID:       edgeID,
		From:     from,
		To:       to,
		Weight:   weight,
		Metadata: make(map[string]interface{}),
	}
	
	// Update graph
	g.edges[edgeID] = edge
	fromNode.OutEdges = append(fromNode.OutEdges, edgeID)
	toNode.InEdges = append(toNode.InEdges, edgeID)
	
	// Update indices
	g.fromIndex[from] = append(g.fromIndex[from], edgeID)
	g.toIndex[to] = append(g.toIndex[to], edgeID)
	
	g.version.Add(1)
	
	// Check for cycles
	if g.hasCycle() {
		// Rollback
		delete(g.edges, edgeID)
		fromNode.OutEdges = fromNode.OutEdges[:len(fromNode.OutEdges)-1]
		toNode.InEdges = toNode.InEdges[:len(toNode.InEdges)-1]
		g.fromIndex[from] = g.fromIndex[from][:len(g.fromIndex[from])-1]
		g.toIndex[to] = g.toIndex[to][:len(g.toIndex[to])-1]
		
		return "", fmt.Errorf("adding edge would create a cycle")
	}
	
	return edgeID, nil
}

// RemoveEdge removes an edge from the graph
func (g *Graph) RemoveEdge(edgeID EdgeID) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	edge, exists := g.edges[edgeID]
	if !exists {
		return fmt.Errorf("edge %s not found", edgeID)
	}
	
	// Update nodes
	if fromNode, ok := g.nodes[edge.From]; ok {
		fromNode.OutEdges = removeEdgeID(fromNode.OutEdges, edgeID)
	}
	if toNode, ok := g.nodes[edge.To]; ok {
		toNode.InEdges = removeEdgeID(toNode.InEdges, edgeID)
	}
	
	// Update indices
	g.fromIndex[edge.From] = removeEdgeID(g.fromIndex[edge.From], edgeID)
	g.toIndex[edge.To] = removeEdgeID(g.toIndex[edge.To], edgeID)
	
	// Remove edge
	delete(g.edges, edgeID)
	g.version.Add(1)
	
	return nil
}

// GetDependencies returns nodes that the given node depends on
func (g *Graph) GetDependencies(nodeID NodeID) []NodeID {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	node, exists := g.nodes[nodeID]
	if !exists {
		return nil
	}
	
	dependencies := make([]NodeID, 0, len(node.InEdges))
	for _, edgeID := range node.InEdges {
		if edge, ok := g.edges[edgeID]; ok {
			dependencies = append(dependencies, edge.From)
		}
	}
	
	return dependencies
}

// GetDependents returns nodes that depend on the given node
func (g *Graph) GetDependents(nodeID NodeID) []NodeID {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	node, exists := g.nodes[nodeID]
	if !exists {
		return nil
	}
	
	dependents := make([]NodeID, 0, len(node.OutEdges))
	for _, edgeID := range node.OutEdges {
		if edge, ok := g.edges[edgeID]; ok {
			dependents = append(dependents, edge.To)
		}
	}
	
	return dependents
}

// TopologicalSort returns nodes in topological order
func (g *Graph) TopologicalSort() ([]NodeID, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	// Count in-degrees
	inDegrees := make(map[NodeID]int)
	for nodeID, node := range g.nodes {
		inDegrees[nodeID] = len(node.InEdges)
	}
	
	// Find nodes with no dependencies
	queue := make([]NodeID, 0)
	for nodeID, degree := range inDegrees {
		if degree == 0 {
			queue = append(queue, nodeID)
		}
	}
	
	result := make([]NodeID, 0, len(g.nodes))
	
	// Process queue
	for len(queue) > 0 {
		// Dequeue
		nodeID := queue[0]
		queue = queue[1:]
		result = append(result, nodeID)
		
		// Process dependents
		node := g.nodes[nodeID]
		for _, edgeID := range node.OutEdges {
			if edge, ok := g.edges[edgeID]; ok {
				inDegrees[edge.To]--
				if inDegrees[edge.To] == 0 {
					queue = append(queue, edge.To)
				}
			}
		}
	}
	
	// Check if all nodes were processed
	if len(result) != len(g.nodes) {
		return nil, fmt.Errorf("graph contains a cycle")
	}
	
	return result, nil
}

// DFS performs depth-first search traversal
func (g *Graph) DFS() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		g.mu.RLock()
		defer g.mu.RUnlock()
		
		visited := make(map[NodeID]bool)
		
		var dfsRecursive func(NodeID) bool
		dfsRecursive = func(nodeID NodeID) bool {
			if visited[nodeID] {
				return true
			}
			visited[nodeID] = true
			
			node, exists := g.nodes[nodeID]
			if !exists {
				return true
			}
			
			if !yield(node) {
				return false
			}
			
			// Visit dependents
			for _, edgeID := range node.OutEdges {
				if edge, ok := g.edges[edgeID]; ok {
					if !dfsRecursive(edge.To) {
						return false
					}
				}
			}
			
			return true
		}
		
		// Start from nodes with no dependencies
		for nodeID, node := range g.nodes {
			if len(node.InEdges) == 0 {
				if !dfsRecursive(nodeID) {
					return
				}
			}
		}
		
		// Visit any remaining unvisited nodes (disconnected components)
		for nodeID := range g.nodes {
			if !visited[nodeID] {
				if !dfsRecursive(nodeID) {
					return
				}
			}
		}
	}
}

// BFS performs breadth-first search traversal
func (g *Graph) BFS() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		g.mu.RLock()
		defer g.mu.RUnlock()
		
		visited := make(map[NodeID]bool)
		queue := make([]NodeID, 0)
		
		// Start from nodes with no dependencies
		for nodeID, node := range g.nodes {
			if len(node.InEdges) == 0 {
				queue = append(queue, nodeID)
				visited[nodeID] = true
			}
		}
		
		// Process queue
		for len(queue) > 0 {
			nodeID := queue[0]
			queue = queue[1:]
			
			node, exists := g.nodes[nodeID]
			if !exists {
				continue
			}
			
			if !yield(node) {
				return
			}
			
			// Add unvisited dependents to queue
			for _, edgeID := range node.OutEdges {
				if edge, ok := g.edges[edgeID]; ok {
					if !visited[edge.To] {
						visited[edge.To] = true
						queue = append(queue, edge.To)
					}
				}
			}
		}
		
		// Process any remaining unvisited nodes (disconnected components)
		for nodeID, node := range g.nodes {
			if !visited[nodeID] {
				visited[nodeID] = true
				if !yield(node) {
					return
				}
			}
		}
	}
}

// ParallelProcess processes independent nodes in parallel
func (g *Graph) ParallelProcess(ctx context.Context, processor func(context.Context, *Node) error) error {
	// Get topological order
	order, err := g.TopologicalSort()
	if err != nil {
		return err
	}
	
	// Group nodes by level (nodes that can be processed in parallel)
	levels := g.groupByLevel(order)
	
	// Process each level
	for _, level := range levels {
		if err := g.processLevel(ctx, level, processor); err != nil {
			return err
		}
	}
	
	return nil
}

// groupByLevel groups nodes that can be processed in parallel
func (g *Graph) groupByLevel(order []NodeID) [][]NodeID {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	levels := make([][]NodeID, 0)
	nodeLevel := make(map[NodeID]int)
	
	// Calculate level for each node
	for _, nodeID := range order {
		level := 0
		node := g.nodes[nodeID]
		
		// Find max level of dependencies
		for _, edgeID := range node.InEdges {
			if edge, ok := g.edges[edgeID]; ok {
				if depLevel, exists := nodeLevel[edge.From]; exists {
					if depLevel >= level {
						level = depLevel + 1
					}
				}
			}
		}
		
		nodeLevel[nodeID] = level
		
		// Ensure we have enough levels
		for len(levels) <= level {
			levels = append(levels, []NodeID{})
		}
		
		levels[level] = append(levels[level], nodeID)
	}
	
	return levels
}

// processLevel processes all nodes in a level in parallel
func (g *Graph) processLevel(ctx context.Context, level []NodeID, processor func(context.Context, *Node) error) error {
	if len(level) == 0 {
		return nil
	}
	
	var wg sync.WaitGroup
	errChan := make(chan error, len(level))
	
	for _, nodeID := range level {
		node, exists := g.nodes[nodeID]
		if !exists {
			continue
		}
		
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			
			if err := processor(ctx, n); err != nil {
				select {
				case errChan <- fmt.Errorf("node %s: %w", n.ID, err):
				case <-ctx.Done():
				}
			}
		}(node)
	}
	
	// Wait for completion
	doneChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneChan)
	}()
	
	select {
	case <-doneChan:
		close(errChan)
		for err := range errChan {
			if err != nil {
				return err
			}
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// hasCycle detects if the graph contains a cycle using DFS
func (g *Graph) hasCycle() bool {
	visited := make(map[NodeID]bool)
	recStack := make(map[NodeID]bool)
	
	var hasCycleDFS func(NodeID) bool
	hasCycleDFS = func(nodeID NodeID) bool {
		visited[nodeID] = true
		recStack[nodeID] = true
		
		node := g.nodes[nodeID]
		for _, edgeID := range node.OutEdges {
			if edge, ok := g.edges[edgeID]; ok {
				if !visited[edge.To] {
					if hasCycleDFS(edge.To) {
						return true
					}
				} else if recStack[edge.To] {
					return true
				}
			}
		}
		
		recStack[nodeID] = false
		return false
	}
	
	for nodeID := range g.nodes {
		if !visited[nodeID] {
			if hasCycleDFS(nodeID) {
				return true
			}
		}
	}
	
	return false
}

// Clone creates a deep copy of the graph
func (g *Graph) Clone() *Graph {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	newGraph := NewGraph()
	
	// Clone nodes
	for nodeID, node := range g.nodes {
		newNode := &Node{
			ID:       node.ID,
			Data:     node.Data,
			Metadata: make(map[string]interface{}),
			InEdges:  make([]EdgeID, len(node.InEdges)),
			OutEdges: make([]EdgeID, len(node.OutEdges)),
			graph:    newGraph,
		}
		
		// Copy metadata
		for k, v := range node.Metadata {
			newNode.Metadata[k] = v
		}
		
		// Copy edge IDs
		copy(newNode.InEdges, node.InEdges)
		copy(newNode.OutEdges, node.OutEdges)
		
		newGraph.nodes[nodeID] = newNode
	}
	
	// Clone edges
	for edgeID, edge := range g.edges {
		newEdge := &Edge{
			ID:       edge.ID,
			From:     edge.From,
			To:       edge.To,
			Weight:   edge.Weight,
			Metadata: make(map[string]interface{}),
		}
		
		// Copy metadata
		for k, v := range edge.Metadata {
			newEdge.Metadata[k] = v
		}
		
		newGraph.edges[edgeID] = newEdge
	}
	
	// Clone indices
	for k, v := range g.fromIndex {
		newGraph.fromIndex[k] = append([]EdgeID{}, v...)
	}
	for k, v := range g.toIndex {
		newGraph.toIndex[k] = append([]EdgeID{}, v...)
	}
	
	return newGraph
}

// NodeCount returns the number of nodes in the graph
func (g *Graph) NodeCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.nodes)
}

// EdgeCount returns the number of edges in the graph
func (g *Graph) EdgeCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.edges)
}

// IsDAG checks if the graph is a directed acyclic graph
func (g *Graph) IsDAG() bool {
	return !g.hasCycle()
}

// Clear removes all nodes and edges from the graph
func (g *Graph) Clear() {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	g.nodes = make(map[NodeID]*Node)
	g.edges = make(map[EdgeID]*Edge)
	g.fromIndex = make(map[NodeID][]EdgeID)
	g.toIndex = make(map[NodeID][]EdgeID)
	g.version.Add(1)
}

// Helper function to remove an EdgeID from a slice
func removeEdgeID(slice []EdgeID, id EdgeID) []EdgeID {
	for i, edgeID := range slice {
		if edgeID == id {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// MultipassProcessor handles multi-pass processing of the graph
type MultipassProcessor struct {
	graph     *Graph
	passes    []PassFunc
	results   map[NodeID]interface{}
	mu        sync.RWMutex
}

// PassFunc defines a processing pass function
type PassFunc func(context.Context, *Node, map[NodeID]interface{}) (interface{}, error)

// NewMultipassProcessor creates a new multipass processor
func NewMultipassProcessor(graph *Graph) *MultipassProcessor {
	return &MultipassProcessor{
		graph:   graph,
		passes:  make([]PassFunc, 0),
		results: make(map[NodeID]interface{}),
	}
}

// AddPass adds a processing pass
func (p *MultipassProcessor) AddPass(pass PassFunc) {
	p.passes = append(p.passes, pass)
}

// Execute runs all passes on the graph
func (p *MultipassProcessor) Execute(ctx context.Context) error {
	for i, pass := range p.passes {
		if err := p.executePass(ctx, pass, i); err != nil {
			return fmt.Errorf("pass %d failed: %w", i, err)
		}
	}
	return nil
}

// executePass runs a single pass on all nodes
func (p *MultipassProcessor) executePass(ctx context.Context, pass PassFunc, passIndex int) error {
	// Get topological order for this pass
	order, err := p.graph.TopologicalSort()
	if err != nil {
		return err
	}
	
	// Group nodes by level for parallel processing
	levels := p.graph.groupByLevel(order)
	
	// Process each level
	for _, level := range levels {
		if err := p.processPassLevel(ctx, level, pass); err != nil {
			return err
		}
	}
	
	return nil
}

// processPassLevel processes a level of nodes for a pass
func (p *MultipassProcessor) processPassLevel(ctx context.Context, level []NodeID, pass PassFunc) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(level))
	
	for _, nodeID := range level {
		node, exists := p.graph.nodes[nodeID]
		if !exists {
			continue
		}
		
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			
			// Get current results snapshot
			p.mu.RLock()
			resultsCopy := make(map[NodeID]interface{})
			for k, v := range p.results {
				resultsCopy[k] = v
			}
			p.mu.RUnlock()
			
			// Execute pass
			result, err := pass(ctx, n, resultsCopy)
			if err != nil {
				select {
				case errChan <- err:
				case <-ctx.Done():
				}
				return
			}
			
			// Store result
			p.mu.Lock()
			p.results[n.ID] = result
			p.mu.Unlock()
		}(node)
	}
	
	// Wait for completion
	doneChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneChan)
	}()
	
	select {
	case <-doneChan:
		close(errChan)
		for err := range errChan {
			if err != nil {
				return err
			}
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// GetResult returns the result for a node
func (p *MultipassProcessor) GetResult(nodeID NodeID) (interface{}, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	result, exists := p.results[nodeID]
	return result, exists
}

// GetAllResults returns all results
func (p *MultipassProcessor) GetAllResults() map[NodeID]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	results := make(map[NodeID]interface{})
	for k, v := range p.results {
		results[k] = v
	}
	return results
}