package build

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCheckTectonic_NotFound(t *testing.T) {
	t.Setenv("PATH", "")
	found, path := CheckTectonic()
	if found {
		t.Fatalf("esperava found=false com PATH vazio, obteve path=%q", path)
	}
}

func TestCheckTectonic_Found(t *testing.T) {
	dir := t.TempDir()
	name := "tectonic"
	if runtime.GOOS == "windows" {
		name = "tectonic.exe"
	}
	fake := filepath.Join(dir, name)
	if err := os.WriteFile(fake, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	found, path := CheckTectonic()
	if !found {
		t.Fatal("esperava found=true com tectonic falso no PATH")
	}
	if path != fake {
		t.Fatalf("path = %q, esperava %q", path, fake)
	}
}
