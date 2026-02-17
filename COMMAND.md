# TRONCLI Command Reference

Generated on Tue, 17 Feb 2026 07:56:03 -03

## Main Help
```
TRONCLI is a production-grade Linux tool for system administration.
It features a TUI (Terminal User Interface) and a comprehensive CLI for automation.

Usage:
  troncli [flags]
  troncli [command]

Available Commands:
  audit       Auditoria de Segurança
  bash        Executar comandos e scripts Bash
  completion  Gerar scripts de autocompletar para shell
  container   Gerenciar containers (Docker/Podman)
  disk        Gerenciamento de Disco
  doctor      Saúde do Sistema
  firewall    Gerenciamento de Firewall
  group       Gerenciamento de Grupos
  help        Help about any command
  network     Gerenciamento de Rede
  pkg         Gerenciador de Pacotes Universal
  plugin      Gerenciar plugins do TRONCLI
  process     Gerenciamento de Processos do Sistema
  remote      Gerenciar conexões remotas SSH
  service     Gerenciar serviços do sistema
  system      Informações e Perfil do Sistema
  user        Gerenciamento de Usuários e Grupos

Flags:
      --dry-run       Simulate execution without making changes
  -h, --help          help for troncli
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli [command] --help" for more information about a command.
```

## Subcommands
### audit
```
Ferramentas para auditoria de segurança, logs e integridade do sistema.

Usage:
  troncli audit [command]

Available Commands:
  commands     Lista comandos executados (via sudo/logs)
  file-changes Monitora alterações recentes em arquivos
  logins       Lista histórico de logins
  sudoers      Lista usuários com privilégios sudo

Flags:
  -h, --help   help for audit

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli audit [command] --help" for more information about a command.
```

### bash
```
Executa comandos Bash diretamente ou scripts de arquivos, gerenciando permissões e execução.

Usage:
  troncli bash [command]

Available Commands:
  run         Executar um comando Bash
  script      Executar um script Bash de arquivo

Flags:
  -h, --help   help for bash

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli bash [command] --help" for more information about a command.
```

### completion
```
Para carregar o autocompletar:

Bash:
  $ source <(troncli completion bash)

Zsh:
  # Se o autocompletar do shell não estiver ativado, adicione ao .zshrc:
  # autoload -U compinit; compinit
  $ source <(troncli completion zsh)

Fish:
  $ troncli completion fish | source

PowerShell:
  PS> troncli completion powershell | Out-String | Invoke-Expression

Usage:
  troncli completion [bash|zsh|fish|powershell]

Flags:
  -h, --help   help for completion

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format
```

### container
```
Gerenciar ciclo de vida de containers, suportando Docker e Podman automaticamente.

Usage:
  troncli container [command]

Available Commands:
  list        Listar containers
  logs        Ver logs de um container
  start       Iniciar um container
  stop        Parar um container

Flags:
  -h, --help   help for container

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli container [command] --help" for more information about a command.
```

### disk
```
Gerencie discos, partições, e uso de espaço (usage, cleanup, health).

Usage:
  troncli disk [command]

Available Commands:
  cleanup     Limpa arquivos temporários e caches
  health      Verifica saúde do disco
  inodes      Exibe uso de inodes
  top-files   Lista maiores arquivos
  usage       Exibe uso do disco

Flags:
  -h, --help   help for disk

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli disk [command] --help" for more information about a command.
```

### doctor
```
Executa verificações de saúde do sistema (Load, Swap, Disco, TCP, etc).

Usage:
  troncli doctor [flags]

Flags:
      --disk       Executa verificações de disco
      --full       Executa todas as verificações
  -h, --help       help for doctor
      --network    Executa verificações de rede
      --security   Executa verificações de segurança

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format
```

### firewall
```
Controlar regras de firewall (ufw, firewalld, iptables, nftables).

Usage:
  troncli firewall [command]

Available Commands:
  allow       Permitir tráfego na porta
  deny        Bloquear tráfego na porta
  disable     Desabilitar firewall
  enable      Habilitar firewall
  list        Listar regras de firewall

Flags:
  -h, --help   help for firewall

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli firewall [command] --help" for more information about a command.
```

