# Luz Writer — Especificação de Construção

> **Instruções para o agente (Claude Code):** Este documento é a fonte da verdade do projeto. Leia-o integralmente antes de escrever qualquer código. A construção é dividida em **etapas** (seção final). **Nunca avance para a etapa seguinte sem aprovação explícita do usuário.** Use o Graphify para se orientar no código e o RTK para executar comandos (seções "Ferramentas de Apoio ao Agente").

---

## 1. Visão Geral

Luz Writer é uma IDE desktop expansível com foco absoluto em escritores (livros, artigos, e-books e roteiros). Construída com **Wails v2** (Go + Vue 3), ela abandona o conceito de "editor de texto comum" para adotar uma arquitetura inspirada no VS Code: um núcleo limpo, edição livre de distrações e um ecossistema modular de configurações.

### Princípio Central

O escritor **nunca toca em código LaTeX**. Ele escreve em uma interface "Modo Fluxo" limpa (WYSIWYG semântico). Toda a configuração estrutural é feita através de painéis visuais interativos. O backend em Go age como um compilador e servidor de validação em tempo real (estilo LSP): valida incompatibilidades editoriais, traduz o documento e orquestra a geração tipográfica via **Tectonic**.

### Decisões Fixas (não reavaliar)

| Decisão | Valor | Justificativa |
|---|---|---|
| Motor LaTeX | **Tectonic** (binário único, baixa dependências automaticamente) | Zero configuração para o usuário final. `pdflatex` NÃO é usado. |
| Framework desktop | Wails **v2** (não v3) | Estabilidade. |
| Frontend | Vue 3 + `<script setup>` + TypeScript | — |
| Editor | Tiptap 2.x com nós customizados | — |
| Componentes | Reka UI (headless) | Controle total do visual. |
| Estado frontend | Pinia (apenas estado de sessão) | O Go é a fonte da verdade em disco. |
| Bundler | Vite | Padrão do template Wails+Vue. |
| Formato dos capítulos | JSON do ProseMirror (nunca HTML, nunca Markdown) | Mapeamento semântico sem ambiguidade. |

### Convenções de Desenvolvimento

- **Separação estrita:** o frontend desconhece a existência do LaTeX. Toda lógica de compilação vive em `internal/` no Go. Nenhuma string `\usepackage`, `\chapter` etc. pode existir em arquivos `.vue`/`.ts`.
- **Design headless:** nada de bibliotecas de componentes estilizadas. CSS próprio sobre Reka UI.
- **Comunicação Go→Vue por eventos Wails** (`runtime.EventsEmit`); **Vue→Go por métodos vinculados** (bindings).
- Código Go em inglês; textos de UI em **pt-BR**.
- Testes: `go test` para o conversor LaTeX e o Rule Engine (obrigatório); Vitest para lógica não trivial do frontend.

---

## 2. Anatomia da Interface (UX / UI)

Layout em zonas de foco, inspirado no VS Code:

1. **Activity Bar (esquerda, ~48px):** ícones verticais — Explorer, Build Targets, Extensions (Plugins), Settings.
2. **Sidebar (painel lateral, redimensionável, ocultável):** renderiza dinamicamente a árvore de capítulos ou os formulários de configuração conforme o ícone ativo.
3. **Editor Central ("Modo Fluxo"):** Tiptap com `width: 100%` e `max-width: 800px` centralizado. Sem bordas de página; apenas texto e marcações semânticas.
4. **Status Bar (rodapé, ~24px):** target ativo (clicável para trocar), contagem de palavras do capítulo, status do analisador (`✓ Sem problemas` / `✗ 2 problemas`).
5. **Problems Panel (painel inferior, recolhível):** lista de `Problem`s emitidos pelo Rule Engine, com severidade, mensagem e origem. Clicar em um problema navega para a origem quando aplicável.

---

## 3. Estrutura do Repositório (código do Luz Writer)

```text
luz-writer/
├── main.go                  # Bootstrap Wails
├── app.go                   # Struct App: métodos expostos ao frontend (só delega para internal/)
├── wails.json
├── go.mod
├── internal/
│   ├── workspace/           # Abrir/criar workspace, IO de project.json, capítulos, targets
│   ├── model/               # Structs compartilhadas (Project, Target, Problem, ...)
│   ├── plugins/             # Registro de plugins, manifests, geração de preâmbulo
│   ├── rules/               # Rule Engine
│   ├── latex/               # Conversor ProseMirror JSON → LaTeX
│   └── build/               # Orquestração do Tectonic, pipeline de compilação
├── frontend/
│   ├── src/
│   │   ├── components/      # ActivityBar.vue, Sidebar.vue, StatusBar.vue, ProblemsPanel.vue ...
│   │   ├── editor/          # Configuração Tiptap + extensões (nós luz*)
│   │   ├── stores/          # Pinia: workspace.ts, editor.ts, problems.ts
│   │   ├── forms/           # Renderizador de formulários schema-driven
│   │   └── App.vue
│   └── ...
└── testdata/                # Workspaces de exemplo para testes Go
```

---

## 4. Estrutura do Workspace do Usuário

Premissa "Single Source of Truth", centralizada na pasta oculta `.luz/`:

```text
meu-livro/
├── .luz/
│   ├── project.json         # Metadados comuns
│   ├── targets/
│   │   ├── amazon-6x9.json
│   │   └── ebook.json
│   ├── overrides/           # Configurações por página, criadas sob demanda (ver 5.6)
│   │   └── 00-dedicatoria.json
│   ├── styles.json          # Estilos personalizados (plugin customStyles; ver 5.7)
│   └── plugins.json         # Plugins habilitados
├── capitulos/               # TODOS os documentos de texto (capítulos E páginas especiais)
│   ├── 00-dedicatoria.json  # role: "dedication"
│   ├── 01-introducao.json   # role: "chapter"
│   ├── 02-metodologia.json  # role: "chapter"
│   └── 99-sobre-o-autor.json# role: "aboutAuthor"
├── imagens/                 # Imagens importadas (o app copia para cá via ImportImage)
│   └── grafico-01.png
├── anexos/                  # PDFs externos importados (capa pronta, ficha da editora...)
│   └── ficha-editora.pdf
├── dist/                    # Saída final (PDF/EPUB). Usuário não edita.
├── .tmp/                    # .tex gerado, cache do Tectonic
└── .gitignore               # Exclui .tmp/ e dist/
```

---

## 5. Schemas dos Arquivos (contratos de dados)

Todos os arquivos JSON gravados pelo Luz Writer incluem `"luzVersion": 1` para migrações futuras. O Go valida os schemas na leitura e devolve erro amigável se corrompidos.

### 5.1 `.luz/project.json`

```json
{
  "luzVersion": 1,
  "title": "Meu Livro",
  "subtitle": "",
  "authors": ["Nome do Autor"],
  "language": "pt-BR",
  "chapterOrder": ["01-introducao", "02-metodologia"],
  "activeTarget": "amazon-6x9",
  "variables": [
    { "name": "protagonista", "value": "Ana Clara" },
    { "name": "cidade-natal", "value": "Iguape" }
  ]
}
```

- `chapterOrder` define a ordem de compilação e da árvore no Explorer. Cada ID corresponde a `capitulos/<id>.json` — inclui **todos** os documentos, tanto capítulos quanto páginas especiais (dedicatória, prefácio etc.; ver seção 5.5).
- `language` é mapeado para `babel`/`polyglossia` pelo Go.
- `variables`: **Variáveis do Projeto** — pares nome/valor reutilizáveis no texto (equivalente amigável ao `\newcommand` de substituição). `name` é um slug único (letras minúsculas, números e hífen); `value` é texto simples. Gerenciadas num painel visual em *Settings* (lista com adicionar/renomear/editar/excluir); inseridas no texto via nó `luzVariable` (seção 6). Na exportação, o Go **expande o valor diretamente** durante a conversão (com o escape normal de LaTeX) — não gera `\newcommand`, evitando colisões de nome e mantendo a validação simples. Renomear ou alterar o valor propaga para todas as ocorrências.

### 5.2 `.luz/targets/<id>.json`

