# brboleto

Biblioteca Go para extrair linhas digitáveis de boletos bancários brasileiros.

O pacote público suporta boletos bancários no formato de 47 dígitos. A entrada deve ser um arquivo PDF ou o conteúdo de um PDF codificado em Base64. A leitura do PDF é feita diretamente em Go, sem `pdftotext` ou outras dependências externas. Boletos de arrecadação de 48 dígitos não fazem parte do escopo atual.

## Requisitos

- Go 1.22 ou superior;
- sistema operacional Windows, macOS ou Linux;
- um PDF com texto selecionável.

## Instalação

```bash
go get github.com/buenoinfo/brboleto
```

## Uso com arquivo RTF

Quando o RTF original estiver disponível, prefira esse formato. A biblioteca extrai o texto diretamente, sem OCR e sem dependências instaladas no sistema:

```go
linhas, err := linha.ExtrairDeRTF("./boleto.rtf")
if err != nil {
    log.Fatal(err)
}
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
    FatorVencimento string
    ValorReal   string
}
```

Também estão disponíveis `item.ValidarDV()` e `item.CodigoBarras44()`. `Valor` contém os 47 dígitos sem espaços ou pontuação. `Vencimento` usa o formato `dd/mm/aa`; `FatorVencimento` preserva o fator original. `ValorReal` usa ponto como separador decimal, por exemplo `100.00`.

## Erros

As funções podem retornar erros quando o PDF não existe ou é inválido, o Base64 é inválido, nenhuma linha válida é encontrada ou os dígitos verificadores são inválidos.

PDFs escaneados, que contêm apenas imagens, não possuem texto para extração e exigem OCR. OCR não faz parte desta versão. Se o RTF original estiver disponível, use `ExtrairDeRTF`.

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

## Teste com um PDF

Para testar manualmente um boleto, informe o caminho do PDF:

```bash
task boleto -- ./caminho/para/boleto.rtf
```

Ou execute diretamente:

```bash
go run ./cmd/brboleto ./caminho/para/boleto.rtf
```

O comando identifica automaticamente arquivos `.rtf` e `.pdf` e imprime a linha digitável, banco, moeda, dígitos verificadores, vencimento, valor e código de barras encontrados.
