# Graph Report - .  (2026-07-02)

## Corpus Check
- Corpus is ~25,835 words - fits in a single context window. You may not need a graph.

## Summary
- 338 nodes · 466 edges · 20 communities (15 shown, 5 thin omitted)
- Extraction: 94% EXTRACTED · 6% INFERRED · 0% AMBIGUOUS · INFERRED: 29 edges (avg confidence: 0.82)
- Token cost: 122,108 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_New Project Creation Flow|New Project Creation Flow]]
- [[_COMMUNITY_Graphify Skill Commands|Graphify Skill Commands]]
- [[_COMMUNITY_Luz Writer Project Overview|Luz Writer Project Overview]]
- [[_COMMUNITY_Sidebar and File Tree UI|Sidebar and File Tree UI]]
- [[_COMMUNITY_LaTeX Parser Tests|LaTeX Parser Tests]]
- [[_COMMUNITY_Frontend Package Dependencies|Frontend Package Dependencies]]
- [[_COMMUNITY_File and Rich Text Editor|File and Rich Text Editor]]
- [[_COMMUNITY_Wails Go Bindings and Demo|Wails Go Bindings and Demo]]
- [[_COMMUNITY_Wails Runtime Package Metadata|Wails Runtime Package Metadata]]
- [[_COMMUNITY_Wails Project Config|Wails Project Config]]
- [[_COMMUNITY_LaTeX Parser Core|LaTeX Parser Core]]
- [[_COMMUNITY_Wails Go Data Models|Wails Go Data Models]]
- [[_COMMUNITY_Wails Runtime Environment Types|Wails Runtime Environment Types]]
- [[_COMMUNITY_Luz Writer File Formats|Luz Writer File Formats]]
- [[_COMMUNITY_Wails Runtime Event Listeners|Wails Runtime Event Listeners]]
- [[_COMMUNITY_Nunito Font License|Nunito Font License]]
- [[_COMMUNITY_Wails Logo Asset|Wails Logo Asset]]
- [[_COMMUNITY_Go Module Root|Go Module Root]]

## God Nodes (most connected - your core abstractions)
1. `Graphify (/graphify skill)` - 26 edges
2. `Luz Writer` - 22 edges
3. `Parse()` - 19 edges
4. `App` - 17 edges
5. `assertNoError()` - 16 edges
6. `assertEqual()` - 16 edges
7. `doc()` - 15 edges
8. `renderBlock()` - 7 edges
9. `renderInline()` - 7 edges
10. `TestParseParagraph()` - 7 edges

## Surprising Connections (you probably didn't know these)
- `Graphify Skill Trigger` --semantically_similar_to--> `Graphify Project Rules`  [INFERRED] [semantically similar]
  .claude/CLAUDE.md → CLAUDE.md
- `Vue 3 + Vite Template` --semantically_similar_to--> `Official Wails Vue Template`  [INFERRED] [semantically similar]
  frontend/README.md → README.md
- `Official Wails Vue Template` --shares_data_with--> `Wails v2`  [INFERRED]
  README.md → .claude/commands/project-overview.md
- `main()` --calls--> `NewApp()`  [INFERRED]
  main.go → app.go
- `Vue 3` --shares_data_with--> `Vue 3 + Vite Template`  [INFERRED]
  .claude/commands/project-overview.md → frontend/README.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Luz Writer Frontend Stack** — claude_commands_project_overview_vue3, claude_commands_project_overview_vite, claude_commands_project_overview_tiptap, claude_commands_project_overview_rekaui [EXTRACTED 1.00]
- **Luz Writer Project Scaffold Generation** — claude_commands_project_overview_createproject, claude_commands_project_overview_luztxt, claude_commands_project_overview_luzprof, claude_commands_project_overview_a4_luzprof [EXTRACTED 1.00]
- **Graphify Subcommand Reference Set** — claude_skills_graphify_skill_graphify, claude_skills_graphify_references_query_query_command, claude_skills_graphify_references_github_and_merge_clone, claude_skills_graphify_references_hooks_post_commit_hook, claude_skills_graphify_references_update_incremental_update, claude_skills_graphify_references_add_watch_graphify_add [INFERRED 0.85]

