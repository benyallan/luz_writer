# Graph Report - Luz-writer  (2026-07-03)

## Corpus Check
- 32 files · ~20,004 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 244 nodes · 239 edges · 26 communities (23 shown, 3 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 1 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Luz Writer — Especificação de Construção|Luz Writer — Especificação de Construção]]
- [[_COMMUNITY_package.json|package.json]]
- [[_COMMUNITY_compilerOptions|compilerOptions]]
- [[_COMMUNITY_package.json|package.json]]
- [[_COMMUNITY_schema.go|schema.go]]
- [[_COMMUNITY_wails.json|wails.json]]
- [[_COMMUNITY_12. Plano de Construção em Etapas|12. Plano de Construção em Etapas]]
- [[_COMMUNITY_App.vue|App.vue]]
- [[_COMMUNITY_App|App]]
- [[_COMMUNITY_5. Schemas dos Arquivos (contratos de dados)|5. Schemas dos Arquivos (contratos de dados)]]
- [[_COMMUNITY_6. O Editor Tiptap — Definição dos Nós e Marks|6. O Editor Tiptap — Definição dos Nós e Marks]]
- [[_COMMUNITY_compilerOptions|compilerOptions]]
- [[_COMMUNITY_chapter.go|chapter.go]]
- [[_COMMUNITY_Target|Target]]
- [[_COMMUNITY_runtime.d.ts|runtime.d.ts]]
- [[_COMMUNITY_README|README]]
- [[_COMMUNITY_Vue 3 + TypeScript + Vite|Vue 3 + TypeScript + Vite]]
- [[_COMMUNITY_problem.go|problem.go]]
- [[_COMMUNITY_EventsOnMultiple|EventsOnMultiple]]
- [[_COMMUNITY_styles.go|styles.go]]
- [[_COMMUNITY_luz-writer|luz-writer]]

