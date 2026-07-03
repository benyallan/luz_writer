package rules

import (
	"encoding/json"

	"luz-writer/internal/model"
	"luz-writer/internal/plugins"
)

// Rule é uma regra do Rule Engine (seção 9): recebe o Context e devolve os
// Problems que encontrar (nil se nenhum).
type Rule func(ctx Context) []model.Problem

// All é o registro central das regras do MVP (R001-R018), na ordem da
// seção 9. Novas regras devem ser adicionadas aqui.
var All = []Rule{
	R001, R002, R003, R004, R005, R006, R007, R008, R009,
	R010, R011, R012, R013, R014, R015, R016, R017, R018,
}

// Validate roda todas as regras registradas e concatena os Problems. Devolve
// sempre um slice não-nil (serializa como "[]", nunca "null" — o frontend
// espera um array de verdade em luz:problems e no retorno de Validate()).
func Validate(ctx Context) []model.Problem {
	problems := make([]model.Problem, 0)
	for _, rule := range All {
		problems = append(problems, rule(ctx)...)
	}
	return problems
}

// effectiveConfig resolve a config de um plugin para o target ativo do
// Context, sem sobrescrita de documento (default ⊕ target).
func effectiveConfig(ctx Context, p plugins.Plugin) (json.RawMessage, error) {
	return plugins.Resolve(p, ctx.Target.PluginConfig[p.Name()], nil)
}

// effectiveConfigWithOverride resolve com uma sobrescrita de documento
// adicional (default ⊕ target ⊕ override).
func effectiveConfigWithOverride(ctx Context, p plugins.Plugin, overrideRaw json.RawMessage) (json.RawMessage, error) {
	return plugins.Resolve(p, ctx.Target.PluginConfig[p.Name()], overrideRaw)
}

func roleLabel(r model.DocumentRole) string {
	switch r {
	case model.RoleDedication:
		return "Dedicatória"
	case model.RoleEpigraph:
		return "Epígrafe"
	case model.RoleAcknowledgments:
		return "Agradecimentos"
	case model.RolePreface:
		return "Prefácio"
	case model.RoleAboutAuthor:
		return "Sobre o Autor"
	case model.RoleAppendix:
		return "Apêndice"
	default:
		return string(r)
	}
}

func documentOverrides(ctx Context, id string) map[string]json.RawMessage {
	raw, err := ctx.Workspace.GetDocumentOverrides(id)
	if err != nil {
		return nil
	}
	var overrides map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &overrides); err != nil {
		return nil
	}
	return overrides
}
