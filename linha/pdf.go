package linha

import (
	"fmt"
	"os/exec"
)

// ExtrairDePDF extrai linhas de um PDF usando o executável pdftotext.
func ExtrairDePDF(caminho string) ([]LinhaDigitavel, error) {
	comando := exec.Command("pdftotext", "-layout", caminho, "-")
	texto, err := comando.Output()
	if err != nil {
		return nil, fmt.Errorf("falha ao extrair texto do PDF: %w", err)
	}
	return extrairDeTexto(string(texto))
}
