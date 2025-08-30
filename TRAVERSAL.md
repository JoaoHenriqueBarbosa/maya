# Sistema de UI Completo em Go/WASM

## 1. Arquitetura Geral

### Camadas do Sistema
```
┌─────────────────────────────────────────┐
│         API de Alto Nível               │
│    (Widgets, Declarative DSL)           │
├─────────────────────────────────────────┤
│         Sistema de Reatividade          │
│    (State, Signals, Effects)            │
├─────────────────────────────────────────┤
│         Layout Engine                   │
│    (Flexbox, Grid, Constraints)         │
├─────────────────────────────────────────┤
│         Rendering Engine                │
│    (Multi-pass, Dirty Checking)         │
├─────────────────────────────────────────┤
│         Canvas/WebGL Backend            │
│    (Drawing Commands, GPU Accel)        │
└─────────────────────────────────────────┘
```

## 2. Algoritmos de Traversal da Árvore

### 2.1 Estruturas Base da Árvore

```go
// Node representa um nó na árvore de widgets
type Node struct {
    ID         NodeID
    Widget     Widget
    Parent     *Node
    Children   []*Node
    Depth      int
    Index      int // índice entre siblings
    
    // Cache para otimizações
    cachedSize     Size
    cachedPosition Offset
    isDirty        bool
    dirtyFlags     DirtyFlags
}

// Tree representa a árvore completa
type Tree struct {
    root        *Node
    nodeMap     map[NodeID]*Node
    dirtyNodes  *PriorityQueue // heap de nós dirty ordenados por profundidade
    version     uint64
}
```

### 2.2 Breadth-First Search (BFS)

```go
// BFS é usado quando precisamos processar níveis uniformemente
// Ideal para: propagação de constraints, temas, contexto
type BFSTraversal struct {
    queue *Queue[*Node]
}

func (t *Tree) BreadthFirstTraversal(visitor func(*Node) error) error {
    if t.root == nil {
        return nil
    }
    
    queue := NewQueue[*Node]()
    queue.Enqueue(t.root)
    
    // Processamento nível por nível
    currentDepth := 0
    nodesInCurrentLevel := 1
    nodesInNextLevel := 0
    
    for !queue.IsEmpty() {
        node := queue.Dequeue()
        
        // Callback com informação de nível
        if err := visitor(node); err != nil {
            return err
        }
        
        // Adiciona filhos à fila
        for _, child := range node.Children {
            queue.Enqueue(child)
            nodesInNextLevel++
        }
        
        nodesInCurrentLevel--
        if nodesInCurrentLevel == 0 {
            // Mudança de nível - útil para sincronização
            currentDepth++
            nodesInCurrentLevel = nodesInNextLevel
            nodesInNextLevel = 0
        }
    }
    
    return nil
}

// BFS otimizado com early termination
func (t *Tree) BFSWithPredicate(
    predicate func(*Node) bool,
    visitor func(*Node) error,
) error {
    queue := NewQueue[*Node]()
    queue.Enqueue(t.root)
    visited := make(map[NodeID]bool)
    
    for !queue.IsEmpty() {
        node := queue.Dequeue()
        
        if visited[node.ID] {
            continue
        }
        visited[node.ID] = true
        
        // Early termination se predicado falhar
        if !predicate(node) {
            continue
        }
        
        if err := visitor(node); err != nil {
            return err
        }
        
        for _, child := range node.Children {
            if !visited[child.ID] {
                queue.Enqueue(child)
            }
        }
    }
    
    return nil
}
```

### 2.3 Depth-First Search (DFS)

