package presets

import (
	"encoding/json"

	"luz-writer/internal/model"
)

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// All são os presets embutidos no binário (seção 8.5) — pontos de partida
// para o diálogo "Novo target". Depois de copiados para .luz/targets/, não
// há vínculo com o original.
var All = []model.Target{
	{
		LuzVersion:    1,
		ID:            "amazon-kdp-6x9",
		Name:          `Amazon KDP 6"x9" (miolo)`,
		Kind:          model.TargetKindPrint,
		DocumentClass: model.DocumentClassBook,
		FontSize:      "11pt",
		IncludeToc:    true,
		PluginConfig: map[string]json.RawMessage{
			"geometry": mustJSON(map[string]any{
				"paperWidth": "6in", "paperHeight": "9in",
				"marginInner": "0.75in", "marginOuter": "0.5in",
				"marginTop": "0.75in", "marginBottom": "0.75in",
				"mirrored": true,
			}),
			"fancyhdr": mustJSON(map[string]any{
				"headerLeft": "author", "headerRight": "chapterTitle",
				"pageNumberPosition": "outer-footer",
			}),
		},
	},
	{
		LuzVersion:    1,
		ID:            "livro-a5-generico",
		Name:          "Livro A5 Genérico",
		Kind:          model.TargetKindPrint,
		DocumentClass: model.DocumentClassBook,
		FontSize:      "11pt",
		IncludeToc:    true,
		PluginConfig: map[string]json.RawMessage{
			"geometry": mustJSON(map[string]any{
				"paperWidth": "14.8cm", "paperHeight": "21cm",
				"marginInner": "2.5cm", "marginOuter": "2.5cm",
				"marginTop": "2.5cm", "marginBottom": "2.5cm",
				"mirrored": true,
			}),
		},
	},
	{
		LuzVersion:    1,
		ID:            "ebook-pdf-fluido",
		Name:          "e-Book (PDF Fluido)",
		Kind:          model.TargetKindEbook,
		DocumentClass: model.DocumentClassBook,
		FontSize:      "12pt",
		IncludeToc:    true,
		PluginConfig: map[string]json.RawMessage{
			"geometry": mustJSON(map[string]any{
				"paperWidth": "14.8cm", "paperHeight": "21cm",
				"marginInner": "1.5cm", "marginOuter": "1.5cm",
				"marginTop": "1.5cm", "marginBottom": "1.5cm",
				"mirrored": false,
			}),
		},
	},
	{
		LuzVersion:    1,
		ID:            "artigo-a4",
		Name:          "Artigo A4",
		Kind:          model.TargetKindArticle,
		DocumentClass: model.DocumentClassArticle,
		FontSize:      "12pt",
		IncludeToc:    false,
		PluginConfig: map[string]json.RawMessage{
			"geometry": mustJSON(map[string]any{
				"paperWidth": "21cm", "paperHeight": "29.7cm",
				"marginInner": "2.5cm", "marginOuter": "2.5cm",
				"marginTop": "2.5cm", "marginBottom": "2.5cm",
				"mirrored": false,
			}),
		},
	},
}