```json
{
  "luzVersion": 1,
  "id": "amazon-6x9",
  "name": "Amazon KDP 6\"x9\"",
  "kind": "print",
  "documentClass": "book",
  "fontSize": "11pt",
  "includeToc": true,
  "pluginConfig": {
    "geometry": {
      "paperWidth": "6in",
      "paperHeight": "9in",
      "marginInner": "0.75in",
      "marginOuter": "0.5in",
      "marginTop": "0.75in",
      "marginBottom": "0.75in",
      "mirrored": true
    },
    "bodyText": {
      "lineSpacing": "1.0",
      "paragraphStyle": "indent"
    },
    "fancyhdr": {
      "headerLeft": "author",
      "headerRight": "chapterTitle",
      "pageNumberPosition": "outer-footer"
    }
  }
}
```

- `kind`: `"print" | "ebook" | "article"`. Controla regras de validação e quais nós a toolbar do editor exibe (ex.: `article` esconde o botão "Capítulo").
- `documentClass`: `"book" | "report" | "article"`.
- `includeToc` (boolean, default `true` para `book`/`report`, `false` para `article`): gera `\tableofcontents` no início do mainmatter, respeitando o atributo `includeInToc` dos títulos. (A UI deste campo é implementada na Etapa 6.)
- `pluginConfig` guarda a configuração dos **módulos do núcleo** (sempre considerados; ver seção 8.0) e dos **plugins opcionais habilitados** em `plugins.json`; chaves de plugins opcionais desabilitados são ignoradas na compilação, mas preservadas no arquivo.

### 5.3 `.luz/plugins.json`

```json
{
  "luzVersion": 1,
  "enabled": ["fancyhdr"]
}
```

- `enabled` lista **apenas plugins opcionais**. Módulos do núcleo (`geometry`, `bodyText`; seção 8.0) nunca aparecem aqui — estão sempre ativos e não podem ser desativados.

### 5.4 `capitulos/<id>.json`

```json
{
  "luzVersion": 1,
  "id": "01-introducao",
  "role": "chapter",
  "language": null,
  "content": {
    "type": "doc",
    "content": [
      {
        "type": "luzChapter",
        "attrs": { "numbered": true, "includeInToc": true },
        "content": [{ "type": "text", "text": "Introdução" }]
      },
      {
        "type": "paragraph",
        "content": [
          { "type": "text", "text": "O conceito de " },
          { "type": "text", "marks": [{ "type": "italic" }], "text": "fluxo" },
          { "type": "text", "text": " é central" },
          {
            "type": "luzFootnote",
            "attrs": { "number": 1, "text": "Ver Csikszentmihalyi, 1990." }
          },
          { "type": "text", "text": "." }
        ]
      }
    ]
  }
}
```

- O arquivo do documento contém **apenas conteúdo** (texto puro). Configurações específicas da página vivem em `.luz/overrides/` (seção 5.6) — mantendo a premissa `.luz/` = configuração, `capitulos/` = texto.
- `language` (opcional, default `null`): idioma **deste documento inteiro**, quando diverso do idioma principal do projeto (ex.: um prefácio em inglês numa obra em português). `null`/ausente = herda o `language` do `project.json`. Idioma é propriedade **semântica do conteúdo** (como o `role`), não configuração de exportação — por isso vive aqui e não em `.luz/`. Editável por um seletor discreto no topo do editor (visível apenas com o plugin `languages` ativo; ver seções 6 e 8.3). No LaTeX vira `\begin{otherlanguage}{...} ... \end{otherlanguage}` envolvendo o documento.

### 5.5 Papéis de Documento (`role`) — Páginas Especiais e Front/Back Matter

Todo documento em `capitulos/` possui um campo `role`, escolhido na criação ("Novo documento" abre um seletor de tipo, com "Capítulo" como default). O `role` determina **automaticamente** em qual bloco do livro o documento é compilado e qual LaTeX o envolve — o escritor nunca escolhe "frontmatter" manualmente.

| `role` | Rótulo na UI | Bloco | LaTeX gerado (envoltório) |
|---|---|---|---|
| `chapter` | Capítulo | `\mainmatter` | conteúdo normal (nós `luzChapter` etc.) |
| `dedication` | Dedicatória | `\frontmatter` | página própria, `\thispagestyle{empty}`, texto em itálico alinhado à direita, deslocado ao terço superior |
| `epigraph` | Epígrafe | `\frontmatter` | página própria, `\thispagestyle{empty}`, citação centralizada; parágrafo final do documento é tratado como atribuição (alinhado à direita, precedido de travessão) |
| `acknowledgments` | Agradecimentos | `\frontmatter` | `\chapter*{Agradecimentos}` (fora do sumário por padrão) |
| `preface` | Prefácio | `\frontmatter` | `\chapter*{Prefácio}` + `\addcontentsline` |
| `aboutAuthor` | Sobre o Autor | `\backmatter` | `\chapter*{Sobre o Autor}` |
| `appendix` | Apêndice | `\backmatter` (após `\appendix`) | `\chapter{...}` → numeração "Apêndice A, B..." |

Regras de comportamento:

- **Ordem de compilação:** `chapterOrder` continua sendo a autoridade de ordem, mas o Go **agrupa por bloco** na hora de compilar: primeiro todos os documentos de frontmatter (na ordem em que aparecem), depois `\mainmatter` com os `chapter`, depois `\backmatter` com `appendix` e `aboutAuthor`. Se o usuário arrastar uma dedicatória para depois de um capítulo no Explorer, a regra R007 avisa que ela será reposicionada na compilação.
- **Explorer agrupado:** a Sidebar exibe três grupos visuais — "Pré-textuais", "Capítulos", "Pós-textuais" — e o arrastar-e-soltar reordena apenas dentro do grupo.
- **Editor sensível ao role:** em documentos com `role ≠ chapter/appendix`, a toolbar esconde os botões de Capítulo/Seção/Subseção (o título da página vem do próprio role); nós desses tipos, se presentes, geram warning na validação. Dedicatória e epígrafe também escondem listas e imagens — são páginas de texto curto.
- **Roles em `documentClass: "article"`:** targets de artigo ignoram o agrupamento matter (classes `article` não têm `\frontmatter`); documentos especiais viram seções não numeradas, e a regra R010 emite warning.

### 5.6 `.luz/overrides/<id-do-documento>.json` — Configurações por Página

Sobrescritas de configuração de plugins **válidas apenas para um documento específico** (ex.: dedicatória com margens maiores que as dos capítulos). Vivem em arquivo próprio dentro de `.luz/` — nunca dentro do arquivo de conteúdo — no mesmo formato da `pluginConfig` do target, porém parcial:

```json
{
  "luzVersion": 1,
  "documentId": "00-dedicatoria",
  "overrides": {
    "geometry": { "marginTop": "4cm", "marginBottom": "4cm" }
  }
}
```

Ciclo de vida (responsabilidade de `internal/workspace`):

- **Criação sob demanda:** o arquivo só passa a existir quando o usuário define a primeira sobrescrita da página; salvar sobrescritas vazias (`{}`) **apaga** o arquivo.
- **Cascata:** renomear (id) ou excluir um documento renomeia/exclui o arquivo de override correspondente na mesma operação.
- **Órfãos:** arquivo em `.luz/overrides/` cujo `documentId` não existe em `capitulos/` gera a regra R014 (não é apagado automaticamente — pode ser um capítulo temporariamente removido do projeto).
- **UI:** o arquivo **não aparece como nó na árvore** do Explorer. A presença de sobrescritas é indicada por um **badge ⚙** ao lado do nome do documento na árvore; o painel de edição é descrito na seção 8.4.

O sistema completo de camadas e precedência está na seção 8.4.

### 5.7 `.luz/styles.json` — Estilos Personalizados

Estilos de texto nomeados, **compostos visualmente** a partir de blocos seguros (equivalente amigável ao `\newcommand` de formatação). Criados/geridos pela UI do plugin `customStyles` (seção 8.3); o arquivo persiste mesmo com o plugin desativado.

```json
{
  "luzVersion": 1,
  "styles": [
    {
      "id": "termo-estrangeiro",
      "name": "Termo Estrangeiro",
      "italic": true,
      "bold": false,
      "smallCaps": false,
      "color": null
    },
    {
      "id": "voz-interior",
      "name": "Voz Interior",
      "italic": true,
      "bold": false,
      "smallCaps": true,
      "color": "#555555"
    }
  ]
}
```