```go
// DFS é usado para operações que precisam de contexto completo do caminho
// Ideal para: cálculo de tamanhos, painting, hit testing
type DFSTraversal struct {
    stack *Stack[*Node]
}

// Pre-order DFS (processa nó antes dos filhos)
// Usado para: aplicar transformações, propagar constraints
func (t *Tree) PreOrderDFS(visitor func(*Node, int) error) error {
    return t.preOrderDFSRecursive(t.root, visitor, 0)
}

func (t *Tree) preOrderDFSRecursive(
    node *Node,
    visitor func(*Node, int) error,
    depth int,
) error {
    if node == nil {
        return nil
    }
    
    // Processa nó atual
    if err := visitor(node, depth); err != nil {
        return err
    }
    
    // Processa filhos
    for _, child := range node.Children {
        if err := t.preOrderDFSRecursive(child, visitor, depth+1); err != nil {
            return err
        }
    }
    
    return nil
}

// Post-order DFS (processa filhos antes do nó)
// Usado para: calcular tamanhos intrínsecos, liberar recursos
func (t *Tree) PostOrderDFS(visitor func(*Node, int) error) error {
    return t.postOrderDFSRecursive(t.root, visitor, 0)
}

func (t *Tree) postOrderDFSRecursive(
    node *Node,
    visitor func(*Node, int) error,
    depth int,
) error {
    if node == nil {
        return nil
    }
    
    // Processa filhos primeiro
    for _, child := range node.Children {
        if err := t.postOrderDFSRecursive(child, visitor, depth+1); err != nil {
            return err
        }
    }
    
    // Depois processa nó atual
    return visitor(node, depth)
}

// In-order DFS para árvores binárias (usado em BST de layers)
func (t *Tree) InOrderDFS(visitor func(*Node) error) error {
    return t.inOrderDFSRecursive(t.root, visitor)
}
```

### 2.4 Priority-Based Traversal

```go
// Traversal baseado em prioridade para dirty checking
// Nós mais profundos são processados primeiro (bottom-up)
type PriorityTraversal struct {
    heap *MinHeap[*Node]
}

func (t *Tree) PriorityTraversal(
    priority func(*Node) int,
    visitor func(*Node) error,
) error {
    heap := NewMinHeap[*Node](func(a, b *Node) bool {
        return priority(a) < priority(b)
    })
    
    // Coleta todos os nós
    t.PreOrderDFS(func(node *Node, depth int) error {
        heap.Push(node)
        return nil
    })
    
    // Processa em ordem de prioridade
    for !heap.IsEmpty() {
        node := heap.Pop()
        if err := visitor(node); err != nil {
            return err
        }
    }
    
    return nil
}
```

## 3. Layout Engine com Multi-pass

### 3.1 Layout Algorithm Pipeline

```go
type LayoutPipeline struct {
    tree           *Tree
    constraintSolver *ConstraintSolver
    flexSolver      *FlexboxSolver
    gridSolver      *GridSolver
}

// Pipeline completo de layout em múltiplas passadas
func (p *LayoutPipeline) PerformLayout(rootConstraints Constraints) error {
    // Pass 1: Marca nós dirty e propaga (BFS)
    if err := p.markDirtyNodes(); err != nil {
        return err
    }
    
    // Pass 2: Calcula dimensões intrínsecas (Post-order DFS)
    if err := p.calculateIntrinsicDimensions(); err != nil {
        return err
    }
    
    // Pass 3: Resolve constraints (Pre-order DFS)
    if err := p.resolveConstraints(rootConstraints); err != nil {
        return err
    }
    
    // Pass 4: Calcula tamanhos finais (Post-order DFS)
    if err := p.calculateFinalSizes(); err != nil {
        return err
    }
    
    // Pass 5: Atribui posições (Pre-order DFS)
    if err := p.assignPositions(); err != nil {
        return err
    }
    
    // Pass 6: Alinhamento de baseline (BFS por nível)
    if err := p.alignBaselines(); err != nil {
        return err
    }
    
    return nil
}
```

### 3.2 Dirty Marking com Propagação

```go
// Marca nós como dirty e propaga para ancestrais
func (p *LayoutPipeline) markDirtyNodes() error {
    // Usa BFS para propagar dirty flags uniformemente
    return p.tree.BreadthFirstTraversal(func(node *Node) error {
        if node.isDirty {
            // Propaga para ancestrais
            p.propagateDirtyToAncestors(node)
        }
        return nil
    })
}

func (p *LayoutPipeline) propagateDirtyToAncestors(node *Node) {
    current := node.Parent
    for current != nil {
        if current.dirtyFlags&LayoutDirty != 0 {
            break // Já marcado, pode parar
        }
        current.dirtyFlags |= LayoutDirty
        current = current.Parent
    }
}
```

### 3.3 Cálculo de Dimensões Intrínsecas (Bottom-up)

