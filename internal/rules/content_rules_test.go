package rules

import (
	"encoding/json"
	"testing"

	"luz-writer/internal/model"
)

func TestR015_LuzLangWithPluginDisabledFires(t *testing.T) {
	w := newTestWorkspace(t)
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[
		{"type":"text","marks":[{"type":"luzLang","attrs":{"lang":"en"}}],"text":"hello"}
	]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R015(mustContext(t, w)); !hasCode(problems, "R015") {
		t.Fatalf("esperava R015, got %v", problems)
	}
}

func TestR015_LuzLangWithPluginEnabledDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.SetPluginEnabled("languages", true); err != nil {
		t.Fatal(err)
	}
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[
		{"type":"text","marks":[{"type":"luzLang","attrs":{"lang":"en"}}],"text":"hello"}
	]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R015(mustContext(t, w)); hasCode(problems, "R015") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}

func TestR016_SoftHyphenWithPluginDisabledFires(t *testing.T) {
	w := newTestWorkspace(t)
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[
		{"type":"text","text":"pseudo"},{"type":"luzSoftHyphen"},{"type":"text","text":"random"}
	]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R016(mustContext(t, w)); !hasCode(problems, "R016") {
		t.Fatalf("esperava R016, got %v", problems)
	}
}

func TestR016_SoftHyphenWithModeOffFires(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.SetPluginEnabled("hyphenation", true); err != nil {
		t.Fatal(err)
	}
	target, err := w.SaveTarget(model.Target{
		Name:         "T",
		PluginConfig: map[string]json.RawMessage{"hyphenation": mustJSON(t, map[string]any{"mode": "off"})},
	})
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
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[
		{"type":"text","text":"pseudo"},{"type":"luzSoftHyphen"},{"type":"text","text":"random"}
	]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R016(mustContext(t, w)); !hasCode(problems, "R016") {
		t.Fatalf("esperava R016 com modo off, got %v", problems)
	}
}

func TestR016_SoftHyphenWithPluginEnabledAutoDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.SetPluginEnabled("hyphenation", true); err != nil {
		t.Fatal(err)
	}
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[
		{"type":"text","text":"pseudo"},{"type":"luzSoftHyphen"},{"type":"text","text":"random"}
	]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R016(mustContext(t, w)); hasCode(problems, "R016") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}

func TestR017_MissingVariableFires(t *testing.T) {
	w := newTestWorkspace(t)
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[
		{"type":"luzVariable","attrs":{"name":"inexistente"}}
	]}]}`); err != nil {
		t.Fatal(err)
	}
	problems := R017(mustContext(t, w))
	if !hasCode(problems, "R017") {
		t.Fatalf("esperava R017, got %v", problems)
	}
	if problems[0].Severity != "error" {
		t.Errorf("R017 deveria ser error, got %s", problems[0].Severity)
	}
}

func TestR017_ExistingVariableDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	project, err := w.Project()
	if err != nil {
		t.Fatal(err)
	}
	project.Variables = []model.Variable{{Name: "heroi", Value: "Ana"}}
	if err := w.SaveProject(project); err != nil {
		t.Fatal(err)
	}
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[
		{"type":"luzVariable","attrs":{"name":"heroi"}}
	]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R017(mustContext(t, w)); hasCode(problems, "R017") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}

func TestR018_UnknownStyleFires(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.SetPluginEnabled("customStyles", true); err != nil {
		t.Fatal(err)
	}
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[
		{"type":"text","marks":[{"type":"luzCustomStyle","attrs":{"styleId":"nao-existe"}}],"text":"x"}
	]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R018(mustContext(t, w)); !hasCode(problems, "R018") {
		t.Fatalf("esperava R018, got %v", problems)
	}
}

func TestR018_PluginDisabledFires(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.SaveStyles([]model.CustomStyle{{ID: "termo-estrangeiro", Name: "Termo Estrangeiro", Italic: true}}); err != nil {
		t.Fatal(err)
	}
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[
		{"type":"text","marks":[{"type":"luzCustomStyle","attrs":{"styleId":"termo-estrangeiro"}}],"text":"x"}
	]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R018(mustContext(t, w)); !hasCode(problems, "R018") {
		t.Fatalf("esperava R018 (plugin desabilitado mesmo com estilo existente), got %v", problems)
	}
}

func TestR018_KnownStyleWithPluginEnabledDoesNotFire(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.SetPluginEnabled("customStyles", true); err != nil {
		t.Fatal(err)
	}
	if err := w.SaveStyles([]model.CustomStyle{{ID: "termo-estrangeiro", Name: "Termo Estrangeiro", Italic: true}}); err != nil {
		t.Fatal(err)
	}
	ch, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SaveChapter(ch.ID, `{"type":"doc","content":[{"type":"paragraph","content":[
		{"type":"text","marks":[{"type":"luzCustomStyle","attrs":{"styleId":"termo-estrangeiro"}}],"text":"x"}
	]}]}`); err != nil {
		t.Fatal(err)
	}
	if problems := R018(mustContext(t, w)); hasCode(problems, "R018") {
		t.Errorf("não deveria disparar: %v", problems)
	}
}
