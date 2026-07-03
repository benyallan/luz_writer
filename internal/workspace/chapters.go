package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"luz-writer/internal/model"
)

// ListChapters lista os documentos em capitulos/ na ordem de chapterOrder,
// com o resumo usado para popular o Explorer.
func (w *Workspace) ListChapters() ([]model.ChapterMeta, error) {
	project, err := w.Project()
	if err != nil {
		return nil, err
	}

	metas := make([]model.ChapterMeta, 0, len(project.ChapterOrder))
	for _, id := range project.ChapterOrder {
		var ch model.Chapter
		if err := readJSON(w.chapterPath(id), &ch); err != nil {
			return nil, err
		}
		metas = append(metas, model.ChapterMeta{
			ID:           id,
			Title:        titleOrFallback(extractTitle(ch.Content)),
			Role:         ch.Role,
			WordCount:    countWords(ch.Content),
			HasOverrides: w.HasDocumentOverrides(id),
		})
	}
	return metas, nil
}

func titleOrFallback(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return "(sem título)"
	}
	return title
}

// Chapters devolve os documentos completos (não só o resumo), na ordem de
// chapterOrder — usado pelo pipeline de compilação (internal/build).
func (w *Workspace) Chapters() ([]model.Chapter, error) {
	project, err := w.Project()
	if err != nil {
		return nil, err
	}

	chapters := make([]model.Chapter, 0, len(project.ChapterOrder))
	for _, id := range project.ChapterOrder {
		var ch model.Chapter
		if err := readJSON(w.chapterPath(id), &ch); err != nil {
			return nil, err
		}
		chapters = append(chapters, ch)
	}
	return chapters, nil
}

// ChaptersTolerant é como Chapters, mas não aborta quando um id listado em
// chapterOrder não tem arquivo correspondente — em vez disso, devolve esse
// id em missingIDs. Usado pelo Rule Engine (R006), que precisa continuar
// validando mesmo com referências quebradas. Arquivos presentes porém
// corrompidos continuam sendo um erro genuíno (propagado normalmente).
func (w *Workspace) ChaptersTolerant() (chapters []model.Chapter, missingIDs []string, err error) {
	project, err := w.Project()
	if err != nil {
		return nil, nil, err
	}

	for _, id := range project.ChapterOrder {
		if _, statErr := os.Stat(w.chapterPath(id)); os.IsNotExist(statErr) {
			missingIDs = append(missingIDs, id)
			continue
		}
		var ch model.Chapter
		if err := readJSON(w.chapterPath(id), &ch); err != nil {
			return nil, nil, err
		}
		chapters = append(chapters, ch)
	}
	return chapters, missingIDs, nil
}

// LoadChapter devolve o conteúdo ProseMirror (campo "content") do documento
// id, como string JSON — pronto para o Tiptap consumir.
func (w *Workspace) LoadChapter(id string) (string, error) {
	var ch model.Chapter
	if err := readJSON(w.chapterPath(id), &ch); err != nil {
		return "", err
	}
	return string(ch.Content), nil
}

// SaveChapter grava um novo conteúdo ProseMirror no documento id, preservando
// luzVersion/role/language.
func (w *Workspace) SaveChapter(id string, contentJSON string) error {
	var ch model.Chapter
	if err := readJSON(w.chapterPath(id), &ch); err != nil {
		return err
	}
	if !json.Valid([]byte(contentJSON)) {
		return fmt.Errorf("conteúdo do capítulo '%s' não é um JSON válido", id)
	}
	ch.Content = json.RawMessage(contentJSON)
	return writeJSON(w.chapterPath(id), ch)
}

// CreateChapter cria um novo documento com id slug "NN-titulo", registra-o em
// chapterOrder e devolve seu resumo.
func (w *Workspace) CreateChapter(title string, role model.DocumentRole) (model.ChapterMeta, error) {
	project, err := w.Project()
	if err != nil {
		return model.ChapterMeta{}, err
	}

	id := nextChapterID(title, project.ChapterOrder)

	seedContent, err := json.Marshal(map[string]any{
		"type": "doc",
		"content": []any{
			map[string]any{
				"type":    "paragraph",
				"content": []any{map[string]any{"type": "text", "text": title}},
			},
		},
	})
	if err != nil {
		return model.ChapterMeta{}, err
	}

	ch := model.Chapter{
		LuzVersion: luzVersion,
		ID:         id,
		Role:       role,
		Language:   nil,
		Content:    seedContent,
	}
	if err := writeJSON(w.chapterPath(id), ch); err != nil {
		return model.ChapterMeta{}, err
	}

	project.ChapterOrder = append(project.ChapterOrder, id)
	if err := w.SaveProject(project); err != nil {
		return model.ChapterMeta{}, err
	}

	return model.ChapterMeta{
		ID:        id,
		Title:     titleOrFallback(title),
		Role:      role,
		WordCount: len(strings.Fields(title)),
	}, nil
}

// DeleteChapter remove o arquivo do documento e sua entrada em chapterOrder.
func (w *Workspace) DeleteChapter(id string) error {
	project, err := w.Project()
	if err != nil {
		return err
	}

	found := false
	kept := make([]string, 0, len(project.ChapterOrder))
	for _, existing := range project.ChapterOrder {
		if existing == id {
			found = true
			continue
		}
		kept = append(kept, existing)
	}
	if !found {
		return fmt.Errorf("documento '%s' não está no projeto", id)
	}

	if err := os.Remove(w.chapterPath(id)); err != nil && !os.IsNotExist(err) {
		return err
	}
	// Cascata: exclui também as sobrescritas de página do documento, se houver.
	if err := os.Remove(w.overridePath(id)); err != nil && !os.IsNotExist(err) {
		return err
	}

	project.ChapterOrder = kept
	return w.SaveProject(project)
}

// ReorderChapters substitui chapterOrder por order, desde que seja uma
// permutação do conjunto atual de documentos.
func (w *Workspace) ReorderChapters(order []string) error {
	project, err := w.Project()
	if err != nil {
		return err
	}

	current := make(map[string]bool, len(project.ChapterOrder))
	for _, id := range project.ChapterOrder {
		current[id] = true
	}
	if len(order) != len(current) {
		return fmt.Errorf("a nova ordem tem %d documento(s), esperava %d", len(order), len(current))
	}
	seen := make(map[string]bool, len(order))
	for _, id := range order {
		if !current[id] {
			return fmt.Errorf("documento '%s' não pertence ao projeto", id)
		}
		if seen[id] {
			return fmt.Errorf("documento '%s' duplicado na nova ordem", id)
		}
		seen[id] = true
	}

	project.ChapterOrder = order
	return w.SaveProject(project)
}
