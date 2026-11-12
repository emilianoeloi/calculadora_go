// Comando calc: calculadora de linha de comando.
//
// Uso: calc <operacao> <a> <b>, com operacao em soma, subtracao,
// multiplicacao ou divisao.
package main

import (
	"os"

	"github.com/emilianoeloi/calculadora_go/internal/cli"
)

func main() {
	os.Exit(cli.Executar(os.Args[1:], os.Stdout, os.Stderr))
}