- `id` é um slug único e estável (as marks no texto referenciam por `id`; renomear `name` não quebra nada).
- Propriedades componíveis no MVP: `italic`, `bold`, `smallCaps` (booleanos) e `color` (hex ou `null`). Nada de LaTeX cru.
- Na compilação, o Go gera um `\newcommand{\luzstyleTermoEstrangeiro}[1]{...}` por estilo no preâmbulo (nomes derivados do `id`, prefixados com `luzstyle` — sem risco de colisão com comandos LaTeX) e converte a mark para o comando. `color` exige `xcolor`, injetado pelo plugin quando algum estilo usa cor.
- Estilos são **semânticos e globais** (como variáveis e idiomas): valem para todos os targets. A aparência final pode variar levemente entre targets apenas pelo contexto tipográfico, nunca por configuração de estilo por target.

---

## 6. O Editor Tiptap — Definição dos Nós e Marks

### 6.0 Princípio: WYSIWYG Semântico — o editor mostra INTENÇÃO, não preview

**O visual do editor não espelha (e não deve tentar espelhar) a exportação.** O Modo Fluxo é uma interface de *captura de intenção autoral*: cada marcação semântica recebe um tratamento visual próprio, cuja função é tornar a intenção **reconhecível e distinta** — não antecipar como ela ficará no PDF. O único lugar onde o resultado tipográfico existe é o PDF em `dist/`.

Consequências práticas (regras rígidas para o frontend):

- O editor **nunca** renderiza geometria do target: sem bordas de página, sem quebras de página, sem margens do target, sem a fonte tipográfica do PDF. Trocar o target ativo **não altera** a aparência do texto no editor (apenas a toolbar, via `kind`).
- Cada nó/mark semântico tem um **visual de intenção** próprio, definido em CSS do editor, estável entre targets. Ele pode (e deve) divergir do resultado final — ex.: a citação inline aparece com fundo destacado no editor, mas no PDF vira apenas `\enquote{...}` com aspas tipográficas.
- Elementos de interação (badges, popovers, marcadores clicáveis) fazem parte do editor e obviamente não existem no PDF.
- É proibido implementar "modo preview de impressão" dentro do editor. Preview = compilar.

**Especificação canônica do visual de intenção da mark `luzInlineQuote`** (o padrão de referência para o estilo dos demais elementos):

```css
/* CSS do editor para a Mark luzInlineQuote */
.luz-inline-quote {
  background-color: rgba(0, 0, 0, 0.05); /* fundo cinza bem suave */
  padding: 0.1rem 0.3rem;
  border-radius: 4px;
  font-style: normal; /* garante que não se confunde com itálico */
}
.luz-inline-quote::before { content: "“"; color: #888; }
.luz-inline-quote::after  { content: "”"; color: #888; }
```

As aspas são **geradas pelo CSS** (pseudo-elementos), nunca digitadas nem armazenadas no conteúdo — assim a citação semântica não se confunde com aspas comuns digitadas pelo autor, e o Go emite as aspas corretas via `csquotes` conforme o idioma.

**Visual de intenção dos demais elementos** (mesma filosofia):

| Elemento | Visual no editor (intenção) | O que NÃO fazer |
|---|---|---|
| `luzChapter` / `luzSection` / `luzSubsection` | Tamanhos de fonte escalonados + um rótulo discreto à esquerda ("Capítulo", "Seção"...) visível ao focar; ícone/badge quando `numbered: false` ou `includeInToc: false` | Não simular a página de abertura de capítulo do PDF |
| `luzFootnote` | Marcador sobrescrito clicável `[1]`, `[2]`... renumerado automaticamente; clique abre Popover com o texto da nota | Não exibir o texto da nota no rodapé do editor |
| `luzImage` | Imagem renderizada com moldura leve + legenda editável abaixo em estilo distinto | Não calcular o tamanho real que terá na página do target |
| `blockquote` | Recuo com barra vertical à esquerda (padrão de editores) | Não aplicar a formatação final do ambiente `quote` |
| `luzLang` (mark) | Sublinhado pontilhado cinza + código do idioma sobrescrito discreto ao fim do trecho, gerado por CSS a partir de `data-lang` | Não trocar fonte/tipografia para "parecer" o idioma; não armazenar o código como texto |
| `luzSoftHyphen` (nó) | Glifo fino visível `‧` em cinza no ponto de quebra sugerido | Não usar o caractere invisível U+00AD (intenção invisível viola o princípio); não simular a quebra de linha |
| `luzVariable` (nó) | *Chip* com fundo suave e cantos arredondados exibindo o valor atual + ícone discreto de variável; hover mostra o `name`; referência quebrada = chip vermelho | Não renderizar como texto comum (o autor precisa distinguir texto digitado de valor propagado) |
| `luzCustomStyle` (mark) | Aplica a composição do estilo literalmente (itálico/negrito/versalete/cor — como o alinhamento, intenção e visual coincidem) + sublinhado pontilhado sutil para distinguir de formatação manual; estilo excluído = fundo vermelho suave | Não permitir edição da composição inline (só no gerenciador do plugin) |
| `bold` / `italic` / `underline` / `strike` | Convencionais (negrito, itálico, sublinhado, tachado) | — |
| Alinhamento de parágrafo (`align`) | Mostrado **literalmente** (texto centralizado aparece centralizado) — caso raro em que intenção e visual coincidem, pois o alinhamento independe da geometria do target | Não recalcular quebras de linha do PDF |

CSS de referência para a mark `luzLang` (mesmo espírito do CSS canônico acima):

```css
.luz-lang { border-bottom: 1px dotted #888; }
.luz-lang::after {
  content: attr(data-lang);
  font-size: 0.6em;
  vertical-align: super;
  color: #888;
  margin-left: 2px;
  text-transform: uppercase;
}
```

A interface descarta headings genéricos (`H1`, `H2`) em prol de nós com nomenclatura editorial. Definições exatas:

### Nós de bloco

| Nó | Grupo | Atributos | LaTeX gerado |
|---|---|---|---|
| `luzChapter` | block, conteúdo inline | `numbered: boolean` (default `true`), `includeInToc: boolean` (default `true`) | `\chapter{...}` / `\chapter*{...}` (+ `\addcontentsline` se `numbered=false` e `includeInToc=true`) |
| `luzSection` | block, conteúdo inline | idem | `\section{...}` / `\section*{...}` |
| `luzSubsection` | block, conteúdo inline | idem | `\subsection{...}` / `\subsection*{...}` |
| `paragraph` | block padrão do Tiptap, estendido | `align: string` — `"justify"` (default) \| `"left"` \| `"center"` \| `"right"` | texto + linha em branco; alinhamento ≠ justify envolve com `\begin{Center}` / `\begin{FlushLeft}` / `\begin{FlushRight}` (ambientes do `ragged2e`, com hifenização preservada) |
| `blockquote` | block padrão | — | `\begin{quote} ... \end{quote}` |
| `bulletList` / `orderedList` / `listItem` | padrão | — | `itemize` / `enumerate` / `\item` |
| `luzImage` | block, atômico | `src: string` (caminho relativo, ex. `imagens/grafico-01.png`), `caption: string` (default `""`), `width: number` (percentual da largura do texto, 10–100, default `80`) | Com caption: `\begin{figure}[htbp]\centering\includegraphics[width=0.NN\linewidth]{...}\caption{...}\end{figure}`. Sem caption: `\begin{center}\includegraphics[...]{...}\end{center}` (sem numeração de figura) |

- Os atributos `numbered`/`includeInToc` são editados por **checkboxes num pequeno popover** que aparece ao focar o nó de título (sem poluir o texto).
- Se o target ativo tiver `kind: "article"`, a toolbar esconde "Capítulo" e o Rule Engine sinaliza `luzChapter` existentes como erro (regra R002).
- **Fluxo de imagem:** o botão "Imagem" da toolbar chama `ImportImage()` (diálogo nativo). O Go copia o arquivo para `imagens/` (renomeando em caso de colisão) e retorna o caminho relativo, que vira o `src` do nó. O editor exibe a imagem via asset handler do Wails; `caption` e `width` são editados em popover ao selecionar a imagem. O nó nunca referencia caminhos absolutos ou externos ao workspace.

