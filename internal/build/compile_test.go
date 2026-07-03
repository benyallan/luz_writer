package build

import (
	"encoding/base64"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"luz-writer/internal/model"
	"luz-writer/internal/workspace"
)

// pngBase64 é um PNG 1x1 válido (mesmo usado nos testes E2E do frontend),
// suficiente para o graphicx/xdvipdfmx do Tectonic processar de verdade.
const pngBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

func requireTectonic(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("tectonic"); err != nil {
		t.Skip("tectonic não encontrado no PATH — pulando teste de compilação real")
	}
}

// setupCompileWorkspace cria, via internal/workspace, um projeto com
// dedicatória, dois capítulos (cobrindo negrito, itálico, sublinhado,
// parágrafo centralizado, citação inline, trecho em inglês, hífen sugerido,
// variável usada duas vezes, blockquote, lista, nota de rodapé e imagem com
// legenda) e "Sobre o Autor" — o cenário do critério de aceite da Etapa 3.
func setupCompileWorkspace(t *testing.T) *workspace.Workspace {
	t.Helper()
	root := t.TempDir()

	info, err := workspace.Create(root, "Livro de Teste", "Autora de Teste", "pt-BR")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	w := &workspace.Workspace{Root: root}

	// Grava as variáveis antes de criar qualquer capítulo: CreateChapter lê e
	// regrava project.json a cada chamada (para atualizar chapterOrder), então
	// salvar um project.go stale por cima depois apagaria o chapterOrder.
	project := info.Project
	project.Variables = []model.Variable{{Name: "cidade", Value: "Iguape"}}
	must(t, w.SaveProject(project))

	if _, err := os.Stat(filepath.Join(root, "imagens")); err != nil {
		t.Fatalf("imagens/ deveria existir: %v", err)
	}
	imgPath := filepath.Join(root, "imagens", "grafico.png")
	pngBytes, err := base64.StdEncoding.DecodeString(pngBase64)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(imgPath, pngBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	dedication, err := w.CreateChapter("Dedicatória", model.RoleDedication)
	if err != nil {
		t.Fatalf("CreateChapter dedicatória: %v", err)
	}
	must(t, w.SaveChapter(dedication.ID, `{"type":"doc","content":[
		{"type":"paragraph","content":[
			{"type":"text","text":"Para "},
			{"type":"luzVariable","attrs":{"name":"cidade"}},
			{"type":"text","text":"."}
		]}
	]}`))

	ch1, err := w.CreateChapter("Introdução", model.RoleChapter)
	if err != nil {
		t.Fatalf("CreateChapter capítulo 1: %v", err)
	}
	must(t, w.SaveChapter(ch1.ID, `{"type":"doc","content":[
		{"type":"luzChapter","attrs":{"numbered":true,"includeInToc":true},"content":[{"type":"text","text":"Introdução"}]},
		{"type":"paragraph","content":[
			{"type":"text","marks":[{"type":"bold"}],"text":"negrito"},
			{"type":"text","text":" "},
			{"type":"text","marks":[{"type":"italic"}],"text":"itálico"},
			{"type":"text","text":" "},
			{"type":"text","marks":[{"type":"underline"}],"text":"sublinhado"}
		]},
		{"type":"paragraph","attrs":{"align":"center"},"content":[{"type":"text","text":"parágrafo centralizado"}]},
		{"type":"paragraph","content":[
			{"type":"text","marks":[{"type":"luzInlineQuote"}],"text":"uma citação"},
			{"type":"text","text":" e um trecho em "},
			{"type":"text","marks":[{"type":"luzLang","attrs":{"lang":"en"}}],"text":"English"},
			{"type":"text","text":", pseudo"},
			{"type":"luzSoftHyphen"},
			{"type":"text","text":"random, com a variável "},
			{"type":"luzVariable","attrs":{"name":"cidade"}},
			{"type":"text","text":" citada aqui."},
			{"type":"luzFootnote","attrs":{"number":1,"text":"Nota de rodapé de teste."}}
		]},
		{"type":"blockquote","content":[{"type":"paragraph","content":[{"type":"text","text":"uma citação em bloco"}]}]},
		{"type":"bulletList","content":[
			{"type":"listItem","content":[{"type":"paragraph","content":[{"type":"text","text":"item um"}]}]},
			{"type":"listItem","content":[{"type":"paragraph","content":[{"type":"text","text":"item dois"}]}]}
		]},
		{"type":"luzImage","attrs":{"src":"imagens/grafico.png","caption":"Um gráfico de teste","width":60}}
	]}`))

	ch2, err := w.CreateChapter("Metodologia", model.RoleChapter)
	if err != nil {
		t.Fatalf("CreateChapter capítulo 2: %v", err)
	}
	must(t, w.SaveChapter(ch2.ID, `{"type":"doc","content":[
		{"type":"luzChapter","attrs":{"numbered":true,"includeInToc":true},"content":[{"type":"text","text":"Metodologia"}]},
		{"type":"paragraph","content":[
			{"type":"text","text":"A cidade de novo: "},
			{"type":"luzVariable","attrs":{"name":"cidade"}},
			{"type":"text","text":"."}
		]}
	]}`))

	about, err := w.CreateChapter("Sobre o Autor", model.RoleAboutAuthor)
	if err != nil {
		t.Fatalf("CreateChapter sobre o autor: %v", err)
	}
	must(t, w.SaveChapter(about.ID, `{"type":"doc","content":[
		{"type":"paragraph","content":[{"type":"text","text":"Biografia da autora."}]}
	]}`))

	return w
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompile_ProducesValidPDFWithCorrectOrder(t *testing.T) {
	requireTectonic(t)
	w := setupCompileWorkspace(t)

	var stages []string
	result, err := Compile(w, func(p model.BuildProgress) {
		stages = append(stages, p.Stage)
	})
	if err != nil {
		t.Fatalf("Compile retornou erro Go inesperado: %v", err)
	}
	if !result.Success {
		t.Fatalf("Compile falhou: problems=%v logTail=%s", result.Problems, result.LogTail)
	}
	if result.OutputPath == "" {
		t.Fatal("OutputPath vazio")
	}
	if _, err := os.Stat(result.OutputPath); err != nil {
		t.Fatalf("PDF não encontrado em %s: %v", result.OutputPath, err)
	}
	if len(stages) == 0 {
		t.Error("esperava eventos de progresso")
	}

	texPath := filepath.Join(w.Root, ".tmp", "main.tex")
	texBytes, err := os.ReadFile(texPath)
	if err != nil {
		t.Fatalf("main.tex não encontrado: %v", err)
	}
	tex := string(texBytes)

	dedicIdx := indexOf(t, tex, "Para Iguape.")
	introIdx := indexOf(t, tex, "\\chapter{Introdução}")
	metodIdx := indexOf(t, tex, "\\chapter{Metodologia}")
	sobreIdx := indexOf(t, tex, "\\chapter*{Sobre o Autor}")

	if !(dedicIdx < introIdx && introIdx < metodIdx && metodIdx < sobreIdx) {
		t.Errorf("ordem incorreta no .tex montado (dedicatória antes, sobre o autor depois dos capítulos):\n%s", tex)
	}
}

func indexOf(t *testing.T, haystack, needle string) int {
	t.Helper()
	idx := -1
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			idx = i
			break
		}
	}
	if idx == -1 {
		t.Fatalf("%q não encontrado em:\n%s", needle, haystack)
	}
	return idx
}

