package workspace

import "path/filepath"

const (
	dirLuz       = ".luz"
	dirTargets   = "targets"
	dirOverrides = "overrides"
	dirChapters  = "capitulos"
	dirImages    = "imagens"
	dirAttach    = "anexos"
	dirDist      = "dist"
	dirTmp       = ".tmp"

	fileProject   = "project.json"
	filePlugins   = "plugins.json"
	fileStyles    = "styles.json"
	fileGitignore = ".gitignore"
)

func (w *Workspace) luzDir() string       { return filepath.Join(w.Root, dirLuz) }
func (w *Workspace) targetsDir() string   { return filepath.Join(w.luzDir(), dirTargets) }
func (w *Workspace) overridesDir() string { return filepath.Join(w.luzDir(), dirOverrides) }
func (w *Workspace) chaptersDir() string  { return filepath.Join(w.Root, dirChapters) }
func (w *Workspace) imagesDir() string    { return filepath.Join(w.Root, dirImages) }
func (w *Workspace) attachDir() string    { return filepath.Join(w.Root, dirAttach) }
func (w *Workspace) distDir() string      { return filepath.Join(w.Root, dirDist) }
func (w *Workspace) tmpDir() string       { return filepath.Join(w.Root, dirTmp) }

func (w *Workspace) projectPath() string { return filepath.Join(w.luzDir(), fileProject) }
func (w *Workspace) pluginsPath() string { return filepath.Join(w.luzDir(), filePlugins) }
func (w *Workspace) stylesPath() string  { return filepath.Join(w.luzDir(), fileStyles) }

func (w *Workspace) chapterPath(id string) string {
	return filepath.Join(w.chaptersDir(), id+".json")
}

func (w *Workspace) overridePath(id string) string {
	return filepath.Join(w.overridesDir(), id+".json")
}

func (w *Workspace) targetPath(id string) string {
	return filepath.Join(w.targetsDir(), id+".json")
}
