---
name: cve
binary: cve-bin-tool
category: security
subcategory: vulnerability-scanning
repo: https://github.com/ossf/cve-bin-tool
website: https://github.com/ossf/cve-bin-tool
language: python
license: GPL-3.0

description: Ferramenta para análise de binários, pacotes e diretórios em busca de componentes vulneráveis e CVEs conhecidas.

install:
  rocky_linux:
    - sudo dnf install -y python3 python3-pip
    - pip3 install --upgrade cve-bin-tool
  ubuntu:
    - sudo apt update
    - sudo apt install -y python3 python3-pip
    - pip3 install --upgrade cve-bin-tool
  pipx:
    - pipx install cve-bin-tool
  container:
    - podman run --rm ghcr.io/ossf/cve-bin-tool/cve-bin-tool --help

commands:
  verify:
    - cve-bin-tool --version
  scan_binary:
    - cve-bin-tool /usr/bin/ssh
  scan_directory:
    - cve-bin-tool /opt/app/
  scan_package_file:
    - cve-bin-tool package.rpm
  output_json:
    - cve-bin-tool /usr/bin/ssh -f json -o report.json
  output_csv:
    - cve-bin-tool /usr/bin/ssh -f csv -o report.csv
  severity_filter:
    - cve-bin-tool /opt/app --severity high
  offline_db:
    - cve-bin-tool --update now

nux:
  install:
    - nux skill install cve
  info:
    - nux skill info cve
  scan:
    - nux cve scan /usr/bin/ssh
  system:
    - nux cve system
  report:
    - nux cve report

tags:
  - cve
  - security
  - binary
  - scanner
  - compliance
  - devsecops
  - linux
---

# CVE Bin Tool

## Visão Geral

O `cve-bin-tool` é uma ferramenta poderosa mantida pela OpenSSF para detectar componentes vulneráveis em:

- Binários compilados
- Diretórios de aplicações
- Imagens extraídas
- Pacotes `.rpm`, `.deb`, `.apk`
- Bibliotecas conhecidas
- Software legado sem SBOM

Excelente para ambientes Linux, servidores, auditoria e hardening.

## Repositório Oficial

https://github.com/ossf/cve-bin-tool

## Casos de Uso

- Auditar servidores Linux
- Verificar binários suspeitos
- Analisar software legado
- Pipeline CI/CD
- Compliance corporativo
- Validação antes de produção

## Comandos Essenciais

```bash
cve-bin-tool --version
cve-bin-tool /usr/bin/ssh
cve-bin-tool /usr/sbin/httpd
cve-bin-tool /opt/app
cve-bin-tool package.rpm
```

## Exportar Relatórios

```bash
cve-bin-tool /usr/bin/ssh -f json -o report.json
cve-bin-tool /usr/bin/ssh -f csv -o report.csv
cve-bin-tool /opt/app -f html -o report.html
```

## Atualizar Base de CVEs

```bash
cve-bin-tool --update now
```

## Integração com NUX

```bash
nux skill install cve
nux cve scan /usr/bin/ssh
nux cve system
nux cve report
```

## Ideias de Módulos NUX

Escanear binários críticos:
```bash
nux cve scan /usr/bin/ssh
nux cve scan /usr/bin/sudo
nux cve scan /usr/sbin/nginx
```

Auditoria completa do host:
```bash
nux cve system
```

Executaria scans em:
- /usr/bin
- /usr/sbin
- Pacotes instalados
- Serviços ativos
- Módulo compliance

```bash
nux cve report --json
nux cve report --html
```

## Boas Práticas

- Atualizar base CVE frequentemente
- Validar falsos positivos
- Integrar em CI/CD
- Priorizar CVSS alto/crítico
- Manter inventário de binários
