## MODIFIED Requirements

### Requirement: Divisão por zero informada ao usuário
O programa MUST escrever em stderr uma mensagem de erro em português quando a divisão tiver divisor zero (positivo ou negativo, inclusive com dividendo zero), MUST NOT escrever resultado em stdout e MUST encerrar com código de saída 1, sem pânico.

#### Scenario: Divisão por zero pela linha de comando
- **WHEN** o usuário invoca o programa com os argumentos `divisao 6 0`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 1
- **AND** nenhum pânico ocorre

#### Scenario: Divisão por zero negativo
- **WHEN** o usuário invoca o programa com os argumentos `divisao 6 -0`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 1
- **AND** nenhum pânico ocorre

#### Scenario: Zero dividido por zero
- **WHEN** o usuário invoca o programa com os argumentos `divisao 0 0`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 1
- **AND** nenhum pânico ocorre

### Requirement: Entrada inválida informada ao usuário
O programa MUST escrever em stderr uma mensagem de erro em português e MUST encerrar com código de saída 2, sem pânico, quando a operação for desconhecida, o número de argumentos for diferente de três ou um operando não for um número aceito pelos requisitos de operandos numéricos (aceitos, especiais rejeitados e fora da faixa). Em todos esses casos MUST NOT escrever nada em stdout.

#### Scenario: Operação desconhecida
- **WHEN** o usuário invoca o programa com os argumentos `potencia 2 3`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

#### Scenario: Operando que não é número
- **WHEN** o usuário invoca o programa com os argumentos `soma abc 3`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

#### Scenario: Operandos ausentes
- **WHEN** o usuário invoca o programa com os argumentos `soma 2`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

#### Scenario: Argumentos em excesso
- **WHEN** o usuário invoca o programa com os argumentos `soma 2 3 4`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

## ADDED Requirements

### Requirement: Operandos numéricos aceitos
O programa SHALL aceitar como operando somente um número decimal finito de `float64`, escrito com sinal opcional, dígitos, ponto como único separador decimal e, opcionalmente, expoente decimal com `e` ou `E` (por exemplo `3`, `-2.5`, `1e3`, `-0`). A validação MUST ocorrer antes de qualquer cálculo, de modo que um operando rejeitado nunca é confundido com divisão por zero, e nenhum operando pode causar pânico.

#### Scenario: Operandos decimais finitos
- **WHEN** o usuário invoca o programa com os argumentos `soma 1e3 -2.5`
- **THEN** o programa escreve em stdout somente o resultado numérico, de valor 997.5
- **AND** o código de saída do processo é 0
- **AND** nada é escrito em stderr

### Requirement: Operandos especiais rejeitados
O programa MUST rejeitar como operando, com mensagem em português em stderr, nada em stdout e código de saída 2, os textos `nan`, `inf` e `infinity` (sem diferenciar maiúsculas de minúsculas, com ou sem sinal), os números hexadecimais (como `0x1p-2`), os números com sublinhados (como `1_000`) e os números com vírgula como separador decimal (como `1,5`). Nenhum desses casos pode causar pânico.

#### Scenario: NaN como operando
- **WHEN** o usuário invoca o programa com os argumentos `divisao 1 nan`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

#### Scenario: NaN com sinal ou em outra caixa
- **WHEN** o usuário invoca o programa com cada um dos argumentos `soma NaN 1`, `soma +nan 1` e `soma -NAN 1`
- **THEN** em cada invocação o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

#### Scenario: NaN como dividendo com divisor zero
- **WHEN** o usuário invoca o programa com os argumentos `divisao nan 0`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

#### Scenario: Infinitos como operandos
- **WHEN** o usuário invoca o programa com os argumentos `soma inf -inf`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

#### Scenario: Infinity por extenso
- **WHEN** o usuário invoca o programa com os argumentos `soma Infinity 1`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

#### Scenario: Infinito com sinal em outra caixa
- **WHEN** o usuário invoca o programa com cada um dos argumentos `soma -INF 1` e `soma +Infinity 1`
- **THEN** em cada invocação o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

#### Scenario: Hexadecimal como operando
- **WHEN** o usuário invoca o programa com os argumentos `soma 0x1p-2 1`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

