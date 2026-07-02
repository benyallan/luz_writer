package latex

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Node represents a ProseMirror/Tiptap document node.
type Node struct {
	Type    string                 `json:"type"`
	Content []Node                 `json:"content,omitempty"`
	Text    string                 `json:"text,omitempty"`
	Marks   []Mark                 `json:"marks,omitempty"`
	Attrs   map[string]interface{} `json:"attrs,omitempty"`
}

// Mark represents inline formatting applied to a text node.
type Mark struct {
	Type  string                 `json:"type"`
	Attrs map[string]interface{} `json:"attrs,omitempty"`
}

// Parse converts a Tiptap/ProseMirror JSON document string to LaTeX source.
func Parse(jsonContent string) (string, error) {
	var doc Node
	if err := json.Unmarshal([]byte(jsonContent), &doc); err != nil {
		return "", fmt.Errorf("invalid document JSON: %w", err)
	}
	if doc.Type != "doc" {
		return "", fmt.Errorf("expected root node type \"doc\", got %q", doc.Type)
	}

	var sb strings.Builder
	for _, node := range doc.Content {
		renderBlock(&sb, node)
	}
	return sb.String(), nil
}

func renderBlock(sb *strings.Builder, node Node) {
	switch node.Type {
	case "paragraph":
		renderInline(sb, node.Content)
		sb.WriteString("\n\n")

	case "heading":
		level := 1
		if l, ok := node.Attrs["level"].(float64); ok {
			level = int(l)
		}
		sb.WriteString(headingCommand(level))
		sb.WriteString("{")
		renderInline(sb, node.Content)
		sb.WriteString("}\n\n")

	case "bulletList":
		sb.WriteString("\\begin{itemize}\n")
		for _, item := range node.Content {
			renderListItem(sb, item)
		}
		sb.WriteString("\\end{itemize}\n\n")

	case "orderedList":
		sb.WriteString("\\begin{enumerate}\n")
		for _, item := range node.Content {
			renderListItem(sb, item)
		}
		sb.WriteString("\\end{enumerate}\n\n")

	case "blockquote":
		sb.WriteString("\\begin{quote}\n")
		for _, child := range node.Content {
			renderBlock(sb, child)
		}
		sb.WriteString("\\end{quote}\n\n")

	case "horizontalRule":
		sb.WriteString("\\hrulefill\n\n")
	}
}

func renderListItem(sb *strings.Builder, item Node) {
	firstParagraph := true
	for _, child := range item.Content {
		switch child.Type {
		case "paragraph":
			if firstParagraph {
				sb.WriteString("  \\item ")
				firstParagraph = false
			} else {
				sb.WriteString("\n  ")
			}
			renderInline(sb, child.Content)
			sb.WriteString("\n")
		default:
			// nested list or other block
			renderBlock(sb, child)
		}
	}
	if firstParagraph {
		sb.WriteString("  \\item\n")
	}
}

func renderInline(sb *strings.Builder, nodes []Node) {
	for _, node := range nodes {
		switch node.Type {
		case "text":
			text := applyMarks(escapeLatex(node.Text), node.Marks)
			sb.WriteString(text)
		case "hardBreak":
			sb.WriteString("\\\\\n")
		}
	}
}

func applyMarks(text string, marks []Mark) string {
	for _, mark := range marks {
		switch mark.Type {
		case "bold":
			text = "\\textbf{" + text + "}"
		case "italic":
			text = "\\textit{" + text + "}"
		case "strike":
			text = "\\sout{" + text + "}"
		}
	}
	return text
}

func headingCommand(level int) string {
	switch level {
	case 1:
		return "\\chapter"
	case 2:
		return "\\section"
	case 3:
		return "\\subsection"
	default:
		return "\\subsubsection"
	}
}

// escapeLatex escapes LaTeX special characters in plain text.
// strings.NewReplacer performs a single left-to-right scan and never
// re-processes replacement strings, so the backslash must come first.
var latexEscaper = strings.NewReplacer(
	`\`, `\textbackslash{}`,
	`&`, `\&`,
	`%`, `\%`,
	`$`, `\$`,
	`#`, `\#`,
	`_`, `\_`,
	`{`, `\{`,
	`}`, `\}`,
	`~`, `\textasciitilde{}`,
	`^`, `\textasciicircum{}`,
)

func escapeLatex(s string) string {
	return latexEscaper.Replace(s)
}
