# 🤖 TRON ROOT AGENT - Guia de Uso

## 🎯 Modo de Uso Simplificado

### Comando Único
```bash
troncli agent
```

**Isso é tudo que você precisa!**

---

## 🚀 Primeira Execução (Automática)

Na primeira vez que você executar `troncli agent`, o sistema faz TUDO automaticamente:

```bash
$ troncli agent

╔════════════════════════════════════════════════════════════╗
║                                                            ║
║  🚀 PRIMEIRA EXECUÇÃO - SETUP AUTOMÁTICO                   ║
║                                                            ║
║  Instalando llama.cpp + Modelo Qwen2.5-Coder-7B            ║
║  Tempo estimado: ~10 minutos                               ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝

📁 Creating directories...
  ✓ /home/user/.troncli
  ✓ /home/user/.troncli/bin
  ✓ /home/user/.troncli/models

📥 Cloning llama.cpp...
✓ AVX2 detected, using optimized build

🔨 Compiling llama.cpp (this may take a few minutes)...
✓ Binary installed to: /home/user/.troncli/bin/llama-cli

📥 Downloading Qwen2.5-Coder-7B model (~4GB)...
  This may take several minutes depending on your connection
✓ Model downloaded (3.8GB)

╔════════════════════════════════════════════════════════════╗
║                                                            ║
║  ✓ SETUP COMPLETE!                                         ║
║                                                            ║
║  TRON ROOT AGENT is ready!                                 ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
```

**Após o setup, o modo interativo inicia automaticamente!**

---

## 💬 Modo Interativo

Após o setup (ou em execuções subsequentes), você verá:

```bash
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║  TRON ROOT AGENT › MODO INTERATIVO                         ║
║                                                            ║
║  Digite comandos em linguagem natural                      ║
║  Digite 'sair' ou 'exit' para encerrar                     ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝

❯ _
```

### Como Usar

1. **Digite sua intenção em linguagem natural**
2. **Pressione ENTER**
3. **O agente analisa e executa**
4. **Veja o resultado em AZUL**

### Exemplos de Comandos

```bash
❯ verificar saúde do sistema

┌── TRON ROOT AGENT › MODO AUTÔNOMO ───────────────────────┐
│                                                           │
│  Intenção › verificar saúde do sistema                    │
│  Modelo   › Qwen2.5-Coder-7B                             │
│  Engine   › llama.cpp                                     │
│  Modo     › Hardcore Linux                               │
│                                                           │
└───────────────────────────────────────────────────────────┘

🧠 Analisando intenção...

┌── ANÁLISE DO AGENTE ─────────────────────────────────────┐
│                                                           │
│  Vou verificar o status do sistema usando comandos       │
│  troncli para coletar informações sobre saúde            │
│                                                           │
└───────────────────────────────────────────────────────────┘

┌── AVALIAÇÃO DE RISCO ────────────────────────────────────┐
│                                                           │
│  Nível de Risco  › LOW                                   │
│  Impacto         › Read-only system check                │
│  Confirmação     › false                                 │
│                                                           │
└───────────────────────────────────────────────────────────┘

┌── COMANDOS A EXECUTAR ───────────────────────────────────┐
│                                                           │
│  1. troncli system info                                  │
│  2. troncli service list                                 │
│  3. troncli disk usage                                   │
│                                                           │
└───────────────────────────────────────────────────────────┘

💡 Raciocínio: Estes comandos fornecem uma visão completa da
saúde do sistema incluindo informações básicas, serviços e
uso de disco.

🚀 Executando comandos...

[1/3] Executando: troncli system info
┌── TRONCLI: SYSTEM INFO ──────────────────────────────────┐
│                                                           │
│  Hostname    › myserver                                   │
│  OS          › Linux                                      │
│  Kernel      › 5.15.0-91-generic                          │
│  Uptime      › 5 days, 3 hours                            │
│                                                           │
└───────────────────────────────────────────────────────────┘
✅ Comando executado com sucesso

[2/3] Executando: troncli service list
...
✅ Comando executado com sucesso

[3/3] Executando: troncli disk usage
...
✅ Comando executado com sucesso

🎉 Todos os comandos executados com sucesso!

❯ _
```

---

## 📝 Mais Exemplos

### Exemplo 1: Instalar Software
```bash
❯ instalar nginx

┌── ANÁLISE DO AGENTE ─────────────────────────────────────┐
│                                                           │
│  Vou instalar o nginx usando o gerenciador de pacotes    │
│  do sistema                                               │
│                                                           │
└───────────────────────────────────────────────────────────┘

┌── AVALIAÇÃO DE RISCO ────────────────────────────────────┐
│                                                           │
│  Nível de Risco  › MEDIUM                                │
│  Impacto         › Instala novo software no sistema      │
│  Confirmação     › true                                  │
│                                                           │
└───────────────────────────────────────────────────────────┘

⚠️  CONFIRMATION REQUIRED
This operation has been classified as MEDIUM risk.

Do you want to proceed? (yes/no): yes

🚀 Executando comandos...
[1/1] Executando: troncli pkg install nginx
✅ Comando executado com sucesso

🎉 Todos os comandos executados com sucesso!
```

