# Proposal

## Why

O repositório ainda não tem nenhum código de produto. Falta uma calculadora de linha de comando em Go que faça as quatro operações aritméticas básicas e trate a divisão por zero como erro explícito, sem derrubar o programa. Esta change também serve de primeira funcionalidade para ativar `go test ./...` como comando de sensores do harness.

## What Changes

- Adicionar a lógica de cálculo com as operações soma, subtração, multiplicação e divisão.
- A divisão por zero passa a devolver um erro, sem `panic`.
- Adicionar uma interface de linha de comando na forma `calc <operacao> <a> <b>` (operação por nome em português sem acento: `soma`, `subtracao`, `multiplicacao`, `divisao`), que escreve só o resultado em stdout (código 0) ou uma mensagem de erro em português em stderr (código 1 na divisão por zero, 2 em uso/entrada inválida).
- Operandos e resultado em `float64`; a divisão não trunca.
- Organizar o código em `cmd/calc/` (ponto de entrada) e `internal/` (lógica), no módulo `github.com/emilianoeloi/calculadora_go`, e cobri-lo com testes executáveis por `go test ./...`.
- As decisões D-A1 a D-A5 e D-A14 foram respondidas pelo humano e estão registradas no `design.md`. Seguem em aberto, sem suposição: confirmação de D-A6 e D-A8, formato de impressão do float, overflow/`Inf`/`NaN`, separador decimal, versão mínima do Go, e idioma dos identificadores.

## Capabilities

### New Capabilities
- `operacoes-aritmeticas`: cálculo de soma, subtração, multiplicação e divisão sobre dois operandos, com erro (e não pânico) na divisão por zero.
- `cli-calculadora`: execução da calculadora pela linha de comando, incluindo a entrega do resultado ou do erro ao usuário.

### Modified Capabilities
<!-- Nenhuma: `openspec/specs/` está vazio. -->

## Impact

- Código novo em `cmd/calc/` e `internal/` (hoje inexistentes). Não há `go.mod`; a criação dele faz parte da implementação, com o módulo `github.com/emilianoeloi/calculadora_go`.
- Sem dependências externas previstas além da biblioteca padrão do Go (não informado em contrário).
- Harness: `go test ./...` será o `sensors_command`, mas o `harness.config.json` NÃO é alterado por esta change de planejamento. A atualização é feita à parte, pelo humano.
- Fora de escopo (não pedido): operações além das quatro, modo interativo, histórico, expressões compostas, números complexos.
