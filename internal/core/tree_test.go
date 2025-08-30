package core

import (
	"errors"
	"sync"
	"testing"
)

// TestTree_Creation tests tree creation and initialization
func TestTree_Creation(t *testing.T) {
	tree := NewTree()
	
	if tree == nil {
		t.Fatal("NewTree returned nil")
	}
	
	if tree.NodeCount() != 0 {
		t.Errorf("New tree should have 0 nodes, got %d", tree.NodeCount())
	}
	
	if tree.GetRoot() != nil {
		t.Error("New tree should have nil root")
	}
	
	if tree.GetVersion() != 0 {
		t.Error("New tree should have version 0")
	}
}

// TestTree_RootManagement tests setting and getting root
func TestTree_RootManagement(t *testing.T) {
	t.Run("set_root", func(t *testing.T) {
		tree := NewTree()
		root := NewNode("root", &mockWidget{})
		
		initialVersion := tree.GetVersion()
		tree.SetRoot(root)
		
		if tree.GetRoot() != root {
			t.Error("Root not set correctly")
		}
		
		if tree.NodeCount() != 1 {
			t.Errorf("Tree should have 1 node, got %d", tree.NodeCount())
		}
		
		if tree.GetVersion() <= initialVersion {
			t.Error("Version should increment after setting root")
		}
	})
	
	t.Run("set_root_with_children", func(t *testing.T) {
		tree := NewTree()
		root := NewNode("root", &mockWidget{})
		child1 := NewNode("child1", &mockWidget{})
		child2 := NewNode("child2", &mockWidget{})
		grandchild := NewNode("grandchild", &mockWidget{})
		
		root.AddChild(child1)
		root.AddChild(child2)
		child1.AddChild(grandchild)
		
		tree.SetRoot(root)
		
		if tree.NodeCount() != 4 {
			t.Errorf("Tree should have 4 nodes, got %d", tree.NodeCount())
		}
		
		// Verify all nodes are in index
		if tree.FindNodeByID("root") != root {
			t.Error("Root not found in index")
		}
		if tree.FindNodeByID("child1") != child1 {
			t.Error("Child1 not found in index")
		}
		if tree.FindNodeByID("child2") != child2 {
			t.Error("Child2 not found in index")
		}
		if tree.FindNodeByID("grandchild") != grandchild {
			t.Error("Grandchild not found in index")
		}
	})
	
	t.Run("replace_root", func(t *testing.T) {
		tree := NewTree()
		oldRoot := NewNode("old-root", &mockWidget{})
		oldChild := NewNode("old-child", &mockWidget{})
		oldRoot.AddChild(oldChild)
		
		tree.SetRoot(oldRoot)
		oldCount := tree.NodeCount()
		
		newRoot := NewNode("new-root", &mockWidget{})
		tree.SetRoot(newRoot)
		
		if tree.GetRoot() != newRoot {
			t.Error("Root not replaced correctly")
		}
		
		if tree.NodeCount() != 1 {
			t.Errorf("Tree should have 1 node after replacement, got %d", tree.NodeCount())
		}
		
		if tree.NodeCount() >= oldCount {
			t.Error("Node count should decrease after replacing root with smaller tree")
		}
		
		// Old nodes should not be in index
		if tree.FindNodeByID("old-root") != nil {
			t.Error("Old root should not be in index")
		}
		if tree.FindNodeByID("old-child") != nil {
			t.Error("Old child should not be in index")
		}
	})
}

