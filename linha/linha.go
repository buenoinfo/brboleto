// Package linha extrai e interpreta linhas digitáveis de boletos bancários.
package linha

import "errors"

var (
	ErrNenhumaLinha   = errors.New("nenhuma linha digitável encontrada")
	ErrLinhaInvalida  = errors.New("linha digitável inválida")
	ErrDVInvalido     = errors.New("dígito verificador inválido")
	ErrBase64Invalido = errors.New("base64 inválido")
)

// LinhaDigitavel representa uma linha digitável bancária de 47 dígitos.
type LinhaDigitavel struct {
	Valor           string
	Banco           string
	Moeda           string
	DV              string
	Vencimento      string
	FatorVencimento string
	ValorReal       string
}

// ValidarDV valida os três dígitos dos campos e o dígito geral do boleto.
func (l LinhaDigitavel) ValidarDV() bool {
	if len(l.Valor) != 47 {
		return false
	}
	return modulo10(l.Valor[0:9]) == l.Valor[9:10] &&
		modulo10(l.Valor[10:20]) == l.Valor[20:21] &&
		modulo10(l.Valor[21:31]) == l.Valor[31:32] &&
		modulo11(l.CodigoBarras44()[:4]+l.CodigoBarras44()[5:]) == l.Valor[32:33]
}

// CodigoBarras44 converte a linha digitável em código de barras.
func (l LinhaDigitavel) CodigoBarras44() string {
	if len(l.Valor) != 47 {
		return ""
	}
	return l.Valor[0:4] + l.Valor[32:33] + l.Valor[33:47] + l.Valor[4:9] + l.Valor[10:20] + l.Valor[21:31]
}
