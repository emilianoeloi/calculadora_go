# Design

## Context

A change `calculadora-basica` (arquivada) entregou `calc <operacao> <a> <b>` sobre `float64`, com código em `cmd/calc/`, `internal/calculo/` e `internal/cli/`, e deixou em aberto D-A9, D-A10 e D-A11. Esta change fecha essas três decisões com a escolha do humano e alinha as specs `operacoes-aritmeticas` e `cli-calculadora` (deltas em `specs/`). As decisões D-A6, D-A8, D-A12 e D-A13 seguem como estavam e estão fora do escopo.

Comportamento atual, reconfirmado em 2026-11-12 lendo `internal/cli/cli.go` e `internal/calculo/calculo.go` e executando o binário:

- `lerOperando` chama `strconv.ParseFloat(texto, 64)`. `ErrRange` vira "operando %q fora da faixa de float64"; qualquer outro erro vira "operando %q não é um número válido"; ambos saem com código 2.
- `ParseFloat` aceita `nan`, `NaN`, `inf`, `-Inf`, `+INFINITY` (maiúsculas indiferentes), hexadecimais com expoente `p` (`0x1p-2` vale 0.25) e sublinhados (`1_000` vale 1000). Resultado observado: `soma nan 1` imprime `NaN`, `soma -Inf 1` imprime `-Inf`, `soma inf -inf` imprime `NaN`, `soma 0x1p-2 1` imprime `1.25`, `soma 1_000 1` imprime `1001`, todos com código 0.
- `nan` com sinal já é rejeitado pelo `ParseFloat` como sintaxe inválida: `soma +nan 1` e `soma -nan 1` saem com código 2. `inf` e `infinity` com sinal são aceitos.
- `divisao nan 0` sai com "divisão por zero" e código 1 (o divisor é 0 e `Dividir` não olha o dividendo). `divisao 1 nan` imprime `NaN` com código 0.
- `1e-999` é aceito como `0` sem erro (`soma 1e-999 0` imprime `0`). `1e999` e `-1e999` dão `ErrRange`, código 2.
- `soma 1e308 1e308`, `multiplicacao 1e200 1e200` e `divisao 1e308 1e-308` imprimem `+Inf` com código 0.
- `formatar` usa `strconv.FormatFloat(valor, 'g', -1, 64)`: `multiplicacao -1 0` e `soma -0 -0` imprimem `-0`; `soma 0.1 0.2` imprime `0.30000000000000004`; `soma 1e22 0` imprime `1e+22`.
- `calculo.Dividir` já devolve `ErrDivisaoPorZero` para `b == 0` (inclusive `-0` e `0/0`); `internal/cli/cli_test.go` já cobre `divisao 6 0`, `6 -0` e `0 0`. `Somar`, `Subtrair` e `Multiplicar` devolvem só `float64`.
- `reportar` é o único ponto que mapeia erro em código: `ErrDivisaoPorZero` vira 1, qualquer outro erro vira 2.

## Goals / Non-Goals

**Goals:**
- Garantir que `Inf` e `NaN` nunca entram como operando nem saem como resultado impresso.
- Distinguir três classes de falha por código: entrada inválida (2), erro de cálculo (1: divisão por zero e resultado fora da faixa) e sucesso (0).
- Fixar o formato de impressão atual e o zero negativo normalizado, com testes que comparam o texto exato.
- Cobrir cada cenário das specs com teste automatizado executável por `go test ./...`, sem nenhum caminho que cause pânico.

**Non-Goals:**
- Aceitar vírgula, separador de milhar, hexadecimal ou qualquer outra notação além do decimal com ponto.
- Alterar o formato de impressão (mudar casas decimais, evitar notação científica).
- Decidir o tratamento de underflow (ver Open Questions).
- Alterar o `harness.config.json`, adicionar dependências externas ou novas operações.

## Decisions

### Decididas pelo humano (registradas, não reabrir)

