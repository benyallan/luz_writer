package plugins

import (
	"encoding/json"

	"luz-writer/internal/model"
)

type languagesPlugin struct{}

// Languages é o portão de UI para conteúdo em idioma diverso: habilita a
// mark luzLang e o seletor de idioma do documento. Não tem configuração por
// target — a detecção de idiomas usados para o babel acontece
// automaticamente em internal/build quando este plugin está habilitado.
var Languages = register(&languagesPlugin{})

func (languagesPlugin) Name() string        { return "languages" }
func (languagesPlugin) DisplayName() string { return "Suporte Multilíngue" }
func (languagesPlugin) Description() string {
	return "Habilita marcação de trechos e documentos em outro idioma."
}
func (languagesPlugin) Core() bool          { return false }
func (languagesPlugin) DocumentScope() bool { return false }
func (languagesPlugin) Schema() model.FormSchema {
	return model.FormSchema{Fields: []model.FormField{}}
}
func (languagesPlugin) DefaultConfig() json.RawMessage { return json.RawMessage("{}") }
func (languagesPlugin) Validate(cfg json.RawMessage, ctx model.BuildContext) []model.Problem {
	return nil
}
func (languagesPlugin) Preamble(cfg json.RawMessage, ctx model.BuildContext) (string, error) {
	return "", nil
}
func (languagesPlugin) ScopedLaTeX(effectiveCfg json.RawMessage, ctx model.BuildContext) (string, string, error) {
	return "", "", nil
}
