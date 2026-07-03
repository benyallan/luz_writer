package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"luz-writer/internal/model"
)

const luzVersion = 1

const gitignoreContent = ".tmp/\ndist/\n"

// Workspace representa um projeto Luz Writer aberto em disco, centralizado
// na pasta raiz Root (seção 4 da spec).
type Workspace struct {
	Root string
}

// Create cria a estrutura de um novo workspace em path (seção 4) e grava o
// project.json inicial.
func Create(path, title, author, language string) (model.WorkspaceInfo, error) {
	if path == "" {
		return model.WorkspaceInfo{}, errors.New("caminho do workspace não pode ser vazio")
	}
	w := &Workspace{Root: path}

	if _, err := os.Stat(w.projectPath()); err == nil {
		return model.WorkspaceInfo{}, fmt.Errorf("já existe um workspace Luz Writer em '%s'", path)
	}

	dirs := []string{
		w.Root,
		w.luzDir(),
		w.targetsDir(),
		w.overridesDir(),
		w.chaptersDir(),
		w.imagesDir(),
		w.attachDir(),
		w.distDir(),
		w.tmpDir(),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return model.WorkspaceInfo{}, fmt.Errorf("não foi possível criar '%s': %w", d, err)
		}
	}

	authors := []string{}
	if author != "" {
		authors = []string{author}
	}
	project := model.Project{
		LuzVersion:   luzVersion,
		Title:        title,
		Subtitle:     "",
		Authors:      authors,
		Language:     language,
		ChapterOrder: []string{},
		ActiveTarget: "",
		Variables:    []model.Variable{},
	}
	if err := writeJSON(w.projectPath(), project); err != nil {
		return model.WorkspaceInfo{}, err
	}
	if err := writeJSON(w.pluginsPath(), model.PluginsConfig{LuzVersion: luzVersion, Enabled: []string{}}); err != nil {
		return model.WorkspaceInfo{}, err
	}
	if err := writeJSON(w.stylesPath(), model.StylesFile{LuzVersion: luzVersion, Styles: []model.CustomStyle{}}); err != nil {
		return model.WorkspaceInfo{}, err
	}
	if err := os.WriteFile(filepath.Join(w.Root, fileGitignore), []byte(gitignoreContent), 0o644); err != nil {
		return model.WorkspaceInfo{}, err
	}

	return model.WorkspaceInfo{Path: w.Root, Project: project}, nil
}

// Open lê um workspace existente em path e valida sua estrutura mínima
// (a presença de .luz/project.json).
func Open(path string) (model.WorkspaceInfo, error) {
	w := &Workspace{Root: path}
	project, err := w.Project()
	if err != nil {
		if os.IsNotExist(err) {
			return model.WorkspaceInfo{}, fmt.Errorf("'%s' não é um workspace Luz Writer válido (project.json não encontrado)", path)
		}
		return model.WorkspaceInfo{}, err
	}
	return model.WorkspaceInfo{Path: w.Root, Project: project}, nil
}
