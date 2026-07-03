package rules

import (
	"luz-writer/internal/model"
	"luz-writer/internal/workspace"
)

// Context carrega tudo que as regras precisam para validar um projeto —
// montado uma única vez por Validate()/Compile() e passado a cada regra.
type Context struct {
	Workspace *workspace.Workspace
	Project   model.Project
	Target    model.Target
	// EnabledOptional lista os plugins opcionais habilitados, na ordem
	// declarada em plugins.json (módulos do núcleo nunca aparecem aqui).
	EnabledOptional []string
	Styles          []model.CustomStyle
	// Chapters são os documentos que existem e puderam ser lidos, na ordem
	// de chapterOrder. MissingChapterIDs são os ids listados em chapterOrder
	// cujo arquivo não existe (R006) — Chapters não os contém.
	Chapters         []model.Chapter
	MissingChapterID []string
}

// NewContext monta o Context a partir do workspace ativo.
func NewContext(ws *workspace.Workspace) (Context, error) {
	project, err := ws.Project()
	if err != nil {
		return Context{}, err
	}
	target, err := ws.ActiveTarget()
	if err != nil {
		return Context{}, err
	}
	pluginsCfg, err := ws.PluginsConfig()
	if err != nil {
		return Context{}, err
	}
	styles, err := ws.ListStyles()
	if err != nil {
		return Context{}, err
	}
	chapters, missing, err := ws.ChaptersTolerant()
	if err != nil {
		return Context{}, err
	}

	return Context{
		Workspace:        ws,
		Project:          project,
		Target:           target,
		EnabledOptional:  pluginsCfg.Enabled,
		Styles:           styles,
		Chapters:         chapters,
		MissingChapterID: missing,
	}, nil
}

func (ctx Context) isEnabled(name string) bool {
	for _, n := range ctx.EnabledOptional {
		if n == name {
			return true
		}
	}
	return false
}