## God Nodes (most connected - your core abstractions)
1. `Luz Writer — Especificação de Construção` - 14 edges
2. `compilerOptions` - 13 edges
3. `12. Plano de Construção em Etapas` - 9 edges
4. `5. Schemas dos Arquivos (contratos de dados)` - 8 edges
5. `6. O Editor Tiptap — Definição dos Nós e Marks` - 7 edges
6. `8. Arquitetura de Plugins (Schema-Driven UI)` - 7 edges
7. `App` - 5 edges
8. `compilerOptions` - 5 edges
9. `Target` - 5 edges
10. `scripts` - 4 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewApp()`  [INFERRED]
  main.go → app.go
- `BuildContext` --references--> `Target`  [EXTRACTED]
  internal/model/plugin.go → internal/model/target.go
- `PluginManifest` --references--> `FormSchema`  [EXTRACTED]
  internal/model/plugin.go → internal/model/schema.go
- `BuildContext` --references--> `Project`  [EXTRACTED]
  internal/model/plugin.go → internal/model/project.go

## Import Cycles
- None detected.

## Communities (26 total, 3 thin omitted)

### Community 1 - "Luz Writer — Especificação de Construção"
Cohesion: 0.08
Nodes (25): 10. Pipeline de Compilação (Go), 11.1 Graphify — mapa navegável do projeto, 11.2 RTK — filtro de saídas de comandos, 11. Ferramentas de Apoio ao Agente, 1. Visão Geral, 2. Anatomia da Interface (UX / UI), 3. Estrutura do Repositório (código do Luz Writer), 4. Estrutura do Workspace do Usuário (+17 more)

### Community 2 - "package.json"
Cohesion: 0.09
Nodes (21): dependencies, pinia, reka-ui, @tiptap/pm, @tiptap/starter-kit, @tiptap/vue-3, vue, devDependencies (+13 more)

### Community 3 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, esModuleInterop, isolatedModules, jsx, lib, module, moduleResolution, resolveJsonModule (+7 more)

### Community 4 - "package.json"
Cohesion: 0.12
Nodes (15): author, bugs, url, description, homepage, keywords, license, main (+7 more)

### Community 5 - "schema.go"
Cohesion: 0.24
Nodes (9): BuildContext, FieldOption, FieldType, FormField, FormSchema, PluginManifest, Project, Variable (+1 more)

### Community 6 - "wails.json"
Cohesion: 0.18
Nodes (10): author, email, name, frontend:build, frontend:dev:serverUrl, frontend:dev:watcher, frontend:install, name (+2 more)

### Community 7 - "12. Plano de Construção em Etapas"
Cohesion: 0.22
Nodes (9): 12. Plano de Construção em Etapas, Etapa 0 — Fundação, Etapa 1 — Workspace e Capítulos (sem editor real), Etapa 2 — Editor Tiptap com Nós Semânticos, Etapa 3 — Conversor LaTeX + Primeira Compilação, Etapa 4 — Targets e Plugins (Schema-Driven UI), Etapa 5 — Rule Engine e Problems Panel, Etapa 6 — Polimento do MVP (+1 more)

### Community 9 - "App"
Cohesion: 0.32
Nodes (4): NewApp(), Context, App, main()

### Community 10 - "5. Schemas dos Arquivos (contratos de dados)"
Cohesion: 0.25
Nodes (8): 5.1 `.luz/project.json`, 5.2 `.luz/targets/<id>.json`, 5.3 `.luz/plugins.json`, 5.4 `capitulos/<id>.json`, 5.5 Papéis de Documento (`role`) — Páginas Especiais e Front/Back Matter, 5.6 `.luz/overrides/<id-do-documento>.json` — Configurações por Página, 5.7 `.luz/styles.json` — Estilos Personalizados, 5. Schemas dos Arquivos (contratos de dados)

### Community 11 - "6. O Editor Tiptap — Definição dos Nós e Marks"
Cohesion: 0.29
Nodes (7): 6.0 Princípio: WYSIWYG Semântico — o editor mostra INTENÇÃO, não preview, 6. O Editor Tiptap — Definição dos Nós e Marks, Escape de LaTeX (obrigatório), Exemplo canônico de conversão (usar como teste), Marks, Nós de bloco, Nós inline

### Community 12 - "compilerOptions"
Cohesion: 0.29
Nodes (6): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, include

### Community 13 - "chapter.go"
Cohesion: 0.53
Nodes (5): RawMessage, Chapter, ChapterMeta, DocumentOverrides, DocumentRole

### Community 14 - "Target"
Cohesion: 0.47
Nodes (5): RawMessage, DocumentClass, PluginsConfig, Target, TargetKind

### Community 15 - "runtime.d.ts"
Cohesion: 0.40
Nodes (4): EnvironmentInfo, Position, Screen, Size

### Community 16 - "README"
Cohesion: 0.40
Nodes (4): About, Building, Live Development, README

### Community 17 - "Vue 3 + TypeScript + Vite"
Cohesion: 0.50
Nodes (3): Recommended IDE Setup, Type Support For `.vue` Imports in TS, Vue 3 + TypeScript + Vite

### Community 18 - "problem.go"
Cohesion: 0.67
Nodes (3): BuildProgress, BuildResult, Problem

### Community 19 - "EventsOnMultiple"
Cohesion: 0.67
Nodes (3): EventsOn(), EventsOnce(), EventsOnMultiple()

## Knowledge Gaps
- **114 isolated node(s):** `name`, `private`, `version`, `type`, `dev` (+109 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Luz Writer — Especificação de Construção` connect `Luz Writer — Especificação de Construção` to `5. Schemas dos Arquivos (contratos de dados)`, `6. O Editor Tiptap — Definição dos Nós e Marks`, `12. Plano de Construção em Etapas`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **Why does `12. Plano de Construção em Etapas` connect `12. Plano de Construção em Etapas` to `Luz Writer — Especificação de Construção`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Why does `5. Schemas dos Arquivos (contratos de dados)` connect `5. Schemas dos Arquivos (contratos de dados)` to `Luz Writer — Especificação de Construção`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **What connects `name`, `private`, `version` to the rest of the system?**
  _114 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `runtime.js` be split into smaller, more focused modules?**
  _Cohesion score 0.0392156862745098 - nodes in this community are weakly interconnected._
- **Should `Luz Writer — Especificação de Construção` be split into smaller, more focused modules?**
  _Cohesion score 0.07692307692307693 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.09090909090909091 - nodes in this community are weakly interconnected._