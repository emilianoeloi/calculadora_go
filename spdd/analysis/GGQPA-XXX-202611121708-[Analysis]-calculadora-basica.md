# SPDD Analysis: Calculadora básica de linha de comando em Go

## Original Business Requirement

Pedido original (verbatim), recebido pela sessão principal:

> Planeje a change `calculadora-basica` no projeto /Users/emiliano.barbosa/Pessoal/calculadora_go (fase `proposal` já declarada). Pedido: calculadora de linha de comando em Go com soma, subtração, multiplicação e divisão; divisão por zero devolve erro, sem pânico. Código em `cmd/` e `internal/`; testes com `go test ./...` (será o sensors_command, mas NÃO altere o harness.config.json agora). Artefatos OpenSpec em português. Entregue a proposta OpenSpec validada (`openspec validate calculadora-basica --strict`), a análise SPDD em spdd/analysis/ e o REASONS Canvas em spdd/prompt/. Não escreva código nem faça commit. O que faltar informar vira "não informado", nunca suposição (ex.: sintaxe exata de argumentos, formato de saída, precisão numérica: se não definido, registre como decisão em aberto para o humano). Ao final, liste os arquivos criados e as decisões em aberto.

Decisões do humano recebidas depois do pedido original (segunda rodada, registradas no `design.md`):

> - D-A1 sintaxe: `calc <operacao> <a> <b>`, operação por nome em português como primeiro argumento: soma, subtracao, multiplicacao, divisao (sem acento). Uma operação por execução, exatamente dois operandos.
> - D-A3 numérico: float64 para operandos e resultado; divisão não trunca. Overflow/Inf/NaN: continua "não informado".
> - D-A2 saída: só o resultado numérico em stdout (exit 0); mensagens de erro em português em stderr. Formato exato de impressão do float (casas decimais): continua "não informado".
> - D-A4 códigos de saída: 0 sucesso; 1 divisão por zero; 2 uso/entrada inválida.
> - D-A5: módulo `github.com/emilianoeloi/calculadora_go`; binário em `cmd/calculadora/main.go`.
> - D-A14 (terceira rodada): o executável se chama `calc`; o ponto de entrada passa a ser `cmd/calc/main.go` (substitui o binário de D-A5).
> - Continuam em aberto: D-A6, D-A8 (cenários de entrada inválida, mantidos como propostos), separador decimal, versão mínima do Go, idioma dos identificadores, formato de impressão do float, overflow/Inf/NaN.

Artefatos OpenSpec relacionados (leitura obrigatória para o canvas): `openspec/changes/calculadora-basica/{proposal,design,tasks}.md` e `specs/{operacoes-aritmeticas,cli-calculadora}/spec.md`.

Critérios de aceite (derivados do pedido, numerados para a análise):
- AC1: a calculadora soma dois operandos.
- AC2: a calculadora subtrai dois operandos.
- AC3: a calculadora multiplica dois operandos.
- AC4: a calculadora divide dois operandos.
- AC5: divisão por zero devolve erro, sem pânico.
- AC6: código em `cmd/` e `internal/`.
- AC7: testes executam com `go test ./...`.
- AC8: o `harness.config.json` não é alterado nesta fase.

## Domain Concept Identification

Exploração do código: o projeto é greenfield. Não há `go.mod`, `cmd/` nem `internal/`; `openspec/specs/` e `spdd/{analysis,prompt}/` estavam vazios. Go 1.27.1 está instalado. O `CLAUDE.md` e o `harness.config.json` definem as pastas de código (`cmd/`, `internal/`), o gate de fase e que o `sensors_command` hoje é vazio. Não existem convenções de código, de erro ou de testes herdadas; as convenções são as idiomáticas de Go (não informado em contrário).

### Existing Concepts (from codebase)
- Fase da change (gate proposal/apply): controla quando `cmd/` e `internal/` podem receber escrita; durante o planejamento, nenhuma escrita de código é permitida.
- `sensors_command` do harness: hoje vazio; passará a ser `go test ./...`, por ação do humano e fora desta change.
- Specs OpenSpec da change: `operacoes-aritmeticas` e `cli-calculadora`, ambas novas.

### New Concepts Required
- Operação aritmética: uma das quatro (soma, subtração, multiplicação, divisão) aplicada a dois operandos; identificada na CLI pelos nomes `soma`, `subtracao`, `multiplicacao`, `divisao`. Relaciona-se ao Resultado e ao Erro de cálculo.
- Operando: número de ponto flutuante de 64 bits (`float64`), decisão D-A3.
- Resultado: valor `float64` produzido por uma operação bem-sucedida; a divisão não trunca.
- Erro de cálculo (divisão por zero): falha de negócio reportada como valor, não como pânico; termina com código de saída 1.
- Invocação de linha de comando: `<operacao> <a> <b>`, exatamente três argumentos, e a resposta apresentada (resultado em stdout ou mensagem em stderr).
- Erro de uso (entrada inválida): operação desconhecida, operando não numérico, argumentos ausentes ou em excesso; termina com código de saída 2. Conceito proposto pelo planejador (D-A8), mantido como proposto.
- Saída do programa: stdout só com o resultado numérico; stderr com mensagens de erro em português; código de saída 0, 1 ou 2.

