// Package cli interpreta os argumentos de `calc <operacao> <a> <b>`, chama o
// pacote calculo e apresenta o resultado em stdout ou o erro em stderr.
//
// Códigos de saída: 0 sucesso; 1 divisão por zero; 2 erro de uso ou entrada
// inválida.
package cli

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/emilianoeloi/calculadora_go/internal/calculo"
)

// Códigos de saída do programa.
const (
	CodigoSucesso          = 0
	CodigoDivisaoPorZero   = 1
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

func lerOperando(texto string) (float64, error) {
	valor, err := strconv.ParseFloat(texto, 64)
	if err != nil {
		return 0, fmt.Errorf("operando %q não é um número válido", texto)
	}
	return valor, nil
}

// calcular assume operação e operandos já validados por interpretar.
func calcular(operacao string, a, b float64) (float64, error) {
	switch operacao {
	case "soma":
		return calculo.Somar(a, b), nil
	case "subtracao":
		return calculo.Subtrair(a, b), nil
	case "multiplicacao":
		return calculo.Multiplicar(a, b), nil
	default: // "divisao"
		return calculo.Dividir(a, b)
	}
}

// formatar imprime o menor texto que representa o float64 sem perda.
func formatar(valor float64) string {
	return strconv.FormatFloat(valor, 'g', -1, 64)
}

// reportar é o único ponto que converte erro em mensagem em stderr e em
// código de saída.
func reportar(stderr io.Writer, err error) int {
	fmt.Fprintf(stderr, "erro: %v\n", err)
	if errors.Is(err, calculo.ErrDivisaoPorZero) {
		return CodigoDivisaoPorZero
	}
	return CodigoUsoInvalido
}
