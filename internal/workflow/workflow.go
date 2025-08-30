package workflow

import (
	"context"
	"fmt"
	"iter"
	"sync"
	"sync/atomic"
	"time"

	"github.com/maya-framework/maya/internal/graph"
)

// Stage represents a stage in the workflow
type Stage struct {
	ID          string
	Name        string
	Description string
	Execute     StageFunc
	Timeout     time.Duration
	MaxRetries  int
	Parallel    bool
	MaxWorkers  int
}

// StageFunc is the function executed by a stage
type StageFunc func(context.Context, *StageContext) error

// StageContext provides context for stage execution
type StageContext struct {
	Stage    *Stage
	Input    interface{}
	Output   interface{}
	Metadata map[string]interface{}
	
	// Dependencies results
	Dependencies map[string]interface{}
	
	mu sync.RWMutex
}

// WorkflowEngine orchestrates workflow execution
type WorkflowEngine struct {
	name        string
	description string
	graph       *graph.Graph
	stages      map[string]*Stage
	
	// Execution state
	running     atomic.Bool
	results     sync.Map // map[string]interface{}
	errors      sync.Map // map[string]error
	
	// Metrics
	metrics     *WorkflowMetrics
	
	// Configuration
	maxConcurrency int
	timeout        time.Duration
	
	mu sync.RWMutex
}

// WorkflowMetrics tracks workflow execution metrics
type WorkflowMetrics struct {
	StartTime      time.Time
	EndTime        time.Time
	StageMetrics   map[string]*StageMetrics
	TotalStages    int
	CompletedStages int
	FailedStages   int
	mu             sync.RWMutex
}

// StageMetrics tracks individual stage metrics
type StageMetrics struct {
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	RetryCount   int
	Success      bool
	Error        error
}

// NewWorkflowEngine creates a new workflow engine
func NewWorkflowEngine(name string) *WorkflowEngine {
	return &WorkflowEngine{
		name:           name,
		graph:          graph.NewGraph(),
		stages:         make(map[string]*Stage),
		maxConcurrency: 10,
		timeout:        30 * time.Minute,
		metrics: &WorkflowMetrics{
			StageMetrics: make(map[string]*StageMetrics),
		},
	}
}

// SetDescription sets the workflow description
func (w *WorkflowEngine) SetDescription(desc string) {
	w.description = desc
}

// SetMaxConcurrency sets the maximum concurrent stages
func (w *WorkflowEngine) SetMaxConcurrency(max int) {
	w.maxConcurrency = max
}

// SetTimeout sets the workflow timeout
func (w *WorkflowEngine) SetTimeout(timeout time.Duration) {
	w.timeout = timeout
}

// AddStage adds a stage to the workflow
func (w *WorkflowEngine) AddStage(stage *Stage) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	if _, exists := w.stages[stage.ID]; exists {
		return fmt.Errorf("stage %s already exists", stage.ID)
	}
	
	// Add to graph
	if err := w.graph.AddNode(graph.NodeID(stage.ID), stage); err != nil {
		return err
	}
	
	w.stages[stage.ID] = stage
	return nil
}

// AddDependency creates a dependency between stages
func (w *WorkflowEngine) AddDependency(from, to string) error {
	_, err := w.graph.AddEdge(graph.NodeID(from), graph.NodeID(to), 1.0)
	return err
}

// Execute runs the workflow
func (w *WorkflowEngine) Execute(ctx context.Context, input interface{}) error {
	if !w.running.CompareAndSwap(false, true) {
		return fmt.Errorf("workflow is already running")
	}
	defer w.running.Store(false)
	
	// Create context with timeout
	if w.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, w.timeout)
		defer cancel()
	}
	
	// Initialize metrics
	w.metrics.StartTime = time.Now()
	w.metrics.TotalStages = len(w.stages)
	
	// Execute using graph parallel processing
	err := w.graph.ParallelProcess(ctx, func(ctx context.Context, node *graph.Node) error {
		stage, ok := node.Data.(*Stage)
		if !ok {
			return fmt.Errorf("invalid stage data for node %s", node.ID)
		}
		
		return w.executeStage(ctx, stage, input)
	})
	
	w.metrics.EndTime = time.Now()
	
	return err
}

