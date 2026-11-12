package calculo

import (
	"errors"
	"math"
	"testing"
)

func TestOperacoesSemErro(t *testing.T) {
	casos := []struct {
		nome     string
		operacao func(a, b float64) (float64, error)
		a, b     float64
		esperado float64
	}{
		{"Soma de dois números", Somar, 2, 3, 5},
		{"Subtração de dois números", Subtrair, 5, 3, 2},
		{"Multiplicação de dois números", Multiplicar, 4, 3, 12},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			got, err := c.operacao(c.a, c.b)
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if got != c.esperado {
				t.Errorf("resultado = %v, esperado %v", got, c.esperado)
			}
		})
	}
}

func TestSomarResultadoGrandeMasFinito(t *testing.T) {
	got, err := Somar(1e308, 1e307)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if math.IsInf(got, 0) || math.IsNaN(got) {
		t.Errorf("resultado = %v, esperado finito", got)
	}
}

func TestOperacoesResultadoForaDaFaixa(t *testing.T) {
	casos := []struct {
		nome     string
		operacao func(a, b float64) (float64, error)
		a, b     float64
	}{
		{"Soma que estoura a faixa", Somar, 1e308, 1e308},
		{"Multiplicação que estoura a faixa", Multiplicar, 1e200, 1e200},
		{"Subtração que estoura a faixa para o lado negativo", Subtrair, -1e308, 1e308},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			got, err := c.operacao(c.a, c.b)
			if !errors.Is(err, ErrResultadoForaDaFaixa) {
				t.Fatalf("erro = %v, esperado ErrResultadoForaDaFaixa", err)
			}
			if got != 0 {
				t.Errorf("resultado = %v, esperado nenhum resultado válido (0)", got)
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
	casos := []struct {
		nome string
		a, b float64
	}{
		{"Divisor igual a zero", 6, 0},
		{"Divisor igual a zero negativo", 6, math.Copysign(0, -1)},
		{"Zero dividido por zero", 0, 0},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			got, err := Dividir(c.a, c.b)
			if !errors.Is(err, ErrDivisaoPorZero) {
				t.Fatalf("erro = %v, esperado ErrDivisaoPorZero", err)
			}
			if errors.Is(err, ErrResultadoForaDaFaixa) {
				t.Errorf("erro = %v, não deveria ser ErrResultadoForaDaFaixa", err)
			}
			if got != 0 {
				t.Errorf("resultado = %v, esperado nenhum resultado válido (0)", got)
			}
		})
	}
}

func TestDividirResultadoForaDaFaixa(t *testing.T) {
	got, err := Dividir(1e308, 1e-308)
	if !errors.Is(err, ErrResultadoForaDaFaixa) {
		t.Fatalf("erro = %v, esperado ErrResultadoForaDaFaixa", err)
	}
	if errors.Is(err, ErrDivisaoPorZero) {
		t.Errorf("erro = %v, não deveria ser ErrDivisaoPorZero", err)
	}
	if got != 0 {
		t.Errorf("resultado = %v, esperado nenhum resultado válido (0)", got)
	}
}
