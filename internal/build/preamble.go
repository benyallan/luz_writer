package build

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"luz-writer/internal/latex"
	"luz-writer/internal/plugins"
)

// BuildPreamble monta o preâmbulo LaTeX completo (seção 8.3): classe e
// tamanho de fonte do target ativo, módulos do núcleo (geometry, bodyText)
// com sua config efetiva, pacotes base sempre presentes, e os plugins
// opcionais habilitados (na ordem declarada em plugins.json) — hyperref
// sempre por último, como o pacote exige.
func BuildPreamble(in BuildInputs, langsUsed map[string]bool) (string, error) {
	ctx := in.buildContext()

	geomCfg, err := in.effectiveConfig(plugins.Geometry, nil)
	if err != nil {
		return "", err
	}
	var geom plugins.GeometryConfig
	if err := json.Unmarshal(geomCfg, &geom); err != nil {
		return "", err
	}

	docClass := string(in.Target.DocumentClass)
	if docClass == "" {
		docClass = "book"
	}
	fontSize := in.Target.FontSize
	if fontSize == "" {
		fontSize = "11pt"
	}
	sidedness := "oneside"
	if geom.Mirrored {
		sidedness = "twoside"
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("\\documentclass[%s,%s]{%s}\n\n", fontSize, sidedness, docClass))
	b.WriteString(babelLine(in.Project.Language, langsUsed))

	for _, p := range plugins.Core() {
		cfg := geomCfg
		if p.Name() != "geometry" {
			cfg, err = in.effectiveConfig(p, nil)
			if err != nil {
				return "", err
			}
		}
		pre, err := p.Preamble(cfg, ctx)
		if err != nil {
			return "", err
		}
		b.WriteString(pre)
	}
	b.WriteString("\n")

	b.WriteString("\\usepackage{csquotes}\n")
	b.WriteString("\\usepackage{ulem}\n")
	b.WriteString("\\usepackage{graphicx}\n\\graphicspath{{imagens/}}\n")
	b.WriteString("\\usepackage{ragged2e}\n")
	b.WriteString("\\usepackage{microtype}\n\n")

	for _, name := range in.EnabledOptional {
		p, ok := plugins.ByName(name)
		if !ok || p.Core() {
			continue
		}
		cfg, err := in.effectiveConfig(p, nil)
		if err != nil {
			return "", err
		}
		pre, err := p.Preamble(cfg, ctx)
		if err != nil {
			return "", err
		}
		b.WriteString(pre)
	}

	// hyperref é sempre o último pacote carregado, como o pacote exige.
	b.WriteString("\\usepackage[hidelinks]{hyperref}\n\n")

	b.WriteString(titleMetadata(in.Project))

	return b.String(), nil
}

func babelLine(mainLang string, langsUsed map[string]bool) string {
	mainName, ok := latex.BabelLangName[mainLang]
	if !ok {
		mainName = "brazilian"
	}

	others := make([]string, 0, len(langsUsed))
	for code := range langsUsed {
		if code == mainLang {
			continue
		}
		if name, ok := latex.BabelLangName[code]; ok {
			others = append(others, name)
		}
	}
	sort.Strings(others)

	all := append(others, mainName) // último idioma listado = idioma principal do babel
	return fmt.Sprintf("\\usepackage[%s]{babel}\n", strings.Join(all, ","))
}