// TestTree_Iterators tests all iterator implementations
func TestTree_Iterators(t *testing.T) {
	// Build test tree:
	//       root
	//      /    \
	//   child1  child2
	//     |       |
	//   gc1      gc2
	
	buildTestTree := func() *Tree {
		tree := NewTree()
		root := NewNode("root", &mockWidget{})
		child1 := NewNode("child1", &mockWidget{})
		child2 := NewNode("child2", &mockWidget{})
		gc1 := NewNode("gc1", &mockWidget{})
		gc2 := NewNode("gc2", &mockWidget{})
		
		root.AddChild(child1)
		root.AddChild(child2)
		child1.AddChild(gc1)
		child2.AddChild(gc2)
		
		tree.SetRoot(root)
		return tree
	}
	
	t.Run("DFS_iterator", func(t *testing.T) {
		tree := buildTestTree()
		
		var visited []string
		for node := range tree.DFS() {
			visited = append(visited, string(node.ID))
		}
		
		expected := []string{"root", "child1", "gc1", "child2", "gc2"}
		if !sliceEqual(visited, expected) {
			t.Errorf("DFS order incorrect. Got %v, expected %v", visited, expected)
		}
	})
	
	t.Run("BFS_iterator", func(t *testing.T) {
		tree := buildTestTree()
		
		var visited []string
		for node := range tree.BFS() {
			visited = append(visited, string(node.ID))
		}
		
		expected := []string{"root", "child1", "child2", "gc1", "gc2"}
		if !sliceEqual(visited, expected) {
			t.Errorf("BFS order incorrect. Got %v, expected %v", visited, expected)
		}
	})
	
	t.Run("PostOrderDFS_iterator", func(t *testing.T) {
		tree := buildTestTree()
		
		var visited []string
		for node := range tree.PostOrderDFS() {
			visited = append(visited, string(node.ID))
		}
		
		expected := []string{"gc1", "child1", "gc2", "child2", "root"}
		if !sliceEqual(visited, expected) {
			t.Errorf("PostOrderDFS order incorrect. Got %v, expected %v", visited, expected)
		}
	})
	
	t.Run("PreOrderDFS_iterator", func(t *testing.T) {
		tree := buildTestTree()
		
		var visited []string
		for node := range tree.PreOrderDFS() {
			visited = append(visited, string(node.ID))
		}
		
		expected := []string{"root", "child1", "gc1", "child2", "gc2"}
		if !sliceEqual(visited, expected) {
			t.Errorf("PreOrderDFS order incorrect. Got %v, expected %v", visited, expected)
		}
	})
	
	t.Run("LevelOrder_iterator", func(t *testing.T) {
		tree := buildTestTree()
		
		type levelNode struct {
			level int
			id    string
		}
		
		var visited []levelNode
		for level, node := range tree.LevelOrder() {
			visited = append(visited, levelNode{level, string(node.ID)})
		}
		
		expected := []levelNode{
			{0, "root"},
			{1, "child1"},
			{1, "child2"},
			{2, "gc1"},
			{2, "gc2"},
		}
		
		if len(visited) != len(expected) {
			t.Fatalf("LevelOrder count incorrect. Got %d, expected %d", len(visited), len(expected))
		}
		
		for i, v := range visited {
			if v.level != expected[i].level || v.id != expected[i].id {
				t.Errorf("LevelOrder[%d] = {%d, %s}, expected {%d, %s}",
					i, v.level, v.id, expected[i].level, expected[i].id)
			}
		}
	})
	
	t.Run("DirtyNodes_iterator", func(t *testing.T) {
		tree := buildTestTree()
		
		// Mark some nodes dirty
		tree.GetRoot().Children[0].MarkDirty(LayoutDirty) // child1
		tree.GetRoot().Children[1].Children[0].MarkDirty(PaintDirty) // gc2
		
		var dirtyNodes []string
		for node := range tree.DirtyNodes() {
			if node.IsDirty() {
				dirtyNodes = append(dirtyNodes, string(node.ID))
			}
		}
		
		// Due to propagation, root should also be dirty
		if len(dirtyNodes) < 3 {
			t.Errorf("Expected at least 3 dirty nodes, got %d", len(dirtyNodes))
		}
	})
	
	t.Run("early_termination", func(t *testing.T) {
		tree := buildTestTree()
		
		count := 0
		for node := range tree.DFS() {
			count++
			if string(node.ID) == "child1" {
				break // Early termination
			}
		}
		
		if count != 2 { // root and child1
			t.Errorf("Early termination failed, processed %d nodes", count)
		}
	})
	
	t.Run("empty_tree_iterators", func(t *testing.T) {
		tree := NewTree()
		
		count := 0
		for range tree.DFS() {
			count++
		}
		if count != 0 {
			t.Error("Empty tree DFS should not iterate")
		}
		
		for range tree.BFS() {
			count++
		}
		if count != 0 {
			t.Error("Empty tree BFS should not iterate")
		}
		
		for range tree.PostOrderDFS() {
			count++
		}
		if count != 0 {
			t.Error("Empty tree PostOrderDFS should not iterate")
		}
		
		for range tree.LevelOrder() {
			count++
		}
		if count != 0 {
			t.Error("Empty tree LevelOrder should not iterate")
		}
	})
}

