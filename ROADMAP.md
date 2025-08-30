# Maya Framework - Roadmap de Desenvolvimento 2025-2026

## Visão Geral

Maya é uma framework UI moderna em Go 1.24+ compilada para WebAssembly, aproveitando as tecnologias mais recentes para criar uma experiência de desenvolvimento superior.

## Atualização Agosto 2025

Com Go 1.24 lançado e WebGPU maduro, ajustamos o roadmap para aproveitar as novas features.

## Estrutura de Épicos

```
📦 Epic 1: Core Foundation & Modern Go Features (3-4 semanas)
📦 Epic 2: Fine-Grained Reactive System (2-3 semanas)
📦 Epic 3: Type-Safe Widget System (3-4 semanas)
📦 Epic 4: GPU-Accelerated Layout Engine (4-5 semanas)
📦 Epic 5: Hybrid Rendering Pipeline (3-4 semanas)
📦 Epic 6: Modern Event & Input System (2-3 semanas)
📦 Epic 7: Physics-Based Animation (2-3 semanas)
📦 Epic 8: Adaptive Design System (2-3 semanas)
📦 Epic 9: Developer Experience & Tools (3-4 semanas)
📦 Epic 10: WASM Optimization & Distribution (2-3 semanas)
```

**Tempo Total Estimado:** 20-30 semanas (5-7 meses) - Reduzido devido às novas features do Go 1.24

---

## Epic 1: Core Foundation & Modern Go Features 🚀

### Objetivo
Estabelecer fundação sólida usando Go 1.24+ features, incluindo generic type aliases (agora oficial!), runtime.AddCleanup, weak pointers, e go:wasmexport.

### Tasks

#### Task 1.1: Project Setup com Go 1.24+
- [ ] **1.1.1** Inicializar módulo Go 1.24 com estrutura moderna
  ```
  /cmd         - Entry points
  /pkg         - Public packages
  /internal    - Internal packages
  /wasm        - WASM specific code
  /shaders     - WGSL shaders
  /examples    - Example apps
  ```
- [ ] **1.1.2** Configurar go.mod com Go 1.24 e tool directives
- [ ] **1.1.3** Setup build system com suporte a TinyGo e WASI 0.2
- [ ] **1.1.4** Configurar GitHub Actions com Go 1.24
- [ ] **1.1.5** Criar Makefile com targets para WASM/TinyGo/go:wasmexport

#### Task 1.2: Core Data Structures com Generics
- [ ] **1.2.1** Implementar Node com unique.Handle para IDs
- [ ] **1.2.2** Criar Tree com iteradores nativos (iter.Seq)
- [ ] **1.2.3** Implementar PriorityQueue genérica
- [ ] **1.2.4** Criar Pool[T] genérico com reset functions
- [ ] **1.2.5** Implementar WeakCache com referências fracas

#### Task 1.3: Modern Tree Traversal
- [ ] **1.3.1** Implementar iteradores DFS/BFS usando iter.Seq
- [ ] **1.3.2** Criar ParallelSubtrees iterator
- [ ] **1.3.3** Implementar iterador com yield para controle fino
- [ ] **1.3.4** Adicionar benchmarks comparando com old approach
- [ ] **1.3.5** Otimizar com SIMD onde possível

#### Task 1.4: Workflow System com Concorrência
- [ ] **1.4.1** Implementar Pipeline[T, R] com generics
- [ ] **1.4.2** Criar WorkerPool genérico
- [ ] **1.4.3** Implementar Supervisor com health monitoring
- [ ] **1.4.4** Criar sistema de backpressure
- [ ] **1.4.5** Adicionar graceful shutdown

#### Task 1.5: Memory Optimization com Go 1.24
- [ ] **1.5.1** Implementar string interning com unique.Handle
- [ ] **1.5.2** Usar weak.Pointer para caches eficientes
- [ ] **1.5.3** Substituir SetFinalizer por runtime.AddCleanup
- [ ] **1.5.4** Implementar pools com Swiss Tables (maps mais rápidos)
- [ ] **1.5.5** Criar benchmarks usando testing.B.Loop

---

## Epic 2: Fine-Grained Reactive System ⚛️

### Objetivo
Construir sistema de reatividade inspirado em Solid.js com signals e effects, superando Virtual DOM.

### Tasks

