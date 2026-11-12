# calculadora_go

Projeto de prova do pacote APM [`emilianoeloi/apm`](https://github.com/emilianoeloi/apm) (spec APM_SETUP V1.2).
Serve para exercitar o harness de desenvolvimento do pacote (skills, agentes `sdlc-*`, prompts `/opsx-*` e hooks) num caso pequeno e conhecido: uma calculadora de linha de comando em Go.

Escopo previsto da calculadora: soma, subtração, multiplicação e divisão. A divisão por zero devolve erro, sem pânico.

> **Estado atual (Onda 0):** a calculadora ainda **não existe**. Esta onda só prepara o repositório e instala o harness. O código entra nas ondas seguintes, pelo fluxo do harness.

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

Não há o que rodar ainda: a calculadora não existe. Quando existir, o ponto de entrada ficará em `cmd/` e a lógica em `internal/`.

## Como testar

Também não há testes ainda. Com a primeira funcionalidade, `go test ./...` passa a ser o comando de sensores do harness (`sensors_command` em `harness.config.json`). Até lá o campo fica vazio (não informado), porque `go test ./...` sem `go.mod` travaria o hook de Stop.
