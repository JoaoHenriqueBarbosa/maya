package reactive

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestEffect_BasicOperations(t *testing.T) {
	t.Run("create_and_run", func(t *testing.T) {
		runCount := 0
		effect := CreateEffect(func() {
			runCount++
		})
		
		if runCount != 1 {
			t.Error("Effect should run immediately on creation")
		}
		
		if !effect.IsActive() {
			t.Error("Effect should be active")
		}
		
		effect.Dispose()
		if effect.IsActive() {
			t.Error("Effect should be inactive after dispose")
		}
	})
	
	t.Run("deferred_effect", func(t *testing.T) {
		runCount := 0
		effect := CreateEffectWithOptions(func() {
			runCount++
		}, EffectOptions{
			Immediate: false,
			Defer:     true,
		})
		
		if runCount != 0 {
			t.Error("Deferred effect should not run immediately")
		}
		
		// Trigger the effect
		effect.invalidate()
		
		// Wait for deferred execution
		time.Sleep(10 * time.Millisecond)
		
		if runCount != 1 {
			t.Error("Deferred effect should run after invalidation")
		}
		
		effect.Dispose()
	})
}

func TestEffect_DependencyTracking(t *testing.T) {
	t.Run("auto_track_signal", func(t *testing.T) {
		sig := NewSignal(0)
		runCount := 0
		values := []int{}
		
		effect := CreateEffect(func() {
			runCount++
			values = append(values, sig.Get())
		})
		
		if runCount != 1 {
			t.Error("Effect should run once on creation")
		}
		
		sig.Set(1)
		if runCount != 2 {
			t.Error("Effect should re-run when signal changes")
		}
		
		sig.Set(2)
		if runCount != 3 {
			t.Error("Effect should re-run on each signal change")
		}
		
		effect.Dispose()
		sig.Set(3)
		
		if runCount != 3 {
			t.Error("Disposed effect should not re-run")
		}
		
		expected := []int{0, 1, 2}
		if len(values) != len(expected) {
			t.Errorf("Expected %v, got %v", expected, values)
		}
		for i, v := range expected {
			if values[i] != v {
				t.Errorf("values[%d]: expected %d, got %d", i, v, values[i])
			}
		}
	})
	
	t.Run("multiple_signal_dependencies", func(t *testing.T) {
		sig1 := NewSignal(1)
		sig2 := NewSignal(2)
		runCount := 0
		sum := 0
		
		effect := CreateEffect(func() {
			runCount++
			sum = sig1.Get() + sig2.Get()
		})
		
		if sum != 3 {
			t.Errorf("Initial sum should be 3, got %d", sum)
		}
		
		sig1.Set(10)
		if sum != 12 {
			t.Errorf("Sum should be 12 after sig1 update, got %d", sum)
		}
		
		sig2.Set(20)
		if sum != 30 {
			t.Errorf("Sum should be 30 after sig2 update, got %d", sum)
		}
		
		if runCount != 3 {
			t.Errorf("Effect should have run 3 times, got %d", runCount)
		}
		
		effect.Dispose()
	})
	
	t.Run("dynamic_dependencies", func(t *testing.T) {
		condition := NewSignal(true)
		sigA := NewSignal("A")
		sigB := NewSignal("B")
		runCount := 0
		result := ""
		
		effect := CreateEffect(func() {
			runCount++
			if condition.Get() {
				result = sigA.Get()
			} else {
				result = sigB.Get()
			}
		})
		
		if result != "A" {
			t.Error("Should initially read from sigA")
		}
		
		// Change sigB - should not trigger since we're not reading it
		sigB.Set("B2")
		if runCount != 1 {
			t.Error("Should not re-run when unused signal changes")
		}
		
		// Change condition to false
		condition.Set(false)
		if result != "B2" {
			t.Error("Should now read from sigB")
		}
		
		// Now sigA changes should not trigger
		sigA.Set("A2")
		if runCount != 2 {
			t.Error("Should not re-run when now-unused signal changes")
		}
		
		// But sigB changes should trigger
		sigB.Set("B3")
		if result != "B3" {
			t.Error("Should update when active dependency changes")
		}
		
		effect.Dispose()
	})
}

