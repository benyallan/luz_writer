package plugins

import (
	"encoding/json"

	"luz-writer/internal/model"
)

// Plugin é a interface implementada por todo módulo do núcleo ou plugin
// opcional (seção 8.1).
type Plugin interface {
	Name() string        // "geometry"
	DisplayName() string // "Geometria de Página" (pt-BR)
	Description() string
	Core() bool // true = módulo do núcleo (seção 8.0): sempre ativo, não desativável
	Schema() model.FormSchema
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
	ScopedLaTeX(effectiveCfg json.RawMessage, ctx model.BuildContext) (before string, after string, err error)
}
