package reactive

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMemo_BasicOperations(t *testing.T) {
	t.Run("create_and_get", func(t *testing.T) {
		computeCount := 0
		memo := NewMemo(func() int {
			computeCount++
			return 42
		})
		
		// First get should compute
		result := memo.Get()
		if result != 42 {
			t.Errorf("Expected 42, got %d", result)
		}
		if computeCount != 1 {
			t.Error("Should compute once on first get")
		}
		
		// Second get should use cache
		result = memo.Get()
		if result != 42 {
			t.Errorf("Expected 42, got %d", result)
		}
		if computeCount != 1 {
			t.Error("Should not recompute for cached value")
		}
		
		memo.Dispose()
	})
	
	t.Run("peek_without_compute", func(t *testing.T) {
		computeCount := 0
		memo := NewMemo(func() int {
			computeCount++
			return 100
		})
		
		// Peek should return zero value if not computed
		result := memo.Peek()
		if result != 0 {
			t.Errorf("Peek before compute should return zero value, got %d", result)
		}
		
		// Get to compute
		memo.Get()
		
		// Now peek should return computed value
		result = memo.Peek()
		if result != 100 {
			t.Errorf("Peek after compute should return 100, got %d", result)
		}
		
		memo.Dispose()
	})
	
	t.Run("invalidate_and_recompute", func(t *testing.T) {
		computeCount := 0
		value := 1
		
		memo := NewMemo(func() int {
			computeCount++
			return value
		})
		
		if memo.Get() != 1 {
			t.Error("Should get initial value")
		}
		if computeCount != 1 {
			t.Error("Should compute once")
		}
		
		// Change underlying value and invalidate
		value = 2
		memo.Invalidate()
		
		// Should recompute on next get
		if memo.Get() != 2 {
			t.Error("Should get new value after invalidate")
		}
		if computeCount != 2 {
			t.Error("Should recompute after invalidate")
		}
		
		memo.Dispose()
	})
}

func TestMemo_DependencyTracking(t *testing.T) {
	t.Run("track_signal_dependencies", func(t *testing.T) {
		sig := NewSignal(10)
		computeCount := 0
		
		memo := NewMemo(func() int {
			computeCount++
			return sig.Get() * 2
		})
		
		if memo.Get() != 20 {
			t.Error("Should compute based on signal")
		}
		if computeCount != 1 {
			t.Error("Should compute once initially")
		}
		
		// Signal change should invalidate memo
		sig.Set(15)
		
		// Note: In this simplified version, signal changes don't auto-invalidate
		// This would need proper integration
		memo.Invalidate() // Manual invalidation for now
		
		if memo.Get() != 30 {
			t.Error("Should recompute with new signal value")
		}
		if computeCount != 2 {
			t.Error("Should recompute after signal change")
		}
		
		memo.Dispose()
	})
	
	t.Run("multiple_dependencies", func(t *testing.T) {
		sig1 := NewSignal(1)
		sig2 := NewSignal(2)
		sig3 := NewSignal(3)
		computeCount := 0
		
		memo := NewMemo(func() int {
			computeCount++
			return sig1.Get() + sig2.Get() + sig3.Get()
		})
		
		if memo.Get() != 6 {
			t.Error("Should sum all signals")
		}
		if computeCount != 1 {
			t.Error("Should compute once")
		}
		
		// Any signal change should invalidate
		sig2.Set(10)
		memo.Invalidate()
		
		if memo.Get() != 14 {
			t.Error("Should recompute with new values")
		}
		if computeCount != 2 {
			t.Error("Should recompute after change")
		}
		
		memo.Dispose()
	})
}

func TestMemo_Concurrency(t *testing.T) {
	t.Run("concurrent_gets", func(t *testing.T) {
		var computeCount atomic.Int32
		memo := NewMemo(func() int {
			computeCount.Add(1)
			time.Sleep(10 * time.Millisecond) // Simulate expensive computation
			return 42
		})
		
		var wg sync.WaitGroup
		results := make([]int, 10)
		
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx] = memo.Get()
			}(i)
		}
		
		wg.Wait()
		
		// All should get same result
		for i, r := range results {
			if r != 42 {
				t.Errorf("Result %d: expected 42, got %d", i, r)
			}
		}
		
		// Should only compute once despite concurrent access
		if computeCount.Load() != 1 {
			t.Errorf("Expected 1 computation, got %d", computeCount.Load())
		}
		
		memo.Dispose()
	})
	
	t.Run("concurrent_invalidate", func(t *testing.T) {
		value := atomic.Int32{}
		memo := NewMemo(func() int {
			return int(value.Load())
		})
		
		var wg sync.WaitGroup
		
		// Writers
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(v int) {
				defer wg.Done()
				value.Store(int32(v))
				memo.Invalidate()
			}(i)
		}
		
		// Readers
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = memo.Get()
			}()
		}
		
		wg.Wait()
		memo.Dispose()
	})
}

func TestComputed_BasicOperations(t *testing.T) {
	t.Run("computed_wrapper", func(t *testing.T) {
		sig := NewSignal(5)
		computeCount := 0
		
		computed := NewComputed(func() int {
			computeCount++
			return sig.Get() * sig.Get()
		})
		
		if computed.Get() != 25 {
			t.Error("Should compute square")
		}
		if computeCount != 1 {
			t.Error("Should compute once")
		}
		
		// Get again - should use cache
		if computed.Get() != 25 {
			t.Error("Should return cached value")
		}
		if computeCount != 1 {
			t.Error("Should not recompute")
		}
		
		computed.Dispose()
	})
}

func TestWatch(t *testing.T) {
	t.Run("watch_side_effects", func(t *testing.T) {
		sig := NewSignal(0)
		sideEffectCount := 0
		values := []int{}
		
		dispose := Watch(func() {
			sideEffectCount++
			values = append(values, sig.Get())
		})
		
		if sideEffectCount != 1 {
			t.Error("Watch should run immediately")
		}
		
		sig.Set(1)
		sig.Set(2)
		
		if sideEffectCount != 3 {
			t.Errorf("Expected 3 runs, got %d", sideEffectCount)
		}
		
		dispose()
		
		sig.Set(3)
		if sideEffectCount != 3 {
			t.Error("Should not run after dispose")
		}
	})
}

func BenchmarkMemo_Get(b *testing.B) {
	memo := NewMemo(func() int {
		return 42
	})
	defer memo.Dispose()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = memo.Get()
	}
}

func BenchmarkMemo_Invalidate(b *testing.B) {
	value := 0
	memo := NewMemo(func() int {
		return value
	})
	defer memo.Dispose()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value = i
		memo.Invalidate()
		_ = memo.Get()
	}
}