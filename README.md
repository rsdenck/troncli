# NUX - Linux CLI Master & Manager

<p align="center">
  <img src="https://img.shields.io/badge/version-0.2.19-green?style=for-the-badge" alt="Version">
  <img src="https://img.shields.io/badge/license-MIT-yellow?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/linux-333?style=for-the-badge&logo=linux" alt="Linux">
</p>

<pre align="center" style="color:#c9a01a; font-family:'Fira Code',monospace; font-size:12px; line-height:1.1; background:#050505; padding:20px; border-radius:6px; border:1px solid #1a1a1a;">
 ███╗   ██╗██╗   ██╗██╗  ██╗
 ████╗  ██║██║   ██║╚██╗██╔╝
 ██╔██╗ ██║██║   ██║ ╚███╔╝
 ██║╚██╗██║██║   ██║ ██╔██╗
 ██║ ╚████║╚██████╔╝██╔╝ ██╗
 ╚═╝  ╚═══╝ ╚═════╝ ╚═╝  ╚═╝
</pre>

<p align="center" style="color:#a0a0a0; font-size:1.1rem; max-width:600px; margin:0 auto 30px;">
Production-grade CLI para administração completa de sistemas Linux.<br>
<strong style="color:#c9a01a;">Multi-distribuição • Skill Engine • IA Integrada • Gestão de Discos & LVM</strong>
</p>

<p align="center">
  <a href="#instalação-rápida" style="background:#c9a01a; color:#000; padding:12px 30px; text-decoration:none; border-radius:4px; font-weight:bold; margin:5px; display:inline-block;">INSTALAR AGORA</a>
  <a href="https://github.com/rsdenck/nux" style="background:#0f0f0f; color:#c9a01a; padding:12px 30px; text-decoration:none; border-radius:4px; font-weight:bold; margin:5px; display:inline-block; border:1px solid #c9a01a;">VER NO GITHUB</a>
</p>

---

## Sobre o NUX

O **NUX** é um CLI Master para Linux que centraliza operações complexas de sysadmin em uma única ferramenta moderna. Projetado para profissionais que precisam gerenciar múltiplas distribuições com eficiência, output padronizado (GCX JSON) e extensibilidade via Skill Engine.

### Por que NUX?

- **Multi-distribuição nativa**: Detecta e opera em qualquer distro (Debian, RHEL, Arch, Alpine, openSUSE)
- **Output GCX padronizado**: JSON estruturado para automação e integração
- **Skill Engine**: 166+ skills prontas, extensíveis via arquivos markdown
- **IA integrada**: Ollama com qwen3-coder para sugestões inteligentes
- **Gestão completa**: Discos, LVM, rede, serviços, pacotes, containers, segurança

---

## Recursos Principais

<table>
<tr>
<td width="50%">

### Core System
- **Service Management**: systemd, openrc, sysvinit, runit
- **Package Managers**: apt, dnf, yum, pacman, apk, zypper
- **Disk Management**: listagem, usage, LVM (pv/vg/lv)
- **Filesystem Ops**: mount, resize, verify
- **Process Control**: monitoramento, sinais, prioridades

</td>
<td width="50%">

### Network & Security
- **Network Config**: interfaces, bridges, bonds, VLANs
- **Firewall**: iptables, nftables, ufw, firewalld
- **Security Audit**: CVE check, hardening, compliance
- **User Management**: users, groups, sudo, permissions
- **SSH & VPN**: configuração e tunnel management

</td>
</tr>
<tr>
<td width="50%">

### Cloud & Containers
- **Docker/Podman**: containers, images, networks
- **Kubernetes**: kubectl integration via skills
- **Cloud APIs**: AWS, GCloud, Azure via skills
- **Terraform**: infraestrutura como código
- **Virtualization**: KVM, QEMU, libvirt

</td>
<td width="50%">

### Developer Tools
- **Shell Integration**: bash, zsh, fish, tmux
- **Git Operations**: via skills
- **AI Agent**: Ollama com qwen3-coder
- **Autocompletion**: bash, zsh, fish
- **CI/CD**: GitHub Actions, GoReleaser

</td>
</tr>
</table>

---

## Instalação Rápida

### Via Package Manager (Recomendado)

```bash
# Debian / Ubuntu
dpkg -i nux_*.deb

# RHEL / Fedora / CentOS
rpm -ivh nux-*.rpm

# Arch Linux
pacman -U nux-*.pkg.tar.zst

# Alpine
apk add --allow-untrusted nux-*.apk
```

### Via Go (Do código fonte)

```bash
git clone https://github.com/rsdenck/nux.git
cd nux
go build -o nux ./cmd/nux
sudo mv nux /usr/local/bin/
```

### Via Release (Binário)