- **D-A10 (decidida). Política "só finitos".** A entrada rejeita `nan`, `inf` e `infinity` (sem diferenciar maiúsculas, com ou sem sinal), hexadecimais (como `0x1p-2`) e sublinhados (como `1_000`), com mensagem em português em stderr, stdout vazio e código 2 (o mesmo de "operando não é um número válido"). Um resultado infinito por overflow com operandos finitos (`soma 1e308 1e308`, `multiplicacao 1e200 1e200`, `divisao 1e308 1e-308`) vira erro "resultado fora da faixa", stdout vazio, código 1 (erro de cálculo, como a divisão por zero). `NaN` como resultado deixa de ser possível. A rejeição de `1e999` como fora da faixa (código 2) continua.
- **D-A9 (decidida). Formato de impressão: manter o atual.** `strconv.FormatFloat` com `'g'` e precisão `-1` (menor texto que reconstrói o valor); inteiro imprime `5`; uma linha com quebra de linha final.
- **D-A9b (decidida). Zero negativo normalizado.** `-0` é impresso como `0` (`multiplicacao -1 0` e `soma -0 -0` imprimem `0`).
- **D-A11 (decidida). Separador decimal só o ponto,** na entrada e na saída. `1,5` continua erro com código 2.
- **Não alteradas:** D-A6, D-A8, D-A12 e D-A13 da change arquivada.

### Decisões do planejador (propostas; o humano pode contestar)

- **D-B1. Erro sentinela de faixa em `internal/calculo`.** Um novo erro exportado (nome sugerido `ErrResultadoForaDaFaixa`, mensagem "resultado fora da faixa") comparável com `errors.Is`, distinto de `ErrDivisaoPorZero`. As quatro operações passam a devolver `(float64, error)` e devolvem o erro quando o resultado não é finito (`math.IsInf` ou `math.IsNaN`). Motivo: a regra é de negócio e fica testável sem a CLI, e o requisito de `operacoes-aritmeticas` fala em "devolver erro". Alternativa descartada: verificar a finitude só em `internal/cli`, que deixaria o pacote de cálculo devolvendo `Inf` como se fosse resultado válido. Consequência: assinatura de `Somar`, `Subtrair` e `Multiplicar` muda, e `internal/calculo/calculo_test.go` e o `switch` de `calcular` em `internal/cli/cli.go` acompanham.
- **D-B2. Precedência dos erros na divisão.** Primeiro o divisor zero (`ErrDivisaoPorZero`), depois a finitude do quociente (`ErrResultadoForaDaFaixa`). Assim `divisao 0 0` e `divisao 6 -0` seguem como divisão por zero (código 1) e `divisao 1e308 1e-308` é "fora da faixa" (código 1).
- **D-B3. Onde rejeitar os especiais na entrada.** Em `lerOperando`, antes de qualquer cálculo, mantendo `reportar` como único ponto que converte erro em código. Ordem de decisão: (1) `ErrRange` do `ParseFloat` mantém a mensagem "fora da faixa" e o código 2; (2) qualquer outro erro do `ParseFloat` mantém "não é um número válido" e o código 2; (3) se o texto foi aceito, rejeitar com código 2 quando o valor for `NaN` ou infinito, quando o texto contiver `_` ou quando, retirado o sinal, começar com `0x`/`0X` (hexadecimal). Os valores `NaN` e infinito só aparecem em texto aceito pelo `ParseFloat` como `nan`, `inf` ou `infinity`, então a checagem do valor cobre todas as caixas e sinais. A mensagem para esses casos é do planejador (por exemplo "operando %q não é um número finito válido"); o texto exato não é contrato, só "em português, em stderr, código 2".
  Alternativa considerada: validar o texto com uma gramática decimal própria antes do `ParseFloat`. Fica a critério da implementação desde que o resultado seja o mesmo: `+5`, `.5`, `5.` e `1e3` (aceitos hoje) continuam aceitos, e a lista de rejeições da spec é respeitada. Esses formatos não são cenários da spec e nenhum deles é alterado.
