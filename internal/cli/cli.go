// Package cli interpreta os argumentos de `calc <operacao> <a> <b>`, chama o
// pacote calculo e apresenta o resultado em stdout ou o erro em stderr.
//
// Códigos de saída: 0 sucesso; 1 erro de cálculo (divisão por zero ou
// resultado fora da faixa); 2 erro de uso ou entrada inválida.
package cli

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/emilianoeloi/calculadora_go/internal/calculo"
)

// Códigos de saída do programa.
const (
	CodigoSucesso          = 0
	CodigoErroDeCalculo    = 1
	CodigoUsoInvalido      = 2
	quantidadeDeArgumentos = 3
)

// Executar roda a calculadora com os argumentos (sem o nome do programa),
// escreve o resultado em stdout ou a mensagem de erro em stderr e devolve o
// código de saída.
func Executar(args []string, stdout, stderr io.Writer) int {
	operacao, a, b, err := interpretar(args)
	if err != nil {
		return reportar(stderr, err)
	}
	resultado, err := calcular(operacao, a, b)
	if err != nil {
		return reportar(stderr, err)
	}
	fmt.Fprintln(stdout, formatar(resultado))
	return CodigoSucesso
}

// interpretar valida os argumentos antes de qualquer cálculo.
func interpretar(args []string) (operacao string, a, b float64, err error) {
	if len(args) != quantidadeDeArgumentos {
		return "", 0, 0, fmt.Errorf(
			"número de argumentos inválido: esperado %d, recebido %d; uso: calc <operacao> <a> <b>",
			quantidadeDeArgumentos, len(args))
	}
	operacao = args[0]
	switch operacao {
	case "soma", "subtracao", "multiplicacao", "divisao":
	default:
		return "", 0, 0, fmt.Errorf(
			"operação desconhecida %q; use soma, subtracao, multiplicacao ou divisao", operacao)
	}
	if a, err = lerOperando(args[1]); err != nil {
		return "", 0, 0, err
	}
	if b, err = lerOperando(args[2]); err != nil {
		return "", 0, 0, err
	}
	return operacao, a, b, nil
}

// lerOperando aceita somente número decimal finito com ponto. A ordem é: erro
// de faixa do ParseFloat (código 2), outro erro do ParseFloat (código 2) e,
// no texto aceito, rejeição de NaN, infinito, sublinhado e hexadecimal.
func lerOperando(texto string) (float64, error) {
	valor, err := strconv.ParseFloat(texto, 64)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			return 0, erroDeOperando{
				mensagem: fmt.Sprintf("operando %q fora da faixa de float64", texto),
				causa:    err,
			}
		}
		return 0, erroDeOperando{
			mensagem: fmt.Sprintf("operando %q não é um número válido", texto),
			causa:    err,
		}
	}
	if math.IsNaN(valor) || math.IsInf(valor, 0) {
		return 0, erroDeOperando{
			mensagem: fmt.Sprintf("operando %q não é um número finito válido", texto),
		}
	}
	if strings.Contains(texto, "_") || ehHexadecimal(texto) {
		return 0, erroDeOperando{
			mensagem: fmt.Sprintf("operando %q não é um número decimal válido (hexadecimal e sublinhado não são aceitos)", texto),
		}
	}
	return valor, nil
}

// ehHexadecimal indica se o texto, sem sinal opcional, começa com 0x ou 0X.
func ehHexadecimal(texto string) bool {
	if strings.HasPrefix(texto, "+") || strings.HasPrefix(texto, "-") {
		texto = texto[1:]
	}
	return strings.HasPrefix(texto, "0x") || strings.HasPrefix(texto, "0X")
}

// erroDeOperando mantém a mensagem em português para o usuário e preserva a
// causa original (strconv) para errors.Is e errors.As.
type erroDeOperando struct {
	mensagem string
	causa    error
}

func (e erroDeOperando) Error() string { return e.mensagem }
func (e erroDeOperando) Unwrap() error { return e.causa }

// calcular assume operação e operandos já validados por interpretar.
func calcular(operacao string, a, b float64) (float64, error) {
	switch operacao {
	case "soma":
		return calculo.Somar(a, b)
	case "subtracao":
		return calculo.Subtrair(a, b)
	case "multiplicacao":
		return calculo.Multiplicar(a, b)
	default: // "divisao"
		return calculo.Dividir(a, b)
	}
}

// formatar imprime o menor texto que representa o float64 sem perda. O zero
// negativo é normalizado para 0 antes da conversão.
func formatar(valor float64) string {
	if valor == 0 {
		valor = 0 // descarta o sinal de -0
	}
	return strconv.FormatFloat(valor, 'g', -1, 64)
}

// reportar é o único ponto que converte erro em mensagem em stderr e em
// código de saída. Erros de cálculo (divisão por zero e resultado fora da
// faixa) saem com código 1; o restante, com código 2.
func reportar(stderr io.Writer, err error) int {
	fmt.Fprintf(stderr, "erro: %v\n", err)
	if errors.Is(err, calculo.ErrDivisaoPorZero) || errors.Is(err, calculo.ErrResultadoForaDaFaixa) {
		return CodigoErroDeCalculo
	}
	return CodigoUsoInvalido
}
