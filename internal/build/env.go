package build

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// forceSystemTectonic desliga a preferência pelo binário embutido — só é
// setado por testes que precisam exercitar o fallback via PATH sem depender
// de rodar (ou não) numa plataforma com binário embutido.
var forceSystemTectonic = false

// CheckTectonic garante que um binário tectonic utilizável está disponível:
// prefere a cópia embutida no app, extraída para o diretório de cache do
// usuário na primeira execução (para o Luz Writer não depender de instalação
// manual — seção 10), e cai para o PATH do sistema como reserva em
// plataformas sem binário embutido ou se a extração falhar por algum motivo.
// Sua ausência nunca deve travar o app — apenas desabilita a exportação,
// sinalizada ao usuário como um Problem persistente.
func CheckTectonic() (found bool, path string) {
	if !forceSystemTectonic {
		if p, err := ensureEmbeddedTectonic(); err == nil {
			return true, p
		}
	}
	p, err := exec.LookPath("tectonic")
	if err != nil {
		return false, ""
	}
	return true, p
}

// ensureEmbeddedTectonic extrai o binário embutido (se houver um para esta
// plataforma) para <cache do usuário>/luz-writer/bin/tectonic[.exe],
// pulando a extração se uma cópia da mesma versão já existe lá.
func ensureEmbeddedTectonic() (string, error) {
	if len(embeddedTectonic) == 0 {
		return "", errors.New("nenhum binário tectonic embutido para esta plataforma")
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	binDir := filepath.Join(cacheDir, "luz-writer", "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return "", err
	}

	name := "tectonic"
	if runtime.GOOS == "windows" {
		name = "tectonic.exe"
	}
	dest := filepath.Join(binDir, name)
	versionMarker := dest + ".version"

	if data, err := os.ReadFile(versionMarker); err == nil && string(data) == embeddedTectonicVersion {
		if info, err := os.Stat(dest); err == nil && info.Size() == int64(len(embeddedTectonic)) {
			return dest, nil
		}
	}

	if err := os.WriteFile(dest, embeddedTectonic, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(versionMarker, []byte(embeddedTectonicVersion), 0o644); err != nil {
		return "", err
	}
	return dest, nil
}
