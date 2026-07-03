package build

import (
	"encoding/json"
	"strings"
	"testing"

	"luz-writer/internal/model"
	"luz-writer/internal/workspace"
)

func TestBuildPreamble_BaseStructure(t *testing.T) {
	in := BuildInputs{
		Workspace: testWorkspace(t),
		Project:   model.Project{Language: "pt-BR"},
		Target:    workspace.DefaultTarget(),
	}
	pre, err := BuildPreamble(in, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}

	mustContain := []string{
		`\documentclass[11pt,twoside]{book}`,
		`\usepackage[paperwidth=14.8cm,paperheight=21cm`,
		`\usepackage{csquotes}`,
		`\usepackage{ulem}`,
		`\usepackage{graphicx}`,
		`\graphicspath{{imagens/}}`,
		`\usepackage{ragged2e}`,
		`\usepackage{microtype}`,
		`\usepackage[hidelinks]{hyperref}`,
	}
	for _, s := range mustContain {
		if !strings.Contains(pre, s) {
			t.Errorf("preâmbulo deveria conter %q:\n%s", s, pre)
		}
	}
}

func TestBuildPreamble_HyperrefLoadedLast(t *testing.T) {
	in := BuildInputs{
		Workspace:       testWorkspace(t),
		Project:         model.Project{Language: "pt-BR"},
		Target:          workspace.DefaultTarget(),
		EnabledOptional: []string{"fancyhdr"},
	}
	pre, err := BuildPreamble(in, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}
	lastPackageIdx := strings.LastIndex(pre, `\usepackage`)
	hyperrefIdx := strings.Index(pre, `\usepackage[hidelinks]{hyperref}`)
	if lastPackageIdx != hyperrefIdx {
		t.Errorf("hyperref deveria ser o último \\usepackage mesmo com plugins opcionais ativos:\n%s", pre)
	}
}

func TestBuildPreamble_BabelMainLanguageLast(t *testing.T) {
	in := BuildInputs{Workspace: testWorkspace(t), Project: model.Project{Language: "pt-BR"}, Target: workspace.DefaultTarget()}
	pre, err := BuildPreamble(in, map[string]bool{"en": true})
	if err != nil {
		t.Fatal(err)
	}
	want := `\usepackage[english,brazilian]{babel}`
	if !strings.Contains(pre, want) {
		t.Errorf("esperava %q no preâmbulo:\n%s", want, pre)
	}
}

func TestBuildPreamble_BabelOnlyMainLanguage(t *testing.T) {
	in := BuildInputs{Workspace: testWorkspace(t), Project: model.Project{Language: "pt-BR"}, Target: workspace.DefaultTarget()}
	pre, err := BuildPreamble(in, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}
	want := `\usepackage[brazilian]{babel}`
	if !strings.Contains(pre, want) {
		t.Errorf("esperava %q no preâmbulo:\n%s", want, pre)
	}
}

func TestBuildPreamble_OptionalPluginContributesPackage(t *testing.T) {
	in := BuildInputs{
		Workspace:       testWorkspace(t),
		Project:         model.Project{Language: "pt-BR", Authors: []string{"Autora"}},
		Target:          workspace.DefaultTarget(),
		EnabledOptional: []string{"fancyhdr"},
	}
	pre, err := BuildPreamble(in, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(pre, "\\usepackage{fancyhdr}") {
		t.Errorf("esperava fancyhdr no preâmbulo (plugin habilitado):\n%s", pre)
	}
}

func TestBuildPreamble_DisabledOptionalPluginNotIncluded(t *testing.T) {
	in := BuildInputs{Workspace: testWorkspace(t), Project: model.Project{Language: "pt-BR"}, Target: workspace.DefaultTarget()}
	pre, err := BuildPreamble(in, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(pre, "fancyhdr") {
		t.Errorf("fancyhdr não deveria aparecer sem estar habilitado:\n%s", pre)
	}
}

func TestBuildPreamble_MirroredFalseUsesOneside(t *testing.T) {
	mirroredFalse, err := json.Marshal(map[string]any{"mirrored": false})
	if err != nil {
		t.Fatal(err)
	}
	target := workspace.DefaultTarget()
	target.PluginConfig = map[string]json.RawMessage{"geometry": mirroredFalse}
	in := BuildInputs{Workspace: testWorkspace(t), Project: model.Project{Language: "pt-BR"}, Target: target}
	pre, err := BuildPreamble(in, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(pre, "oneside") {
		t.Errorf("esperava oneside na classe do documento com mirrored=false:\n%s", pre)
	}
	if !strings.Contains(pre, "left=") || !strings.Contains(pre, "right=") {
		t.Errorf("esperava left/right na geometry com mirrored=false:\n%s", pre)
	}
}
