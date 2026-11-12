package cli

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func executar(args ...string) (stdout, stderr string, codigo int) {
	var out, errOut bytes.Buffer
	codigo = Executar(args, &out, &errOut)
	return out.String(), errOut.String(), codigo
}

func TestExecutarSucesso(t *testing.T) {
	casos := []struct {
		nome     string
		args     []string
		esperado float64
	}{
		{"Operação válida", []string{"soma", "2", "3"}, 5},
		{"Ordem dos operandos", []string{"subtracao", "5", "3"}, 2},
		{"Divisão sem truncamento", []string{"divisao", "7", "2"}, 3.5},
		{"Multiplicação", []string{"multiplicacao", "4", "3"}, 12},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			stdout, stderr, codigo := executar(c.args...)
			if codigo != 0 {
				t.Errorf("código = %d, esperado 0", codigo)
			}
			if stderr != "" {
				t.Errorf("stderr = %q, esperado vazio", stderr)
			}
			// Compara o valor numérico lido; o texto exato está em TestExecutarFormatoDaSaida.
			got, err := strconv.ParseFloat(strings.TrimSpace(stdout), 64)
			if err != nil {
				t.Fatalf("stdout %q não é um número: %v", stdout, err)
			}
			if got != c.esperado {
				t.Errorf("stdout = %v, esperado %v", got, c.esperado)
			}
		})
	}
}

func TestExecutarDivisaoPorZero(t *testing.T) {
	casos := []struct {
		nome string
		args []string
	}{
		{"Divisor zero", []string{"divisao", "6", "0"}},
		{"Divisor zero negativo", []string{"divisao", "6", "-0"}},
		{"Zero dividido por zero", []string{"divisao", "0", "0"}},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			stdout, stderr, codigo := executar(c.args...)
			if codigo != 1 {
				t.Errorf("código = %d, esperado 1", codigo)
			}
			if stdout != "" {
				t.Errorf("stdout = %q, esperado vazio", stdout)
			}
			if !strings.Contains(stderr, "divisão por zero") {
				t.Errorf("stderr = %q, esperada mensagem de divisão por zero", stderr)
			}
		})
	}
}

func TestExecutarOperandoForaDaFaixa(t *testing.T) {
	stdout, stderr, codigo := executar("soma", "1e999", "1")
	if codigo != 2 {
		t.Errorf("código = %d, esperado 2", codigo)
	}
	if stdout != "" {
		t.Errorf("stdout = %q, esperado vazio", stdout)
	}
	if !strings.Contains(stderr, "fora da faixa") {
		t.Errorf("stderr = %q, esperada mensagem de operando fora da faixa", stderr)
	}
}

func TestExecutarEntradaInvalida(t *testing.T) {
	casos := []struct {
		nome string
		args []string
	}{
		{"Operação desconhecida", []string{"potencia", "2", "3"}},
		{"Operando que não é número", []string{"soma", "abc", "3"}},
		{"Segundo operando que não é número", []string{"soma", "2", "abc"}},
		{"Operandos ausentes", []string{"soma", "2"}},
		{"Argumentos em excesso", []string{"soma", "2", "3", "4"}},
		{"Sem argumentos", nil},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			stdout, stderr, codigo := executar(c.args...)
			if codigo != 2 {
				t.Errorf("código = %d, esperado 2", codigo)
			}
			if stdout != "" {
				t.Errorf("stdout = %q, esperado vazio", stdout)
			}
			if stderr == "" {
				t.Error("stderr vazio, esperada mensagem de erro em português")
			}
		})
	}
}

// entradasEspeciais são operandos que a política "só finitos" rejeita: são
// usados para provar código 2, stdout vazio e ausência de pânico.
var entradasEspeciais = []string{
	"nan", "NaN", "+nan", "-NAN",
	"inf", "-inf", "Infinity", "-INF", "+Infinity",
	"0x1p-2", "-0x1p-2",
	"1_000",
	"1,5",
	"1e999", "-1e999",
}

func TestExecutarOperandoEspecialRejeitado(t *testing.T) {
	casos := []struct {
		nome string
		args []string
	}{
		{"NaN como operando", []string{"divisao", "1", "nan"}},
		{"NaN com sinal ou em outra caixa (NaN)", []string{"soma", "NaN", "1"}},
		{"NaN com sinal ou em outra caixa (+nan)", []string{"soma", "+nan", "1"}},
		{"NaN com sinal ou em outra caixa (-NAN)", []string{"soma", "-NAN", "1"}},
		{"NaN como dividendo com divisor zero", []string{"divisao", "nan", "0"}},
		{"Infinitos como operandos", []string{"soma", "inf", "-inf"}},
		{"Infinity por extenso", []string{"soma", "Infinity", "1"}},
		{"Infinito com sinal em outra caixa (-INF)", []string{"soma", "-INF", "1"}},
		{"Infinito com sinal em outra caixa (+Infinity)", []string{"soma", "+Infinity", "1"}},
		{"Hexadecimal como operando", []string{"soma", "0x1p-2", "1"}},
		{"Hexadecimal com sinal", []string{"soma", "-0x1p-2", "1"}},
		{"Hexadecimal com sinal de mais", []string{"soma", "+0x1p-2", "1"}},
		{"Sublinhado como operando", []string{"soma", "1_000", "1"}},
		{"Sublinhado no segundo operando", []string{"soma", "1", "1_000"}},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			stdout, stderr, codigo := executar(c.args...)
			if codigo != 2 {
				t.Errorf("código = %d, esperado 2", codigo)
			}
			if stdout != "" {
				t.Errorf("stdout = %q, esperado vazio", stdout)
			}
			if stderr == "" {
				t.Error("stderr vazio, esperada mensagem de erro em português")
			}
			if strings.Contains(stderr, "divisão por zero") {
				t.Errorf("stderr = %q, operando rejeitado não pode ser divisão por zero", stderr)
			}
		})
	}
}

