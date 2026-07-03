package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportImage_CopiesAndReturnsRelativePath(t *testing.T) {
	w := newTestWorkspace(t)

	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "grafico.png")
	if err := os.WriteFile(src, []byte("fake-png-bytes"), 0o644); err != nil {
		t.Fatal(err)
	}

	rel, err := w.ImportImage(src)
	if err != nil {
		t.Fatalf("ImportImage: %v", err)
	}
	if rel != "imagens/grafico.png" {
		t.Errorf("rel = %q, esperava %q", rel, "imagens/grafico.png")
	}

	data, err := os.ReadFile(filepath.Join(w.Root, "imagens", "grafico.png"))
	if err != nil {
		t.Fatalf("arquivo copiado não encontrado: %v", err)
	}
	if string(data) != "fake-png-bytes" {
		t.Errorf("conteúdo copiado = %q", data)
	}
}

func TestImportImage_RenamesOnCollision(t *testing.T) {
	w := newTestWorkspace(t)

	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "grafico.png")
	if err := os.WriteFile(src, []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}

	first, err := w.ImportImage(src)
	if err != nil {
		t.Fatalf("primeiro ImportImage: %v", err)
	}
	if err := os.WriteFile(src, []byte("v2"), 0o644); err != nil {
		t.Fatal(err)
	}
	second, err := w.ImportImage(src)
	if err != nil {
		t.Fatalf("segundo ImportImage: %v", err)
	}

	if first == second {
		t.Fatalf("esperava caminhos diferentes, ambos foram %q", first)
	}
	if second != "imagens/grafico-1.png" {
		t.Errorf("second = %q, esperava %q", second, "imagens/grafico-1.png")
	}
}

func TestImportAttachment_CopiesToAnexos(t *testing.T) {
	w := newTestWorkspace(t)

	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "capa.pdf")
	if err := os.WriteFile(src, []byte("fake-pdf-bytes"), 0o644); err != nil {
		t.Fatal(err)
	}

	rel, err := w.ImportAttachment(src)
	if err != nil {
		t.Fatalf("ImportAttachment: %v", err)
	}
	if rel != "anexos/capa.pdf" {
		t.Errorf("rel = %q, esperava %q", rel, "anexos/capa.pdf")
	}
}