Baixe a última release em: [https://github.com/rsdenck/nux/releases](https://github.com/rsdenck/nux/releases)

---

## Primeiros Passos

### 1. Configuração Inicial (Onboard)

```bash
nux onboard
```

Este comando configura:
- Vault seguro em `~/.skills/.nux.json` (permissões 0600)
- Conexão com Ollama (local ou remoto)
- Detecção automática da distribuição
- Skill Engine inicializado

### 2. Verificação de Saúde

```bash
nux doctor
```

Retorna status JSON:
```json
{
  "status": "success",
  "data": {
    "vault": "ok",
    "ollama": "connected",
    "distribution": "rhel",
    "package_manager": "dnf"
  }
}
```

### 3. Listar e Instalar Skills

```bash
# Listar todas as 166+ skills disponíveis
nux skill list

# Instalar skill específica
nux skill install docker

# Habilitar skill
nux skill enable docker
```

---

## Exemplos de Uso

### Gestão de Discos e LVM

```bash
# Listar dispositivos de disco
nux disk list

# Verificar uso de disco
nux disk usage /

# Criar volume LVM (simulado)
nux disk lvm create /dev/sdb 50GB

# Rescan SCSI para novos discos
nux disk rescan
```

### Operações com Pacotes

```bash
# Instalar pacote
nux pkg install nginx

# Atualizar sistema
nux pkg update

# Listar pacotes instalados
nux pkg list --json
```

### Rede e Firewall

```bash
# Listar interfaces
nux network list

# Configurar IP estático
nux network set eth0 --ip 192.168.1.100 --netmask 255.255.255.0

# Regras de firewall
nux firewall add --port 80 --protocol tcp
```

### Usando o Agente IA

```bash
# Consultar Ollama
nux agent ask "como criar um volume lvm de 100gb"

# Status do agente
nux agent status
```

---

## Formato de Output (GCX)

Todos os comandos seguem o padrão GCX JSON:

```json
{
  "status": "success",
  "data": {
    "items": [...]
  },
  "total": 10,
  "message": "operação realizada com sucesso"
}
```

Flags globais disponíveis:
- `--json`: Força output JSON
- `--yaml`: Output em YAML
- `--quiet`: Suprime output
- `--dry-run`: Simula execução
- `--verbose`: Log detalhado
- `--no-color`: Desabilita cores
- `--log-file`: Arquivo de log

---

## Skill Engine

O NUX usa arquivos `.md` na pasta `skills/` para definir skills externas. Cada skill contém:

```markdown
# Skill: docker

- **Repository**: https://github.com/docker/cli
- **Description**: Docker container management
- **Install Command**: dnf install docker-ce
- **Type**: container
- **Dependencies**: systemd
```

### Categorias de Skills Disponíveis

| Categoria | Exemplos |
|-----------|---------|
| **Shell** | bash, zsh, fish, tmux |
| **System** | docker, kubectl, systemd, lvm |
| **Network** | curl, wget, nmap, tcpdump |
| **Cloud** | aws, gcloud, azure, terraform |
| **AI** | ollama, openai, claude, aider |
| **Security** | fail2ban, auditd, clamav |
| **Database** | mysql, postgres, redis, mongodb |
| **Monitoring** | prometheus, grafana, nagios |

---

## Vault de Configuração

O NUX usa um vault seguro em `~/.skills/.nux.json`:

```json
{
  "version": "1.0.0",
  "installed_skills": ["docker", "kubectl"],
  "enabled_skills": ["docker"],
  "api_keys": {},
  "ollama": {
    "host": "http://192.168.130.25:11434",
    "model": "qwen3-coder",
    "enabled": true
  },
  "vault_mode": true
}
```

Permissões: `0600` (apenas o usuário tem acesso)

---

## Suporte Multi-Distribuição

| Distribuição | Package Manager | Init System | Status |
|--------------|-----------------|-------------|--------|
| Debian/Ubuntu | apt | systemd | ✅ Testado |
| RHEL/Fedora/CentOS | dnf/yum | systemd | ✅ Testado |
| Arch Linux | pacman | systemd | ✅ Testado |
| Alpine | apk | openrc | ✅ Testado |
| openSUSE | zypper | systemd | ✅ Testado |

---

## CI/CD e Releases

O NUX utiliza:
- **GitHub Actions**: Builds automáticos multi-plataforma
- **GoReleaser**: Geração de pacotes (.deb, .rpm, .pkg.tar.zst, .apk)
- **Dependabot**: Atualizações automáticas de dependências

### Criar uma Release

```bash
# Atualizar versão
vim .goreleaser.yaml  # alterar version
git commit -am "Release v0.2.19"
git tag v0.2.19
git push origin v0.2.19
```

O GoReleaser criará automaticamente os pacotes para todas as distribuições.

---

## Contribuindo

Contribuições são bem-vindas! Por favor:

1. Fork o repositório
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

### Padrões de Código

- Linguagem: Go 1.21+
- Output: Sempre usar pacote `internal/output` (GCX format)
- Commits: Seguir conventional commits
- Testes: Implementar testes para novos comandos
- Skills: Adicionar arquivo .md na pasta skills/

---

## Licença

Distribuído sob a licença **MIT**. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

## Links Úteis

- **Repositório**: [https://github.com/rsdenck/nux](https://github.com/rsdenck/nux)
- **Issues**: [https://github.com/rsdenck/nux/issues](https://github.com/rsdenck/nux/issues)
- **Releases**: [https://github.com/rsdenck/nux/releases](https://github.com/rsdenck/nux/releases)
- **GitHub Pages**: [https://rsdenck.github.io/nux/](https://rsdenck.github.io/nux/)
- **Documentação**: [docs/](docs/)

---

<p align="center" style="color:#4a4a4a; font-size:0.9rem; margin-top:50px;">
  Desenvolvido para Linux por <a href="https://github.com/rsdenck" style="color:#c9a01a;">@rsdenck</a><br>
  <span style="color:#c9a01a;">NUX v0.2.19</span> • Produção • 2026
</p>
