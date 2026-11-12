// Package calculo implementa as quatro operações aritméticas básicas sobre
// float64. O pacote não faz E/S: erros são devolvidos como valor, nunca como
// pânico.
package calculo

import "errors"

// ErrDivisaoPorZero é devolvido por Dividir quando o divisor é zero
// (inclusive -0). Compare com errors.Is.
var ErrDivisaoPorZero = errors.New("divisão por zero")

// Somar devolve a + b.
func Somar(a, b float64) float64 {
	return a + b
}

// Subtrair devolve a - b (a é o minuendo, b o subtraendo).
func Subtrair(a, b float64) float64 {
	return a - b
}

// Multiplicar devolve a * b.
func Multiplicar(a, b float64) float64 {
	return a * b
}

// Dividir devolve a / b, sem truncar a parte fracionária. Se b for zero
// (inclusive -0), devolve ErrDivisaoPorZero e nenhum resultado válido.
func Dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisaoPorZero
	}
	return a / b, nil
}
