package plugins

import (
	"encoding/json"
	"fmt"
	"strings"

	"luz-writer/internal/model"
)

type customStylesPlugin struct{}

// CustomStyles habilita o gerenciador de Estilos Personalizados (seção 5.7)
// e a mark luzCustomStyle. Não tem configuração por target — os estilos são
// globais, em .luz/styles.json (carregados via BuildContext.Styles).
var CustomStyles = register(&customStylesPlugin{})

func (customStylesPlugin) Name() string        { return "customStyles" }
func (customStylesPlugin) DisplayName() string { return "Estilos Personalizados" }
func (customStylesPlugin) Description() string {
	return "Estilos de texto nomeados, compostos a partir de itálico, negrito, versalete e cor."
}
func (customStylesPlugin) Core() bool          { return false }
func (customStylesPlugin) DocumentScope() bool { return false }

func (customStylesPlugin) Schema() model.FormSchema {
	return model.FormSchema{Fields: []model.FormField{}}
}

func (customStylesPlugin) DefaultConfig() json.RawMessage { return json.RawMessage("{}") }

func (customStylesPlugin) Validate(cfg json.RawMessage, ctx model.BuildContext) []model.Problem {
	return nil
}

func (customStylesPlugin) Preamble(cfg json.RawMessage, ctx model.BuildContext) (string, error) {
	if len(ctx.Styles) == 0 {
		return "", nil
	}

	needsColor := false
	for _, s := range ctx.Styles {
		if s.Color != nil && *s.Color != "" {
			needsColor = true
			break
		}
	}

	var b strings.Builder
	if needsColor {
		b.WriteString("\\usepackage{xcolor}\n")
	}
	for _, s := range ctx.Styles {
		b.WriteString(CustomStyleCommand(s))
	}
	return b.String(), nil
}

func (customStylesPlugin) ScopedLaTeX(effectiveCfg json.RawMessage, ctx model.BuildContext) (string, string, error) {
	return "", "", nil
}

// CustomStyleCommandName deriva o nome do comando LaTeX a partir do id do
// estilo (seção 5.7): "termo-estrangeiro" → "luzstyleTermoEstrangeiro". A
// mesma convenção é usada (duplicada, para não acoplar internal/latex a
// internal/plugins) em internal/latex/marks.go para converter luzCustomStyle.
func CustomStyleCommandName(id string) string {
	var b strings.Builder
	for _, part := range strings.Split(id, "-") {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		b.WriteString(part[1:])
	}
	return "luzstyle" + b.String()
}

// CustomStyleCommand gera o \newcommand de um estilo (seção 5.7).
func CustomStyleCommand(s model.CustomStyle) string {
	name := CustomStyleCommandName(s.ID)
	body := "#1"
	if s.Italic {
		body = "\\textit{" + body + "}"
	}
	if s.Bold {
		body = "\\textbf{" + body + "}"
	}
	if s.SmallCaps {
		body = "\\textsc{" + body + "}"
	}
	if s.Color != nil && *s.Color != "" {
		body = fmt.Sprintf("\\textcolor[HTML]{%s}{%s}", strings.TrimPrefix(*s.Color, "#"), body)
	}
	return fmt.Sprintf("\\newcommand{\\%s}[1]{%s}\n", name, body)
}