#### Scenario: Sublinhado como operando
- **WHEN** o usuário invoca o programa com os argumentos `soma 1_000 1`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

#### Scenario: Vírgula como separador decimal
- **WHEN** o usuário invoca o programa com os argumentos `soma 1,5 2`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

### Requirement: Operando fora da faixa rejeitado
O programa MUST rejeitar com código de saída 2, nada em stdout e mensagem em português em stderr que indique "fora da faixa" um operando decimal cujo valor esteja fora da faixa de `float64` (como `1e999`), sem pânico.

#### Scenario: Operando decimal fora da faixa de float64
- **WHEN** o usuário invoca o programa com os argumentos `soma 1e999 1`
- **THEN** o programa escreve em stderr uma mensagem em português que indica que o operando está fora da faixa
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 2

### Requirement: Resultado fora da faixa informado ao usuário
O programa MUST escrever em stderr uma mensagem de erro em português que indique que o resultado está fora da faixa quando uma operação com operandos aceitos produzir um resultado não finito (overflow para infinito), MUST NOT escrever resultado em stdout e MUST encerrar com código de saída 1 (erro de cálculo, o mesmo da divisão por zero), sem pânico. O programa MUST NOT imprimir `Inf`, `+Inf`, `-Inf` nem `NaN` em stdout em nenhuma situação.

#### Scenario: Soma que estoura a faixa
- **WHEN** o usuário invoca o programa com os argumentos `soma 1e308 1e308`
- **THEN** o programa escreve em stderr uma mensagem em português que indica resultado fora da faixa
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 1
- **AND** nenhum pânico ocorre

#### Scenario: Multiplicação que estoura a faixa
- **WHEN** o usuário invoca o programa com os argumentos `multiplicacao 1e200 1e200`
- **THEN** o programa escreve em stderr uma mensagem em português que indica resultado fora da faixa
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 1

#### Scenario: Divisão que estoura a faixa
- **WHEN** o usuário invoca o programa com os argumentos `divisao 1e308 1e-308`
- **THEN** o programa escreve em stderr uma mensagem em português que indica resultado fora da faixa
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 1

### Requirement: Formato de impressão do resultado
O programa SHALL imprimir o resultado em stdout como o menor texto decimal que reconstrói exatamente o valor `float64` (formato geral com precisão mínima, o equivalente a `strconv.FormatFloat` com `'g'` e precisão `-1`), com ponto como único separador decimal, sem sufixo de tipo, em uma única linha terminada por uma quebra de linha. Um resultado inteiro MUST ser impresso sem parte fracionária (`5`, não `5.0`). O programa MUST imprimir o zero negativo como `0`, nunca como `-0`.

#### Scenario: Resultado inteiro sem parte fracionária
- **WHEN** o usuário invoca o programa com os argumentos `soma 2 3`
- **THEN** o programa escreve em stdout exatamente `5` seguido de uma quebra de linha
- **AND** o código de saída do processo é 0

#### Scenario: Menor texto que reconstrói o valor
- **WHEN** o usuário invoca o programa com os argumentos `soma 0.1 0.2`
- **THEN** o programa escreve em stdout exatamente `0.30000000000000004` seguido de uma quebra de linha
- **AND** o código de saída do processo é 0

#### Scenario: Notação científica para magnitudes grandes
- **WHEN** o usuário invoca o programa com os argumentos `soma 1e22 0`
- **THEN** o programa escreve em stdout exatamente `1e+22` seguido de uma quebra de linha
- **AND** o código de saída do processo é 0

#### Scenario: Zero negativo produzido por multiplicação
- **WHEN** o usuário invoca o programa com os argumentos `multiplicacao -1 0`
- **THEN** o programa escreve em stdout exatamente `0` seguido de uma quebra de linha
- **AND** o código de saída do processo é 0

#### Scenario: Zero negativo digitado como operando
- **WHEN** o usuário invoca o programa com os argumentos `soma -0 -0`
- **THEN** o programa escreve em stdout exatamente `0` seguido de uma quebra de linha
- **AND** o código de saída do processo é 0
