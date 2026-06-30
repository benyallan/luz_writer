package main

import (
	"context"
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
		if entry.Name()[0] == '.' {
			continue
		}
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
	return projectPath, nil
}

func (a *App) QuitApp() {
	runtime.Quit(a.ctx)
}
