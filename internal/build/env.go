package build

import "os/exec"

// CheckTectonic verifica se o binário tectonic está disponível no PATH.
// Sua ausência nunca deve travar o app (seção 10 da spec) — apenas desabilita
// a exportação, sinalizada ao usuário como um Problem persistente.
func CheckTectonic() (found bool, path string) {
	p, err := exec.LookPath("tectonic")
	if err != nil {
		return false, ""
	}
	return true, p
}
