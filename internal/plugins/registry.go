package plugins

// registry mantém o catálogo de todos os plugins conhecidos (núcleo +
// opcionais), na ordem em que aparecem na aba Extensions.
var registry []Plugin

func register(p Plugin) Plugin {
	registry = append(registry, p)
	return p
}

// All devolve o catálogo completo (núcleo + opcionais), na ordem de registro.
func All() []Plugin {
	out := make([]Plugin, len(registry))
	copy(out, registry)
	return out
}

// Core devolve apenas os módulos do núcleo.
func Core() []Plugin {
	var out []Plugin
	for _, p := range registry {
		if p.Core() {
			out = append(out, p)
		}
	}
	return out
}

// Optional devolve apenas os plugins opcionais.
func Optional() []Plugin {
	var out []Plugin
	for _, p := range registry {
		if !p.Core() {
			out = append(out, p)
		}
	}
	return out
}

// ByName procura um plugin (núcleo ou opcional) pelo nome.
func ByName(name string) (Plugin, bool) {
	for _, p := range registry {
		if p.Name() == name {
			return p, true
		}
	}
	return nil, false
}
