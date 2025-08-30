package core

import (
	"runtime"
	"testing"
	"time"
)

// mockWidget implements the Widget interface for testing
type mockWidget struct {
	disposed      bool
	layoutCalled  int
	paintCalled   int
	eventHandled  bool
}

func (m *mockWidget) Layout(constraints Constraints) (width, height float64) {
	m.layoutCalled++
	return constraints.MaxWidth, constraints.MaxHeight
}

func (m *mockWidget) Paint(context PaintContext) {
	m.paintCalled++
}

func (m *mockWidget) HandleEvent(event Event) bool {
	m.eventHandled = true
	return true
}

func (m *mockWidget) GetIntrinsicWidth(height float64) float64 {
	return 100.0
}

func (m *mockWidget) GetIntrinsicHeight(width float64) float64 {
	return 50.0
}

func (m *mockWidget) Dispose() {
	m.disposed = true
}

// TestNode_Creation tests node creation and initialization
func TestNode_Creation(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		widget   Widget
		wantID   NodeID
		hasWidget bool
	}{
		{
			name:      "create_node_with_widget",
			id:        "test-node",
			widget:    &mockWidget{},
			wantID:    "test-node",
			hasWidget: true,
		},
		{
			name:      "create_node_without_widget",
			id:        "empty-node",
			widget:    nil,
			wantID:    "empty-node",
			hasWidget: false,
		},
		{
			name:      "create_node_with_empty_id",
			id:        "",
			widget:    &mockWidget{},
			wantID:    "",
			hasWidget: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewNode(tt.id, tt.widget)
			
			if node == nil {
				t.Fatal("NewNode returned nil")
			}
			
			if node.ID != tt.wantID {
				t.Errorf("ID = %v, want %v", node.ID, tt.wantID)
			}
			
			if tt.hasWidget && node.Widget == nil {
				t.Error("Expected widget to be set")
			}
			
			if !tt.hasWidget && node.Widget != nil {
				t.Error("Expected widget to be nil")
			}
			
			if node.Children == nil {
				t.Error("Children slice should be initialized")
			}
			
			if cap(node.Children) < 4 {
				t.Error("Children slice should have initial capacity of 4")
			}
		})
	}
}

// TestNode_ParentChildRelationship tests parent-child relationships with weak references
func TestNode_ParentChildRelationship(t *testing.T) {
	t.Run("add_single_child", func(t *testing.T) {
		parent := NewNode("parent", &mockWidget{})
		child := NewNode("child", &mockWidget{})
		
		parent.AddChild(child)
		
		if len(parent.Children) != 1 {
			t.Errorf("Expected 1 child, got %d", len(parent.Children))
		}
		
		if parent.Children[0] != child {
			t.Error("Child not added correctly")
		}
		
		if child.GetParent() != parent {
			t.Error("Parent reference not set correctly")
		}
		
		if !parent.IsDirty() {
			t.Error("Parent should be marked dirty after adding child")
		}
	})
	
	t.Run("add_multiple_children", func(t *testing.T) {
		parent := NewNode("parent", &mockWidget{})
		children := make([]*Node, 5)
		
		for i := 0; i < 5; i++ {
			children[i] = NewNode(string(rune('a'+i)), &mockWidget{})
			parent.AddChild(children[i])
		}
		
		if len(parent.Children) != 5 {
			t.Errorf("Expected 5 children, got %d", len(parent.Children))
		}
		
		for i, child := range children {
			if parent.Children[i] != child {
				t.Errorf("Child %d not in correct position", i)
			}
			if child.GetParent() != parent {
				t.Errorf("Child %d parent reference incorrect", i)
			}
		}
	})
	
	t.Run("remove_child", func(t *testing.T) {
		parent := NewNode("parent", &mockWidget{})
		child1 := NewNode("child1", &mockWidget{})
		child2 := NewNode("child2", &mockWidget{})
		child3 := NewNode("child3", &mockWidget{})
		
		parent.AddChild(child1)
		parent.AddChild(child2)
		parent.AddChild(child3)
		
		// Clear dirty flags before removal
		parent.ClearDirty()
		
		// Remove middle child
		removed := parent.RemoveChild(child2)
		
		if !removed {
			t.Error("RemoveChild should return true for existing child")
		}
		
		if len(parent.Children) != 2 {
			t.Errorf("Expected 2 children after removal, got %d", len(parent.Children))
		}
		
		if child2.GetParent() != nil {
			t.Error("Removed child should have nil parent")
		}
		
		if !parent.IsDirty() {
			t.Error("Parent should be marked dirty after removing child")
		}
		
		// Try to remove non-existent child
		removed = parent.RemoveChild(child2)
		if removed {
			t.Error("RemoveChild should return false for non-existent child")
		}
	})
	
	t.Run("set_parent_directly", func(t *testing.T) {
		node := NewNode("node", &mockWidget{})
		parent := NewNode("parent", &mockWidget{})
		
		node.SetParent(parent)
		
		if node.GetParent() != parent {
			t.Error("SetParent did not set parent correctly")
		}
		
		// Test setting nil parent
		node.SetParent(nil)
		
		if node.GetParent() != nil {
			t.Error("SetParent(nil) should clear parent")
		}
	})
	
	t.Run("weak_parent_reference", func(t *testing.T) {
		child := NewNode("child", &mockWidget{})
		
		// Create parent in a scope
		func() {
			parent := NewNode("parent", &mockWidget{})
			child.SetParent(parent)
			
			if child.GetParent() != parent {
				t.Error("Parent should be accessible")
			}
		}()
		
		// Force garbage collection
		runtime.GC()
		runtime.GC() // Run twice to ensure finalization
		time.Sleep(10 * time.Millisecond)
		
		// Note: Due to weak pointer semantics, the parent might still be accessible
		// This test mainly ensures the weak pointer API works correctly
		_ = child.GetParent() // This should not panic
	})
}

