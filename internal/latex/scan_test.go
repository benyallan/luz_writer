package latex

import "testing"

func TestScanContent_DetectsEverything(t *testing.T) {
	content := []byte(`{"type":"doc","content":[
		{"type":"luzChapter","attrs":{"numbered":true,"includeInToc":true},"content":[{"type":"text","text":"Título"}]},
		{"type":"paragraph","content":[
			{"type":"text","text":"Olá "},
			{"type":"text","marks":[{"type":"luzLang","attrs":{"lang":"en"}}],"text":"world"},
			{"type":"luzSoftHyphen"},
			{"type":"luzVariable","attrs":{"name":"heroi"}},
			{"type":"text","marks":[{"type":"luzCustomStyle","attrs":{"styleId":"termo-estrangeiro"}}],"text":"lato sensu"}
		]},
		{"type":"luzImage","attrs":{"src":"imagens/grafico.png"}}
	]}`)

	u, err := ScanContent(content)
	if err != nil {
		t.Fatal(err)
	}
	if !u.HasLuzChapter {
		t.Error("esperava HasLuzChapter=true")
	}
	if len(u.ImageSrcs) != 1 || u.ImageSrcs[0] != "imagens/grafico.png" {
		t.Errorf("ImageSrcs = %v", u.ImageSrcs)
	}
	if len(u.VariableNames) != 1 || u.VariableNames[0] != "heroi" {
		t.Errorf("VariableNames = %v", u.VariableNames)
	}
	if len(u.LangCodes) != 1 || u.LangCodes[0] != "en" {
		t.Errorf("LangCodes = %v", u.LangCodes)
	}
	if !u.HasSoftHyphen {
		t.Error("esperava HasSoftHyphen=true")
	}
	if len(u.CustomStyleIDs) != 1 || u.CustomStyleIDs[0] != "termo-estrangeiro" {
		t.Errorf("CustomStyleIDs = %v", u.CustomStyleIDs)
	}
	if u.IsBlank() {
		t.Error("não deveria estar em branco")
	}
}

func TestScanContent_BlankDocument(t *testing.T) {
	content := []byte(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"   "}]}]}`)
	u, err := ScanContent(content)
	if err != nil {
		t.Fatal(err)
	}
	if !u.IsBlank() {
		t.Error("esperava documento em branco")
	}
}

func TestScanContent_EmptyDocument(t *testing.T) {
	content := []byte(`{"type":"doc","content":[]}`)
	u, err := ScanContent(content)
	if err != nil {
		t.Fatal(err)
	}
	if !u.IsBlank() {
		t.Error("documento sem nenhum bloco deveria estar em branco")
	}
}
