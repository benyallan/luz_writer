package main

import (
	"context"
	"os"
	"path/filepath"
	"sort"

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

func (a *App) QuitApp() {
	runtime.Quit(a.ctx)
}
