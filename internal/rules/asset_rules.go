package rules

import (
	"encoding/json"
	"fmt"

	"luz-writer/internal/latex"
	"luz-writer/internal/model"
	"luz-writer/internal/plugins"
)

// R008: nó luzImage cujo src não existe em imagens/.
func R008(ctx Context) []model.Problem {
	var problems []model.Problem
	for _, ch := range ctx.Chapters {
		usage, err := latex.ScanContent(ch.Content)
		if err != nil {
			continue
		}
		for _, src := range usage.ImageSrcs {
			if src == "" || ctx.Workspace.FileExistsRelative(src) {
				continue
			}
			problems = append(problems, model.Problem{
				Severity: "error",
				Code:     "R008",
				Message:  fmt.Sprintf("Imagem '%s' não encontrada no capítulo '%s'.", src, ch.ID),
				Source:   "chapter:" + ch.ID,
			})
		}
	}
	return problems
}

// R009: plugin pdfpages habilitado com filePath vazio ou inexistente.
func R009(ctx Context) []model.Problem {
	if !ctx.isEnabled("pdfpages") {
		return nil
	}
	cfg, err := effectiveConfig(ctx, plugins.Pdfpages)
	if err != nil {
		return nil
	}
	var c plugins.PdfpagesConfig
	if err := json.Unmarshal(cfg, &c); err != nil {
		return nil
	}
	if c.FilePath != "" && ctx.Workspace.FileExistsRelative(c.FilePath) {
		return nil
	}
	return []model.Problem{{
		Severity: "error",
		Code:     "R009",
		Message:  "O anexo PDF configurado não foi encontrado.",
		Source:   "plugin:pdfpages",
	}}
}
