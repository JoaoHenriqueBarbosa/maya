package reactive

import (
	"sync"
	"time"
)

// UpdateBatcher batches reactive updates for efficiency
type UpdateBatcher struct {
	updates   chan func()
	pending   []func()
	ticker    *time.Ticker
	threshold int
	mu        sync.Mutex
	running   bool
}

// NewUpdateBatcher creates a new update batcher
func NewUpdateBatcher() *UpdateBatcher {
	return &UpdateBatcher{
		updates:   make(chan func(), 100),
		pending:   make([]func(), 0),
		threshold: 10,
	}
}

// Start starts the batching process
func (b *UpdateBatcher) Start() {
	b.mu.Lock()
	if b.running {
		b.mu.Unlock()
		return
	}
	b.running = true
	b.ticker = time.NewTicker(16 * time.Millisecond) // 60 FPS
	b.mu.Unlock()
	
	go b.processBatches()
}

// Stop stops the batching process
func (b *UpdateBatcher) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if b.running && b.ticker != nil {
		b.ticker.Stop()
		b.running = false
	}
}

// Add adds an update to the batch
func (b *UpdateBatcher) Add(update func()) {
	select {
	case b.updates <- update:
	default:
		// Channel full, execute immediately
		update()
	}
}

// processBatches processes batched updates
func (b *UpdateBatcher) processBatches() {
	for {
		select {
		case update := <-b.updates:
			b.mu.Lock()
			b.pending = append(b.pending, update)
			b.mu.Unlock()
			
		case <-b.ticker.C:
			b.mu.Lock()
			if len(b.pending) > 0 {
				batch := b.pending
				b.pending = make([]func(), 0)
				b.mu.Unlock()
				
				// Execute batch
				for _, update := range batch {
					update()
				}
			} else {
				b.mu.Unlock()
			}
		}
	}
}

// Flush immediately processes all pending updates
func (b *UpdateBatcher) Flush() {
	b.mu.Lock()
	batch := b.pending
	b.pending = make([]func(), 0)
	b.mu.Unlock()
	
	for _, update := range batch {
		update()
	}
}