package build

import (
	"encoding/json"
	"strings"

	"luz-writer/internal/latex"
	"luz-writer/internal/model"
	"luz-writer/internal/plugins"
	"luz-writer/internal/workspace"
)

// BuildInputs agrupa tudo que AssembleBody/BuildPreamble precisam para uma
// compilação: o target ativo (real, seção 8.4) e a lista de plugins
// opcionais habilitados, na ordem declarada em plugins.json.
type BuildInputs struct {
	Workspace       *workspace.Workspace
	Project         model.Project
	Target          model.Target
	EnabledOptional []string
	Styles          []model.CustomStyle
	Chapters        []model.Chapter
}

// activeModules devolve os módulos ativos (núcleo primeiro, depois os
// opcionais habilitados na ordem declarada) — a mesma ordem vale para
// Preamble() e para ScopedLaTeX().
func (in BuildInputs) activeModules() []plugins.Plugin {
	out := append([]plugins.Plugin{}, plugins.Core()...)
	for _, name := range in.EnabledOptional {
		if p, ok := plugins.ByName(name); ok && !p.Core() {
			out = append(out, p)
		}
	}
	return out
}

func (in BuildInputs) isEnabled(name string) bool {
	for _, n := range in.EnabledOptional {
		if n == name {
			return true
		}
	}
	return false
}

func (in BuildInputs) buildContext() model.BuildContext {
	names := make([]string, 0)
	for _, p := range in.activeModules() {
		names = append(names, p.Name())
	}
	return model.BuildContext{Project: in.Project, Target: in.Target, ActiveModules: names, Styles: in.Styles}
}

// effectiveConfig resolve a config de um plugin para o target ativo, com uma
// sobrescrita de documento opcional (seção 8.4).
func (in BuildInputs) effectiveConfig(p plugins.Plugin, docOverrides map[string]json.RawMessage) (json.RawMessage, error) {
	targetCfg := in.Target.PluginConfig[p.Name()]
	var overrideCfg json.RawMessage
	if docOverrides != nil {
		overrideCfg = docOverrides[p.Name()]
	}
	return plugins.Resolve(p, targetCfg, overrideCfg)
}

// AssembleBody converte e agrupa todos os documentos do projeto, na ordem de
// chapterOrder, envolvendo-os pelos marcadores de bloco do livro (seção 5.5,
// passo 3 da seção 10): frontmatter → mainmatter (capítulos) → backmatter
// (apêndices e sobre o autor, com \appendix antes do primeiro apêndice), com
// o PDF externo do plugin pdfpages splicado na posição configurada.
func AssembleBody(in BuildInputs) (body string, problems []model.Problem, langsUsed map[string]bool, err error) {
	vars := make(map[string]string, len(in.Project.Variables))
	for _, v := range in.Project.Variables {
		vars[v.Name] = v.Value
	}

	var front, main, back []model.Chapter
	for _, ch := range in.Chapters {
		switch ch.Role {
		case model.RoleDedication, model.RoleEpigraph, model.RoleAcknowledgments, model.RolePreface:
			front = append(front, ch)
		case model.RoleAboutAuthor, model.RoleAppendix:
			back = append(back, ch)
		default:
			main = append(main, ch)
		}
	}

	// Não-nil mesmo sem nenhum warning: Problem trafega até o frontend via
	// JSON, onde "null" quebraria código que espera array (seção 7.2).
	problems = make([]model.Problem, 0)
	langsUsed = map[string]bool{}
	ctx := in.buildContext()

	renderGroup := func(docs []model.Chapter, withAppendixMarker bool) (string, error) {
		var b strings.Builder
		appendixStarted := false
		for _, ch := range docs {
			if withAppendixMarker && ch.Role == model.RoleAppendix && !appendixStarted {
				b.WriteString("\\appendix\n\n")
				appendixStarted = true
			}
			tex, probs, langs, rerr := in.renderChapter(ch, vars, ctx)
			if rerr != nil {
				return "", rerr
			}
			if strings.TrimSpace(tex) != "" {
				b.WriteString(tex)
				b.WriteString("\n\n")
			}
			problems = append(problems, probs...)
			for code := range langs {
				langsUsed[code] = true
			}
		}
		return b.String(), nil
	}

	frontStr, err := renderGroup(front, false)
	if err != nil {
		return "", nil, nil, err
	}
	mainStr, err := renderGroup(main, false)
	if err != nil {
		return "", nil, nil, err
	}
	backStr, err := renderGroup(back, true)
	if err != nil {
		return "", nil, nil, err
	}

	includePDF, placement, err := in.pdfpagesInjection()
	if err != nil {
		return "", nil, nil, err
	}

	var b strings.Builder
	if placement == "beforeFrontmatter" {
		b.WriteString(includePDF)
	}
	b.WriteString("\\frontmatter\n\n")
	b.WriteString(frontStr)
	if placement == "afterFrontmatter" {
		b.WriteString(includePDF)
	}
	b.WriteString("\\mainmatter\n\n")
	b.WriteString(mainStr)
	b.WriteString("\\backmatter\n\n")
	b.WriteString(backStr)
	if placement == "afterBackmatter" {
		b.WriteString(includePDF)
	}

	return b.String(), problems, langsUsed, nil
}

