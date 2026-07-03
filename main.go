package main

import (
	"embed"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// workspaceFilesPrefix é o prefixo de URL usado pelo editor para exibir
// imagens importadas (seção 6, luzImage) — resolvido dinamicamente contra o
// workspace ativo, nunca embutido no binário.
const workspaceFilesPrefix = "/workspace-files/"

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "luz-writer",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: workspaceFilesHandler(app),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// workspaceFilesHandler serve arquivos de dentro do workspace ativo (por
// exemplo imagens/grafico.png) sob /workspace-files/, para que o editor
// possa exibir imagens importadas via <img src="/workspace-files/...">.
func workspaceFilesHandler(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, workspaceFilesPrefix) || app.ws == nil {
			http.NotFound(w, r)
			return
		}
		rel := strings.TrimPrefix(r.URL.Path, workspaceFilesPrefix)
		full := filepath.Join(app.ws.Root, rel)
		if !strings.HasPrefix(full, filepath.Clean(app.ws.Root)+string(filepath.Separator)) {
			http.Error(w, "caminho inválido", http.StatusForbidden)
			return
		}
		http.ServeFile(w, r, full)
	})
}