### Exemplo 2: Listar Serviços
```bash
❯ listar serviços ativos

┌── ANÁLISE DO AGENTE ─────────────────────────────────────┐
│                                                           │
│  Vou listar todos os serviços ativos do sistema          │
│                                                           │
└───────────────────────────────────────────────────────────┘

┌── COMANDOS A EXECUTAR ───────────────────────────────────┐
│                                                           │
│  1. troncli service list                                 │
│                                                           │
└───────────────────────────────────────────────────────────┘

🚀 Executando comandos...
...
```

### Exemplo 3: Verificar Uso de Disco
```bash
❯ mostrar uso de disco

┌── COMANDOS A EXECUTAR ───────────────────────────────────┐
│                                                           │
│  1. troncli disk usage                                   │
│                                                           │
└───────────────────────────────────────────────────────────┘
```

### Exemplo 4: Listar Processos
```bash
❯ listar processos

┌── COMANDOS A EXECUTAR ───────────────────────────────────┐
│                                                           │
│  1. troncli process tree                                 │
│                                                           │
└───────────────────────────────────────────────────────────┘
```

---

## 🎨 Cores e Formatação

### Todas as saídas do agente são em AZUL:
- 🔵 Análise do agente
- 🔵 Comandos a executar
- 🔵 Raciocínio
- 🔵 Mensagens de progresso
- 🔵 Confirmações de sucesso

### Cores de Risco:
- 🟢 **LOW** - Verde (operações seguras, read-only)
- 🟡 **MEDIUM** - Amarelo (instalações, modificações)
- 🔴 **HIGH** - Vermelho (operações críticas)
- 🔴 **CRITICAL** - Vermelho Bold (operações perigosas)

---

## 🚪 Sair do Modo Interativo

Para sair, digite qualquer um destes comandos:
```bash
❯ sair
❯ exit
❯ quit
```

Você verá:
```bash
👋 Até logo!
```

---

## 🎯 Modo Direto (Sem Interatividade)

Se preferir executar um comando único sem entrar no modo interativo:

```bash
troncli agent "verificar saúde do sistema"
```

O agente executa e retorna ao shell imediatamente.

---

## ⚙️ Configuração

### Localização dos Arquivos

Tudo é instalado em `~/.troncli/`:

```
~/.troncli/
├── bin/
│   └── llama-cli          # Binário do llama.cpp
├── models/
│   └── qwen2.5-coder-7b-instruct-q4_0.gguf  # Modelo (~4GB)
└── llama.cpp/             # Código fonte (opcional)
```

### Requisitos de Sistema

**Mínimo:**
- Linux Kernel 3.10+
- 4GB RAM
- 8GB espaço em disco
- CPU x86_64 ou aarch64

**Recomendado:**
- Linux Kernel 5.0+
- 8GB RAM
- CPU com AVX2
- 16GB espaço em disco

### Performance

**Tempo de Resposta (CPU):**
- Análise simples: ~5-10 segundos
- Análise complexa: ~15-30 segundos

**Tempo de Resposta (GPU):**
- Análise simples: ~1-2 segundos
- Análise complexa: ~3-5 segundos

---

## 🐛 Troubleshooting

### Problema: Setup falha no download
```bash
# Re-executar setup manualmente
rm -rf ~/.troncli
troncli agent
```

### Problema: Modelo não carrega
```bash
# Verificar tamanho do arquivo
ls -lh ~/.troncli/models/*.gguf

# Deve mostrar ~3.8GB
# Se menor, re-baixar:
cd ~/.troncli/models/
rm -f *.gguf
wget https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_0.gguf
```

### Problema: llama-cli não encontrado
```bash
# Verificar instalação
ls -la ~/.troncli/bin/llama-cli

# Se não existir, re-compilar:
cd ~/.troncli/llama.cpp
make LLAMA_NATIVE=1
cp llama-cli ~/.troncli/bin/
```

### Problema: Respostas lentas
```bash
# Verificar se AVX2 está habilitado
grep avx2 /proc/cpuinfo

# Se sim, re-compilar com otimização:
cd ~/.troncli/llama.cpp
make clean
make LLAMA_NATIVE=1
cp llama-cli ~/.troncli/bin/
```

---

## 💡 Dicas de Uso

### 1. Seja Específico
❌ "fazer algo"
✅ "verificar saúde do sistema"

### 2. Use Linguagem Natural
✅ "instalar nginx"
✅ "listar serviços ativos"
✅ "mostrar uso de disco"
✅ "verificar processos rodando"

### 3. Confirme Operações Perigosas
O agente sempre pede confirmação para operações de risco MEDIUM ou superior.

### 4. Use Modo Direto para Scripts
```bash
#!/bin/bash
troncli agent "verificar saúde do sistema" > health_report.txt
troncli agent "listar serviços" > services.txt
```

---

## 🎉 Resumo

**Um único comando:**
```bash
troncli agent
```

**Faz tudo:**
- ✅ Setup automático na primeira execução
- ✅ Download de llama.cpp
- ✅ Download do modelo (~4GB)
- ✅ Modo interativo com prompt simples
- ✅ Análise de risco automática
- ✅ Execução com confirmação
- ✅ Output em azul
- ✅ Zero configuração manual

**Pronto para usar!** 🚀
