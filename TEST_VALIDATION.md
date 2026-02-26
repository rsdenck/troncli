# TRONCLI - Guia de Validação e Testes

## 🎯 Objetivo
Validar que a TRONCLI está COMPLETAMENTE FUNCIONAL na CLI Linux com o TRON ROOT AGENT.

---

## 📋 Pré-requisitos

### Sistema Operacional
- Linux (Ubuntu, Debian, Fedora, Arch, Alpine, openSUSE, Gentoo, Void)
- Kernel 3.10+
- Arquitetura: x86_64 ou aarch64

### Ferramentas Necessárias
```bash
# Verificar se estão instaladas
which git
which make
which gcc
which g++
which wget  # ou curl
which go    # Go 1.24+
```

### Instalar dependências (se necessário)
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y git build-essential wget golang

# Fedora/RHEL
sudo dnf install -y git gcc gcc-c++ make wget golang

# Arch Linux
sudo pacman -S git base-devel wget go

# Alpine
sudo apk add git build-base wget go

# Gentoo
sudo emerge --ask dev-vcs/git sys-devel/gcc sys-devel/make net-misc/wget dev-lang/go

# Void Linux
sudo xbps-install -S git gcc make wget go
```

---

## 🚀 PASSO 1: Clonar e Compilar TRONCLI

```bash
# Clonar repositório
git clone https://github.com/rsdenck/troncli.git
cd troncli

# Checkout branch dev
git checkout dev

# Verificar que está na branch correta
git branch

# Compilar TRONCLI
go build -o troncli cmd/troncli/main.go

# Verificar compilação
./troncli --version
./troncli --help
```

**✅ Validação Esperada:**
- Binário `troncli` criado
- Comando `--help` mostra todos os comandos
- Sem erros de compilação

---

## 🧠 PASSO 2: Instalar llama.cpp

### Opção A: Instalação Automática (RECOMENDADO)
```bash
# Usar o comando de setup da TRONCLI
./troncli agent setup
```

**O que o comando faz:**
1. Cria diretórios em `~/.troncli/`
2. Clona llama.cpp
3. Detecta AVX2 e compila otimizado
4. Baixa modelo Qwen2.5-Coder-7B (~4GB)
5. Instala binário em `~/.troncli/bin/llama-cli`

### Opção B: Instalação Manual
```bash
# Criar diretórios
mkdir -p ~/.troncli/{bin,models}

# Clonar llama.cpp
git clone https://github.com/ggerganov/llama.cpp ~/.troncli/llama.cpp
cd ~/.troncli/llama.cpp

# Verificar suporte AVX2
grep -q avx2 /proc/cpuinfo && echo "AVX2 suportado" || echo "AVX2 não suportado"

# Compilar (com AVX2 se disponível)
if grep -q avx2 /proc/cpuinfo; then
    make LLAMA_NATIVE=1
else
    make
fi

# Copiar binário
cp llama-cli ~/.troncli/bin/ || cp main ~/.troncli/bin/llama-cli

# Verificar instalação
~/.troncli/bin/llama-cli --version
```

**✅ Validação Esperada:**
- Binário `llama-cli` criado
- Comando `--version` funciona
- Compilação sem erros

---

## 📥 PASSO 3: Baixar Modelo GGUF

### Opção A: Modelo Qwen2.5-Coder-7B (RECOMENDADO)
```bash
# Baixar modelo quantizado Q4_0 (~4GB)
cd ~/.troncli/models/
wget https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_0.gguf

# Verificar download
ls -lh qwen2.5-coder-7b-instruct-q4_0.gguf
```

### Opção B: Modelo Qwen3-Coder (Alternativo)
```bash
# Baixar Qwen3-Coder Q4_K_M
cd ~/.troncli/models/
wget https://huggingface.co/Qwen/Qwen3-Coder-GGUF/resolve/main/qwen3-coder-q4_k_m.gguf

# Verificar download
ls -lh qwen3-coder-q4_k_m.gguf
```

**✅ Validação Esperada:**
- Arquivo `.gguf` baixado (~4GB)
- Sem erros de download

---

## 🧪 PASSO 4: Teste Manual do llama.cpp

```bash
# Teste básico
~/.troncli/bin/llama-cli \
  -m ~/.troncli/models/qwen2.5-coder-7b-instruct-q4_0.gguf \
  -p "Explique o comando systemctl restart nginx" \
  -n 256 \
  --temp 0.2