### group
```
Gerenciamento de Grupos

Usage:
  troncli group [command]

Available Commands:
  add         Adicionar grupo
  del         Remover grupo
  list        Listar grupos

Flags:
  -h, --help   help for group

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli group [command] --help" for more information about a command.
```

### network
```
Gerencie interfaces, rotas, DNS e configurações de rede.

Usage:
  troncli network [command]

Available Commands:
  capture     Capturar pacotes (tcpdump)
  dig         Consultar DNS (dig)
  info        Informações detalhadas de rede
  scan        Escanear portas (nmap)
  set-state   Alterar estado da interface
  sockets     Listar sockets abertos (ss)
  trace       Executar traceroute

Flags:
  -h, --help   help for network

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli network [command] --help" for more information about a command.
```

### pkg
```
Instala, remove e gerencia pacotes de forma transparente em apt, dnf, yum, pacman, apk e zypper.

Usage:
  troncli pkg [command]

Available Commands:
  install     Instala um pacote
  remove      Remove um pacote
  search      Pesquisa por pacotes
  update      Atualiza a lista de pacotes
  upgrade     Atualiza todos os pacotes do sistema

Flags:
  -h, --help   help for pkg

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli pkg [command] --help" for more information about a command.
```

### plugin
```
Instalar, listar e remover plugins (scripts ou binários) do TRONCLI.

Usage:
  troncli plugin [command]

Available Commands:
  install     Instalar um plugin
  list        Listar plugins instalados
  remove      Remover um plugin

Flags:
  -h, --help   help for plugin

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli plugin [command] --help" for more information about a command.
```

### process
```
Visualiza, finaliza e gerencia prioridade de processos em execução.

Usage:
  troncli process [command]

Available Commands:
  kill        Envia sinal para um processo (default SIGTERM)
  listening   Lista todas as portas em escuta no sistema
  open-files  Lista arquivos abertos por um processo
  ports       Lista portas ouvidas por um processo
  renice      Altera a prioridade de um processo
  tree        Exibe a árvore de processos
  zombies     Elimina processos zumbis

Flags:
  -h, --help   help for process

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli process [command] --help" for more information about a command.
```

### remote
```
Conectar, executar comandos e transferir arquivos via SSH.

Usage:
  troncli remote [command]

Available Commands:
  connect     Conectar a um host remoto (interativo)
  copy        Copiar arquivo para host remoto (SCP)
  exec        Executar comando em host remoto
  list        Listar perfis SSH configurados

Flags:
  -h, --help   help for remote

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli remote [command] --help" for more information about a command.
```

### service
```
Controlar serviços (systemd, sysvinit, openrc, runit) de forma unificada.

Usage:
  troncli service [command]

Available Commands:
  disable     Desabilitar serviço na inicialização
  enable      Habilitar serviço na inicialização
  list        Listar serviços
  logs        Ver logs do serviço
  restart     Reiniciar um serviço
  start       Iniciar um serviço
  status      Ver status detalhado de um serviço
  stop        Parar um serviço

Flags:
  -h, --help   help for service

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli service [command] --help" for more information about a command.
```

### system
```
Exibe informações detalhadas sobre o sistema, kernel, uptime e ambiente.

Usage:
  troncli system [command]

Available Commands:
  info        Exibe informações gerais do sistema
  kernel      Exibe versão do kernel
  profile     Exibe o perfil completo do sistema (JSON)

Flags:
  -h, --help   help for system

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli system [command] --help" for more information about a command.
```

### user
```
Gerencie usuários e grupos do sistema (add, del, modify, list).

Usage:
  troncli user [command]

Available Commands:
  add         Adicionar novo usuário
  del         Remover usuário
  list        Listar usuários
  modify      Modificar usuário existente

Flags:
  -h, --help   help for user

Global Flags:
      --dry-run       Simulate execution without making changes
      --json          Output in JSON format
      --no-color      Disable color output
      --quiet         Suppress output
      --timeout int   Timeout in seconds (default 30)
      --verbose       Enable verbose logging
      --yaml          Output in YAML format

Use "troncli user [command] --help" for more information about a command.
```

