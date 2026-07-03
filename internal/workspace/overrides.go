package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"luz-writer/internal/model"
)

// GetDocumentOverrides devolve as sobrescritas de plugins de um documento
// (seção 5.6), como string JSON — "{}" se o arquivo não existir.
func (w *Workspace) GetDocumentOverrides(id string) (string, error) {
	path := w.overridePath(id)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "{}", nil
	}

	var f model.DocumentOverrides
	if err := readJSON(path, &f); err != nil {
		return "", err
	}
	if len(f.Overrides) == 0 {
		return "{}", nil
	}
	b, err := json.Marshal(f.Overrides)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// SaveDocumentOverrides grava as sobrescritas de um documento. overridesJSON
// vazio/"{}" apaga o arquivo (criação sob demanda — seção 5.6).
func (w *Workspace) SaveDocumentOverrides(id string, overridesJSON string) error {
	if overridesJSON == "" {
		overridesJSON = "{}"
	}
	var overrides map[string]json.RawMessage
	if err := json.Unmarshal([]byte(overridesJSON), &overrides); err != nil {
		return fmt.Errorf("sobrescritas do documento '%s' não são um JSON válido: %w", id, err)
	}

	path := w.overridePath(id)
	if len(overrides) == 0 {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(w.overridesDir(), 0o755); err != nil {
		return err
	}
	return writeJSON(path, model.DocumentOverrides{
		LuzVersion: luzVersion,
		DocumentID: id,
		Overrides:  overrides,
	})
}

// HasDocumentOverrides indica se o documento id tem sobrescritas gravadas —
// usado para o badge ⚙ do Explorer.
func (w *Workspace) HasDocumentOverrides(id string) bool {
	_, err := os.Stat(w.overridePath(id))
	return err == nil
}

// OrphanedOverrideIDs lista os arquivos em .luz/overrides/ cujo id não está
// em knownIDs (tipicamente project.ChapterOrder) — usado pela regra R014.
func (w *Workspace) OrphanedOverrideIDs(knownIDs []string) ([]string, error) {
	entries, err := os.ReadDir(w.overridesDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	known := make(map[string]bool, len(knownIDs))
	for _, id := range knownIDs {
		known[id] = true
	}

	var orphans []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		if !known[id] {
			orphans = append(orphans, id)
		}
	}
	return orphans, nil
}