func (in BuildInputs) pdfpagesInjection() (command string, placement string, err error) {
	if !in.isEnabled("pdfpages") {
		return "", "", nil
	}
	cfg, err := in.effectiveConfig(plugins.Pdfpages, nil)
	if err != nil {
		return "", "", err
	}
	cmd, err := plugins.IncludePDFCommand(cfg)
	if err != nil {
		return "", "", err
	}
	if cmd == "" {
		return "", "", nil
	}
	var c plugins.PdfpagesConfig
	if err := json.Unmarshal(cfg, &c); err != nil {
		return "", "", err
	}
	return cmd, c.Placement, nil
}

// renderChapter converte um documento e aplica o ScopedLaTeX de qualquer
// plugin com DocumentScope==true que tenha sobrescrita para ele (seção 8.4),
// aninhando na ordem núcleo→opcionais declarados (o último a abrir é o
// primeiro a fechar).
func (in BuildInputs) renderChapter(ch model.Chapter, vars map[string]string, ctx model.BuildContext) (string, []model.Problem, map[string]bool, error) {
	overridesJSON, err := in.Workspace.GetDocumentOverrides(ch.ID)
	if err != nil {
		return "", nil, nil, err
	}
	var overrides map[string]json.RawMessage
	if err := json.Unmarshal([]byte(overridesJSON), &overrides); err != nil {
		return "", nil, nil, err
	}

	opts := latex.Options{
		LanguagesEnabled:   in.isEnabled("languages"),
		HyphenationEnabled: in.hyphenationEnabledFor(overrides),
	}

	tex, problems, langs := latex.RenderDocument(ch, vars, opts)

	for _, p := range in.activeModules() {
		raw, has := overrides[p.Name()]
		if !has || !p.DocumentScope() {
			continue
		}
		if !p.Core() && !in.isEnabled(p.Name()) {
			continue
		}
		effCfg, rerr := plugins.Resolve(p, in.Target.PluginConfig[p.Name()], raw)
		if rerr != nil {
			return "", nil, nil, rerr
		}
		before, after, rerr := p.ScopedLaTeX(effCfg, ctx)
		if rerr != nil {
			return "", nil, nil, rerr
		}
		tex = before + "\n" + tex + "\n" + after
	}

	return tex, problems, langs, nil
}

func (in BuildInputs) hyphenationEnabledFor(overrides map[string]json.RawMessage) bool {
	if !in.isEnabled("hyphenation") {
		return false
	}
	cfg, err := in.effectiveConfig(plugins.Hyphenation, overrides)
	if err != nil {
		return false
	}
	var c plugins.HyphenationConfig
	if err := json.Unmarshal(cfg, &c); err != nil {
		return false
	}
	return c.Mode != "off"
}
