//go:build !(linux && amd64)

package build

// embeddedTectonic fica vazio em plataformas sem binário embutido ainda
// (seção 10) — CheckTectonic() cai para o PATH do sistema nesse caso. Para
// adicionar uma plataforma, crie embedded_<GOOS>_<GOARCH>.go seguindo o
// padrão de embedded_linux_amd64.go e ajuste a build tag deste arquivo para
// continuar excluindo-a.
var embeddedTectonic []byte

const embeddedTectonicVersion = ""
