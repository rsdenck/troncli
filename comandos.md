# Comandos NUX

Documentação completa de todos os comandos disponíveis na CLI NUX.

## Visão Geral

A CLI NUX oferece os seguintes comandos para administração de sistemas Linux:

### Agent
IA para análise e automação de sistemas.
Uso: `nux agent [query]`

### Ask
Fazer perguntas a provedores de IA (Ollama, OpenAI, Claude, NVIDIA Build).
Uso: `nux ask <pergunta>`

### Audit
Auditoria de segurança.
Uso: `nux audit`

### Bash
Executar comandos Bash.
Uso: `nux bash <comando>`

### Clean
Limpa cache de pacotes.
Uso: `nux clean`

### Completion
Gerar scripts de autocompletar para shell.
Uso: `nux completion [shell]`

### Container
Gerenciamento de containers.
Uso: `nux container <subcomando> [args]`

### Disk
Gerenciamento de discos.
Uso: `nux disk <subcomando> [args]`

### Doctor
Verificação de saúde do sistema.
Uso: `nux doctor`

### Firewall
Gerenciamento de firewall.
Uso: `nux firewall <subcomando> [args]`

### GeoIP
Geolocalização e análise de IP.
Uso: `nux geoip <ip>`

### Help
Ajuda sobre qualquer comando.
Uso: `nux help [comando]`

### Install
Instala pacotes (universal, suporta apt, dnf, yum, pacman, zypper, apk).
Uso: `nux install <pacotes>`

### Network
Gerenciamento de rede.
Uso: `nux network <subcomando> [args]`

### Onboard
Experiência premium de onboarding.
Uso: `nux onboard`

### Plugin
Gerenciamento de plugins.
Uso: `nux plugin <subcomando> [args]`

### Process
Gerenciamento de processos.
Uso: `nux process <subcomando> [args]`

### Proton
Integração com ecossistema Proton (VPN, Pass, Drive).
Uso: `nux proton <subcomando> [args]`

### Remote
Conexões SSH remotas.
Uso: `nux remote <subcomando> [args]`

### Remove
Remove pacotes (universal).
Uso: `nux remove <pacotes>`

### Restart
Reinicia um serviço (alias para 'service restart').
Uso: `nux restart <serviço>`

### Search
Pesquisa em arquivos, pacotes, processos, serviços e mais.
Uso: `nux search <termo>`

### Service
Gerenciamento de serviços.
Uso: `nux service <subcomando> [args]`

### Skill
Gerencia skills NUX (integrações CLI externas).
Uso: `nux skill <subcomando> [args]`

### Start
Inicia um serviço (alias para 'service start').
Uso: `nux start <serviço>`

### Status
Mostra status de serviço (alias para 'service status').
Uso: `nux status <serviço>`

### Stop
Para um serviço (alias para 'service stop').
Uso: `nux stop <serviço>`

### System
Informações do sistema.
Uso: `nux system`

### Update
Atualiza lista de pacotes.
Uso: `nux update`

### Upgrade
Atualiza pacotes instalados.
Uso: `nux upgrade`

### Users
Gerenciamento de usuários.
Uso: `nux users <subcomando> [args]`

### Vault
Gerencia vault NUX para segredos e API keys.
Uso: `nux vault <subcomando> [args]`

## Flags Globais

- `--dry-run`: Simula execução sem fazer alterações
- `--json`: Saída em formato JSON
- `--log-file string`: Caminho do arquivo de log
- `--no-color`: Desativa saída colorida
- `--quiet`: Suprime saída
- `--timeout int`: Tempo limite em segundos (default 30)
- `--verbose`: Ativa log detalhado
- `-v, --version`: Mostra versão
- `--yaml`: Saída em formato YAML
