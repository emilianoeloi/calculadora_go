package calculo

import (
	"errors"
	"math"
	"testing"
)

func TestOperacoesSemErro(t *testing.T) {
	casos := []struct {
		nome     string
		operacao func(a, b float64) float64
		a, b     float64
		esperado float64
	}{
		{"Soma de dois números", Somar, 2, 3, 5},
		{"Subtração de dois números", Subtrair, 5, 3, 2},
		{"Multiplicação de dois números", Multiplicar, 4, 3, 12},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if got := c.operacao(c.a, c.b); got != c.esperado {
				t.Errorf("resultado = %v, esperado %v", got, c.esperado)
			}
		})
	}
}

func TestDividir(t *testing.T) {
	casos := []struct {
		nome     string
		a, b     float64
		esperado float64
	}{
		{"Divisão exata", 6, 3, 2},
		{"Divisão com parte fracionária", 7, 2, 3.5},
		{"Dividendo igual a zero com divisor diferente de zero", 0, 3, 0},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			got, err := Dividir(c.a, c.b)
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if got != c.esperado {
				t.Errorf("resultado = %v, esperado %v", got, c.esperado)
			}
		})
	}
}

func TestDividirDivisorZero(t *testing.T) {
	// Um pânico faria o próprio go test falhar; o teste confirma o erro devolvido.
	divisores := map[string]float64{
		"Divisor igual a zero":          0,
		"Divisor igual a zero negativo": math.Copysign(0, -1),
	}
	for nome, b := range divisores {
		t.Run(nome, func(t *testing.T) {
			got, err := Dividir(6, b)
			if !errors.Is(err, ErrDivisaoPorZero) {
				t.Fatalf("erro = %v, esperado ErrDivisaoPorZero", err)
			}
			if got != 0 {
				t.Errorf("resultado = %v, esperado nenhum resultado válido (0)", got)
			}
		})
	}
}
