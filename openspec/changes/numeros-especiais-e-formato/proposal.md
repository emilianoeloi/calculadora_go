# Proposal

## Why

A change arquivada `calculadora-basica` deixou em aberto D-A9 (formato de impressão do float), D-A10 (overflow, `Inf` e `NaN`) e D-A11 (separador decimal). O comportamento atual, observado em 2026-11-12 no código e executando o binário, é surpreendente: `strconv.ParseFloat` aceita `nan`, `inf`, `infinity` e hexadecimais com expoente `p` (como `0x1p-2`) e sublinhados (como `1_000`); `soma nan 1` imprime `NaN` com código 0; `soma 1e308 1e308` imprime `+Inf` com código 0; `multiplicacao -1 0` imprime `-0`; e `divisao nan 0` sai com "divisão por zero" (código 1) em vez de acusar a entrada inválida. O humano decidiu a política para esses casos; esta change a especifica e a torna testável.

## What Changes

- **Política "só finitos" na entrada (D-A10).** O programa passa a rejeitar `nan`, `inf` e `infinity` (sem diferenciar maiúsculas, com ou sem sinal), hexadecimais (como `0x1p-2`) e sublinhados (como `1_000`), com mensagem em português em stderr, stdout vazio e código de saída 2 (o mesmo de "operando não é um número válido").
- **Resultado fora da faixa vira erro de cálculo (D-A10).** Quando os operandos são finitos mas o resultado estoura para infinito (`soma 1e308 1e308`, `multiplicacao 1e200 1e200`, `divisao 1e308 1e-308`), a operação devolve erro, o programa escreve "resultado fora da faixa" em stderr, nada em stdout, e sai com código 1 (o mesmo da divisão por zero). `NaN` como resultado deixa de ser possível.
- **Operando fora da faixa continua código 2.** `soma 1e999 1` segue como erro de entrada, com a mensagem de "fora da faixa" (sem mudança).
- **Formato de impressão mantido (D-A9).** O menor texto decimal que reconstrói o `float64` (`strconv.FormatFloat` com `'g'` e precisão `-1`): inteiro imprime `5`, `soma 0.1 0.2` imprime `0.30000000000000004`, `soma 1e22 0` imprime `1e+22`; uma linha com quebra de linha final.
- **Zero negativo normalizado na saída (D-A9b).** `multiplicacao -1 0` e `soma -0 -0` passam a imprimir `0`, não `-0`.
- **Separador decimal só o ponto (D-A11).** Na entrada e na saída; `soma 1,5 2` continua sendo erro com código 2.
- **Divisão por zero reafirmada.** `divisao 6 -0` e `divisao 0 0` ficam explícitos como cenários (erro, código 1, stdout vazio); o comportamento já existe e já tem teste, mas não estava nas specs.
- Mudanças de comportamento observáveis em relação a hoje: entradas `nan`/`inf`/`infinity`/hexadecimal/sublinhado passam de código 0 (ou 1, no caso de `divisao nan 0`) para código 2; overflow passa de código 0 para código 1; `-0` impresso passa a `0`.
- Sem alteração (herdadas da change arquivada): D-A6 (uma operação e dois operandos por execução), D-A8 (entrada inválida com código 2), D-A12 (versão mínima do Go) e D-A13 (idioma dos identificadores).
- Em aberto, sem decisão: comportamento de underflow (ver Open Questions no `design.md`). O comportamento atual é mantido.

## Capabilities

### New Capabilities
<!-- Nenhuma: as regras novas cabem nas capabilities existentes. -->

### Modified Capabilities
- `operacoes-aritmeticas`: divisão por zero ganha os cenários `-0` e `0/0`; novo requisito de resultado não finito (overflow) devolvido como erro.
- `cli-calculadora`: novos requisitos de operandos numéricos aceitos (só finitos, só ponto), formato de impressão do resultado (inclui zero negativo) e resultado fora da faixa informado ao usuário; a entrada inválida e a divisão por zero ganham cenários.

## Impact

- Código a alterar na fase `apply`: `internal/calculo/calculo.go` (erro sentinela de faixa e verificação de resultado finito), `internal/cli/cli.go` (`lerOperando`, `formatar`, `reportar`) e os testes `internal/calculo/calculo_test.go` e `internal/cli/cli_test.go`. Se as operações passarem a devolver `error`, os chamadores e testes existentes do pacote `calculo` mudam de assinatura.
- Documentação: `README.md` (seção que lista como "não informados" o formato, `Inf`/`NaN` e o separador decimal) e `CHANGELOG.md`.
- Sem dependências novas além da biblioteca padrão do Go. O `harness.config.json` NÃO é alterado.
- Fora de escopo: novas operações, modo interativo, aceitar vírgula, separadores de milhar, notação hexadecimal, tratamento de underflow, mudança do formato de impressão.
