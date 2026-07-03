package latex

import (
	"encoding/json"
	"strings"
)

// ContentUsage resume o que aparece em um documento ProseMirror — usado pelo
// Rule Engine (internal/rules) para as regras que dependem do conteúdo
// (R002, R005, R008, R015-R018), sem duplicar a árvore de nós desse pacote.
type ContentUsage struct {
	HasLuzChapter  bool
	ImageSrcs      []string
	VariableNames  []string
	LangCodes      []string
	HasSoftHyphen  bool
	CustomStyleIDs []string
	PlainText      string
}

// ScanContent varre o campo "content" de um documento (capitulos/<id>.json).
func ScanContent(content json.RawMessage) (ContentUsage, error) {
	doc, err := parseDoc(content)
	if err != nil {
		return ContentUsage{}, err
	}
	var u ContentUsage
	scanNode(doc, &u)
	return u, nil
}

func scanNode(n node, u *ContentUsage) {
	switch n.Type {
	case "luzChapter":
		u.HasLuzChapter = true
	case "luzImage":
		u.ImageSrcs = append(u.ImageSrcs, attrString(n.Attrs, "src", ""))
	case "luzVariable":
		u.VariableNames = append(u.VariableNames, attrString(n.Attrs, "name", ""))
	case "luzSoftHyphen":
		u.HasSoftHyphen = true
	case "text":
		u.PlainText += n.Text
		for _, m := range n.Marks {
			switch m.Type {
			case "luzLang":
				if code := attrString(m.Attrs, "lang", ""); code != "" {
					u.LangCodes = append(u.LangCodes, code)
				}
			case "luzCustomStyle":
				if id := attrString(m.Attrs, "styleId", ""); id != "" {
					u.CustomStyleIDs = append(u.CustomStyleIDs, id)
				}
			}
		}
	}

	for _, c := range n.Content {
		scanNode(c, u)
	}
}

// IsBlank indica se o documento não tem texto visível (seção 9, R005).
func (u ContentUsage) IsBlank() bool {
	return strings.TrimSpace(u.PlainText) == ""
}
