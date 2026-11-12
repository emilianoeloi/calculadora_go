## MODIFIED Requirements

### Requirement: Divisão por zero devolve erro
O sistema MUST devolver um erro quando o divisor da divisão for zero (positivo ou negativo, inclusive quando o dividendo também for zero), e MUST NOT entrar em pânico nem encerrar o programa de forma abrupta nesse caso.

#### Scenario: Divisor igual a zero
- **WHEN** a operação de divisão é pedida com os operandos 6 e 0
- **THEN** um erro é devolvido
- **AND** nenhum resultado numérico é apresentado como válido
- **AND** nenhum pânico ocorre

#### Scenario: Dividendo igual a zero com divisor diferente de zero
- **WHEN** a operação de divisão é pedida com os operandos 0 e 3
- **THEN** o resultado devolvido é 0
- **AND** nenhum erro é devolvido

#### Scenario: Divisor igual a zero negativo
- **WHEN** a operação de divisão é pedida com os operandos 6 e -0
- **THEN** o erro de divisão por zero é devolvido
- **AND** nenhum resultado numérico é apresentado como válido
- **AND** nenhum pânico ocorre

#### Scenario: Zero dividido por zero
- **WHEN** a operação de divisão é pedida com os operandos 0 e 0
- **THEN** o erro de divisão por zero é devolvido
- **AND** nenhum resultado numérico é apresentado como válido
- **AND** nenhum pânico ocorre

## ADDED Requirements

### Requirement: Resultado não finito devolve erro de faixa
O sistema MUST devolver um erro de resultado fora da faixa, distinto do erro de divisão por zero, quando o resultado de uma operação não for um número finito (infinito por overflow ou não numérico), e MUST NOT apresentar esse resultado como válido nem entrar em pânico. Com operandos finitos, o infinito só pode surgir por overflow; um resultado `NaN` nunca é devolvido como valor válido. A divisão por zero continua sendo reportada como divisão por zero, e não como resultado fora da faixa.

#### Scenario: Soma que estoura a faixa
- **WHEN** a operação de soma é pedida com os operandos 1e308 e 1e308
- **THEN** o erro de resultado fora da faixa é devolvido
- **AND** nenhum resultado numérico é apresentado como válido
- **AND** nenhum pânico ocorre

#### Scenario: Multiplicação que estoura a faixa
- **WHEN** a operação de multiplicação é pedida com os operandos 1e200 e 1e200
- **THEN** o erro de resultado fora da faixa é devolvido
- **AND** nenhum resultado numérico é apresentado como válido

#### Scenario: Divisão que estoura a faixa
- **WHEN** a operação de divisão é pedida com os operandos 1e308 e 1e-308
- **THEN** o erro de resultado fora da faixa é devolvido
- **AND** o erro de divisão por zero não é devolvido
- **AND** nenhum resultado numérico é apresentado como válido

#### Scenario: Subtração que estoura a faixa para o lado negativo
- **WHEN** a operação de subtração é pedida com os operandos -1e308 e 1e308
- **THEN** o erro de resultado fora da faixa é devolvido
- **AND** nenhum resultado numérico é apresentado como válido

#### Scenario: Resultado grande mas finito
- **WHEN** a operação de soma é pedida com os operandos 1e308 e 1e307
- **THEN** o resultado devolvido é finito
- **AND** nenhum erro é devolvido
