# Maya Reactive System - Teoria e Implementação

## 📚 Fundamentos Teóricos (Agosto 2025)

### Fine-Grained Reactivity

Fine-grained reactivity é um paradigma onde mudanças de estado propagam automaticamente apenas para os consumidores específicos daquele estado, sem re-renderizar componentes inteiros. Inspirado em SolidJS, este sistema oferece:

- **O(1) updates** ao invés de O(n) diffing do Virtual DOM
- **Tracking automático** de dependências
- **Sincronização garantida** - nunca observa estado inconsistente
- **Execução síncrona** - mudanças aplicam imediatamente

### Os 3 Primitivos Fundamentais

#### 1. Signals (Estado Observável)
```go
// Signal é um valor reativo com getter e setter
type Signal[T any] struct {
    value    T
    version  atomic.Uint64
    observers map[uint64]*Effect
    equals   func(a, b T) bool
}
```

**Características:**
- Event emitters otimizados
- Getters executam código arbitrário para tracking
- Setters notificam observers automaticamente
- Versionamento para otimização

#### 2. Effects/Reactions (Computações Reativas)
```go
// Effect é uma função que re-executa quando suas dependências mudam
type Effect struct {
    fn           func()
    dependencies []*Signal
    cleanups     []func()
    state        EffectState // CLEAN | STALE | PENDING
}
```

**Características:**
- Auto-tracking de dependências na primeira execução
- Re-execução automática quando signals mudam
- Cleanup automático entre execuções
- Prevenção de loops infinitos

#### 3. Memos/Derivations (Valores Computados)
```go
// Memo é um signal derivado com cache
type Memo[T any] struct {
    Signal[T]
    fn       func() T
    stale    bool
}
```

**Características:**
- Lazy evaluation - só calcula quando lido
- Caching automático do resultado
- Próprio grafo de dependências
- Garantia de consistência

## 🔄 Ciclo de Vida Reativo

### 1. Fase de Tracking (Coleta de Dependências)

Quando um Effect executa, automaticamente registra todos os Signals lidos:

```go
// Pseudo-código do tracking
var currentEffect *Effect

func (s *Signal[T]) Get() T {
    if currentEffect != nil {
        // Registra este signal como dependência
        currentEffect.addDependency(s)
        s.addObserver(currentEffect)
    }
    return s.value
}
```

### 2. Fase de Invalidação

Quando um Signal muda, marca seus observers como STALE:

```go
func (s *Signal[T]) Set(value T) {
    if !s.equals(s.value, value) {
        s.value = value
        s.version++
        
        // Marca observers como stale
        for _, observer := range s.observers {
            observer.invalidate()
        }
    }
}
```

### 3. Fase de Re-execução

Effects marcados como STALE são re-executados em ordem topológica:

```go
func (e *Effect) invalidate() {
    if e.state == CLEAN {
        e.state = STALE
        scheduler.enqueue(e)
    }
}
```

## 🎯 Características Críticas

### Glitch-Free (Sem Inconsistências)

O sistema garante que nunca observamos estado inconsistente:

```go
// PROBLEMA: Estado inconsistente
a := Signal(1)
b := Signal(2)
sum := Memo(func() int { return a.Get() + b.Get() })

// Se a=2 e b=3 aplicados "simultaneamente"
// NUNCA devemos ver sum=4 (a novo + b velho)
```

**Solução:** Batching e execução síncrona

```go
func Batch(fn func()) {
    startBatch()
    defer endBatch() // Aplica todas mudanças atomicamente
    fn()
}
```

### Dynamic Dependencies (Dependências Dinâmicas)

Dependências são reconstruídas a cada execução:

```go
showFull := Signal(true)
firstName := Signal("John")
lastName := Signal("Doe")

name := Memo(func() string {
    if !showFull.Get() {
        return firstName.Get() // lastName NÃO é dependência aqui
    }
    return firstName.Get() + " " + lastName.Get()
})

// Quando showFull=false, mudanças em lastName NÃO triggeram name
```

### Cleanup Automático

Limpeza de recursos entre execuções:

```go
func CreateEffect(fn func()) {
    effect := &Effect{fn: fn}
    
    execute := func() {
        // Limpa dependências antigas
        effect.cleanupDependencies()
        
        // Executa com tracking
        prevEffect := currentEffect
        currentEffect = effect
        defer func() { currentEffect = prevEffect }()
        
        fn()
    }
    
    execute() // Execução inicial
}
```

## 💫 Padrões Avançados

### 1. Untrack (Leitura sem Tracking)

```go
func Untrack[T any](fn func() T) T {
    prev := currentEffect
    currentEffect = nil
    defer func() { currentEffect = prev }()
    return fn()
}

// Uso: ler signal sem criar dependência
effect := CreateEffect(func() {
    important := signal1.Get() // tracked
    debug := Untrack(func() string {
        return debugSignal.Get() // NOT tracked
    })
})
```

### 2. Batch Updates

```go
func Batch(updates func()) {
    if isInBatch() {
        updates() // Nested batch, apenas executa
        return
    }
    
    startBatch()
    updates()
    flushBatch() // Aplica todas mudanças de uma vez
}
```

### 3. Lazy Memos

```go
type LazyMemo[T any] struct {
    fn    func() T
    cache *T
    deps  []*Signal
}

func (m *LazyMemo[T]) Get() T {
    if m.cache == nil || m.isStale() {
        m.cache = &m.fn()
        m.updateDeps()
    }
    return *m.cache
}
```

## 🚀 Implementação para Maya

### Arquitetura Proposta

