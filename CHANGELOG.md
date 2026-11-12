# Changelog

Formato baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/).

## [Não lançado]

### Onda 0 — preparação do repositório e do harness

- Adicionados `README.md` e `CHANGELOG.md` iniciais.
- Adicionado `apm.yml` (`calculadora_go` 0.0.1, hosts `claude` e `copilot`, dependência `emilianoeloi/apm#v1.0.0-alpha.1`) e `apm.lock.yaml`.
- Harness implantado por `apm install --target claude,copilot`: skills, agentes `sdlc-*`, prompts `/opsx-*` e hooks (`.claude/`, `.github/`, `.agents/`). Os scripts de hook foram revisados antes do commit.
- Estado do OpenSpec criado com `openspec init --tools none` (idioma: português) e pastas `spdd/analysis/` e `spdd/prompt/`.
- Adicionados `harness.config.json`, `CLAUDE.md` e `scripts/harness/phase_cli.py`, a partir dos modelos do skill `install-dev-harness`.
- `.gitignore` com `apm_modules/` e `.harness/`.
- Código da calculadora: nenhum (decisão desta onda).
- Não informado: `sensors_command` (vazio até haver código; `go test ./...` vira o comando com a primeira funcionalidade), lint, build e cobertura.

#### Versão do harness instalado

| Componente | Versão |
| --- | --- |
| Pacote `emilianoeloi/apm` | `v1.0.0-alpha.1` (commit `8c3f94a`) |
| `apm` (CLI) | 0.35.0 |
| `openspec` (CLI) | 1.14.1 |
| `python3` | 3.14.5 |
| `go` | 1.27.1 |

#### Resultado da verificação do skill `install-dev-harness`

| Verificação | Resultado |
| --- | --- |
| `session-start.py` imprime JSON com `additionalContext` | Passou |
| `openspec list` | Passou (nenhuma change ativa, saída 0) |
| `phase_cli.py show` | Passou (fase não declarada) |
| `spdd-*/SKILL.md` nos dois hosts (`.claude/skills/`, `.agents/skills/`) | Passou (5 skills em cada) |
| `opsx-*` nos dois hosts (`.claude/commands/`, `.github/prompts/`) | Passou (6 em cada) |
| Cada evento de hook uma vez em `.claude/settings.json` | Passou (SessionStart 1, PreToolUse 2 com matchers distintos, Stop 1) |
| Stop com `sensors_command` vazio | Passou (saída 0) |
| Guardas: `git push` e escrita fora da raiz | Bloqueados (saída 2) |
| Arquivos legados (`.claude/commands/opsx/`, `spdd-*.md`, `.github/copilot-prompts/`) | Ausentes |
