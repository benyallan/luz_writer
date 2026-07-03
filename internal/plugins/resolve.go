package plugins

import "encoding/json"

// Resolve calcula a config efetiva de um plugin para um documento (seção
// 8.4):
//
//	efetiva = DefaultConfig()  ⊕  targetCfg (pluginConfig do target ativo)  ⊕  overrideCfg (documento)
//
// ⊕ é merge raso campo a campo: um campo definido na camada mais específica
// vence o da mais genérica; campos ausentes herdam. targetCfg e overrideCfg
// podem ser nil/vazios.
//
// O "portão" da camada 2 (plugin opcional desabilitado ⇒ ignora as camadas 3
// e 4) é responsabilidade de quem chama Resolve (internal/build), que é
// quem sabe quais plugins estão habilitados em plugins.json — Resolve só
// faz o merge.
func Resolve(p Plugin, targetCfg json.RawMessage, overrideCfg json.RawMessage) (json.RawMessage, error) {
	merged := map[string]any{}

	if err := unmarshalLayer(p.DefaultConfig(), merged); err != nil {
		return nil, err
	}
	if err := unmarshalLayer(targetCfg, merged); err != nil {
		return nil, err
	}
	if err := unmarshalLayer(overrideCfg, merged); err != nil {
		return nil, err
	}

	return json.Marshal(merged)
}

func unmarshalLayer(raw json.RawMessage, dst map[string]any) error {
	if len(raw) == 0 {
		return nil
	}
	var layer map[string]any
	if err := json.Unmarshal(raw, &layer); err != nil {
		return err
	}
	for k, v := range layer {
		dst[k] = v
	}
	return nil
}
