# Tasks

## 1. Pré-requisitos

- [ ] 1.1 Conferir no `design.md` o registro das decisões D-A10, D-A9, D-A9b e D-A11 como decididas, D-B1 a D-B7 como propostas do planejador e o underflow em Open Questions; se o humano resolver o underflow ou contestar alguma decisão D-B antes da implementação, atualizar `design.md` e as specs correspondentes e só então seguir
- [ ] 1.2 Rodar `go test ./...` e `go vet ./...` antes de alterar código e verificar que ambos terminam com código 0 (linha de base)

## 2. Cálculo (spec `operacoes-aritmeticas`)

- [ ] 2.1 Adicionar em `internal/calculo` o erro sentinela de resultado fora da faixa (D-B1), fazer `Somar`, `Subtrair` e `Multiplicar` devolverem `(float64, error)` com esse erro quando o resultado não for finito, e verificar com testes de unidade que cobrem "Soma que estoura a faixa" (1e308 e 1e308), "Multiplicação que estoura a faixa" (1e200 e 1e200), "Subtração que estoura a faixa para o lado negativo" (-1e308 e 1e308) e "Resultado grande mas finito" (1e308 e 1e307), usando `errors.Is` e conferindo que os cenários de soma, subtração e multiplicação existentes (2 e 3, 5 e 3, 4 e 3) seguem sem erro (`go test ./...` passa)
- [ ] 2.2 Fazer `Dividir` checar primeiro o divisor zero e depois a finitude do quociente (D-B2) e verificar com testes de unidade que cobrem "Divisão que estoura a faixa" (1e308 e 1e-308: `ErrResultadoForaDaFaixa` e não `ErrDivisaoPorZero`), "Divisor igual a zero negativo" (6 e -0), "Zero dividido por zero" (0 e 0) e os cenários existentes de divisão, sem pânico (`go test ./...` passa)

## 3. Entrada da linha de comando (spec `cli-calculadora`)

- [ ] 3.1 Em `lerOperando`, rejeitar com mensagem em português, nada em stdout e código 2 os operandos `nan`, `inf` e `infinity` (qualquer caixa, com ou sem sinal), mantendo a precedência de D-B3, e verificar com testes de tabela sobre `Executar` que cobrem "NaN como operando" (`divisao 1 nan`), "NaN com sinal ou em outra caixa" (`NaN`, `+nan`, `-NAN`), "NaN como dividendo com divisor zero" (`divisao nan 0`: código 2, e não 1), "Infinitos como operandos" (`soma inf -inf`), "Infinity por extenso" (`soma Infinity 1`) e "Infinito com sinal em outra caixa" (`-INF`, `+Infinity`) (`go test ./...` passa)
- [ ] 3.2 Rejeitar hexadecimais e sublinhados com a mesma mensagem em português, stdout vazio e código 2, e verificar com testes que cobrem "Hexadecimal como operando" (`soma 0x1p-2 1`) e "Sublinhado como operando" (`soma 1_000 1`), incluindo uma variante com sinal (`-0x1p-2`) e com o operando na segunda posição (`soma 1 1_000`) (`go test ./...` passa)
- [ ] 3.3 Verificar com testes que "Vírgula como separador decimal" (`soma 1,5 2`) e "Operando decimal fora da faixa de float64" (`soma 1e999 1` e `soma -1e999 1`: código 2, mensagem com "fora da faixa", stdout vazio) seguem funcionando, e que "Operandos decimais finitos" (`soma 1e3 -2.5` imprime `997.5`, além de `+5`, `.5` e `5.` aceitos como hoje) não regride; adicionar uma tabela única com todas as entradas especiais desta seção que confirma que `Executar` retorna sem pânico (`go test ./...` passa)

## 4. Saída da linha de comando (spec `cli-calculadora`)

- [ ] 4.1 Fazer `calcular` propagar o erro de faixa e `reportar` mapear `ErrResultadoForaDaFaixa` para o código 1 com mensagem em português que contenha "fora da faixa" (D-B4), e verificar com testes sobre `Executar` que cobrem "Soma que estoura a faixa" (`soma 1e308 1e308`), "Multiplicação que estoura a faixa" (`multiplicacao 1e200 1e200`) e "Divisão que estoura a faixa" (`divisao 1e308 1e-308`): stdout vazio, stderr com a mensagem, código 1; conferir também que a divisão por zero segue com código 1 e que nenhuma saída contém `Inf` nem `NaN` (`go test ./...` passa)
- [ ] 4.2 Reafirmar a divisão por zero na CLI com testes que cobrem "Divisão por zero pela linha de comando" (`divisao 6 0`), "Divisão por zero negativo" (`divisao 6 -0`) e "Zero dividido por zero" (`divisao 0 0`): erro em português em stderr, stdout vazio, código 1 (os casos `-0` e `0 0` já existem em `cli_test.go`; manter) (`go test ./...` passa)
- [ ] 4.3 Em `formatar`, normalizar o zero negativo para `0` (D-B5) mantendo `'g'` com precisão `-1` (D-A9), e verificar com testes que comparam o texto exato de stdout, incluindo a quebra de linha final, em "Resultado inteiro sem parte fracionária" (`soma 2 3` → `5`), "Menor texto que reconstrói o valor" (`soma 0.1 0.2` → `0.30000000000000004`), "Notação científica para magnitudes grandes" (`soma 1e22 0` → `1e+22`), "Zero negativo produzido por multiplicação" (`multiplicacao -1 0` → `0`) e "Zero negativo digitado como operando" (`soma -0 -0` → `0`) (`go test ./...` passa)

## 5. Integração e documentação

- [ ] 5.1 Rodar `go vet ./...`, `gofmt -l .` e `go test ./...` e verificar que o primeiro e o terceiro terminam com código 0 e que o `gofmt` não lista arquivos
- [ ] 5.2 Atualizar o `README.md` (tabela de códigos de saída: 1 passa a incluir resultado fora da faixa; remover da linha "Não informado" o formato, `Inf`/`NaN` e o separador decimal; registrar que só decimais finitos com ponto são aceitos e o formato de impressão) e verificar executando `go run ./cmd/calc soma 1e308 1e308` (código 1), `go run ./cmd/calc soma nan 1` (código 2), `go run ./cmd/calc multiplicacao -1 0` (imprime `0`) e `go run ./cmd/calc soma 0.1 0.2` (imprime `0.30000000000000004`)
- [ ] 5.3 Registrar a change no `CHANGELOG.md`, incluindo as mudanças de comportamento (entradas especiais passam a código 2, overflow passa a código 1, `-0` impresso como `0`) e o underflow como pergunta em aberto

## Workflow follow-up

- O `harness.config.json` NÃO é alterado por esta change.
- Arquivar a change depois que a revisão for aprovada (`openspec archive numeros-especiais-e-formato`).
