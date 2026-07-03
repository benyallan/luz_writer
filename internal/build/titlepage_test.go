package build

import (
	"strings"
	"testing"

	"luz-writer/internal/model"
	"luz-writer/internal/workspace"
)

func TestTitleMetadata_TitleSubtitleAuthors(t *testing.T) {
	project := model.Project{
		Title:    "Meu Livro",
		Subtitle: "Uma Jornada",
		Authors:  []string{"Fulana", "Beltrano"},
	}
	got := titleMetadata(project)
	if !strings.Contains(got, "\\title{Meu Livro\\\\[0.5em]\\large Uma Jornada}") {
		t.Errorf("título/subtítulo inesperados:\n%s", got)
	}
	if !strings.Contains(got, "\\author{Fulana \\and Beltrano}") {
		t.Errorf("autores inesperados:\n%s", got)
	}
}

func TestTitleMetadata_EscapesSpecialChars(t *testing.T) {
	project := model.Project{Title: "50% & Cia", Authors: []string{"A & B"}}
	got := titleMetadata(project)
	if !strings.Contains(got, `50\% \& Cia`) {
		t.Errorf("título não escapado:\n%s", got)
	}
	if !strings.Contains(got, `A \& B`) {
		t.Errorf("autor não escapado:\n%s", got)
	}
}

func TestTitleMetadata_NoSubtitleOmitsLineBreak(t *testing.T) {
	project := model.Project{Title: "Só Título"}
	got := titleMetadata(project)
	if strings.Contains(got, "\\\\") {
		t.Errorf("não deveria ter quebra de linha sem subtítulo:\n%s", got)
	}
}

func TestTitlePageBody_IncludeTocRespected(t *testing.T) {
	withToc := titlePageBody(model.Target{IncludeToc: true})
	if !strings.Contains(withToc, "\\maketitle") || !strings.Contains(withToc, "\\tableofcontents") {
		t.Errorf("esperava maketitle+tableofcontents:\n%s", withToc)
	}

	withoutToc := titlePageBody(model.Target{IncludeToc: false})
	if !strings.Contains(withoutToc, "\\maketitle") {
		t.Errorf("esperava maketitle:\n%s", withoutToc)
	}
	if strings.Contains(withoutToc, "\\tableofcontents") {
		t.Errorf("não deveria ter tableofcontents com IncludeToc=false:\n%s", withoutToc)
	}
}

func TestBuildPreamble_IncludesTitleMetadata(t *testing.T) {
	in := BuildInputs{
		Workspace: testWorkspace(t),
		Project:   model.Project{Language: "pt-BR", Title: "Livro X", Authors: []string{"Autora"}},
		Target:    workspace.DefaultTarget(),
	}
	pre, err := BuildPreamble(in, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(pre, "\\title{Livro X}") || !strings.Contains(pre, "\\author{Autora}") {
		t.Errorf("preâmbulo deveria conter title/author:\n%s", pre)
	}
}
