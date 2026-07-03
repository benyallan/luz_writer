package workspace

import (
	"testing"

	"luz-writer/internal/model"
)

func TestSaveTarget_GeneratesIDWhenEmpty(t *testing.T) {
	w := newTestWorkspace(t)
	saved, err := w.SaveTarget(model.Target{Name: "Amazon KDP 6x9"})
	if err != nil {
		t.Fatalf("SaveTarget: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("esperava id gerado")
	}

	targets, err := w.ListTargets()
	if err != nil {
		t.Fatalf("ListTargets: %v", err)
	}
	if len(targets) != 1 || targets[0].ID != saved.ID {
		t.Fatalf("ListTargets = %+v", targets)
	}
}

func TestSaveTarget_UpdatesExisting(t *testing.T) {
	w := newTestWorkspace(t)
	saved, err := w.SaveTarget(model.Target{Name: "Meu Target", FontSize: "11pt"})
	if err != nil {
		t.Fatal(err)
	}

	saved.FontSize = "12pt"
	if _, err := w.SaveTarget(saved); err != nil {
		t.Fatal(err)
	}

	targets, err := w.ListTargets()
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 || targets[0].FontSize != "12pt" {
		t.Fatalf("esperava atualização in-place, got %+v", targets)
	}
}

func TestDeleteTarget_ClearsActiveTarget(t *testing.T) {
	w := newTestWorkspace(t)
	saved, err := w.SaveTarget(model.Target{Name: "Target A"})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.SetActiveTarget(saved.ID); err != nil {
		t.Fatal(err)
	}

	if err := w.DeleteTarget(saved.ID); err != nil {
		t.Fatal(err)
	}

	project, err := w.Project()
	if err != nil {
		t.Fatal(err)
	}
	if project.ActiveTarget != "" {
		t.Errorf("activeTarget deveria ser limpo, got %q", project.ActiveTarget)
	}
}

func TestSetActiveTarget_UnknownID(t *testing.T) {
	w := newTestWorkspace(t)
	if err := w.SetActiveTarget("nao-existe"); err == nil {
		t.Fatal("esperava erro para target inexistente")
	}
}

func TestSetPluginEnabled_ToggleAndPersist(t *testing.T) {
	w := newTestWorkspace(t)

	if err := w.SetPluginEnabled("fancyhdr", true); err != nil {
		t.Fatal(err)
	}
	cfg, err := w.PluginsConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Enabled) != 1 || cfg.Enabled[0] != "fancyhdr" {
		t.Fatalf("esperava ['fancyhdr'], got %v", cfg.Enabled)
	}

	// idempotente ao habilitar de novo.
	if err := w.SetPluginEnabled("fancyhdr", true); err != nil {
		t.Fatal(err)
	}
	cfg, _ = w.PluginsConfig()
	if len(cfg.Enabled) != 1 {
		t.Fatalf("não deveria duplicar: %v", cfg.Enabled)
	}

	if err := w.SetPluginEnabled("fancyhdr", false); err != nil {
		t.Fatal(err)
	}
	cfg, _ = w.PluginsConfig()
	if len(cfg.Enabled) != 0 {
		t.Fatalf("esperava lista vazia, got %v", cfg.Enabled)
	}
}

func TestStyles_SaveAndList(t *testing.T) {
	w := newTestWorkspace(t)
	color := "#555555"
	styles := []model.CustomStyle{
		{ID: "termo-estrangeiro", Name: "Termo Estrangeiro", Italic: true},
		{ID: "voz-interior", Name: "Voz Interior", Italic: true, SmallCaps: true, Color: &color},
	}
	if err := w.SaveStyles(styles); err != nil {
		t.Fatal(err)
	}

	got, err := w.ListStyles()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[1].Color == nil || *got[1].Color != color {
		t.Fatalf("ListStyles = %+v", got)
	}
}

func TestStyles_EmptyByDefault(t *testing.T) {
	w := newTestWorkspace(t)
	got, err := w.ListStyles()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("esperava vazio, got %+v", got)
	}
}
