package linha

import (
	"encoding/base64"
	"fmt"
	"os"
)

// ExtrairDeBase64 extrai linhas de um PDF codificado em Base64.
func ExtrairDeBase64(pdfBase64 string) ([]LinhaDigitavel, error) {
	dados, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBase64Invalido, err)
	}
	arquivo, err := os.CreateTemp("", "brboleto-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("criar PDF temporário: %w", err)
	}
	nome := arquivo.Name()
	defer os.Remove(nome)
	if _, err = arquivo.Write(dados); err != nil {
		arquivo.Close()
		return nil, fmt.Errorf("escrever PDF temporário: %w", err)
	}
	if err = arquivo.Close(); err != nil {
		return nil, fmt.Errorf("fechar PDF temporário: %w", err)
	}
	return ExtrairDePDF(nome)
}