### Nós inline

| Nó | Atributos | Comportamento na UI | LaTeX |
|---|---|---|---|
| `luzFootnote` | `number: number` (renumerado automaticamente pelo editor na ordem do documento), `text: string` (texto simples da nota) | Renderizado como marcador sobrescrito `[n]`. Clique abre Popover (Reka UI) com textarea para editar `text`. | `\footnote{...}` |
| `luzSoftHyphen` | — (atômico) | Ponto de quebra de palavra **sugerido** pelo autor, para palavras que o compilador não sabe hifenizar (nomes próprios, termos técnicos). Inserido pela toolbar ou atalho `Ctrl/Cmd+Shift+-`, exibido como glifo `‧` (seção 6.0). Visível apenas com o plugin `hyphenation` ativo. | `\-` |
| `luzVariable` | `name: string` (slug da variável em `project.json`) | **Variável do Projeto** (seção 5.1). Renderizada como *chip* exibindo o **valor atual** da variável; inserida via dropdown "Inserir variável" na toolbar (lista as variáveis definidas). Alterar o valor no painel de Settings atualiza todos os chips instantaneamente. Variável inexistente → chip em vermelho (referência quebrada). | valor expandido pelo Go como texto normal (com escape) — sem `\newcommand` |

### Marks

| Mark | LaTeX |
|---|---|
| `bold` | `\textbf{...}` |
| `italic` | `\textit{...}` |
| `underline` | `\uline{...}` (requer `ulem` — sempre incluído no preâmbulo base) |
| `strike` | `\sout{...}` (requer `ulem` — sempre incluído no preâmbulo base) |
| `luzInlineQuote` | `\enquote{...}` (requer `csquotes` — sempre incluído). Na UI: visual de intenção definido na seção 6.0 (fundo suave + aspas via CSS). |
| `luzLang` (atributo `lang: string`, código BCP-47 ex. `"en"`, `"fr"`, `"la"`) | `\foreignlanguage{<nome babel>}{...}`. O Go mantém a tabela de mapeamento código → nome babel (`en`→`english`, `fr`→`french`, `de`→`german` (ngerman), `es`→`spanish`, `it`→`italian`, `la`→`latin`, `pt-BR`→`brazilian`); esses são os idiomas suportados no MVP e listados no dropdown da toolbar. Botão visível apenas com o plugin `languages` ativo. |
| `luzCustomStyle` (atributo `styleId: string`, referencia `.luz/styles.json`) | `\luzstyle<IdEmCamelCase>{...}` — comando gerado pelo Go no preâmbulo a partir da composição do estilo (seção 5.7). Dropdown "Estilos" na toolbar, visível apenas com o plugin `customStyles` ativo. |

### Escape de LaTeX (obrigatório)

Todo texto vindo do usuário passa por escape antes de entrar no `.tex`. Tabela mínima: `\ { } $ & # ^ _ % ~` → `\textbackslash{}`, `\{`, `\}`, `\$`, `\&`, `\#`, `\^{}`, `\_`, `\%`, `\textasciitilde{}`. Escrever testes unitários para isso.

### Exemplo canônico de conversão (usar como teste)

O capítulo da seção 5.4 deve gerar exatamente:

```latex
\chapter{Introdução}

O conceito de \textit{fluxo} é central\footnote{Ver Csikszentmihalyi, 1990.}.
```

---

## 7. Contrato da Ponte Wails (Go ⇄ Vue)

### 7.1 Métodos expostos (Vue → Go)

Todos no struct `App` (arquivo `app.go`), delegando para `internal/`. Assinaturas exatas:

```go
// Workspace
func (a *App) OpenWorkspaceDialog() (model.WorkspaceInfo, error)      // abre diálogo nativo de pasta
func (a *App) CreateWorkspace(path string, title string, author string, language string) (model.WorkspaceInfo, error)
func (a *App) GetProject() (model.Project, error)
func (a *App) SaveProject(p model.Project) error

// Capítulos
func (a *App) ListChapters() ([]model.ChapterMeta, error)             // {id, title, role, wordCount}
func (a *App) LoadChapter(id string) (string, error)                  // JSON ProseMirror como string
func (a *App) SaveChapter(id string, contentJSON string) error
func (a *App) CreateChapter(title string, role string) (model.ChapterMeta, error) // gera id slug "NN-titulo"; role da seção 5.5
func (a *App) DeleteChapter(id string) error
func (a *App) ReorderChapters(order []string) error

// Configurações por página (seção 5.6)
func (a *App) GetDocumentOverrides(id string) (string, error)               // JSON das sobrescritas; "{}" se o arquivo não existir
func (a *App) SaveDocumentOverrides(id string, overridesJSON string) error  // "{}" apaga o arquivo em .luz/overrides/

// Assets
func (a *App) ImportImage() (string, error)        // diálogo nativo → copia para imagens/ → retorna caminho relativo
func (a *App) ImportAttachmentPDF() (string, error) // diálogo nativo → copia para anexos/ → retorna caminho relativo (usado pelo plugin pdfpages)

// Targets
func (a *App) ListTargets() ([]model.Target, error)
func (a *App) ListTargetPresets() ([]model.Target, error) // modelos embutidos no binário (ver seção 8.5)
func (a *App) SaveTarget(t model.Target) error
func (a *App) DeleteTarget(id string) error
func (a *App) SetActiveTarget(id string) error

// Plugins
func (a *App) ListAvailablePlugins() ([]model.PluginManifest, error)  // catálogo embutido no binário
func (a *App) SetPluginEnabled(name string, enabled bool) error

// Estilos personalizados (seção 5.7; UI visível apenas com o plugin customStyles ativo)
func (a *App) ListStyles() ([]model.CustomStyle, error)
func (a *App) SaveStyles(styles []model.CustomStyle) error            // substitui o conjunto inteiro em .luz/styles.json

// Validação e Build
func (a *App) Validate() ([]model.Problem, error)
func (a *App) Compile() (model.BuildResult, error)                    // usa o target ativo; assíncrono internamente
```

### 7.2 Eventos (Go → Vue)

| Evento | Payload | Quando |
|---|---|---|
| `luz:problems` | `[]Problem` | Após qualquer `Validate()` ou save que dispare revalidação |
| `luz:build:progress` | `{ stage: string, percent: number }` | Durante a compilação (`stage`: `"validating" \| "generating" \| "compiling" \| "done"`) |
| `luz:build:done` | `BuildResult` | Fim da compilação |

### 7.3 Structs compartilhadas (`internal/model`)

```go
type Problem struct {
    Severity string `json:"severity"` // "error" | "warning" | "info"
    Code     string `json:"code"`     // ex.: "R001"
    Message  string `json:"message"`  // em pt-BR
    Source   string `json:"source"`   // "project" | "target:<id>" | "chapter:<id>" | "plugin:<name>" | "styles" | "override:<docId>"
}

type BuildResult struct {
    Success    bool      `json:"success"`
    OutputPath string    `json:"outputPath"` // caminho em dist/
    Problems   []Problem `json:"problems"`
    LogTail    string    `json:"logTail"`    // últimas ~40 linhas do log do Tectonic em caso de falha
}
```

Regra: **erros (`severity: "error"`) bloqueiam a compilação**; warnings não.

---

## 8. Arquitetura de Plugins (Schema-Driven UI)

Pacotes LaTeX são apresentados como **Módulos de Publicação**, divididos em duas categorias (seção 8.0). O fluxo dos opcionais:

1. Na aba *Extensions*, o usuário liga/desliga um módulo → `plugins.json` é atualizado via `SetPluginEnabled`.
2. O Vue relê o estado e renderiza uma nova aba de configuração em *Build Targets*, gerada a partir do `FormSchema` do plugin (formulários abstraídos: switches, dropdowns, inputs numéricos com unidade — nunca LaTeX cru).
3. Na compilação, o Go itera sobre os módulos ativos (padrão *Builder*) e injeta `\usepackage{}` + configurações no preâmbulo.

### 8.0 Módulos do Núcleo vs Plugins Opcionais

