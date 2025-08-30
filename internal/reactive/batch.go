package reactive

import (
	"sync"
	"sync/atomic"
)

var (
	// Batch tracking
	batchDepth    atomic.Int32
	batchMutex    sync.Mutex
	pendingSignals = make(map[SignalInterface]struct{})
	pendingEffects = make(map[*Effect]struct{})
	
	// Effect scheduler
	scheduledEffects   = make(map[*Effect]struct{})
	scheduleMutex      sync.Mutex
	scheduleRunning    atomic.Bool
)

// Batch defers all signal updates until the batch completes
func Batch(fn func()) {
	startBatch()
	defer endBatch()
	fn()
}

// BatchValue runs a function in a batch and returns its result
func BatchValue[T any](fn func() T) T {
	startBatch()
	defer endBatch()
	return fn()
}

// startBatch begins a new batch operation
func startBatch() {
	batchDepth.Add(1)
}

// endBatch completes a batch and flushes pending updates
func endBatch() {
	if batchDepth.Add(-1) == 0 {
		flushBatch()
	}
}

// isInBatch returns true if currently in a batch
func isInBatch() bool {
	return batchDepth.Load() > 0
}

// addPendingSignal queues a signal for batch notification
func addPendingSignal(s SignalInterface) {
	batchMutex.Lock()
	defer batchMutex.Unlock()
	pendingSignals[s] = struct{}{}
}

// addPendingEffect queues an effect for batch execution
func addPendingEffect(e *Effect) {
	batchMutex.Lock()
	defer batchMutex.Unlock()
	pendingEffects[e] = struct{}{}
}

// flushBatch processes all pending updates
func flushBatch() {
	batchMutex.Lock()
	
	// Copy and clear pending signals
	signals := make([]SignalInterface, 0, len(pendingSignals))
	for s := range pendingSignals {
		signals = append(signals, s)
	}
	pendingSignals = make(map[SignalInterface]struct{})
	
	// Copy and clear pending effects  
	effects := make([]*Effect, 0, len(pendingEffects))
	for e := range pendingEffects {
		effects = append(effects, e)
	}
	pendingEffects = make(map[*Effect]struct{})
	
	batchMutex.Unlock()
	
	// First collect all observers from signals
	for _, s := range signals {
		observers := s.getObservers()
		for _, obs := range observers {
			addPendingEffect(obs)
		}
	}
	
	// Collect all pending effects again after signal notifications
	batchMutex.Lock()
	for e := range pendingEffects {
		effects = append(effects, e)
	}
	pendingEffects = make(map[*Effect]struct{})
	batchMutex.Unlock()
	
	// Run all unique effects once
	seen := make(map[uint64]bool)
	for _, e := range effects {
		if e.IsActive() && !seen[e.id] {
			seen[e.id] = true
			e.run()
		}
	}
}

// scheduleEffect queues an effect for deferred execution
func scheduleEffect(e *Effect) {
	scheduleMutex.Lock()
	scheduledEffects[e] = struct{}{}
	needsRun := len(scheduledEffects) == 1
	scheduleMutex.Unlock()
	
	if needsRun && scheduleRunning.CompareAndSwap(false, true) {
		go runScheduledEffects()
	}
}

// runScheduledEffects processes all scheduled effects
func runScheduledEffects() {
	defer scheduleRunning.Store(false)
	
	for {
		scheduleMutex.Lock()
		if len(scheduledEffects) == 0 {
			scheduleMutex.Unlock()
			break
		}
		
		// Take all scheduled effects
		effects := make([]*Effect, 0, len(scheduledEffects))
		for e := range scheduledEffects {
			effects = append(effects, e)
		}
		scheduledEffects = make(map[*Effect]struct{})
		scheduleMutex.Unlock()
		
		// Run effects
		for _, e := range effects {
			if e.IsActive() && e.dirty.Load() {
				e.run()
			}
		}
	}
}

// Transaction runs multiple operations as a single batch
type Transaction struct {
	completed atomic.Bool
}

// NewTransaction creates a new transaction
func NewTransaction() *Transaction {
	t := &Transaction{}
	startBatch()
	return t
}

// Commit completes the transaction
func (t *Transaction) Commit() {
	if t.completed.CompareAndSwap(false, true) {
		endBatch()
	}
}

// Rollback cancels the transaction (in this simple version, just ends the batch)
func (t *Transaction) Rollback() {
	if t.completed.CompareAndSwap(false, true) {
		// In a real implementation, this would restore previous values
		endBatch()
	}
}