// executeStage executes a single stage with retry logic
func (w *WorkflowEngine) executeStage(ctx context.Context, stage *Stage, input interface{}) error {
	metrics := &StageMetrics{
		StartTime: time.Now(),
	}
	
	// Store metrics
	w.metrics.mu.Lock()
	w.metrics.StageMetrics[stage.ID] = metrics
	w.metrics.mu.Unlock()
	
	// Create stage context
	stageCtx := &StageContext{
		Stage:        stage,
		Input:        input,
		Metadata:     make(map[string]interface{}),
		Dependencies: w.getDependencyResults(stage.ID),
	}
	
	// Apply timeout if specified
	if stage.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, stage.Timeout)
		defer cancel()
	}
	
	// Execute with retries
	var err error
	for attempt := 0; attempt <= stage.MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			select {
			case <-time.After(time.Second * time.Duration(attempt)):
			case <-ctx.Done():
				err = ctx.Err()
				break
			}
		}
		
		metrics.RetryCount = attempt
		
		// Execute stage
		err = stage.Execute(ctx, stageCtx)
		if err == nil {
			// Success
			w.results.Store(stage.ID, stageCtx.Output)
			metrics.Success = true
			break
		}
	}
	
	// Update metrics
	metrics.EndTime = time.Now()
	metrics.Duration = metrics.EndTime.Sub(metrics.StartTime)
	metrics.Error = err
	
	if err != nil {
		w.errors.Store(stage.ID, err)
		w.metrics.mu.Lock()
		w.metrics.FailedStages++
		w.metrics.mu.Unlock()
		return fmt.Errorf("stage %s failed: %w", stage.ID, err)
	}
	
	w.metrics.mu.Lock()
	w.metrics.CompletedStages++
	w.metrics.mu.Unlock()
	
	return nil
}

// getDependencyResults gets results from dependency stages
func (w *WorkflowEngine) getDependencyResults(stageID string) map[string]interface{} {
	dependencies := w.graph.GetDependencies(graph.NodeID(stageID))
	results := make(map[string]interface{})
	
	for _, depID := range dependencies {
		if value, ok := w.results.Load(string(depID)); ok {
			results[string(depID)] = value
		}
	}
	
	return results
}

// GetResult returns the result of a stage
func (w *WorkflowEngine) GetResult(stageID string) (interface{}, bool) {
	return w.results.Load(stageID)
}

// GetError returns the error of a stage
func (w *WorkflowEngine) GetError(stageID string) (error, bool) {
	if err, ok := w.errors.Load(stageID); ok {
		return err.(error), true
	}
	return nil, false
}

// GetMetrics returns workflow metrics
func (w *WorkflowEngine) GetMetrics() *WorkflowMetrics {
	w.metrics.mu.RLock()
	defer w.metrics.mu.RUnlock()
	
	// Create a copy
	metrics := &WorkflowMetrics{
		StartTime:       w.metrics.StartTime,
		EndTime:         w.metrics.EndTime,
		TotalStages:     w.metrics.TotalStages,
		CompletedStages: w.metrics.CompletedStages,
		FailedStages:    w.metrics.FailedStages,
		StageMetrics:    make(map[string]*StageMetrics),
	}
	
	for id, sm := range w.metrics.StageMetrics {
		metrics.StageMetrics[id] = &StageMetrics{
			StartTime:  sm.StartTime,
			EndTime:    sm.EndTime,
			Duration:   sm.Duration,
			RetryCount: sm.RetryCount,
			Success:    sm.Success,
			Error:      sm.Error,
		}
	}
	
	return metrics
}

// Pipeline represents a linear workflow pipeline
type Pipeline struct {
	name   string
	stages []*PipelineStage
	
	// Execution control
	semaphore chan struct{}
	
	mu sync.RWMutex
}

// PipelineStage represents a stage in a pipeline
type PipelineStage struct {
	Name      string
	Process   func(context.Context, interface{}) (interface{}, error)
	Parallel  bool
	Workers   int
}

// NewPipeline creates a new pipeline
func NewPipeline(name string, maxConcurrency int) *Pipeline {
	return &Pipeline{
		name:      name,
		stages:    make([]*PipelineStage, 0),
		semaphore: make(chan struct{}, maxConcurrency),
	}
}

// AddStage adds a stage to the pipeline
func (p *Pipeline) AddStage(stage *PipelineStage) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stages = append(p.stages, stage)
}

// Execute runs the pipeline
func (p *Pipeline) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	current := input
	
	for i, stage := range p.stages {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		result, err := p.executeStage(ctx, stage, current)
		if err != nil {
			return nil, fmt.Errorf("stage %d (%s) failed: %w", i, stage.Name, err)
		}
		
		current = result
	}
	
	return current, nil
}

