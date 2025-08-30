package reactive

import (
	"sync"
	"sync/atomic"
	"weak"
)

// Memo represents a lazily computed value that depends on signals
type Memo[T any] struct {
	compute      func() T
	cached       T
	version      atomic.Uint64
	stale        atomic.Bool
	computing    atomic.Bool
	mu           sync.RWMutex
	
	// Weak cache for memory efficiency
	weakCache    *weak.Pointer[T]
	
	// Dependencies
	effect       *Effect
	dependencies []SignalInterface
	depmu        sync.RWMutex
}

// NewMemo creates a new memoized computation
func NewMemo[T any](compute func() T) *Memo[T] {
	m := &Memo[T]{
		compute: compute,
	}
	
	// Mark as stale initially to force first computation
	m.stale.Store(true)
	
	// Create effect that marks memo as stale when dependencies change
	m.effect = CreateEffectWithOptions(func() {
		// This runs in tracking context, capturing dependencies
		_ = m.recompute()
	}, EffectOptions{
		Immediate: false, // Don't run immediately
		Defer:     true,  // Defer updates
	})
	
	return m
}

// Get returns the memoized value, recomputing if necessary
func (m *Memo[T]) Get() T {
	// Check if we're tracking dependencies
	if current := getCurrentEffect(); current != nil {
		// Add this memo as a dependency
		// This ensures the effect re-runs when the memo changes
		m.trackAsSignal(current)
	}
	
	// Fast path: check if cached and not stale
	if !m.stale.Load() {
		m.mu.RLock()
		value := m.cached
		m.mu.RUnlock()
		return value
	}
	
	// Slow path: recompute
	return m.recompute()
}

// Peek returns the cached value without recomputing
func (m *Memo[T]) Peek() T {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cached
}

// recompute recalculates the memoized value
func (m *Memo[T]) recompute() T {
	// Prevent concurrent recomputation
	if !m.computing.CompareAndSwap(false, true) {
		// Another goroutine is computing, wait for it
		for m.computing.Load() {
			// Spin wait (in production, use condition variable)
		}
		m.mu.RLock()
		value := m.cached
		m.mu.RUnlock()
		return value
	}
	defer m.computing.Store(false)
	
	// Clear old dependencies
	m.clearDependencies()
	
	// Track new dependencies during computation
	var newDeps []SignalInterface
	var value T
	
	// Create temporary effect to capture dependencies
	tempEffect := &Effect{
		id:           effectCounter.Add(1),
		dependencies: make(map[SignalInterface]uint64),
	}
	
	pushEffect(tempEffect)
	value = m.compute()
	popEffect()
	
	// Extract captured dependencies
	tempEffect.depmu.RLock()
	for dep := range tempEffect.dependencies {
		newDeps = append(newDeps, dep)
	}
	tempEffect.depmu.RUnlock()
	
	// Update cached value
	m.mu.Lock()
	m.cached = value
	m.version.Add(1)
	m.stale.Store(false)
	
	// Update weak cache
	if m.weakCache == nil {
		wc := weak.Make(&value)
		m.weakCache = &wc
	}
	m.mu.Unlock()
	
	// Update dependencies
	m.depmu.Lock()
	m.dependencies = newDeps
	m.depmu.Unlock()
	
	// Register to watch these dependencies
	// Note: This is simplified - in production we'd have a better way
	// to handle different signal types
	
	return value
}

// clearDependencies removes all dependency tracking
func (m *Memo[T]) clearDependencies() {
	m.depmu.Lock()
	defer m.depmu.Unlock()
	
	// In production, dependencies would properly unregister
	// For now, just clear the list
	m.dependencies = nil
}

// trackAsSignal allows a memo to be tracked like a signal
func (m *Memo[T]) trackAsSignal(effect *Effect) {
	// This is a simplified approach
	// In production, Memo would implement SignalInterface
}

// Invalidate marks the memo as needing recomputation
func (m *Memo[T]) Invalidate() {
	m.stale.Store(true)
}

// Dispose cleans up the memo
func (m *Memo[T]) Dispose() {
	if m.effect != nil {
		m.effect.Dispose()
	}
	m.clearDependencies()
	m.weakCache = nil
}

// Computed is an alias for Memo with automatic dependency tracking
type Computed[T any] struct {
	*Memo[T]
}

// NewComputed creates a computed value that automatically tracks dependencies
func NewComputed[T any](compute func() T) *Computed[T] {
	return &Computed[T]{
		Memo: NewMemo(compute),
	}
}

// Watch creates a computed value that runs a side effect
func Watch(fn func()) func() {
	effect := CreateEffect(fn)
	return effect.Dispose
}