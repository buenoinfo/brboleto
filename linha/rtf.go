package linha

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ExtrairDeRTF extrai linhas de um arquivo RTF sem dependências externas.
func ExtrairDeRTF(caminho string) ([]LinhaDigitavel, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return nil, fmt.Errorf("ler RTF: %w", err)
	}
	texto := decodificarRTF(string(dados))
	return extrairDeTexto(texto)
}

func decodificarRTF(rtf string) string {
	var resultado strings.Builder
	grupos := []bool{false}
	skipFallback := 0

	for i := 0; i < len(rtf); {
		switch rtf[i] {
		case '{':
			grupos = append(grupos, grupos[len(grupos)-1])
			i++
		case '}':
			if len(grupos) > 1 {
				grupos = grupos[:len(grupos)-1]
			}
			i++
		case '\\':
			i++
			if i >= len(rtf) || grupos[len(grupos)-1] {
				continue
			}
			if rtf[i] == '\'' && i+2 < len(rtf) {
				if valor, err := strconv.ParseUint(rtf[i+1:i+3], 16, 8); err == nil {
					resultado.WriteByte(byte(valor))
				}
				i += 3
				continue
			}
			if rtf[i] == '\\' || rtf[i] == '{' || rtf[i] == '}' {
				resultado.WriteByte(rtf[i])
				i++
				continue
			}
			inicio := i
			for i < len(rtf) && ((rtf[i] >= 'a' && rtf[i] <= 'z') || (rtf[i] >= 'A' && rtf[i] <= 'Z')) {
				i++
			}
			comando := rtf[inicio:i]
			negativo := false
			if i < len(rtf) && rtf[i] == '-' {
				negativo = true
				i++
			}
			numeroInicio := i
			for i < len(rtf) && rtf[i] >= '0' && rtf[i] <= '9' {
				i++
			}
			numero := rtf[numeroInicio:i]
			if i < len(rtf) && rtf[i] == ' ' {
				i++
			}

			switch comando {
			case "fonttbl", "colortbl", "stylesheet", "info", "pict", "object", "header", "footer":
				grupos[len(grupos)-1] = true
			case "par", "line", "cell", "row":
				resultado.WriteByte('\n')
			case "tab":
				resultado.WriteByte('\t')
			case "u":
				if numero != "" {
					n, _ := strconv.Atoi(numero)
					if negativo {
						n = -n
					}
					resultado.WriteRune(rune(n))
					skipFallback = 1
				}
			}
		case '\n', '\r':
			i++
		default:
			if skipFallback > 0 {
				_, tamanho := utf8.DecodeRuneInString(rtf[i:])
				i += tamanho
				skipFallback--
				continue
			}
			if !grupos[len(grupos)-1] {
				resultado.WriteByte(rtf[i])
			}
			i++
		}
	}
	return resultado.String()
}
