package latex

import "strings"

// BabelLangName mapeia o código BCP-47 usado nos nós/marks de idioma para o
// nome de idioma do babel (seção 6, mark luzLang). Também usado para nomear
// o idioma principal do projeto no preâmbulo (internal/build).
var BabelLangName = map[string]string{
	"en":    "english",
	"fr":    "french",
	"de":    "ngerman",
	"es":    "spanish",
	"it":    "italian",
	"la":    "latin",
	"pt-BR": "brazilian",
}

// markPriority define a ordem de aninhamento quando várias marks incidem
// sobre o mesmo texto — determinística, para a saída ser estável e testável.
var markPriority = []string{"bold", "italic", "underline", "strike", "luzInlineQuote", "luzLang", "luzCustomStyle"}

// wrapMarks aplica as marks de um nó de texto ao seu conteúdo já escapado,
// na ordem de markPriority (a primeira da lista fica mais externa).
func wrapMarks(escaped string, marks []mark, ctx *Context) string {
	byType := make(map[string]mark, len(marks))
	for _, m := range marks {
		byType[m.Type] = m
	}

	out := escaped
	// Aplica de dentro para fora: percorre markPriority ao contrário para que
	// o primeiro tipo da lista termine sendo o comando mais externo.
	for i := len(markPriority) - 1; i >= 0; i-- {
		t := markPriority[i]
		m, ok := byType[t]
		if !ok {
			continue
		}
		out = wrapMark(t, m, out, ctx)
	}
	return out
}

func wrapMark(t string, m mark, content string, ctx *Context) string {
	switch t {
	case "bold":
		return `\textbf{` + content + `}`
	case "italic":
		return `\textit{` + content + `}`
	case "underline":
		return `\uline{` + content + `}`
	case "strike":
		return `\sout{` + content + `}`
	case "luzInlineQuote":
		return `\enquote{` + content + `}`
	case "luzLang":
		// R015: com o plugin languages desativado, a marcação é ignorada na
		// exportação (texto simples, sem \foreignlanguage).
		if !ctx.Options.LanguagesEnabled {
			return content
		}
		code := attrString(m.Attrs, "lang", "")
		name, ok := BabelLangName[code]
		if !ok {
			return content
		}
		ctx.LangsUsed[code] = true
		return `\foreignlanguage{` + name + `}{` + content + `}`
	case "luzCustomStyle":
		styleID := attrString(m.Attrs, "styleId", "")
		if styleID == "" {
			return content
		}
		return `\` + customStyleCommandName(styleID) + `{` + content + `}`
	default:
		return content
	}
}

// customStyleCommandName segue a mesma convenção de
// internal/plugins.CustomStyleCommandName (duplicada de propósito — este
// pacote não deve depender de internal/plugins).
func customStyleCommandName(id string) string {
	var b strings.Builder
	for _, part := range strings.Split(id, "-") {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		b.WriteString(part[1:])
	}
	return "luzstyle" + b.String()
}