# Teste com JSON
~/.troncli/bin/llama-cli \
  -m ~/.troncli/models/qwen2.5-coder-7b-instruct-q4_0.gguf \
  -p 'Responda em JSON: {"analysis": "teste", "commands": ["ls -la"]}' \
  -n 128 \
  --temp 0.2
```

**✅ Validação Esperada:**
- Modelo carrega sem erros
- Resposta gerada em texto
- Tempo de resposta < 30 segundos (CPU) ou < 5 segundos (GPU)

---

## 🎯 PASSO 5: Testar TRONCLI CLI (Comandos Básicos)

### Teste 1: System Info
```bash
./troncli system info
```
**Esperado:** Box-drawing perfeito com informações do sistema

### Teste 2: Service List
```bash
./troncli service list
```
**Esperado:** Lista de serviços com formatação profissional

### Teste 3: Process Tree
```bash
./troncli process tree
```
**Esperado:** Árvore de processos formatada

### Teste 4: Network Interfaces
```bash
./troncli network interfaces
```
**Esperado:** Lista de interfaces de rede

### Teste 5: Disk Usage
```bash
./troncli disk usage
```
**Esperado:** Uso de disco com box-drawing

### Teste 6: Package Manager Detection
```bash
./troncli pkg list
```
**Esperado:** Detecta gerenciador de pacotes correto

---

## 🤖 PASSO 6: Testar TRON ROOT AGENT

### Teste 1: Agent Status
```bash
./troncli agent status
```
**Esperado:** Status do agente configurado

### Teste 2: Root Agent - Comando Simples
```bash
./troncli agent root "verificar saúde do sistema"
```

**Esperado:**
```
┌── TRON ROOT AGENT › AUTONOMOUS MODE ─────────────┐
│                                                   │
│  Intent    › verificar saúde do sistema          │
│  Model     › Qwen2.5-Coder-7B                    │
│  Engine    › llama.cpp                           │
│  Mode      › Hardcore Linux                      │
│                                                   │
└───────────────────────────────────────────────────┘

🧠 Analyzing intent...

┌── AGENT ANALYSIS ────────────────────────────────┐
│                                                   │
│  Análise do que precisa ser feito                │
│                                                   │
└───────────────────────────────────────────────────┘

┌── RISK ASSESSMENT ───────────────────────────────┐
│                                                   │
│  Risk Level      › LOW                           │
│  Impact          › Read-only system check        │
│  Confirmation    › false                         │
│                                                   │
└───────────────────────────────────────────────────┘

┌── COMMANDS TO EXECUTE ───────────────────────────┐
│                                                   │
│  1. troncli system info                          │
│  2. troncli service list                         │
│                                                   │
└───────────────────────────────────────────────────┘

🚀 Executing commands...

[1/2] Executing: troncli system info
✅ Command completed successfully

[2/2] Executing: troncli service list
✅ Command completed successfully

