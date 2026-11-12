# SPDD Analysis: Números especiais, faixa e formato de impressão da calculadora

## Original Business Requirement

Pedido original (verbatim), recebido pela sessão principal:

> Planeje a change `numeros-especiais-e-formato` em /Users/emiliano.barbosa/Pessoal/calculadora_go (branch feat/numeros-especiais-e-formato; fase `proposal` já declarada; o esqueleto openspec/changes/numeros-especiais-e-formato/.openspec.yaml já existe). Entregue proposta, specs (por deltas sobre as capabilities existentes `operacoes-aritmeticas` e `cli-calculadora` em openspec/specs/; nova capability só se justificar), design.md e tasks.md em português, a análise SPDD em spdd/analysis/ e o REASONS Canvas em spdd/prompt/ (nomes no padrão da change anterior, ticket GGQPA-XXX). Valide com `openspec validate numeros-especiais-e-formato --strict`. Não escreva código, não altere harness.config.json, não faça git add/commit. Se algum gate bloquear, não contorne: pare e relate.

Decisões do humano já tomadas (registradas como decididas em `design.md`):

> - D-A10, política "só finitos": a entrada rejeita `nan`, `inf`, `infinity` (sem diferenciar maiúsculas, com ou sem sinal), hexadecimais (como `0x1p-2`) e sublinhados (como `1_000`) com mensagem em português em stderr, stdout vazio e código 2 (o mesmo de "operando não é um número válido"). Resultado infinito por overflow com operandos finitos (`soma 1e308 1e308`, `multiplicacao 1e200 1e200`, `divisao 1e308 1e-308`) vira erro "resultado fora da faixa", stdout vazio, código 1 (erro de cálculo, como a divisão por zero). NaN como resultado deixa de ser possível. A rejeição de `1e999` como fora da faixa (código 2) continua.
> - D-A9, formato de impressão: manter o atual (`strconv.FormatFloat` com `'g'` e `-1`; inteiro imprime `5`; uma linha com quebra de linha final).
> - D-A9b: zero negativo normalizado para `0` na saída (`multiplicacao -1 0` imprime `0`, não `-0`).
> - D-A11: separador decimal só o ponto, na entrada e na saída (`1,5` continua erro com código 2).
> - Em aberto, não decidido: underflow (`soma 1e-999 0` hoje imprime `0` sem erro); manter o comportamento atual e registrar a pergunta.
> - Fora do escopo, não alteradas: D-A6, D-A8, D-A12 e D-A13 da change arquivada.

Artefatos OpenSpec relacionados (leitura obrigatória para o canvas): `openspec/changes/numeros-especiais-e-formato/{proposal,design,tasks}.md` e `specs/{operacoes-aritmeticas,cli-calculadora}/spec.md` (deltas sobre `openspec/specs/`).

Critérios de aceite (derivados do pedido, numerados para a análise):
- AC1: `nan`, `inf` e `infinity` (qualquer caixa, com ou sem sinal) são rejeitados como operando, código 2, stdout vazio.
- AC2: hexadecimais (`0x1p-2`) e sublinhados (`1_000`) são rejeitados como operando, código 2.
- AC3: vírgula decimal (`1,5`) continua erro, código 2.
- AC4: operando decimal fora da faixa (`1e999`) continua erro "fora da faixa", código 2.
- AC5: overflow de resultado com operandos finitos vira erro "resultado fora da faixa", código 1, stdout vazio.
- AC6: `NaN` e `Inf` nunca são impressos em stdout.
- AC7: divisão por zero (`6 0`, `6 -0`, `0 0`) é erro, código 1, stdout vazio.
- AC8: o formato de impressão atual é mantido (`0.30000000000000004`, `1e+22`, `5`).
- AC9: zero negativo é impresso como `0`.
- AC10: nenhum caminho de código causa pânico.
- AC11: o underflow permanece como está e vira pergunta em aberto.
- AC12: o `harness.config.json` não é alterado.

## Domain Concept Identification

Exploração do código (2026-11-12): `internal/cli/cli.go` tem `Executar`, `interpretar`, `lerOperando` (chama `strconv.ParseFloat(texto, 64)`; `ErrRange` vira "fora da faixa", outro erro vira "não é um número válido"; ambos código 2), `calcular`, `formatar` (`strconv.FormatFloat(valor, 'g', -1, 64)`) e `reportar` (único ponto erro → mensagem e código: `ErrDivisaoPorZero` → 1, resto → 2). `internal/calculo/calculo.go` tem `Somar`, `Subtrair`, `Multiplicar` (devolvem só `float64`) e `Dividir` (devolve `ErrDivisaoPorZero` quando `b == 0`, inclusive `-0`). Os testes em `internal/cli/cli_test.go` e `internal/calculo/calculo_test.go` já cobrem `divisao 6 0`, `6 -0`, `0 0` e `1e999`.

