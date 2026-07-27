package linha

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/ledongthuc/pdf"
)

// ExtrairDePDF extrai linhas de um arquivo PDF sem dependências externas.
func ExtrairDePDF(caminho string) ([]LinhaDigitavel, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return nil, fmt.Errorf("ler PDF: %w", err)
	}
	return extrairDePDFBytes(dados)
}

func extrairDePDFBytes(dados []byte) ([]LinhaDigitavel, error) {
	leitor, err := pdf.NewReader(bytes.NewReader(dados), int64(len(dados)))
	if err != nil {
		return nil, fmt.Errorf("abrir PDF: %w", err)
	}
	texto, err := leitor.GetPlainText()
	if err != nil {
		return nil, fmt.Errorf("extrair texto do PDF: %w", err)
	}
	conteudo, err := io.ReadAll(texto)
	if err != nil {
		return nil, fmt.Errorf("ler texto do PDF: %w", err)
	}
	return extrairDeTexto(string(conteudo))
}
