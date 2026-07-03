package latex

import (
	"encoding/json"
	"fmt"
	"strings"

	"luz-writer/internal/model"
)

// Options controla comportamentos de conversão que dependem de plugins
// habilitados (seção 8.4) — o conversor permanece "puro", apenas obedecendo
// às flags que internal/build calcula a partir de plugins.json.
type Options struct {
	LanguagesEnabled   bool // plugin "languages": honra luzLang e o language do documento
	HyphenationEnabled bool // plugin "hyphenation" (modo != "off"): honra luzSoftHyphen
}

// ConvertContent converte um documento ProseMirror completo (campo "content"
// de capitulos/<id>.json) para LaTeX. docID identifica a origem dos Problems
// gerados; vars é o mapa nome→valor das Variáveis do Projeto (seção 5.1).
func ConvertContent(content json.RawMessage, docID string, vars map[string]string, opts Options) (tex string, problems []model.Problem, langsUsed map[string]bool) {
	doc, err := parseDoc(content)
	if err != nil {
		return "", []model.Problem{{
			Severity: "error",
			Code:     "LATEX",
			Message:  fmt.Sprintf("Não foi possível interpretar o conteúdo do documento '%s'.", docID),
			Source:   "chapter:" + docID,
		}}, nil
	}

	ctx := newContext(docID, vars, opts)
	tex = convertBlocks(doc.Content, ctx)
	return tex, ctx.Problems, ctx.LangsUsed
}

// Context carrega o estado compartilhado durante a conversão de um documento.
type Context struct {
	DocumentID string
	Variables  map[string]string // name -> value já resolvido de project.json
	Options    Options
	LangsUsed  map[string]bool
	Problems   []model.Problem
}

func newContext(docID string, vars map[string]string, opts Options) *Context {
	return &Context{
		DocumentID: docID,
		Variables:  vars,
		Options:    opts,
		LangsUsed:  map[string]bool{},
		// Problems começa como slice vazio (não nil): Problem trafega até o
		// frontend via JSON, onde "null" quebraria código que espera array.
		Problems: []model.Problem{},
	}
}

func (c *Context) warnf(format string, args ...any) {
	c.Problems = append(c.Problems, model.Problem{
		Severity: "warning",
		Code:     "LATEX",
		Message:  fmt.Sprintf(format, args...),
		Source:   "chapter:" + c.DocumentID,
	})
}

// convertBlocks converte uma sequência de nós de bloco, unindo-os com linha
// em branco. Nós que geram string vazia (ex.: desconhecidos, ignorados) não
// deixam blocos em branco extras.
func convertBlocks(nodes []node, ctx *Context) string {
	parts := make([]string, 0, len(nodes))
	for _, n := range nodes {
		s := convertBlock(n, ctx)
		if strings.TrimSpace(s) != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, "\n\n")
}

func convertBlock(n node, ctx *Context) string {
	switch n.Type {
	case "luzChapter":
		return convertHeading(n, ctx, "chapter")
	case "luzSection":
		return convertHeading(n, ctx, "section")
	case "luzSubsection":
		return convertHeading(n, ctx, "subsection")
	case "paragraph":
		return convertParagraph(n, ctx)
	case "blockquote":
		inner := convertBlocks(n.Content, ctx)
		return "\\begin{quote}\n" + inner + "\n\\end{quote}"
	case "bulletList":
		return convertList(n, ctx, "itemize")
	case "orderedList":
		return convertList(n, ctx, "enumerate")
	case "luzImage":
		return convertImage(n, ctx)
	default:
		ctx.warnf("Nó desconhecido '%s' foi ignorado na exportação.", n.Type)
		return ""
	}
}

func convertHeading(n node, ctx *Context, cmd string) string {
	numbered := attrBool(n.Attrs, "numbered", true)
	includeInToc := attrBool(n.Attrs, "includeInToc", true)
	title := convertInline(n.Content, ctx)

	if numbered {
		return fmt.Sprintf(`\%s{%s}`, cmd, title)
	}
	out := fmt.Sprintf(`\%s*{%s}`, cmd, title)
	if includeInToc {
		out += fmt.Sprintf("\n\\addcontentsline{toc}{%s}{%s}", cmd, title)
	}
	return out
}

func convertParagraph(n node, ctx *Context) string {
	align := attrString(n.Attrs, "align", "justify")
	text := convertInline(n.Content, ctx)
	switch align {
	case "left":
		return "\\begin{FlushLeft}\n" + text + "\n\\end{FlushLeft}"
	case "center":
		return "\\begin{Center}\n" + text + "\n\\end{Center}"
	case "right":
		return "\\begin{FlushRight}\n" + text + "\n\\end{FlushRight}"
	default:
		return text
	}
}

func convertList(n node, ctx *Context, env string) string {
	var b strings.Builder
	b.WriteString("\\begin{" + env + "}\n")
	for _, item := range n.Content {
		if item.Type != "listItem" {
			ctx.warnf("Item de lista inesperado '%s' foi ignorado.", item.Type)
			continue
		}
		inner := strings.TrimSpace(convertBlocks(item.Content, ctx))
		b.WriteString("\\item " + inner + "\n")
	}
	b.WriteString("\\end{" + env + "}")
	return b.String()
}

func convertImage(n node, ctx *Context) string {
	src := attrString(n.Attrs, "src", "")
	caption := attrString(n.Attrs, "caption", "")
	width := attrNumber(n.Attrs, "width", 80)
	widthTex := fmt.Sprintf(`%.2f\linewidth`, width/100)

	if src == "" {
		ctx.warnf("Imagem sem arquivo definido foi ignorada.")
		return ""
	}

	if caption != "" {
		return fmt.Sprintf("\\begin{figure}[htbp]\n\\centering\n\\includegraphics[width=%s]{%s}\n\\caption{%s}\n\\end{figure}",
			widthTex, src, EscapeText(caption))
	}
	return fmt.Sprintf("\\begin{center}\n\\includegraphics[width=%s]{%s}\n\\end{center}", widthTex, src)
}

func convertInline(nodes []node, ctx *Context) string {
	var b strings.Builder
	for _, n := range nodes {
		b.WriteString(convertInlineNode(n, ctx))
	}
	return b.String()
}

func convertInlineNode(n node, ctx *Context) string {
	switch n.Type {
	case "text":
		return wrapMarks(EscapeText(n.Text), n.Marks, ctx)
	case "luzFootnote":
		text := attrString(n.Attrs, "text", "")
		return `\footnote{` + EscapeText(text) + `}`
	case "luzSoftHyphen":
		// R016: com o plugin hyphenation desativado (ou modo "off"), pontos
		// de hifenização sugeridos são ignorados na exportação.
		if !ctx.Options.HyphenationEnabled {
			return ""
		}
		return `\-`
	case "luzVariable":
		name := attrString(n.Attrs, "name", "")
		value, ok := ctx.Variables[name]
		if !ok {
			ctx.warnf("A variável '%s' usada no documento não existe mais no projeto.", name)
			return ""
		}
		return EscapeText(value)
	default:
		ctx.warnf("Nó inline desconhecido '%s' foi ignorado na exportação.", n.Type)
		return ""
	}
}