#### Task 2.1: Signal System
- [ ] **2.1.1** Implementar Signal[T] com canonicalização
- [ ] **2.1.2** Criar sistema de versioning atômico
- [ ] **2.1.3** Implementar auto-tracking de dependências
- [ ] **2.1.4** Adicionar batching de updates
- [ ] **2.1.5** Criar testes de stress com 10k signals

#### Task 2.2: Memo & Computed
- [ ] **2.2.1** Implementar Memo[T] com lazy evaluation
- [ ] **2.2.2** Criar Computed values com cache
- [ ] **2.2.3** Implementar invalidação seletiva
- [ ] **2.2.4** Adicionar dependency pruning
- [ ] **2.2.5** Otimizar recomputação

#### Task 2.3: Effect System
- [ ] **2.3.1** Implementar Effect com cleanup automático
- [ ] **2.3.2** Criar EffectScheduler com prioridades
- [ ] **2.3.3** Implementar effect batching
- [ ] **2.3.4** Adicionar error boundaries
- [ ] **2.3.5** Criar effect devtools

#### Task 2.4: Eliminar Virtual DOM
- [ ] **2.4.1** Implementar direct DOM manipulation
- [ ] **2.4.2** Criar fine-grained updates
- [ ] **2.4.3** Implementar surgical DOM patches
- [ ] **2.4.4** Adicionar benchmarks vs VDOM
- [ ] **2.4.5** Otimizar para <1ms updates

#### Task 2.5: State Management
- [ ] **2.5.1** Implementar Store[S] genérico
- [ ] **2.5.2** Criar sistema de reducers type-safe
- [ ] **2.5.3** Implementar middleware pipeline
- [ ] **2.5.4** Adicionar time-travel debugging
- [ ] **2.5.5** Criar persistence layer

---

## Epic 3: Type-Safe Widget System 🎨

### Objetivo
Desenvolver sistema de widgets declarativo com type safety completo usando generics.

### Tasks

#### Task 3.1: Widget Core
- [ ] **3.1.1** Definir Widget interface com generics
- [ ] **3.1.2** Implementar BaseWidget com signals
- [ ] **3.1.3** Criar Props system type-safe
- [ ] **3.1.4** Implementar widget lifecycle
- [ ] **3.1.5** Adicionar widget pooling

#### Task 3.2: Builder Pattern
- [ ] **3.2.1** Implementar WidgetBuilder[W, P] genérico
- [ ] **3.2.2** Criar fluent API com type inference
- [ ] **3.2.3** Adicionar prop validation em compile-time
- [ ] **3.2.4** Implementar children management
- [ ] **3.2.5** Criar macro system para reduzir boilerplate

#### Task 3.3: Component System
- [ ] **3.3.1** Implementar functional components
- [ ] **3.3.2** Criar hooks system (useState, useEffect, etc)
- [ ] **3.3.3** Implementar context API type-safe
- [ ] **3.3.4** Adicionar component memoization
- [ ] **3.3.5** Criar component devtools

#### Task 3.4: Core Widgets
- [ ] **3.4.1** Implementar Container, Text, Image
- [ ] **3.4.2** Criar Button, Input, Select
- [ ] **3.4.3** Implementar List, Grid views
- [ ] **3.4.4** Adicionar Modal, Dialog, Tooltip
- [ ] **3.4.5** Criar Navigation components

#### Task 3.5: Layout Widgets
- [ ] **3.5.1** Implementar Row/Column com flex
- [ ] **3.5.2** Criar Stack para overlays
- [ ] **3.5.3** Implementar Grid layout
- [ ] **3.5.4** Adicionar Wrap para flow layout
- [ ] **3.5.5** Criar responsive containers

---

## Epic 4: GPU-Accelerated Layout Engine 📐

### Objetivo
Implementar layout engine com WebGPU compute shaders para performance máxima.

### Tasks

#### Task 4.1: WebGPU Integration
- [ ] **4.1.1** Criar abstração WebGPU para Go/WASM
- [ ] **4.1.2** Implementar shader compiler/validator
- [ ] **4.1.3** Criar buffer management system
- [ ] **4.1.4** Implementar compute pipeline
- [ ] **4.1.5** Adicionar fallback para Canvas2D

