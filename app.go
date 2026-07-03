package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"luz-writer/internal/build"
	"luz-writer/internal/model"
	"luz-writer/internal/plugins"
	"luz-writer/internal/plugins/presets"
	"luz-writer/internal/rules"
	"luz-writer/internal/workspace"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context

	tectonicFound bool
	tectonicPath  string

	ws *workspace.Workspace
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.tectonicFound, a.tectonicPath = build.CheckTectonic()
	if a.tectonicFound {
		log.Printf("tectonic encontrado em %s", a.tectonicPath)
	} else {
		log.Println("tectonic não encontrado no PATH — exportação ficará desabilitada")
	}
}

// requireWorkspace garante que um workspace está aberto antes de qualquer
// operação de projeto/capítulos.
func (a *App) requireWorkspace() (*workspace.Workspace, error) {
	if a.ws == nil {
		return nil, errors.New("nenhum workspace aberto")
	}
	return a.ws, nil
}

// --- Workspace ---------------------------------------------------------

// PickDirectory abre um diálogo nativo de pasta e devolve o caminho
// escolhido, sem tentar abri-lo como workspace (usado pelo fluxo de criação
// de projeto para escolher o local).
func (a *App) PickDirectory() (string, error) {
	path, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Escolha a pasta do projeto",
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// OpenWorkspaceDialog abre um diálogo nativo de pasta e abre o workspace
// escolhido.
func (a *App) OpenWorkspaceDialog() (model.WorkspaceInfo, error) {
	path, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Abrir workspace do Luz Writer",
	})
	if err != nil {
		return model.WorkspaceInfo{}, err
	}
	if path == "" {
		return model.WorkspaceInfo{}, errors.New("nenhuma pasta selecionada")
	}
	return a.OpenWorkspace(path)
}

// OpenWorkspace abre o workspace em path sem diálogo — usado para restaurar
// o último projeto aberto ao reiniciar o app.
func (a *App) OpenWorkspace(path string) (model.WorkspaceInfo, error) {
	info, err := workspace.Open(path)
	if err != nil {
		return model.WorkspaceInfo{}, err
	}
	a.ws = &workspace.Workspace{Root: path}
	a.revalidate()
	return info, nil
}

// CreateWorkspace cria um novo workspace em path e o torna o workspace ativo.
func (a *App) CreateWorkspace(path string, title string, author string, language string) (model.WorkspaceInfo, error) {
	info, err := workspace.Create(path, title, author, language)
	if err != nil {
		return model.WorkspaceInfo{}, err
	}
	a.ws = &workspace.Workspace{Root: path}
	a.revalidate()
	return info, nil
}

// GetProject lê o project.json do workspace ativo.
func (a *App) GetProject() (model.Project, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return model.Project{}, err
	}
	return w.Project()
}

// SaveProject grava o project.json do workspace ativo.
func (a *App) SaveProject(p model.Project) error {
	w, err := a.requireWorkspace()
	if err != nil {
		return err
	}
	if err := w.SaveProject(p); err != nil {
		return err
	}
	a.revalidate()
	return nil
}

// --- Capítulos -----------------------------------------------------------

// ListChapters lista os documentos do workspace ativo, na ordem de exibição.
func (a *App) ListChapters() ([]model.ChapterMeta, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return nil, err
	}
	return w.ListChapters()
}

// LoadChapter devolve o conteúdo ProseMirror (JSON) do documento id.
func (a *App) LoadChapter(id string) (string, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return "", err
	}
	return w.LoadChapter(id)
}

// SaveChapter grava um novo conteúdo ProseMirror no documento id.
func (a *App) SaveChapter(id string, contentJSON string) error {
	w, err := a.requireWorkspace()
	if err != nil {
		return err
	}
	if err := w.SaveChapter(id, contentJSON); err != nil {
		return err
	}
	a.revalidate()
	return nil
}

// CreateChapter cria um novo documento com o role informado (seção 5.5).
func (a *App) CreateChapter(title string, role string) (model.ChapterMeta, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return model.ChapterMeta{}, err
	}
	docRole := model.DocumentRole(role)
	if !isValidRole(docRole) {
		return model.ChapterMeta{}, fmt.Errorf("papel de documento inválido: '%s'", role)
	}
	meta, err := w.CreateChapter(title, docRole)
	if err != nil {
		return model.ChapterMeta{}, err
	}
	a.revalidate()
	return meta, nil
}

