package workspace

import "luz-writer/internal/model"

// ListStyles lê .luz/styles.json (seção 5.7).
func (w *Workspace) ListStyles() ([]model.CustomStyle, error) {
	var f model.StylesFile
	if err := readJSON(w.stylesPath(), &f); err != nil {
		return nil, err
	}
	return f.Styles, nil
}

// SaveStyles substitui o conjunto inteiro de estilos em .luz/styles.json.
func (w *Workspace) SaveStyles(styles []model.CustomStyle) error {
	return writeJSON(w.stylesPath(), model.StylesFile{LuzVersion: luzVersion, Styles: styles})
}
