# Spec Delta

## Purpose

Define o comportamento das quatro operações aritméticas básicas (soma, subtração, multiplicação e divisão) sobre números de ponto flutuante de 64 bits e o tratamento da divisão por zero como erro devolvido, nunca como pânico.

## ADDED Requirements

### Requirement: Soma de dois operandos
O sistema SHALL calcular a soma de dois operandos numéricos de ponto flutuante de 64 bits (float64) e devolver o resultado também em float64.

#### Scenario: Soma de dois números
- **WHEN** a operação de soma é pedida com os operandos 2 e 3
- **THEN** o resultado devolvido é 5
- **AND** nenhum erro é devolvido

### Requirement: Subtração de dois operandos
O sistema SHALL calcular a diferença entre o primeiro e o segundo operando e devolver o resultado em float64.

#### Scenario: Subtração de dois números
- **WHEN** a operação de subtração é pedida com os operandos 5 e 3
- **THEN** o resultado devolvido é 2
- **AND** nenhum erro é devolvido

### Requirement: Multiplicação de dois operandos
O sistema SHALL calcular o produto de dois operandos numéricos e devolver o resultado em float64.

#### Scenario: Multiplicação de dois números
- **WHEN** a operação de multiplicação é pedida com os operandos 4 e 3
- **THEN** o resultado devolvido é 12
- **AND** nenhum erro é devolvido

### Requirement: Divisão de dois operandos sem truncamento
O sistema SHALL calcular o quociente do primeiro operando pelo segundo quando o segundo operando for diferente de zero, sem truncar a parte fracionária.

#### Scenario: Divisão exata
- **WHEN** a operação de divisão é pedida com os operandos 6 e 3
- **THEN** o resultado devolvido é 2
- **AND** nenhum erro é devolvido

#### Scenario: Divisão com parte fracionária
- **WHEN** a operação de divisão é pedida com os operandos 7 e 2
- **THEN** o resultado devolvido é 3.5
- **AND** nenhum erro é devolvido

### Requirement: Divisão por zero devolve erro
O sistema MUST devolver um erro quando o divisor da divisão for zero, e MUST NOT entrar em pânico nem encerrar o programa de forma abrupta nesse caso.

#### Scenario: Divisor igual a zero
- **WHEN** a operação de divisão é pedida com os operandos 6 e 0
- **THEN** um erro é devolvido
- **AND** nenhum resultado numérico é apresentado como válido
- **AND** nenhum pânico ocorre

#### Scenario: Dividendo igual a zero com divisor diferente de zero
- **WHEN** a operação de divisão é pedida com os operandos 0 e 3
- **THEN** o resultado devolvido é 0
- **AND** nenhum erro é devolvido