Comportamento atual reconfirmado executando o binário: `soma nan 1` imprime `NaN`; `soma -Inf 1` imprime `-Inf`; `soma inf -inf` imprime `NaN`; `soma Infinity 1` imprime `+Inf`; `soma 0x1p-2 1` imprime `1.25`; `soma 1_000 1` imprime `1001`; `divisao 1 nan` imprime `NaN`; todos com código 0. `soma +nan 1` e `soma -nan 1` já saem com código 2 (o `ParseFloat` não aceita `nan` com sinal), mas `nan` sem sinal é aceito. `divisao nan 0` sai com "divisão por zero" e código 1. `soma 1e999 1` sai com código 2, "fora da faixa". `soma 1e-999 0` imprime `0`, código 0. `soma 1e308 1e308`, `multiplicacao 1e200 1e200` e `divisao 1e308 1e-308` imprimem `+Inf`, código 0. `multiplicacao -1 0` e `soma -0 -0` imprimem `-0`. `soma 0.1 0.2` imprime `0.30000000000000004`; `soma 1e22 0` imprime `1e+22`; `soma 500000 500000` imprime `1e+06`; `soma 1,5 2` sai com código 2.

### Existing Concepts (from codebase)
- Operando: `float64` lido de argumento por `lerOperando`, hoje com a gramática inteira do `ParseFloat` (decimal, hexadecimal com `p`, sublinhados, `nan`, `inf`, `infinity`).
- Erro de operando (`erroDeOperando`): mensagem em português com a causa original embrulhada; código 2.
- Erro sentinela `ErrDivisaoPorZero`: único erro de cálculo hoje; código 1.
- Ponto único de conversão erro → stderr/código (`reportar`).
- Formatação de saída (`formatar`): menor texto que reconstrói o `float64`.
- Códigos de saída 0 / 1 / 2.

### New Concepts Required
- Número finito: único tipo de operando e de resultado válido; `NaN` e `±Inf` nunca entram nem saem.
- Resultado fora da faixa (erro de cálculo): resultado não finito produzido por overflow com operandos finitos; sentinela próprio (`ErrResultadoForaDaFaixa`, nome sugerido), distinto da divisão por zero; código 1.
- Operando especial rejeitado: `nan`/`inf`/`infinity`, hexadecimal e sublinhado; erro de entrada, código 2, mensagem distinta de "fora da faixa".
- Zero negativo normalizado: `-0` é impresso como `0`.
- Classes de falha por código: 2 = entrada inválida (inclui operando especial e operando fora da faixa); 1 = erro de cálculo (divisão por zero e resultado fora da faixa).

### Key Business Rules
- Só finitos: nenhum `NaN` ou `Inf` entra como operando nem sai como resultado ou texto em stdout.
- A validação do operando vem antes do cálculo: `divisao nan 0` passa a ser entrada inválida (código 2) e não "divisão por zero" (código 1).
- Na divisão, o divisor zero é checado primeiro; só depois a finitude do quociente (`divisao 1e308 1e-308` é "fora da faixa", `divisao 0 0` é divisão por zero).
- Em erro, stdout fica vazio; em sucesso, stderr fica vazio.
- Formato de impressão inalterado, com a única exceção de `-0` → `0`. Consequência registrada: o formato `'g'` com `-1` usa notação científica quando o expoente decimal é menor que -4 ou maior ou igual a 6 (`500000 + 500000` imprime `1e+06`).
- Separador decimal só o ponto.
- Underflow não é tratado como erro nesta change (em aberto).
- Nenhum `panic`/`recover`; D-A6, D-A8, D-A12, D-A13 e o `harness.config.json` não são tocados.

## Strategic Approach

### Solution Direction
- Rejeitar os operandos especiais em `lerOperando`: depois do `ParseFloat`, descartar valor `NaN`/infinito, textos com `_` e textos hexadecimais (`0x`/`0X` após o sinal), mantendo `ErrRange` → "fora da faixa" e demais erros → "não é um número válido", sempre com código 2. A mensagem exata dos casos novos é do planejador (não é contrato).
- Em `internal/calculo`, criar `ErrResultadoForaDaFaixa` e fazer as quatro operações devolverem `(float64, error)`, com erro quando o resultado não é finito; na divisão, divisor zero antes da checagem de finitude.
- Em `internal/cli`, propagar o erro por `calcular`, mapear os dois sentinelas para o código 1 em `reportar` e normalizar `-0` em `formatar`, mantendo `'g'` com `-1`.
- Testes de tabela sobre `Executar` comparando o texto exato de stdout e de código; uma tabela de entradas especiais confirmando ausência de pânico.
- README e CHANGELOG atualizados; nenhuma dependência nova.

### Key Design Decisions
- Regra de faixa no pacote de cálculo, não na CLI: o contrato de `operacoes-aritmeticas` fica verificável sem E/S; custo: mudança de assinatura das operações e dos testes existentes → recomendado (D-B1).
- Rejeição de especiais na entrada, antes do cálculo: evita que `divisao nan 0` pareça divisão por zero → decisão do humano (D-A10) com mecanismo proposto (D-B3).
- Código 1 passa a significar "erro de cálculo" em geral → decisão do humano; a mensagem distingue os casos (D-B4).
- Normalizar `-0` em `formatar`: um ponto só, cobre qualquer caminho que produza `-0` (D-B5).
- Sem nova capability (D-B6): as regras pertencem a `cli-calculadora` e `operacoes-aritmeticas`.