| | Módulos do **núcleo** | Plugins **opcionais** |
|---|---|---|
| Exemplos | `geometry`, `bodyText` | `fancyhdr`, `catalogRecord`, `pdfpages`, `languages`, `hyphenation` |
| Podem ser desativados? | **Não** — sem eles não existe livro (toda página tem dimensão; todo corpo de texto tem espaçamento) | Sim, via aba *Extensions* |
| Aparecem em `plugins.json`? | Nunca | Sim, quando habilitados |
| Configuráveis por target (`pluginConfig`)? | Sim | Sim |
| Override por página? | Conforme `DocumentScope()` | Conforme `DocumentScope()` |
| Na aba *Extensions* | Listados numa seção "Essenciais", com badge e switch travado (transparência sem desativação) | Seção normal, com switch |
| Abas em *Build Targets* | Sempre presentes (fixas) | Só quando habilitados |

Ambas as categorias implementam a **mesma interface `Plugin`** (seção 8.1) — a distinção é o método `Core() bool`. Módulos do núcleo ignoram o portão da camada 2 (seção 8.4): são sempre considerados habilitados pelo resolvedor, pela validação e pela compilação.

Além dos módulos, o **preâmbulo base** carrega pacotes invisíveis ao usuário (sem UI alguma): `microtype` (refinamento tipográfico automático), `hyperref` com `hidelinks` (sumário clicável e bookmarks no PDF — carregado por último, como o pacote exige), `ragged2e` (alinhamentos com hifenização preservada), `csquotes`, `ulem`, `graphicx` e `babel`. Ver lista completa na seção 8.3.

### 8.1 Interface Go de todo plugin

```go
type Plugin interface {
    Name() string                                   // "geometry"
    DisplayName() string                            // "Geometria de Página" (pt-BR)
    Description() string
    Core() bool                                     // true = módulo do núcleo (seção 8.0): sempre ativo, não desativável
    Schema() model.FormSchema                       // define o formulário no frontend
    DefaultConfig() json.RawMessage
    Validate(cfg json.RawMessage, ctx model.BuildContext) []model.Problem
    Preamble(cfg json.RawMessage, ctx model.BuildContext) (string, error)

    // Suporte a sobrescrita por documento (seção 8.4).
    // DocumentScope() == false → o plugin só atua globalmente (ex.: catalogRecord, pdfpages)
    // e o painel "Página" não o exibe.
    DocumentScope() bool
    // Chamado apenas quando o documento possui override para este plugin.
    // Recebe a config EFETIVA (default + target + override já mesclados) e devolve
    // o LaTeX a inserir antes e depois do conteúdo do documento.
    // Ex. geometry: before = `\newgeometry{...}`, after = `\restoregeometry`.
    // Ex. fancyhdr: before = pagestyle temporário, after = restauração do pagestyle do target.
    ScopedLaTeX(effectiveCfg json.RawMessage, ctx model.BuildContext) (before string, after string, err error)
}
```

`BuildContext` carrega `Project`, `Target` ativo e a lista de módulos ativos (núcleo + opcionais habilitados — para regras cruzadas).

### 8.2 FormSchema (contrato do renderizador de formulários)

```json
{
  "fields": [
    { "key": "paperWidth", "label": "Largura da página", "type": "dimension", "default": "6in" },
    { "key": "mirrored",   "label": "Margens espelhadas (frente/verso)", "type": "switch", "default": true },
    { "key": "pageNumberPosition", "label": "Posição do número de página", "type": "select",
      "options": [
        { "value": "outer-footer", "label": "Rodapé externo" },
        { "value": "center-footer", "label": "Rodapé central" }
      ]
    }
  ]
}
```

Tipos de campo suportados no MVP: `switch`, `select`, `dimension` (número + unidade `in|cm|mm|pt`), `text`, `number`. O componente Vue `SchemaForm.vue` renderiza qualquer `FormSchema` genericamente com Reka UI.

### 8.3 Módulos do MVP

**Módulos do núcleo** (`Core() == true` — sempre ativos):

| Módulo | Pacote(s) LaTeX | Escopo no MVP |
|---|---|---|
| `geometry` | `geometry` | Dimensões de página, margens, espelhamento. `DocumentScope() == true`. Defaults sensatos (A5, margens 2,5cm, espelhado) para o target funcionar sem configuração. |
| `bodyText` | `setspace` + comandos de parágrafo | **Corpo do Texto.** `FormSchema`: `lineSpacing` (select: `"1.0"` \| `"1.15"` \| `"1.5"` \| `"2.0"`; default `"1.0"` → `\setstretch{}`), `paragraphStyle` (select: `"indent"` = recuo na 1ª linha, padrão em livros no Brasil \| `"spaced"` = sem recuo, com espaço entre parágrafos → `\parindent`/`\parskip`). `DocumentScope() == true` — ex.: uma página especial com espaçamento próprio. |

**Plugins opcionais** (liga/desliga em *Extensions*):

| Plugin | Pacote(s) LaTeX | Escopo no MVP |
|---|---|---|
| `fancyhdr` | `fancyhdr` | Cabeçalhos/rodapés, posição de numeração |
| `catalogRecord` | ambiente próprio (minipage emoldurada no padrão ABNT) | Ficha catalográfica **gerada pelo app** a partir de formulário. Renderizada no verso da página de rosto (frontmatter). `FormSchema`: `isbn` (text), `publisher` (text, "Editora"), `publisherCity` (text, "Cidade"), `year` (number, "Ano"), `cdd` (text, "CDD"), `cdu` (text, "CDU", opcional), `subjectEntries` (text, "Assuntos — separados por ponto-e-vírgula"), `preparedBy` (text, "Elaborada por (nome e CRB)", opcional). Autor e título vêm de `project.json` automaticamente. |
| `pdfpages` | `pdfpages` | Inserção de **PDFs externos prontos** (ex.: capa diagramada, ficha catalográfica fornecida pela editora). `FormSchema`: `filePath` (text somente-leitura preenchido via botão que chama `ImportAttachmentPDF()`), `placement` (select: `"beforeFrontmatter"` = antes de tudo \| `"afterFrontmatter"` = após pré-textuais \| `"afterBackmatter"` = última(s) página(s)), `pages` (text, default `"-"` = todas). Gera `\includepdf[pages={...}]{anexos/...}`. Se o usuário já tem a ficha da editora em PDF, usa este plugin **em vez** do `catalogRecord`. |
| `languages` | `babel` (idiomas adicionais) | **Suporte multilíngue.** É o portão de UI para conteúdo em idioma diverso: habilita o botão/dropdown de idioma na toolbar (mark `luzLang`) e o seletor de idioma do documento (campo `language`, seção 5.4). `FormSchema`: vazia (`fields: []`) — não há o que configurar por target: na compilação o Go **escaneia o conteúdo**, coleta todos os idiomas usados (marks + documentos) e injeta os idiomas adicionais no `babel` automaticamente. `DocumentScope() == false`. |
| `hyphenation` | `hyphenat` | **Controle de hifenização.** A hifenização automática por idioma é comportamento nativo do LaTeX; este plugin a controla e habilita o botão de hífen sugerido (`luzSoftHyphen`) na toolbar. `FormSchema`: `mode` (select: `"auto"` = hifenização automática \| `"off"` = desligada → `\usepackage[none]{hyphenat}`; default `"auto"`). `DocumentScope() == true` — permite, por exemplo, desligar hifenização apenas numa página especial via painel "Configurações desta página". |
| `customStyles` | `xcolor` (apenas se algum estilo usar cor) | **Estilos Personalizados** (seção 5.7) — o equivalente amigável ao `\newcommand` de formatação. Habilita: (a) o **gerenciador de estilos**, aberto pelo botão "Gerenciar estilos..." no card do plugin na aba *Extensions* (lista de estilos; cada um editado com switches `italic`/`bold`/`smallCaps` e seletor de cor opcional — composição visual, nunca LaTeX cru); (b) o dropdown "Estilos" na toolbar (mark `luzCustomStyle`). `FormSchema` de target: vazia — estilos são globais (`.luz/styles.json`), não por target. Na compilação gera um `\newcommand` por estilo no preâmbulo. `DocumentScope() == false`. |

O preâmbulo **base** (sempre presente, sem UI) inclui: `documentclass` do target, `fontenc`/`inputenc` ou setup nativo do XeLaTeX/Tectonic, `babel` conforme `language`, `csquotes`, `ulem` (sublinhado/tachado), `graphicx` (com `\graphicspath{{imagens/}}` — necessário para o nó `luzImage`), `ragged2e` (alinhamentos de parágrafo), `microtype` (refinamento tipográfico automático) e `hyperref` com opção `hidelinks` (sumário clicável e bookmarks — **carregado por último**, como o pacote exige).

