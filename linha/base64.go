package linha

import (
	"encoding/base64"
	"fmt"
)

// ExtrairDeBase64 extrai linhas de um PDF codificado em Base64.
func ExtrairDeBase64(pdfBase64 string) ([]LinhaDigitavel, error) {
	dados, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBase64Invalido, err)
	}
	return extrairDePDFBytes(dados)
}
