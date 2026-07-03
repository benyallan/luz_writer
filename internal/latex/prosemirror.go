package latex

import "encoding/json"

// node é a representação mínima de um nó ProseMirror necessária para a
// conversão (seção 6). Vive só neste pacote — internal/workspace tem sua
// própria versão minimalista para extrair título/contagem de palavras.
type node struct {
	Type    string         `json:"type"`
	Text    string         `json:"text,omitempty"`
	Marks   []mark         `json:"marks,omitempty"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Content []node         `json:"content,omitempty"`
}

type mark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

func parseDoc(content json.RawMessage) (node, error) {
	var n node
	if err := json.Unmarshal(content, &n); err != nil {
		return node{}, err
	}
	return n, nil
}

func attrString(attrs map[string]any, key, def string) string {
	if v, ok := attrs[key].(string); ok {
		return v
	}
	return def
}

func attrBool(attrs map[string]any, key string, def bool) bool {
	if v, ok := attrs[key].(bool); ok {
		return v
	}
	return def
}

func attrNumber(attrs map[string]any, key string, def float64) float64 {
	if v, ok := attrs[key].(float64); ok {
		return v
	}
	return def
}
