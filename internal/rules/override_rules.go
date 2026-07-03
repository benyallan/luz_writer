package rules

import (
	"fmt"

	"luz-writer/internal/model"
	"luz-writer/internal/plugins"
)

// R012: documento com override para plugin desabilitado em plugins.json.
func R012(ctx Context) []model.Problem {
	var problems []model.Problem
	for _, ch := range ctx.Chapters {
		overrides := documentOverrides(ctx, ch.ID)
		for name := range overrides {
			p, ok := plugins.ByName(name)
			if !ok || p.Core() || ctx.isEnabled(name) {
				continue
			}
			problems = append(problems, model.Problem{
				Severity: "warning",
				Code:     "R012",
				Message:  fmt.Sprintf("O documento '%s' sobrescreve configurações de '%s', que está desativado — as sobrescritas serão ignoradas.", ch.ID, name),
				Source:   "override:" + ch.ID,
			})
		}
	}
	return problems
}

// R013: documento com override de geometry em target kind "ebook".
func R013(ctx Context) []model.Problem {
	if ctx.Target.Kind != model.TargetKindEbook {
		return nil
	}
	var problems []model.Problem
	for _, ch := range ctx.Chapters {
		overrides := documentOverrides(ctx, ch.ID)
		if _, has := overrides["geometry"]; !has {
			continue
		}
		problems = append(problems, model.Problem{
			Severity: "warning",
			Code:     "R013",
			Message:  "Sobrescritas de geometria por página têm pouco efeito em e-books de fluxo contínuo.",
			Source:   "override:" + ch.ID,
		})
	}
	return problems
}

// R014: arquivo em .luz/overrides/ cujo documentId não existe mais.
func R014(ctx Context) []model.Problem {
	orphans, err := ctx.Workspace.OrphanedOverrideIDs(ctx.Project.ChapterOrder)
	if err != nil {
		return nil
	}
	var problems []model.Problem
	for _, id := range orphans {
		problems = append(problems, model.Problem{
			Severity: "warning",
			Code:     "R014",
			Message:  fmt.Sprintf("Há configurações de página para o documento '%s', que não existe mais no projeto.", id),
			Source:   "override:" + id,
		})
	}
	return problems
}
