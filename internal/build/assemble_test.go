package build

import (
	"strings"
	"testing"

	"luz-writer/internal/model"
	"luz-writer/internal/workspace"
)

func ch(id string, role model.DocumentRole, text string) model.Chapter {
	return model.Chapter{
		LuzVersion: 1,
		ID:         id,
		Role:       role,
		Content:    []byte(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"` + text + `"}]}]}`),
	}
}

// testWorkspace cria um workspace real vazio (necessário porque AssembleBody
// consulta GetDocumentOverrides por documento) — os capítulos de teste são
// passados diretamente em memória, sem precisar existir em capitulos/.
func testWorkspace(t *testing.T) *workspace.Workspace {
	t.Helper()
	root := t.TempDir()
	if _, err := workspace.Create(root, "Livro de Teste", "Autora", "pt-BR"); err != nil {
		t.Fatalf("workspace.Create: %v", err)
	}
	return &workspace.Workspace{Root: root}
}

func TestAssembleBody_OrdersByBlockRegardlessOfInputOrder(t *testing.T) {
	in := BuildInputs{
		Workspace: testWorkspace(t),
		Project:   model.Project{Language: "pt-BR"},
		Target:    workspace.DefaultTarget(),
		Chapters: []model.Chapter{
			ch("01-intro", model.RoleChapter, "IntroTexto"),
			ch("00-dedic", model.RoleDedication, "DedicTexto"),
			ch("99-sobre", model.RoleAboutAuthor, "SobreTexto"),
		},
	}

	body, problems, _, err := AssembleBody(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Fatalf("problems inesperados: %v", problems)
	}

	frontIdx := strings.Index(body, "\\frontmatter")
	dedicIdx := strings.Index(body, "DedicTexto")
	mainIdx := strings.Index(body, "\\mainmatter")
	introIdx := strings.Index(body, "IntroTexto")
	backIdx := strings.Index(body, "\\backmatter")
	sobreIdx := strings.Index(body, "SobreTexto")

	for name, idx := range map[string]int{
		"frontmatter": frontIdx, "dedicatoria": dedicIdx, "mainmatter": mainIdx,
		"intro": introIdx, "backmatter": backIdx, "sobre": sobreIdx,
	} {
		if idx == -1 {
			t.Fatalf("%s não encontrado no corpo montado:\n%s", name, body)
		}
	}

	if !(frontIdx < dedicIdx && dedicIdx < mainIdx && mainIdx < introIdx && introIdx < backIdx && backIdx < sobreIdx) {
		t.Errorf("ordem incorreta dos blocos:\n%s", body)
	}
}

func TestAssembleBody_AppendixMarkerBeforeFirstAppendix(t *testing.T) {
	in := BuildInputs{
		Workspace: testWorkspace(t),
		Project:   model.Project{Language: "pt-BR"},
		Target:    workspace.DefaultTarget(),
		Chapters: []model.Chapter{
			ch("01-intro", model.RoleChapter, "IntroTexto"),
			ch("90-apendice-a", model.RoleAppendix, "ApendiceA"),
			ch("91-apendice-b", model.RoleAppendix, "ApendiceB"),
			ch("99-sobre", model.RoleAboutAuthor, "SobreTexto"),
		},
	}

	body, _, _, err := AssembleBody(in)
	if err != nil {
		t.Fatal(err)
	}

	appendixIdx := strings.Index(body, "\\appendix")
	aIdx := strings.Index(body, "ApendiceA")
	bIdx := strings.Index(body, "ApendiceB")
	sobreIdx := strings.Index(body, "SobreTexto")

	if appendixIdx == -1 || aIdx == -1 || bIdx == -1 || sobreIdx == -1 {
		t.Fatalf("marcadores/textos ausentes:\n%s", body)
	}
	if !(appendixIdx < aIdx && aIdx < bIdx && bIdx < sobreIdx) {
		t.Errorf("ordem incorreta ao redor de \\appendix:\n%s", body)
	}
	if strings.Count(body, "\\appendix") != 1 {
		t.Errorf("\\appendix deveria aparecer uma única vez:\n%s", body)
	}
}

func TestAssembleBody_NoAppendixNoMarker(t *testing.T) {
	in := BuildInputs{
		Workspace: testWorkspace(t),
		Project:   model.Project{Language: "pt-BR"},
		Target:    workspace.DefaultTarget(),
		Chapters: []model.Chapter{
			ch("01-intro", model.RoleChapter, "IntroTexto"),
			ch("99-sobre", model.RoleAboutAuthor, "SobreTexto"),
		},
	}
	body, _, _, err := AssembleBody(in)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(body, "\\appendix") {
		t.Errorf("não deveria haver \\appendix sem documentos desse role:\n%s", body)
	}
}

func TestAssembleBody_ExpandsProjectVariables(t *testing.T) {
	chapter := model.Chapter{
		ID:   "01-intro",
		Role: model.RoleChapter,
		Content: []byte(`{"type":"doc","content":[
			{"type":"paragraph","content":[
				{"type":"text","text":"Olá "},
				{"type":"luzVariable","attrs":{"name":"heroi"}}
			]}
		]}`),
	}
	in := BuildInputs{
		Workspace: testWorkspace(t),
		Project: model.Project{
			Language:  "pt-BR",
			Variables: []model.Variable{{Name: "heroi", Value: "Ana Clara"}},
		},
		Target:   workspace.DefaultTarget(),
		Chapters: []model.Chapter{chapter},
	}
	body, problems, _, err := AssembleBody(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Fatalf("problems inesperados: %v", problems)
	}
	if !strings.Contains(body, "Olá Ana Clara") {
		t.Errorf("variável não expandida:\n%s", body)
	}
}

func TestAssembleBody_CollectsLanguagesUsedWhenPluginEnabled(t *testing.T) {
	chapter := model.Chapter{
		ID:   "01-intro",
		Role: model.RoleChapter,
		Content: []byte(`{"type":"doc","content":[
			{"type":"paragraph","content":[
				{"type":"text","marks":[{"type":"luzLang","attrs":{"lang":"en"}}],"text":"hello"}
			]}
		]}`),
	}
	in := BuildInputs{
		Workspace:       testWorkspace(t),
		Project:         model.Project{Language: "pt-BR"},
		Target:          workspace.DefaultTarget(),
		EnabledOptional: []string{"languages"},
		Chapters:        []model.Chapter{chapter},
	}
	body, _, langs, err := AssembleBody(in)
	if err != nil {
		t.Fatal(err)
	}
	if !langs["en"] {
		t.Errorf("esperava 'en' em langsUsed: %v", langs)
	}
	if !strings.Contains(body, "\\foreignlanguage{english}{hello}") {
		t.Errorf("esperava luzLang honrado com plugin habilitado:\n%s", body)
	}
}

func TestAssembleBody_IgnoresLanguageMarkupWhenPluginDisabled(t *testing.T) {
	chapter := model.Chapter{
		ID:   "01-intro",
		Role: model.RoleChapter,
		Content: []byte(`{"type":"doc","content":[
			{"type":"paragraph","content":[
				{"type":"text","marks":[{"type":"luzLang","attrs":{"lang":"en"}}],"text":"hello"}
			]}
		]}`),
	}
	in := BuildInputs{
		Workspace: testWorkspace(t),
		Project:   model.Project{Language: "pt-BR"},
		Target:    workspace.DefaultTarget(),
		Chapters:  []model.Chapter{chapter},
	}
	body, _, langs, err := AssembleBody(in)
	if err != nil {
		t.Fatal(err)
	}
	if langs["en"] {
		t.Errorf("não deveria registrar idioma com o plugin desabilitado: %v", langs)
	}
	if strings.Contains(body, "\\foreignlanguage") {
		t.Errorf("não deveria gerar \\foreignlanguage com o plugin desabilitado:\n%s", body)
	}
}
