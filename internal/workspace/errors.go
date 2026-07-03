package workspace

import "fmt"

// CorruptedFileError é retornado quando um arquivo JSON do workspace não pode
// ser decodificado, com uma mensagem amigável em pt-BR (a spec exige que o Go
// nunca vaze erros crus de encoding/json para o usuário).
type CorruptedFileError struct {
	Path string
	Err  error
}

func (e *CorruptedFileError) Error() string {
	return fmt.Sprintf("o arquivo '%s' está corrompido e não pôde ser lido: %v", e.Path, e.Err)
}

func (e *CorruptedFileError) Unwrap() error { return e.Err }
