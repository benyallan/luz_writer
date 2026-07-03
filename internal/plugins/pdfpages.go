package plugins

import (
	"encoding/json"
	"fmt"

	"luz-writer/internal/model"
)

type PdfpagesConfig struct {
	FilePath  string `json:"filePath"`
	Placement string `json:"placement"` // beforeFrontmatter | afterFrontmatter | afterBackmatter
	Pages     string `json:"pages"`
}

var PdfpagesDefaults = PdfpagesConfig{
	FilePath:  "",
	Placement: "beforeFrontmatter",
	Pages:     "-",
}

type pdfpagesPlugin struct{}

// Pdfpages é o plugin opcional de inserção de PDFs externos prontos (capa
// diagramada, ficha catalográfica fornecida pela editora etc.). A injeção
// posicional de \includepdf acontece em internal/build, que conhece os
// limites de frontmatter/backmatter da montagem — este plugin só contribui
// o \usepackage.
var Pdfpages = register(&pdfpagesPlugin{})

func (pdfpagesPlugin) Name() string        { return "pdfpages" }
func (pdfpagesPlugin) DisplayName() string { return "Inserir PDF Externo" }
func (pdfpagesPlugin) Description() string {
	return "Insere um PDF pronto (capa, ficha da editora) em uma posição do livro."
}
func (pdfpagesPlugin) Core() bool          { return false }
func (pdfpagesPlugin) DocumentScope() bool { return false }

func (pdfpagesPlugin) Schema() model.FormSchema {
	return model.FormSchema{Fields: []model.FormField{
		{Key: "filePath", Label: "Arquivo PDF", Type: model.FieldTypeText, Default: ""},
		{
			Key: "placement", Label: "Posição", Type: model.FieldTypeSelect, Default: PdfpagesDefaults.Placement,
			Options: []model.FieldOption{
				{Value: "beforeFrontmatter", Label: "Antes de tudo"},
				{Value: "afterFrontmatter", Label: "Após pré-textuais"},
				{Value: "afterBackmatter", Label: "Última(s) página(s)"},
			},
		},
		{Key: "pages", Label: "Páginas (ex.: \"-\" para todas)", Type: model.FieldTypeText, Default: PdfpagesDefaults.Pages},
	}}
}

func (pdfpagesPlugin) DefaultConfig() json.RawMessage {
	b, _ := json.Marshal(PdfpagesDefaults)
	return b
}

func (pdfpagesPlugin) Validate(cfg json.RawMessage, ctx model.BuildContext) []model.Problem {
	return nil
}

func (pdfpagesPlugin) Preamble(cfg json.RawMessage, ctx model.BuildContext) (string, error) {
	return "\\usepackage{pdfpages}\n", nil
}

func (pdfpagesPlugin) ScopedLaTeX(effectiveCfg json.RawMessage, ctx model.BuildContext) (string, string, error) {
	return "", "", nil
}

// IncludePDFCommand monta o \includepdf para a config resolvida — usado por
// internal/build para splicar o comando na posição certa da montagem.
func IncludePDFCommand(cfg json.RawMessage) (string, error) {
	c, err := parsePdfpagesConfig(cfg)
	if err != nil {
		return "", err
	}
	if c.FilePath == "" {
		return "", nil
	}
	pages := c.Pages
	if pages == "" {
		pages = "-"
	}
	return fmt.Sprintf("\\includepdf[pages={%s}]{%s}\n", pages, c.FilePath), nil
}

func parsePdfpagesConfig(cfg json.RawMessage) (PdfpagesConfig, error) {
	c := PdfpagesDefaults
	if len(cfg) == 0 {
		return c, nil
	}
	if err := json.Unmarshal(cfg, &c); err != nil {
		return PdfpagesConfig{}, fmt.Errorf("configuração de pdfpages inválida: %w", err)
	}
	return c, nil
}
