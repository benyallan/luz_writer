package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCreate_GeneratesStructureFromSpecSection4(t *testing.T) {
	root := filepath.Join(t.TempDir(), "meu-livro")

	info, err := Create(root, "Meu Livro", "Nome do Autor", "pt-BR")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if info.Project.Title != "Meu Livro" {
		t.Errorf("Title = %q, esperava %q", info.Project.Title, "Meu Livro")
	}
	if info.Project.LuzVersion != 1 {
		t.Errorf("LuzVersion = %d, esperava 1", info.Project.LuzVersion)
	}

	wantDirs := []string{
		".luz",
		filepath.Join(".luz", "targets"),
		filepath.Join(".luz", "overrides"),
		"capitulos",
		"imagens",
		"anexos",
		"dist",
		".tmp",
	}
	for _, d := range wantDirs {
		fi, err := os.Stat(filepath.Join(root, d))
		if err != nil {
			t.Errorf("diretório %q ausente: %v", d, err)
			continue
		}
		if !fi.IsDir() {
			t.Errorf("%q deveria ser um diretório", d)
		}
	}

	wantFiles := []string{
		filepath.Join(".luz", "project.json"),
		filepath.Join(".luz", "plugins.json"),
		filepath.Join(".luz", "styles.json"),
		".gitignore",
	}
	for _, f := range wantFiles {
		if _, err := os.Stat(filepath.Join(root, f)); err != nil {
			t.Errorf("arquivo %q ausente: %v", f, err)
		}
	}
}

func TestCreate_FailsIfWorkspaceAlreadyExists(t *testing.T) {
	root := t.TempDir()
	if _, err := Create(root, "Livro", "Autor", "pt-BR"); err != nil {
		t.Fatalf("primeiro Create: %v", err)
	}
	if _, err := Create(root, "Livro", "Autor", "pt-BR"); err == nil {
		t.Fatal("esperava erro ao criar workspace sobre um existente")
	}
}

func TestOpen_RestoresState(t *testing.T) {
	root := t.TempDir()
	if _, err := Create(root, "Meu Livro", "Autora", "pt-BR"); err != nil {
		t.Fatalf("Create: %v", err)
	}

	info, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if info.Project.Title != "Meu Livro" {
		t.Errorf("Title = %q, esperava %q", info.Project.Title, "Meu Livro")
	}
	if len(info.Project.Authors) != 1 || info.Project.Authors[0] != "Autora" {
		t.Errorf("Authors = %v", info.Project.Authors)
	}
}

func TestOpen_NotAWorkspace(t *testing.T) {
	if _, err := Open(t.TempDir()); err == nil {
		t.Fatal("esperava erro ao abrir diretório sem project.json")
	}
}

func TestProject_CorruptedJSONReturnsFriendlyError(t *testing.T) {
	root := t.TempDir()
	if _, err := Create(root, "Livro", "Autor", "pt-BR"); err != nil {
		t.Fatalf("Create: %v", err)
	}

	projectFile := filepath.Join(root, ".luz", "project.json")
	if err := os.WriteFile(projectFile, []byte("{ isso não é json"), 0o644); err != nil {
		t.Fatal(err)
	}

	w := &Workspace{Root: root}
	_, err := w.Project()
	if err == nil {
		t.Fatal("esperava erro para project.json corrompido")
	}
	var corrupted *CorruptedFileError
	if !errors.As(err, &corrupted) {
		t.Fatalf("erro = %T, esperava *CorruptedFileError", err)
	}
}
