package core

import (
	"iter"
	"sync"
	"sync/atomic"
)

// Tree represents the complete UI tree structure
type Tree struct {
	root      *Node
	nodeCount atomic.Int64
	version   atomic.Uint64
	mu        sync.RWMutex
	
	// Index for fast lookups (uses Swiss Tables internally in Go 1.24)
	nodeIndex map[NodeID]*Node
	
	// Stats for monitoring
	stats TreeStats
}

// TreeStats contains performance statistics
type TreeStats struct {
	TotalNodes      int64
	DirtyNodes      int64
	LastTraversalMs int64
	LastLayoutMs    int64
	LastPaintMs     int64
}

// NewTree creates a new UI tree
func NewTree() *Tree {
	return &Tree{
		nodeIndex: make(map[NodeID]*Node),
	}
}

// SetRoot sets the root node of the tree
func (t *Tree) SetRoot(root *Node) {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	t.root = root
	t.version.Add(1)
	t.rebuildIndex()
}

// GetRoot returns the root node
func (t *Tree) GetRoot() *Node {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.root
}

// NodeCount returns the total number of nodes
func (t *Tree) NodeCount() int64 {
	return t.nodeCount.Load()
}

// GetVersion returns the current tree version
func (t *Tree) GetVersion() uint64 {
	return t.version.Load()
}

// rebuildIndex rebuilds the node index
func (t *Tree) rebuildIndex() {
	t.nodeIndex = make(map[NodeID]*Node)
	count := int64(0)
	
	if t.root != nil {
		t.visitAll(t.root, func(n *Node) {
			t.nodeIndex[n.ID] = n
			count++
		})
	}
	
	t.nodeCount.Store(count)
}

// visitAll internal helper for traversing all nodes
func (t *Tree) visitAll(node *Node, visitor func(*Node)) {
	visitor(node)
	for _, child := range node.Children {
		t.visitAll(child, visitor)
	}
}

// FindNodeByID finds a node by its ID
func (t *Tree) FindNodeByID(id NodeID) *Node {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.nodeIndex[id]
}

// =============================================================================
// Go 1.24 Native Iterators using iter package
// =============================================================================

// DFS returns an iterator for depth-first traversal
func (t *Tree) DFS() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		t.mu.RLock()
		root := t.root
		t.mu.RUnlock()
		
		if root == nil {
			return
		}
		
		var traverse func(*Node) bool
		traverse = func(n *Node) bool {
			if !yield(n) {
				return false
			}
			for _, child := range n.Children {
				if !traverse(child) {
					return false
				}
			}
			return true
		}
		
		traverse(root)
	}
}

// BFS returns an iterator for breadth-first traversal
func (t *Tree) BFS() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		t.mu.RLock()
		root := t.root
		t.mu.RUnlock()
		
		if root == nil {
			return
		}
		
		queue := []*Node{root}
		
		for len(queue) > 0 {
			node := queue[0]
			queue = queue[1:]
			
			if !yield(node) {
				return // Early termination
			}
			
			queue = append(queue, node.Children...)
		}
	}
}

// PostOrderDFS returns an iterator for post-order depth-first traversal
func (t *Tree) PostOrderDFS() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		t.mu.RLock()
		root := t.root
		t.mu.RUnlock()
		
		if root == nil {
			return
		}
		
		var traverse func(*Node) bool
		traverse = func(n *Node) bool {
			for _, child := range n.Children {
				if !traverse(child) {
					return false
				}
			}
			return yield(n)
		}
		
		traverse(root)
	}
}

// PreOrderDFS returns an iterator for pre-order depth-first traversal
func (t *Tree) PreOrderDFS() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		// PreOrderDFS is the same as DFS
		for node := range t.DFS() {
			if !yield(node) {
				return
			}
		}
	}
}

// LevelOrder returns an iterator that yields nodes with their depth level
func (t *Tree) LevelOrder() iter.Seq2[int, *Node] {
	return func(yield func(int, *Node) bool) {
		t.mu.RLock()
		root := t.root
		t.mu.RUnlock()
		
		if root == nil {
			return
		}
		
		type nodeDepth struct {
			node  *Node
			depth int
		}
		
		queue := []nodeDepth{{root, 0}}
		
		for len(queue) > 0 {
			item := queue[0]
			queue = queue[1:]
			
			if !yield(item.depth, item.node) {
				return
			}
			
			for _, child := range item.node.Children {
				queue = append(queue, nodeDepth{child, item.depth + 1})
			}
		}
	}
}

