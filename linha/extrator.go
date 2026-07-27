package linha

import (
	"fmt"
	"regexp"
	"strings"
)

var reLinhaDigitavel = regexp.MustCompile(`\b\d(?:[ \t.-]*\d){46,47}\b`)

func extrairDeTexto(texto string) ([]LinhaDigitavel, error) {
	var linhas []LinhaDigitavel
	for _, encontrado := range reLinhaDigitavel.FindAllString(texto, -1) {
		digitos := somenteDigitos(encontrado)
		if len(digitos) != 47 {
			continue
		}
		linha, err := interpretar(digitos)
		if err == nil {
			linhas = append(linhas, linha)
		}
	}
	if len(linhas) == 0 {
		return nil, ErrNenhumaLinha
	}
	return linhas, nil
}

func interpretar(valor string) (LinhaDigitavel, error) {
	if len(valor) != 47 {
		return LinhaDigitavel{}, ErrLinhaInvalida
	}
	linha := LinhaDigitavel{
		Valor: valor, Banco: valor[0:3], Moeda: valor[3:4],
		DV: valor[32:33], Vencimento: formatarVencimento(valor[33:37]),
		FatorVencimento: valor[33:37],
		ValorReal:       formatarValor(valor[37:47]),
	}
	if !linha.ValidarDV() {
		return LinhaDigitavel{}, fmt.Errorf("%w: %s", ErrDVInvalido, valor)
	}
	return linha, nil
}

func somenteDigitos(valor string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, valor)
}
