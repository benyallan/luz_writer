package latex_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"luz-writer/internal/latex"
)

func doc(nodes ...string) string {
	return `{"type":"doc","content":[` + join(nodes...) + `]}`
}

func join(parts ...string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ","
		}
		result += p
	}
	return result
}

// jsonStr returns the JSON encoding of s (quoted and escaped).
func jsonStr(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func para(text string) string {
	return `{"type":"paragraph","content":[{"type":"text","text":` + jsonStr(text) + `}]}`
}

func markedText(text, mark string) string {
	return `{"type":"text","text":` + jsonStr(text) + `,"marks":[{"type":"` + mark + `"}]}`
}

// ── Paragraph ────────────────────────────────────────────────────────────────

func TestParseParagraph(t *testing.T) {
	got, err := latex.Parse(doc(para("Hello world")))
	assertNoError(t, err)
	assertEqual(t, "Hello world\n\n", got)
}

func TestParseEmptyDoc(t *testing.T) {
	got, err := latex.Parse(`{"type":"doc","content":[]}`)
	assertNoError(t, err)
	assertEqual(t, "", got)
}

// ── Marks ────────────────────────────────────────────────────────────────────

func TestParseBold(t *testing.T) {
	node := `{"type":"paragraph","content":[` + markedText("strong", "bold") + `]}`
	got, err := latex.Parse(doc(node))
	assertNoError(t, err)
	assertEqual(t, `\textbf{strong}`+"\n\n", got)
}

func TestParseItalic(t *testing.T) {
	node := `{"type":"paragraph","content":[` + markedText("slanted", "italic") + `]}`
	got, err := latex.Parse(doc(node))
	assertNoError(t, err)
	assertEqual(t, `\textit{slanted}`+"\n\n", got)
}

func TestParseStrike(t *testing.T) {
	node := `{"type":"paragraph","content":[` + markedText("deleted", "strike") + `]}`
	got, err := latex.Parse(doc(node))
	assertNoError(t, err)
	assertEqual(t, `\sout{deleted}`+"\n\n", got)
}

func TestParseBoldItalicCombined(t *testing.T) {
	node := `{"type":"paragraph","content":[{"type":"text","text":"bi","marks":[{"type":"bold"},{"type":"italic"}]}]}`
	got, err := latex.Parse(doc(node))
	assertNoError(t, err)
	// bold applied first, then italic wraps it
	assertEqual(t, `\textit{\textbf{bi}}`+"\n\n", got)
}

// ── Headings ─────────────────────────────────────────────────────────────────

func TestParseHeadings(t *testing.T) {
	cases := []struct {
		level int
		want  string
	}{
		{1, `\chapter{My Title}` + "\n\n"},
		{2, `\section{My Title}` + "\n\n"},
		{3, `\subsection{My Title}` + "\n\n"},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("H%d", tc.level), func(t *testing.T) {
			node := fmt.Sprintf(
				`{"type":"heading","attrs":{"level":%d},"content":[{"type":"text","text":"My Title"}]}`,
				tc.level,
			)
			got, err := latex.Parse(doc(node))
			assertNoError(t, err)
			assertEqual(t, tc.want, got)
		})
	}
}

// ── Lists ────────────────────────────────────────────────────────────────────

func TestParseBulletList(t *testing.T) {
	input := doc(`{"type":"bulletList","content":[
		{"type":"listItem","content":[{"type":"paragraph","content":[{"type":"text","text":"Alpha"}]}]},
		{"type":"listItem","content":[{"type":"paragraph","content":[{"type":"text","text":"Beta"}]}]}
	]}`)
	got, err := latex.Parse(input)
	assertNoError(t, err)
	want := "\\begin{itemize}\n  \\item Alpha\n  \\item Beta\n\\end{itemize}\n\n"
	assertEqual(t, want, got)
}

func TestParseOrderedList(t *testing.T) {
	input := doc(`{"type":"orderedList","content":[
		{"type":"listItem","content":[{"type":"paragraph","content":[{"type":"text","text":"First"}]}]},
		{"type":"listItem","content":[{"type":"paragraph","content":[{"type":"text","text":"Second"}]}]}
	]}`)
	got, err := latex.Parse(input)
	assertNoError(t, err)
	want := "\\begin{enumerate}\n  \\item First\n  \\item Second\n\\end{enumerate}\n\n"
	assertEqual(t, want, got)
}

// ── Blockquote ───────────────────────────────────────────────────────────────

func TestParseBlockquote(t *testing.T) {
	input := doc(`{"type":"blockquote","content":[` + para("To be or not to be") + `]}`)
	got, err := latex.Parse(input)
	assertNoError(t, err)
	want := "\\begin{quote}\nTo be or not to be\n\n\\end{quote}\n\n"
	assertEqual(t, want, got)
}

// ── HorizontalRule ───────────────────────────────────────────────────────────

func TestParseHorizontalRule(t *testing.T) {
	got, err := latex.Parse(doc(`{"type":"horizontalRule"}`))
	assertNoError(t, err)
	assertEqual(t, "\\hrulefill\n\n", got)
}

// ── HardBreak ────────────────────────────────────────────────────────────────

func TestParseHardBreak(t *testing.T) {
	input := doc(`{"type":"paragraph","content":[
		{"type":"text","text":"Line one"},
		{"type":"hardBreak"},
		{"type":"text","text":"Line two"}
	]}`)
	got, err := latex.Parse(input)
	assertNoError(t, err)
	assertEqual(t, "Line one\\\\\nLine two\n\n", got)
}

// ── LaTeX escaping ───────────────────────────────────────────────────────────

func TestParseSpecialChars(t *testing.T) {
	cases := []struct{ in, want string }{
		{`100% done`, `100\% done`},
		{`price $5`, `price \$5`},
		{`a & b`, `a \& b`},
		{`item #1`, `item \#1`},
		{`foo_bar`, `foo\_bar`},
		{`{value}`, `\{value\}`},
		{`x^2`, `x\textasciicircum{}2`},
		{`~home`, `\textasciitilde{}home`},
		{`back\slash`, `back\textbackslash{}slash`},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, err := latex.Parse(doc(para(tc.in)))
			assertNoError(t, err)
			assertEqual(t, tc.want+"\n\n", got)
		})
	}
}

// ── Error cases ──────────────────────────────────────────────────────────────

func TestParseInvalidJSON(t *testing.T) {
	_, err := latex.Parse(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParseWrongRootType(t *testing.T) {
	_, err := latex.Parse(`{"type":"paragraph"}`)
	if err == nil {
		t.Fatal("expected error for wrong root type, got nil")
	}
}

// ── Mixed document ───────────────────────────────────────────────────────────

func TestParseMixedDocument(t *testing.T) {
	input := doc(
		`{"type":"heading","attrs":{"level":1},"content":[{"type":"text","text":"Chapter One"}]}`,
		para("Introduction text."),
		`{"type":"bulletList","content":[
			{"type":"listItem","content":[{"type":"paragraph","content":[{"type":"text","text":"Point A"}]}]}
		]}`,
	)
	got, err := latex.Parse(input)
	assertNoError(t, err)
	want := "\\chapter{Chapter One}\n\nIntroduction text.\n\n\\begin{itemize}\n  \\item Point A\n\\end{itemize}\n\n"
	assertEqual(t, want, got)
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertEqual(t *testing.T, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("\nwant: %q\n got: %q", want, got)
	}
}
