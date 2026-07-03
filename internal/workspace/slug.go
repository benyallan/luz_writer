package workspace

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var slugNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify normaliza s para um slug ASCII minúsculo (letras, números e
// hífens) — usado tanto para ids de capítulo quanto para nomes de arquivo
// de saída (internal/build).
func Slugify(s string) string {
	return slugify(s)
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = removeDiacritics(s)
	s = slugNonAlnum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "documento"
	}
	return s
}

func removeDiacritics(s string) string {
	var b strings.Builder
	for _, r := range norm.NFD.String(s) {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// nextChapterID gera um id "NN-slug" único dentro de existing, a partir do
// título fornecido.
func nextChapterID(title string, existing []string) string {
	base := slugify(title)
	used := make(map[string]bool, len(existing))
	for _, id := range existing {
		used[id] = true
	}
	for n := len(existing); ; n++ {
		id := fmt.Sprintf("%02d-%s", n, base)
		if !used[id] {
			return id
		}
	}
}
