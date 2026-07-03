package rules

import (
	"testing"

	"luz-writer/internal/model"
)

func TestR012_OverrideForDisabledPluginFires(t *testing.T) {
	w := newTestWorkspace(t)
	ch, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveDocumentOverrides(ch.ID, `{"fancyhdr":{"headerLeft":"none"}}`); err != nil {
		t.Fatal(err)
	}
	if problems := R012(mustContext(t, w)); !hasCode(problems, "R012") {
		t.Fatalf("esperava R012, got %v", problems)
	}
}

func TestR012_OverrideForEnabledPluginDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.SetPluginEnabled("fancyhdr", true); err != nil {
		t.Fatal(err)
	}
	ch, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveDocumentOverrides(ch.ID, `{"fancyhdr":{"headerLeft":"none"}}`); err != nil {
		t.Fatal(err)
	}
	if problems := R012(mustContext(t, w)); hasCode(problems, "R012") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}

func TestR012_OverrideForCorePluginDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	ch, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveDocumentOverrides(ch.ID, `{"geometry":{"marginTop":"4cm"}}`); err != nil {
		t.Fatal(err)
	}
	if problems := R012(mustContext(t, w)); hasCode(problems, "R012") {
		t.Errorf("núcleo nunca está 'desabilitado': %v", problems)
	}
}

func TestR013_GeometryOverrideInEbookFires(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{Name: "Ebook", Kind: model.TargetKindEbook})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	ch, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveDocumentOverrides(ch.ID, `{"geometry":{"marginTop":"4cm"}}`); err != nil {
		t.Fatal(err)
	}
	if problems := R013(mustContext(t, w)); !hasCode(problems, "R013") {
		t.Fatalf("esperava R013, got %v", problems)
	}
}

func TestR013_GeometryOverrideInPrintDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{Name: "Print", Kind: model.TargetKindPrint})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	ch, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveDocumentOverrides(ch.ID, `{"geometry":{"marginTop":"4cm"}}`); err != nil {
		t.Fatal(err)
	}
	if problems := R013(mustContext(t, w)); hasCode(problems, "R013") {
		t.Errorf("não deveria disparar em target print: %v", problems)
	}
}

func TestR014_OrphanedOverrideFires(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.SaveDocumentOverrides("orfao", `{"geometry":{"marginTop":"4cm"}}`); err != nil {
		t.Fatal(err)
	}
	if problems := R014(mustContext(t, w)); !hasCode(problems, "R014") {
		t.Fatalf("esperava R014, got %v", problems)
	}
}

func TestR014_ValidOverrideDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	ch, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveDocumentOverrides(ch.ID, `{"geometry":{"marginTop":"4cm"}}`); err != nil {
		t.Fatal(err)
	}
	if problems := R014(mustContext(t, w)); hasCode(problems, "R014") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}