func TestCompile_TectonicFailureReturnsLogTailWithoutError(t *testing.T) {
	requireTectonic(t)
	root := t.TempDir()
	if _, err := workspace.Create(root, "Livro Quebrado", "Autor", "pt-BR"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	w := &workspace.Workspace{Root: root}

	// Uma imagem ausente seria pega antes pelo Rule Engine (R008), sem chegar
	// a rodar o Tectonic — para exercitar uma falha real do compilador,
	// usamos um arquivo que existe (passa R008) mas não é uma imagem válida.
	if err := os.MkdirAll(filepath.Join(root, "imagens"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "imagens", "corrompida.png"), []byte("isto não é um PNG"), 0o644); err != nil {
		t.Fatal(err)
	}

	ch, err := w.CreateChapter("Capítulo", model.RoleChapter)
	if err != nil {
		t.Fatalf("CreateChapter: %v", err)
	}
	must(t, w.SaveChapter(ch.ID, `{"type":"doc","content":[
		{"type":"luzImage","attrs":{"src":"imagens/corrompida.png","width":50}}
	]}`))

	result, err := Compile(w, nil)
	if err != nil {
		t.Fatalf("Compile não deveria retornar erro Go: %v", err)
	}
	if result.Success {
		t.Fatal("esperava falha de compilação")
	}
	if result.LogTail == "" {
		t.Error("esperava LogTail preenchido")
	}
}
