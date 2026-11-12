# calculadora_go

Projeto de prova do pacote APM [`emilianoeloi/apm`](https://github.com/emilianoeloi/apm) (spec APM_SETUP V1.2).
Serve para exercitar o harness de desenvolvimento do pacote (skills, agentes `sdlc-*`, prompts `/opsx-*` e hooks) num caso pequeno e conhecido: uma calculadora de linha de comando em Go.

Escopo previsto da calculadora: soma, subtração, multiplicação e divisão. A divisão por zero devolve erro, sem pânico.

> **Estado atual:** a Onda 0 preparou o repositório e instalou o harness; a change `calculadora-basica` implementa a calculadora (código em `cmd/calc/` e `internal/`).

## Pré-requisitos

| Ferramenta | Versão |
| --- | --- |
| Go | versão estável atual (testado com 1.27.1) |
| `apm` | 0.35.0 |
| `openspec` | >= 1 (testado com 1.14.1) |
| `python3` | 3.x (testado com 3.14) |

## Instalar o harness

```bash
apm install --target claude,copilot
```

A dependência está fixada em `apm.yml` (`emilianoeloi/apm#v1.0.0-alpha.1`). A instalação implanta skills, agentes, prompts e hooks para os hosts `claude` e `copilot`. Hooks são código executável: revise antes de confiar neles.

## Como rodar

```bash
go run ./cmd/calc <operacao> <a> <b>
# ou compilar o executável:
go build -o calc ./cmd/calc && ./calc <operacao> <a> <b>
```

Operações (nomes sem acento): `soma`, `subtracao`, `multiplicacao`, `divisao`. Os operandos são números `float64`. Uma operação por execução, com exatamente dois operandos.

```bash
go run ./cmd/calc soma 2 3        # 5
go run ./cmd/calc divisao 7 2     # 3.5
go run ./cmd/calc divisao 6 0     # erro: divisão por zero (código 1)
```

O resultado vai para `stdout`. Erros vão para `stderr`, em português, sem pânico.

| Código de saída | Significado |
| --- | --- |
| 0 | Sucesso |
| 1 | Divisão por zero |
| 2 | Uso ou entrada inválida (operação desconhecida, operando não numérico, argumentos ausentes ou em excesso) |

Não informado: formato de impressão do float além do que o Go produz por padrão, tratamento de overflow, `Inf` e `NaN`, separador decimal (hoje só o ponto é aceito).

## Como testar

```bash
go vet ./...
go test ./...
```

`go test ./...` é o `sensors_command` do harness (`harness.config.json`), executado pelo hook de Stop.
