# ROADMAP TRONCLI

> "A ferramenta mais completa já criada para Sysadmins Linux."

Nossa visão é construir uma TUI unificada e de nível empresarial que centralize todas as tarefas críticas de administração Linux.

---

## Fase 1: Fundação (Atual)
- [x] **Estrutura do Projeto**: Implementação da Clean Architecture.
- [x] **Pipeline CI/CD**: Workflow completo do GitHub Actions.
- [x] **Módulos Principais**: Sistema, Disco, Rede, Usuários.
- [x] **Framework TUI**: Integração tview/tcell com tema TRON.

## Fase 2: Integração Profunda com Kernel
- [ ] **Gerenciador LVM**: Criar/Redimensionar PVs, VGs, LVs.
- [ ] **Sistema de Auditoria**: Parsear `/var/log/auth.log` e systemd journal.
- [ ] **Análise de Rede**: Estatísticas de soquetes (ss), integração nftables.
- [ ] **Matador de Processos**: Kill/renice interativo.

## Fase 3: Recursos Empresariais
- [ ] **Gerenciamento Remoto**: Túneis SSH, integração rsd-sshm.
- [ ] **Gerenciamento de Containers**: TUI para Podman/Docker.
- [ ] **Unidades Systemd**: Iniciar/Parar/Habilitar/Editar serviços.
- [ ] **Cron/Timers**: Gerenciamento de agendamentos.

## Fase 4: Polimento Final
- [ ] **Documentação**: Man pages, Wiki.
- [ ] **Empacotamento**: DEB, RPM, AUR, Snap.
- [ ] **Temas**: Esquemas de cores personalizados (Cyberpunk, Matrix).

---
TRONCLI | EXECUÇÃO FUTURA
