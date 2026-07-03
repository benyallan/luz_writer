package plugins

import (
	"strings"
	"testing"

	"luz-writer/internal/model"
)

func TestCustomStyleCommandName(t *testing.T) {
	if got := CustomStyleCommandName("termo-estrangeiro"); got != "luzstyleTermoEstrangeiro" {
		t.Errorf("got %q", got)
	}
	if got := CustomStyleCommandName("voz-interior"); got != "luzstyleVozInterior" {
		t.Errorf("got %q", got)
	}
}

func TestCustomStyles_PreambleGeneratesCommands(t *testing.T) {
	color := "#555555"
	ctx := model.BuildContext{Styles: []model.CustomStyle{
		{ID: "termo-estrangeiro", Name: "Termo Estrangeiro", Italic: true},
		{ID: "voz-interior", Name: "Voz Interior", Italic: true, SmallCaps: true, Color: &color},
	}}

	pre, err := CustomStyles.Preamble(CustomStyles.DefaultConfig(), ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(pre, "\\usepackage{xcolor}") {
		t.Errorf("esperava xcolor (um estilo usa cor):\n%s", pre)
	}
	if !strings.Contains(pre, "\\newcommand{\\luzstyleTermoEstrangeiro}[1]{\\textit{#1}}") {
		t.Errorf("comando do primeiro estilo ausente:\n%s", pre)
	}
	if !strings.Contains(pre, "\\textcolor[HTML]{555555}") {
		t.Errorf("cor ausente:\n%s", pre)
	}
}

func TestCustomStyles_NoStylesNoPreamble(t *testing.T) {
	pre, err := CustomStyles.Preamble(CustomStyles.DefaultConfig(), model.BuildContext{})
	if err != nil {
		t.Fatal(err)
	}
	if pre != "" {
		t.Errorf("esperava preâmbulo vazio sem estilos, got %q", pre)
	}
}
