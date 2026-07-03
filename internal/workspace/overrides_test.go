package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"luz-writer/internal/model"
)

func TestGetDocumentOverrides_EmptyWhenNoFile(t *testing.T) {
	w := newTestWorkspace(t)
	got, err := w.GetDocumentOverrides("00-dedicatoria")
	if err != nil {
		t.Fatal(err)
	}
	if got != "{}" {
		t.Errorf("got %q, want {}", got)
	}
}

func TestSaveDocumentOverrides_CreatesFileOnDemand(t *testing.T) {
	w := newTestWorkspace(t)
	meta, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatal(err)
	}

	overridesJSON := `{"geometry":{"marginTop":"4cm","marginBottom":"4cm"}}`
	if err := w.SaveDocumentOverrides(meta.ID, overridesJSON); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(w.Root, ".luz", "overrides", meta.ID+".json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("arquivo de override não foi criado: %v", err)
	}

	got, err := w.GetDocumentOverrides(meta.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(t, got, overridesJSON) {
		t.Errorf("got %s, want (estruturalmente) %s", got, overridesJSON)
	}

	if !w.HasDocumentOverrides(meta.ID) {
		t.Error("HasDocumentOverrides deveria ser true")
	}
}

func TestSaveDocumentOverrides_EmptyDeletesFile(t *testing.T) {
	w := newTestWorkspace(t)
	meta, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatal(err)
	}

	if err := w.SaveDocumentOverrides(meta.ID, `{"geometry":{"marginTop":"4cm"}}`); err != nil {
		t.Fatal(err)
	}
	if err := w.SaveDocumentOverrides(meta.ID, `{}`); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(w.Root, ".luz", "overrides", meta.ID+".json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("arquivo de override deveria ter sido apagado")
	}
	if w.HasDocumentOverrides(meta.ID) {
		t.Error("HasDocumentOverrides deveria ser false")
	}
}

func TestDeleteChapter_CascadesOverrides(t *testing.T) {
	w := newTestWorkspace(t)
	meta, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveDocumentOverrides(meta.ID, `{"geometry":{"marginTop":"4cm"}}`); err != nil {
		t.Fatal(err)
	}

	if err := w.DeleteChapter(meta.ID); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(w.Root, ".luz", "overrides", meta.ID+".json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("override deveria ter sido excluído em cascata")
	}
}

func TestListChapters_HasOverridesFlag(t *testing.T) {
	w := newTestWorkspace(t)
	meta, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatal(err)
	}

	metas, err := w.ListChapters()
	if err != nil {
		t.Fatal(err)
	}
	if metas[0].HasOverrides {
		t.Error("não deveria ter overrides ainda")
	}

	if err := w.SaveDocumentOverrides(meta.ID, `{"geometry":{"marginTop":"4cm"}}`); err != nil {
		t.Fatal(err)
	}

	metas, err = w.ListChapters()
	if err != nil {
		t.Fatal(err)
	}
	if !metas[0].HasOverrides {
		t.Error("esperava HasOverrides=true após salvar sobrescritas")
	}
}