// executeStage executes a single pipeline stage
func (p *Pipeline) executeStage(ctx context.Context, stage *PipelineStage, input interface{}) (interface{}, error) {
	if !stage.Parallel {
		return stage.Process(ctx, input)
	}
	
	// Parallel processing for collection inputs
	if items, ok := input.([]interface{}); ok {
		return p.processParallel(ctx, stage, items)
	}
	
	return stage.Process(ctx, input)
}

// processParallel processes items in parallel
func (p *Pipeline) processParallel(ctx context.Context, stage *PipelineStage, items []interface{}) (interface{}, error) {
	results := make([]interface{}, len(items))
	errChan := make(chan error, len(items))
	
	workers := stage.Workers
	if workers <= 0 {
		workers = len(items)
	}
	if workers > len(items) {
		workers = len(items)
	}
	
	// Work queue
	workChan := make(chan int, len(items))
	for i := range items {
		workChan <- i
	}
	close(workChan)
	
	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for i := range workChan {
				select {
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				default:
				}
				
				result, err := stage.Process(ctx, items[i])
				if err != nil {
					errChan <- err
					return
				}
				
				results[i] = result
			}
		}()
	}
	
	// Wait for completion
	wg.Wait()
	close(errChan)
	
	// Check for errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}
	
	return results, nil
}

// Stream provides streaming workflow processing
type Stream struct {
	name      string
	processor StreamProcessor
	
	// Channels
	input  chan interface{}
	output chan interface{}
	errors chan error
	
	// Control
	running atomic.Bool
	done    chan struct{}
	
	// Configuration
	bufferSize int
	workers    int
}

// StreamProcessor processes stream items
type StreamProcessor func(context.Context, interface{}) (interface{}, error)

// NewStream creates a new stream processor
func NewStream(name string, processor StreamProcessor) *Stream {
	return &Stream{
		name:       name,
		processor:  processor,
		bufferSize: 100,
		workers:    1,
	}
}

// SetBufferSize sets the buffer size for channels
func (s *Stream) SetBufferSize(size int) {
	s.bufferSize = size
}

// SetWorkers sets the number of workers
func (s *Stream) SetWorkers(workers int) {
	s.workers = workers
}

// Start starts the stream processing
func (s *Stream) Start(ctx context.Context) error {
	if !s.running.CompareAndSwap(false, true) {
		return fmt.Errorf("stream is already running")
	}
	
	// Initialize channels
	s.input = make(chan interface{}, s.bufferSize)
	s.output = make(chan interface{}, s.bufferSize)
	s.errors = make(chan error, s.bufferSize)
	s.done = make(chan struct{})
	
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.worker(ctx)
		}()
	}
	
	// Cleanup goroutine
	go func() {
		wg.Wait()
		close(s.output)
		close(s.errors)
		close(s.done)
		s.running.Store(false)
	}()
	
	return nil
}

// worker processes items from the input channel
func (s *Stream) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-s.input:
			if !ok {
				return
			}
			
			result, err := s.processor(ctx, item)
			if err != nil {
				select {
				case s.errors <- err:
				case <-ctx.Done():
					return
				}
				continue
			}
			
			select {
			case s.output <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

// Send sends an item to the stream
func (s *Stream) Send(ctx context.Context, item interface{}) error {
	if !s.running.Load() {
		return fmt.Errorf("stream is not running")
	}
	
	select {
	case s.input <- item:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Receive receives a processed item from the stream
func (s *Stream) Receive(ctx context.Context) (interface{}, error) {
	select {
	case result, ok := <-s.output:
		if !ok {
			return nil, fmt.Errorf("stream closed")
		}
		return result, nil
	case err := <-s.errors:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Close closes the input channel and waits for processing to complete
func (s *Stream) Close() {
	close(s.input)
	<-s.done
}

// Process provides an iterator for stream processing
func (s *Stream) Process(ctx context.Context, items iter.Seq[interface{}]) iter.Seq2[interface{}, error] {
	return func(yield func(interface{}, error) bool) {
		// Start the stream
		if err := s.Start(ctx); err != nil {
			yield(nil, err)
			return
		}
		defer s.Close()
		
		// Send items
		go func() {
			for item := range items {
				if err := s.Send(ctx, item); err != nil {
					return
				}
			}
			close(s.input)
		}()
		
		// Receive results
		for {
			result, err := s.Receive(ctx)
			if err != nil {
				if err.Error() == "stream closed" {
					return
				}
				if !yield(nil, err) {
					return
				}
				continue
			}
			
			if !yield(result, nil) {
				return
			}
		}
	}
}