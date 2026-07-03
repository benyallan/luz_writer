package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"luz-writer/internal/model"
	"luz-writer/internal/plugins/presets"
	"luz-writer/internal/workspace"
)

// TestCompile_DifferentPresetsProduceDifferentGeometries cobre o critério de
// aceite da Etapa 4: compilar com o preset Amazon KDP 6x9 e com o preset
// e-Book gera preâmbulos (e portanto PDFs) com geometrias visivelmente
// diferentes.
func TestCompile_DifferentPresetsProduceDifferentGeometries(t *testing.T) {
	requireTectonic(t)

	root := t.TempDir()
	if _, err := workspace.Create(root, "Livro de Presets", "Autora", "pt-BR"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	w := &workspace.Workspace{Root: root}

	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatalf("CreateChapter: %v", err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[
		{"type":"luzChapter","attrs":{"numbered":true,"includeInToc":true},"content":[{"type":"text","text":"Introdução"}]},
		{"type":"paragraph","content":[{"type":"text","text":"Texto de teste."}]}
	]}`); err != nil {
		t.Fatalf("SaveChapter: %v", err)
	}

	var kdp, ebook model.Target
	for _, p := range presets.All {
		switch p.ID {
		case "amazon-kdp-6x9":
			kdp = p
		case "ebook-pdf-fluido":
			ebook = p
		}
	}
	if kdp.ID == "" || ebook.ID == "" {
		t.Fatal("presets esperados não encontrados em presets.All")
	}

	if _, err := w.SaveTarget(kdp); err != nil {
		t.Fatalf("SaveTarget kdp: %v", err)
	}
	if _, err := w.SaveTarget(ebook); err != nil {
		t.Fatalf("SaveTarget ebook: %v", err)
	}

	compileWith := func(targetID string) string {
		t.Helper()
		if err := w.SetActiveTarget(targetID); err != nil {
			t.Fatalf("SetActiveTarget(%s): %v", targetID, err)
		}
		result, err := Compile(w, nil)
		if err != nil {
			t.Fatalf("Compile(%s): %v", targetID, err)
		}
		if !result.Success {
			t.Fatalf("Compile(%s) falhou: problems=%v logTail=%s", targetID, result.Problems, result.LogTail)
		}
		if _, err := os.Stat(result.OutputPath); err != nil {
			t.Fatalf("PDF de %s não encontrado: %v", targetID, err)
		}
		texBytes, err := os.ReadFile(filepath.Join(root, ".tmp", "main.tex"))
		if err != nil {
			t.Fatalf("main.tex de %s não encontrado: %v", targetID, err)
		}
		return string(texBytes)
	}

	kdpTex := compileWith(kdp.ID)
	ebookTex := compileWith(ebook.ID)

	if !strings.Contains(kdpTex, "paperwidth=6in,paperheight=9in") {
		t.Errorf("preâmbulo do KDP não tem a geometria esperada:\n%s", kdpTex)
	}
	if !strings.Contains(kdpTex, "twoside") {
		t.Errorf("preâmbulo do KDP deveria ser twoside (espelhado):\n%s", kdpTex)
	}
	if !strings.Contains(ebookTex, "paperwidth=14.8cm,paperheight=21cm") {
		t.Errorf("preâmbulo do e-Book não tem a geometria esperada:\n%s", ebookTex)
	}
	if !strings.Contains(ebookTex, "oneside") {
		t.Errorf("preâmbulo do e-Book deveria ser oneside (sem espelhamento):\n%s", ebookTex)
	}
	if kdpTex == ebookTex {
		t.Error("os dois .tex deveriam ser diferentes")
	}

	// Nomes de arquivo distintos por target (seção 10, passo 6).
	kdpPDF := filepath.Join(root, "dist", "livro-de-presets-amazon-kdp-6x9.pdf")
	ebookPDF := filepath.Join(root, "dist", "livro-de-presets-ebook-pdf-fluido.pdf")
	if _, err := os.Stat(kdpPDF); err != nil {
		t.Errorf("PDF do KDP não encontrado em %s: %v", kdpPDF, err)
	}
	if _, err := os.Stat(ebookPDF); err != nil {
		t.Errorf("PDF do e-Book não encontrado em %s: %v", ebookPDF, err)
	}
}