// DirtyNodes returns an iterator for all dirty nodes
func (t *Tree) DirtyNodes() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for node := range t.DFS() {
			if node.IsDirty() {
				if !yield(node) {
					return
				}
			}
		}
	}
}

// Ancestors returns an iterator for a node's ancestors
func (n *Node) Ancestors() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		current := n.GetParent()
		for current != nil {
			if !yield(current) {
				return
			}
			current = current.GetParent()
		}
	}
}

// Descendants returns an iterator for all descendants
func (n *Node) Descendants() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		var traverse func(*Node) bool
		traverse = func(node *Node) bool {
			for _, child := range node.Children {
				if !yield(child) {
					return false
				}
				if !traverse(child) {
					return false
				}
			}
			return true
		}
		
		traverse(n)
	}
}

// Siblings returns an iterator for a node's siblings
func (n *Node) Siblings() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		parent := n.GetParent()
		if parent == nil {
			return
		}
		
		for _, sibling := range parent.Children {
			if sibling != n {
				if !yield(sibling) {
					return
				}
			}
		}
	}
}

// =============================================================================
// Parallel Traversal
// =============================================================================

// ParallelSubtrees processes subtrees in parallel
func (t *Tree) ParallelSubtrees(processor func(*Node) error) error {
	t.mu.RLock()
	root := t.root
	t.mu.RUnlock()
	
	if root == nil {
		return nil
	}
	
	// Process root's children in parallel
	var wg sync.WaitGroup
	errChan := make(chan error, len(root.Children))
	
	for _, child := range root.Children {
		wg.Add(1)
		go func(node *Node) {
			defer wg.Done()
			if err := t.processSubtree(node, processor); err != nil {
				errChan <- err
			}
		}(child)
	}
	
	wg.Wait()
	close(errChan)
	
	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	
	// Process root node
	return processor(root)
}

// processSubtree processes a subtree recursively
func (t *Tree) processSubtree(node *Node, processor func(*Node) error) error {
	// Process children first (post-order)
	for _, child := range node.Children {
		if err := t.processSubtree(child, processor); err != nil {
			return err
		}
	}
	
	return processor(node)
}

// =============================================================================
// Tree Manipulation
// =============================================================================

// InsertNode inserts a node at a specific position
func (t *Tree) InsertNode(parent *Node, child *Node, index int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if parent == nil || child == nil {
		return
	}
	
	// Ensure index is valid
	if index < 0 {
		index = 0
	}
	if index > len(parent.Children) {
		index = len(parent.Children)
	}
	
	// Insert at position
	parent.Children = append(parent.Children, nil)
	copy(parent.Children[index+1:], parent.Children[index:])
	parent.Children[index] = child
	child.SetParent(parent)
	
	// Update index
	t.nodeIndex[child.ID] = child
	t.nodeCount.Add(1)
	t.version.Add(1)
	
	// Mark parent as dirty
	parent.MarkDirty(ChildrenDirty)
}

// RemoveNode removes a node from the tree
func (t *Tree) RemoveNode(node *Node) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if node == nil || node == t.root {
		return false
	}
	
	parent := node.GetParent()
	if parent == nil {
		return false
	}
	
	// Remove from parent
	if parent.RemoveChild(node) {
		// Remove from index
		t.removeFromIndex(node)
		t.version.Add(1)
		return true
	}
	
	return false
}

// removeFromIndex removes a node and its descendants from the index
func (t *Tree) removeFromIndex(node *Node) {
	delete(t.nodeIndex, node.ID)
	t.nodeCount.Add(-1)
	
	for _, child := range node.Children {
		t.removeFromIndex(child)
	}
}

// GetStats returns current tree statistics
func (t *Tree) GetStats() TreeStats {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	stats := t.stats
	stats.TotalNodes = t.nodeCount.Load()
	
	// Count dirty nodes
	dirtyCount := int64(0)
	for node := range t.DirtyNodes() {
		dirtyCount++
		_ = node // Use node to avoid warning
	}
	stats.DirtyNodes = dirtyCount
	
	return stats
}