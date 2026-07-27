# Extração de Linha Digitável de Boletos

Documentação técnica do processo de extração da linha digitável (código de barras composto) a partir de PDFs de boleto bancário, com foco na criação de uma biblioteca reutilizável.

---

## 1. O que é a Linha Digitável

A **linha digitável** (também chamada de "código de barras composto" ou "nosso número formatado") é uma sequência de **47 ou 48 dígitos** que codifica todas as informações necessárias para pagamento de um boleto bancário brasileiro.

### Estrutura (ISO 646 / FEBRABAN)

```
Posições  Campo                          Dígitos
─────────────────────────────────────────────────
01-03     Código do banco (COMPE)        3
04-04     Código da moeda (9 = BRL)      1
05-05     Dígito verificador geral (DV)  1
06-19     Fator de vencimento + valor    14
20-44     Campo livre (varia por banco)  25
─────────────────────────────────────────────────
Total                                    47 ou 48
```

**Exemplo real (47 dígitos):**
```
74891126280760260812206704471041413740000629866
```

### Diferença entre Código de Barras e Linha Digitável

| Formato | Qtd dígitos | Uso |
|---------|-------------|-----|
| **Código de barras** | 44 dígitos | Leitura óptica (scanner/boleto registrado) |
| **Linha digitável** | 47-48 dígitos | Digitação manual / colar no app bancário |

A linha digitável é o código de barras **codificado** com dígitos verificadores por campo. Para pagamento, a linha digitável é o que o cliente cola no app bancário.

---

## 2. Fluxo de Extração

```
PDF do boleto (filesystem)
  │
  ▼
pdftotext (poppler-utils)
  │
  ▼
Texto bruto do PDF
  │
  ▼
Regex: /(\d{47,48})/
  │
  ▼
Linha digitável extraída
  │
  ▼
Validação: dígito verificador (módulo 10)
  │
  ▼
Linha digitável válida (ou erro)
```

### Dependência Externa

O `pdftotext` pertence ao pacote **poppler-utils**:

```bash
# Ubuntu/Debian
sudo apt-get install poppler-utils

# macOS
brew install poppler

# CentOS/RHEL
sudo yum install poppler-utils
```

### Alternativa: extração sem pdftotext

Se `pdftotext` não estiver disponível, é possível usar bibliotecas nativas:

| Linguagem | Biblioteca | Nota |
|-----------|-----------|------|
| Go | `pdfcpu`, `unidoc/unipdf` | Pago para uso comercial |
| Python | `PyPDF2`, `pdfplumber` | `pdfplumber` mais confiável |
| Node.js | `pdf-parse`, `pdfjs-dist` | `pdf-parse` wrapper do `pdftotext` |
| Rust | `lopdf`, `pdf-extract` | `pdf-extract` usa `pdftotext` internamente |

---

## 3. Regex de Extração

### Padrão principal

```regex
\b(\d{47,48})\b
```

**Por que funciona:**
- A linha digitável é uma sequência contínua de 47-48 dígitos
- Delimitada por espaços, quebras de linha ou início/fim de texto
- `\b` assegura limites de palavra (evita falsos positivos em meio a textos)

### Padrão alternativo (com separadores)

Alguns boletos formatam a linha com pontos e espaços:

```
74891.12628 07602.608122 06704.471041 41374.000062 9866
```

Regex para este caso:

```regex
\b(\d{5}\.\d{5}\s\d{5}\.\d{6}\s\d{5}\.\d{6}\s\d{5}\.\d{6}\s\d{5}\.\d{4})\b
```

Ou normalizar antes de buscar:

```regex
# Remove pontos e espaços, depois busca 47-48 dígitos
[\d\.]+ → remove tudo que não é dígito
```

### Padrão robusto (recomendado para biblioteca)

```regex
\b(?:\d[\. ]?){47,48}\b
```

Este padrão aceita:
- Sequência pura: `74891126280760260812206704471041413740000629866`
- Com pontos: `74891.12628.07602.608122...`
- Com espaços: `74891 12628 07602 608122...`
- Com ambos: `74891.12628 07602.608122...`

