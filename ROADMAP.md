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
- [x] **Gerenciador LVM**: Criar/Redimensionar PVs, VGs, LVs.
- [x] **Sistema de Auditoria**: Parsear `/var/log/auth.log` e systemd journal.
- [x] **Análise de Rede**: Estatísticas de soquetes (ss), integração nftables.
- [x] **Matador de Processos**: Kill/renice interativo e limpeza automática de Zumbis.

## Fase 3: Recursos Empresariais
- [x] **Gerenciamento Remoto**: Túneis SSH, integração rsd-sshm.
- [x] **Gerenciamento de Containers**: TUI para Podman/Docker.
- [x] **Unidades Systemd**: Iniciar/Parar/Habilitar/Editar serviços.
- [x] **Cron/Timers**: Gerenciamento de agendamentos.

## Fase 4: Abstração Universal (Multi-Distro)
- [x] **Gerenciador de Pacotes Universal**: Camada unificada para `apt`, `dnf`, `yum`, `pacman`, `zypper`, `apk` (install, remove, update, search).
- [x] **Firewall Unificado**: Abstração para `nftables`, `iptables`, `firewalld`, `ufw` (allow, block, list).
- [x] **Rede Universal**: Configuração de IP, DHCP, DNS, Gateway, Rotas (suporte a `netplan`, `ifcfg`, `interfaces`, `NM`, `systemd-networkd`).
- [x] **Usuários e Permissões Avançadas**: Auditoria de UIDs, senhas expiradas, grupos privilegiados, chaves SSH inválidas.
- [x] **Modo Auditor Universal**: Detecção automática de logs (`auth.log`, `secure`, `journald`) e análise de falhas/login.
- [x] **Processos Avançados (Cross-Distro)**: Árvore de processos, zombies, openfiles, ports (via `ps`, `ss`, `lsof`, `/proc`).
- [x] **Disco Universal (Growth Analysis)**: Detecção de crescimento de diretórios, top arquivos, órfãos, inodes críticos (sem snapshot).
- [x] **Modo "Compatibilidade Bash"**: Auditoria de `.bashrc`, aliases perigosos, PATH, permissões.
- [x] **Gestão de Serviços Multi-Init**: Abstração para `systemd`, `service`, `rc-service`, `runit` (start, enable, status).
- [x] **Modo Doctor Multi-Distro**: Checks universais de saúde (Load, Swap, TCP, Disco, Kernel, Permissões).
- [x] **Compatibilidade de Ambientes**: Adaptação para WSL, Docker, Kubernetes Node, VM, Bare Metal.
- [ ] **Sistema de Plugins**: Instalação de plugins específicos por distro (`troncli plugin install arch`).

## Fase 5: Polimento e Distribuição
- [ ] **Documentação**: Man pages, Wiki.
- [ ] **Empacotamento**: DEB, RPM, AUR, Snap.
- [ ] **Temas**: Esquemas de cores personalizados (Cyberpunk, Matrix).

---
TRONCLI | EXECUÇÃO FUTURA
