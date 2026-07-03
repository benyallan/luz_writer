package model

// Variable é um par nome/valor reutilizável no texto (seção 5.1).
// Expandida como texto simples na exportação — nunca vira \newcommand.
type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Project é o conteúdo de .luz/project.json (seção 5.1).
type Project struct {
	LuzVersion   int        `json:"luzVersion"`
	Title        string     `json:"title"`
	Subtitle     string     `json:"subtitle"`
	Authors      []string   `json:"authors"`
	Language     string     `json:"language"`
	ChapterOrder []string   `json:"chapterOrder"`
	ActiveTarget string     `json:"activeTarget"`
	Variables    []Variable `json:"variables"`
}

// WorkspaceInfo é o retorno de OpenWorkspaceDialog/CreateWorkspace (seção 7.1).
type WorkspaceInfo struct {
	Path    string  `json:"path"`
	Project Project `json:"project"`
}
