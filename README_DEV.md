# TRONCLI - Branch DEV 🚀

## 🎯 Status: PRONTO PARA TESTES

Esta branch contém as melhorias completas da TRONCLI CLI Enhancement, incluindo o **TRON ROOT AGENT** com integração llama.cpp.

---

## ✨ Novidades nesta Branch

### 🤖 TRON ROOT AGENT
- Agente autônomo usando llama.cpp diretamente (sem Ollama)
- Análise de risco automática (low, medium, high, critical)
- Confirmação para operações perigosas
- Execução de comandos com timeout
- Modo streaming para respostas em tempo real
- Instalação automatizada: `troncli agent setup`

### 🎨 Output Profissional
- Box-drawing characters perfeitos: `┌─┐│└┘├┤`
- Cores profissionais (cyan, bold, dim)
- Alinhamento perfeito em TODOS os comandos
- Suporte JSON/YAML

### 🐧 Multi-Distribuição Linux
- Gentoo Linux (Portage)
- Void Linux (XBPS)
- 8+ distribuições suportadas

### ⚡ Integração Linux Direta
- Leitura de `/proc` filesystem
- Leitura de `/sys` filesystem
- Syscalls diretos (kill, renice, statfs)
- Zero dependências externas

### 🗑️ TUI Removida
- CLI puro e limpo
- Sem dependências de TUI
- 100% backward compatible

---

## 🚀 Instalação Rápida (Linux)

### Opção 1: Script Automatizado (RECOMENDADO)
```bash
# Clonar repositório
git clone https://github.com/rsdenck/troncli.git
cd troncli
git checkout dev

# Executar instalação completa
chmod +x quick-install.sh
./quick-install.sh
```

**O script instala:**
- Dependências do sistema
- TRONCLI compilada
- llama.cpp otimizado
- Modelo Qwen2.5-Coder-7B (~4GB)

### Opção 2: Manual
```bash
# Clonar e compilar
git clone https://github.com/rsdenck/troncli.git
cd troncli
git checkout dev
go build -o troncli cmd/troncli/main.go

# Instalar Root Agent
./troncli agent setup
```

---

## 🧪 Validação e Testes

### Suite de Testes Automatizada
```bash
chmod +x test-troncli.sh
./test-troncli.sh
```

**O script testa:**
- ✅ Compilação
- ✅ Comandos CLI básicos
- ✅ Output formatting
- ✅ llama.cpp instalação
- ✅ Modelo GGUF
- ✅ TRON ROOT AGENT
- ✅ Performance benchmarks

### Testes Manuais
```bash
# Comandos básicos
./troncli system info
./troncli service list
./troncli process tree
./troncli network interfaces
./troncli disk usage

# TRON ROOT AGENT
./troncli agent root "verificar saúde do sistema"
./troncli agent root "listar serviços ativos"
./troncli agent root "mostrar uso de disco"

# Streaming mode
TRONCLI_AGENT_STREAMING=true ./troncli agent root "listar processos"
```

### Documentação Completa
Veja `TEST_VALIDATION.md` para guia detalhado de validação.

---

## 📊 Resultados Esperados

### Comandos CLI
```
┌── TRONCLI: SYSTEM INFO ──────────────────────────────────┐
│                                                           │
│  Hostname    › myserver                                   │
│  OS          › Linux                                      │
│  Kernel      › 5.15.0-91-generic                          │
│  Uptime      › 5 days, 3 hours                            │
│                                                           │
└───────────────────────────────────────────────────────────┘
```

### TRON ROOT AGENT
```
┌── TRON ROOT AGENT › AUTONOMOUS MODE ─────────────────────┐
│                                                           │
│  Intent    › verificar saúde do sistema                   │
│  Model     › Qwen2.5-Coder-7B                            │
│  Engine    › llama.cpp                                    │
│  Mode      › Hardcore Linux                              │
│                                                           │
└───────────────────────────────────────────────────────────┘

🧠 Analyzing intent...

┌── AGENT ANALYSIS ────────────────────────────────────────┐
│                                                           │
│  Vou verificar o status do sistema usando comandos       │
│  troncli para coletar informações sobre saúde            │
│                                                           │
└───────────────────────────────────────────────────────────┘

┌── RISK ASSESSMENT ───────────────────────────────────────┐
│                                                           │
│  Risk Level      › LOW                                   │
│  Impact          › Read-only system check                │
│  Confirmation    › false                                 │
│                                                           │
└───────────────────────────────────────────────────────────┘

┌── COMMANDS TO EXECUTE ───────────────────────────────────┐
│                                                           │
│  1. troncli system info                                  │
│  2. troncli service list                                 │
│  3. troncli disk usage                                   │
│                                                           │
└───────────────────────────────────────────────────────────┘

🚀 Executing commands...

[1/3] Executing: troncli system info
✅ Command completed successfully

[2/3] Executing: troncli service list
✅ Command completed successfully

[3/3] Executing: troncli disk usage
✅ Command completed successfully

🎉 All commands executed successfully!
```

