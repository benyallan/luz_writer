# Luz Writer — Visão Geral do Projeto

Luz Writer é uma IDE desktop para escritores (livros, artigos, e-books, roteiros). Construída com **Wails v2** (Go + Vue 3), oferece a experiência de um editor de texto rico (como LibreOffice Writer) mas com o poder do LaTeX por baixo.

## Princípio central

O escritor nunca toca LaTeX. Ele escreve em uma interface WYSIWYG rica. O backend Go converte o documento para LaTeX e usa essa representação para gerar PDFs com qualidade tipográfica profissional, sumários automáticos, configurações de impressão (livro A5, artigo A4, e-book reflowable, etc.) e tudo que o ecossistema LaTeX oferece.

## Stack

### Frontend (`frontend/`)
- **Vue 3** + Composition API
- **Tiptap** — editor de texto rico (extensível via extensões próprias do Tiptap)
- **Reka UI** — componentes acessíveis sem estilo forçado (headless)
- **Vite** como bundler

### Backend (`*.go`)
- **Go** + **Wails v2** para bridge entre frontend e backend
- Conversão de HTML/JSON (ProseMirror) → LaTeX
- Orquestração de compilação LaTeX (via `tectonic` ou `pdflatex` instalado no sistema)
- Sistema de plugins via interface Go

## Estrutura do repositório (código-fonte do app)

```
luz-writer/                  # Raiz do repositório do app
├── main.go                  # Entrada Wails
├── app.go                   # App struct e todos os métodos expostos ao frontend
├── internal/
│   ├── document/            # Modelo de documento, serialização
│   ├── latex/               # Conversão ProseMirror JSON → LaTeX
│   ├── compiler/            # Invocação do compilador LaTeX
│   ├── plugin/              # Runtime de plugins
│   └── config/              # Configurações de usuário/app
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── layout/      # Shell do app (MenuBar, Sidebar)
│   │   │   ├── sidebar/     # FileTree, FileTreeNode
│   │   │   └── dialogs/     # Modais (NewProjectDialog, ...)
│   │   ├── stores/          # Pinia stores (workspace.js, ...)
│   │   ├── editor/          # Extensões Tiptap customizadas (a criar)
│   │   └── panels/          # Painéis laterais/inferiores (a criar)
│   └── wailsjs/             # Gerado pelo Wails — NÃO editar à mão exceto App.js/App.d.ts/models.ts
├── build/                   # Assets de build (ícones, manifests)
└── .claude/
    └── commands/            # Skills deste projeto
```

## Estrutura de um projeto do escritor (criado pelo app)

Quando o usuário cria um novo projeto pelo Luz Writer, a seguinte estrutura é
gerada em disco. Toda a lógica de scaffold está em `app.go → CreateProject`.

```
meu-livro/                   # Pasta raiz escolhida pelo usuário
├── src/                     # CONTEÚDO — arquivos que o escritor edita
│   │                        # Capítulos, seções, artigos em formato próprio
│   │                        # do Luz Writer (JSON ProseMirror salvo em .lwx ou similar)
│   └── (capítulos, partes)
│
├── targets/                 # PERFIS DE COMPILAÇÃO
│   │                        # Cada arquivo descreve um destino de saída:
│   │                        # classe LaTeX, margens, fonte, tamanho de página,
│   │                        # pacotes extras, metadados (título, autor, ISBN…)
│   │                        # Exemplos: livro-a5.toml, ebook.toml, artigo-abnt.toml
│   └── (perfis .toml ou similar)
│
├── dist/                    # SAÍDA — arquivos gerados (não versionar)
│   │                        # PDFs compilados, EPUBs, HTML exportados
│   │                        # Sobrescrito a cada compilação; o usuário não edita aqui
│   └── (PDFs, EPUBs…)
│
├── .tmp/                    # TEMPORÁRIOS de compilação (não versionar)
│   │                        # Arquivos .tex gerados, .aux, .log, fontes baixadas
│   │                        # pelo Tectonic, cache de compilação incremental
│   │                        # Pode ser apagado a qualquer momento sem perda
│   └── (artefatos LaTeX)
│
└── .gitignore               # Gerado automaticamente; exclui .tmp/ e dist/
```

### Regras sobre o projeto do escritor

- **Nunca ler nem escrever em `dist/` diretamente** — é território exclusivo do compilador.
- **`.tmp/` é descartável** — o app deve recriá-la se não existir antes de compilar.
- **`targets/`** define *o quê* compilar e *como*; `src/` define *o conteúdo*. São orthogonais.
- O formato dos arquivos de `src/` ainda não está definido — candidatos: `.lwx` (JSON ProseMirror), `.md` com front-matter, ou formato binário próprio.
- O formato dos perfis em `targets/` ainda não está definido — candidato: TOML por ser legível e simples.

## Modos de exportação LaTeX

| Modo       | Classe LaTeX | Caso de uso                     |
|------------|--------------|---------------------------------|
| `book`     | `book`       | Livro A5/A6, capítulos, partes  |
| `article`  | `article`    | Artigos acadêmicos, A4          |
| `memoir`   | `memoir`     | Livros com controle fino        |
| `beamer`   | `beamer`     | Apresentações                   |
| `ebook`    | custom       | E-pub via pandoc ou htlatex     |

## Sistema de plugins

Plugins são Go structs que implementam a interface `Plugin`:

```go
type Plugin interface {
    ID()          string
    Name()        string
    Version()     string
    Register(app *App) error
    Shutdown() error
}
```

Plugins podem:
- Registrar novos comandos Wails (métodos Go expostos ao frontend)
- Injetar extensões Tiptap (via manifesto JSON)
- Adicionar classes/pacotes LaTeX ao pipeline de compilação
- Adicionar painéis laterais (componentes Vue carregados dinamicamente)

### Plugins planejados
- `lang-foreign` — suporte a línguas estrangeiras (hifenização, fontes)
- `research-library` — importar e referenciar livros/artigos (BibTeX)
- `word-count` — metas de escrita e estatísticas
- `distraction-free` — modo foco com tema escuro e sem painéis

## Convenções de desenvolvimento

- Go: pacotes em `internal/`, nunca exportar tipos que não precisem ser públicos
- Vue: Composition API com `<script setup>`, sem Options API
- Nomes de arquivos Vue: PascalCase para componentes, kebab-case para páginas
- Stores Pinia: um arquivo por domínio (`document.ts`, `editor.ts`, `settings.ts`)
- Bindings Wails: sempre regenerar com `wails generate module` após mudar `app.go`
- Nunca expor detalhes LaTeX na interface Go pública — é detalhe de implementação

## Fluxo de compilação LaTeX

1. Frontend envia JSON do documento (formato ProseMirror) via Wails binding
2. `internal/latex` converte para string `.tex` com preâmbulo correto
3. `internal/compiler` grava o `.tex` em diretório temporário e invoca o compilador
4. O PDF resultante é retornado como caminho ou bytes para o frontend exibir

## Referências úteis

- Wails v2 docs: https://wails.io/docs/introduction
- Tiptap docs: https://tiptap.dev/docs
- Reka UI: https://reka-ui.com
- Tectonic (compilador LaTeX moderno): https://tectonic-typesetting.github.io
