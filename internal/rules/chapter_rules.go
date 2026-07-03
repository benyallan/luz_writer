package rules

import (
	"fmt"

	"luz-writer/internal/latex"
	"luz-writer/internal/model"
)

// R005: capítulo referenciado em chapterOrder com conteúdo vazio.
func R005(ctx Context) []model.Problem {
	var problems []model.Problem
	for _, ch := range ctx.Chapters {
		usage, err := latex.ScanContent(ch.Content)
		if err != nil || !usage.IsBlank() {
			continue
		}
		problems = append(problems, model.Problem{
			Severity: "warning",
			Code:     "R005",
			Message:  fmt.Sprintf("O capítulo '%s' está vazio.", ch.ID),
			Source:   "chapter:" + ch.ID,
		})
	}
	return problems
}

// R006: chapterOrder referencia arquivo inexistente em capitulos/.
func R006(ctx Context) []model.Problem {
	var problems []model.Problem
	for _, id := range ctx.MissingChapterID {
		problems = append(problems, model.Problem{
			Severity: "error",
			Code:     "R006",
			Message:  fmt.Sprintf("Capítulo '%s' listado no projeto mas o arquivo não existe.", id),
			Source:   "project",
		})
	}
	return problems
}

// R007: documento de frontmatter posicionado depois de um chapter.
func R007(ctx Context) []model.Problem {
	var problems []model.Problem
	seenChapter := false
	for _, ch := range ctx.Chapters {
		switch ch.Role {
		case model.RoleChapter:
			seenChapter = true
		case model.RoleDedication, model.RoleEpigraph, model.RoleAcknowledgments, model.RolePreface:
			if seenChapter {
				problems = append(problems, model.Problem{
					Severity: "warning",
					Code:     "R007",
					Message:  fmt.Sprintf("A '%s' será reposicionada para o início do livro na compilação.", roleLabel(ch.Role)),
					Source:   "chapter:" + ch.ID,
				})
			}
		}
	}
	return problems
}