// DeleteChapter remove o documento id do workspace ativo.
func (a *App) DeleteChapter(id string) error {
	w, err := a.requireWorkspace()
	if err != nil {
		return err
	}
	if err := w.DeleteChapter(id); err != nil {
		return err
	}
	a.revalidate()
	return nil
}

// ReorderChapters substitui a ordem de exibição/compilação dos documentos.
func (a *App) ReorderChapters(order []string) error {
	w, err := a.requireWorkspace()
	if err != nil {
		return err
	}
	if err := w.ReorderChapters(order); err != nil {
		return err
	}
	a.revalidate()
	return nil
}

// --- Assets --------------------------------------------------------------

// ImportImage abre um diálogo nativo, copia a imagem escolhida para
// imagens/ e devolve o caminho relativo (para o atributo src do nó luzImage).
func (a *App) ImportImage() (string, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return "", err
	}
	src, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Importar imagem",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Imagens (*.png, *.jpg, *.jpeg, *.gif, *.webp)", Pattern: "*.png;*.jpg;*.jpeg;*.gif;*.webp"},
		},
	})
	if err != nil {
		return "", err
	}
	if src == "" {
		return "", errors.New("nenhuma imagem selecionada")
	}
	return w.ImportImage(src)
}

// --- Targets ---------------------------------------------------------------

// ListTargets lista os targets salvos no workspace ativo.
func (a *App) ListTargets() ([]model.Target, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return nil, err
	}
	return w.ListTargets()
}

// ListTargetPresets devolve os presets embutidos no binário (seção 8.5),
// pontos de partida para o diálogo "Novo target".
func (a *App) ListTargetPresets() ([]model.Target, error) {
	return presets.All, nil
}

// SaveTarget cria (t.ID vazio) ou atualiza um target do workspace ativo.
func (a *App) SaveTarget(t model.Target) error {
	w, err := a.requireWorkspace()
	if err != nil {
		return err
	}
	if _, err := w.SaveTarget(t); err != nil {
		return err
	}
	a.revalidate()
	return nil
}

// DeleteTarget remove um target do workspace ativo.
func (a *App) DeleteTarget(id string) error {
	w, err := a.requireWorkspace()
	if err != nil {
		return err
	}
	if err := w.DeleteTarget(id); err != nil {
		return err
	}
	a.revalidate()
	return nil
}

// SetActiveTarget torna id o target ativo do projeto.
func (a *App) SetActiveTarget(id string) error {
	w, err := a.requireWorkspace()
	if err != nil {
		return err
	}
	if err := w.SetActiveTarget(id); err != nil {
		return err
	}
	a.revalidate()
	return nil
}

// --- Plugins -----------------------------------------------------------

// ListAvailablePlugins devolve o catálogo embutido de módulos do núcleo e
// plugins opcionais (seção 8.0), para a aba Extensions. Enabled reflete o
// plugins.json do workspace ativo.
func (a *App) ListAvailablePlugins() ([]model.PluginManifest, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return nil, err
	}
	pluginsCfg, err := w.PluginsConfig()
	if err != nil {
		return nil, err
	}
	enabled := make(map[string]bool, len(pluginsCfg.Enabled))
	for _, name := range pluginsCfg.Enabled {
		enabled[name] = true
	}

	all := plugins.All()
	out := make([]model.PluginManifest, 0, len(all))
	for _, p := range all {
		out = append(out, model.PluginManifest{
			Name:          p.Name(),
			DisplayName:   p.DisplayName(),
			Description:   p.Description(),
			Core:          p.Core(),
			DocumentScope: p.DocumentScope(),
			Schema:        p.Schema(),
			Enabled:       p.Core() || enabled[p.Name()],
		})
	}
	return out, nil
}

// SetPluginEnabled liga/desliga um plugin opcional do workspace ativo.
func (a *App) SetPluginEnabled(name string, enabled bool) error {
	w, err := a.requireWorkspace()
	if err != nil {
		return err
	}
	if err := w.SetPluginEnabled(name, enabled); err != nil {
		return err
	}
	a.revalidate()
	return nil
}