```go
func (p *LayoutPipeline) calculateIntrinsicDimensions() error {
    // Post-order DFS para calcular de baixo para cima
    return p.tree.PostOrderDFS(func(node *Node, depth int) error {
        widget := node.Widget
        
        // Calcula baseado nos filhos (que já foram processados)
        childrenWidths := make([]float64, len(node.Children))
        childrenHeights := make([]float64, len(node.Children))
        
        for i, child := range node.Children {
            childrenWidths[i] = child.Widget.GetIntrinsicWidth(math.Inf(1))
            childrenHeights[i] = child.Widget.GetIntrinsicHeight(math.Inf(1))
        }
        
        // Calcula as próprias dimensões intrínsecas
        node.intrinsicWidth = widget.CalculateIntrinsicWidth(childrenWidths)
        node.intrinsicHeight = widget.CalculateIntrinsicHeight(childrenHeights)
        
        return nil
    })
}
```

### 3.4 Resolução de Constraints (Top-down)

```go
func (p *LayoutPipeline) resolveConstraints(rootConstraints Constraints) error {
    // Pre-order DFS para propagar constraints de cima para baixo
    return p.tree.PreOrderDFS(func(node *Node, depth int) error {
        var constraints Constraints
        
        if node.Parent == nil {
            constraints = rootConstraints
        } else {
            // Pega constraints do pai e ajusta para este widget
            parentConstraints := node.Parent.resolvedConstraints
            constraints = node.Widget.AdjustConstraints(parentConstraints)
        }
        
        node.resolvedConstraints = constraints
        
        // Para layouts especiais (flex, grid), resolve constraints dos filhos
        switch widget := node.Widget.(type) {
        case *FlexWidget:
            p.flexSolver.ResolveFlexConstraints(node, constraints)
        case *GridWidget:
            p.gridSolver.ResolveGridConstraints(node, constraints)
        }
        
        return nil
    })
}
```

## 4. Algoritmos Especializados para Layout

### 4.1 Flexbox Algorithm

```go
type FlexboxSolver struct {
    mainAxisSize  float64
    crossAxisSize float64
    direction     FlexDirection
}

func (s *FlexboxSolver) ResolveFlexConstraints(node *Node, constraints Constraints) {
    flex := node.Widget.(*FlexWidget)
    children := node.Children
    
    // Fase 1: Calcula tamanhos base e collect flex factors
    var totalFlex float64
    var inflexibleSpace float64
    flexChildren := make([]*FlexChild, 0)
    
    for _, child := range children {
        if flexChild, ok := child.Widget.(*FlexChild); ok && flexChild.Flex > 0 {
            flexChildren = append(flexChildren, flexChild)
            totalFlex += flexChild.Flex
        } else {
            // Calcula tamanho do filho inflexível
            size := s.calculateInflexibleSize(child, constraints)
            inflexibleSpace += s.getMainAxisSize(size)
        }
    }
    
    // Fase 2: Distribui espaço disponível
    availableSpace := s.getMainAxisConstraint(constraints) - inflexibleSpace
    if totalFlex > 0 && availableSpace > 0 {
        flexUnit := availableSpace / totalFlex
        
        for _, flexChild := range flexChildren {
            allocatedSpace := flexUnit * flexChild.Flex
            s.assignFlexChildConstraints(flexChild, allocatedSpace)
        }
    }
    
    // Fase 3: Resolve alinhamento no cross axis
    s.resolveCrossAxisAlignment(node, constraints)
}
```

### 4.2 Grid Layout Algorithm

```go
type GridSolver struct {
    columns []GridTrack
    rows    []GridTrack
    areas   map[string]GridArea
}

func (s *GridSolver) ResolveGridConstraints(node *Node, constraints Constraints) {
    grid := node.Widget.(*GridWidget)
    
    // Fase 1: Calcula tracks sizes
    columnSizes := s.calculateTrackSizes(grid.Columns, constraints.MaxWidth)
    rowSizes := s.calculateTrackSizes(grid.Rows, constraints.MaxHeight)
    
    // Fase 2: Posiciona items usando placement algorithm
    placements := s.autoPlacement(node.Children, len(columnSizes), len(rowSizes))
    
    // Fase 3: Resolve spanning items
    for _, placement := range placements {
        if placement.ColumnSpan > 1 || placement.RowSpan > 1 {
            s.resolveSpanning(placement, columnSizes, rowSizes)
        }
    }
    
    // Fase 4: Atribui constraints finais
    for i, child := range node.Children {
        placement := placements[i]
        childConstraints := s.getConstraintsForPlacement(
            placement,
            columnSizes,
            rowSizes,
        )
        child.resolvedConstraints = childConstraints
    }
}

// Auto-placement algorithm (similar ao CSS Grid)
func (s *GridSolver) autoPlacement(
    children []*Node,
    cols, rows int,
) []GridPlacement {
    grid := make([][]bool, rows)
    for i := range grid {
        grid[i] = make([]bool, cols)
    }
    
    placements := make([]GridPlacement, len(children))
    cursor := GridPosition{Row: 0, Col: 0}
    
    for i, child := range children {
        // Find next available position
        placement := s.findNextAvailable(grid, cursor, child)
        placements[i] = placement
        
        // Mark cells as occupied
        s.markOccupied(grid, placement)
        
        // Update cursor
        cursor = s.updateCursor(cursor, placement, cols)
    }
    
    return placements
}
```

