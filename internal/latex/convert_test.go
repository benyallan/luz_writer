package latex

import (
	"strings"
	"testing"
)

var optsAllEnabled = Options{LanguagesEnabled: true, HyphenationEnabled: true}

// TestConvertContent_CanonicalExample é o teste obrigatório da seção 6: o
// capítulo de exemplo da seção 5.4 deve gerar exatamente este LaTeX.
func TestConvertContent_CanonicalExample(t *testing.T) {
	content := []byte(`{
		"type": "doc",
		"content": [
			{
				"type": "luzChapter",
				"attrs": {"numbered": true, "includeInToc": true},
				"content": [{"type": "text", "text": "Introdução"}]
			},
			{
				"type": "paragraph",
				"content": [
					{"type": "text", "text": "O conceito de "},
					{"type": "text", "marks": [{"type": "italic"}], "text": "fluxo"},
					{"type": "text", "text": " é central"},
					{"type": "luzFootnote", "attrs": {"number": 1, "text": "Ver Csikszentmihalyi, 1990."}},
					{"type": "text", "text": "."}
				]
			}
		]
	}`)

	want := "\\chapter{Introdução}\n\n" +
		"O conceito de \\textit{fluxo} é central\\footnote{Ver Csikszentmihalyi, 1990.}."

	got, problems, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	if got != want {
		t.Errorf("ConvertContent() =\n%q\nesperava\n%q", got, want)
	}
	if len(problems) != 0 {
		t.Errorf("problems inesperados: %v", problems)
	}
}

func TestConvertContent_HeadingUnnumberedNotInToc(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"luzSection","attrs":{"numbered":false,"includeInToc":false},"content":[{"type":"text","text":"Notas"}]}
	]}`)
	got, _, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	want := `\section*{Notas}`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConvertContent_HeadingUnnumberedInToc(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"luzSection","attrs":{"numbered":false,"includeInToc":true},"content":[{"type":"text","text":"Notas"}]}
	]}`)
	got, _, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	want := "\\section*{Notas}\n\\addcontentsline{toc}{section}{Notas}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConvertContent_ParagraphAlign(t *testing.T) {
	cases := map[string]string{
		"left":    "\\begin{FlushLeft}\ntexto\n\\end{FlushLeft}",
		"center":  "\\begin{Center}\ntexto\n\\end{Center}",
		"right":   "\\begin{FlushRight}\ntexto\n\\end{FlushRight}",
		"justify": "texto",
	}
	for align, want := range cases {
		content := []byte(`{"type":"doc","content":[
			{"type":"paragraph","attrs":{"align":"` + align + `"},"content":[{"type":"text","text":"texto"}]}
		]}`)
		got, _, _ := ConvertContent(content, "test", nil, optsAllEnabled)
		if got != want {
			t.Errorf("align=%s: got %q, want %q", align, got, want)
		}
	}
}

func TestConvertContent_Marks(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"paragraph","content":[
			{"type":"text","marks":[{"type":"bold"}],"text":"negrito"},
			{"type":"text","text":" "},
			{"type":"text","marks":[{"type":"underline"}],"text":"sublinhado"},
			{"type":"text","text":" "},
			{"type":"text","marks":[{"type":"strike"}],"text":"tachado"},
			{"type":"text","text":" "},
			{"type":"text","marks":[{"type":"luzInlineQuote"}],"text":"citado"}
		]}
	]}`)
	got, _, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	want := `\textbf{negrito} \uline{sublinhado} \sout{tachado} \enquote{citado}`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConvertContent_LuzLang(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"paragraph","content":[
			{"type":"text","marks":[{"type":"luzLang","attrs":{"lang":"en"}}],"text":"hello"}
		]}
	]}`)
	got, _, langs := ConvertContent(content, "test", nil, optsAllEnabled)
	want := `\foreignlanguage{english}{hello}`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if !langs["en"] {
		t.Errorf("langsUsed não contém 'en': %v", langs)
	}
}