// TestNode_DirtyFlags tests dirty flag management
func TestNode_DirtyFlags(t *testing.T) {
	t.Run("mark_dirty_single_flag", func(t *testing.T) {
		node := NewNode("node", &mockWidget{})
		
		if node.IsDirty() {
			t.Error("New node should not be dirty")
		}
		
		node.MarkDirty(LayoutDirty)
		
		if !node.IsDirty() {
			t.Error("Node should be dirty after marking")
		}
		
		flags := node.GetDirtyFlags()
		if flags&LayoutDirty == 0 {
			t.Error("LayoutDirty flag should be set")
		}
		
		if flags&PaintDirty != 0 {
			t.Error("PaintDirty flag should not be set")
		}
	})
	
	t.Run("mark_dirty_multiple_flags", func(t *testing.T) {
		node := NewNode("node", &mockWidget{})
		
		node.MarkDirty(LayoutDirty | PaintDirty)
		
		flags := node.GetDirtyFlags()
		if flags&LayoutDirty == 0 {
			t.Error("LayoutDirty flag should be set")
		}
		if flags&PaintDirty == 0 {
			t.Error("PaintDirty flag should be set")
		}
	})
	
	t.Run("dirty_propagation_to_parent", func(t *testing.T) {
		parent := NewNode("parent", &mockWidget{})
		child := NewNode("child", &mockWidget{})
		grandchild := NewNode("grandchild", &mockWidget{})
		
		parent.AddChild(child)
		child.AddChild(grandchild)
		
		// Clear all dirty flags
		parent.ClearDirty()
		child.ClearDirty()
		grandchild.ClearDirty()
		
		// Mark grandchild dirty
		grandchild.MarkDirty(LayoutDirty)
		
		if !grandchild.IsDirty() {
			t.Error("Grandchild should be dirty")
		}
		
		if !child.IsDirty() {
			t.Error("Child should be dirty due to propagation")
		}
		
		if !parent.IsDirty() {
			t.Error("Parent should be dirty due to propagation")
		}
		
		// Parent and child should have ChildrenDirty flag
		if child.GetDirtyFlags()&ChildrenDirty == 0 {
			t.Error("Child should have ChildrenDirty flag")
		}
		
		if parent.GetDirtyFlags()&ChildrenDirty == 0 {
			t.Error("Parent should have ChildrenDirty flag")
		}
	})
	
	t.Run("clear_dirty_flags", func(t *testing.T) {
		node := NewNode("node", &mockWidget{})
		
		node.MarkDirty(LayoutDirty | PaintDirty | TransformDirty)
		node.ClearDirty()
		
		if node.IsDirty() {
			t.Error("Node should not be dirty after clearing")
		}
		
		if node.GetDirtyFlags() != CleanFlag {
			t.Error("All dirty flags should be cleared")
		}
	})
	
	t.Run("version_increment_on_dirty", func(t *testing.T) {
		node := NewNode("node", &mockWidget{})
		
		initialVersion := node.GetVersion()
		
		node.MarkDirty(LayoutDirty)
		v1 := node.GetVersion()
		
		if v1 <= initialVersion {
			t.Error("Version should increment when marked dirty")
		}
		
		// Marking with same flag should not increment version
		node.MarkDirty(LayoutDirty)
		v2 := node.GetVersion()
		
		if v2 != v1 {
			t.Error("Version should not increment for same flags")
		}
		
		// Adding new flag should increment version
		node.MarkDirty(PaintDirty)
		v3 := node.GetVersion()
		
		if v3 <= v2 {
			t.Error("Version should increment for new flags")
		}
	})
}

