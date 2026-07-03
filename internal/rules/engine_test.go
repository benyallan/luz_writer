package rules

import (
	"testing"

	"luz-writer/internal/model"
)

func TestValidate_RunsAllRulesAndAggregates(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{Name: "Ebook", Kind: model.TargetKindEbook})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(target.ID); err != nil {
		t.Fatal(err)
	}
	// R001 dispara (mirrored=true é o default do módulo geometry).

	problems := Validate(mustContext(t, w))
	if !hasCode(problems, "R001") {
		t.Fatalf("esperava R001 no resultado agregado, got %v", problems)
	}
}

func TestValidate_CleanProjectHasNoProblems(t *testing.T) {
	w := newTestWorkspace(t)
	target, err := w.SaveTarget(model.Target{Name: "Print", Kind: model.TargetKindPrint})
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
		{"type":"luzChapter","attrs":{"numbered":true},"content":[{"type":"text","text":"Introdução"}]},
		{"type":"paragraph","content":[{"type":"text","text":"Texto real."}]}
	]}`); err != nil {
		t.Fatal(err)
	}

	if problems := Validate(mustContext(t, w)); len(problems) != 0 {
		t.Errorf("esperava projeto limpo sem problems, got %v", problems)
	}
}
