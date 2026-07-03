package model

import "encoding/json"

// DocumentRole determina em qual bloco do livro o documento é compilado
// e qual LaTeX o envolve (seção 5.5). Escolhido na criação; nunca "frontmatter"
// diretamente pelo escritor.
type DocumentRole string

const (
	RoleChapter         DocumentRole = "chapter"
	RoleDedication      DocumentRole = "dedication"
	RoleEpigraph        DocumentRole = "epigraph"
	RoleAcknowledgments DocumentRole = "acknowledgments"
	RolePreface         DocumentRole = "preface"
	RoleAboutAuthor     DocumentRole = "aboutAuthor"
	RoleAppendix        DocumentRole = "appendix"
)

// Chapter é o conteúdo de capitulos/<id>.json (seção 5.4).
// Contém apenas texto — configurações de página vivem em .luz/overrides/.
type Chapter struct {
	LuzVersion int             `json:"luzVersion"`
	ID         string          `json:"id"`
	Role       DocumentRole    `json:"role"`
	Language   *string         `json:"language"` // nil = herda o language do project.json
	Content    json.RawMessage `json:"content"`  // documento ProseMirror
}

// ChapterMeta é o resumo usado para popular o Explorer (ListChapters, seção 7.1).
// HasOverrides alimenta o badge ⚙ da árvore (seção 5.6/8.4) sem exigir uma
// chamada por documento a GetDocumentOverrides.
type ChapterMeta struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	Role         DocumentRole `json:"role"`
	WordCount    int          `json:"wordCount"`
	HasOverrides bool         `json:"hasOverrides"`
}

// DocumentOverrides é o conteúdo de .luz/overrides/<id>.json (seção 5.6).
// Existe apenas quando o documento tem ao menos uma sobrescrita; salvar
// Overrides vazio apaga o arquivo (regra de internal/workspace).
type DocumentOverrides struct {
	LuzVersion int                        `json:"luzVersion"`
	DocumentID string                     `json:"documentId"`
	Overrides  map[string]json.RawMessage `json:"overrides"`
}
