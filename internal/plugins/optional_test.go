package plugins

import (
	"strings"
	"testing"

	"luz-writer/internal/model"
)

func TestFancyhdr_PreambleAndScoped(t *testing.T) {
	ctx := model.BuildContext{Project: model.Project{Authors: []string{"Fulana"}}}
	pre, err := Fancyhdr.Preamble(Fancyhdr.DefaultConfig(), ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(pre, "\\usepackage{fancyhdr}") || !strings.Contains(pre, "\\pagestyle{luzfancy}") {
		t.Errorf("preâmbulo inesperado: %s", pre)
	}
	if !strings.Contains(pre, "Fulana") {
		t.Errorf("esperava autor no cabeçalho: %s", pre)
	}

	before, after, err := Fancyhdr.ScopedLaTeX(Fancyhdr.DefaultConfig(), ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(before, "\\pagestyle{luzfancyoverride}") {
		t.Errorf("before inesperado: %s", before)
	}
	if !strings.Contains(after, "\\pagestyle{luzfancy}") {
		t.Errorf("after inesperado: %s", after)
	}
}

func TestCatalogRecord_Preamble(t *testing.T) {
	cfg := []byte(`{"isbn":"123","publisher":"Editora X","publisherCity":"São Paulo","year":2024,"cdd":"800","subjectEntries":"Literatura; Ficção"}`)
	ctx := model.BuildContext{Project: model.Project{Title: "Meu Livro", Authors: []string{"Autora"}}}
	pre, err := CatalogRecord.Preamble(cfg, ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"AtBeginDocument", "Meu Livro", "Autora", "ISBN 123", "CDD 800", "1. Literatura", "1. Ficção"} {
		if !strings.Contains(pre, want) {
			t.Errorf("esperava %q no preâmbulo:\n%s", want, pre)
		}
	}
}

func TestPdfpages_IncludePDFCommand(t *testing.T) {
	cfg := []byte(`{"filePath":"anexos/capa.pdf","placement":"beforeFrontmatter","pages":"1-2"}`)
	cmd, err := IncludePDFCommand(cfg)
	if err != nil {
		t.Fatal(err)
	}
	want := "\\includepdf[pages={1-2}]{anexos/capa.pdf}\n"
	if cmd != want {
		t.Errorf("got %q, want %q", cmd, want)
	}
}

func TestPdfpages_IncludePDFCommand_EmptyFilePath(t *testing.T) {
	cmd, err := IncludePDFCommand(Pdfpages.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "" {
		t.Errorf("esperava comando vazio sem filePath, got %q", cmd)
	}
}

func TestHyphenation_OffAddsPackageAndScoped(t *testing.T) {
	ctx := model.BuildContext{}
	pre, err := Hyphenation.Preamble([]byte(`{"mode":"off"}`), ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(pre, "hyphenat") {
		t.Errorf("esperava pacote hyphenat: %s", pre)
	}

	preAuto, err := Hyphenation.Preamble([]byte(`{"mode":"auto"}`), ctx)
	if err != nil {
		t.Fatal(err)
	}
	if preAuto != "" {
		t.Errorf("modo auto não deveria adicionar pacote: %q", preAuto)
	}

	before, after, err := Hyphenation.ScopedLaTeX([]byte(`{"mode":"off"}`), ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(before, "10000") || after != "\\endgroup\n" {
		t.Errorf("before/after inesperados: %q / %q", before, after)
	}
}

func TestLanguages_NoOpPluginRegistered(t *testing.T) {
	p, ok := ByName("languages")
	if !ok {
		t.Fatal("languages não registrado")
	}
	if p.Core() {
		t.Error("languages não deveria ser núcleo")
	}
	if len(p.Schema().Fields) != 0 {
		t.Error("languages não deveria ter campos de configuração")
	}
}