#### Task 4.2: GPU Layout Algorithms
- [ ] **4.2.1** Portar Flexbox para WGSL
- [ ] **4.2.2** Implementar Grid layout em GPU
- [ ] **4.2.3** Criar constraint solver paralelo
- [ ] **4.2.4** Implementar text layout na GPU
- [ ] **4.2.5** Otimizar para 10k+ widgets

#### Task 4.3: Multi-pass Pipeline
- [ ] **4.3.1** Implementar 6-phase layout pipeline
- [ ] **4.3.2** Criar dirty tracking inteligente
- [ ] **4.3.3** Implementar incremental layout
- [ ] **4.3.4** Adicionar layout caching
- [ ] **4.3.5** Criar profiler para layout

#### Task 4.4: CPU Fallback
- [ ] **4.4.1** Implementar Flexbox em Go puro
- [ ] **4.4.2** Criar Grid solver otimizado
- [ ] **4.4.3** Implementar auto-switch GPU/CPU
- [ ] **4.4.4** Adicionar benchmarks comparativos
- [ ] **4.4.5** Otimizar para mobile devices

#### Task 4.5: Advanced Features
- [ ] **4.5.1** Implementar custom layout protocol
- [ ] **4.5.2** Criar layout animations
- [ ] **4.5.3** Adicionar viewport culling
- [ ] **4.5.4** Implementar virtualization
- [ ] **4.5.5** Criar layout debugger visual

---

## Epic 5: Hybrid Rendering Pipeline 🎯

### Objetivo
Construir pipeline de renderização híbrido WebGPU/Canvas2D com detecção automática.

### Tasks

#### Task 5.1: Renderer Architecture
- [ ] **5.1.1** Implementar Renderer interface unificada
- [ ] **5.1.2** Criar WebGPU renderer
- [ ] **5.1.3** Implementar Canvas2D renderer
- [ ] **5.1.4** Adicionar auto-detection de capacidades
- [ ] **5.1.5** Criar switching em runtime

#### Task 5.2: WebGPU Rendering
- [ ] **5.2.1** Implementar render pipeline
- [ ] **5.2.2** Criar vertex/fragment shaders
- [ ] **5.2.3** Implementar instanced rendering
- [ ] **5.2.4** Adicionar texture atlas
- [ ] **5.2.5** Otimizar draw calls

#### Task 5.3: Command Buffer
- [ ] **5.3.1** Implementar command pattern
- [ ] **5.3.2** Criar command merging
- [ ] **5.3.3** Implementar state sorting
- [ ] **5.3.4** Adicionar command replay
- [ ] **5.3.5** Criar command profiler

#### Task 5.4: Optimization
- [ ] **5.4.1** Implementar occlusion culling
- [ ] **5.4.2** Criar dirty rectangle tracking
- [ ] **5.4.3** Implementar layer compositing
- [ ] **5.4.4** Adicionar render caching
- [ ] **5.4.5** Otimizar para 60+ FPS

#### Task 5.5: Effects & Filters
- [ ] **5.5.1** Implementar blur, shadows
- [ ] **5.5.2** Criar color filters
- [ ] **5.5.3** Adicionar blend modes
- [ ] **5.5.4** Implementar masks/clips
- [ ] **5.5.5** Criar custom shader support

---

## Epic 6: Modern Event & Input System 🖱️

### Objetivo
Sistema de eventos moderno com gesture recognition e processamento otimizado.

### Tasks

#### Task 6.1: Event Core
- [ ] **6.1.1** Implementar typed event system
- [ ] **6.1.2** Criar event bus com generics
- [ ] **6.1.3** Implementar event bubbling/capturing
- [ ] **6.1.4** Adicionar event delegation
- [ ] **6.1.5** Criar event replay system

#### Task 6.2: Input Processing
- [ ] **6.2.1** Implementar mouse/keyboard handlers
- [ ] **6.2.2** Criar touch event processing
- [ ] **6.2.3** Implementar pointer events
- [ ] **6.2.4** Adicionar gamepad support
- [ ] **6.2.5** Criar input throttling/debouncing

#### Task 6.3: Gesture Recognition
- [ ] **6.3.1** Implementar tap, double-tap, long-press
- [ ] **6.3.2** Criar swipe, pan, fling gestures
- [ ] **6.3.3** Implementar pinch, rotate, zoom
- [ ] **6.3.4** Adicionar custom gesture API
- [ ] **6.3.5** Criar gesture conflict resolution