### Múltiplas linhas digitáveis

Alguns boletos de carnê (grupos de parcelas) geram **várias linhas digitáveis** no mesmo PDF. A regex deve capturar todas:

```go
re := regexp.MustCompile(`\b(?:\d[\. ]?){47,48}\b`)
matches := re.FindAllString(text, -1)
```

---

## 4. Implementação em Go (Referência)

### Função core de extração

```go
package boleto

import (
    "fmt"
    "os/exec"
    "regexp"
    "strings"
    "unicode"
)

// LinhaDigitavel representa uma linha digitável extraída de um boleto
type LinhaDigitavel struct {
    Valor    string // 47-48 dígitos limpos
    Banco    string // Código COMPE (3 dígitos)
    Moeda    string // 9 = BRL
    DV       string // Dígito verificador geral
    Vencimento string // Fator de vencimento
    ValorReal string // Valor em reais (com.centavos)
}

var reLinhaDigitavel = regexp.MustCompile(`\b(?:\d[\. ]?){47,48}\b`)
var reApenasDigitos = regexp.MustCompile(`\D`)

// ExtrairDePDF extrai linhas digitáveis de um PDF usando pdftotext
func ExtrairDePDF(pdfPath string) ([]LinhaDigitavel, error) {
    // 1. Converter PDF para texto
    texto, err := pdfParaTexto(pdfPath)
    if err != nil {
        return nil, fmt.Errorf("falha ao extrair texto do PDF: %w", err)
    }

    // 2. Buscar linhas digitáveis
    return ExtrairDeTexto(texto)
}

// ExtrairDeTexto busca linhas digitáveis em texto já extraído
func ExtrairDeTexto(texto string) ([]LinhaDigitavel, error) {
    matches := reLinhaDigitavel.FindAllString(texto, -1)

    var linhas []LinhaDigitavel
    for _, m := range matches {
        // Limpar: manter apenas dígitos
        limpa := reApenasDigitos.ReplaceAllString(m, "")

        if len(limpa) < 47 || len(limpa) > 48 {
            continue
        }

        linha := parseLinhaDigitavel(limpa)
        if linha != nil {
            linhas = append(linhas, *linha)
        }
    }

    if len(linhas) == 0 {
        return nil, fmt.Errorf("nenhuma linha digitável encontrada")
    }

    return linhas, nil
}

// ExtrairDeBase64 extrai de um PDF em base64
func ExtrairDeBase64(base64PDF string) ([]LinhaDigitavel, error) {
    // Decodificar base64 para bytes
    data, err := base64.StdEncoding.DecodeString(base64PDF)
    if err != nil {
        return nil, fmt.Errorf("base64 inválido: %w", err)
    }

    // Salvar em arquivo temporário
    tmpFile, err := os.CreateTemp("", "boleto-*.pdf")
    if err != nil {
        return nil, fmt.Errorf("falha ao criar arquivo temporário: %w", err)
    }
    defer os.Remove(tmpFile.Name())

    if _, err := tmpFile.Write(data); err != nil {
        tmpFile.Close()
        return nil, fmt.Errorf("falha ao escrever PDF temporário: %w", err)
    }
    tmpFile.Close()

    return ExtrairDePDF(tmpFile.Name())
}

func pdfParaTexto(pdfPath string) (string, error) {
    cmd := exec.Command("pdftotext", "-layout", pdfPath, "-")
    out, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("pdftotext falhou: %s — %w", string(out), err)
    }
    return string(out), nil
}

func parseLinhaDigitavel(valor string) *LinhaDigitavel {
    if len(valor) < 47 {
        return nil
    }

    // Validar que todos são dígitos
    for _, r := range valor {
        if !unicode.IsDigit(r) {
            return nil
        }
    }

    // Validar dígito verificador geral (posição 5)
    dvCalculado := modulo10(valor[:4] + valor[5:])
    if valor[4:5] != dvCalculado {
        return nil // DV inválido
    }

    // Fator de vencimento (posições 6-9)
    fator := valor[5:9]

    // Valor (posições 10-19) — último dígito são centavos
    valorRaw := valor[9:19]
    valorFormatado := formatarValor(valorRaw)

    return &LinhaDigitavel{
        Valor:      valor,
        Banco:      valor[0:3],
        Moeda:      valor[3:4],
        DV:         valor[4:5],
        Vencimento: fator,
        ValorReal:  valorFormatado,
    }
}

// modulo10 implements FEBRABAN módulo 10
func modulo10(campo string) string {
    pesos := []int{2, 1}
    soma := 0
    idx := 0

    for i := len(campo) - 1; i >= 0; i-- {
        d := int(campo[i] - '0')
        p := pesos[idx%2]
        produto := d * p

        if produto > 9 {
            produto = (produto / 10) + (produto % 10)
        }

        soma += produto
        idx++
    }

    resto := soma % 10
    resultado := 10 - resto
    if resultado == 10 {
        resultado = 0
    }

    return fmt.Sprintf("%d", resultado)
}

func formatarValor(valorRaw string) string {
    // Último dígito = centavos
    inteiro := valorRaw[:len(valorRaw)-2]
    centavos := valorRaw[len(valorRaw)-2:]

    // Remover zeros à esquerda do inteiro
    inteiro = strings.TrimLeft(inteiro, "0")
    if inteiro == "" {
        inteiro = "0"
    }

    return inteiro + "." + centavos
}
```