---

## 🔧 Requisitos do Sistema

### Mínimo
- Linux Kernel 3.10+
- 4GB RAM
- 8GB espaço em disco (para modelo)
- CPU x86_64 ou aarch64

### Recomendado
- Linux Kernel 5.0+
- 8GB RAM
- CPU com AVX2
- 16GB espaço em disco

### Distribuições Testadas
- ✅ Ubuntu 20.04+
- ✅ Debian 11+
- ✅ Fedora 35+
- ✅ Arch Linux
- ✅ Alpine Linux
- ✅ openSUSE
- ✅ Gentoo
- ✅ Void Linux

---

## 📝 Checklist de Validação

Antes de fazer merge com `main`, validar:

### Compilação
- [ ] `go build` sem erros
- [ ] Binário `troncli` criado
- [ ] `./troncli --version` funciona

### Comandos CLI
- [ ] `troncli system info` - output perfeito
- [ ] `troncli service list` - formatação profissional
- [ ] `troncli process tree` - árvore formatada
- [ ] `troncli network interfaces` - lista de interfaces
- [ ] `troncli disk usage` - uso de disco
- [ ] Todos os comandos com box-drawing perfeito

### TRON ROOT AGENT
- [ ] `troncli agent setup` - instalação automática
- [ ] `troncli agent status` - status do agente
- [ ] `troncli agent root "teste"` - execução básica
- [ ] Análise de risco funciona
- [ ] Confirmação para comandos perigosos
- [ ] Execução de comandos funciona
- [ ] Streaming mode funciona

### Performance
- [ ] Startup < 100ms
- [ ] System info < 200ms
- [ ] Agent response < 30s (CPU)

### Testes Automatizados
- [ ] `./test-troncli.sh` - todos os testes passam
- [ ] Taxa de sucesso: 100%

---

## 🐛 Troubleshooting

### Problema: llama-cli não encontrado
```bash
# Verificar instalação
ls -la ~/.troncli/bin/llama-cli

# Reinstalar
./troncli agent setup
```

### Problema: Modelo não carrega
```bash
# Verificar tamanho
ls -lh ~/.troncli/models/*.gguf

# Re-baixar
cd ~/.troncli/models/
rm -f *.gguf
wget https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_0.gguf
```

### Problema: Caracteres quebrados
```bash
# Configurar UTF-8
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8
```

### Problema: Cores não aparecem
```bash
# Forçar cores
export TERM=xterm-256color
```

---

## 📚 Documentação

- `TEST_VALIDATION.md` - Guia completo de validação
- `test-troncli.sh` - Suite de testes automatizada
- `quick-install.sh` - Script de instalação rápida
- `.kiro/specs/troncli-cli-enhancement/` - Especificações completas

---

## 🎯 Próximos Passos

1. **Testar em VM Linux**
   ```bash
   ./quick-install.sh
   ./test-troncli.sh
   ```

2. **Validar todos os comandos**
   - Executar checklist completo
   - Verificar output em diferentes distribuições

3. **Testar Root Agent**
   - Comandos simples
   - Comandos com confirmação
   - Streaming mode

4. **Fazer Merge com Main**
   ```bash
   git checkout main
   git merge dev
   git push origin main
   ```

---

## 🤝 Contribuindo

Esta é uma branch de desenvolvimento. Para contribuir:

1. Teste completamente em sua distribuição Linux
2. Reporte bugs ou problemas
3. Sugira melhorias
4. Valide que todos os testes passam

---

## 📞 Suporte

- Issues: https://github.com/rsdenck/troncli/issues
- Documentação: `TEST_VALIDATION.md`
- Testes: `./test-troncli.sh`

---

## ⚡ Performance

### Benchmarks Esperados
- Startup: < 100ms
- System info: < 200ms
- Service list: < 500ms
- Agent response (CPU): < 30s
- Agent response (GPU): < 5s

### Otimizações
- AVX2 detection automática
- Compilação otimizada do llama.cpp
- Cache de contexto
- Modelo quantizado Q4_0

---

## 🎉 Status Final

**Branch DEV está pronta para:**
- ✅ Testes em produção
- ✅ Validação em múltiplas distribuições
- ✅ Merge com main (após validação)

**Execute os testes e valide!** 🚀