## 5. Algoritmos de Rendering

### 5.1 Dirty Rectangle Algorithm

```go
type DirtyRectTracker struct {
    dirtyRegions []Rect
    framebuffer  *Framebuffer
}

// Coalesce overlapping dirty rectangles
func (d *DirtyRectTracker) CoalesceDirtyRects() []Rect {
    if len(d.dirtyRegions) == 0 {
        return nil
    }
    
    // Sort by area para otimizar merging
    sort.Slice(d.dirtyRegions, func(i, j int) bool {
        return d.dirtyRegions[i].Area() > d.dirtyRegions[j].Area()
    })
    
    merged := make([]Rect, 0)
    merged = append(merged, d.dirtyRegions[0])
    
    for i := 1; i < len(d.dirtyRegions); i++ {
        rect := d.dirtyRegions[i]
        wasMerged := false
        
        for j := range merged {
            if merged[j].Intersects(rect) {
                // Merge rectangles
                merged[j] = merged[j].Union(rect)
                wasMerged = true
                break
            }
        }
        
        if !wasMerged {
            merged = append(merged, rect)
        }
    }
    
    return merged
}
```

### 5.2 Z-Order Sorting (Painter's Algorithm)

```go
type LayerTree struct {
    root   *Layer
    zIndex map[*Layer]int
}

// Topological sort para rendering order
func (t *LayerTree) GetRenderOrder() []*Layer {
    layers := make([]*Layer, 0)
    visited := make(map[*Layer]bool)
    
    // DFS modificado para respeitar z-index
    var visit func(*Layer)
    visit = func(layer *Layer) {
        if visited[layer] {
            return
        }
        visited[layer] = true
        
        // Ordena filhos por z-index
        children := make([]*Layer, len(layer.Children))
        copy(children, layer.Children)
        sort.Slice(children, func(i, j int) bool {
            return t.zIndex[children[i]] < t.zIndex[children[j]]
        })
        
        // Visita filhos em ordem
        for _, child := range children {
            visit(child)
        }
        
        layers = append(layers, layer)
    }
    
    visit(t.root)
    return layers
}
```

### 5.3 Occlusion Culling

```go
type OcclusionCuller struct {
    viewport Rect
    zBuffer  [][]float64
}

// Quadtree para spatial partitioning
type Quadtree struct {
    bounds   Rect
    nodes    []*Node
    children [4]*Quadtree // NW, NE, SW, SE
    maxNodes int
    maxDepth int
    depth    int
}

func (q *Quadtree) Insert(node *Node) {
    if !q.bounds.Contains(node.Bounds()) {
        return
    }
    
    if len(q.nodes) < q.maxNodes || q.depth >= q.maxDepth {
        q.nodes = append(q.nodes, node)
        return
    }
    
    // Subdivide if necessary
    if q.children[0] == nil {
        q.subdivide()
    }
    
    // Insert into children
    for i := range q.children {
        q.children[i].Insert(node)
    }
}

func (q *Quadtree) Query(rect Rect) []*Node {
    result := make([]*Node, 0)
    
    if !q.bounds.Intersects(rect) {
        return result
    }
    
    for _, node := range q.nodes {
        if rect.Intersects(node.Bounds()) {
            result = append(result, node)
        }
    }
    
    if q.children[0] != nil {
        for i := range q.children {
            childResults := q.children[i].Query(rect)
            result = append(result, childResults...)
        }
    }
    
    return result
}
```

## 6. Hit Testing Algorithm

