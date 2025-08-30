package reactive

import (
	"testing"
)

func TestTracking_Untrack(t *testing.T) {
	t.Run("untrack_prevents_dependency", func(t *testing.T) {
		sig := NewSignal(0)
		tracked := false
		notTracked := false
		
		effect := CreateEffect(func() {
			// This should be tracked
			_ = sig.Get()
			tracked = true
			
			// This should NOT be tracked
			Untrack(func() int {
				sig.Get()
				return 0
			})
			notTracked = true
		})
		
		if !tracked || !notTracked {
			t.Error("Effect should have run")
		}
		
		// Reset flags
		tracked = false
		notTracked = false
		
		// Change signal - effect should re-run
		sig.Set(1)
		
		if !tracked {
			t.Error("Effect should re-run on signal change")
		}
		
		effect.Dispose()
	})
	
	t.Run("untrack_void", func(t *testing.T) {
		sig := NewSignal(0)
		runCount := 0
		
		effect := CreateEffect(func() {
			runCount++
			_ = sig.Get()
			
			// Use UntrackVoid for functions without return
			UntrackVoid(func() {
				// This read should not create dependency
				_ = sig.Get()
			})
		})
		
		if runCount != 1 {
			t.Error("Effect should run once initially")
		}
		
		sig.Set(1)
		if runCount != 2 {
			t.Error("Effect should re-run when tracked signal changes")
		}
		
		effect.Dispose()
	})
	
	t.Run("nested_untrack", func(t *testing.T) {
		sig1 := NewSignal(1)
		sig2 := NewSignal(2)
		sig3 := NewSignal(3)
		
		deps := []int{}
		
		effect := CreateEffect(func() {
			deps = []int{}
			
			// Tracked
			deps = append(deps, sig1.Get())
			
			// Not tracked
			Untrack(func() int {
				v2 := sig2.Get()
				
				// Nested untrack
				v3 := Untrack(func() int {
					return sig3.Get()
				})
				
				return v2 + v3
			})
			
			// Tracked again
			deps = append(deps, sig1.Get())
		})
		
		initialDeps := len(deps)
		
		// sig2 and sig3 changes should not trigger
		sig2.Set(20)
		sig3.Set(30)
		
		if len(deps) != initialDeps {
			t.Error("Untracked signals should not trigger effect")
		}
		
		// sig1 change should trigger
		sig1.Set(10)
		
		if deps[0] != 10 || deps[1] != 10 {
			t.Error("Tracked signal should trigger effect")
		}
		
		effect.Dispose()
	})
}

func TestTracking_GetGoroutineID(t *testing.T) {
	t.Run("unique_ids", func(t *testing.T) {
		ids := make(chan uint64, 10)
		
		for i := 0; i < 10; i++ {
			go func() {
				ids <- getGoroutineID()
			}()
		}
		
		seen := make(map[uint64]bool)
		for i := 0; i < 10; i++ {
			id := <-ids
			if id == 0 {
				t.Error("Goroutine ID should not be 0")
			}
			if seen[id] {
				t.Error("Goroutine IDs should be unique")
			}
			seen[id] = true
		}
	})
	
	t.Run("consistent_id", func(t *testing.T) {
		id1 := getGoroutineID()
		id2 := getGoroutineID()
		
		if id1 != id2 {
			t.Error("Same goroutine should have consistent ID")
		}
		
		if id1 == 0 {
			t.Error("ID should not be 0")
		}
	})
}

func TestTracking_EffectStack(t *testing.T) {
	t.Run("effect_stack_management", func(t *testing.T) {
		// Initially no effect
		if getCurrentEffect() != nil {
			t.Error("Should have no current effect initially")
		}
		
		effect1 := &Effect{id: 1}
		pushEffect(effect1)
		
		if getCurrentEffect() != effect1 {
			t.Error("Should get pushed effect")
		}
		
		effect2 := &Effect{id: 2}
		pushEffect(effect2)
		
		if getCurrentEffect() != effect2 {
			t.Error("Should get most recent effect")
		}
		
		popEffect()
		
		if getCurrentEffect() != effect1 {
			t.Error("Should get previous effect after pop")
		}
		
		popEffect()
		
		if getCurrentEffect() != nil {
			t.Error("Should have no effect after popping all")
		}
	})
	
	t.Run("stack_cleanup", func(t *testing.T) {
		// Push and pop multiple times
		for i := 0; i < 5; i++ {
			effect := &Effect{id: uint64(i)}
			pushEffect(effect)
		}
		
		// Pop all
		for i := 0; i < 5; i++ {
			popEffect()
		}
		
		// Extra pops should not panic
		popEffect()
		popEffect()
		
		if getCurrentEffect() != nil {
			t.Error("Stack should be empty")
		}
	})
}