#### Task 6.4: Hit Testing
- [ ] **6.4.1** Implementar R-Tree para spatial indexing
- [ ] **6.4.2** Criar hit test caching
- [ ] **6.4.3** Implementar custom hit areas
- [ ] **6.4.4** Adicionar hit test debugging
- [ ] **6.4.5** Otimizar para complex shapes

#### Task 6.5: Accessibility
- [ ] **6.5.1** Implementar keyboard navigation
- [ ] **6.5.2** Criar screen reader support
- [ ] **6.5.3** Adicionar focus management
- [ ] **6.5.4** Implementar ARIA attributes
- [ ] **6.5.5** Criar a11y testing tools

---

## Epic 7: Physics-Based Animation 🎬

### Objetivo
Sistema de animação com spring physics e timeline control para animações naturais.

### Tasks

#### Task 7.1: Animation Core
- [ ] **7.1.1** Implementar Animation base
- [ ] **7.1.2** Criar AnimationController
- [ ] **7.1.3** Implementar easing functions
- [ ] **7.1.4** Adicionar animation state machine
- [ ] **7.1.5** Criar animation pooling

#### Task 7.2: Spring Physics
- [ ] **7.2.1** Implementar spring solver
- [ ] **7.2.2** Criar damping system
- [ ] **7.2.3** Implementar velocity tracking
- [ ] **7.2.4** Adicionar spring presets
- [ ] **7.2.5** Criar spring debugger

#### Task 7.3: Timeline System
- [ ] **7.3.1** Implementar Timeline controller
- [ ] **7.3.2** Criar keyframe system
- [ ] **7.3.3** Implementar animation blending
- [ ] **7.3.4** Adicionar timeline events
- [ ] **7.3.5** Criar timeline editor

#### Task 7.4: Advanced Animations
- [ ] **7.4.1** Implementar morph animations
- [ ] **7.4.2** Criar particle system
- [ ] **7.4.3** Adicionar skeletal animation
- [ ] **7.4.4** Implementar animation sequences
- [ ] **7.4.5** Criar animation templates

#### Task 7.5: Performance
- [ ] **7.5.1** Implementar RAF scheduling
- [ ] **7.5.2** Criar animation batching
- [ ] **7.5.3** Adicionar animation culling
- [ ] **7.5.4** Implementar GPU animations
- [ ] **7.5.5** Otimizar para 120+ FPS

---

## Epic 8: Adaptive Design System 🎨

### Objetivo
Sistema de design moderno, responsivo e acessível com temas dinâmicos.

### Tasks

#### Task 8.1: Theme System
- [ ] **8.1.1** Implementar Theme com design tokens
- [ ] **8.1.2** Criar color system com palettes
- [ ] **8.1.3** Implementar typography scales
- [ ] **8.1.4** Adicionar spacing system
- [ ] **8.1.5** Criar theme inheritance

#### Task 8.2: Component Library
- [ ] **8.2.1** Implementar Material Design components
- [ ] **8.2.2** Criar iOS-style components
- [ ] **8.2.3** Adicionar custom component API
- [ ] **8.2.4** Implementar component variants
- [ ] **8.2.5** Criar component playground

#### Task 8.3: Responsive Design
- [ ] **8.3.1** Implementar breakpoint system
- [ ] **8.3.2** Criar responsive utilities
- [ ] **8.3.3** Adicionar container queries
- [ ] **8.3.4** Implementar adaptive layouts
- [ ] **8.3.5** Criar responsive debugger

#### Task 8.4: Icons & Assets
- [ ] **8.4.1** Implementar icon system
- [ ] **8.4.2** Criar SVG support
- [ ] **8.4.3** Adicionar icon fonts
- [ ] **8.4.4** Implementar lazy loading
- [ ] **8.4.5** Criar asset optimization

#### Task 8.5: Dark Mode & Themes
- [ ] **8.5.1** Implementar dark/light modes
- [ ] **8.5.2** Criar theme switching
- [ ] **8.5.3** Adicionar system preference detection
- [ ] **8.5.4** Implementar custom themes
- [ ] **8.5.5** Criar theme builder tool

---

## Epic 9: Developer Experience & Tools 🛠️

### Objetivo
Ferramentas de desenvolvimento superiores com hot reload, debugging e testing.

### Tasks