func TestEffect_Cleanup(t *testing.T) {
	t.Run("cleanup_on_rerun", func(t *testing.T) {
		sig := NewSignal(0)
		cleanupCount := 0
		runCount := 0
		
		effect := CreateEffect(func() {
			runCount++
			current := getCurrentEffect()
			if current != nil {
				current.OnCleanup(func() {
					cleanupCount++
				})
			}
			_ = sig.Get()
		})
		
		sig.Set(1)
		if cleanupCount != 1 {
			t.Error("Cleanup should run before re-execution")
		}
		
		sig.Set(2)
		if cleanupCount != 2 {
			t.Error("Cleanup should run on each re-execution")
		}
		
		effect.Dispose()
		if cleanupCount != 3 {
			t.Error("Cleanup should run on dispose")
		}
	})
	
	t.Run("multiple_cleanups", func(t *testing.T) {
		sig := NewSignal(0)
		cleanups := []int{}
		
		effect := CreateEffect(func() {
			value := sig.Get()
			current := getCurrentEffect()
			if current != nil {
				current.OnCleanup(func() {
					cleanups = append(cleanups, value)
				})
			}
		})
		
		sig.Set(1)
		sig.Set(2)
		effect.Dispose()
		
		expected := []int{0, 1, 2}
		if len(cleanups) != len(expected) {
			t.Errorf("Expected %v cleanups, got %v", expected, cleanups)
		}
	})
}

func TestEffect_NestedEffects(t *testing.T) {
	t.Run("nested_effect_tracking", func(t *testing.T) {
		outer := NewSignal(1)
		inner := NewSignal(2)
		outerRuns := 0
		innerRuns := 0
		
		outerEffect := CreateEffect(func() {
			outerRuns++
			_ = outer.Get()
			
			CreateEffect(func() {
				innerRuns++
				_ = inner.Get()
			})
		})
		
		if outerRuns != 1 || innerRuns != 1 {
			t.Error("Both effects should run initially")
		}
		
		inner.Set(3)
		if innerRuns != 2 {
			t.Error("Inner effect should re-run on inner signal change")
		}
		if outerRuns != 1 {
			t.Error("Outer effect should not re-run on inner signal change")
		}
		
		outer.Set(2)
		if outerRuns != 2 {
			t.Error("Outer effect should re-run on outer signal change")
		}
		// Note: This creates a new inner effect, so innerRuns increases
		if innerRuns != 3 {
			t.Error("New inner effect should be created and run")
		}
		
		outerEffect.Dispose()
	})
}

func TestEffect_PreventRecursion(t *testing.T) {
	t.Run("prevent_infinite_loop", func(t *testing.T) {
		sig := NewSignal(0)
		runCount := 0
		
		effect := CreateEffect(func() {
			runCount++
			value := sig.Get()
			if value < 10 {
				sig.Set(value + 1) // This could cause infinite loop
			}
		})
		
		// Should not cause stack overflow
		time.Sleep(10 * time.Millisecond)
		
		if runCount > 100 {
			t.Error("Runaway recursion detected")
		}
		
		effect.Dispose()
	})
}

func TestEffect_Concurrency(t *testing.T) {
	t.Run("concurrent_signal_updates", func(t *testing.T) {
		sig := NewSignal(0)
		var runCount atomic.Int32
		
		effect := CreateEffect(func() {
			runCount.Add(1)
			_ = sig.Get()
		})
		
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				sig.Set(val)
			}(i)
		}
		
		wg.Wait()
		effect.Dispose()
		
		// Should have run at least once (initial) and at most 101 times
		count := runCount.Load()
		if count < 1 || count > 101 {
			t.Errorf("Unexpected run count: %d", count)
		}
	})
}

func BenchmarkEffect_Creation(b *testing.B) {
	sig := NewSignal(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		effect := CreateEffect(func() {
			_ = sig.Get()
		})
		effect.Dispose()
	}
}

func BenchmarkEffect_Invalidation(b *testing.B) {
	sig := NewSignal(0)
	effect := CreateEffect(func() {
		_ = sig.Get()
	})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sig.Set(i)
	}
	
	b.StopTimer()
	effect.Dispose()
}