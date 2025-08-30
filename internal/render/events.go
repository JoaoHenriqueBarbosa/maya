//go:build wasm
// +build wasm

package render

import (
	"sync"
)

// EventRegistry manages event callbacks without js.FuncOf
type EventRegistry struct {
	callbacks map[int32]func()
	nextID    int32
	mu        sync.RWMutex
}

var globalEventRegistry = &EventRegistry{
	callbacks: make(map[int32]func()),
}

// RegisterCallback registers a callback and returns its ID
func RegisterCallback(callback func()) int32 {
	globalEventRegistry.mu.Lock()
	defer globalEventRegistry.mu.Unlock()
	
	id := globalEventRegistry.nextID
	globalEventRegistry.nextID++
	globalEventRegistry.callbacks[id] = callback
	
	return id
}

// UnregisterCallback removes a callback
func UnregisterCallback(id int32) {
	globalEventRegistry.mu.Lock()
	defer globalEventRegistry.mu.Unlock()
	
	delete(globalEventRegistry.callbacks, id)
}

// GetCallback retrieves a callback by ID
func GetCallback(id int32) func() {
	globalEventRegistry.mu.RLock()
	defer globalEventRegistry.mu.RUnlock()
	
	return globalEventRegistry.callbacks[id]
}

// Note: The actual exported functions are in maya/exports.go
// because go:wasmexport needs to be in the main package