## Communities (20 total, 5 thin omitted)

### Community 1 - "New Project Creation Flow"
Cohesion: 0.07
Nodes (23): NewApp(), recentsFilePath(), Context, chooseParent(), errorMsg, isSubmitting, parentPath, projectName (+15 more)

### Community 2 - "Graphify Skill Commands"
Cohesion: 0.07
Nodes (34): Graphify Skill Trigger, GRAPH_REPORT.md, Graphify Project Rules, /graphify add <url>, --watch (folder watcher), FalkorDB Export, GraphML Export, MCP Server (--mcp) (+26 more)

### Community 3 - "Luz Writer Project Overview"
Cohesion: 0.07
Nodes (34): Escritor Nunca Toca LaTeX (Princípio Central), distraction-free Plugin (planned), Go, internal/compiler, internal/config, internal/document, internal/latex, internal/plugin (+26 more)

### Community 4 - "Sidebar and File Tree UI"
Cohesion: 0.08
Nodes (20): validateItemName(), cancelRootCreate(), confirmRootCreate(), rootCreatingInput, rootCreatingType, rootNewItemName, workspace, workspace (+12 more)

### Community 5 - "LaTeX Parser Tests"
Cohesion: 0.33
Nodes (25): Parse(), assertEqual(), assertNoError(), doc(), join(), jsonStr(), markedText(), para() (+17 more)

### Community 6 - "Frontend Package Dependencies"
Cohesion: 0.08
Nodes (23): dependencies, @lucide/vue, pinia, reka-ui, @tiptap/extension-character-count, @tiptap/extension-placeholder, @tiptap/extension-typography, @tiptap/pm (+15 more)

### Community 7 - "File and Rich Text Editor"
Cohesion: 0.12
Nodes (9): editor, isRichText, store, workspace, EMPTY_DOC, store, tiptap, app (+1 more)

### Community 9 - "Wails Runtime Package Metadata"
Cohesion: 0.12
Nodes (15): author, bugs, url, description, homepage, keywords, license, main (+7 more)

### Community 10 - "Wails Project Config"
Cohesion: 0.18
Nodes (10): author, email, name, frontend:build, frontend:dev:serverUrl, frontend:dev:watcher, frontend:install, name (+2 more)

### Community 11 - "LaTeX Parser Core"
Cohesion: 0.49
Nodes (9): Builder, applyMarks(), escapeLatex(), headingCommand(), renderBlock(), renderInline(), renderListItem(), Mark (+1 more)

### Community 13 - "Wails Runtime Environment Types"
Cohesion: 0.40
Nodes (4): EnvironmentInfo, Position, Screen, Size

### Community 14 - "Luz Writer File Formats"
Cohesion: 0.67
Nodes (4): templates/a4.luzprof, CreateProject (app.go), .luzprof File Format, .luztxt File Format

### Community 15 - "Wails Runtime Event Listeners"
Cohesion: 0.67
Nodes (3): EventsOn(), EventsOnce(), EventsOnMultiple()

## Knowledge Gaps
- **126 isolated node(s):** `name`, `private`, `version`, `type`, `dev` (+121 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Parse()` connect `LaTeX Parser Tests` to `New Project Creation Flow`, `LaTeX Parser Core`?**
  _High betweenness centrality (0.060) - this node is a cross-community bridge._
- **Why does `App` connect `New Project Creation Flow` to `Sidebar and File Tree UI`?**
  _High betweenness centrality (0.058) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `Parse()` (e.g. with `TestParseBlockquote()` and `TestParseBold()`) actually correct?**
  _`Parse()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `name`, `private`, `version` to the rest of the system?**
  _130 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Wails Runtime JS API` be split into smaller, more focused modules?**
  _Cohesion score 0.0392156862745098 - nodes in this community are weakly interconnected._
- **Should `New Project Creation Flow` be split into smaller, more focused modules?**
  _Cohesion score 0.06707317073170732 - nodes in this community are weakly interconnected._
- **Should `Graphify Skill Commands` be split into smaller, more focused modules?**
  _Cohesion score 0.0659536541889483 - nodes in this community are weakly interconnected._