func TestExecutarRegressoesDeEntrada(t *testing.T) {
	t.Run("Vírgula como separador decimal", func(t *testing.T) {
		stdout, stderr, codigo := executar("soma", "1,5", "2")
		if codigo != 2 {
			t.Errorf("código = %d, esperado 2", codigo)
		}
		if stdout != "" {
			t.Errorf("stdout = %q, esperado vazio", stdout)
		}
		if stderr == "" {
			t.Error("stderr vazio, esperada mensagem de erro em português")
		}
	})
	for _, texto := range []string{"1e999", "-1e999"} {
		t.Run("Operando decimal fora da faixa "+texto, func(t *testing.T) {
			stdout, stderr, codigo := executar("soma", texto, "1")
			if codigo != 2 {
				t.Errorf("código = %d, esperado 2", codigo)
			}
			if stdout != "" {
				t.Errorf("stdout = %q, esperado vazio", stdout)
			}
			if !strings.Contains(stderr, "fora da faixa") {
				t.Errorf("stderr = %q, esperada mensagem de operando fora da faixa", stderr)
			}
		})
	}
	decimaisFinitos := []struct {
		nome     string
		args     []string
		esperado string
	}{
		{"Decimais finitos com expoente e sinal", []string{"soma", "1e3", "-2.5"}, "997.5\n"},
		{"Sinal de mais aceito", []string{"soma", "+5", "0"}, "5\n"},
		{"Decimal sem parte inteira aceito", []string{"soma", ".5", "0"}, "0.5\n"},
		{"Decimal sem parte fracionária aceito", []string{"soma", "5.", "0"}, "5\n"},
	}
	for _, c := range decimaisFinitos {
		t.Run(c.nome, func(t *testing.T) {
			stdout, stderr, codigo := executar(c.args...)
			if codigo != 0 {
				t.Errorf("código = %d, esperado 0", codigo)
			}
			if stderr != "" {
				t.Errorf("stderr = %q, esperado vazio", stderr)
			}
			if stdout != c.esperado {
				t.Errorf("stdout = %q, esperado %q", stdout, c.esperado)
			}
		})
	}
}

func TestExecutarNaoEntraEmPanico(t *testing.T) {
	// Um pânico derrubaria o teste; cada combinação deve retornar com código 2.
	operacoes := []string{"soma", "subtracao", "multiplicacao", "divisao"}
	for _, op := range operacoes {
		for _, especial := range entradasEspeciais {
			for _, args := range [][]string{
				{op, especial, "1"},
				{op, "1", especial},
				{op, especial, "0"},
			} {
				t.Run(strings.Join(args, " "), func(t *testing.T) {
					stdout, _, codigo := executar(args...)
					if codigo != 2 {
						t.Errorf("código = %d, esperado 2", codigo)
					}
					if stdout != "" {
						t.Errorf("stdout = %q, esperado vazio", stdout)
					}
				})
			}
		}
	}
}

func TestExecutarResultadoForaDaFaixa(t *testing.T) {
	casos := []struct {
		nome string
		args []string
	}{
		{"Soma que estoura a faixa", []string{"soma", "1e308", "1e308"}},
		{"Multiplicação que estoura a faixa", []string{"multiplicacao", "1e200", "1e200"}},
		{"Divisão que estoura a faixa", []string{"divisao", "1e308", "1e-308"}},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			stdout, stderr, codigo := executar(c.args...)
			if codigo != 1 {
				t.Errorf("código = %d, esperado 1", codigo)
			}
			if stdout != "" {
				t.Errorf("stdout = %q, esperado vazio", stdout)
			}
			if !strings.Contains(stderr, "fora da faixa") {
				t.Errorf("stderr = %q, esperada mensagem de resultado fora da faixa", stderr)
			}
			if strings.Contains(stdout+stderr, "Inf") || strings.Contains(stdout+stderr, "NaN") {
				t.Errorf("saída %q não pode conter Inf nem NaN", stdout+stderr)
			}
		})
	}
}

func TestExecutarFormatoDaSaida(t *testing.T) {
	casos := []struct {
		nome     string
		args     []string
		esperado string
	}{
		{"Resultado inteiro sem parte fracionária", []string{"soma", "2", "3"}, "5\n"},
		{"Menor texto que reconstrói o valor", []string{"soma", "0.1", "0.2"}, "0.30000000000000004\n"},
		{"Notação científica para magnitudes grandes", []string{"soma", "1e22", "0"}, "1e+22\n"},
		{"Zero negativo produzido por multiplicação", []string{"multiplicacao", "-1", "0"}, "0\n"},
		{"Zero negativo digitado como operando", []string{"soma", "-0", "-0"}, "0\n"},
		{"Zero negativo como dividendo", []string{"divisao", "-0", "5"}, "0\n"},
		{"Zero negativo por subtração", []string{"subtracao", "-0", "0"}, "0\n"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			stdout, stderr, codigo := executar(c.args...)
			if codigo != 0 {
				t.Errorf("código = %d, esperado 0", codigo)
			}
			if stderr != "" {
				t.Errorf("stderr = %q, esperado vazio", stderr)
			}
			if stdout != c.esperado {
				t.Errorf("stdout = %q, esperado %q", stdout, c.esperado)
			}
		})
	}
}