// TestNode_Iterators tests node-specific iterators
func TestNode_Iterators(t *testing.T) {
	t.Run("Ancestors_iterator", func(t *testing.T) {
		root := NewNode("root", &mockWidget{})
		parent := NewNode("parent", &mockWidget{})
		node := NewNode("node", &mockWidget{})
		
		root.AddChild(parent)
		parent.AddChild(node)
		
		var ancestors []string
		for ancestor := range node.Ancestors() {
			ancestors = append(ancestors, string(ancestor.ID))
		}
		
		expected := []string{"parent", "root"}
		if !sliceEqual(ancestors, expected) {
			t.Errorf("Ancestors incorrect. Got %v, expected %v", ancestors, expected)
		}
	})
	
	t.Run("Descendants_iterator", func(t *testing.T) {
		node := NewNode("node", &mockWidget{})
		child1 := NewNode("child1", &mockWidget{})
		child2 := NewNode("child2", &mockWidget{})
		grandchild := NewNode("grandchild", &mockWidget{})
		
		node.AddChild(child1)
		node.AddChild(child2)
		child1.AddChild(grandchild)
		
		var descendants []string
		for desc := range node.Descendants() {
			descendants = append(descendants, string(desc.ID))
		}
		
		expected := []string{"child1", "grandchild", "child2"}
		if !sliceEqual(descendants, expected) {
			t.Errorf("Descendants incorrect. Got %v, expected %v", descendants, expected)
		}
	})
	
	t.Run("Siblings_iterator", func(t *testing.T) {
		parent := NewNode("parent", &mockWidget{})
		node := NewNode("node", &mockWidget{})
		sibling1 := NewNode("sibling1", &mockWidget{})
		sibling2 := NewNode("sibling2", &mockWidget{})
		
		parent.AddChild(sibling1)
		parent.AddChild(node)
		parent.AddChild(sibling2)
		
		var siblings []string
		for sib := range node.Siblings() {
			siblings = append(siblings, string(sib.ID))
		}
		
		expected := []string{"sibling1", "sibling2"}
		if !sliceEqual(siblings, expected) {
			t.Errorf("Siblings incorrect. Got %v, expected %v", siblings, expected)
		}
	})
	
	t.Run("node_without_parent", func(t *testing.T) {
		node := NewNode("orphan", &mockWidget{})
		
		count := 0
		for range node.Ancestors() {
			count++
		}
		
		if count != 0 {
			t.Error("Orphan node should have no ancestors")
		}
		
		for range node.Siblings() {
			count++
		}
		
		if count != 0 {
			t.Error("Orphan node should have no siblings")
		}
	})
}

