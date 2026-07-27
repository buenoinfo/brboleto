# brboleto

Biblioteca Go para extrair linhas digitáveis de boletos bancários brasileiros.

O pacote público suporta boletos bancários no formato de 47 dígitos. A entrada deve ser um arquivo PDF ou o conteúdo de um PDF codificado em Base64. Boletos de arrecadação de 48 dígitos não fazem parte do escopo atual.

## Requisitos

- Go 1.22 ou superior;
- `pdftotext` disponível no `PATH`.

Instalação do `pdftotext`:

```bash
# macOS
brew install poppler

# Debian/Ubuntu
sudo apt-get install poppler-utils
```

No Windows, instale o Poppler para Windows e adicione o diretório que contém `pdftotext.exe` ao `PATH`.

## Instalação

```bash
go get github.com/buenoinfo/brboleto
```

## Uso com arquivo PDF

```go
package main

import (
    "fmt"
    "log"
    "github.com/buenoinfo/brboleto/linha"
)

func main() {
    linhas, err := linha.ExtrairDePDF("./boleto.pdf")
    if err != nil { log.Fatal(err) }
    for _, item := range linhas {
        fmt.Println(item.Valor, item.ValorReal, item.CodigoBarras44())
    }
}
```

## Uso com PDF em Base64

```go
linhas, err := linha.ExtrairDeBase64(pdfBase64)
if err != nil {
    log.Fatal(err)
}
```

## Resultado

Cada item retornado é uma `linha.LinhaDigitavel`:

```go
type LinhaDigitavel struct {
    Valor       string
    Banco       string
    Moeda       string
    DV          string
    Vencimento  string
    ValorReal   string
}
```

Também estão disponíveis `item.ValidarDV()` e `item.CodigoBarras44()`. `Valor` contém os 47 dígitos sem espaços ou pontuação. `ValorReal` usa ponto como separador decimal, por exemplo `100.00`.

## Erros

As funções podem retornar erros quando o PDF não existe, `pdftotext` não está instalado, o Base64 é inválido, nenhuma linha válida é encontrada ou os dígitos verificadores são inválidos.

## Desenvolvimento

```bash
task check       # formata, testa, executa vet e compila
task test        # executa os testes
task test-verbose
task coverage    # gera cobertura
task tidy        # organiza o go.mod
```

Sem Task:

```bash
gofmt -w .
go test ./...
go vet ./...
go build ./...
```