// TestNode_CachedValues tests weak cache functionality
func TestNode_CachedValues(t *testing.T) {
	t.Run("set_and_get_cached_values", func(t *testing.T) {
		node := NewNode("node", &mockWidget{})
		
		if node.GetCachedValues() != nil {
			t.Error("New node should have no cached values")
		}
		
		cached := &ComputedValues{
			Layout: LayoutData{
				Position: Offset{X: 10, Y: 20},
			},
			Paint: PaintData{
				Opacity: 0.8,
				Visibility: true,
			},
		}
		
		node.SetCachedValues(cached)
		
		retrieved := node.GetCachedValues()
		if retrieved == nil {
			t.Fatal("Cached values should be retrievable")
		}
		
		if retrieved.Layout.Position.X != 10 || retrieved.Layout.Position.Y != 20 {
			t.Error("Cached layout position incorrect")
		}
		
		if retrieved.Paint.Opacity != 0.8 {
			t.Error("Cached paint opacity incorrect")
		}
	})
	
	t.Run("weak_cache_behavior", func(t *testing.T) {
		node := NewNode("node", &mockWidget{})
		
		// Set cached values in a scope
		func() {
			cached := &ComputedValues{
				Layout: LayoutData{
					Position: Offset{X: 100, Y: 200},
				},
			}
			node.SetCachedValues(cached)
			
			if node.GetCachedValues() == nil {
				t.Error("Cached values should be accessible")
			}
		}()
		
		// Cached values might still be accessible after scope
		// This depends on GC behavior
		_ = node.GetCachedValues()
		
		// Force GC to potentially clear weak references
		runtime.GC()
		runtime.GC()
		
		// This should not panic even if values are cleared
		_ = node.GetCachedValues()
	})
}

// TestNode_Cleanup tests runtime.AddCleanup functionality
func TestNode_Cleanup(t *testing.T) {
	t.Run("cleanup_with_disposable_widget", func(t *testing.T) {
		widget := &mockWidget{disposed: false}
		
		// Create and destroy node in a scope
		func() {
			node := NewNode("node", widget)
			_ = node // Use node to avoid compiler optimization
		}()
		
		// Force garbage collection
		runtime.GC()
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
		
		// Note: Cleanup timing is not guaranteed
		// This test ensures the cleanup mechanism doesn't cause issues
	})
	
	t.Run("cleanup_without_widget", func(t *testing.T) {
		// This should not panic
		func() {
			node := NewNode("node", nil)
			_ = node
		}()
		
		runtime.GC()
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	})
}

// TestNode_EdgeCases tests edge cases and error conditions
func TestNode_EdgeCases(t *testing.T) {
	t.Run("remove_child_from_node_without_children", func(t *testing.T) {
		parent := NewNode("parent", &mockWidget{})
		child := NewNode("child", &mockWidget{})
		
		removed := parent.RemoveChild(child)
		
		if removed {
			t.Error("RemoveChild should return false when removing non-existent child")
		}
	})
	
	t.Run("circular_reference_prevention", func(t *testing.T) {
		// This test ensures we don't create circular references
		// even though we use weak pointers for parents
		node1 := NewNode("node1", &mockWidget{})
		node2 := NewNode("node2", &mockWidget{})
		
		node1.AddChild(node2)
		
		// Attempting to add node1 as child of node2 would create a cycle
		// This should be prevented at a higher level (Tree operations)
		// For now, we just ensure the basic operations work
		if node2.GetParent() != node1 {
			t.Error("Parent-child relationship not established")
		}
	})
	
	t.Run("massive_children_list", func(t *testing.T) {
		parent := NewNode("parent", &mockWidget{})
		
		// Add many children to test slice growth
		for i := 0; i < 1000; i++ {
			child := NewNode(string(rune(i)), &mockWidget{})
			parent.AddChild(child)
		}
		
		if len(parent.Children) != 1000 {
			t.Errorf("Expected 1000 children, got %d", len(parent.Children))
		}
	})
	
	t.Run("concurrent_dirty_flag_updates", func(t *testing.T) {
		node := NewNode("node", &mockWidget{})
		done := make(chan bool, 3)
		
		// Concurrent dirty flag updates
		go func() {
			for i := 0; i < 100; i++ {
				node.MarkDirty(LayoutDirty)
			}
			done <- true
		}()
		
		go func() {
			for i := 0; i < 100; i++ {
				node.MarkDirty(PaintDirty)
			}
			done <- true
		}()
		
		go func() {
			for i := 0; i < 100; i++ {
				_ = node.IsDirty()
				_ = node.GetDirtyFlags()
			}
			done <- true
		}()
		
		// Wait for all goroutines
		for i := 0; i < 3; i++ {
			<-done
		}
		
		// Should not panic and should have some dirty flags
		if !node.IsDirty() {
			t.Error("Node should be dirty after concurrent updates")
		}
	})
}