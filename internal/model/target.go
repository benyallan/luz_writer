package model

import "encoding/json"

// TargetKind é o tipo de saída do target (seção 5.2).
type TargetKind string

const (
	TargetKindPrint   TargetKind = "print"
	TargetKindEbook   TargetKind = "ebook"
	TargetKindArticle TargetKind = "article"
)

// DocumentClass é a classe LaTeX associada ao target (seção 5.2).
type DocumentClass string

const (
	DocumentClassBook    DocumentClass = "book"
	DocumentClassReport  DocumentClass = "report"
	DocumentClassArticle DocumentClass = "article"
)

// Target é o conteúdo de .luz/targets/<id>.json (seção 5.2).
type Target struct {
	LuzVersion    int                        `json:"luzVersion"`
	ID            string                     `json:"id"`
	Name          string                     `json:"name"`
	Kind          TargetKind                 `json:"kind"`
	DocumentClass DocumentClass              `json:"documentClass"`
	FontSize      string                     `json:"fontSize"`
	IncludeToc    bool                       `json:"includeToc"`
	PluginConfig  map[string]json.RawMessage `json:"pluginConfig"`
}

// PluginsConfig é o conteúdo de .luz/plugins.json (seção 5.3).
// Enabled lista apenas plugins OPCIONAIS; módulos do núcleo nunca aparecem aqui.
type PluginsConfig struct {
	LuzVersion int      `json:"luzVersion"`
	Enabled    []string `json:"enabled"`
}
