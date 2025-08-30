package core

import (
	"errors"
	"testing"
)

// TestEdgeCases_MissingCoverage tests edge cases to achieve 100% coverage
func TestEdgeCases_MissingCoverage(t *testing.T) {
	t.Run("iterators_with_nil_yields", func(t *testing.T) {
		tree := buildTestTree()
		
		// Test early termination in all iterators
		count := 0
		for range tree.PreOrderDFS() {
			count++
			if count >= 2 {
				break
			}
		}
		
		count = 0
		for range tree.PostOrderDFS() {
			count++
			if count >= 2 {
				break
			}
		}
		
		count = 0
		for range tree.BFS() {
			count++
			if count >= 2 {
				break
			}
		}
		
		count = 0
		for _, _ = range tree.LevelOrder() {
			count++
			if count >= 2 {
				break
			}
		}
		
		// Test DirtyNodes with early termination
		tree.GetRoot().MarkDirty(LayoutDirty)
		count = 0
		for range tree.DirtyNodes() {
			count++
			if count >= 1 {
				break
			}
		}
		
		// Test node iterators with early termination
		node := tree.GetRoot().Children[0]
		count = 0
		for range node.Ancestors() {
			count++
			if count >= 1 {
				break
			}
		}
		
		count = 0
		for range node.Descendants() {
			count++
			if count >= 1 {
				break
			}
		}
		
		count = 0
		for range node.Siblings() {
			count++
			if count >= 1 {
				break
			}
		}
	})
	
	t.Run("processSubtree_with_error_in_child", func(t *testing.T) {
		tree := NewTree()
		root := NewNode("root", &mockWidget{})
		child := NewNode("child", &mockWidget{})
		grandchild := NewNode("grandchild", &mockWidget{})
		
		root.AddChild(child)
		child.AddChild(grandchild)
		tree.SetRoot(root)
		
		// Test processSubtree error handling
		err := tree.processSubtree(child, func(n *Node) error {
			if n.ID == "grandchild" {
				return errors.New("test error")
			}
			return nil
		})
		
		if err == nil {
			t.Error("processSubtree should propagate error from child")
		}
	})
	
	t.Run("RemoveNode_with_multiple_children", func(t *testing.T) {
		tree := NewTree()
		root := NewNode("root", &mockWidget{})
		child1 := NewNode("child1", &mockWidget{})
		child2 := NewNode("child2", &mockWidget{})
		
		root.AddChild(child1)
		root.AddChild(child2)
		tree.SetRoot(root)
		
		// Remove child2 which is not the first child
		removed := tree.RemoveNode(child2)
		if !removed {
			t.Error("Should be able to remove second child")
		}
		
		if tree.FindNodeByID("child2") != nil {
			t.Error("Removed node should not be in index")
		}
	})
}

func buildTestTree() *Tree {
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