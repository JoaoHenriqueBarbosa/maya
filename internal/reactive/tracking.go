package reactive

import (
	"runtime"
)

// Goroutine-local effect tracking using runtime.Goid
// This is a simplified approach - production would use context.Context

// getGoroutineID returns the current goroutine ID
func getGoroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	
	// Parse "goroutine N [..."
	var id uint64
	inNumber := false
	for i := 0; i < n; i++ {
		if buf[i] >= '0' && buf[i] <= '9' {
			inNumber = true
			id = id*10 + uint64(buf[i]-'0')
		} else if inNumber {
			// We've finished reading the number
			break
		}
	}
	
	if id == 0 {
		// Fallback: use a unique ID based on current time
		// This is not ideal but ensures we don't return 0
		id = uint64(runtime.NumGoroutine())
	}
	return id
}

// effectStackEntry represents a stack entry for nested effects
type effectStackEntry struct {
	effect *Effect
	prev   *effectStackEntry
}

// getCurrentEffect returns the currently executing effect
func getCurrentEffect() *Effect {
	gid := getGoroutineID()
	if v, ok := effectStack.Load(gid); ok {
		if entry, ok := v.(*effectStackEntry); ok && entry != nil {
			return entry.effect
		}
	}
	return nil
}

// pushEffect adds an effect to the current goroutine's stack
func pushEffect(e *Effect) {
	gid := getGoroutineID()
	var prev *effectStackEntry
	if v, ok := effectStack.Load(gid); ok {
		prev, _ = v.(*effectStackEntry)
	}
	
	entry := &effectStackEntry{
		effect: e,
		prev:   prev,
	}
	effectStack.Store(gid, entry)
}

// popEffect removes the top effect from the current goroutine's stack
func popEffect() {
	gid := getGoroutineID()
	if v, ok := effectStack.Load(gid); ok {
		if entry, ok := v.(*effectStackEntry); ok && entry != nil {
			if entry.prev != nil {
				effectStack.Store(gid, entry.prev)
			} else {
				effectStack.Delete(gid)
			}
		}
	}
}

// Untrack runs a function without dependency tracking
func Untrack[T any](fn func() T) T {
	gid := getGoroutineID()
	
	// Save current stack
	var saved *effectStackEntry
	if v, ok := effectStack.Load(gid); ok {
		saved, _ = v.(*effectStackEntry)
	}
	
	// Clear stack temporarily
	effectStack.Delete(gid)
	
	// Run function
	result := fn()
	
	// Restore stack
	if saved != nil {
		effectStack.Store(gid, saved)
	}
	
	return result
}

// UntrackVoid is like Untrack but for functions that don't return a value
func UntrackVoid(fn func()) {
	Untrack(func() any {
		fn()
		return nil
	})
}