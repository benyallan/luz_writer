package rules

import (
	"encoding/json"
	"testing"

	"luz-writer/internal/model"
)

func mustJSON(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestR001_MirroredEbookFires(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{
		Name: "Ebook", Kind: model.TargetKindEbook,
		PluginConfig: map[string]json.RawMessage{"geometry": mustJSON(t, map[string]any{"mirrored": true})},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}

	problems := R001(mustContext(t, w))
	if !hasCode(problems, "R001") {
		t.Fatalf("esperava R001, got %v", problems)
	}
	if problems[0].Severity != "error" {
		t.Errorf("R001 deveria ser error, got %s", problems[0].Severity)
	}
}

func TestR001_NotMirroredDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{
		Name: "Ebook", Kind: model.TargetKindEbook,
		PluginConfig: map[string]json.RawMessage{"geometry": mustJSON(t, map[string]any{"mirrored": false})},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}

	if problems := R001(mustContext(t, w)); hasCode(problems, "R001") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}

func TestR001_PrintTargetDoesNotFireEvenIfMirrored(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{
		Name: "Print", Kind: model.TargetKindPrint,
		PluginConfig: map[string]json.RawMessage{"geometry": mustJSON(t, map[string]any{"mirrored": true})},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	if problems := R001(mustContext(t, w)); hasCode(problems, "R001") {
		t.Errorf("não deveria disparar em target print: %v", problems)
	}
}

func TestR002_LuzChapterInArticleFires(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{Name: "Artigo", DocumentClass: model.DocumentClassArticle})
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
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"luzChapter","attrs":{"numbered":true},"content":[{"type":"text","text":"Introdução"}]}]}`); err != nil {
		t.Fatal(err)
	}

	problems := R002(mustContext(t, w))
	if !hasCode(problems, "R002") {
		t.Fatalf("esperava R002, got %v", problems)
	}
}

func TestR002_NoArticleDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"luzChapter","attrs":{"numbered":true},"content":[{"type":"text","text":"Introdução"}]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R002(mustContext(t, w)); hasCode(problems, "R002") {
		t.Errorf("não deveria disparar sem target article: %v", problems)
	}
}

func TestR003_CatalogRecordInEbookFires(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{Name: "Ebook", Kind: model.TargetKindEbook})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	if err := w.SetPluginEnabled("catalogRecord", true); err != nil {
		t.Fatal(err)
	}
	if problems := R003(mustContext(t, w)); !hasCode(problems, "R003") {
		t.Fatalf("esperava R003, got %v", problems)
	}
}

func TestR003_DisabledDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{Name: "Ebook", Kind: model.TargetKindEbook})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	if problems := R003(mustContext(t, w)); hasCode(problems, "R003") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}

func TestR004_PrintNotMirroredFires(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{
		Name: "Print", Kind: model.TargetKindPrint,
		PluginConfig: map[string]json.RawMessage{"geometry": mustJSON(t, map[string]any{"mirrored": false})},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	if problems := R004(mustContext(t, w)); !hasCode(problems, "R004") {
		t.Fatalf("esperava R004, got %v", problems)
	}
}

func TestR004_PrintMirroredDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{
		Name: "Print", Kind: model.TargetKindPrint,
		PluginConfig: map[string]json.RawMessage{"geometry": mustJSON(t, map[string]any{"mirrored": true})},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	if problems := R004(mustContext(t, w)); hasCode(problems, "R004") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}

func TestR010_SpecialRoleInArticleFires(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{Name: "Artigo", DocumentClass: model.DocumentClassArticle})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := w.CreateChapter("Dedicatória", model.RoleDedication); err != nil {
		t.Fatal(err)
	}
	if problems := R010(mustContext(t, w)); !hasCode(problems, "R010") {
		t.Fatalf("esperava R010, got %v", problems)
	}
}

func TestR010_ChapterRoleDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{Name: "Artigo", DocumentClass: model.DocumentClassArticle})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := w.CreateChapter("Capítulo", model.RoleChapter); err != nil {
		t.Fatal(err)
	}
	if problems := R010(mustContext(t, w)); hasCode(problems, "R010") {
		t.Errorf("não deveria disparar para role chapter: %v", problems)
	}
}

func TestR011_CatalogRecordAndFrontmatterPdfpagesFires(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{
		Name: "T",
		PluginConfig: map[string]json.RawMessage{
			"pdfpages": mustJSON(t, map[string]any{"filePath": "anexos/x.pdf", "placement": "beforeFrontmatter"}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	if err := w.SetPluginEnabled("catalogRecord", true); err != nil {
		t.Fatal(err)
	}
	if err := w.SetPluginEnabled("pdfpages", true); err != nil {
		t.Fatal(err)
	}
	if problems := R011(mustContext(t, w)); !hasCode(problems, "R011") {
		t.Fatalf("esperava R011, got %v", problems)
	}
}

func TestR011_AfterBackmatterDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{
		Name: "T",
		PluginConfig: map[string]json.RawMessage{
			"pdfpages": mustJSON(t, map[string]any{"filePath": "anexos/x.pdf", "placement": "afterBackmatter"}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	if err := w.SetPluginEnabled("catalogRecord", true); err != nil {
		t.Fatal(err)
	}
	if err := w.SetPluginEnabled("pdfpages", true); err != nil {
		t.Fatal(err)
	}
	if problems := R011(mustContext(t, w)); hasCode(problems, "R011") {
		t.Errorf("não deveria disparar com placement afterBackmatter: %v", problems)
	}
}

func TestR011_OnlyOnePluginDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.SetPluginEnabled("catalogRecord", true); err != nil {
		t.Fatal(err)
	}
	if problems := R011(mustContext(t, w)); hasCode(problems, "R011") {
		t.Errorf("não deveria disparar com só um plugin: %v", problems)
	}
}
