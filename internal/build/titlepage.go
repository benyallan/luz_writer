package build

import (
	"strings"

	"luz-writer/internal/latex"
	"luz-writer/internal/model"
)

// titleMetadata monta \title/\author a partir de project.json, para o
// \maketitle customizado (seção 12, Etapa 6). Sem \date explícito — o
// projeto não tem um campo de data, então o LaTeX usa o padrão (\today).
func titleMetadata(project model.Project) string {
	var b strings.Builder

	title := latex.EscapeText(project.Title)
	if project.Subtitle != "" {
		title += "\\\\[0.5em]\\large " + latex.EscapeText(project.Subtitle)
	}
	b.WriteString("\\title{" + title + "}\n")

	authors := make([]string, len(project.Authors))
	for i, a := range project.Authors {
		authors[i] = latex.EscapeText(a)
	}
	b.WriteString("\\author{" + strings.Join(authors, " \\and ") + "}\n")

	return b.String()
}

// titlePageBody é inserido logo após \begin{document}: a página de título
// (\maketitle) e, se o target pedir (seção 5.2, includeToc), o sumário —
// respeitando includeInToc de cada título (seção 6, convertHeading).
func titlePageBody(target model.Target) string {
	var b strings.Builder
	b.WriteString("\\maketitle\n\n")
	if target.IncludeToc {
		b.WriteString("\\tableofcontents\n\n")
	}
	return b.String()
}