### 8.4 Camadas de Configuração e Precedência

Existem **quatro camadas** de configuração, cada uma com um dono e um arquivo. Esta hierarquia é lei para todo o sistema (UI, validação e compilação):

| # | Camada | Onde vive | Escopo | Exemplo |
|---|---|---|---|---|
| 1 | **Projeto** | `.luz/project.json` | Vale para todos os targets e documentos | Título, autores, idioma |
| 2 | **Disponibilidade** (plugins habilitados) | `.luz/plugins.json` | Define **quais configurações existem** na UI (só para opcionais; núcleo sempre existe) | `fancyhdr` ligado → formulários de cabeçalho passam a existir |
| 3 | **Target** | `.luz/targets/<id>.json` → `pluginConfig` | Valores por formato de exportação | "Amazon KDP 6×9" tem página 6×9in; "e-Book" não tem margem espelhada |
| 4 | **Documento** | `.luz/overrides/<id>.json` (seção 5.6) | Sobrescritas para uma página/documento específico | Dedicatória com margens maiores |

**Resolução da configuração efetiva** (por plugin, por documento, calculada no Go em `internal/plugins/resolve.go`):

```text
efetiva = DefaultConfig()  ⊕  pluginConfig[plugin] (target ativo)  ⊕  overrides[plugin] (documento)
```

onde `⊕` é merge raso campo a campo: um campo definido na camada mais específica vence o da camada mais genérica; campos ausentes herdam. Regras derivadas:

- A camada 2 é um **portão** para plugins opcionais: se o plugin não está em `enabled`, as camadas 3 e 4 daquele plugin são ignoradas na compilação (mas preservadas nos arquivos, para o usuário poder religar sem perder valores). **Módulos do núcleo (`Core() == true`) ignoram o portão** — são sempre considerados habilitados.
- `Validate()` e `Preamble()`/`ScopedLaTeX()` sempre recebem a config **efetiva**, nunca camadas cruas — plugins não conhecem a hierarquia.
- Overrides são independentes do target (aplicam-se sobre qualquer target ativo). Se um override não faz sentido num target, é papel do Rule Engine avisar, não da UI esconder silenciosamente.
- Documentos **não** podem sobrescrever a camada 1 nem ligar/desligar plugins.

**UI da camada 4 — painel "Configurações desta página":** o painel lista somente os módulos **ativos** (núcleo + opcionais habilitados) com `DocumentScope() == true`, reutilizando o mesmo `SchemaForm.vue` em **modo override**: cada campo tem um toggle "sobrescrever"; desligado, o campo aparece esmaecido exibindo o valor herdado do target ativo (com a indicação "herdado de <nome do target>"); ligado, o valor digitado é gravado via `SaveDocumentOverrides`. Um botão "Limpar sobrescritas" salva `{}` (o que apaga o arquivo). Pontos de acesso: (a) **badge ⚙** ao lado do documento na árvore do Explorer — presente apenas quando existem sobrescritas, clicável; (b) item "Configurações desta página..." no **menu de contexto** de qualquer documento da árvore; (c) botão no topo do editor com o documento aberto. O arquivo de override em si nunca é exposto como nó da árvore — o escritor interage só com o painel.

**Compilação:** para documentos com `overrides` não vazio, o pipeline envolve o conteúdo com `before`/`after` de `ScopedLaTeX()` de cada plugin sobrescrito (na ordem de `plugins.json`), aninhando corretamente (o último a abrir é o primeiro a fechar).

### 8.5 Presets de Target

Para o usuário não montar targets do zero, o binário embute um catálogo de **presets** (arquivos JSON de target completos, em `internal/plugins/presets/`), expostos via `ListTargetPresets()`. O diálogo "Novo target" oferece: começar de um preset (que é copiado para `.luz/targets/` e vira 100% editável) ou começar vazio. Presets do MVP:

| Preset | kind | Conteúdo |
|---|---|---|
| Amazon KDP 6×9 (miolo) | print | `book`, 11pt, geometry 6×9in espelhado, fancyhdr padrão |
| Livro A5 genérico | print | `book`, 11pt, geometry A5 espelhado |
| e-Book (PDF fluido) | ebook | `book`, 12pt, geometry margens uniformes estreitas, sem espelhamento |
| Artigo A4 | article | `article`, 12pt, geometry A4 |

Presets são apenas pontos de partida — depois de copiados, não há vínculo com o original (sem "atualização de preset").

---

## 9. Rule Engine — Regras do MVP

O Rule Engine roda em `Validate()` e antes de toda compilação. Cada regra tem código estável, severidade e mensagem em pt-BR. Regras iniciais:

| Código | Severidade | Condição | Mensagem |
|---|---|---|---|
| R001 | error | `geometry.mirrored = true` em target `kind: "ebook"` | "Margens espelhadas não fazem sentido em e-books (fluxo contínuo)." |
| R002 | error | Nó `luzChapter` presente com target `documentClass: "article"` | "Artigos não possuem capítulos. Converta para Seção ou mude o target." |
| R003 | warning | Plugin `catalogRecord` habilitado em target `kind: "ebook"` | "Ficha Catalográfica ativada em um target de e-Book." |
| R004 | warning | Target `kind: "print"` com `geometry.mirrored: false` | "Livros impressos frente-e-verso geralmente usam margens espelhadas — confirme se a impressão será apenas frente." |
| R005 | warning | Capítulo referenciado em `chapterOrder` com conteúdo vazio | "O capítulo '<id>' está vazio." |
| R006 | error | `chapterOrder` referencia arquivo inexistente em `capitulos/` | "Capítulo '<id>' listado no projeto mas o arquivo não existe." |
| R007 | warning | Documento de frontmatter (`dedication`, `epigraph`, `acknowledgments`, `preface`) posicionado depois de um `chapter` em `chapterOrder` | "A '<rótulo>' será reposicionada para o início do livro na compilação." |
| R008 | error | Nó `luzImage` cujo `src` não existe em `imagens/` | "Imagem '<src>' não encontrada no capítulo '<id>'." |
| R009 | error | Plugin `pdfpages` habilitado com `filePath` vazio ou apontando para arquivo inexistente em `anexos/` | "O anexo PDF configurado não foi encontrado." |
| R010 | warning | Documento com role especial (`dedication`, `epigraph`, `aboutAuthor`...) em target `documentClass: "article"` | "Artigos normalmente não possuem <rótulo>; o conteúdo será inserido como seção simples." |
| R011 | warning | Plugins `catalogRecord` e `pdfpages` (com placement frontmatter) ativos simultaneamente | "Você tem uma ficha catalográfica gerada e um PDF anexo pré-textual ao mesmo tempo — verifique se não haverá duplicidade." |
| R012 | warning | Documento com `overrides` para plugin desabilitado em `plugins.json` | "O documento '<id>' sobrescreve configurações de '<plugin>', que está desativado — as sobrescritas serão ignoradas." |
| R013 | warning | Documento com override de `geometry` em target `kind: "ebook"` | "Sobrescritas de geometria por página têm pouco efeito em e-books de fluxo contínuo." |
| R014 | warning | Arquivo em `.luz/overrides/` cujo `documentId` não existe em `capitulos/` | "Há configurações de página para o documento '<id>', que não existe mais no projeto." |
| R015 | warning | Mark `luzLang` ou documento com `language` diverso do projeto, com o plugin `languages` desativado | "Há conteúdo marcado em outro idioma, mas o Suporte Multilíngue está desativado — as marcações serão ignoradas na exportação." |
| R016 | warning | Nó `luzSoftHyphen` presente com o plugin `hyphenation` desativado ou em modo `"off"` | "Há pontos de hifenização sugeridos, mas a hifenização está desativada — eles serão ignorados na exportação." |
| R017 | error | Nó `luzVariable` cujo `name` não existe em `project.json → variables` | "A variável '<name>' usada no documento '<id>' não existe mais no projeto." |
| R018 | warning | Mark `luzCustomStyle` cujo `styleId` não existe em `.luz/styles.json`, ou presente com o plugin `customStyles` desativado | "O estilo '<styleId>' não está disponível — o trecho será exportado como texto sem formatação de estilo." |

