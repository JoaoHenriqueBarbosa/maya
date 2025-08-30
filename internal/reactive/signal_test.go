package reactive

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSignal_BasicOperations(t *testing.T) {
	t.Run("create_and_get", func(t *testing.T) {
		sig := NewSignal(42)
		
		if value := sig.Get(); value != 42 {
			t.Errorf("Expected 42, got %v", value)
		}
		
		if value := sig.Peek(); value != 42 {
			t.Errorf("Peek: Expected 42, got %v", value)
		}
	})
	
	t.Run("set_and_get", func(t *testing.T) {
		sig := NewSignal("hello")
		sig.Set("world")
		
		if value := sig.Get(); value != "world" {
			t.Errorf("Expected 'world', got %v", value)
		}
	})
	
	t.Run("update_function", func(t *testing.T) {
		sig := NewSignal(10)
		sig.Update(func(v int) int {
			return v * 2
		})
		
		if value := sig.Get(); value != 20 {
			t.Errorf("Expected 20, got %v", value)
		}
	})
	
	t.Run("version_tracking", func(t *testing.T) {
		sig := NewSignal(1)
		v1 := sig.Version()
		
		sig.Set(2)
		v2 := sig.Version()
		
		if v2 <= v1 {
			t.Error("Version should increment on set")
		}
		
		sig.Set(2) // Same value with no equality checker
		v3 := sig.Version()
		
		if v3 <= v2 {
			t.Error("Version should increment even for same value without equality")
		}
	})
}

func TestSignal_Equality(t *testing.T) {
	t.Run("with_equality_checker", func(t *testing.T) {
		updateCount := 0
		sig := NewSignalWithEquals(42, func(a, b int) bool {
			return a == b
		})
		
		sig.Subscribe(func(v int) {
			updateCount++
		})
		
		// Initial subscription call
		if updateCount != 1 {
			t.Error("Subscribe should call immediately")
		}
		
		// Set same value - should not trigger
		sig.Set(42)
		if updateCount != 1 {
			t.Error("Setting same value should not trigger update")
		}
		
		// Set different value - should trigger
		sig.Set(43)
		if updateCount != 2 {
			t.Error("Setting different value should trigger update")
		}
	})
	
	t.Run("without_equality_checker", func(t *testing.T) {
		updateCount := 0
		sig := NewSignal(42)
		
		sig.Subscribe(func(v int) {
			updateCount++
		})
		
		// Set same value - should trigger without equality
		sig.Set(42)
		if updateCount != 2 {
			t.Error("Without equality, same value should trigger update")
		}
	})
}

func TestSignal_Subscribe(t *testing.T) {
	t.Run("basic_subscription", func(t *testing.T) {
		sig := NewSignal(0)
		values := []int{}
		
		unsubscribe := sig.Subscribe(func(v int) {
			values = append(values, v)
		})
		
		sig.Set(1)
		sig.Set(2)
		sig.Set(3)
		
		if len(values) != 4 { // Initial + 3 updates
			t.Errorf("Expected 4 values, got %d", len(values))
		}
		
		unsubscribe()
		sig.Set(4)
		
		if len(values) != 4 {
			t.Error("Should not receive updates after unsubscribe")
		}
	})
	
	t.Run("multiple_subscribers", func(t *testing.T) {
		sig := NewSignal("a")
		count1, count2 := 0, 0
		
		unsub1 := sig.Subscribe(func(v string) {
			count1++
		})
		
		unsub2 := sig.Subscribe(func(v string) {
			count2++
		})
		
		sig.Set("b")
		
		if count1 != 2 || count2 != 2 {
			t.Error("Both subscribers should receive updates")
		}
		
		unsub1()
		sig.Set("c")
		
		if count1 != 2 {
			t.Error("Unsubscribed listener should not receive updates")
		}
		if count2 != 3 {
			t.Error("Active subscriber should still receive updates")
		}
		
		unsub2()
	})
}

func TestSignal_Concurrency(t *testing.T) {
	t.Run("concurrent_reads", func(t *testing.T) {
		sig := NewSignal(42)
		var wg sync.WaitGroup
		
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					_ = sig.Get()
					_ = sig.Peek()
				}
			}()
		}
		
		wg.Wait()
	})
	
	t.Run("concurrent_writes", func(t *testing.T) {
		sig := NewSignal(0)
		var wg sync.WaitGroup
		var finalValue atomic.Int32
		
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				sig.Set(val)
				finalValue.Store(int32(val))
			}(i)
		}
		
		wg.Wait()
		
		// Final value should be one of the written values
		final := sig.Get()
		if final < 0 || final >= 100 {
			t.Errorf("Unexpected final value: %d", final)
		}
	})
	
	t.Run("concurrent_subscribe_unsubscribe", func(t *testing.T) {
		sig := NewSignal(0)
		var wg sync.WaitGroup
		
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				unsub := sig.Subscribe(func(v int) {
					// Do nothing
				})
				time.Sleep(time.Microsecond)
				unsub()
			}()
		}
		
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				sig.Set(val)
			}(i)
		}
		
		wg.Wait()
	})
}

func TestSignal_MemoryLeaks(t *testing.T) {
	t.Run("dispose_cleans_observers", func(t *testing.T) {
		sig := NewSignal(0)
		
		// Add multiple subscribers
		for i := 0; i < 10; i++ {
			sig.Subscribe(func(v int) {})
		}
		
		if len(sig.observers) != 10 {
			t.Error("Should have 10 observers")
		}
		
		sig.Dispose()
		
		if sig.observers != nil {
			t.Error("Observers should be nil after dispose")
		}
	})
	
	t.Run("unsubscribe_removes_observer", func(t *testing.T) {
		sig := NewSignal(0)
		
		unsub1 := sig.Subscribe(func(v int) {})
		unsub2 := sig.Subscribe(func(v int) {})
		
		if len(sig.observers) != 2 {
			t.Error("Should have 2 observers")
		}
		
		unsub1()
		if len(sig.observers) != 1 {
			t.Error("Should have 1 observer after unsubscribe")
		}
		
		unsub2()
		if len(sig.observers) != 0 {
			t.Error("Should have 0 observers after all unsubscribed")
		}
	})
}

func BenchmarkSignal_Get(b *testing.B) {
	sig := NewSignal(42)
	
	b.ResetTimer()
	for b.Loop() {
		_ = sig.Get()
	}
}

func BenchmarkSignal_Set(b *testing.B) {
	sig := NewSignal(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sig.Set(i)
	}
}

func BenchmarkSignal_Subscribe(b *testing.B) {
	sig := NewSignal(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		unsub := sig.Subscribe(func(v int) {})
		unsub()
	}
}

func BenchmarkSignal_ConcurrentReadWrite(b *testing.B) {
	sig := NewSignal(0)
	
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				sig.Set(i)
			} else {
				_ = sig.Get()
			}
			i++
		}
	})
}