package plugins

import (
	"encoding/json"
	"fmt"

	"luz-writer/internal/model"
)

// GeometryDefaults é o valor padrão do módulo geometry (seção 8.3): A5,
// margens 2,5cm, espelhado — suficiente para qualquer target funcionar sem
// configuração.
type GeometryConfig struct {
	PaperWidth   string `json:"paperWidth"`
	PaperHeight  string `json:"paperHeight"`
	MarginInner  string `json:"marginInner"`
	MarginOuter  string `json:"marginOuter"`
	MarginTop    string `json:"marginTop"`
	MarginBottom string `json:"marginBottom"`
	Mirrored     bool   `json:"mirrored"`
}

var GeometryDefaults = GeometryConfig{
	PaperWidth:   "14.8cm",
	PaperHeight:  "21cm",
	MarginInner:  "2.5cm",
	MarginOuter:  "2.5cm",
	MarginTop:    "2.5cm",
	MarginBottom: "2.5cm",
	Mirrored:     true,
}

type geometryPlugin struct{}

// Geometry é o módulo do núcleo de dimensões de página e margens.
var Geometry = register(&geometryPlugin{})

func (geometryPlugin) Name() string        { return "geometry" }
func (geometryPlugin) DisplayName() string { return "Geometria de Página" }
func (geometryPlugin) Description() string {
	return "Dimensões de página, margens e espelhamento frente/verso."
}
func (geometryPlugin) Core() bool          { return true }
func (geometryPlugin) DocumentScope() bool { return true }

func (geometryPlugin) Schema() model.FormSchema {
	return model.FormSchema{Fields: []model.FormField{
		{Key: "paperWidth", Label: "Largura da página", Type: model.FieldTypeDimension, Default: GeometryDefaults.PaperWidth},
		{Key: "paperHeight", Label: "Altura da página", Type: model.FieldTypeDimension, Default: GeometryDefaults.PaperHeight},
		{Key: "marginInner", Label: "Margem interna", Type: model.FieldTypeDimension, Default: GeometryDefaults.MarginInner},
		{Key: "marginOuter", Label: "Margem externa", Type: model.FieldTypeDimension, Default: GeometryDefaults.MarginOuter},
		{Key: "marginTop", Label: "Margem superior", Type: model.FieldTypeDimension, Default: GeometryDefaults.MarginTop},
		{Key: "marginBottom", Label: "Margem inferior", Type: model.FieldTypeDimension, Default: GeometryDefaults.MarginBottom},
		{Key: "mirrored", Label: "Margens espelhadas (frente/verso)", Type: model.FieldTypeSwitch, Default: GeometryDefaults.Mirrored},
	}}
}

func (geometryPlugin) DefaultConfig() json.RawMessage {
	b, _ := json.Marshal(GeometryDefaults)
	return b
}

func (geometryPlugin) Validate(cfg json.RawMessage, ctx model.BuildContext) []model.Problem {
	return nil
}

func (p geometryPlugin) Preamble(cfg json.RawMessage, ctx model.BuildContext) (string, error) {
	c, err := parseGeometryConfig(cfg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\\usepackage[%s]{geometry}\n", geometryOptions(c)), nil
}

func (p geometryPlugin) ScopedLaTeX(effectiveCfg json.RawMessage, ctx model.BuildContext) (before string, after string, err error) {
	c, err := parseGeometryConfig(effectiveCfg)
	if err != nil {
		return "", "", err
	}
	before = fmt.Sprintf("\\newgeometry{%s}\n", geometryOptions(c))
	after = "\\restoregeometry\n"
	return before, after, nil
}

func parseGeometryConfig(cfg json.RawMessage) (GeometryConfig, error) {
	c := GeometryDefaults
	if len(cfg) == 0 {
		return c, nil
	}
	if err := json.Unmarshal(cfg, &c); err != nil {
		return GeometryConfig{}, fmt.Errorf("configuração de geometry inválida: %w", err)
	}
	return c, nil
}

// geometryOptions monta as opções do pacote geometry. "oneside" não é uma
// chave válida do geometry.sty (a ausência de "twoside" já implica oneside)
// — só twoside é emitido, e apenas quando mirrored=true.
func geometryOptions(c GeometryConfig) string {
	base := fmt.Sprintf("paperwidth=%s,paperheight=%s,top=%s,bottom=%s", c.PaperWidth, c.PaperHeight, c.MarginTop, c.MarginBottom)
	if c.Mirrored {
		return base + fmt.Sprintf(",inner=%s,outer=%s,twoside", c.MarginInner, c.MarginOuter)
	}
	return base + fmt.Sprintf(",left=%s,right=%s", c.MarginInner, c.MarginOuter)
}