```go
type HitTester struct {
    tree     *Tree
    quadtree *Quadtree
}

// Hit testing com early termination
func (h *HitTester) HitTest(point Point) *Node {
    // Primeiro usa quadtree para filtering espacial
    candidates := h.quadtree.Query(Rect{
        X:      point.X - 1,
        Y:      point.Y - 1,
        Width:  2,
        Height: 2,
    })
    
    // Ordena por z-index (reverso para pegar o top-most)
    sort.Slice(candidates, func(i, j int) bool {
        return candidates[i].ZIndex > candidates[j].ZIndex
    })
    
    // Testa cada candidato
    for _, node := range candidates {
        if h.hitTestNode(node, point) {
            // Testa filhos recursivamente (DFS)
            if child := h.hitTestChildren(node, point); child != nil {
                return child
            }
            return node
        }
    }
    
    return nil
}

func (h *HitTester) hitTestChildren(node *Node, point Point) *Node {
    // Itera filhos em ordem reversa (top-most primeiro)
    for i := len(node.Children) - 1; i >= 0; i-- {
        child := node.Children[i]
        
        // Transforma ponto para coordenadas locais do filho
        localPoint := point.Transform(child.Transform.Inverse())
        
        if h.hitTestNode(child, localPoint) {
            // Recursão para testar filhos do filho
            if grandchild := h.hitTestChildren(child, localPoint); grandchild != nil {
                return grandchild
            }
            return child
        }
    }
    
    return nil
}
```

## 7. Sistema de Reatividade com Dependency Tracking

### 7.1 Dependency Graph

```go
type DependencyGraph struct {
    nodes    map[NodeID]*DependencyNode
    edges    map[NodeID][]NodeID
    circular map[NodeID]bool
}

type DependencyNode struct {
    ID           NodeID
    Dependencies []NodeID
    Dependents   []NodeID
    Value        any
    Computer     func() any
    Version      uint64
}

// Topological sort para ordem de atualização
func (g *DependencyGraph) GetUpdateOrder() ([]NodeID, error) {
    visited := make(map[NodeID]bool)
    recStack := make(map[NodeID]bool)
    result := make([]NodeID, 0)
    
    var visit func(NodeID) error
    visit = func(id NodeID) error {
        visited[id] = true
        recStack[id] = true
        
        for _, dep := range g.edges[id] {
            if !visited[dep] {
                if err := visit(dep); err != nil {
                    return err
                }
            } else if recStack[dep] {
                return fmt.Errorf("circular dependency detected: %v -> %v", id, dep)
            }
        }
        
        recStack[id] = false
        result = append([]NodeID{id}, result...) // Prepend
        return nil
    }
    
    for id := range g.nodes {
        if !visited[id] {
            if err := visit(id); err != nil {
                return nil, err
            }
        }
    }
    
    return result, nil
}
```

### 7.2 Batch Updates Algorithm

```go
type UpdateBatcher struct {
    pending     map[NodeID]*Update
    queue       *PriorityQueue
    isScheduled bool
    mu          sync.Mutex
}

func (b *UpdateBatcher) ScheduleUpdate(update *Update) {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    // Coalesce updates para o mesmo nó
    if existing, ok := b.pending[update.NodeID]; ok {
        existing.Merge(update)
    } else {
        b.pending[update.NodeID] = update
        b.queue.Push(update)
    }
    
    if !b.isScheduled {
        b.isScheduled = true
        requestAnimationFrame(b.flush)
    }
}

func (b *UpdateBatcher) flush() {
    b.mu.Lock()
    updates := make([]*Update, 0, len(b.pending))
    for !b.queue.IsEmpty() {
        updates = append(updates, b.queue.Pop().(*Update))
    }
    b.pending = make(map[NodeID]*Update)
    b.isScheduled = false
    b.mu.Unlock()
    
    // Processa updates em batch
    b.processUpdates(updates)
}
```

## 8. Memory Management e Recycling

### 8.1 Object Pool com Generics

```go
type Pool[T any] struct {
    pool     chan T
    factory  func() T
    reset    func(T) T
    maxSize  int
}

func NewPool[T any](maxSize int, factory func() T, reset func(T) T) *Pool[T] {
    return &Pool[T]{
        pool:    make(chan T, maxSize),
        factory: factory,
        reset:   reset,
        maxSize: maxSize,
    }
}

func (p *Pool[T]) Get() T {
    select {
    case item := <-p.pool:
        return p.reset(item)
    default:
        return p.factory()
    }
}

func (p *Pool[T]) Put(item T) {
    select {
    case p.pool <- item:
        // Successfully returned to pool
    default:
        // Pool is full, let GC handle it
    }
}

// Uso específico para nós
var nodePool = NewPool[*Node](1000,
    func() *Node { return &Node{} },
    func(n *Node) *Node {
        n.Children = n.Children[:0]
        n.isDirty = false
        n.dirtyFlags = 0
        return n
    },
)
```

