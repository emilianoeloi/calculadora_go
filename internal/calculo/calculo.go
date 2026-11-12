// Package calculo implementa as quatro operações aritméticas básicas sobre
// float64. O pacote não faz E/S: erros são devolvidos como valor, nunca como
// pânico. Nenhum resultado não finito (Inf ou NaN) é devolvido como válido.
package calculo

import (
	"errors"
	"math"
)

// ErrDivisaoPorZero é devolvido por Dividir quando o divisor é zero
// (inclusive -0). Compare com errors.Is.
var ErrDivisaoPorZero = errors.New("divisão por zero")

// ErrResultadoForaDaFaixa é devolvido quando o resultado de uma operação
// com operandos finitos não é finito (overflow para infinito). Compare com
// errors.Is.
var ErrResultadoForaDaFaixa = errors.New("resultado fora da faixa")

// Somar devolve a + b, ou ErrResultadoForaDaFaixa se o resultado não for
// finito.
func Somar(a, b float64) (float64, error) {
	return verificarFinito(a + b)
}

// Subtrair devolve a - b (a é o minuendo, b o subtraendo), ou
// ErrResultadoForaDaFaixa se o resultado não for finito.
func Subtrair(a, b float64) (float64, error) {
	return verificarFinito(a - b)
}

// Multiplicar devolve a * b, ou ErrResultadoForaDaFaixa se o resultado não
// for finito.
func Multiplicar(a, b float64) (float64, error) {
	return verificarFinito(a * b)
}

// Dividir devolve a / b, sem truncar a parte fracionária. Se b for zero
// (inclusive -0), devolve ErrDivisaoPorZero, verificado antes do quociente.
// Se o quociente não for finito, devolve ErrResultadoForaDaFaixa.
func Dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisaoPorZero
	}
	return verificarFinito(a / b)
}

// verificarFinito devolve r sem erro quando r é finito; caso contrário
// devolve 0 e ErrResultadoForaDaFaixa, sem resultado numérico válido.
func verificarFinito(r float64) (float64, error) {
	if math.IsInf(r, 0) || math.IsNaN(r) {
		return 0, ErrResultadoForaDaFaixa
	}
	return r, nil
}