#### Task 9.1: DevTools
- [ ] **9.1.1** Implementar Widget Inspector
- [ ] **9.1.2** Criar Performance Profiler
- [ ] **9.1.3** Adicionar Signal Debugger
- [ ] **9.1.4** Implementar Layout Visualizer
- [ ] **9.1.5** Criar Memory Analyzer

#### Task 9.2: Hot Module Replacement
- [ ] **9.2.1** Implementar file watcher
- [ ] **9.2.2** Criar incremental compiler
- [ ] **9.2.3** Adicionar state preservation
- [ ] **9.2.4** Implementar error recovery
- [ ] **9.2.5** Otimizar para <500ms reload

#### Task 9.3: Testing Framework
- [ ] **9.3.1** Implementar WidgetTester
- [ ] **9.3.2** Criar snapshot testing
- [ ] **9.3.3** Adicionar visual regression
- [ ] **9.3.4** Implementar e2e testing
- [ ] **9.3.5** Criar test coverage tools

#### Task 9.4: Documentation
- [ ] **9.4.1** Gerar API docs automáticos
- [ ] **9.4.2** Criar interactive tutorials
- [ ] **9.4.3** Adicionar code examples
- [ ] **9.4.4** Implementar playground
- [ ] **9.4.5** Criar video tutorials

#### Task 9.5: CLI & Tooling
- [ ] **9.5.1** Criar maya CLI tool
- [ ] **9.5.2** Implementar project generator
- [ ] **9.5.3** Adicionar component scaffolding
- [ ] **9.5.4** Criar migration tools
- [ ] **9.5.5** Implementar VS Code extension

---

## Epic 10: WASM Optimization & Distribution 📦

### Objetivo
Otimizar para produção com bundle mínimo e performance máxima.

### Tasks

#### Task 10.1: WASM Optimization
- [ ] **10.1.1** Implementar dual build (Go/TinyGo)
- [ ] **10.1.2** Criar tree shaking
- [ ] **10.1.3** Adicionar code splitting
- [ ] **10.1.4** Implementar lazy loading
- [ ] **10.1.5** Otimizar para <100KB gzipped

#### Task 10.2: Build Pipeline
- [ ] **10.2.1** Criar build orchestrator
- [ ] **10.2.2** Implementar parallel compilation
- [ ] **10.2.3** Adicionar asset optimization
- [ ] **10.2.4** Criar source maps
- [ ] **10.2.5** Implementar CI/CD pipeline

#### Task 10.3: Performance
- [ ] **10.3.1** Implementar streaming compilation
- [ ] **10.3.2** Criar shared memory support
- [ ] **10.3.3** Adicionar Web Workers
- [ ] **10.3.4** Implementar caching strategies
- [ ] **10.3.5** Criar performance budgets

#### Task 10.4: Distribution
- [ ] **10.4.1** Criar NPM package
- [ ] **10.4.2** Implementar CDN distribution
- [ ] **10.4.3** Adicionar module federation
- [ ] **10.4.4** Criar standalone builds
- [ ] **10.4.5** Implementar auto-updates

#### Task 10.5: Platform Support
- [ ] **10.5.1** Validar em Chrome, Firefox, Safari
- [ ] **10.5.2** Testar em mobile browsers
- [ ] **10.5.3** Adicionar polyfills necessários
- [ ] **10.5.4** Criar compatibility matrix
- [ ] **10.5.5** Implementar feature detection

---

## Cronograma por Quarter

### Q4 2024 (Out-Dez)
- **Epic 1:** Core Foundation & Modern Go Features
- **Epic 2:** Fine-Grained Reactive System
- **Epic 3:** Type-Safe Widget System (início)

### Q1 2025 (Jan-Mar)
- **Epic 3:** Type-Safe Widget System (conclusão)
- **Epic 4:** GPU-Accelerated Layout Engine
- **Epic 5:** Hybrid Rendering Pipeline

### Q2 2025 (Abr-Jun)
- **Epic 6:** Modern Event & Input System
- **Epic 7:** Physics-Based Animation
- **Epic 8:** Adaptive Design System
- **Epic 9:** Developer Experience & Tools (início)

### Q3 2025 (Jul-Set)
- **Epic 9:** Developer Experience & Tools (conclusão)
- **Epic 10:** WASM Optimization & Distribution
- **Beta Release**
- **Community Feedback**

### Q4 2025 (Out-Dez)
- **Performance Optimization**
- **Bug Fixes**
- **Documentation Completion**
- **1.0 Release**

