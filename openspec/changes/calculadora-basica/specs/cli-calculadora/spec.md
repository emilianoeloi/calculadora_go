# Spec Delta

## Purpose

Define como a calculadora é usada pela linha de comando: recebe a operação e os dois operandos, executa o cálculo e entrega ao usuário o resultado (em stdout) ou o erro (em stderr), com códigos de saída distintos e sem pânico.

## ADDED Requirements

### Requirement: Sintaxe de invocação
O programa SHALL aceitar exatamente três argumentos na forma `calc <operacao> <a> <b>`, em que `<operacao>` é um dos nomes `soma`, `subtracao`, `multiplicacao` ou `divisao` (em português, sem acento) e `<a>` e `<b>` são os operandos, nessa ordem. Uma operação é executada por invocação.

#### Scenario: Operação válida
- **WHEN** o usuário invoca o programa com os argumentos `soma 2 3`
- **THEN** o programa escreve em stdout somente o resultado numérico, de valor 5
- **AND** o código de saída do processo é 0
- **AND** nada é escrito em stderr

#### Scenario: Ordem dos operandos
- **WHEN** o usuário invoca o programa com os argumentos `subtracao 5 3`
- **THEN** o programa escreve em stdout somente o resultado numérico, de valor 2
- **AND** o código de saída do processo é 0

#### Scenario: Divisão sem truncamento
- **WHEN** o usuário invoca o programa com os argumentos `divisao 7 2`
- **THEN** o programa escreve em stdout somente o resultado numérico, de valor 3.5
- **AND** o código de saída do processo é 0

### Requirement: Divisão por zero informada ao usuário
O programa MUST escrever em stderr uma mensagem de erro em português quando a divisão tiver divisor zero, MUST NOT escrever resultado em stdout e MUST encerrar com código de saída 1, sem pânico.

#### Scenario: Divisão por zero pela linha de comando
- **WHEN** o usuário invoca o programa com os argumentos `divisao 6 0`
- **THEN** o programa escreve uma mensagem de erro em português em stderr
- **AND** nada é escrito em stdout
- **AND** o código de saída do processo é 1
- **AND** nenhum pânico ocorre

### Requirement: Entrada inválida informada ao usuário
O programa MUST escrever em stderr uma mensagem de erro em português e MUST encerrar com código de saída 2, sem pânico, quando a operação for desconhecida, o número de argumentos for diferente de três ou um operando não for um número válido.

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