// TestTree_Manipulation tests tree manipulation operations
func TestTree_Manipulation(t *testing.T) {
	t.Run("insert_node", func(t *testing.T) {
		tree := NewTree()
		root := NewNode("root", &mockWidget{})
		child1 := NewNode("child1", &mockWidget{})
		child2 := NewNode("child2", &mockWidget{})
		
		tree.SetRoot(root)
		tree.InsertNode(root, child1, 0)
		tree.InsertNode(root, child2, 1)
		
		newChild := NewNode("inserted", &mockWidget{})
		
		// Insert at index 1 (between child1 and child2)
		tree.InsertNode(root, newChild, 1)
		
		if len(root.Children) != 3 {
			t.Errorf("Expected 3 children, got %d", len(root.Children))
		}
		
		if root.Children[1] != newChild {
			t.Error("Node not inserted at correct position")
		}
		
		if newChild.GetParent() != root {
			t.Error("Inserted node parent not set")
		}
		
		if tree.FindNodeByID("inserted") != newChild {
			t.Error("Inserted node not in index")
		}
		
		if tree.NodeCount() != 4 {
			t.Errorf("Tree should have 4 nodes, got %d", tree.NodeCount())
		}
	})
	
	t.Run("insert_node_edge_cases", func(t *testing.T) {
		tree := NewTree()
		root := NewNode("root", &mockWidget{})
		tree.SetRoot(root)
		
		child := NewNode("child", &mockWidget{})
		
		// Insert at negative index (should insert at 0)
		tree.InsertNode(root, child, -5)
		if root.Children[0] != child {
			t.Error("Node with negative index should insert at beginning")
		}
		
		// Insert at index beyond bounds
		child2 := NewNode("child2", &mockWidget{})
		tree.InsertNode(root, child2, 100)
		if root.Children[len(root.Children)-1] != child2 {
			t.Error("Node with large index should insert at end")
		}
		
		// Insert with nil parent (should do nothing)
		child3 := NewNode("child3", &mockWidget{})
		oldCount := tree.NodeCount()
		tree.InsertNode(nil, child3, 0)
		if tree.NodeCount() != oldCount {
			t.Error("Insert with nil parent should do nothing")
		}
		
		// Insert nil child (should do nothing)
		tree.InsertNode(root, nil, 0)
		if len(root.Children) != 2 {
			t.Error("Insert nil child should do nothing")
		}
	})
	
	t.Run("remove_node", func(t *testing.T) {
		tree := NewTree()
		root := NewNode("root", &mockWidget{})
		child := NewNode("child", &mockWidget{})
		grandchild := NewNode("grandchild", &mockWidget{})
		
		root.AddChild(child)
		child.AddChild(grandchild)
		tree.SetRoot(root)
		
		// Remove child (and its subtree)
		removed := tree.RemoveNode(child)
		
		if !removed {
			t.Error("RemoveNode should return true for existing node")
		}
		
		if len(root.Children) != 0 {
			t.Error("Child should be removed from parent")
		}
		
		if tree.FindNodeByID("child") != nil {
			t.Error("Removed node should not be in index")
		}
		
		if tree.FindNodeByID("grandchild") != nil {
			t.Error("Descendant of removed node should not be in index")
		}
		
		if tree.NodeCount() != 1 {
			t.Errorf("Tree should have 1 node after removal, got %d", tree.NodeCount())
		}
	})
	
	t.Run("remove_node_edge_cases", func(t *testing.T) {
		tree := NewTree()
		root := NewNode("root", &mockWidget{})
		tree.SetRoot(root)
		
		// Try to remove root (should fail)
		removed := tree.RemoveNode(root)
		if removed {
			t.Error("Should not be able to remove root")
		}
		
		// Try to remove nil (should fail)
		removed = tree.RemoveNode(nil)
		if removed {
			t.Error("Should not be able to remove nil")
		}
		
		// Try to remove orphan node (should fail)
		orphan := NewNode("orphan", &mockWidget{})
		removed = tree.RemoveNode(orphan)
		if removed {
			t.Error("Should not be able to remove orphan node")
		}
	})
}

// TestTree_ParallelProcessing tests parallel subtree processing
func TestTree_ParallelProcessing(t *testing.T) {
	t.Run("parallel_subtrees", func(t *testing.T) {
		tree := NewTree()
		root := NewNode("root", &mockWidget{})
		
		// Create multiple subtrees
		for i := 0; i < 5; i++ {
			child := NewNode(string(rune('a'+i)), &mockWidget{})
			for j := 0; j < 3; j++ {
				// Use unique IDs for grandchildren
				grandchildID := string(rune('a'+i)) + "-" + string(rune('0'+j))
				grandchild := NewNode(grandchildID, &mockWidget{})
				child.AddChild(grandchild)
			}
			root.AddChild(child)
		}
		
		tree.SetRoot(root)
		
		processedNodes := make(map[NodeID]bool)
		var mu sync.Mutex
		
		err := tree.ParallelSubtrees(func(node *Node) error {
			mu.Lock()
			processedNodes[node.ID] = true
			mu.Unlock()
			return nil
		})
		
		if err != nil {
			t.Errorf("ParallelSubtrees returned error: %v", err)
		}
		
		if len(processedNodes) != int(tree.NodeCount()) {
			t.Errorf("Not all nodes processed. Got %d, expected %d",
				len(processedNodes), tree.NodeCount())
		}
	})
	
	t.Run("parallel_with_error", func(t *testing.T) {
		tree := NewTree()
		root := NewNode("root", &mockWidget{})
		child1 := NewNode("child1", &mockWidget{})
		child2 := NewNode("error-node", &mockWidget{})
		
		root.AddChild(child1)
		root.AddChild(child2)
		tree.SetRoot(root)
		
		err := tree.ParallelSubtrees(func(node *Node) error {
			if node.ID == "error-node" {
				return errors.New("test error")
			}
			return nil
		})
		
		if err == nil {
			t.Error("ParallelSubtrees should return error when processor fails")
		}
	})
	
	t.Run("parallel_empty_tree", func(t *testing.T) {
		tree := NewTree()
		
		err := tree.ParallelSubtrees(func(node *Node) error {
			t.Error("Should not process any nodes for empty tree")
			return nil
		})
		
		if err != nil {
			t.Error("ParallelSubtrees should not error on empty tree")
		}
	})
}

