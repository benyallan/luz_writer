package latex

import (
	"fmt"
	"strings"

	"luz-writer/internal/model"
)

// RenderDocument converte um documento inteiro (capitulos/<id>.json) para
// LaTeX, aplicando o envoltório do seu role (seção 5.5) e, se houver, o
// idioma de documento (seção 5.4). Não inclui os marcadores de bloco
// \frontmatter/\mainmatter/\backmatter/\appendix — esses são inseridos na
// montagem de múltiplos documentos, em internal/build.
func RenderDocument(ch model.Chapter, vars map[string]string, opts Options) (tex string, problems []model.Problem, langsUsed map[string]bool) {
	doc, err := parseDoc(ch.Content)
	if err != nil {
		return "", []model.Problem{{
			Severity: "error",
			Code:     "LATEX",
			Message:  fmt.Sprintf("Não foi possível interpretar o conteúdo do documento '%s'.", ch.ID),
			Source:   "chapter:" + ch.ID,
		}}, nil
	}

	ctx := newContext(ch.ID, vars, opts)
	body := renderByRole(ch.Role, doc, ctx)

	// R015: com o plugin languages desativado, o idioma de documento também
	// é ignorado na exportação.
	if opts.LanguagesEnabled && ch.Language != nil && *ch.Language != "" {
		name, ok := BabelLangName[*ch.Language]
		if !ok {
			ctx.warnf("Idioma de documento desconhecido '%s'.", *ch.Language)
		} else {
			ctx.LangsUsed[*ch.Language] = true
			body = "\\begin{otherlanguage}{" + name + "}\n" + body + "\n\\end{otherlanguage}"
		}
	}

	return body, ctx.Problems, ctx.LangsUsed
}

func renderByRole(role model.DocumentRole, doc node, ctx *Context) string {
	switch role {
	case model.RoleChapter, model.RoleAppendix, "":
		// O título vem do próprio nó luzChapter dentro do conteúdo — sem
		// envoltório adicional. A numeração especial de apêndice
		// ("Apêndice A, B...") vem do \appendix inserido na montagem.
		return convertBlocks(doc.Content, ctx)
	case model.RoleAcknowledgments:
		return "\\chapter*{Agradecimentos}\n\n" + convertBlocks(doc.Content, ctx)
	case model.RolePreface:
		return "\\chapter*{Prefácio}\n\\addcontentsline{toc}{chapter}{Prefácio}\n\n" + convertBlocks(doc.Content, ctx)
	case model.RoleAboutAuthor:
		return "\\chapter*{Sobre o Autor}\n\n" + convertBlocks(doc.Content, ctx)
	case model.RoleDedication:
		return renderDedication(doc, ctx)
	case model.RoleEpigraph:
		return renderEpigraph(doc, ctx)
	default:
		ctx.warnf("Papel de documento desconhecido '%s' — tratado como capítulo comum.", role)
		return convertBlocks(doc.Content, ctx)
	}
}

// renderDedication: página própria, sem numeração, texto em itálico
// alinhado à direita, deslocado ao terço superior da página (seção 5.5).
func renderDedication(doc node, ctx *Context) string {
	body := convertBlocks(doc.Content, ctx)
	return "\\clearpage\n\\thispagestyle{empty}\n\\vspace*{0.33\\textheight}\n" +
		"\\begin{flushright}\n\\itshape\n" + body + "\n\\end{flushright}\n\\clearpage"
}

// renderEpigraph: página própria, citação centralizada; o último parágrafo
// do documento é tratado como atribuição (alinhado à direita, precedido de
// travessão) — seção 5.5.
func renderEpigraph(doc node, ctx *Context) string {
	blocks := doc.Content
	var quote, attribution string

	if len(blocks) > 0 {
		last := blocks[len(blocks)-1]
		quote = convertBlocks(blocks[:len(blocks)-1], ctx)
		attribution = convertInline(last.Content, ctx)
	}

	var b strings.Builder
	if strings.TrimSpace(quote) != "" {
		b.WriteString("\\begin{center}\n" + quote + "\n\\end{center}\n")
	}
	if strings.TrimSpace(attribution) != "" {
		b.WriteString("\\begin{flushright}\n---" + attribution + "\n\\end{flushright}\n")
	}

	return "\\clearpage\n\\thispagestyle{empty}\n\\vspace*{0.33\\textheight}\n" + strings.TrimRight(b.String(), "\n") + "\n\\clearpage"
}
