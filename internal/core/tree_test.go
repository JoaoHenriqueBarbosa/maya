package core

import (
	"testing"
)

// TestTreeCreation tests basic tree creation
func TestTreeCreation(t *testing.T) {
	tree := NewTree()
	if tree == nil {
		t.Fatal("NewTree returned nil")
	}
	
	if tree.NodeCount() != 0 {
		t.Errorf("Expected empty tree, got %d nodes", tree.NodeCount())
	}
}

// TestNodeCreation tests node creation with weak pointers and cleanup
func TestNodeCreation(t *testing.T) {
	// Create a simple widget for testing
	widget := &testWidget{}
	
	node := NewNode("test-node", widget)
	if node == nil {
		t.Fatal("NewNode returned nil")
	}
	
	if node.ID != "test-node" {
		t.Errorf("Expected ID 'test-node', got %s", node.ID)
	}
	
	if node.Widget != widget {
		t.Error("Widget not set correctly")
	}
}

// TestTreeTraversal tests the iterator-based traversal
func TestTreeTraversal(t *testing.T) {
	tree := NewTree()
	root := NewNode("root", &testWidget{})
	child1 := NewNode("child1", &testWidget{})
	child2 := NewNode("child2", &testWidget{})
	grandchild := NewNode("grandchild", &testWidget{})
	
	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild)
	
	tree.SetRoot(root)
	
	// Test DFS traversal
	var dfsNodes []string
	for node := range tree.DFS() {
		dfsNodes = append(dfsNodes, string(node.ID))
	}
	
	expectedDFS := []string{"root", "child1", "grandchild", "child2"}
	if !sliceEqual(dfsNodes, expectedDFS) {
		t.Errorf("DFS traversal incorrect. Got %v, expected %v", dfsNodes, expectedDFS)
	}
	
	// Test BFS traversal
	var bfsNodes []string
	for node := range tree.BFS() {
		bfsNodes = append(bfsNodes, string(node.ID))
	}
	
	expectedBFS := []string{"root", "child1", "child2", "grandchild"}
	if !sliceEqual(bfsNodes, expectedBFS) {
		t.Errorf("BFS traversal incorrect. Got %v, expected %v", bfsNodes, expectedBFS)
	}
}

// TestWeakParentReference tests that parent references are weak
func TestWeakParentReference(t *testing.T) {
	parent := NewNode("parent", &testWidget{})
	child := NewNode("child", &testWidget{})
	
	parent.AddChild(child)
	
	// Verify parent-child relationship
	if child.GetParent() != parent {
		t.Error("Parent not set correctly")
	}
	
	// Verify weak reference works
	if child.Parent != nil && child.Parent.Value() != parent {
		t.Error("Weak parent reference not working")
	}
}

// TestDirtyFlagPropagation tests dirty flag propagation
func TestDirtyFlagPropagation(t *testing.T) {
	tree := NewTree()
	root := NewNode("root", &testWidget{})
	child := NewNode("child", &testWidget{})
	grandchild := NewNode("grandchild", &testWidget{})
	
	root.AddChild(child)
	child.AddChild(grandchild)
	tree.SetRoot(root)
	
	// Mark grandchild as dirty
	grandchild.MarkDirty(LayoutDirty)
	
	if !grandchild.IsDirty() {
		t.Error("Grandchild should be dirty")
	}
	
	// Parent should also be marked dirty (with ChildrenDirty)
	if !child.IsDirty() {
		t.Error("Child should be dirty due to propagation")
	}
	
	if !root.IsDirty() {
		t.Error("Root should be dirty due to propagation")
	}
}

// TestB.Loop tests the new testing.B.Loop feature from Go 1.24
func BenchmarkTreeTraversal(b *testing.B) {
	// Setup tree once
	tree := buildLargeTree(100) // 100 nodes
	
	b.ResetTimer()
	
	// Use the new b.Loop() from Go 1.24
	for b.Loop() {
		count := 0
		for range tree.DFS() {
			count++
		}
		if count != 100 {
			b.Fatalf("Expected 100 nodes, got %d", count)
		}
	}
}

// Helper functions

type testWidget struct{}

func (w *testWidget) Layout(constraints Constraints) (width, height float64) {
	return 100, 100
}

func (w *testWidget) Paint(context PaintContext) {}

func (w *testWidget) HandleEvent(event Event) bool {
	return false
}

func (w *testWidget) GetIntrinsicWidth(height float64) float64 {
	return 100
}

func (w *testWidget) GetIntrinsicHeight(width float64) float64 {
	return 100
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildLargeTree(nodeCount int) *Tree {
	tree := NewTree()
	root := NewNode("root", &testWidget{})
	tree.SetRoot(root)
	
	// Build a tree with the specified number of nodes
	current := root
	for i := 1; i < nodeCount; i++ {
		node := NewNode(string(rune('a'+i%26)), &testWidget{})
		if i%3 == 0 && current != root {
			// Go back to root every 3 nodes to create branches
			current = root
		}
		current.AddChild(node)
		if i%2 == 0 {
			current = node // Make some nodes have children
		}
	}
	
	return tree
}