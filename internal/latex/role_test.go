package latex

import (
	"strings"
	"testing"

	"luz-writer/internal/model"
)

func chapterWithContent(t *testing.T, id string, role model.DocumentRole, contentJSON string) model.Chapter {
	t.Helper()
	return model.Chapter{
		LuzVersion: 1,
		ID:         id,
		Role:       role,
		Content:    []byte(contentJSON),
	}
}

func TestRenderDocument_ChapterIsUnwrapped(t *testing.T) {
	ch := chapterWithContent(t, "01-intro", model.RoleChapter, `{"type":"doc","content":[
		{"type":"luzChapter","attrs":{"numbered":true,"includeInToc":true},"content":[{"type":"text","text":"Introdução"}]}
	]}`)
	got, _, _ := RenderDocument(ch, nil, optsAllEnabled)
	want := `\chapter{Introdução}`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRenderDocument_Acknowledgments(t *testing.T) {
	ch := chapterWithContent(t, "00-agradecimentos", model.RoleAcknowledgments, `{"type":"doc","content":[
		{"type":"paragraph","content":[{"type":"text","text":"Obrigado a todos."}]}
	]}`)
	got, _, _ := RenderDocument(ch, nil, optsAllEnabled)
	if !strings.HasPrefix(got, "\\chapter*{Agradecimentos}") {
		t.Errorf("esperava prefixo \\chapter*{Agradecimentos}, got %q", got)
	}
	if strings.Contains(got, "addcontentsline") {
		t.Errorf("agradecimentos não deveria entrar no sumário por padrão: %q", got)
	}
}

func TestRenderDocument_Preface(t *testing.T) {
	ch := chapterWithContent(t, "00-prefacio", model.RolePreface, `{"type":"doc","content":[
		{"type":"paragraph","content":[{"type":"text","text":"Texto."}]}
	]}`)
	got, _, _ := RenderDocument(ch, nil, optsAllEnabled)
	want := "\\chapter*{Prefácio}\n\\addcontentsline{toc}{chapter}{Prefácio}\n\nTexto."
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRenderDocument_AboutAuthor(t *testing.T) {
	ch := chapterWithContent(t, "99-sobre", model.RoleAboutAuthor, `{"type":"doc","content":[
		{"type":"paragraph","content":[{"type":"text","text":"Bio."}]}
	]}`)
	got, _, _ := RenderDocument(ch, nil, optsAllEnabled)
	want := "\\chapter*{Sobre o Autor}\n\nBio."
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRenderDocument_Dedication(t *testing.T) {
	ch := chapterWithContent(t, "00-dedicatoria", model.RoleDedication, `{"type":"doc","content":[
		{"type":"paragraph","content":[{"type":"text","text":"Para minha família."}]}
	]}`)
	got, _, _ := RenderDocument(ch, nil, optsAllEnabled)
	if !strings.Contains(got, "\\thispagestyle{empty}") {
		t.Errorf("esperava \\thispagestyle{empty}: %q", got)
	}
	if !strings.Contains(got, "\\begin{flushright}") || !strings.Contains(got, "\\itshape") {
		t.Errorf("esperava flushright + itshape: %q", got)
	}
	if !strings.Contains(got, "Para minha família.") {
		t.Errorf("texto ausente: %q", got)
	}
}

func TestRenderDocument_Epigraph(t *testing.T) {
	ch := chapterWithContent(t, "00-epigrafe", model.RoleEpigraph, `{"type":"doc","content":[
		{"type":"paragraph","content":[{"type":"text","text":"A vida é breve."}]},
		{"type":"paragraph","content":[{"type":"text","text":"Sêneca"}]}
	]}`)
	got, _, _ := RenderDocument(ch, nil, optsAllEnabled)
	if !strings.Contains(got, "\\begin{center}\nA vida é breve.\n\\end{center}") {
		t.Errorf("esperava citação centralizada: %q", got)
	}
	if !strings.Contains(got, "\\begin{flushright}\n---Sêneca\n\\end{flushright}") {
		t.Errorf("esperava atribuição com travessão à direita: %q", got)
	}
}

func TestRenderDocument_DocumentLanguageWrapsInOtherlanguage(t *testing.T) {
	lang := "en"
	ch := model.Chapter{
		ID:       "05-prefacio-em-ingles",
		Role:     model.RoleChapter,
		Language: &lang,
		Content:  []byte(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Hello."}]}]}`),
	}
	got, _, langs := RenderDocument(ch, nil, optsAllEnabled)
	want := "\\begin{otherlanguage}{english}\nHello.\n\\end{otherlanguage}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if !langs["en"] {
		t.Errorf("langsUsed deveria conter 'en': %v", langs)
	}
}

func TestRenderDocument_CorruptedContent(t *testing.T) {
	ch := model.Chapter{ID: "broken", Role: model.RoleChapter, Content: []byte(`not json`)}
	_, problems, _ := RenderDocument(ch, nil, optsAllEnabled)
	if len(problems) != 1 || problems[0].Severity != "error" {
		t.Errorf("esperava 1 error para conteúdo corrompido, got %v", problems)
	}
}