### Servidor HTTP (para N8N)

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/seu-org/boleto-lib"
)

type BoletoRequest struct {
    FileBase64 string `json:"file_base64"`
    Filename   string `json:"filename"`
}

type BoletoResponse struct {
    Success          bool     `json:"success"`
    LinhasDigitaveis []string `json:"linhas_digitaveis"`
    Filename         string   `json:"filename"`
    Error            string   `json:"error,omitempty"`
}

func handleBoleto(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(BoletoResponse{Error: "method not allowed"})
        return
    }

    var req BoletoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(BoletoResponse{Error: "invalid JSON"})
        return
    }

    if req.FileBase64 == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(BoletoResponse{Error: "file_base64 required"})
        return
    }

    linhas, err := boleto.ExtrairDeBase64(req.FileBase64)
    if err != nil {
        w.WriteHeader(http.StatusUnprocessableEntity)
        json.NewEncoder(w).Encode(BoletoResponse{
            Success:  false,
            Error:    err.Error(),
            Filename: req.Filename,
        })
        return
    }

    var linhasStr []string
    for _, l := range linhas {
        linhasStr = append(linhasStr, l.Valor)
    }

    json.NewEncoder(w).Encode(BoletoResponse{
        Success:          true,
        LinhasDigitaveis: linhasStr,
        Filename:         req.Filename,
    })
}

