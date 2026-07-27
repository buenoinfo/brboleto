package linha

import "testing"

func TestExtrairDeTexto(t *testing.T) {
	valor := "00190500954014481606906809350314337370000000100"
	linhas, err := extrairDeTexto("Boleto: 00190.50095 40144.816069 06809.350314 3 37370000000100")
	if err != nil {
		t.Fatal(err)
	}
	if len(linhas) != 1 || linhas[0].Valor != valor {
		t.Fatalf("resultado inesperado: %#v", linhas)
	}
	if linhas[0].Banco != "001" || linhas[0].Moeda != "9" {
		t.Fatalf("campos inesperados: %#v", linhas[0])
	}
}

func TestCodigoBarras44(t *testing.T) {
	l := LinhaDigitavel{Valor: "00190500954014481606906809350314337370000000100"}
	want := "00193373700000001000500940144816060680935031"
	if got := l.CodigoBarras44(); got != want {
		t.Fatalf("código = %s, esperado %s", got, want)
	}
}

func TestRejeitaLinhaDeArrecadacao(t *testing.T) {
	if _, err := extrairDeTexto("123456789012345678901234567890123456789012345678"); err == nil {
		t.Fatal("esperava erro")
	}
}