---

## Métricas de Sucesso

### Performance Targets
| Métrica | Target | Atual | Status |
|---------|--------|-------|--------|
| First Paint | < 50ms | - | 🔴 |
| Re-render (1000 nodes) | < 16ms | - | 🔴 |
| Memory (10k widgets) | < 20MB | - | 🔴 |
| Bundle Size (gzipped) | < 100KB | - | 🔴 |
| Layout Computation | < 1ms | - | 🔴 |
| 60 FPS Consistency | > 95% | - | 🔴 |

### Developer Experience
| Métrica | Target | Atual | Status |
|---------|--------|-------|--------|
| Hot Reload Time | < 500ms | - | 🔴 |
| Test Coverage | > 90% | - | 🔴 |
| API Documentation | 100% | - | 🔴 |
| Time to First App | < 5min | - | 🔴 |
| CLI Commands | > 20 | - | 🔴 |

### Adoption Metrics
| Métrica | Target | Atual | Status |
|---------|--------|-------|--------|
| GitHub Stars | > 5000 | 0 | 🔴 |
| NPM Downloads/month | > 10k | 0 | 🔴 |
| Active Contributors | > 50 | 1 | 🔴 |
| Discord Members | > 1000 | 0 | 🔴 |
| Production Apps | > 100 | 0 | 🔴 |

---

## Riscos e Mitigações

### Riscos Técnicos

1. **WebGPU Browser Support**
   - Risco: Safari ainda sem algumas features (Memory64, Multiple Memories)
   - Mitigação: Fallback robusto para Canvas2D
   - Status 2025: Chrome/Firefox totalmente suportados, Safari 90%

2. **WASM Bundle Size**
   - Risco: Go produz bundles grandes (~1.8MB mínimo)
   - Mitigação: TinyGo (~8KB hello world) + go:wasmexport
   - Solução 2025: go:wasmexport reduz overhead significativamente

3. **Performance em Mobile**
   - Risco: Devices low-end podem ter problemas
   - Mitigação: Adaptive rendering
   - Contingência: Lite version para mobile

4. **Complexidade de Signals**
   - Risco: Debugging difícil
   - Mitigação: DevTools comprehensive
   - Contingência: Virtual DOM híbrido opcional

### Riscos de Projeto

1. **Scope Creep**
   - Mitigação: Reviews quinzenais rigorosas
   - MVPs incrementais por epic

2. **Competição (React, Flutter)**
   - Mitigação: Focar em diferenciais (Go, WebGPU)
   - Parcerias com comunidade Go

3. **Adoção Lenta**
   - Mitigação: Documentação excepcional
   - Templates e starters ricos

---

## Milestones Principais

### 🎯 M1: Alpha Release (Jan 2025)
- Core reactive system funcionando
- 10 widgets básicos
- Hot reload operacional
- Documentação inicial

### 🎯 M2: Beta Release (Jul 2025)
- WebGPU rendering ativo
- 50+ widgets
- DevTools completo
- 90% test coverage

### 🎯 M3: RC Release (Set 2025)
- Performance targets atingidos
- Documentação completa
- CLI tools finalizadas
- Community feedback incorporated

### 🎯 M4: 1.0 Release (Mar 2026)
- Production ready
- Stable API com go:wasmexport
- WebGPU compute shaders para layout
- WASI 0.3 support (async)
- Ecosystem estabelecido
- Case studies publicados

---

## Próximos Passos Imediatos

### Semana 1-2 (Set 2025)
1. ⬜ Setup repositório GitHub
2. ⬜ Configurar Go 1.24 environment
3. ⬜ Implementar generic type aliases
4. ⬜ Criar PoC com go:wasmexport
5. ⬜ Setup CI/CD com Go 1.24

### Semana 3-4 (Set 2025)
1. ⬜ Implementar weak.Pointer para caches
2. ⬜ Usar runtime.AddCleanup em widgets
3. ⬜ PoC WebGPU compute shaders
4. ⬜ Benchmarks com testing.B.Loop
5. ⬜ Publicar roadmap atualizado

### Mês 2
1. ⬜ Alpha interno funcional
2. ⬜ 5 exemplos rodando
3. ⬜ Performance baseline estabelecido
4. ⬜ Comunidade Discord criada
5. ⬜ Primeiro contributor externo

---
