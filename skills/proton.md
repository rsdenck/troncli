---
name: proton
category: security
subcategory: vpn-privacy
repo: https://github.com/ProtonVPN
website: https://protonvpn.com
description: Integração com serviços Proton (VPN, Pass, Drive) para gerenciamento seguro, privacidade e auditoria de viagens.
install:
  rocky_linux:
    - "echo 'Instale o Proton VPN CLI via: https://protonvpn.com/download-linux'"
  ubuntu:
    - "echo 'Instale o Proton VPN CLI via: https://protonvpn.com/download-linux'"
  check:
    - "which protonvpn-cli || echo 'ProtonVPN CLI não instalado'"

commands:
  status:
    - "protonvpn-cli status || echo 'ProtonVPN não instalado'"
  vpn_connect:
    - "protonvpn-cli connect || echo 'Falha ao conectar'"
  vpn_fastest:
    - "protonvpn-cli connect --fastest || echo 'Falha ao conectar no mais rápido'"
  vpn_disconnect:
    - "protonvpn-cli disconnect || echo 'Não conectado'"
  open_mail:
    - "xdg-open https://mail.protonmail.com || echo 'Falha ao abrir Proton Mail'"
  open_drive:
    - "xdg-open https://drive.protonmail.com || echo 'Falha ao abrir Proton Drive'"
  vault_sync:
    - "nux vault sync proton || echo 'Vault não configurado'"
  pass_lookup:
    - "nux proton pass lookup $1 || echo 'Proton Pass não configurado'"
  secure_travel:
    - "nux proton vpn fastest"
    - "nux firewall list"
    - "systemctl disable --now telnet 2>/dev/null || true"
    - "echo 'Modo de viagem ativado: VPN conectado, firewall ativo, serviços inseguros desabilitados'"

nux:
  install:
    - "nux skill install proton"
  enable:
    - "nux skill enable proton"
  status:
    - "nux proton status"
  vpn:
    - "nux proton vpn connect"
    - "nux proton vpn fastest"
    - "nux proton vpn disconnect"
  open:
    - "nux proton open mail"
    - "nux proton open drive"
  sync:
    - "nux proton sync"
  travel:
    - "nux proton secure travel-mode"

tags:
  - proton
  - vpn
  - privacy
  - security
  - vault
  - travel-mode
---

# Proton Skill para NUX

## Visão Geral
Integração com ecossistema Proton para administradores Linux focados em privacidade e segurança.

## Funcionalidades

### Fase 1 — Integração Local/Segura
- Status da conexão VPN
- Conectar VPN (mais rápido disponível)
- Abrir Proton Mail/Drive via navegador
- Gerenciamento local de status

### Fase 2 — Vault Complementar
- Sincronização opcional com Proton Pass
- Lookup de senhas via linha de comando

### Fase 3 — Security Bundle (Travel Mode)
Executa sequência de segurança:
1. Conecta Proton VPN
2. Ativa DNS seguro
3. Checa firewall
4. Desabilita serviços inseguros (telnet, etc)

## Comandos NUX
```bash
nux proton status
nux proton vpn fastest
nux proton open mail
nux proton sync
nux proton secure travel-mode
```

## Instalação
```bash
nux skill install proton
nux skill enable proton
```

## Requisitos
- Conta Proton ativa
- ProtonVPN CLI instalado (opcional para comandos VPN)
- Navegador padrão configurado (para abrir Mail/Drive)
