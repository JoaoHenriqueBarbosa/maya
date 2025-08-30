package reactive

import (
	"sync"
	"testing"
	"time"
)

func TestBatch_BasicOperations(t *testing.T) {
	t.Run("batch_updates", func(t *testing.T) {
		sig1 := NewSignal(0)
		sig2 := NewSignal(0)
		updateCount := 0
		
		CreateEffect(func() {
			updateCount++
			_ = sig1.Get() + sig2.Get()
		})
		
		// Reset counter after initial run
		updateCount = 0
		
		// Without batch - should trigger 2 updates
		sig1.Set(1)
		sig2.Set(2)
		
		if updateCount != 2 {
			t.Errorf("Without batch: expected 2 updates, got %d", updateCount)
		}
		
		// With batch - should trigger 1 update
		updateCount = 0
		Batch(func() {
			sig1.Set(10)
			sig2.Set(20)
		})
		
		// Allow time for batch to complete
		time.Sleep(10 * time.Millisecond)
		
		if updateCount != 1 {
			t.Errorf("With batch: expected 1 update, got %d", updateCount)
		}
	})
	
	t.Run("batch_value_return", func(t *testing.T) {
		sig := NewSignal(5)
		
		result := BatchValue(func() int {
			sig.Set(10)
			return sig.Get() * 2
		})
		
		if result != 20 {
			t.Errorf("Expected 20, got %d", result)
		}
	})
	
	t.Run("nested_batches", func(t *testing.T) {
		sig := NewSignal(0)
		updateCount := 0
		
		CreateEffect(func() {
			updateCount++
			_ = sig.Get()
		})
		
		updateCount = 0
		
		Batch(func() {
			sig.Set(1)
			
			Batch(func() {
				sig.Set(2)
				sig.Set(3)
			})
			
			sig.Set(4)
		})
		
		// Allow batch to complete
		time.Sleep(10 * time.Millisecond)
		
		// Should only update once after all batches complete
		if updateCount != 1 {
			t.Errorf("Nested batch: expected 1 update, got %d", updateCount)
		}
		
		if sig.Get() != 4 {
			t.Errorf("Final value should be 4, got %d", sig.Get())
		}
	})
}

func TestBatch_MultipleSignals(t *testing.T) {
	t.Run("batch_multiple_signals", func(t *testing.T) {
		signals := make([]*Signal[int], 10)
		for i := range signals {
			signals[i] = NewSignal(0)
		}
		
		updateCount := 0
		CreateEffect(func() {
			updateCount++
			sum := 0
			for _, sig := range signals {
				sum += sig.Get()
			}
		})
		
		updateCount = 0
		
		Batch(func() {
			for i, sig := range signals {
				sig.Set(i + 1)
			}
		})
		
		// Allow batch to complete
		time.Sleep(10 * time.Millisecond)
		
		// Should only trigger one update for all signals
		if updateCount != 1 {
			t.Errorf("Expected 1 update for batch, got %d", updateCount)
		}
	})
}

func TestBatch_Concurrency(t *testing.T) {
	t.Run("concurrent_batches", func(t *testing.T) {
		sig := NewSignal(0)
		updateCount := 0
		var mu sync.Mutex
		
		CreateEffect(func() {
			mu.Lock()
			updateCount++
			mu.Unlock()
			_ = sig.Get()
		})
		
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				Batch(func() {
					for j := 0; j < 10; j++ {
						sig.Set(val*10 + j)
					}
				})
			}(i)
		}
		
		wg.Wait()
		time.Sleep(50 * time.Millisecond)
		
		mu.Lock()
		count := updateCount
		mu.Unlock()
		
		// Should have significantly fewer updates than 100
		if count >= 100 {
			t.Errorf("Batching ineffective: got %d updates", count)
		}
	})
}

func TestTransaction(t *testing.T) {
	t.Run("transaction_commit", func(t *testing.T) {
		sig1 := NewSignal(0)
		sig2 := NewSignal(0)
		updateCount := 0
		
		CreateEffect(func() {
			updateCount++
			_ = sig1.Get() + sig2.Get()
		})
		
		updateCount = 0
		
		tx := NewTransaction()
		sig1.Set(5)
		sig2.Set(10)
		tx.Commit()
		
		time.Sleep(10 * time.Millisecond)
		
		if updateCount != 1 {
			t.Errorf("Transaction should batch updates, got %d updates", updateCount)
		}
	})
	
	t.Run("transaction_rollback", func(t *testing.T) {
		sig := NewSignal(0)
		updateCount := 0
		
		CreateEffect(func() {
			updateCount++
			_ = sig.Get()
		})
		
		updateCount = 0
		
		tx := NewTransaction()
		sig.Set(5)
		tx.Rollback() // In this simple implementation, just ends batch
		
		time.Sleep(10 * time.Millisecond)
		
		// Rollback in this implementation doesn't restore values
		// It just completes the batch
		if sig.Get() != 5 {
			t.Error("Value should still be updated")
		}
	})
	
	t.Run("transaction_double_commit", func(t *testing.T) {
		tx := NewTransaction()
		tx.Commit()
		tx.Commit() // Should not panic
		
		// Also test rollback after commit
		tx2 := NewTransaction()
		tx2.Commit()
		tx2.Rollback() // Should not panic
	})
}

func TestScheduledEffects(t *testing.T) {
	t.Run("deferred_effect_scheduling", func(t *testing.T) {
		sig := NewSignal(0)
		runTimes := []time.Time{}
		var mu sync.Mutex
		
		CreateEffectWithOptions(func() {
			mu.Lock()
			runTimes = append(runTimes, time.Now())
			mu.Unlock()
			_ = sig.Get()
		}, EffectOptions{
			Immediate: true,
			Defer:     true,
		})
		
		// Trigger multiple updates rapidly
		for i := 1; i <= 5; i++ {
			sig.Set(i)
			time.Sleep(time.Millisecond)
		}
		
		// Wait for deferred execution
		time.Sleep(50 * time.Millisecond)
		
		mu.Lock()
		defer mu.Unlock()
		
		// Should have run once initially, then batched the updates
		if len(runTimes) > 3 {
			t.Errorf("Deferred effects should batch, got %d runs", len(runTimes))
		}
	})
}

func BenchmarkBatch_SingleSignal(b *testing.B) {
	sig := NewSignal(0)
	CreateEffect(func() {
		_ = sig.Get()
	})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Batch(func() {
			for j := 0; j < 10; j++ {
				sig.Set(j)
			}
		})
	}
}

func BenchmarkBatch_MultipleSignals(b *testing.B) {
	signals := make([]*Signal[int], 10)
	for i := range signals {
		signals[i] = NewSignal(0)
	}
	
	CreateEffect(func() {
		sum := 0
		for _, sig := range signals {
			sum += sig.Get()
		}
	})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Batch(func() {
			for j, sig := range signals {
				sig.Set(i + j)
			}
		})
	}
}