package plugins

import (
	"encoding/json"
	"fmt"

	"luz-writer/internal/model"
)

type BodyTextConfig struct {
	LineSpacing    string `json:"lineSpacing"`
	ParagraphStyle string `json:"paragraphStyle"`
}

var BodyTextDefaults = BodyTextConfig{
	LineSpacing:    "1.0",
	ParagraphStyle: "indent",
}

type bodyTextPlugin struct{}

// BodyText é o módulo do núcleo de espaçamento entre linhas e estilo de
// parágrafo.
var BodyText = register(&bodyTextPlugin{})

func (bodyTextPlugin) Name() string        { return "bodyText" }
func (bodyTextPlugin) DisplayName() string { return "Corpo do Texto" }
func (bodyTextPlugin) Description() string {
	return "Espaçamento entre linhas e estilo de parágrafo."
}
func (bodyTextPlugin) Core() bool          { return true }
func (bodyTextPlugin) DocumentScope() bool { return true }

func (bodyTextPlugin) Schema() model.FormSchema {
	return model.FormSchema{Fields: []model.FormField{
		{
			Key: "lineSpacing", Label: "Espaçamento entre linhas", Type: model.FieldTypeSelect, Default: BodyTextDefaults.LineSpacing,
			Options: []model.FieldOption{
				{Value: "1.0", Label: "Simples (1.0)"},
				{Value: "1.15", Label: "1.15"},
				{Value: "1.5", Label: "1.5"},
				{Value: "2.0", Label: "Duplo (2.0)"},
			},
		},
		{
			Key: "paragraphStyle", Label: "Estilo de parágrafo", Type: model.FieldTypeSelect, Default: BodyTextDefaults.ParagraphStyle,
			Options: []model.FieldOption{
				{Value: "indent", Label: "Recuo na primeira linha"},
				{Value: "spaced", Label: "Sem recuo, com espaço entre parágrafos"},
			},
		},
	}}
}

func (bodyTextPlugin) DefaultConfig() json.RawMessage {
	b, _ := json.Marshal(BodyTextDefaults)
	return b
}

func (bodyTextPlugin) Validate(cfg json.RawMessage, ctx model.BuildContext) []model.Problem {
	return nil
}

func (bodyTextPlugin) Preamble(cfg json.RawMessage, ctx model.BuildContext) (string, error) {
	c, err := parseBodyTextConfig(cfg)
	if err != nil {
		return "", err
	}
	var b string
	b += "\\usepackage{setspace}\n"
	b += fmt.Sprintf("\\setstretch{%s}\n", c.LineSpacing)
	if c.ParagraphStyle == "spaced" {
		b += "\\setlength{\\parindent}{0pt}\n\\setlength{\\parskip}{1em}\n"
	}
	return b, nil
}

func (bodyTextPlugin) ScopedLaTeX(effectiveCfg json.RawMessage, ctx model.BuildContext) (before string, after string, err error) {
	c, err := parseBodyTextConfig(effectiveCfg)
	if err != nil {
		return "", "", err
	}
	before = fmt.Sprintf("\\begingroup\\setstretch{%s}\n", c.LineSpacing)
	if c.ParagraphStyle == "spaced" {
		before += "\\setlength{\\parindent}{0pt}\\setlength{\\parskip}{1em}\n"
	}
	after = "\\endgroup\n"
	return before, after, nil
}

func parseBodyTextConfig(cfg json.RawMessage) (BodyTextConfig, error) {
	c := BodyTextDefaults
	if len(cfg) == 0 {
		return c, nil
	}
	if err := json.Unmarshal(cfg, &c); err != nil {
		return BodyTextConfig{}, fmt.Errorf("configuração de bodyText inválida: %w", err)
	}
	return c, nil
}