func main() {
    http.HandleFunc("/boleto", handleBoleto)
    log.Println("Boleto extraction server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Testes

```go
package boleto_test

import (
    "testing"
    "github.com/seu-org/boleto-lib"
)

func TestExtrairDeTexto_LinhaValida(t *testing.T) {
    texto := `
    Boleto Bancário
    Banco Itaú

    Pagável até: 15/01/2025

    74891.12628 07602.608122 06704.471041 41374.000062 9866

    Valor do documento: R$ 629,86
    `

    linhas, err := boleto.ExtrairDeTexto(texto)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(linhas) != 1 {
        t.Fatalf("expected 1 linha, got %d", len(linhas))
    }

    if linhas[0].Valor != "74891126280760260812206704471041413740000629866" {
        t.Errorf("wrong value: %s", linhas[0].Valor)
    }

    if linhas[0].Banco != "748" {
        t.Errorf("wrong banco: %s", linhas[0].Banco)
    }
}

func TestExtrairDeTexto_SemLinha(t *testing.T) {
    texto := "Este é um documento sem linha digitável"
    _, err := boleto.ExtrairDeTexto(texto)
    if err == nil {
        t.Fatal("expected error for text without linha digitavel")
    }
}

func TestExtrairDeTexto_MultiplasLinhas(t *testing.T) {
    texto := `
    Parcela 1: 74891126280760260812206704471041413740000629866
    Parcela 2: 74891126280760260812206704471041413740000629874
    `

    linhas, err := boleto.ExtrairDeTexto(texto)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(linhas) != 2 {
        t.Fatalf("expected 2 linhas, got %d", len(linhas))
    }
}

func TestModulo10(t *testing.T) {
    // Teste com valor conhecido
    campo := "7489112628076026081220670447104141374000062986"
    dv := boleto.Modulo10(campo)
    if dv != "6" {
        t.Errorf("modulo10(%s) = %s, want 6", campo, dv)
    }
}

func TestExtrairDeBase64(t *testing.T) {
    // Base64 de um PDF de teste (placeholder)
    // Em produção, usar um PDF real de boleto
    base64PDF := "JVBERi0xLjQKMSAwIG9iago8PAovVHlwZSAvQ2F0YWxvZwovUGFn..."
    _, err := boleto.ExtrairDeBase64(base64PDF)
    // O teste pode falhar se o pdftotext não estiver instalado
    if err != nil {
        t.Skip("pdftotext not available, skipping base64 test")
    }
}
```

---

## 5. Especificação da Biblioteca (Projeto)

### Nome sugerido

- **Go:** `github.com/seu-org/boleto-extract`
- **Python:** `boleto_extract`
- **Node.js:** `@seu-org/boleto-extract`

### Interface pública (Go)

```go
// Package boleto extrai linhas digitáveis de PDFs de boleto.
package boleto

// ExtrairDePDF extrai linhas digitáveis de um arquivo PDF.
func ExtrairDePDF(pdfPath string) ([]LinhaDigitavel, error)

// ExtrairDeBase64 extrai de um PDF codificado em base64.
func ExtrairDeBase64(base64PDF string) ([]LinhaDigitavel, error)

// ExtrairDeTexto busca linhas digitáveis em texto já extraído.
func ExtrairDeTexto(texto string) ([]LinhaDigitavel, error)

// LinhaDigitavel é o resultado da extração.
type LinhaDigitavel struct {
    Valor      string // 47-48 dígitos limpos
    Banco      string // Código COMPE
    Moeda      string // 9 = BRL
    DV         string // Dígito verificador geral
    Vencimento string // Fator de vencimento
    ValorReal  string // Valor formatado
}

// ValidarDV verifica o dígito verificador geral (módulo 10).
func (l *LinhaDigitavel) ValidarDV() bool

// CodigoBarras44 converte linha digitável para código de barras de 44 dígitos.
func (l *LinhaDigitavel) CodigoBarras44() string
```

### Interface pública (Python)

```python
from boleto_extract import ExtratorLinhaDigitavel

extrator = ExtratorLinhaDigitavel()

# De arquivo
linhas = extrator.de_pdf("/caminho/boleto.pdf")

# De base64
linhas = extrator.de_base64(pdf_base64)

# De texto
linhas = extrator.de_texto(texto_extraido)

for linha in linhas:
    print(f"Linha: {linha.valor}")
    print(f"Banco: {linha.banco}")
    print(f"Valor: R$ {linha.valor_real}")
```

### Interface pública (Node.js)

```javascript
const { extrairDePDF, extrairDeBase64, extrairDeTexto } = require('@seu-org/boleto-extract');

// De arquivo
const linhas = await extrairDePDF('/caminho/boleto.pdf');

// De base64
const linhas = await extrairDeBase64(pdfBase64);

// Síncrono (se já tem o texto)
const linhas = extrairDeTexto(textoExtraido);
```

### Dependências

| Dependência | Obrigatória | Uso |
|-------------|-------------|-----|
| `poppler-utils` (pdftotext) | Sim | Conversão PDF → texto |
| Regex nativa | Sim | Extração da linha |
| `encoding/base64` | Não | Apenas se usar interface base64 |

### Configuração

```yaml
# config.yml (para servidor HTTP)
server:
  port: 8080
  timeout: 30s

extraction:
  # Usar pdftotext (true) ou biblioteca nativa (false)
  use_pdftotext: true
  # pdftotext flags adicionais
  pdftotext_flags: ["-layout"]
  # Timeout para pdftotext
  timeout: 10s
```

---

## 6. Integração com N8N

### Workflow N8N

```
Webhook (POST /boleto)
  │
  ▼
HTTP Request → POST /extract (servidor Go)
  │
  ▼
Code (JS): Montar mensagem com linha digitável
  │
  ▼
HTTP Request → WAHA (enviar WhatsApp)
```

### Nó HTTP Request no N8N

```json
{
  "method": "POST",
  "url": "http://localhost:8080/boleto",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "file_base64": "={{ $json.boleto_pdf_base64 }}",
    "filename": "={{ $json.boleto_pdf_filename }}"
  }
}
```

### Nó Code (montar mensagem com linha)

```javascript
const extracao = $node["ExtrairLinha"].json;

if (!extracao.success || !extracao.linhas_digitaveis || extracao.linhas_digitaveis.length === 0) {
  throw new Error("Linha digitável não encontrada no boleto");
}

const linhaDigitavel = extracao.linhas_digitaveis[0];
const dados = $node["Decodificador"].json;

const mensagem = `${dados.mensagem}

*Código para pagamento:*
\`${linhaDigitavel}\`

Copie e cole no app do seu banco.`;

return [{
  json: {
    ...dados,
    mensagem_completa: mensagem,
    linha_digitavel: linhaDigitavel
  }
}];
```

---

## 7. Casos Especiais e Robustez

### Problemas comuns

| Problema | Causa | Solução |
|----------|-------|---------|
| pdftotext não encontra texto | PDF escaneado (imagem) | Usar OCR (Tesseract) antes |
| Linha com separadores | Boletos formatados | Regex com `[\. ]?` entre dígitos |
| Múltiplas linhas | Carnê / parcelas | `FindAll` em vez de `Find` |
| DV inválido | Extração incorreta | Validar e reportar, não retornar |
| PDF protegido | Criptografia | Reportar erro específico |
| pdftotext não instalado | Deploy sem poppler | Fallback para biblioteca nativa |

### PDF escaneado (imagem)

Se o PDF é uma imagem escaneada, `pdftotext` retorna vazio. Solução:

```bash
# Converter imagem para texto com Tesseract
# 1. Converter PDF para imagem
pdftoppm -png -r 300 boleto.pdf boleto

# 2. OCR
tesseract boleto-1.png boleto_texto

# 3. Extrair linha do texto OCR
```

Em Go, usar `github.com/otiai10/gosseract2` para Tesseract.

### Fallback sem pdftotext

```go
func extrairComFallback(pdfPath string) ([]LinhaDigitavel, error) {
    // Tentar pdftotext primeiro
    linhas, err := ExtrairDePDF(pdfPath)
    if err == nil {
        return linhas, nil
    }

    // Fallback: usar biblioteca nativa
    linhas, err = extrairComPDFCPU(pdfPath)
    if err == nil {
        return linhas, nil
    }

    return nil, fmt.Errorf("extração falhou com pdftotext e pdfcpu: %w", err)
}
```

---

## 8. Checklist de Implementação

- [ ] Função `ExtrairDePDF` com pdftotext
- [ ] Função `ExtrairDeTexto` (regex pura)
- [ ] Função `ExtrairDeBase64` (decode + temporário + extração)
- [ ] Validação de DV (módulo 10)
- [ ] Validação de tamanho (47 ou 48 dígitos)
- [ ] Testes com PDFs reais de bancos: Itaú, Bradesco, Banco do Brasil, Caixa, Santander
- [ ] Servidor HTTP `POST /boleto`
- [ ] Dockerfile com poppler-utils
- [ ] README com exemplos de uso
- [ ] publicação como pacote Go/Python/Node

---

## 9. Referências

- [FEBRABAN - Leitura Alternativa do Código de Barras](https://portal.febraban.org.br/pagina/315498/ler/codigo-de-barras)
- [ISO 646 - Representação de caracteres](https://www.iso.org/standard/30769.html)
- [poppler-utils pdftotext](https://poppler.freedesktop.org/)
- Documentação interna: `FLUXO_ENVIO_NF_BOLETO.md` seção 5.6
