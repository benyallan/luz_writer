package workspace

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ImportImage copia srcPath para imagens/ (renomeando em caso de colisão) e
// devolve o caminho relativo ao workspace, com separadores "/" — pronto para
// o atributo src do nó luzImage e para a compilação LaTeX.
func (w *Workspace) ImportImage(srcPath string) (string, error) {
	return w.copyIntoDir(srcPath, w.imagesDir(), dirImages)
}

// ImportAttachment copia srcPath para anexos/ (renomeando em caso de
// colisão) e devolve o caminho relativo — usado pelo plugin pdfpages.
func (w *Workspace) ImportAttachment(srcPath string) (string, error) {
	return w.copyIntoDir(srcPath, w.attachDir(), dirAttach)
}

// FileExistsRelative indica se relPath (ex.: "imagens/grafico.png") existe
// dentro do workspace — usado pelo Rule Engine (R008, R009).
func (w *Workspace) FileExistsRelative(relPath string) bool {
	if relPath == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(w.Root, relPath))
	return err == nil
}

func (w *Workspace) copyIntoDir(srcPath, destDir, destDirName string) (string, error) {
	if srcPath == "" {
		return "", fmt.Errorf("nenhum arquivo selecionado")
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	base := filepath.Base(srcPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	dest := filepath.Join(destDir, base)
	for n := 1; ; n++ {
		if _, err := os.Stat(dest); os.IsNotExist(err) {
			break
		}
		dest = filepath.Join(destDir, fmt.Sprintf("%s-%d%s", name, n, ext))
	}

	if err := copyFile(srcPath, dest); err != nil {
		return "", fmt.Errorf("não foi possível copiar '%s': %w", base, err)
	}

	return destDirName + "/" + filepath.Base(dest), nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