```go
package reactive

// Signal - Estado observável
type Signal[T comparable] struct {
    value     T
    version   atomic.Uint64
    mu        sync.RWMutex
    observers map[uint64]*Effect
    equals    func(a, b T) bool
}

// Effect - Computação reativa
type Effect struct {
    id           uint64
    fn           func()
    dependencies map[*SignalInterface]struct{}
    state        atomic.Uint32 // 0=CLEAN, 1=STALE, 2=RUNNING
    cleanup      []func()
}

// Memo - Valor derivado com cache
type Memo[T any] struct {
    signal *Signal[T]
    effect *Effect
}

// UpdateBatcher - Agrupa atualizações
type UpdateBatcher struct {
    pending  []*Effect
    flushing atomic.Bool
}
```

### Integração com Render Pipeline

```go
// 1. Widget cria signals para estado
counter := reactive.NewSignal(0)

// 2. TextSignal cria widget reativo
text := widgets.NewText("counter", "")

// 3. Effect conecta signal ao widget E ao render
reactive.CreateEffect(func() {
    value := counter.Get()
    text.SetText(fmt.Sprintf("Count: %d", value))
    
    // CRITICAL: Agendar re-render do DOM
    if app != nil {
        app.scheduleRender()
    }
})

// 4. Render pipeline reconstrói árvore quando necessário
func (app *App) scheduleRender() {
    app.batcher.Add(func() {
        // Reconstrói widget tree (captura novos valores)
        newTree := app.buildTree()
        
        // Diff e patch DOM (otimização futura)
        app.pipeline.Update(newTree)
    })
}
```

### Prevenção de Loops Infinitos

**PROBLEMA ATUAL:** Effects criados durante render causam loop infinito

```go
// ERRADO - Cria novo effect a cada render
func TextSignal[T](signal *Signal[T]) Widget {
    text := NewText(signal.Get())
    CreateEffect(func() { // NOVO EFFECT A CADA RENDER!
        text.SetText(signal.Get())
        scheduleRender() // LOOP INFINITO!
    })
    return text
}
```

**SOLUÇÃO 1:** Single Root Effect

```go
// Um único effect na raiz que reconstrói tudo
func (app *App) Run() {
    app.rootEffect = CreateEffect(func() {
        // Esta função rastreia TODOS signals usados
        newTree := app.root() // Reconstrói widgets
        app.updateDOM(newTree)
    })
}

// TextSignal apenas lê o signal (sem criar effect)
func TextSignal[T](signal *Signal[T]) Widget {
    return NewText(signal.Get()) // Apenas lê, tracking automático
}
```

**SOLUÇÃO 2:** Effect Registry

```go
// Registra effects uma vez, reutiliza nas re-renders
var effectRegistry = make(map[string]*Effect)

func TextSignal[T](signal *Signal[T], id string) Widget {
    text := NewText(signal.Get())
    
    // Cria effect apenas se não existe
    if _, exists := effectRegistry[id]; !exists {
        effectRegistry[id] = CreateEffect(func() {
            text.SetText(signal.Get())
            scheduleRender()
        })
    }
    
    return text
}
```

**SOLUÇÃO 3:** Tracking Context

```go
// Desabilita tracking durante render
func (app *App) render() {
    Untrack(func() {
        // Reconstrói tree sem criar dependências
        newTree := app.buildTree()
        app.updateDOM(newTree)
    })
}
```

## 📊 Comparação com Outras Abordagens

| Aspecto | Virtual DOM (React) | Proxy (Vue 3) | Signals (SolidJS/Maya) |
|---------|-------------------|---------------|----------------------|
| Update Complexity | O(n) | O(components) | O(effects) |
| Memory Overhead | 2x tree | Proxy wrappers | Signal nodes |
| Tracking | Manual deps | Auto via proxy | Auto via execution |
| Batching | Required | Optional | Optional |
| Debugging | React DevTools | Vue DevTools | Signal graph |
| TypeScript | Good | Good | Excellent |

## 🔮 Otimizações Futuras

### 1. Lazy Effects
```go
// Effects que só executam quando necessário
type LazyEffect struct {
    Effect
    priority int
    deferred bool
}
```

### 2. Weak Signals
```go
// Signals que podem ser coletados pelo GC
type WeakSignal[T any] struct {
    value *weak.Pointer[T]
    // ...
}
```

### 3. Compile-Time Optimization
```go
//go:generate maya-optimizer
// Análise estática para eliminar tracking desnecessário
```

## 📚 Bibliografia

1. **"A Hands-on Introduction to Fine-Grained Reactivity"** - Ryan Carniato (2021)
2. **"The Fundamental Principles Behind MobX"** - Michel Weststrate (2017)
3. **"SolidJS: Reactivity to Rendering"** - Ryan Carniato (2022)
4. **"Building a Reactive Library from Scratch"** - Milo Davis (2023)
5. **S.js** - Adam Haile (Original implementation)
6. **Solid.js Source Code** - github.com/solidjs/solid
7. **"Reactive Programming Patterns"** - André Staltz (2024)
8. **"From Observer Pattern to Reactive Systems"** - Martin Fowler (2024)

## 🎯 Próximos Passos

1. **Implementar Single Root Effect** - Evitar loops infinitos
2. **Adicionar Batch Updates** - Otimizar múltiplas mudanças
3. **Implementar Diff/Patch** - Atualizar DOM eficientemente
4. **Adicionar Suspense** - Lidar com async
5. **Criar DevTools** - Visualizar grafo de dependências

---

*Este documento baseia-se nas melhores práticas de 2025 para sistemas reativos, incorporando lições aprendidas de SolidJS, Vue 3, Svelte 5 e outras frameworks modernas.*