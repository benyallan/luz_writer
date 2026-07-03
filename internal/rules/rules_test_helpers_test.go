package rules

import (
	"testing"

	"luz-writer/internal/model"
	"luz-writer/internal/workspace"
)

func newTestWorkspace(t *testing.T) *workspace.Workspace {
	t.Helper()
	root := t.TempDir()
	if _, err := workspace.Create(root, "Livro de Teste", "Autora", "pt-BR"); err != nil {
		t.Fatalf("workspace.Create: %v", err)
	}
	return &workspace.Workspace{Root: root}
}

func mustContext(t *testing.T, w *workspace.Workspace) Context {
	t.Helper()
	ctx, err := NewContext(w)
	if err != nil {
		t.Fatalf("NewContext: %v", err)
	}
	return ctx
}

func hasCode(problems []model.Problem, code string) bool {
	for _, p := range problems {
		if p.Code == code {
			return true
		}
	}
	return false
}
