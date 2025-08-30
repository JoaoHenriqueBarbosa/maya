package reactive

import (
	"sync"
	"sync/atomic"
)

var (
	// Global effect tracking
	effectStack   = &sync.Map{} // goroutine ID -> *Effect stack
	effectCounter atomic.Uint64
)

// Effect represents a reactive computation that runs when dependencies change
type Effect struct {
	id           uint64
	fn           func()
	dependencies map[SignalInterface]uint64 // signal -> version at last run
	depmu        sync.RWMutex
	
	// State
	active       atomic.Bool
	dirty        atomic.Bool
	running      atomic.Bool
	
	// Cleanup
	cleanups     []func()
	cleanupmu    sync.Mutex
	
	// Options
	immediate    bool
	defer_       bool
}

// CreateEffect creates and runs a new effect
func CreateEffect(fn func()) *Effect {
	return CreateEffectWithOptions(fn, EffectOptions{
		Immediate: true,
	})
}

// EffectOptions configures effect behavior
type EffectOptions struct {
	Immediate bool // Run immediately on creation
	Defer     bool // Defer execution to next tick
}

// CreateEffectWithOptions creates an effect with custom options
func CreateEffectWithOptions(fn func(), opts EffectOptions) *Effect {
	e := &Effect{
		id:           effectCounter.Add(1),
		fn:           fn,
		dependencies: make(map[SignalInterface]uint64),
		immediate:    opts.Immediate,
		defer_:       opts.Defer,
	}
	
	e.active.Store(true)
	
	if opts.Immediate && !opts.Defer {
		e.run()
	} else if opts.Defer {
		scheduleEffect(e)
	}
	
	return e
}

// run executes the effect function
func (e *Effect) run() {
	if !e.active.Load() {
		return
	}
	
	// Prevent recursive runs
	if !e.running.CompareAndSwap(false, true) {
		return
	}
	defer e.running.Store(false)
	
	// Clear dirty flag
	e.dirty.Store(false)
	
	// Clear old dependencies
	e.clearDependencies()
	
	// Push to effect stack for dependency tracking
	pushEffect(e)
	defer popEffect()
	
	// Run cleanups from previous run
	e.runCleanups()
	
	// Execute the effect function
	e.fn()
}

// invalidate marks the effect as needing re-execution
func (e *Effect) invalidate() {
	if !e.active.Load() {
		return
	}
	
	if !e.dirty.CompareAndSwap(false, true) {
		return // Already dirty
	}
	
	// If we're in a batch, queue for later
	if isInBatch() {
		addPendingEffect(e)
		return
	}
	
	if e.defer_ {
		scheduleEffect(e)
	} else {
		e.run()
	}
}

// addDependency registers a signal as a dependency
func (e *Effect) addDependency(signal SignalInterface) {
	e.depmu.Lock()
	defer e.depmu.Unlock()
	e.dependencies[signal] = signal.Version()
}

// removeDependency unregisters a signal dependency
func (e *Effect) removeDependency(signal SignalInterface) {
	e.depmu.Lock()
	defer e.depmu.Unlock()
	delete(e.dependencies, signal)
}

// clearDependencies removes all dependencies
func (e *Effect) clearDependencies() {
	e.depmu.Lock()
	defer e.depmu.Unlock()
	
	for dep := range e.dependencies {
		dep.removeObserver(e)
	}
	e.dependencies = make(map[SignalInterface]uint64)
}

// OnCleanup registers a cleanup function to run before next execution
func (e *Effect) OnCleanup(cleanup func()) {
	e.cleanupmu.Lock()
	defer e.cleanupmu.Unlock()
	e.cleanups = append(e.cleanups, cleanup)
}

// runCleanups executes and clears all cleanup functions
func (e *Effect) runCleanups() {
	e.cleanupmu.Lock()
	cleanups := e.cleanups
	e.cleanups = nil
	e.cleanupmu.Unlock()
	
	for _, cleanup := range cleanups {
		cleanup()
	}
}

// Dispose stops the effect and runs cleanups
func (e *Effect) Dispose() {
	if !e.active.CompareAndSwap(true, false) {
		return // Already disposed
	}
	
	e.clearDependencies()
	e.runCleanups()
}

// IsActive returns whether the effect is still active
func (e *Effect) IsActive() bool {
	return e.active.Load()
}