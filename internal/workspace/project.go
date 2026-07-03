package workspace

import "luz-writer/internal/model"

// Project lê .luz/project.json.
func (w *Workspace) Project() (model.Project, error) {
	var p model.Project
	if err := readJSON(w.projectPath(), &p); err != nil {
		return model.Project{}, err
	}
	return p, nil
}

// SaveProject grava .luz/project.json.
func (w *Workspace) SaveProject(p model.Project) error {
	return writeJSON(w.projectPath(), p)
}
