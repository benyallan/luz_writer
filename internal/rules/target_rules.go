package rules

import (
	"encoding/json"
	"fmt"

	"luz-writer/internal/latex"
	"luz-writer/internal/model"
	"luz-writer/internal/plugins"
)

// R001: geometry.mirrored=true em target kind "ebook".
func R001(ctx Context) []model.Problem {
	if ctx.Target.Kind != model.TargetKindEbook {
		return nil
	}
	cfg, err := effectiveConfig(ctx, plugins.Geometry)
	if err != nil {
		return nil
	}
	var g plugins.GeometryConfig
	if err := json.Unmarshal(cfg, &g); err != nil || !g.Mirrored {
		return nil
	}
	return []model.Problem{{
		Severity: "error",
		Code:     "R001",
		Message:  "Margens espelhadas não fazem sentido em e-books (fluxo contínuo).",
		Source:   "target:" + ctx.Target.ID,
	}}
}

// R002: nó luzChapter presente com target documentClass "article".
func R002(ctx Context) []model.Problem {
	if ctx.Target.DocumentClass != model.DocumentClassArticle {
		return nil
	}
	var problems []model.Problem
	for _, ch := range ctx.Chapters {
		usage, err := latex.ScanContent(ch.Content)
		if err != nil || !usage.HasLuzChapter {
			continue
		}
		problems = append(problems, model.Problem{
			Severity: "error",
			Code:     "R002",
			Message:  "Artigos não possuem capítulos. Converta para Seção ou mude o target.",
			Source:   "chapter:" + ch.ID,
		})
	}
	return problems
}

// R003: plugin catalogRecord habilitado em target kind "ebook".
func R003(ctx Context) []model.Problem {
	if ctx.Target.Kind != model.TargetKindEbook || !ctx.isEnabled("catalogRecord") {
		return nil
	}
	return []model.Problem{{
		Severity: "warning",
		Code:     "R003",
		Message:  "Ficha Catalográfica ativada em um target de e-Book.",
		Source:   "plugin:catalogRecord",
	}}
}

// R004: target kind "print" com geometry.mirrored=false.
func R004(ctx Context) []model.Problem {
	if ctx.Target.Kind != model.TargetKindPrint {
		return nil
	}
	cfg, err := effectiveConfig(ctx, plugins.Geometry)
	if err != nil {
		return nil
	}
	var g plugins.GeometryConfig
	if err := json.Unmarshal(cfg, &g); err != nil || g.Mirrored {
		return nil
	}
	return []model.Problem{{
		Severity: "warning",
		Code:     "R004",
		Message:  "Livros impressos frente-e-verso geralmente usam margens espelhadas — confirme se a impressão será apenas frente.",
		Source:   "target:" + ctx.Target.ID,
	}}
}

// R010: documento com role especial em target documentClass "article".
func R010(ctx Context) []model.Problem {
	if ctx.Target.DocumentClass != model.DocumentClassArticle {
		return nil
	}
	var problems []model.Problem
	for _, ch := range ctx.Chapters {
		if ch.Role == model.RoleChapter || ch.Role == model.RoleAppendix || ch.Role == "" {
			continue
		}
		problems = append(problems, model.Problem{
			Severity: "warning",
			Code:     "R010",
			Message:  fmt.Sprintf("Artigos normalmente não possuem %s; o conteúdo será inserido como seção simples.", roleLabel(ch.Role)),
			Source:   "chapter:" + ch.ID,
		})
	}
	return problems
}

// R011: catalogRecord e pdfpages (com placement pré-textual) ativos juntos.
func R011(ctx Context) []model.Problem {
	if !ctx.isEnabled("catalogRecord") || !ctx.isEnabled("pdfpages") {
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
	if c.Placement != "beforeFrontmatter" && c.Placement != "afterFrontmatter" {
		return nil
	}
	return []model.Problem{{
		Severity: "warning",
		Code:     "R011",
		Message:  "Você tem uma ficha catalográfica gerada e um PDF anexo pré-textual ao mesmo tempo — verifique se não haverá duplicidade.",
		Source:   "plugin:catalogRecord",
	}}
}