// TestTree_Stats tests statistics gathering
func TestTree_Stats(t *testing.T) {
	tree := NewTree()
	root := NewNode("root", &mockWidget{})
	child1 := NewNode("child1", &mockWidget{})
	child2 := NewNode("child2", &mockWidget{})
	
	root.AddChild(child1)
	root.AddChild(child2)
	tree.SetRoot(root)
	
	// Mark some nodes dirty
	child1.MarkDirty(LayoutDirty)
	
	stats := tree.GetStats()
	
	if stats.TotalNodes != 3 {
		t.Errorf("TotalNodes = %d, expected 3", stats.TotalNodes)
	}
	
	// Due to propagation, at least 2 nodes should be dirty
	if stats.DirtyNodes < 2 {
		t.Errorf("DirtyNodes = %d, expected at least 2", stats.DirtyNodes)
	}
}

// TestTree_ConcurrentAccess tests thread safety
func TestTree_ConcurrentAccess(t *testing.T) {
	tree := NewTree()
	root := NewNode("root", &mockWidget{})
	tree.SetRoot(root)
	
	// Add initial children
	for i := 0; i < 10; i++ {
		child := NewNode(string(rune('a'+i)), &mockWidget{})
		root.AddChild(child)
	}
	tree.SetRoot(root) // Rebuild index
	
	done := make(chan bool, 4)
	
	// Concurrent reads
	go func() {
		for i := 0; i < 100; i++ {
			_ = tree.GetRoot()
			_ = tree.NodeCount()
			_ = tree.GetVersion()
			for range tree.DFS() {
				break
			}
		}
		done <- true
	}()
	
	// Concurrent finds
	go func() {
		for i := 0; i < 100; i++ {
			_ = tree.FindNodeByID(NodeID(string(rune('a'+i%10))))
		}
		done <- true
	}()
	
	// Concurrent stats
	go func() {
		for i := 0; i < 100; i++ {
			_ = tree.GetStats()
		}
		done <- true
	}()
	
	// Concurrent modifications
	go func() {
		for i := 0; i < 10; i++ {
			child := NewNode(string(rune('z'-i)), &mockWidget{})
			tree.InsertNode(root, child, 0)
		}
		done <- true
	}()
	
	// Wait for all goroutines
	for i := 0; i < 4; i++ {
		<-done
	}
	
	// Tree should still be in valid state
	if tree.GetRoot() == nil {
		t.Error("Tree root became nil during concurrent access")
	}
	
	if tree.NodeCount() < 10 {
		t.Error("Tree lost nodes during concurrent access")
	}
}

// Benchmark tests using Go 1.24's testing.B.Loop
func BenchmarkTree_DFSTraversal(b *testing.B) {
	tree := buildBenchmarkTree(100)
	
	b.ResetTimer()
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

func BenchmarkTree_BFSTraversal(b *testing.B) {
	tree := buildBenchmarkTree(100)
	
	b.ResetTimer()
	for b.Loop() {
		count := 0
		for range tree.BFS() {
			count++
		}
		if count != 100 {
			b.Fatalf("Expected 100 nodes, got %d", count)
		}
	}
}

func BenchmarkTree_FindNodeByID(b *testing.B) {
	tree := buildBenchmarkTree(1000)
	
	b.ResetTimer()
	for b.Loop() {
		node := tree.FindNodeByID("node-500")
		if node == nil {
			b.Fatal("Node not found")
		}
	}
}

func BenchmarkTree_InsertNode(b *testing.B) {
	tree := buildBenchmarkTree(10)
	root := tree.GetRoot()
	
	b.ResetTimer()
	for b.Loop() {
		node := NewNode("bench", &mockWidget{})
		tree.InsertNode(root, node, 0)
		tree.RemoveNode(node) // Clean up for next iteration
	}
}

// Helper functions
func buildBenchmarkTree(nodeCount int) *Tree {
	tree := NewTree()
	root := NewNode("root", &mockWidget{})
	tree.SetRoot(root)
	
	// Build a balanced tree
	nodes := []*Node{root}
	for i := 1; i < nodeCount; i++ {
		node := NewNode("node-"+string(rune(i)), &mockWidget{})
		parent := nodes[i%len(nodes)]
		parent.AddChild(node)
		if i%3 == 0 {
			nodes = append(nodes, node)
		}
	}
	
	tree.SetRoot(root) // Rebuild index
	return tree
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