### Alternatives Considered
- Verificar a finitude do resultado só na CLI: deixaria o pacote de cálculo devolvendo `Inf` como resultado válido; descartada.
- Gramática decimal própria no lugar do `ParseFloat` mais rejeições: aceitável se mantiver `+5`, `.5`, `5.` e `1e3` aceitos; fica a critério da implementação.
- Aceitar vírgula, milhar, hexadecimal: fora do escopo (D-A11 só ponto).
- Mudar o formato para evitar notação científica: fora do escopo (D-A9 mantém o atual).
- Tratar underflow como erro: pergunta em aberto, não decidida.

## Risk & Gap Analysis

### Requirement Ambiguities
- Underflow (em aberto): `soma 1e-999 0` e `multiplicacao 1e-200 1e-200` imprimem `0` com código 0; falta saber se isso permanece, se `1e-999` vira "fora da faixa" (código 2) e/ou se underflow de resultado vira erro (código 1). Mantido como está; sem cenário nem teste.
- Texto exato das mensagens: não informado; as specs exigem português, stderr e os trechos "fora da faixa" e "divisão por zero".
- Formatos que o `ParseFloat` aceita e que a spec não cita (`+5`, `.5`, `5.`, `0_1` com sublinhado): `+5`, `.5` e `5.` continuam aceitos sem cenário próprio; `0_1` cai na regra do sublinhado.
- D-A12 e D-A13: não informados e fora do escopo.

### Edge Cases
- `nan` com sinal já era rejeitado pelo `ParseFloat` (código 2 por sintaxe inválida); sem sinal e `NaN` eram aceitos.
- `-0` como divisor e `0 0` já são divisão por zero; continuam assim.
- `divisao nan 0` muda de código 1 para 2 (rejeição de entrada antes do cálculo).
- `-0` como operando continua aceito; só a saída é normalizada (`soma -0 -0` → `0`).
- Hexadecimal sem expoente `p` (`0x10`) já é rejeitado pelo `ParseFloat`; com `p` (`0x1p-2`) era aceito. Hexadecimal com expoente enorme pode dar `ErrRange`; ambos terminam em código 2.
- Resultado `-0` por underflow (`multiplicacao -1e-200 1e-200`) também sai como `0`; não decide o underflow.
- Subtração com sinais opostos que estoura (`subtracao -1e308 1e308`) vira erro de faixa, código 1.
- Resultados como `Inf - Inf` ou `0 * Inf` não ocorrem pela CLI (operandos finitos), mas a função de cálculo, chamada direto com `NaN`/`Inf`, também devolve erro de faixa para resultado não finito.

### Technical Risks
- Alterar a assinatura de `Somar`/`Subtrair`/`Multiplicar` quebra `calculo_test.go` e o `switch` de `calcular`: mitigação, tarefas 2.x e 4.1 atualizam código e testes juntos.
- Mudança de comportamento observável (entradas especiais: 0 → 2; overflow: 0 → 1; `divisao nan 0`: 1 → 2; `-0` → `0`): intencional (D-A10, D-A9b); sinalizada na proposta e no CHANGELOG.
- Notação científica inesperada para inteiros de 7 dígitos (`1e+06`) e números pequenos (`1e-05`): consequência de manter D-A9; documentar no README.
- Dependência do que o `ParseFloat` aceita: mitigação com testes que fixam as rejeições e a aceitação de decimais comuns.
- Gate de fase: `cmd/` e `internal/` só recebem escrita em `apply`; a transição é da sessão principal depois da aprovação humana.

### Acceptance Criteria Coverage
| AC# | Description | Addressable? | Gaps/Notes |
|-----|-------------|--------------|------------|
| AC1 | `nan`/`inf`/`infinity` rejeitados | Yes | Cenários em `cli-calculadora`; tarefa 3.1 |
| AC2 | Hexadecimal e sublinhado rejeitados | Yes | Tarefa 3.2 |
| AC3 | Vírgula decimal é erro | Yes | Tarefa 3.3 |
| AC4 | `1e999` fora da faixa, código 2 | Yes | Sem mudança; tarefa 3.3 |
| AC5 | Overflow vira erro, código 1 | Yes | Tarefas 2.1, 2.2, 4.1 |
| AC6 | `NaN`/`Inf` nunca em stdout | Yes | Tarefa 4.1 confere stdout |
| AC7 | Divisão por zero código 1 | Yes | Já implementado; tarefas 2.2 e 4.2 reafirmam |
| AC8 | Formato atual mantido | Yes | Tarefa 4.3; consequência `1e+06` documentada |
| AC9 | `-0` impresso como `0` | Yes | Tarefa 4.3 |
| AC10 | Sem pânico | Yes | Tabela de ausência de pânico (3.3) |
| AC11 | Underflow mantido e em aberto | Yes | Pergunta no `design.md` |
| AC12 | `harness.config.json` intocado | Yes | Respeitado |
