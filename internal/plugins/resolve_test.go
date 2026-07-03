package plugins

import (
	"encoding/json"
	"testing"
)

func TestResolve_DefaultOnly(t *testing.T) {
	cfg, err := Resolve(Geometry, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	var got GeometryConfig
	if err := json.Unmarshal(cfg, &got); err != nil {
		t.Fatal(err)
	}
	if got != GeometryDefaults {
		t.Errorf("got %+v, want defaults %+v", got, GeometryDefaults)
	}
}

func TestResolve_TargetOverridesDefault(t *testing.T) {
	targetCfg := []byte(`{"paperWidth":"6in","paperHeight":"9in"}`)
	cfg, err := Resolve(Geometry, targetCfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	var got GeometryConfig
	if err := json.Unmarshal(cfg, &got); err != nil {
		t.Fatal(err)
	}
	if got.PaperWidth != "6in" || got.PaperHeight != "9in" {
		t.Errorf("campos do target não venceram: %+v", got)
	}
	// campos ausentes no target herdam do default.
	if got.MarginTop != GeometryDefaults.MarginTop {
		t.Errorf("MarginTop deveria herdar do default: %+v", got)
	}
}

func TestResolve_OverrideWinsOverTarget(t *testing.T) {
	targetCfg := []byte(`{"marginTop":"2.5cm","marginBottom":"2.5cm"}`)
	overrideCfg := []byte(`{"marginTop":"4cm"}`)
	cfg, err := Resolve(Geometry, targetCfg, overrideCfg)
	if err != nil {
		t.Fatal(err)
	}
	var got GeometryConfig
	if err := json.Unmarshal(cfg, &got); err != nil {
		t.Fatal(err)
	}
	if got.MarginTop != "4cm" {
		t.Errorf("override deveria vencer: MarginTop=%q", got.MarginTop)
	}
	if got.MarginBottom != "2.5cm" {
		t.Errorf("campo não sobrescrito deveria vir do target: MarginBottom=%q", got.MarginBottom)
	}
	// campo nem no target nem no override: herda do default.
	if got.Mirrored != GeometryDefaults.Mirrored {
		t.Errorf("Mirrored deveria herdar do default: %v", got.Mirrored)
	}
}

func TestResolve_EmptyOverrideDoesNotClobber(t *testing.T) {
	targetCfg := []byte(`{"paperWidth":"6in"}`)
	cfg, err := Resolve(Geometry, targetCfg, []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	var got GeometryConfig
	if err := json.Unmarshal(cfg, &got); err != nil {
		t.Fatal(err)
	}
	if got.PaperWidth != "6in" {
		t.Errorf("override vazio não deveria apagar o valor do target: %+v", got)
	}
}
