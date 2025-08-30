package reactive

import (
	"sync"
	"sync/atomic"
	"weak"
	"github.com/maya-framework/maya/internal/logger"
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
	
	// Dependencies and observers
	effect       *Effect
	dependencies []SignalInterface
	depmu        sync.RWMutex
	observers    map[uint64]*Effect // Effects watching this memo
	obsmu        sync.RWMutex
}

// NewMemo creates a new memoized computation
func NewMemo[T any](compute func() T) *Memo[T] {
	logger.Trace("MEMO", "Creating new memo")
	m := &Memo[T]{
		compute:   compute,
		observers: make(map[uint64]*Effect),
	}
	
	// Mark as stale initially to force first computation
	m.stale.Store(true)
	
	// Do initial computation to capture dependencies
	// This will set up the effect that watches the signals
	logger.Trace("MEMO", "Doing initial computation")
	defer func() {
		if r := recover(); r != nil {
			logger.Error("MEMO", "Panic during initial computation: %v", r)
		}
	}()
	
	// Force initial computation to set up dependencies
	m.recompute()
	logger.Trace("MEMO", "Initial computation complete")
	
	return m
}

// Get returns the memoized value, recomputing if necessary
func (m *Memo[T]) Get() T {
	// Check if we're tracking dependencies
	if current := getCurrentEffect(); current != nil {
		// Register this effect as an observer of the memo
		m.obsmu.Lock()
		m.observers[current.id] = current
		m.obsmu.Unlock()
		
		// Also track in the effect's dependencies
		current.depmu.Lock()
		current.dependencies[m] = m.version.Load()
		current.depmu.Unlock()
		
		logger.Trace("MEMO", "Registered observer effect %d", current.id)
	}
	
	// Fast path: check if cached and not stale
	if !m.stale.Load() {
		m.mu.RLock()
		value := m.cached
		m.mu.RUnlock()
		logger.Trace("MEMO", "Returning cached value")
		return value
	}
	
	// Slow path: recompute
	logger.Trace("MEMO", "Recomputing stale value")
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
	var value T
	
	// Dispose old effect if exists
	if m.effect != nil {
		m.effect.Dispose()
	}
	
	// The effect just marks the memo as stale when dependencies change
	// We don't recompute here - that happens lazily in Get()
	memoEffect := CreateEffectWithOptions(func() {
		logger.Trace("MEMO", "Dependencies changed, marking memo as stale")
		m.stale.Store(true)
		m.notifyObservers()
	}, EffectOptions{
		Immediate: false, // Don't run immediately - we'll set up deps manually
		Defer: false,
	})
	
	// Now compute the value while the effect is active in the stack
	// This allows the compute function to register dependencies
	pushEffect(memoEffect)
	value = m.compute()
	popEffect()
	
	logger.Trace("MEMO", "Computed value, effect has %d dependencies", len(memoEffect.dependencies))
	
	// Store the effect for future cleanup
	m.effect = memoEffect
	
	// Mark the effect as active so it can respond to changes
	memoEffect.active.Store(true)
	
	// The effect was already registered with signals during compute()
	logger.Trace("MEMO", "Effect registered with %d dependencies", len(memoEffect.dependencies))
	
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
	
	// Notify observers if value changed
	// Note: This assumes T is comparable. For non-comparable types,
	// we'd need a different approach
	m.notifyObservers()
	
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

// Implement SignalInterface for Memo
func (m *Memo[T]) AddObserver(effect *Effect) {
	m.obsmu.Lock()
	defer m.obsmu.Unlock()
	m.observers[effect.id] = effect
}

func (m *Memo[T]) RemoveObserver(effectID uint64) {
	m.obsmu.Lock()
	defer m.obsmu.Unlock()
	delete(m.observers, effectID)
}

func (m *Memo[T]) GetVersion() uint64 {
	return m.version.Load()
}

func (m *Memo[T]) Version() uint64 {
	return m.version.Load()
}

// notifyObservers notifies all observers that the memo has changed
func (m *Memo[T]) notifyObservers() {
	m.obsmu.RLock()
	observers := make([]*Effect, 0, len(m.observers))
	for _, obs := range m.observers {
		observers = append(observers, obs)
	}
	m.obsmu.RUnlock()
	
	for _, obs := range observers {
		obs.Invalidate()
	}
}

// Implement remaining SignalInterface methods
func (m *Memo[T]) notify() {
	m.notifyObservers()
}

func (m *Memo[T]) removeObserver(effect *Effect) {
	m.RemoveObserver(effect.id)
}

func (m *Memo[T]) getObservers() []*Effect {
	m.obsmu.RLock()
	defer m.obsmu.RUnlock()
	
	observers := make([]*Effect, 0, len(m.observers))
	for _, obs := range m.observers {
		observers = append(observers, obs)
	}
	return observers
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