func TestConvertContent_LuzCustomStyle(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"paragraph","content":[
			{"type":"text","marks":[{"type":"luzCustomStyle","attrs":{"styleId":"termo-estrangeiro"}}],"text":"lato sensu"}
		]}
	]}`)
	got, _, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	want := `\luzstyleTermoEstrangeiro{lato sensu}`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConvertContent_SoftHyphen(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"paragraph","content":[
			{"type":"text","text":"pseudo"},
			{"type":"luzSoftHyphen"},
			{"type":"text","text":"random"}
		]}
	]}`)
	got, _, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	want := `pseudo\-random`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConvertContent_Variable(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"paragraph","content":[
			{"type":"text","text":"Olá, "},
			{"type":"luzVariable","attrs":{"name":"protagonista"}},
			{"type":"text","text":"!"}
		]}
	]}`)
	got, problems, _ := ConvertContent(content, "test", map[string]string{"protagonista": "Ana Clara"}, optsAllEnabled)
	want := `Olá, Ana Clara!`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if len(problems) != 0 {
		t.Errorf("problems inesperados: %v", problems)
	}
}

func TestConvertContent_MissingVariableWarns(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"paragraph","content":[
			{"type":"luzVariable","attrs":{"name":"inexistente"}}
		]}
	]}`)
	got, problems, _ := ConvertContent(content, "test", map[string]string{}, optsAllEnabled)
	if got != "" {
		t.Errorf("got %q, want empty (variável ausente)", got)
	}
	if len(problems) != 1 || problems[0].Severity != "warning" {
		t.Errorf("esperava 1 warning, got %v", problems)
	}
}

func TestConvertContent_Image(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"luzImage","attrs":{"src":"imagens/grafico.png","caption":"Um gráfico","width":60}}
	]}`)
	got, _, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	want := "\\begin{figure}[htbp]\n\\centering\n\\includegraphics[width=0.60\\linewidth]{imagens/grafico.png}\n\\caption{Um gráfico}\n\\end{figure}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConvertContent_ImageWithoutCaption(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"luzImage","attrs":{"src":"imagens/grafico.png","width":80}}
	]}`)
	got, _, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	want := "\\begin{center}\n\\includegraphics[width=0.80\\linewidth]{imagens/grafico.png}\n\\end{center}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConvertContent_Lists(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"bulletList","content":[
			{"type":"listItem","content":[{"type":"paragraph","content":[{"type":"text","text":"um"}]}]},
			{"type":"listItem","content":[{"type":"paragraph","content":[{"type":"text","text":"dois"}]}]}
		]}
	]}`)
	got, _, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	want := "\\begin{itemize}\n\\item um\n\\item dois\n\\end{itemize}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConvertContent_Blockquote(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"blockquote","content":[{"type":"paragraph","content":[{"type":"text","text":"citação"}]}]}
	]}`)
	got, _, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	want := "\\begin{quote}\ncitação\n\\end{quote}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConvertContent_UnknownNodeWarnsAndSkips(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"paragraph","content":[{"type":"text","text":"antes"}]},
		{"type":"noSeiOQueEIsso"},
		{"type":"paragraph","content":[{"type":"text","text":"depois"}]}
	]}`)
	got, problems, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	if !strings.Contains(got, "antes") || !strings.Contains(got, "depois") {
		t.Errorf("esperava blocos válidos preservados, got %q", got)
	}
	if len(problems) != 1 || problems[0].Severity != "warning" {
		t.Errorf("esperava 1 warning para nó desconhecido, got %v", problems)
	}
}

func TestConvertContent_EscapesUserText(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"paragraph","content":[{"type":"text","text":"50% & R$100"}]}
	]}`)
	got, _, _ := ConvertContent(content, "test", nil, optsAllEnabled)
	want := `50\% \& R\$100`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// R015: com o plugin languages desativado, luzLang é ignorado.
func TestConvertContent_LuzLang_DisabledIgnoresMarkup(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"paragraph","content":[
			{"type":"text","marks":[{"type":"luzLang","attrs":{"lang":"en"}}],"text":"hello"}
		]}
	]}`)
	got, _, langs := ConvertContent(content, "test", nil, Options{LanguagesEnabled: false})
	if got != "hello" {
		t.Errorf("got %q, want %q (marcação ignorada)", got, "hello")
	}
	if langs["en"] {
		t.Errorf("langsUsed não deveria registrar 'en' com o plugin desativado: %v", langs)
	}
}

// R016: com o plugin hyphenation desativado, luzSoftHyphen é ignorado.
func TestConvertContent_SoftHyphen_DisabledIgnoresNode(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"paragraph","content":[
			{"type":"text","text":"pseudo"},
			{"type":"luzSoftHyphen"},
			{"type":"text","text":"random"}
		]}
	]}`)
	got, _, _ := ConvertContent(content, "test", nil, Options{HyphenationEnabled: false})
	want := `pseudorandom`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
