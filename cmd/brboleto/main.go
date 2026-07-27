package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/buenoinfo/brboleto/linha"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "uso: brboleto <caminho-do-pdf>\n")
		os.Exit(2)
	}

	caminho := os.Args[1]
	var linhas []linha.LinhaDigitavel
	var err error
	if filepath.Ext(caminho) == ".rtf" {
		linhas, err = linha.ExtrairDeRTF(caminho)
	} else {
		linhas, err = linha.ExtrairDePDF(caminho)
	}
	if err != nil {
		log.Fatal(err)
	}

	for indice, item := range linhas {
		fmt.Printf("Linha %d:\n", indice+1)
		fmt.Printf("  Linha digitável: %s\n", item.Valor)
		fmt.Printf("  Banco: %s\n", item.Banco)
		fmt.Printf("  Moeda: %s\n", item.Moeda)
		fmt.Printf("  DV: %s\n", item.DV)
		fmt.Printf("  Vencimento: %s\n", item.Vencimento)
		fmt.Printf("  Fator de vencimento: %s\n", item.FatorVencimento)
		fmt.Printf("  Valor: %s\n", item.ValorReal)
		fmt.Printf("  Código de barras: %s\n", item.CodigoBarras44())
	}
}