Novas regras devem seguir o mesmo padrão (arquivo por regra em `internal/rules/`, registradas num slice central, cada uma com teste unitário).

---

## 10. Pipeline de Compilação (Go)

`Compile()` executa, emitindo progresso via eventos:

1. **Validação:** roda o Rule Engine. Qualquer `error` → aborta e retorna `BuildResult{Success: false, Problems: ...}`.
2. **Geração do preâmbulo:** base + iteração sobre os módulos ativos — núcleo primeiro (`geometry`, `bodyText`), depois os opcionais na ordem declarada em `plugins.json` (`Preamble()` de cada um, recebendo a config **efetiva** da seção 8.4 — sem override, pois preâmbulo é global). Com o plugin `languages` ativo, o Go escaneia previamente todo o conteúdo, coleta os idiomas usados (marks `luzLang` + campos `language` dos documentos) e injeta os idiomas adicionais no `babel`, mantendo o idioma principal do projeto como default. Com o plugin `customStyles` ativo, um `\newcommand` por estilo de `.luz/styles.json` é gerado (seção 5.7).
3. **Parsing do conteúdo:** os documentos de `chapterOrder` são agrupados por bloco conforme o `role` (seção 5.5): frontmatter → `\mainmatter` + capítulos → `\backmatter` (+ `\appendix` antes do primeiro apêndice). Cada documento é convertido de JSON ProseMirror para LaTeX (`internal/latex`), com o envoltório do seu role; nós `luzVariable` são expandidos para o valor atual da variável (com escape normal); documentos com arquivo de sobrescritas em `.luz/overrides/` ganham adicionalmente o envoltório `before`/`after` de `ScopedLaTeX()` (seção 8.4). Nós desconhecidos → `Problem` warning + nó ignorado (nunca pânico).
4. **Escrita:** `.tmp/main.tex` consolidado.
5. **Execução silenciosa:** `tectonic .tmp/main.tex --outdir dist/` via `os/exec`, **com diretório de trabalho na raiz do workspace** (para que `\graphicspath{{imagens/}}` e `anexos/` resolvam), timeout de 120s. Stdout/stderr capturados; em falha, as últimas ~40 linhas vão para `BuildResult.LogTail` (o usuário nunca vê LaTeX, mas o log fica disponível num painel "detalhes técnicos" recolhido).
6. **Entrega:** renomeia para `dist/<title-slug>-<target-id>.pdf` e emite `luz:build:done`.

Pré-requisito de ambiente: o binário `tectonic` deve estar no PATH. Na inicialização, o Go verifica com `exec.LookPath("tectonic")`; ausente → `Problem` info persistente na Status Bar ("Tectonic não encontrado — exportação desabilitada"), sem travar o app.

---

## 11. Ferramentas de Apoio ao Agente

O ambiente de desenvolvimento tem duas ferramentas instaladas que você (agente) **deve** usar. Elas existem para você não se perder no código e não estourar contexto com saídas verbosas.

### 11.1 Graphify — mapa navegável do projeto

O Graphify transforma o repositório em um grafo de conhecimento consultável (`graphify-out/graph.json` + `GRAPH_REPORT.md`).

**Regras de uso:**

1. **No início do projeto (Etapa 0):** rode `/graphify .` para indexar este documento de especificação e o esqueleto inicial.
2. **Ao final de cada etapa concluída:** rode `/graphify . --update` para reindexar apenas o que mudou, mantendo o grafo sincronizado antes da próxima etapa.
3. **Antes de procurar código com grep/leitura de arquivos, prefira consultar o grafo:**
   - `graphify query "onde o preâmbulo LaTeX é montado?"`
   - `graphify query "o que conecta o Rule Engine ao Problems Panel?"`
   - `graphify explain "Plugin"` — para entender uma entidade específica.
   - `graphify path "Compile" "Tectonic"` — para rastrear um fluxo entre dois pontos.
4. **Ao retomar trabalho em uma nova sessão:** leia `graphify-out/GRAPH_REPORT.md` antes de tocar no código, para recuperar a visão de arquitetura sem reler todos os arquivos.
5. Consultas escopadas (`graphify query`) são preferíveis a ler o relatório inteiro; use o relatório apenas para revisão ampla de arquitetura.

### 11.2 RTK — filtro de saídas de comandos

O RTK comprime a saída de comandos de terminal em 60–90% antes de ela chegar ao seu contexto. Ele está instalado com o hook de reescrita automática (`rtk init -g`), então comandos Bash comuns (`git status`, `ls`, `grep`, `go test`...) já são reescritos para `rtk <cmd>` transparentemente.

**Regras de uso:**

1. **Prefira comandos de shell** (`cat`, `rg`, `find`, `git ...`) em vez das ferramentas nativas Read/Grep/Glob quando a saída for potencialmente grande — o hook do RTK só intercepta chamadas Bash, e é aí que a economia acontece. Para leituras pontuais, use explicitamente `rtk read <arquivo>`, `rtk grep <padrão> .` e `rtk find "*.go" .`.
2. Para builds e testes deste projeto, use as formas filtradas:
   - `rtk go test ./...` — mostra apenas falhas.
   - `rtk test wails build` e `rtk test npm run build` — wrapper genérico que reporta só erros.
   - `rtk lint` para ESLint no frontend.
3. **Quando um comando falhar,** o RTK salva a saída completa em `~/.local/share/rtk/tee/…` e imprime o caminho. Leia esse arquivo de log (com `rtk read`) em vez de reexecutar o comando para "ver a saída inteira".
4. Nunca desative ou contorne o hook do RTK para "ver mais saída"; use o log tee do item anterior.

---

## 12. Plano de Construção em Etapas

> **Protocolo obrigatório:** implemente **UMA etapa por vez**. Ao concluir uma etapa: (a) rode os critérios de aceite; (b) rode `/graphify . --update`; (c) apresente um resumo curto do que foi feito e como o usuário pode testar manualmente; (d) **PARE e aguarde a aprovação do usuário** antes de iniciar a etapa seguinte. Se um critério de aceite não puder ser cumprido, explique o bloqueio em vez de improvisar mudanças de escopo.

### Etapa 0 — Fundação
- Scaffold do projeto: `wails init` com template Vue+TS, estrutura de pastas da seção 3, `internal/model` com todas as structs deste documento (mesmo as ainda não usadas).
- Configurar Pinia, Reka UI, Tiptap (instalação, sem extensões ainda).
- Verificação de ambiente: checagem de `tectonic` no PATH ao iniciar.
- Rodar `/graphify .` pela primeira vez.
- **Aceite:** `wails dev` abre uma janela com layout vazio das 5 zonas (Activity Bar, Sidebar, Editor placeholder, Status Bar, Problems Panel recolhido); `rtk go test ./...` passa (mesmo que com poucos testes).

### Etapa 1 — Workspace e Capítulos (sem editor real)
- Implementar `internal/workspace`: criar/abrir workspace, leitura/escrita de `project.json` e capítulos com validação de schema.
- Métodos da ponte: todo o bloco "Workspace" e "Capítulos" da seção 7.1 (incluindo o parâmetro `role` em `CreateChapter`).
- UI: diálogo de criar/abrir projeto; "Novo documento" com seletor de tipo (roles da seção 5.5); Explorer na Sidebar com os três grupos ("Pré-textuais", "Capítulos", "Pós-textuais"), com criar, renomear (título), excluir e reordenar por arrastar dentro do grupo.
- **Aceite:** criar um projeto novo gera exatamente a estrutura da seção 4; criar uma dedicatória, dois capítulos e um "Sobre o Autor" os posiciona nos grupos corretos e reflete em disco; reabrir o app restaura o estado; testes Go cobrem IO do workspace (incluindo JSON corrompido → erro amigável).

