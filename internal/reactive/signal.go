package reactive

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Signal represents a reactive value that can be observed for changes
type Signal[T any] struct {
	value    T
	version  atomic.Uint64
	mu       sync.RWMutex
	
	// Tracking
	observers map[uint64]*Effect
	obsmu     sync.RWMutex
	
	// Equality checker for optimization
	equals    func(a, b T) bool
}

// NewSignal creates a new signal with an initial value
func NewSignal[T any](initial T) *Signal[T] {
	s := &Signal[T]{
		value:     initial,
		observers: make(map[uint64]*Effect),
		equals:    nil, // Default to always update
	}
	
	// Add default equality check for comparable types
	switch any(initial).(type) {
	case bool:
		s.equals = func(a, b T) bool {
			aBool, aOk := any(a).(bool)
			bBool, bOk := any(b).(bool)
			return aOk && bOk && aBool == bBool
		}
	case int:
		s.equals = func(a, b T) bool {
			aInt, aOk := any(a).(int)
			bInt, bOk := any(b).(int)
			return aOk && bOk && aInt == bInt
		}
	case string:
		s.equals = func(a, b T) bool {
			aStr, aOk := any(a).(string)
			bStr, bOk := any(b).(string)
			return aOk && bOk && aStr == bStr
		}
	}
	
	return s
}

// NewSignalWithEquals creates a signal with custom equality checking
func NewSignalWithEquals[T any](initial T, equals func(a, b T) bool) *Signal[T] {
	s := NewSignal(initial)
	s.equals = equals
	return s
}

// Get returns the current value and tracks dependencies
func (s *Signal[T]) Get() T {
	// Track this signal as a dependency of the current effect
	if current := getCurrentEffect(); current != nil {
		println("[SIGNAL] Tracking read by effect ID:", current.id)
		s.addObserver(current)
		current.addDependency(s)
		println("[SIGNAL] Added effect", current.id, "as observer, now have", len(s.observers), "observers")
	} else {
		println("[SIGNAL] No current effect to track")
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// Peek returns the value without tracking dependencies
func (s *Signal[T]) Peek() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// Set updates the signal value and notifies observers
func (s *Signal[T]) Set(value T) {
	s.mu.Lock()
	
	// Check if value actually changed
	if s.equals != nil && s.equals(s.value, value) {
		s.mu.Unlock()
		return
	}
	
	oldValue := s.value
	s.value = value
	s.version.Add(1)
	s.mu.Unlock()
	
	println("[SIGNAL] Value changed from", fmt.Sprint(oldValue), "to", fmt.Sprint(value))
	
	// Check if we're in a batch
	if isInBatch() {
		println("[SIGNAL] In batch, queuing notification")
		addPendingSignal(s)
		return
	}
	
	println("[SIGNAL] Notifying", len(s.getObservers()), "observers")
	s.notify()
}

// Update modifies the value using a function
func (s *Signal[T]) Update(fn func(T) T) {
	s.mu.Lock()
	newValue := fn(s.value)
	s.mu.Unlock()
	s.Set(newValue)
}

// Subscribe adds a callback that runs when the signal changes
func (s *Signal[T]) Subscribe(callback func(T)) func() {
	effect := CreateEffect(func() {
		callback(s.Get())
	})
	
	// Return unsubscribe function
	return func() {
		effect.Dispose()
	}
}

// Version returns the current version number
func (s *Signal[T]) Version() uint64 {
	return s.version.Load()
}

// addObserver registers an effect as an observer
func (s *Signal[T]) addObserver(effect *Effect) {
	s.obsmu.Lock()
	defer s.obsmu.Unlock()
	s.observers[effect.id] = effect
}

// removeObserver unregisters an effect
func (s *Signal[T]) removeObserver(effect *Effect) {
	s.obsmu.Lock()
	defer s.obsmu.Unlock()
	delete(s.observers, effect.id)
}

// notify triggers all observer effects
func (s *Signal[T]) notify() {
	// If we're in a batch, don't notify yet
	if isInBatch() {
		addPendingSignal(s)
		return
	}
	
	s.obsmu.RLock()
	observers := make([]*Effect, 0, len(s.observers))
	for _, obs := range s.observers {
		observers = append(observers, obs)
	}
	s.obsmu.RUnlock()
	
	// Notify outside of lock to prevent deadlocks
	for _, obs := range observers {
		obs.invalidate()
	}
}

// getObservers returns a copy of the observers list
func (s *Signal[T]) getObservers() []*Effect {
	s.obsmu.RLock()
	defer s.obsmu.RUnlock()
	
	observers := make([]*Effect, 0, len(s.observers))
	for _, obs := range s.observers {
		observers = append(observers, obs)
	}
	return observers
}

// Dispose cleans up the signal
func (s *Signal[T]) Dispose() {
	s.obsmu.Lock()
	defer s.obsmu.Unlock()
	
	// Clear all observers
	for _, obs := range s.observers {
		obs.removeDependency(s)
	}
	s.observers = nil
}

// SignalInterface allows type-erased signal operations
type SignalInterface interface {
	notify()
	removeObserver(*Effect)
	getObservers() []*Effect
	Version() uint64
}

// Ensure Signal implements SignalInterface
var _ SignalInterface = (*Signal[int])(nil)