🎉 All commands executed successfully!
```

### Teste 3: Root Agent - Comando com Confirmação
```bash
./troncli agent root "instalar nginx"
```

**Esperado:**
- Análise de risco: MEDIUM ou HIGH
- Solicitação de confirmação
- Comandos: `troncli pkg install nginx`

### Teste 4: Root Agent - Streaming Mode
```bash
TRONCLI_AGENT_STREAMING=true ./troncli agent root "listar processos"
```

**Esperado:**
- Resposta em tempo real (streaming)
- Mesma funcionalidade do modo normal

---

## 🔍 PASSO 7: Validação de Output

### Verificar Box-Drawing Characters
```bash
./troncli system info | cat -A
```
**Esperado:** Caracteres UTF-8 corretos: `┌─┐│└┘├┤`

### Verificar Cores
```bash
./troncli system info
```
**Esperado:**
- Bordas: Cyan brilhante
- Títulos: Bold + White
- Footers: Dim

### Verificar JSON Output
```bash
./troncli system info --json | jq .
```
**Esperado:** JSON válido e bem formatado

### Verificar YAML Output
```bash
./troncli system info --yaml
```
**Esperado:** YAML válido

---

## 🐛 PASSO 8: Testes de Erro

### Teste 1: Modelo não encontrado
```bash
./troncli agent root "teste" 2>&1 | grep "Model not found"
```
**Esperado:** Mensagem de erro clara com instruções

### Teste 2: llama-cli não encontrado
```bash
mv ~/.troncli/bin/llama-cli ~/.troncli/bin/llama-cli.bak
./troncli agent root "teste" 2>&1 | grep "binary not found"
mv ~/.troncli/bin/llama-cli.bak ~/.troncli/bin/llama-cli
```
**Esperado:** Mensagem de erro com instruções de instalação

### Teste 3: Comando inválido
```bash
./troncli invalid-command 2>&1
```
**Esperado:** Mensagem de erro + sugestões

---

## 📊 PASSO 9: Testes de Performance

### Benchmark Startup
```bash
time ./troncli --version
```
**Esperado:** < 100ms

### Benchmark System Info
```bash
time ./troncli system info
```
**Esperado:** < 200ms

### Benchmark Agent Response
```bash
time ./troncli agent root "listar serviços"
```
**Esperado:**
- CPU only: < 30 segundos
- GPU: < 5 segundos

---

## ✅ Checklist de Validação Final

### Compilação
- [ ] `go build` sem erros
- [ ] Binário `troncli` criado
- [ ] `./troncli --version` funciona

### llama.cpp
- [ ] `llama-cli` compilado
- [ ] AVX2 detectado (se disponível)
- [ ] Modelo GGUF baixado (~4GB)
- [ ] Teste manual funciona

### Comandos CLI
- [ ] `troncli system info` - box-drawing perfeito
- [ ] `troncli service list` - formatação profissional
- [ ] `troncli process tree` - árvore formatada
- [ ] `troncli network interfaces` - lista de interfaces
- [ ] `troncli disk usage` - uso de disco
- [ ] `troncli pkg list` - detecção de gerenciador

### TRON ROOT AGENT
- [ ] `troncli agent setup` - instalação automática
- [ ] `troncli agent status` - status do agente
- [ ] `troncli agent root "teste"` - execução básica
- [ ] Análise de risco funciona
- [ ] Confirmação para comandos perigosos
- [ ] Execução de comandos funciona
- [ ] Streaming mode funciona

### Output
- [ ] Box-drawing characters corretos (┌─┐│└┘├┤)
- [ ] Cores funcionam (cyan, bold, dim)
- [ ] Alinhamento perfeito
- [ ] JSON output válido
- [ ] YAML output válido

### Performance
- [ ] Startup < 100ms
- [ ] System info < 200ms
- [ ] Agent response < 30s (CPU)

---

## 🎉 Resultado Esperado

Se todos os testes passarem:
```
✅ TRONCLI está COMPLETAMENTE FUNCIONAL!
✅ TRON ROOT AGENT está operacional!
✅ Output profissional com box-drawing perfeito!
✅ Multi-distribuição Linux suportada!
✅ Pronto para merge com main!
```

---

## 🚨 Troubleshooting

### Problema: llama-cli não encontrado
**Solução:**
```bash
# Verificar paths
ls -la ~/.troncli/bin/llama-cli
ls -la /usr/local/bin/llama-cli
ls -la /usr/bin/llama-cli

# Adicionar ao PATH
export PATH="$HOME/.troncli/bin:$PATH"
```

### Problema: Modelo não carrega
**Solução:**
```bash
# Verificar tamanho do arquivo
ls -lh ~/.troncli/models/*.gguf

# Re-baixar se necessário
cd ~/.troncli/models/
rm -f *.gguf
wget https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_0.gguf
```

### Problema: Caracteres quebrados
**Solução:**
```bash
# Verificar locale
locale

# Configurar UTF-8
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8
```

### Problema: Cores não aparecem
**Solução:**
```bash
# Verificar suporte a cores
echo $TERM

# Forçar cores
export TERM=xterm-256color
```

---

## 📝 Relatório de Testes

Após executar todos os testes, preencha:

```
Data: _______________
Sistema: _______________
Distribuição: _______________
Kernel: _______________

Compilação: [ ] OK [ ] FALHOU
llama.cpp: [ ] OK [ ] FALHOU
Comandos CLI: [ ] OK [ ] FALHOU
TRON ROOT AGENT: [ ] OK [ ] FALHOU
Output: [ ] OK [ ] FALHOU
Performance: [ ] OK [ ] FALHOU

Observações:
_________________________________
_________________________________
_________________________________

Status Final: [ ] APROVADO [ ] REPROVADO
```
