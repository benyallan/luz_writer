package workspace

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"luz-writer/internal/model"
)

func jsonEqual(t *testing.T, a, b string) bool {
	t.Helper()
	var va, vb any
	if err := json.Unmarshal([]byte(a), &va); err != nil {
		t.Fatalf("json inválido: %s", a)
	}
	if err := json.Unmarshal([]byte(b), &vb); err != nil {
		t.Fatalf("json inválido: %s", b)
	}
	return reflect.DeepEqual(va, vb)
}

func newTestWorkspace(t *testing.T) *Workspace {
	t.Helper()
	root := t.TempDir()
	if _, err := Create(root, "Livro de Teste", "Autora", "pt-BR"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	return &Workspace{Root: root}
}

func TestCreateChapter_DedicationTwoChaptersAndAboutAuthor(t *testing.T) {
	w := newTestWorkspace(t)

	dedication, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatalf("CreateChapter dedicatória: %v", err)
	}
	intro, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatalf("CreateChapter introdução: %v", err)
	}
	method, err := w.CreateChapter("Metodologia", model.RoleChapter)
	if err != nil {
		t.Fatalf("CreateChapter metodologia: %v", err)
	}
	about, err := w.CreateChapter("Sobre o Autor", model.RoleAboutAuthor)
	if err != nil {
		t.Fatalf("CreateChapter sobre o autor: %v", err)
	}

	ids := map[string]bool{dedication.ID: true, intro.ID: true, method.ID: true, about.ID: true}
	if len(ids) != 4 {
		t.Fatalf("ids não são únicos: %v", []string{dedication.ID, intro.ID, method.ID, about.ID})
	}

	for _, id := range []string{dedication.ID, intro.ID, method.ID, about.ID} {
		if _, err := os.Stat(filepath.Join(w.Root, "capitulos", id+".json")); err != nil {
			t.Errorf("arquivo do documento %q não foi criado: %v", id, err)
		}
	}

	metas, err := w.ListChapters()
	if err != nil {
		t.Fatalf("ListChapters: %v", err)
	}
	if len(metas) != 4 {
		t.Fatalf("len(metas) = %d, esperava 4", len(metas))
	}

	want := []struct {
		id    string
		role  model.DocumentRole
		title string
	}{
		{dedication.ID, model.RoleDedication, "Dedicatória"},
		{intro.ID, model.RoleChapter, "Introdução"},
		{method.ID, model.RoleChapter, "Metodologia"},
		{about.ID, model.RoleAboutAuthor, "Sobre o Autor"},
	}
	for i, w2 := range want {
		if metas[i].ID != w2.id || metas[i].Role != w2.role || metas[i].Title != w2.title {
			t.Errorf("metas[%d] = %+v, esperava id=%s role=%s title=%s", i, metas[i], w2.id, w2.role, w2.title)
		}
	}
}

func TestSaveChapterAndLoadChapter_RoundTrip(t *testing.T) {
	w := newTestWorkspace(t)
	meta, err := w.CreateChapter("Capítulo Um", model.RoleChapter)
	if err != nil {
		t.Fatalf("CreateChapter: %v", err)
	}

	newContent := `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Capítulo Um Editado"}]}]}`
	if err := w.SaveChapter(meta.ID, newContent); err != nil {
		t.Fatalf("SaveChapter: %v", err)
	}

	got, err := w.LoadChapter(meta.ID)
	if err != nil {
		t.Fatalf("LoadChapter: %v", err)
	}
	if !jsonEqual(t, got, newContent) {
		t.Errorf("LoadChapter = %s, esperava (estruturalmente) %s", got, newContent)
	}

	// role/luzVersion devem ser preservados.
	metas, err := w.ListChapters()
	if err != nil {
		t.Fatalf("ListChapters: %v", err)
	}
	if metas[0].Role != model.RoleChapter {
		t.Errorf("Role = %q após SaveChapter, esperava chapter", metas[0].Role)
	}
}

func TestSaveChapter_InvalidJSON(t *testing.T) {
	w := newTestWorkspace(t)
	meta, err := w.CreateChapter("Capítulo Um", model.RoleChapter)
	if err != nil {
		t.Fatalf("CreateChapter: %v", err)
	}
	if err := w.SaveChapter(meta.ID, "isso não é json"); err == nil {
		t.Fatal("esperava erro para conteúdo inválido")
	}
}

func TestDeleteChapter_RemovesFileAndOrder(t *testing.T) {
	w := newTestWorkspace(t)
	meta, err := w.CreateChapter("Capítulo Um", model.RoleChapter)
	if err != nil {
		t.Fatalf("CreateChapter: %v", err)
	}

	if err := w.DeleteChapter(meta.ID); err != nil {
		t.Fatalf("DeleteChapter: %v", err)
	}

	if _, err := os.Stat(filepath.Join(w.Root, "capitulos", meta.ID+".json")); !os.IsNotExist(err) {
		t.Errorf("arquivo do capítulo ainda existe após DeleteChapter")
	}

	project, err := w.Project()
	if err != nil {
		t.Fatalf("Project: %v", err)
	}
	for _, id := range project.ChapterOrder {
		if id == meta.ID {
			t.Errorf("chapterOrder ainda contém %q após DeleteChapter", meta.ID)
		}
	}
}

func TestDeleteChapter_UnknownID(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.DeleteChapter("99-inexistente"); err == nil {
		t.Fatal("esperava erro ao excluir documento inexistente")
	}
}

func TestReorderChapters(t *testing.T) {
	w := newTestWorkspace(t)
	a, _ := w.CreateChapter("A", model.RoleChapter)
	b, _ := w.CreateChapter("B", model.RoleChapter)
	c, _ := w.CreateChapter("C", model.RoleChapter)

	newOrder := []string{c.ID, a.ID, b.ID}
	if err := w.ReorderChapters(newOrder); err != nil {
		t.Fatalf("ReorderChapters: %v", err)
	}

	project, err := w.Project()
	if err != nil {
		t.Fatalf("Project: %v", err)
	}
	for i, id := range newOrder {
		if project.ChapterOrder[i] != id {
			t.Errorf("ChapterOrder[%d] = %q, esperava %q", i, project.ChapterOrder[i], id)
		}
	}
}

func TestReorderChapters_RejectsInvalidPermutation(t *testing.T) {
	w := newTestWorkspace(t)
	a, _ := w.CreateChapter("A", model.RoleChapter)
	b, _ := w.CreateChapter("B", model.RoleChapter)

	if err := w.ReorderChapters([]string{a.ID}); err == nil {
		t.Fatal("esperava erro para ordem com tamanho diferente")
	}
	if err := w.ReorderChapters([]string{a.ID, "99-inexistente"}); err == nil {
		t.Fatal("esperava erro para id desconhecido")
	}
	if err := w.ReorderChapters([]string{a.ID, a.ID}); err == nil {
		t.Fatal("esperava erro para id duplicado")
	}
	_ = b
}

func TestListChapters_CorruptedChapterFile(t *testing.T) {
	w := newTestWorkspace(t)
	meta, err := w.CreateChapter("Capítulo Um", model.RoleChapter)
	if err != nil {
		t.Fatalf("CreateChapter: %v", err)
	}

	chapterFile := filepath.Join(w.Root, "capitulos", meta.ID+".json")
	if err := os.WriteFile(chapterFile, []byte("{ não é json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err = w.ListChapters()
	if err == nil {
		t.Fatal("esperava erro para capítulo corrompido")
	}
	var corrupted *CorruptedFileError
	if !errors.As(err, &corrupted) {
		t.Fatalf("erro = %T, esperava *CorruptedFileError", err)
	}
}
