package build

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"luz-writer/internal/model"
	"luz-writer/internal/workspace"
)

// buildFixturePDF compila um .tex trivial com o tectonic real para produzir
// um PDF de anexo válido (um handcrafted em base64 falha ao ser reaberto
// pelo próprio tectonic/xdvipdfmx no \includepdf).
func buildFixturePDF(t *testing.T) []byte {
	t.Helper()
	dir := t.TempDir()
	texPath := filepath.Join(dir, "fixture.tex")
	if err := os.WriteFile(texPath, []byte("\\documentclass{article}\\begin{document}capa\\end{document}"), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("tectonic", "fixture.tex")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("falha ao gerar PDF de fixture: %v\n%s", err, out)
	}
	data, err := os.ReadFile(filepath.Join(dir, "fixture.pdf"))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// TestCompile_CatalogRecordAndPdfpages cobre os critérios de aceite: o
// formulário da ficha catalográfica gera a página no início do documento, e
// o PDF externo do pdfpages é inserido via \includepdf.
func TestCompile_CatalogRecordAndPdfpages(t *testing.T) {
	requireTectonic(t)

	root := t.TempDir()
	if _, err := workspace.Create(root, "Livro Com Plugins", "Autora Teste", "pt-BR"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	w := &workspace.Workspace{Root: root}

	pdfBytes := buildFixturePDF(t)
	if err := os.MkdirAll(filepath.Join(root, "anexos"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "anexos", "capa.pdf"), pdfBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := w.SetPluginEnabled("catalogRecord", true); err != nil {
		t.Fatal(err)
	}
	if err := w.SetPluginEnabled("pdfpages", true); err != nil {
		t.Fatal(err)
	}

	target, err := w.SaveTarget(model.Target{
		Name:          "Target com Plugins",
		Kind:          model.TargetKindPrint,
		DocumentClass: model.DocumentClassBook,
		FontSize:      "11pt",
		IncludeToc:    true,
		PluginConfig: map[string]json.RawMessage{
			"catalogRecord": jsonMarshal(t, map[string]any{
				"isbn": "978-0-00-000000-0", "publisher": "Editora Teste", "publisherCity": "São Paulo",
				"year": 2026, "cdd": "800",
			}),
			"pdfpages": jsonMarshal(t, map[string]any{
				"filePath": "anexos/capa.pdf", "placement": "beforeFrontmatter", "pages": "-",
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}

	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[
		{"type":"luzChapter","attrs":{"numbered":true,"includeInToc":true},"content":[{"type":"text","text":"Introdução"}]},
		{"type":"paragraph","content":[{"type":"text","text":"Texto."}]}
	]}`); err != nil {
		t.Fatal(err)
	}

	result, err := Compile(w, nil)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	if !result.Success {
		t.Fatalf("Compile falhou: problems=%v logTail=%s", result.Problems, result.LogTail)
	}

	tex, err := os.ReadFile(filepath.Join(root, ".tmp", "main.tex"))
	if err != nil {
		t.Fatal(err)
	}
	texStr := string(tex)

	includePDFIdx := strings.Index(texStr, "\\includepdf[pages={-}]{anexos/capa.pdf}")
	frontmatterIdx := strings.Index(texStr, "\\frontmatter")
	atBeginDocIdx := strings.Index(texStr, "\\AtBeginDocument")

	if includePDFIdx == -1 {
		t.Errorf("esperava \\includepdf no .tex:\n%s", texStr)
	}
	if atBeginDocIdx == -1 || !strings.Contains(texStr, "Editora Teste") || !strings.Contains(texStr, "978-0-00-000000-0") {
		t.Errorf("esperava ficha catalográfica com os dados preenchidos:\n%s", texStr)
	}
	if includePDFIdx == -1 || frontmatterIdx == -1 || includePDFIdx >= frontmatterIdx {
		t.Errorf("placement 'beforeFrontmatter' deveria colocar \\includepdf antes de \\frontmatter:\n%s", texStr)
	}

	if _, err := os.Stat(result.OutputPath); err != nil {
		t.Errorf("PDF não encontrado: %v", err)
	}
}

func jsonMarshal(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
