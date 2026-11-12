# Tasks

## 1. Pré-requisitos e estrutura

- [x] 1.1 Conferir no `design.md` o registro das decisões resolvidas D-A1 a D-A5 e D-A14 (sintaxe, saída, float64, códigos de saída, módulo, binário `calc` em `cmd/calc/`) e verificar que os itens ainda em aberto (D-A6, D-A8 a D-A13) continuam listados em Open Questions; se o humano resolver algum deles antes da implementação, atualizar o `design.md` e as specs correspondentes
- [x] 1.2 Criar o `go.mod` com o módulo `github.com/emilianoeloi/calculadora_go` e as pastas `cmd/calc/` e `internal/`, e verificar que `go build ./...` e `go test ./...` rodam sem erro de configuração (versão do Go na diretiva `go`: D-A12, não informada)

## 2. Operações aritméticas (spec `operacoes-aritmeticas`)

- [x] 2.1 Implementar soma, subtração e multiplicação sobre `float64` em `internal/` e verificar com testes de unidade que cobrem os cenários "Soma de dois números", "Subtração de dois números" e "Multiplicação de dois números" (`go test ./...` passa)
- [x] 2.2 Implementar a divisão sobre `float64`, sem truncamento, com erro sentinela para divisor zero (inclusive `-0`), sem `panic`, e verificar com testes de unidade que cobrem "Divisão exata", "Divisão com parte fracionária", "Divisor igual a zero" (com `errors.Is`) e "Dividendo igual a zero com divisor diferente de zero" (`go test ./...` passa)

## 3. Linha de comando (spec `cli-calculadora`)

- [x] 3.1 Implementar em `internal/` a função que recebe os argumentos no formato `<operacao> <a> <b>` e os escritores de `stdout` e `stderr`, escreve só o resultado em `stdout` com código 0, escreve em português em `stderr` na divisão por zero com código 1, e verificar com testes que cobrem "Operação válida", "Ordem dos operandos", "Divisão sem truncamento" e "Divisão por zero pela linha de comando", conferindo `stdout`, `stderr` e o código de saída (o formato de impressão do float segue D-A9: comparar o valor numérico lido, não o texto exato)
- [x] 3.2 Tratar operação desconhecida, operando não numérico, operandos ausentes e argumentos em excesso sem pânico, com mensagem em português em `stderr` e código 2 (D-A8, a confirmar), e verificar com testes que cobrem os quatro cenários de entrada inválida
- [x] 3.3 Criar o ponto de entrada `cmd/calc/main.go` apenas chamando a função de 3.1 com os argumentos do processo, `os.Stdout`, `os.Stderr` e encerrando com o código devolvido, e verificar executando `go run ./cmd/calc soma 2 3` e `go run ./cmd/calc divisao 6 0`, conferindo a saída, o código 0 e o código 1 e a ausência de pânico

## 4. Integração

- [x] 4.1 Rodar `go vet ./...` e `go test ./...` e verificar que ambos terminam com código 0
- [x] 4.2 Documentar no `README.md` o uso da calculadora (sintaxe `<operacao> <a> <b>`, nomes das operações, códigos de saída e exemplos; executável `calc`, D-A14) e verificar que um exemplo documentado roda como escrito

## Workflow follow-up

- O humano atualiza `sensors_command` no `harness.config.json` para `go test ./...` e o texto de "Quality gates mínimos" no `CLAUDE.md` (não faz parte desta change).
- Arquivar a change depois que a revisão for aprovada (`openspec archive calculadora-basica`).