// --- Estilos Personalizados (seção 5.7) -------------------------------------

// ListStyles lê os estilos personalizados do workspace ativo.
func (a *App) ListStyles() ([]model.CustomStyle, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return nil, err
	}
	return w.ListStyles()
}

// SaveStyles substitui o conjunto inteiro de estilos personalizados.
func (a *App) SaveStyles(styles []model.CustomStyle) error {
	w, err := a.requireWorkspace()
	if err != nil {
		return err
	}
	if err := w.SaveStyles(styles); err != nil {
		return err
	}
	a.revalidate()
	return nil
}

// --- Configurações por Página (seção 5.6) -----------------------------------

// GetDocumentOverrides lê as sobrescritas de um documento ("{}" se não houver).
func (a *App) GetDocumentOverrides(id string) (string, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return "", err
	}
	return w.GetDocumentOverrides(id)
}

// SaveDocumentOverrides grava as sobrescritas de um documento ("{}" apaga o
// arquivo).
func (a *App) SaveDocumentOverrides(id string, overridesJSON string) error {
	w, err := a.requireWorkspace()
	if err != nil {
		return err
	}
	if err := w.SaveDocumentOverrides(id, overridesJSON); err != nil {
		return err
	}
	a.revalidate()
	return nil
}

// ImportAttachmentPDF abre um diálogo nativo, copia o PDF escolhido para
// anexos/ e devolve o caminho relativo (usado pelo plugin pdfpages).
func (a *App) ImportAttachmentPDF() (string, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return "", err
	}
	src, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Importar PDF",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "PDF (*.pdf)", Pattern: "*.pdf"},
		},
	})
	if err != nil {
		return "", err
	}
	if src == "" {
		return "", errors.New("nenhum arquivo selecionado")
	}
	return w.ImportAttachment(src)
}

// --- Validação e Build -----------------------------------------------------

// Validate roda o Rule Engine (seção 9) contra o workspace ativo e emite o
// evento luz:problems.
func (a *App) Validate() ([]model.Problem, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return nil, err
	}
	problems, err := a.validate(w)
	if err != nil {
		return nil, err
	}
	wailsRuntime.EventsEmit(a.ctx, "luz:problems", problems)
	return problems, nil
}

func (a *App) validate(w *workspace.Workspace) ([]model.Problem, error) {
	ctx, err := rules.NewContext(w)
	if err != nil {
		return nil, err
	}
	return rules.Validate(ctx), nil
}

// revalidate roda o Rule Engine e emite luz:problems, sem retornar erro para
// o chamador — usada após salvar projeto/capítulo/target/plugins/overrides
// (seção 7.2), para não quebrar a operação principal se a validação falhar.
func (a *App) revalidate() {
	if a.ws == nil {
		return
	}
	problems, err := a.validate(a.ws)
	if err != nil {
		log.Printf("revalidação falhou: %v", err)
		return
	}
	wailsRuntime.EventsEmit(a.ctx, "luz:problems", problems)
}

// Compile executa o pipeline de compilação (seção 10), emitindo eventos de
// progresso (luz:build:progress) e o evento final (luz:build:done).
func (a *App) Compile() (model.BuildResult, error) {
	w, err := a.requireWorkspace()
	if err != nil {
		return model.BuildResult{}, err
	}

	result, err := build.Compile(w, func(p model.BuildProgress) {
		wailsRuntime.EventsEmit(a.ctx, "luz:build:progress", p)
	})
	if err != nil {
		return model.BuildResult{}, err
	}

	wailsRuntime.EventsEmit(a.ctx, "luz:build:done", result)
	wailsRuntime.EventsEmit(a.ctx, "luz:problems", result.Problems)
	return result, nil
}

func isValidRole(r model.DocumentRole) bool {
	switch r {
	case model.RoleChapter, model.RoleDedication, model.RoleEpigraph,
		model.RoleAcknowledgments, model.RolePreface, model.RoleAboutAuthor, model.RoleAppendix:
		return true
	default:
		return false
	}
}
