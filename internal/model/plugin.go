package model

// PluginManifest é o resumo de um plugin exposto via ListAvailablePlugins
// (seção 7.1), usado pela aba Extensions para listar módulos do núcleo e
// plugins opcionais (seção 8.0). Enabled reflete o plugins.json do workspace
// ativo (sempre true para módulos do núcleo).
type PluginManifest struct {
	Name          string     `json:"name"`
	DisplayName   string     `json:"displayName"`
	Description   string     `json:"description"`
	Core          bool       `json:"core"`          // true = módulo do núcleo, sempre ativo (seção 8.0)
	DocumentScope bool       `json:"documentScope"` // se aceita override por documento (seção 5.6/8.4)
	Schema        FormSchema `json:"schema"`
	Enabled       bool       `json:"enabled"`
}

// BuildContext carrega o contexto usado por Validate/Preamble/ScopedLaTeX de
// cada plugin (seção 8.1): projeto, target ativo e módulos ativos (núcleo +
// opcionais habilitados), para permitir regras cruzadas entre plugins.
// Styles carrega .luz/styles.json — necessário pelo plugin customStyles, que
// não tem configuração própria por target (seção 8.3).
type BuildContext struct {
	Project       Project       `json:"project"`
	Target        Target        `json:"target"`
	ActiveModules []string      `json:"activeModules"`
	Styles        []CustomStyle `json:"styles"`
}
