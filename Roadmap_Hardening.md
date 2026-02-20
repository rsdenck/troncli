# Roadmap de Hardening da TronCLI

Este documento define o plano de melhorias urgentes para garantir a segurança, estabilidade e manutenibilidade da TronCLI.

## 🚨 Fase 1 — Segurança Imediata (Blocking Release)
**Objetivo:** Garantir que a execução de plugins e extensões seja segura e verificável.

- [x] **Remover execução direta de plugins externos:** Não permitir execução arbitrária de URLs sem validação.
- [x] **Exigir SHA256 obrigatório:** O registro de plugins deve conter o hash de verificação.
- [x] **Validar HTTPS obrigatório:** Bloquear qualquer URL que não use HTTPS.
- [x] **Implementar verificação de assinatura (Opcional):** Suporte a GPG ou Cosign no futuro.
- [x] **Bloquear execução se checksum não bater:** Abortar imediatamente a instalação/execução.
- [x] **Permissões mínimas:** Garantir que arquivos salvos tenham permissões restritivas (0755 ou 0700).

## 🤖 Fase 2 — Remover Mocks de IA
**Objetivo:** Eliminar código simulado e implementar clientes reais ou desabilitar funcionalidades.

- [x] **Implementação Real:** Substituir `// TODO` e strings fixas por clientes HTTP reais para Ollama, OpenAI e Claude.
- [x] **Feature Flags:** Se a chave de API ou endpoint não estiver configurado, o adaptador deve retornar erro ou não ser inicializado, em vez de fingir sucesso.
- [x] **Remoção de Mocks:** Código morto ou simulado deve ser removido da base de código.

## ⏱️ Fase 3 — Corrigir Cron Race Condition
**Objetivo:** Evitar corrupção do crontab em ambientes concorrentes.

- [x] **File Locking:** Implementar `flock` (syscall ou utilitário) ao ler e escrever no crontab.
- [x] **Atomicidade:** Garantir que a escrita seja atômica (escrever em temp e mover).

## 🛠️ Fase 4 — Refatorar Parsers
**Objetivo:** Tornar a leitura de status do sistema robusta e independente de formatação visual.

- [x] **Systemd JSON:** Usar `systemctl list-units ... --output=json`.
- [x] **Fallback Robusto:** Se JSON não disponível, usar regex estrito e não `strings.Fields` em colunas fixas.

## 👁️ Fase 5 — Auditoria Real
**Objetivo:** Auditoria baseada em dados estruturados do sistema, não em logs de texto.

- [x] **Journald JSON:** Usar `journalctl -o json`.
- [x] **Campos Estruturados:** Filtrar por `PRIORITY`, `SYSLOG_IDENTIFIER`, `_UID`, etc.
- [x] **Decodificação Struct:** Mapear JSON para structs Go fortemente tipadas.

## 🏗️ Mudanças Arquiteturais Futuras
1. [x] **Plugin Sandbox:** Execução isolada de plugins (limites de output/tempo).
2. [x] **Capability Registry Fail-Fast:** Abortar inicialização se regras de segurança falharem.
3. [x] **Configuração Externa:** Mover registros hardcoded para arquivos de configuração JSON/YAML.
4. [x] **Logging Estruturado:** Migrar para `slog` ou `zerolog`.
5. [x] **Cobertura de Testes:** Atingir >80% de cobertura.
