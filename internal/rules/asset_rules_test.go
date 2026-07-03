package rules

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"luz-writer/internal/model"
)

func TestR008_MissingImageFires(t *testing.T) {
	w := newTestWorkspace(t)
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"luzImage","attrs":{"src":"imagens/nao-existe.png"}}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R008(mustContext(t, w)); !hasCode(problems, "R008") {
		t.Fatalf("esperava R008, got %v", problems)
	}
}

func TestR008_ExistingImageDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	if err := os.MkdirAll(filepath.Join(w.Root, "imagens"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.Root, "imagens", "grafico.png"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"luzImage","attrs":{"src":"imagens/grafico.png"}}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R008(mustContext(t, w)); hasCode(problems, "R008") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}

func TestR009_PdfpagesMissingFileFires(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.SetPluginEnabled("pdfpages", true); err != nil {
		t.Fatal(err)
	}
	target, err := w.SaveTarget(model.Target{
		Name:         "T",
		PluginConfig: map[string]json.RawMessage{"pdfpages": mustJSON(t, map[string]any{"filePath": "anexos/nao-existe.pdf"})},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	if problems := R009(mustContext(t, w)); !hasCode(problems, "R009") {
		t.Fatalf("esperava R009, got %v", problems)
	}
}

func TestR009_PdfpagesExistingFileDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	if err := os.MkdirAll(filepath.Join(w.Root, "anexos"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.Root, "anexos", "capa.pdf"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := w.SetPluginEnabled("pdfpages", true); err != nil {
		t.Fatal(err)
	}
	target, err := w.SaveTarget(model.Target{
		Name:         "T",
		PluginConfig: map[string]json.RawMessage{"pdfpages": mustJSON(t, map[string]any{"filePath": "anexos/capa.pdf"})},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	if problems := R009(mustContext(t, w)); hasCode(problems, "R009") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}

func TestR009_PdfpagesDisabledDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	if problems := R009(mustContext(t, w)); hasCode(problems, "R009") {
		t.Errorf("não deveria disparar com plugin desabilitado: %v", problems)
	}
}
