package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

type FileNode struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
}

type RecentProject struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

const maxRecents = 7

func recentsFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(configDir, "luz-writer")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "recents.json"), nil
}

func (a *App) GetRecentProjects() []RecentProject {
	p, err := recentsFilePath()
	if err != nil {
		return []RecentProject{}
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return []RecentProject{}
	}
	var projects []RecentProject
	if err := json.Unmarshal(data, &projects); err != nil {
		return []RecentProject{}
	}
	return projects
}

func (a *App) AddRecentProject(projectPath string) {
	current := a.GetRecentProjects()
	entry := RecentProject{Name: filepath.Base(projectPath), Path: projectPath}

	updated := []RecentProject{entry}
	for _, r := range current {
		if r.Path != projectPath {
			updated = append(updated, r)
		}
	}
	if len(updated) > maxRecents {
		updated = updated[:maxRecents]
	}

	p, err := recentsFilePath()
	if err != nil {
		return
	}
	data, _ := json.Marshal(updated)
	os.WriteFile(p, data, 0o644)
}

func (a *App) ClearRecentProjects() {
	p, err := recentsFilePath()
	if err != nil {
		return
	}
	os.WriteFile(p, []byte("[]"), 0o644)
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OpenFolder() string {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Abrir pasta",
	})
	if err != nil {
		return ""
	}
	return path
}

func (a *App) ReadDirectory(dirPath string) []FileNode {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return []FileNode{}
	}

	var nodes []FileNode
	for _, entry := range entries {
		nodes = append(nodes, FileNode{
			Name:  entry.Name(),
			Path:  filepath.Join(dirPath, entry.Name()),
			IsDir: entry.IsDir(),
		})
	}

	// Directories first, then files, each group sorted alphabetically
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].IsDir != nodes[j].IsDir {
			return nodes[i].IsDir
		}
		return nodes[i].Name < nodes[j].Name
	})

	return nodes
}

func validateItemName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("o nome não pode estar vazio")
	}
	if strings.ContainsAny(name, `/\:*?"<>|`) {
		return fmt.Errorf("o nome contém caracteres inválidos")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("nome inválido")
	}
	return nil
}

func (a *App) CreateFile(parentPath, name string) (string, error) {
	if err := validateItemName(name); err != nil {
		return "", err
	}
	filePath := filepath.Join(parentPath, strings.TrimSpace(name))
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return "", fmt.Errorf("já existe um arquivo com esse nome")
		}
		return "", fmt.Errorf("não foi possível criar o arquivo: %w", err)
	}
	f.Close()
	return filePath, nil
}

func (a *App) CreateDirectory(parentPath, name string) (string, error) {
	if err := validateItemName(name); err != nil {
		return "", err
	}
	dirPath := filepath.Join(parentPath, strings.TrimSpace(name))
	if err := os.Mkdir(dirPath, 0o755); err != nil {
		if os.IsExist(err) {
			return "", fmt.Errorf("já existe uma pasta com esse nome")
		}
		return "", fmt.Errorf("não foi possível criar a pasta: %w", err)
	}
	return dirPath, nil
}

func (a *App) GetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

func (a *App) CreateProject(name, parentPath string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("o nome do projeto não pode estar vazio")
	}
	if strings.ContainsAny(name, `/\:*?"<>|`) {
		return "", fmt.Errorf("o nome contém caracteres inválidos")
	}
	if name == "." || name == ".." {
		return "", fmt.Errorf("nome de projeto inválido")
	}
	projectPath := filepath.Join(parentPath, name)
	if err := os.MkdirAll(projectPath, 0o755); err != nil {
		return "", fmt.Errorf("não foi possível criar o projeto: %w", err)
	}

	// Scaffold project structure
	dirs := []string{".tmp", "dist", "src", "targets"}
	for _, dir := range dirs {
		if err := os.Mkdir(filepath.Join(projectPath, dir), 0o755); err != nil {
			return "", fmt.Errorf("não foi possível criar a pasta %q: %w", dir, err)
		}
	}

	gitignore := ".tmp/\ndist/\n"
	if err := os.WriteFile(filepath.Join(projectPath, ".gitignore"), []byte(gitignore), 0o644); err != nil {
		return "", fmt.Errorf("não foi possível criar o .gitignore: %w", err)
	}

	return projectPath, nil
}

func (a *App) QuitApp() {
	runtime.Quit(a.ctx)
}
