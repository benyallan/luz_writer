package workspace

import "luz-writer/internal/model"

// PluginsConfig lê .luz/plugins.json.
func (w *Workspace) PluginsConfig() (model.PluginsConfig, error) {
	var cfg model.PluginsConfig
	if err := readJSON(w.pluginsPath(), &cfg); err != nil {
		return model.PluginsConfig{}, err
	}
	return cfg, nil
}

// SetPluginEnabled liga/desliga um plugin OPCIONAL em .luz/plugins.json
// (módulos do núcleo nunca aparecem nessa lista — seção 8.0).
func (w *Workspace) SetPluginEnabled(name string, enabled bool) error {
	cfg, err := w.PluginsConfig()
	if err != nil {
		return err
	}

	has := false
	kept := make([]string, 0, len(cfg.Enabled))
	for _, n := range cfg.Enabled {
		if n == name {
			has = true
			if !enabled {
				continue
			}
		}
		kept = append(kept, n)
	}
	if enabled && !has {
		kept = append(kept, name)
	}

	cfg.Enabled = kept
	return writeJSON(w.pluginsPath(), cfg)
}
