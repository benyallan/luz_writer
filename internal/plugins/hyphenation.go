package plugins

import (
	"encoding/json"
	"fmt"

	"luz-writer/internal/model"
)

type HyphenationConfig struct {
	Mode string `json:"mode"` // auto | off
}

var HyphenationDefaults = HyphenationConfig{Mode: "auto"}

type hyphenationPlugin struct{}

// Hyphenation controla a hifenização automática e habilita o botão de hífen
// sugerido (luzSoftHyphen) na toolbar.
var Hyphenation = register(&hyphenationPlugin{})

func (hyphenationPlugin) Name() string        { return "hyphenation" }
func (hyphenationPlugin) DisplayName() string { return "Hifenização" }
func (hyphenationPlugin) Description() string {
	return "Controla a hifenização automática e habilita hífens sugeridos."
}
func (hyphenationPlugin) Core() bool          { return false }
func (hyphenationPlugin) DocumentScope() bool { return true }

func (hyphenationPlugin) Schema() model.FormSchema {
	return model.FormSchema{Fields: []model.FormField{
		{
			Key: "mode", Label: "Hifenização", Type: model.FieldTypeSelect, Default: HyphenationDefaults.Mode,
			Options: []model.FieldOption{
				{Value: "auto", Label: "Automática"},
				{Value: "off", Label: "Desativada"},
			},
		},
	}}
}

func (hyphenationPlugin) DefaultConfig() json.RawMessage {
	b, _ := json.Marshal(HyphenationDefaults)
	return b
}

func (hyphenationPlugin) Validate(cfg json.RawMessage, ctx model.BuildContext) []model.Problem {
	return nil
}

func (hyphenationPlugin) Preamble(cfg json.RawMessage, ctx model.BuildContext) (string, error) {
	c, err := parseHyphenationConfig(cfg)
	if err != nil {
		return "", err
	}
	if c.Mode == "off" {
		return "\\usepackage[none]{hyphenat}\n", nil
	}
	return "", nil
}

func (hyphenationPlugin) ScopedLaTeX(effectiveCfg json.RawMessage, ctx model.BuildContext) (before string, after string, err error) {
	c, err := parseHyphenationConfig(effectiveCfg)
	if err != nil {
		return "", "", err
	}
	before = "\\begingroup"
	if c.Mode == "off" {
		before += "\\hyphenpenalty=10000\\exhyphenpenalty=10000"
	}
	before += "\n"
	after = "\\endgroup\n"
	return before, after, nil
}

func parseHyphenationConfig(cfg json.RawMessage) (HyphenationConfig, error) {
	c := HyphenationDefaults
	if len(cfg) == 0 {
		return c, nil
	}
	if err := json.Unmarshal(cfg, &c); err != nil {
		return HyphenationConfig{}, fmt.Errorf("configuração de hyphenation inválida: %w", err)
	}
	return c, nil
}
