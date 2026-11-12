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
			// D-A9 em aberto: compara o valor numérico lido, não o texto.
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
