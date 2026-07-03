package build

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// withSystemTectonicOnly desliga a preferência pelo binário embutido para o
// teste corrente, restaurando ao final — necessário porque nesta máquina
// (linux/amd64) há um binário embutido de verdade, que do contrário venceria
// qualquer configuração de PATH que os testes de fallback tentem simular.
func withSystemTectonicOnly(t *testing.T) {
	t.Helper()
	forceSystemTectonic = true
	t.Cleanup(func() { forceSystemTectonic = false })
}

func TestCheckTectonic_NotFound(t *testing.T) {
	withSystemTectonicOnly(t)
	t.Setenv("PATH", "")
	found, path := CheckTectonic()
	if found {
		t.Fatalf("esperava found=false com PATH vazio, obteve path=%q", path)
	}
}

func TestCheckTectonic_Found(t *testing.T) {
	withSystemTectonicOnly(t)
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

// TestCheckTectonic_PrefersEmbeddedOverPath cobre o comportamento default
// (seção 10): mesmo com um tectonic diferente no PATH, o binário embutido
// tem precedência — só roda em plataformas com binário embutido (ver
// embedded_linux_amd64.go / embedded_other.go).
func TestCheckTectonic_PrefersEmbeddedOverPath(t *testing.T) {
	if len(embeddedTectonic) == 0 {
		t.Skip("sem binário embutido nesta plataforma")
	}
	t.Setenv("HOME", t.TempDir())  // isola o diretório de cache do usuário
	t.Setenv("XDG_CACHE_HOME", "") // força os.UserCacheDir() a derivar de HOME

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
		t.Fatal("esperava found=true")
	}
	if path == fake {
		t.Fatalf("esperava preferir o binário embutido (extraído), obteve o do PATH: %q", path)
	}
}

func TestEnsureEmbeddedTectonic_ExtractsAndReusesCopy(t *testing.T) {
	if len(embeddedTectonic) == 0 {
		t.Skip("sem binário embutido nesta plataforma")
	}
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CACHE_HOME", "")

	path1, err := ensureEmbeddedTectonic()
	if err != nil {
		t.Fatalf("ensureEmbeddedTectonic: %v", err)
	}
	info, err := os.Stat(path1)
	if err != nil {
		t.Fatalf("binário não foi extraído: %v", err)
	}
	if info.Size() != int64(len(embeddedTectonic)) {
		t.Errorf("tamanho extraído = %d, esperava %d", info.Size(), len(embeddedTectonic))
	}

	// Segunda chamada deve reconhecer a cópia já extraída (mesmo caminho,
	// sem precisar regravar) em vez de falhar ou duplicar.
	path2, err := ensureEmbeddedTectonic()
	if err != nil {
		t.Fatalf("ensureEmbeddedTectonic (segunda chamada): %v", err)
	}
	if path1 != path2 {
		t.Errorf("caminho mudou entre chamadas: %q vs %q", path1, path2)
	}
}