### Key Business Rules
- Divisor zero na divisão é erro devolvido e nunca pânico (código de saída 1): governa Operação e Erro de cálculo. Como o tipo é `float64`, um divisor `-0` também é zero (consequência direta, sinalizada).
- Dividendo zero com divisor não zero resulta em zero, sem erro.
- Subtração e divisão não são comutativas: `<a>` é o minuendo/dividendo e `<b>` o subtraendo/divisor.
- Em sucesso, stdout recebe só o resultado e stderr fica vazio; em erro, stdout fica vazio.
- Nenhum caminho de falha pode encerrar o programa por pânico (regra implícita do pedido, reforçada pela regra do `CLAUDE.md` de não capturar exceções de forma ampla).
- Quem muda o `harness.config.json` é o humano, não esta change (AC8).

## Strategic Approach

### Solution Direction
- Duas camadas: lógica de cálculo isolada em `internal/` (sem E/S, testável) e ponto de entrada fino em `cmd/calc/main.go`, entre elas uma camada de interpretação de argumentos e apresentação testável sem executar processo externo.
- Fluxo: argumentos → interpretação (`<operacao> <a> <b>`) → cálculo em `float64` → resultado ou erro → stdout/stderr e código de saída 0/1/2.
- Erros como valor (`error`), com erro sentinela identificável para divisão por zero; sem `panic` e sem `recover`.
- Apenas a biblioteca padrão do Go; módulo `github.com/emilianoeloi/calculadora_go`.
- Testes de unidade por cenário das specs, executáveis por `go test ./...`.

### Key Design Decisions
- Separar cálculo de E/S: ganha testabilidade e isola a regra de divisão por zero; custo de um pouco mais de estrutura para um programa pequeno → recomendado.
- Erro sentinela em vez de `recover`: torna o caso identificável e evita captura ampla → recomendado.
- `float64` (decisão do humano): simples e cobre divisão sem truncar; custo é o erro de representação e a possibilidade de `Inf`/`NaN` → aceito; o tratamento de `Inf`/`NaN` fica em aberto (D-A10).
- Código de saída único por classe de falha (1 divisão por zero, 2 uso): permite a scripts distinguir os casos → decisão do humano.
- Formato de impressão do float, separador decimal, versão do Go, idioma dos identificadores: NÃO decididos; o planejamento não supõe valor.

### Alternatives Considered
- Tudo em um único `main.go`: rejeitado por misturar E/S e regra, dificultando testes.
- Capturar pânico com `recover`: rejeitado porque o pedido exige ausência de pânico e o `CLAUDE.md` restringe captura ampla.
- Modo interativo (REPL): fora do escopo; o pedido fala em calculadora de linha de comando com operações simples e não menciona interação contínua.
- Inteiros ou decimal exato: preteridos pelo humano em favor de `float64`.

## Risk & Gap Analysis

### Requirement Ambiguities
- Formato de impressão do float em stdout (D-A9): casas decimais, notação científica, `5` versus `5.0`; não informado.
- Overflow, `Inf` e `NaN` (D-A10): resultado `Inf`/`NaN` é erro ou valor impresso; operandos digitados como `NaN`/`Inf` são aceitos ou rejeitados; não informado.
- Separador decimal na entrada (D-A11): não informado.
- Versão mínima do Go (D-A12): não informado.
- Idioma dos identificadores no código (D-A13): não informado.
- Uma operação por execução e dois operandos (D-A6): tratada como escopo, sem confirmação explícita.
- Tratamento de entrada inválida (D-A8): fora do pedido literal; mantido como proposto pelo planejador, com código de saída 2.

### Edge Cases
- Divisor `-0`: igual a zero em `float64`, tratado como divisão por zero (consequência direta, sinalizada).
- `strconv.ParseFloat` aceita `NaN`, `Inf` e notação científica (por exemplo `1e400` gera erro de faixa): comportamento depende de D-A10.
- Operandos negativos (`-5`) como argumento: não são flags porque o programa lê posicionalmente, mas cabe teste.
- Overflow em soma e multiplicação produz `Inf` sem erro nativo: depende de D-A10.
- Erro de representação (0.1 + 0.2): aceito pela escolha de `float64`; os testes devem usar valores exatos ou tolerância definida (não informada).
- Argumentos com contagem diferente de três: coberto pelo código 2 (ausentes e excesso).

### Technical Risks
- Retrabalho nos testes de saída se D-A9 mudar depois de implementado: mitigação, resolver antes da tarefa 3.1 ou comparar o valor numérico lido.
- Sem `go.mod`, nada compila nem testa: mitigação, criar na tarefa 1.2 com o módulo `github.com/emilianoeloi/calculadora_go`.
- Gate de fase: escrever em `cmd/` e `internal/` só é possível em `apply`; a transição é da sessão principal depois da aprovação humana.
- `go test ./...` só vira sensor quando o humano atualizar o `harness.config.json`; até lá o sensor de Stop não valida o código.

### Acceptance Criteria Coverage
| AC# | Description | Addressable? | Gaps/Notes |
|-----|-------------|--------------|------------|
| AC1 | Soma | Yes | `float64`; sem lacunas |
| AC2 | Subtração | Yes | Ordem `<a>` menos `<b>` resolvida em D-A1 |
| AC3 | Multiplicação | Yes | Overflow para `Inf` depende de D-A10 |
| AC4 | Divisão | Yes | Não trunca (D-A3); impressão do float depende de D-A9 |
| AC5 | Divisão por zero devolve erro, sem pânico | Yes | Código 1 e mensagem em português em stderr resolvidos; texto exato da mensagem não informado |
| AC6 | Código em `cmd/` e `internal/` | Yes | Binário `cmd/calc/main.go` e módulo resolvidos (D-A5, D-A14) |
| AC7 | Testes via `go test ./...` | Yes | Nenhuma lacuna |
| AC8 | `harness.config.json` não alterado | Yes | Respeitado: o arquivo não é tocado nesta fase |
