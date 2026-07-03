package rules

import (
	"testing"

	"luz-writer/internal/model"
)

func TestR005_EmptyChapterFires(t *testing.T) {
	w := newTestWorkspace(t)
	ch, err := w.CreateChapter("", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"   "}]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R005(mustContext(t, w)); !hasCode(problems, "R005") {
		t.Fatalf("esperava R005, got %v", problems)
	}
}

func TestR005_NonEmptyChapterDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"texto real"}]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R005(mustContext(t, w)); hasCode(problems, "R005") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}

func TestR006_MissingChapterFileFires(t *testing.T) {
	w := newTestWorkspace(t)
	project, err := w.Project()
	if err != nil {
		t.Fatal(err)
	}
	project.ChapterOrder = append(project.ChapterOrder, "99-fantasma")
	if err := w.SaveProject(project); err != nil {
		t.Fatal(err)
	}
	if problems := R006(mustContext(t, w)); !hasCode(problems, "R006") {
		t.Fatalf("esperava R006, got %v", problems)
	}
}

func TestR006_AllPresentDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	if _, err := w.CreateChapter("Introdução", model.RoleChapter); err != nil {
		t.Fatal(err)
	}
	if problems := R006(mustContext(t, w)); hasCode(problems, "R006") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}

func TestR007_DedicationAfterChapterFires(t *testing.T) {
	w := newTestWorkspace(t)
	if _, err := w.CreateChapter("Capítulo", model.RoleChapter); err != nil {
		t.Fatal(err)
	}
	if _, err := w.CreateChapter("Dedicatória", model.RoleDedication); err != nil {
		t.Fatal(err)
	}
	if problems := R007(mustContext(t, w)); !hasCode(problems, "R007") {
		t.Fatalf("esperava R007, got %v", problems)
	}
}

func TestR007_DedicationBeforeChapterDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	if _, err := w.CreateChapter("Dedicatória", model.RoleDedication); err != nil {
		t.Fatal(err)
	}
	if _, err := w.CreateChapter("Capítulo", model.RoleChapter); err != nil {
		t.Fatal(err)
	}
	if problems := R007(mustContext(t, w)); hasCode(problems, "R007") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}
