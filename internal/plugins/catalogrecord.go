package plugins

import (
	"encoding/json"
	"fmt"
	"strings"

	"luz-writer/internal/latex"
	"luz-writer/internal/model"
)

type CatalogRecordConfig struct {
	ISBN           string `json:"isbn"`
	Publisher      string `json:"publisher"`
	PublisherCity  string `json:"publisherCity"`
	Year           int    `json:"year"`
	CDD            string `json:"cdd"`
	CDU            string `json:"cdu"`
	SubjectEntries string `json:"subjectEntries"`
	PreparedBy     string `json:"preparedBy"`
}

var CatalogRecordDefaults = CatalogRecordConfig{}

type catalogRecordPlugin struct{}

// CatalogRecord é o plugin opcional de Ficha Catalográfica (padrão ABNT),
// renderizada no verso da página de rosto.
var CatalogRecord = register(&catalogRecordPlugin{})

func (catalogRecordPlugin) Name() string        { return "catalogRecord" }
func (catalogRecordPlugin) DisplayName() string { return "Ficha Catalográfica" }
func (catalogRecordPlugin) Description() string {
	return "Gera a ficha catalográfica (padrão ABNT) a partir de um formulário."
}
func (catalogRecordPlugin) Core() bool          { return false }
func (catalogRecordPlugin) DocumentScope() bool { return false }

func (catalogRecordPlugin) Schema() model.FormSchema {
	return model.FormSchema{Fields: []model.FormField{
		{Key: "isbn", Label: "ISBN", Type: model.FieldTypeText, Default: ""},
		{Key: "publisher", Label: "Editora", Type: model.FieldTypeText, Default: ""},
		{Key: "publisherCity", Label: "Cidade", Type: model.FieldTypeText, Default: ""},
		{Key: "year", Label: "Ano", Type: model.FieldTypeNumber, Default: 0},
		{Key: "cdd", Label: "CDD", Type: model.FieldTypeText, Default: ""},
		{Key: "cdu", Label: "CDU (opcional)", Type: model.FieldTypeText, Default: ""},
		{Key: "subjectEntries", Label: "Assuntos — separados por ponto-e-vírgula", Type: model.FieldTypeText, Default: ""},
		{Key: "preparedBy", Label: "Elaborada por (nome e CRB, opcional)", Type: model.FieldTypeText, Default: ""},
	}}
}

func (catalogRecordPlugin) DefaultConfig() json.RawMessage {
	b, _ := json.Marshal(CatalogRecordDefaults)
	return b
}

func (catalogRecordPlugin) Validate(cfg json.RawMessage, ctx model.BuildContext) []model.Problem {
	return nil
}

func (catalogRecordPlugin) Preamble(cfg json.RawMessage, ctx model.BuildContext) (string, error) {
	c, err := parseCatalogRecordConfig(cfg)
	if err != nil {
		return "", err
	}
	box := buildCatalogRecordBox(c, ctx.Project)
	return "\\AtBeginDocument{\\clearpage\\thispagestyle{empty}\n" + box + "\\clearpage}\n", nil
}

func (catalogRecordPlugin) ScopedLaTeX(effectiveCfg json.RawMessage, ctx model.BuildContext) (string, string, error) {
	return "", "", nil
}

func buildCatalogRecordBox(c CatalogRecordConfig, project model.Project) string {
	esc := latex.EscapeText
	var lines []string

	if project.Title != "" {
		lines = append(lines, esc(project.Title))
	}
	if len(project.Authors) > 0 {
		lines = append(lines, esc(strings.Join(project.Authors, "; ")))
	}
	if c.ISBN != "" {
		lines = append(lines, "ISBN "+esc(c.ISBN))
	}
	place := strings.TrimSpace(c.PublisherCity)
	pub := strings.TrimSpace(c.Publisher)
	if place != "" || pub != "" {
		lines = append(lines, esc(strings.TrimSpace(strings.TrimSuffix(place+" : "+pub, " : "))))
	}
	if c.Year > 0 {
		lines = append(lines, fmt.Sprintf("%d", c.Year))
	}
	if c.CDD != "" {
		lines = append(lines, "CDD "+esc(c.CDD))
	}
	if c.CDU != "" {
		lines = append(lines, "CDU "+esc(c.CDU))
	}
	if c.SubjectEntries != "" {
		for _, s := range strings.Split(c.SubjectEntries, ";") {
			s = strings.TrimSpace(s)
			if s != "" {
				lines = append(lines, "1. "+esc(s))
			}
		}
	}
	if c.PreparedBy != "" {
		lines = append(lines, "Ficha elaborada por "+esc(c.PreparedBy))
	}

	body := strings.Join(lines, " \\\\\n")
	return "\\begin{center}\n\\fbox{\\begin{minipage}{0.8\\linewidth}\n\\footnotesize\n" + body + "\n\\end{minipage}}\n\\end{center}\n"
}

func parseCatalogRecordConfig(cfg json.RawMessage) (CatalogRecordConfig, error) {
	c := CatalogRecordDefaults
	if len(cfg) == 0 {
		return c, nil
	}
	if err := json.Unmarshal(cfg, &c); err != nil {
		return CatalogRecordConfig{}, fmt.Errorf("configuração de catalogRecord inválida: %w", err)
	}
	return c, nil
}
