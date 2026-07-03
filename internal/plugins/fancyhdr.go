package plugins

import (
	"encoding/json"
	"fmt"
	"strings"

	"luz-writer/internal/latex"
	"luz-writer/internal/model"
)

type FancyhdrConfig struct {
	HeaderLeft         string `json:"headerLeft"`
	HeaderRight        string `json:"headerRight"`
	PageNumberPosition string `json:"pageNumberPosition"`
}

var FancyhdrDefaults = FancyhdrConfig{
	HeaderLeft:         "author",
	HeaderRight:        "chapterTitle",
	PageNumberPosition: "outer-footer",
}

type fancyhdrPlugin struct{}

// Fancyhdr é o plugin opcional de cabeçalhos/rodapés e posição de numeração.
var Fancyhdr = register(&fancyhdrPlugin{})

func (fancyhdrPlugin) Name() string        { return "fancyhdr" }
func (fancyhdrPlugin) DisplayName() string { return "Cabeçalhos e Rodapés" }
func (fancyhdrPlugin) Description() string {
	return "Cabeçalhos, rodapés e posição da numeração de página."
}
func (fancyhdrPlugin) Core() bool          { return false }
func (fancyhdrPlugin) DocumentScope() bool { return true }

var headerContentOptions = []model.FieldOption{
	{Value: "author", Label: "Autor"},
	{Value: "chapterTitle", Label: "Título do capítulo"},
	{Value: "none", Label: "Nenhum"},
}

func (fancyhdrPlugin) Schema() model.FormSchema {
	return model.FormSchema{Fields: []model.FormField{
		{Key: "headerLeft", Label: "Cabeçalho esquerdo", Type: model.FieldTypeSelect, Default: FancyhdrDefaults.HeaderLeft, Options: headerContentOptions},
		{Key: "headerRight", Label: "Cabeçalho direito", Type: model.FieldTypeSelect, Default: FancyhdrDefaults.HeaderRight, Options: headerContentOptions},
		{
			Key: "pageNumberPosition", Label: "Posição do número de página", Type: model.FieldTypeSelect, Default: FancyhdrDefaults.PageNumberPosition,
			Options: []model.FieldOption{
				{Value: "outer-footer", Label: "Rodapé externo"},
				{Value: "center-footer", Label: "Rodapé central"},
			},
		},
	}}
}

func (fancyhdrPlugin) DefaultConfig() json.RawMessage {
	b, _ := json.Marshal(FancyhdrDefaults)
	return b
}

func (fancyhdrPlugin) Validate(cfg json.RawMessage, ctx model.BuildContext) []model.Problem {
	return nil
}

func (fancyhdrPlugin) Preamble(cfg json.RawMessage, ctx model.BuildContext) (string, error) {
	c, err := parseFancyhdrConfig(cfg)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString("\\usepackage{fancyhdr}\n")
	b.WriteString("\\fancypagestyle{luzfancy}{\n")
	b.WriteString(fancyhdrBody(c, ctx))
	b.WriteString("}\n\\pagestyle{luzfancy}\n")
	return b.String(), nil
}

func (fancyhdrPlugin) ScopedLaTeX(effectiveCfg json.RawMessage, ctx model.BuildContext) (before string, after string, err error) {
	c, err := parseFancyhdrConfig(effectiveCfg)
	if err != nil {
		return "", "", err
	}
	before = "\\fancypagestyle{luzfancyoverride}{\n" + fancyhdrBody(c, ctx) + "}\n\\pagestyle{luzfancyoverride}\n"
	after = "\\pagestyle{luzfancy}\n"
	return before, after, nil
}

func fancyhdrBody(c FancyhdrConfig, ctx model.BuildContext) string {
	var b strings.Builder
	b.WriteString("\\fancyhf{}\n")
	b.WriteString(fmt.Sprintf("\\fancyhead[L]{%s}\n", fancyhdrHeaderContent(c.HeaderLeft, ctx)))
	b.WriteString(fmt.Sprintf("\\fancyhead[R]{%s}\n", fancyhdrHeaderContent(c.HeaderRight, ctx)))
	if c.PageNumberPosition == "center-footer" {
		b.WriteString("\\fancyfoot[C]{\\thepage}\n")
	} else {
		b.WriteString("\\fancyfoot[LE,RO]{\\thepage}\n")
	}
	return b.String()
}

func fancyhdrHeaderContent(kind string, ctx model.BuildContext) string {
	switch kind {
	case "author":
		if len(ctx.Project.Authors) == 0 {
			return ""
		}
		return latex.EscapeText(strings.Join(ctx.Project.Authors, ", "))
	case "chapterTitle":
		return "\\leftmark"
	default:
		return ""
	}
}

func parseFancyhdrConfig(cfg json.RawMessage) (FancyhdrConfig, error) {
	c := FancyhdrDefaults
	if len(cfg) == 0 {
		return c, nil
	}
	if err := json.Unmarshal(cfg, &c); err != nil {
		return FancyhdrConfig{}, fmt.Errorf("configuração de fancyhdr inválida: %w", err)
	}
	return c, nil
}