### 8.2 Weak References para Cache

```go
type WeakRef[T any] struct {
    ptr    uintptr
    marker *T
}

type WeakCache[K comparable, V any] struct {
    entries map[K]*WeakRef[V]
    mu      sync.RWMutex
}

func (c *WeakCache[K, V]) Get(key K) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if ref, ok := c.entries[key]; ok {
        if val := ref.Get(); val != nil {
            return *val, true
        }
        // Reference foi coletada pelo GC
        delete(c.entries, key)
    }
    
    var zero V
    return zero, false
}
```

## 9. Design System Agnóstico

### 9.1 Sistema de Temas Flexível

```go
// Theme system genérico sem dependência de design específico
type Theme struct {
    Name       string
    Variables  map[string]any
    Tokens     *DesignTokens
    Components map[ComponentType]ComponentTheme
}

type DesignTokens struct {
    // Spacing
    Spacing    SpacingScale
    
    // Typography
    Typography TypographyScale
    
    // Colors
    Colors     ColorPalette
    
    // Borders
    Borders    BorderScale
    
    // Shadows
    Shadows    ShadowScale
    
    // Motion
    Motion     MotionScale
}

// Sistema de design tokens escalável
type SpacingScale struct {
    Base     float64
    Scale    []float64 // [0.25, 0.5, 1, 1.5, 2, 3, 4, 6, 8, 12, 16]
}

func (s SpacingScale) Get(level int) float64 {
    if level < 0 || level >= len(s.Scale) {
        return s.Base
    }
    return s.Base * s.Scale[level]
}
```

### 9.2 Component Factory Pattern

```go
// Factory para criar componentes com temas customizados
type ComponentFactory struct {
    builders map[ComponentType]ComponentBuilder
    theme    *Theme
}

type ComponentBuilder func(props Props, theme *Theme) Widget

func (f *ComponentFactory) Register(
    componentType ComponentType,
    builder ComponentBuilder,
) {
    f.builders[componentType] = builder
}

func (f *ComponentFactory) Build(
    componentType ComponentType,
    props Props,
) Widget {
    if builder, ok := f.builders[componentType]; ok {
        return builder(props, f.theme)
    }
    return nil
}

// Exemplo de uso
factory.Register(ButtonType, func(props Props, theme *Theme) Widget {
    return &Button{
        Style: theme.GetComponentStyle(ButtonType),
        Props: props,
    }
})
```

## 10. Performance Profiling

### 10.1 Frame Budget Tracking

```go
type FrameProfiler struct {
    targetFPS    int
    frameBudget  time.Duration
    measurements []FrameMeasurement
    current      *FrameMeasurement
}

type FrameMeasurement struct {
    StartTime    time.Time
    LayoutTime   time.Duration
    PaintTime    time.Duration
    CompositeTime time.Duration
    TotalTime    time.Duration
    DroppedFrame bool
}

func (p *FrameProfiler) StartFrame() {
    p.current = &FrameMeasurement{
        StartTime: time.Now(),
    }
}

func (p *FrameProfiler) EndPhase(phase FramePhase) {
    elapsed := time.Since(p.current.StartTime)
    
    switch phase {
    case LayoutPhase:
        p.current.LayoutTime = elapsed
    case PaintPhase:
        p.current.PaintTime = elapsed - p.current.LayoutTime
    case CompositePhase:
        p.current.CompositeTime = elapsed - p.current.PaintTime - p.current.LayoutTime
    }
}

func (p *FrameProfiler) EndFrame() {
    p.current.TotalTime = time.Since(p.current.StartTime)
    p.current.DroppedFrame = p.current.TotalTime > p.frameBudget
    
    p.measurements = append(p.measurements, *p.current)
    
    // Keep only last N measurements
    if len(p.measurements) > 1000 {
        p.measurements = p.measurements[1:]
    }
}
```

Este sistema fornece uma base sólida e agnóstica de design específico, com algoritmos otimizados de traversal e layout para construir uma UI framework completa em Go/WASM.