### Etapa 2 — Editor Tiptap com Nós Semânticos
- Extensões Tiptap: `luzChapter`, `luzSection`, `luzSubsection`, `luzFootnote`, `luzImage`, `luzSoftHyphen`, `luzVariable`, marks `luzInlineQuote`, `luzLang`, `bold`, `italic`, `underline` e `strike`, atributo `align` do parágrafo (botões de alinhamento na toolbar), mais os nós padrão da seção 6. (Nesta etapa os botões de idioma e hífen ficam sempre visíveis; o gating pelos plugins `languages`/`hyphenation` entra na Etapa 4. A mark `luzCustomStyle` fica para a Etapa 4, junto do seu plugin.)
- **Variáveis do Projeto** (seção 5.1): painel de gestão em *Settings* (adicionar/renomear/editar/excluir) + dropdown "Inserir variável" na toolbar + chips reativos no editor (alterar o valor atualiza todas as ocorrências).
- Popover de atributos (`numbered`, `includeInToc`) nos títulos; popover de edição de nota de rodapé; renumeração automática das notas; popover de `caption`/`width` na imagem, com importação via `ImportImage()`.
- Toolbar sensível ao `role` do documento aberto (seção 5.5): páginas especiais escondem títulos estruturais; dedicatória/epígrafe escondem também listas e imagens.
- Salvar/carregar capítulo via ponte, com autosave (debounce ~800ms) e indicador "salvo/salvando" na Status Bar; contagem de palavras.
- Toolbar mínima sensível ao `kind` do target ativo (esconder "Capítulo" em `article` — o target pode estar hardcoded como o default até a Etapa 4).
- Visuais de intenção conforme a seção 6.0 (CSS canônico da `luzInlineQuote`, rótulos discretos nos títulos, marcador `[n]` das notas).
- **Aceite:** o JSON salvo de um capítulo de teste corresponde estruturalmente ao exemplo da seção 5.4; a citação inline exibe fundo suave e aspas via CSS sem que nenhum caractere de aspas exista no conteúdo salvo; imagem importada aparece no editor e é copiada para `imagens/`; fechar e reabrir preserva tudo, incluindo notas de rodapé, imagens e atributos dos títulos.

### Etapa 3 — Conversor LaTeX + Primeira Compilação
- `internal/latex`: conversão completa da tabela da seção 6 (incluindo `luzImage`, `luzLang` → `\foreignlanguage`, `language` de documento → `otherlanguage`, `luzSoftHyphen` → `\-`, expansão de `luzVariable`), escape de caracteres e os envoltórios de `role` da seção 5.5 (frontmatter/mainmatter/backmatter), com testes unitários (o exemplo canônico da seção 6 é um teste obrigatório). Detecção de idiomas usados para o preâmbulo (seção 10, passo 2).
- `internal/build`: pipeline da seção 10 usando um **target fixo embutido** (book, apenas módulos do núcleo com seus defaults: geometry A5 espelhado + bodyText padrão) — o sistema de targets vem na próxima etapa. Preâmbulo base completo da seção 8.3 (microtype, hyperref etc.).
- Botão "Exportar PDF" + eventos de progresso + abertura da pasta `dist/` ao concluir.
- **Aceite:** um projeto com dedicatória, 2 capítulos (com negrito, itálico, sublinhado, um parágrafo centralizado, citação inline, um trecho marcado em inglês, um hífen sugerido, uma variável do projeto usada duas vezes, blockquote, lista, nota de rodapé e uma imagem com legenda) e "Sobre o Autor" compila num PDF correto em `dist/` (com o valor da variável expandido), com a dedicatória antes e o "Sobre o Autor" depois dos capítulos; falha do Tectonic exibe `LogTail` sem travar o app; `rtk go test ./...` verde.

### Etapa 4 — Targets e Plugins (Schema-Driven UI)
- `internal/plugins`: interface da seção 8.1 (incluindo `Core()`), catálogo com os módulos do núcleo (`geometry`, `bodyText` — abas fixas em *Build Targets*, seção "Essenciais" com switch travado em *Extensions*) e os plugins opcionais (`fancyhdr`, `catalogRecord` com o formulário completo da seção 8.3 e página no verso da página de rosto, `pdfpages` com importação via `ImportAttachmentPDF()`, `languages`, `hyphenation` e `customStyles` com gerenciador de estilos, mark `luzCustomStyle` e geração de `\newcommand` — seções 5.7 e 8.3); geração de preâmbulo por Builder. Gating de toolbar: botões de idioma, hífen sugerido e estilos (e o seletor de idioma do documento) só aparecem com os respectivos plugins ativos.
- **Resolvedor de camadas** (`internal/plugins/resolve.go`) implementando a precedência da seção 8.4, com testes unitários do merge (default → target → override).
- CRUD de targets (métodos da seção 7.1) + **presets da seção 8.5** no diálogo "Novo target" + seletor de target ativo na Status Bar.
- Aba *Extensions* (liga/desliga plugins) e `SchemaForm.vue` renderizando os `FormSchema` nas abas de *Build Targets*, persistindo em `pluginConfig`.
- **Painel "Configurações desta página"** (seções 5.6 e 8.4): `SchemaForm.vue` em modo override, com toggle por campo e valores herdados esmaecidos; persistência em `.luz/overrides/<id>.json` via `Get/SaveDocumentOverrides` (criação sob demanda, `{}` apaga); badge ⚙ na árvore do Explorer + item no menu de contexto + botão no editor; cascata de renomear/excluir documento sobre o arquivo de override; `ScopedLaTeX()` implementado em `geometry` e `fancyhdr` e aplicado na compilação.
- Compilação passa a usar o target ativo real.
- **Aceite:** criar um target a partir do preset "Amazon KDP 6×9" e outro do preset "e-Book" e compilar cada um gera PDFs com geometrias visivelmente diferentes; preencher o formulário da ficha catalográfica gera a página no verso da página de rosto; importar um PDF via `pdfpages` o insere na posição configurada; definir margens maiores só na dedicatória (via painel "Configurações desta página") cria `.luz/overrides/00-dedicatoria.json`, exibe o badge ⚙ na árvore e altera apenas essa página no PDF, com os capítulos herdando o target; "Limpar sobrescritas" apaga o arquivo e o badge; alterar um campo no formulário persiste no JSON correspondente e afeta o PDF seguinte.

### Etapa 5 — Rule Engine e Problems Panel
- `internal/rules` com as 18 regras da seção 9, cada uma com teste.
- Revalidação automática ao salvar projeto/target/capítulo/plugins; evento `luz:problems`; Problems Panel funcional com navegação; contagem na Status Bar; erros bloqueando `Compile()`.
- **Aceite:** habilitar margens espelhadas no target ebook exibe R001 e impede a exportação com mensagem clara; corrigir o problema limpa o painel e libera a compilação.

### Etapa 6 — Polimento do MVP
- Sumário dinâmico: `\tableofcontents` opcional por target (campo no target: `includeToc: boolean`), respeitando `includeInToc` dos títulos.
- Página de título gerada a partir de `project.json` (`\maketitle` customizado).
- Atalhos de teclado (Ctrl/Cmd+S salvar, Ctrl/Cmd+E exportar), estados de erro amigáveis, revisão de todo texto de UI em pt-BR.
- Revisão final: `/graphify .` completo + leitura do `GRAPH_REPORT.md` procurando inconsistências arquiteturais (ex.: LaTeX vazando para o frontend).
- **Aceite:** fluxo completo de ponta a ponta — criar projeto, escrever 3 capítulos com todos os recursos, configurar 2 targets, resolver um problema apontado pelo Rule Engine e exportar ambos os PDFs — sem tocar em nenhum arquivo manualmente.

### Fora de escopo do MVP (não implementar sem pedido explícito)
Exportação EPUB, sistema de plugins de terceiros (o catálogo é embutido), roteiros/screenplays, colaboração, corretor ortográfico, temas de UI, referências bibliográficas (BibTeX), edição de imagem (recorte/filtros — apenas inserção com legenda e largura), índice remissivo, glossário. **Comandos LaTeX crus escritos pelo usuário** (`\newcommand` manual, preâmbulo customizado, blocos de LaTeX inline): explicitamente excluídos — LaTeX arbitrário do usuário anularia as garantias do Rule Engine e a promessa de erros sempre explicáveis no Problems Panel. As necessidades cobertas por `\newcommand` são atendidas por Variáveis do Projeto (seção 5.1) e Estilos Personalizados (seção 5.7); um eventual "Modo Perito" opt-in, com aviso claro de perda de garantias, é candidato a plugin pós-MVP.

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
