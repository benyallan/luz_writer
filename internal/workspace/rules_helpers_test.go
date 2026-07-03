package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"luz-writer/internal/model"
)

func TestActiveTarget_FallsBackToFirstThenDefault(t *testing.T) {
	w := newTestWorkspace(t)

	// Sem nenhum target salvo: cai para o embutido.
	got, err := w.ActiveTarget()
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != DefaultTarget().Name {
		t.Errorf("esperava target embutido, got %+v", got)
	}

	// Cria um target sem marcá-lo ativo: cai para o primeiro salvo.
	saved, err := w.SaveTarget(model.Target{Name: "Meu Target"})
	if err != nil {
		t.Fatal(err)
	}
	got, err = w.ActiveTarget()
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != saved.ID {
		t.Errorf("esperava cair para o primeiro target salvo, got %+v", got)
	}

	// Marca como ativo explicitamente.
	if err := w.SetActiveTarget(saved.ID); err != nil {
		t.Fatal(err)
	}
	got, err = w.ActiveTarget()
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != saved.ID {
		t.Errorf("esperava o target ativo explícito, got %+v", got)
	}
}

func TestChaptersTolerant_ReportsMissingWithoutAborting(t *testing.T) {
	w := newTestWorkspace(t)
	meta, err := w.CreateChapter("Capítulo Um", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}

	// Injeta uma referência quebrada diretamente no project.json.
	project, err := w.Project()
	if err != nil {
		t.Fatal(err)
	}
	project.ChapterOrder = append(project.ChapterOrder, "99-fantasma")
	if err := w.SaveProject(project); err != nil {
		t.Fatal(err)
	}

	chapters, missing, err := w.ChaptersTolerant()
	if err != nil {
		t.Fatalf("ChaptersTolerant não deveria abortar: %v", err)
	}
	if len(chapters) != 1 || chapters[0].ID != meta.ID {
		t.Errorf("esperava só o capítulo existente, got %+v", chapters)
	}
	if len(missing) != 1 || missing[0] != "99-fantasma" {
		t.Errorf("esperava ['99-fantasma'] em missing, got %v", missing)
	}
}

func TestOrphanedOverrideIDs(t *testing.T) {
	w := newTestWorkspace(t)
	meta, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveDocumentOverrides(meta.ID, `{"geometry":{"marginTop":"4cm"}}`); err != nil {
		t.Fatal(err)
	}
	if err := w.SaveDocumentOverrides("orfao-nao-existe", `{"geometry":{"marginTop":"4cm"}}`); err != nil {
		t.Fatal(err)
	}

	orphans, err := w.OrphanedOverrideIDs([]string{meta.ID})
	if err != nil {
		t.Fatal(err)
	}
	if len(orphans) != 1 || orphans[0] != "orfao-nao-existe" {
		t.Errorf("esperava ['orfao-nao-existe'], got %v", orphans)
	}
}

func TestFileExistsRelative(t *testing.T) {
	w := newTestWorkspace(t)
	if w.FileExistsRelative("imagens/nao-existe.png") {
		t.Error("não deveria existir")
	}
	if err := os.MkdirAll(filepath.Join(w.Root, "imagens"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.Root, "imagens", "existe.png"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !w.FileExistsRelative("imagens/existe.png") {
		t.Error("deveria existir")
	}
	if w.FileExistsRelative("") {
		t.Error("caminho vazio nunca existe")
	}
}
