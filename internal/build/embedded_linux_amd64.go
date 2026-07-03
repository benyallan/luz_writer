//go:build linux && amd64

package build

import _ "embed"

// embeddedTectonic é o binário oficial do Tectonic empacotado no próprio
// executável do Luz Writer (seção 10): o app não depende do usuário instalar
// nada — CheckTectonic() extrai esta cópia para um diretório de dados do
// usuário na primeira execução. Atualizar a versão: baixar o release
// correspondente em https://github.com/tectonic-typesetting/tectonic e
// substituir internal/build/assets/tectonic_linux_amd64, atualizando
// embeddedTectonicVersion abaixo.
//
//go:embed assets/tectonic_linux_amd64
var embeddedTectonic []byte

const embeddedTectonicVersion = "0.16.9"
