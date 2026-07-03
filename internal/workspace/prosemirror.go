package workspace

import (
	"encoding/json"
	"strings"
)

type pmNode struct {
	Type    string   `json:"type"`
	Text    string   `json:"text,omitempty"`
	Content []pmNode `json:"content,omitempty"`
}

// extractTitle deriva um título legível a partir do primeiro nó de bloco do
// conteúdo ProseMirror, concatenando o texto de seus filhos inline. Funciona
// tanto para o "paragraph" usado como placeholder nesta etapa quanto para os
// nós semânticos luzChapter/luzSection da Etapa 2 — ambos guardam o título
// como texto inline dentro do primeiro nó do documento.
func extractTitle(content json.RawMessage) string {
	var doc pmNode
	if err := json.Unmarshal(content, &doc); err != nil || len(doc.Content) == 0 {
		return ""
	}
	return collectText(doc.Content[0])
}

func collectText(n pmNode) string {
	if n.Text != "" {
		return n.Text
	}
	var b strings.Builder
	for _, c := range n.Content {
		b.WriteString(collectText(c))
	}
	return b.String()
}

// countWords conta as palavras de todo o texto do documento (para a Status
// Bar e o Explorer).
func countWords(content json.RawMessage) int {
	var doc pmNode
	if err := json.Unmarshal(content, &doc); err != nil {
		return 0
	}
	return len(strings.Fields(collectAllText(doc)))
}

func collectAllText(n pmNode) string {
	var b strings.Builder
	if n.Text != "" {
		b.WriteString(n.Text)
		b.WriteString(" ")
	}
	for _, c := range n.Content {
		b.WriteString(collectAllText(c))
	}
	return b.String()
}
