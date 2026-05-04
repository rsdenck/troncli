# NUX CLI

![NUX Logo](https://img.shields.io/badge/NUX-Professional-orange?style=for-the-badge)
![Version](https://img.shields.io/badge/version-0.3.0--beta-orange?style=for-the-badge)
![License](https://img.shields.io/badge/license-MIT-black?style=for-the-badge)
![Go Version](https://img.shields.io/badge/go-1.21+-orange?style=for-the-badge)
![Platform](https://img.shields.io/badge/platform-Linux-black?style=for-the-badge)

---

<p align="center">
  <img src="https://img.shields.io/badge/NUX-Big_Boss_for_Linux-orange?style=for-the-badge&labelColor=black">
</p>

<h1 align="center">🖥️ NUX CLI - Big Boss para Linux & SysAdmins</h1>

<p align="center">
  <b>CLI profissional para administração de sistemas Linux com integração nativa de IA</b><br/>
  <i>Ollama, NVIDIA Build, OpenAI, Claude - Tudo em um único lugar</i>
</p>

---

## 🚀 Sobre o Projeto

O **NUX** é uma CLI (Command Line Interface) de próxima geração para administração de sistemas Linux. Construído em Go, oferece uma arquitetura módular com suporte nativo a múltiplos provedores de IA.

### ✨ Características Principais

- 🤖 **Provedores de IA Integrados**: Ollama, OpenAI, Claude, NVIDIA Build
- 🔧 **Gerenciamento de Skills**: +120 skills para ferramentas externas
- 📦 **Instalação Automática**: Download e deploy automático de ferramentas
- 🎨 **Output Profissional**: Formato GCX JSON padronizado
- 🔒 **Vault Seguro**: Gerenciamento de tokens e configurações
- 🌍 **Multi-distribuição**: Compatível com apt, yum, dnf, pacman

---

## 📦 Instalação

```bash
# Download binário (em breve)
wget https://github.com/rsdenck/nux/releases/latest/download/nux-linux-amd64 -O nux
chmod +x nux
sudo mv nux /usr/local/bin/
nux --version
```

---

## 🚀 Uso Básico

### Comandos de IA

```bash
# Consultar Ollama (padrão)
nux ask "Como listar arquivos no Linux?"

# Usar NVIDIA Build
nux ask query --provider nvidia --model minimaxai/minimax-m2.7 "Consulta NVIDIA"

# Configurar provedores
nux ask config --provider ollama --host http://192.168.130.25:11434
nux ask config --provider nvidia --api-key "seu-token"
```

### Gerenciamento de Skills

```bash
# Listar skills disponíveis
nux skill list

# Instalar skill
nux skill install ansible

# Info da skill
nux skill info terraform
```

---

## 🏗️ Arquitetura

```
nux/
├── cmd/nux/commands/    # Comandos Cobra
├── internal/            # Núcleo (core, modules, vault)
├── skills/              # Definições de skills (.md)
└── tests/               # Testes Go e k6
```

---

## 🤖 Provedores Suportados

| Provedor | Status | Modelo Padrão |
|----------|--------|---------------|
| **Ollama** | ✅ Funcional | `gemma:2b` |
| **NVIDIA Build** | ✅ Funcional | `minimaxai/minimax-m2.7` |
| **OpenAI** | ⚠️ Configurável | `gpt-3.5-turbo` |
| **Claude** | ⚠️ Configurável | `claude-3-5-sonnet` |

---

## 📚 Skills Disponíveis

O NUX possui **120+ skills** organizadas por categorias:

- `infrastructure/` - Terraform, Pulumi, Packer
- `automation/` - Ansible, Puppet, Chef
- `ci-cd/` - CircleCI, Jenkins, GitLab
- `git/` - GitHub CLI, GitLab CLI
- `testing/` - k6, JMeter
- `kubernetes/` - kubectl, Helm, K9s
- `containers/` - Docker, Podman, Buildah
- `security/` - Nmap, Metasploit, John

Scripts de instalação: [rsdenck/skillnux](https://github.com/rsdenck/skillnux)

---

## 🧪 Testes

```bash
# Testes Go
go test ./tests/...

# Testes de carga (k6)
k6 run tests/k6-load-test.js
```

---

## 🤝 Contribuindo

1. Fork o projeto
2. Crie sua branch (`git checkout -b feature/nova-skill`)
3. Commit suas mudanças (`git commit -m 'feat: nova skill'`)
4. Push para a branch (`git push origin feature/nova-skill`)
5. Abra um Pull Request

---

## 📄 Licença

Este projeto está sob a licença **MIT** - veja o arquivo [LICENSE](LICENSE) para detalhes.

---

<p align="center">
  <img src="https://img.shields.io/badge/Made_with-❤️-black?style=for-the-badge&labelColor=orange">
  <img src="https://img.shields.io/badge/For-Linux-orange?style=for-the-badge&labelColor=black">
</p>
