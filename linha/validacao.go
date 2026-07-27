package linha

import (
	"strconv"
	"time"
)

func formatarVencimento(fator string) string {
	valor, err := strconv.Atoi(fator)
	if err != nil || valor == 0 {
		return ""
	}
	base := time.Date(1997, 10, 7, 0, 0, 0, 0, time.UTC)
	if valor >= 1000 {
		base = time.Date(2025, 2, 22, 0, 0, 0, 0, time.UTC)
		valor -= 1000
	}
	data := base.AddDate(0, 0, valor)
	return data.Format("02/01/06")
}

func modulo10(valor string) string {
	soma, peso := 0, 2
	for i := len(valor) - 1; i >= 0; i-- {
		digito := int(valor[i]-'0') * peso
		if digito > 9 {
			digito = digito/10 + digito%10
		}
		soma += digito
		if peso == 2 {
			peso = 1
		} else {
			peso = 2
		}
	}
	return strconv.Itoa((10 - soma%10) % 10)
}

func modulo11(valor string) string {
	soma, peso := 0, 2
	for i := len(valor) - 1; i >= 0; i-- {
		soma += int(valor[i]-'0') * peso
		peso++
		if peso > 9 {
			peso = 2
		}
	}
	dv := 11 - soma%11
	if dv == 0 || dv == 10 || dv == 11 {
		dv = 1
	}
	return strconv.Itoa(dv)
}

func formatarValor(valor string) string {
	inteiro := valor[:len(valor)-2]
	for len(inteiro) > 1 && inteiro[0] == '0' {
		inteiro = inteiro[1:]
	}
	return inteiro + "." + valor[len(valor)-2:]
}
