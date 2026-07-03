package rules

import (
	"encoding/json"
	"fmt"

	"luz-writer/internal/latex"
	"luz-writer/internal/model"
	"luz-writer/internal/plugins"
)

// R015: mark luzLang ou language de documento, com plugin languages desativado.
func R015(ctx Context) []model.Problem {
	if ctx.isEnabled("languages") {
		return nil
	}
	var problems []model.Problem
	for _, ch := range ctx.Chapters {
		usage, err := latex.ScanContent(ch.Content)
		if err != nil {
			continue
		}
		hasDocLang := ch.Language != nil && *ch.Language != "" && *ch.Language != ctx.Project.Language
		if len(usage.LangCodes) == 0 && !hasDocLang {
			continue
		}
		problems = append(problems, model.Problem{
			Severity: "warning",
			Code:     "R015",
			Message:  "Há conteúdo marcado em outro idioma, mas o Suporte Multilíngue está desativado — as marcações serão ignoradas na exportação.",
			Source:   "chapter:" + ch.ID,
		})
	}
	return problems
}

// R016: nó luzSoftHyphen presente com plugin hyphenation desativado/off.
func R016(ctx Context) []model.Problem {
	var problems []model.Problem
	for _, ch := range ctx.Chapters {
		usage, err := latex.ScanContent(ch.Content)
		if err != nil || !usage.HasSoftHyphen {
			continue
		}
		if hyphenationActiveFor(ctx, ch.ID) {
			continue
		}
		problems = append(problems, model.Problem{
			Severity: "warning",
			Code:     "R016",
			Message:  "Há pontos de hifenização sugeridos, mas a hifenização está desativada — eles serão ignorados na exportação.",
			Source:   "chapter:" + ch.ID,
		})
	}
	return problems
}

func hyphenationActiveFor(ctx Context, chapterID string) bool {
	if !ctx.isEnabled("hyphenation") {
		return false
	}
	overrides := documentOverrides(ctx, chapterID)
	cfg, err := effectiveConfigWithOverride(ctx, plugins.Hyphenation, overrides["hyphenation"])
	if err != nil {
		return true // config inválida não é o que R016 quer sinalizar
	}
	var c plugins.HyphenationConfig
	if err := json.Unmarshal(cfg, &c); err != nil {
		return true
	}
	return c.Mode != "off"
}

// R017: nó luzVariable cujo name não existe em project.json → variables.
func R017(ctx Context) []model.Problem {
	names := make(map[string]bool, len(ctx.Project.Variables))
	for _, v := range ctx.Project.Variables {
		names[v.Name] = true
	}

	var problems []model.Problem
	for _, ch := range ctx.Chapters {
		usage, err := latex.ScanContent(ch.Content)
		if err != nil {
			continue
		}
		for _, name := range usage.VariableNames {
			if names[name] {
				continue
			}
			problems = append(problems, model.Problem{
				Severity: "error",
				Code:     "R017",
				Message:  fmt.Sprintf("A variável '%s' usada no documento '%s' não existe mais no projeto.", name, ch.ID),
				Source:   "chapter:" + ch.ID,
			})
		}
	}
	return problems
}

// R018: mark luzCustomStyle cujo styleId não existe, ou plugin desativado.
func R018(ctx Context) []model.Problem {
	styleIDs := make(map[string]bool, len(ctx.Styles))
	for _, s := range ctx.Styles {
		styleIDs[s.ID] = true
	}
	enabled := ctx.isEnabled("customStyles")

	var problems []model.Problem
	for _, ch := range ctx.Chapters {
		usage, err := latex.ScanContent(ch.Content)
		if err != nil {
			continue
		}
		for _, id := range usage.CustomStyleIDs {
			if enabled && styleIDs[id] {
				continue
			}
			problems = append(problems, model.Problem{
				Severity: "warning",
				Code:     "R018",
				Message:  fmt.Sprintf("O estilo '%s' não está disponível — o trecho será exportado como texto sem formatação de estilo.", id),
				Source:   "chapter:" + ch.ID,
			})
		}
	}
	return problems
}
