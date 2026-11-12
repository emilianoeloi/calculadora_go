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
go run ./cmd/calc soma 1e308 1e308  # erro: resultado fora da faixa (código 1)
go run ./cmd/calc multiplicacao -1 0 # 0
go run ./cmd/calc soma 0.1 0.2    # 0.30000000000000004
```

O resultado vai para `stdout`. Erros vão para `stderr`, em português, sem pânico.

| Código de saída | Significado |
| --- | --- |
| 0 | Sucesso |
| 1 | Erro de cálculo: divisão por zero, ou resultado fora da faixa de `float64` (overflow) |
| 2 | Uso ou entrada inválida: operação desconhecida, operando não numérico, `nan`, `inf`, `Infinity`, hexadecimais (`0x1p-2`), sublinhados (`1_000`), vírgula decimal, argumentos ausentes ou em excesso, ou operando fora da faixa de `float64` (`1e999`) |

Os operandos aceitos são decimais finitos com ponto (`+5`, `.5`, `5.`, `1e3`). Infinitos e `NaN` não são aceitos como entrada nem produzidos como resultado. O resultado é impresso no menor texto que reconstrói o valor (`formatar` com `'g'` e precisão `-1`): `5`, `0.30000000000000004`, `1e+22`; zero negativo sai como `0`. Por consequência do formato `'g'`, magnitudes grandes saem em notação científica: `soma 500000 500000` imprime `1e+06`.

Não informado: comportamento de underflow (por exemplo `soma 1e-999 0` imprime `0` sem erro) e o texto exato das mensagens de erro.

## Como testar

```bash
go vet ./...
go test ./...
```

`go test ./...` é o `sensors_command` do harness (`harness.config.json`), executado pelo hook de Stop.
