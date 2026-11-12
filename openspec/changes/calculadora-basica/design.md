# Design

## Context

O repositório não tem código Go: não há `go.mod`, `cmd/` nem `internal/` (estado observado em 2026-11-12). O ambiente tem Go 1.27.1 instalado. As pastas de código permitidas pelo harness são `cmd/` e `internal/`. Os testes rodam com `go test ./...`. Motivação e escopo em `proposal.md`; requisitos em `specs/`. As decisões D-A1 a D-A5 e D-A14 foram respondidas pelo humano e estão registradas abaixo.

## Goals / Non-Goals

**Goals:**
- Separar a lógica de cálculo (testável, sem E/S) da camada de linha de comando.
- Devolver erros como valores (`error`), sem `panic`, em qualquer caminho de falha.
- Cobrir cada cenário das specs com teste automatizado executável por `go test ./...`.

**Non-Goals:**
- Qualquer operação além das quatro pedidas, modo interativo, histórico ou expressões compostas.
- Alterar o `harness.config.json` (a troca do `sensors_command` é feita à parte, pelo humano).
- Dependências externas ao Go padrão (não informado em contrário).

## Decisions

### D1. Duas camadas: pacote de cálculo em `internal/` e entrada em `cmd/`
O cálculo fica em um pacote sob `internal/` e a função `main` em `cmd/calc/main.go`. O `main` só liga os argumentos, `stdout`, `stderr` e o código de saída do processo à camada de CLI. Motivo: a lógica fica testável sem processo externo e a divisão por zero é verificada na camada de cálculo.
Alternativa descartada: tudo em um único `main.go` (dificulta testes e mistura E/S com regra).

### D2. Erros como valor, com erro sentinela para divisão por zero
A divisão devolve `error` quando o divisor é zero. Um erro sentinela exportado permite ao chamador e aos testes identificar o caso com `errors.Is`. Nenhum `panic` nem `recover` é usado.
Alternativa descartada: capturar pânico com `recover` (o harness pede para não capturar exceções de forma ampla, e o pedido exige ausência de pânico).

### D3. Testes
Testes de unidade na camada de cálculo (um por cenário de `operacoes-aritmeticas`) e testes da camada de linha de comando chamando uma função que recebe os argumentos e escritores de saída, para evitar executar um processo externo. A forma exata dessa função fica para o `spdd-reasons-canvas`.

### Decisões resolvidas pelo humano

- **D-A1 (resolvida). Sintaxe dos argumentos**: `calc <operacao> <a> <b>`. A operação é o primeiro argumento, por nome em português e sem acento: `soma`, `subtracao`, `multiplicacao`, `divisao`. Uma operação por execução e exatamente dois operandos; `<a>` é o minuendo/dividendo e `<b>` o subtraendo/divisor. Isto também resolve D-A7 (nomes das operações). Refletida em `specs/cli-calculadora`.
- **D-A2 (resolvida). Saída**: somente o resultado numérico em `stdout` quando há sucesso; mensagens de erro, em português, em `stderr`. O formato exato de impressão do float continua em aberto (D-A9).
- **D-A3 (resolvida). Tipo numérico**: `float64` para operandos e resultado; a divisão não trunca. Overflow, `Inf` e `NaN` continuam em aberto (D-A10). Consequência direta e sinalizada: com `float64`, um divisor `-0` é igual a zero na comparação, então também é tratado como divisão por zero (erro).
- **D-A4 (resolvida). Códigos de saída**: 0 em sucesso; 1 em divisão por zero; 2 em erro de uso ou entrada inválida (operação desconhecida, número de argumentos diferente de três, operando não numérico).
- **D-A5 (resolvida). Módulo e binário**: módulo Go `github.com/emilianoeloi/calculadora_go`; ponto de entrada em `cmd/calc/main.go` (ajustado por D-A14).
- **D-A14 (resolvida). Nome do executável**: `calc`, igual ao usado na sintaxe de D-A1; por isso o ponto de entrada é `cmd/calc/main.go` e não `cmd/calculadora`.

## Risks / Trade-offs

- [`float64` acumula erro de representação (por exemplo 0.1 + 0.2) e pode gerar `Inf`/`NaN`] → Aceito pela decisão do humano; os cenários das specs usam valores exatamente representáveis (2, 3, 5, 6, 7, 3.5, 0). Tratamento de `Inf`/`NaN` fica para D-A10.
- [`strconv.ParseFloat` aceita textos como `NaN`, `Inf` e notação científica, o que pode produzir operandos que o usuário não imaginava] → Registrado em D-A10; sem decisão, o comportamento não é especificado nem testado.
- [Formato de impressão do float indefinido pode gerar retrabalho nos testes de saída] → Fixar D-A9 antes da tarefa 3.1; até lá os testes só podem comparar o valor numérico lido da saída.

## Open Questions

Itens que seguem em aberto (não informado). Nenhum foi suposto; a implementação NÃO deve decidir por conta própria os que afetam o comportamento observável.

- **D-A6. Escopo da execução**: os dois operandos por execução foram assumidos só como escopo na resposta do humano; falta confirmação explícita de que não se quer encadear operações nem aceitar mais operandos.
- **D-A8. Cenários de entrada inválida** (operação desconhecida, operando não numérico, operandos ausentes ou em excesso): mantidos como propostos pelo planejador, sem confirmação final. O código de saída 2 vem de D-A4.
- **D-A9. Formato de impressão do float** em `stdout` (número de casas decimais, notação científica, `5` versus `5.0`): não informado.
- **D-A10. Overflow, `Inf` e `NaN`**: não informado se um resultado `Inf`/`NaN` (por exemplo, multiplicação que estoura) deve ser tratado como erro, nem se operandos `NaN`/`Inf` digitados pelo usuário são aceitos ou rejeitados como entrada inválida.
- **D-A11. Separador decimal** na entrada (ponto, vírgula ou ambos): não informado.
- **D-A12. Versão mínima do Go** (diretiva `go` do `go.mod`): não informado.
- **D-A13. Idioma dos identificadores** no código (português ou inglês): não informado. As mensagens de erro ao usuário são em português (D-A2).
