package latex

import "strings"

// EscapeText escapa os caracteres especiais do LaTeX (tabela obrigatória da
// seção 6). Percorre a string uma única vez para nunca re-escapar as barras
// invertidas introduzidas pelas próprias substituições.
func EscapeText(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\textbackslash{}`)
		case '{':
			b.WriteString(`\{`)
		case '}':
			b.WriteString(`\}`)
		case '$':
			b.WriteString(`\$`)
		case '&':
			b.WriteString(`\&`)
		case '#':
			b.WriteString(`\#`)
		case '^':
			b.WriteString(`\^{}`)
		case '_':
			b.WriteString(`\_`)
		case '%':
			b.WriteString(`\%`)
		case '~':
			b.WriteString(`\textasciitilde{}`)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
