package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"luz-writer/internal/model"
)

// ListTargets lista os targets salvos em .luz/targets/, ordenados por id.
func (w *Workspace) ListTargets() ([]model.Target, error) {
	entries, err := os.ReadDir(w.targetsDir())
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Target{}, nil
		}
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)

	targets := make([]model.Target, 0, len(names))
	for _, name := range names {
		var t model.Target
		if err := readJSON(filepath.Join(w.targetsDir(), name), &t); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	return targets, nil
}

// SaveTarget cria (se t.ID estiver vazio) ou atualiza um target.
func (w *Workspace) SaveTarget(t model.Target) (model.Target, error) {
	if t.ID == "" {
		t.ID = w.nextTargetID(t.Name)
	}
	if t.LuzVersion == 0 {
		t.LuzVersion = luzVersion
	}
	if err := writeJSON(w.targetPath(t.ID), t); err != nil {
		return model.Target{}, err
	}
	return t, nil
}

// DeleteTarget remove um target; se era o target ativo, limpa activeTarget.
func (w *Workspace) DeleteTarget(id string) error {
	if err := os.Remove(w.targetPath(id)); err != nil && !os.IsNotExist(err) {
		return err
	}

	project, err := w.Project()
	if err != nil {
		return err
	}
	if project.ActiveTarget == id {
		project.ActiveTarget = ""
		return w.SaveProject(project)
	}
	return nil
}

// SetActiveTarget torna id o target ativo do projeto.
func (w *Workspace) SetActiveTarget(id string) error {
	if _, err := os.Stat(w.targetPath(id)); err != nil {
		return fmt.Errorf("target '%s' não encontrado", id)
	}
	project, err := w.Project()
	if err != nil {
		return err
	}
	project.ActiveTarget = id
	return w.SaveProject(project)
}

// ActiveTarget devolve o target marcado como ativo no projeto; se ausente ou
// não encontrado, cai para o primeiro target salvo; se não houver nenhum,
// devolve um target embutido com os defaults dos módulos do núcleo — para o
// projeto sempre poder compilar/validar, mesmo antes do usuário criar um
// target de verdade.
func (w *Workspace) ActiveTarget() (model.Target, error) {
	project, err := w.Project()
	if err != nil {
		return model.Target{}, err
	}
	targets, err := w.ListTargets()
	if err != nil {
		return model.Target{}, err
	}
	if project.ActiveTarget != "" {
		for _, t := range targets {
			if t.ID == project.ActiveTarget {
				return t, nil
			}
		}
	}
	if len(targets) > 0 {
		return targets[0], nil
	}
	return DefaultTarget(), nil
}

// DefaultTarget é o target embutido usado quando o workspace ainda não tem
// nenhum target salvo (book, 11pt, apenas módulos do núcleo com defaults).
func DefaultTarget() model.Target {
	return model.Target{
		LuzVersion:    1,
		Name:          "Padrão",
		Kind:          model.TargetKindPrint,
		DocumentClass: model.DocumentClassBook,
		FontSize:      "11pt",
		IncludeToc:    true,
		PluginConfig:  map[string]json.RawMessage{},
	}
}

func (w *Workspace) nextTargetID(name string) string {
	base := Slugify(name)
	if _, err := os.Stat(w.targetPath(base)); os.IsNotExist(err) {
		return base
	}
	for n := 1; ; n++ {
		id := fmt.Sprintf("%s-%d", base, n)
		if _, err := os.Stat(w.targetPath(id)); os.IsNotExist(err) {
			return id
		}
	}
}