- **D-B4. Mapeamento de código de saída.** `reportar` passa a mapear `ErrDivisaoPorZero` e `ErrResultadoForaDaFaixa` para o código 1 e o resto para 2. A constante `CodigoDivisaoPorZero` pode ser renomeada (por exemplo `CodigoErroDeCalculo`) pela implementação; o valor 1 não muda.
- **D-B5. Zero negativo.** `formatar` normaliza o zero negativo (`valor == 0` produz `0`) antes de `FormatFloat`. Isso cobre também qualquer outro caminho que gere `-0` (por exemplo `multiplicacao -1 0`, `soma -0 -0` e `subtracao -0 0`).
- **D-B6. Sem nova capability.** As regras novas descrevem a entrada, a saída e os erros da CLI (`cli-calculadora`) e o contrato de resultado do cálculo (`operacoes-aritmeticas`). Criar uma capability só para "números especiais" ou "formato" fragmentaria o que já tem dono.
- **D-B7. Sem pânico.** Nenhum `panic`/`recover`; toda rejeição é valor de erro. A cobertura é uma tabela de testes que executa todas as entradas especiais da spec e verifica que `Executar` retorna (um pânico derrubaria o teste).

## Risks / Trade-offs

- [O formato `'g'` com precisão `-1` passa para notação científica quando o expoente decimal é menor que -4 ou maior ou igual a 6, e não só para números enormes: verificado em 2026-11-12 que `soma 500000 500000` imprime `1e+06`, `soma 123456789 0` imprime `1.23456789e+08` e `soma 0.00001 0` imprime `1e-05`, enquanto `soma 100000 1` imprime `100001`] → É consequência direta de D-A9 (manter o formato atual). Registrado para que o humano saiba que inteiros grandes e números pequenos saem em notação científica; não é cenário de spec além de `1e+22`. Mudar isso é nova decisão.
- [Quem dependia de `inf`, `nan`, hexadecimal ou `1_000` como entrada, ou de `+Inf`/`NaN` na saída com código 0, passa a receber código 1 ou 2] → Mudança intencional (D-A10); sinalizada em `proposal.md` e no CHANGELOG. O projeto não tem outros consumidores conhecidos (não informado).
- [Alterar a assinatura de `Somar`, `Subtrair` e `Multiplicar` quebra os testes existentes de `calculo`] → Esperado; as tarefas 2.x atualizam testes junto com o código.
- [O código 1 passa a significar "erro de cálculo" (divisão por zero ou resultado fora da faixa); scripts que liam 1 como "só divisão por zero" perdem precisão] → Decisão do humano (D-A10: "código 1, erro de cálculo, como a divisão por zero"); a mensagem em stderr distingue os casos.
- [`ParseFloat` pode mudar o que aceita em versões futuras do Go] → A spec lista o que deve ser rejeitado e os testes fixam esses casos; um teste de regressão de aceitação de decimais comuns protege o outro lado.
- [Resultado `-0` produzido por underflow (por exemplo `multiplicacao -1e-200 1e-200`) também sai como `0` por D-B5] → Coerente com D-A9b; não decide o underflow em si (abaixo).

## Open Questions

Nenhuma suposição foi feita sobre estes itens; a implementação NÃO deve decidi-los.

- **Underflow (em aberto, para o humano).** Hoje, operando que sofre underflow é aceito como `0` sem erro (`soma 1e-999 0` imprime `0`; `1e-324` também vira `0`), e um resultado que sofre underflow também sai como `0` (por exemplo `multiplicacao 1e-200 1e-200`), sempre com código 0. Pergunta: isso deve continuar assim, ou o underflow de operando (como `1e-999`) deve ser rejeitado como "fora da faixa" (código 2, em paralelo a `1e999`) e/ou o underflow de resultado deve virar erro "resultado fora da faixa" (código 1)? Esta change mantém o comportamento atual e não cria cenário de spec nem teste para underflow.
- D-A6, D-A8, D-A12 e D-A13: seguem como na change arquivada (não alteradas aqui). D-A12 (versão mínima do Go; `go.mod` diz `go 1.27.1`) e D-A13 (idioma dos identificadores; hoje português) continuam "não informados".
- Texto exato das mensagens de erro: não informado; as specs exigem apenas português, stderr e os trechos "fora da faixa" e "divisão por zero" onde indicado.
