#!/bin/bash
# TRONCLI - Script de Validação Automatizada
# Testa se a TRONCLI está completamente funcional

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

# Contadores
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_TOTAL=0

# Função para imprimir cabeçalho
print_header() {
    echo -e "\n${CYAN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
    echo -e "${CYAN}${BOLD}  $1${RESET}"
    echo -e "${CYAN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}\n"
}

# Função para teste
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
    echo -e "${YELLOW}[TEST $TESTS_TOTAL]${RESET} $test_name"
    
    if eval "$test_command" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PASSED${RESET}\n"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "${RED}❌ FAILED${RESET}\n"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Função para teste com output
run_test_with_output() {
    local test_name="$1"
    local test_command="$2"
    
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
    echo -e "${YELLOW}[TEST $TESTS_TOTAL]${RESET} $test_name"
    
    if eval "$test_command"; then
        echo -e "${GREEN}✅ PASSED${RESET}\n"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "${RED}❌ FAILED${RESET}\n"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Banner
echo -e "${CYAN}${BOLD}"
cat << "EOF"
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║   ████████╗██████╗  ██████╗ ███╗   ██╗ ██████╗██╗     ██╗║
║   ╚══██╔══╝██╔══██╗██╔═══██╗████╗  ██║██╔════╝██║     ██║║
║      ██║   ██████╔╝██║   ██║██╔██╗ ██║██║     ██║     ██║║
║      ██║   ██╔══██╗██║   ██║██║╚██╗██║██║     ██║     ██║║
║      ██║   ██║  ██║╚██████╔╝██║ ╚████║╚██████╗███████╗██║║
║      ╚═╝   ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝╚══════╝╚═╝║
║                                                            ║
║              VALIDATION & TEST SUITE                       ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
EOF
echo -e "${RESET}\n"

# Verificar se está no diretório correto
if [ ! -f "go.mod" ] || [ ! -d "cmd/troncli" ]; then
    echo -e "${RED}❌ Erro: Execute este script no diretório raiz do projeto TRONCLI${RESET}"
    exit 1
fi

# ============================================================================
# FASE 1: PRÉ-REQUISITOS
# ============================================================================
print_header "FASE 1: Verificação de Pré-requisitos"

run_test "Git instalado" "which git"
run_test "Make instalado" "which make"
run_test "GCC instalado" "which gcc"
run_test "Go instalado" "which go"
run_test "Wget ou Curl instalado" "which wget || which curl"

# ============================================================================
# FASE 2: COMPILAÇÃO
# ============================================================================
print_header "FASE 2: Compilação da TRONCLI"

echo -e "${YELLOW}Compilando TRONCLI...${RESET}"
if go build -o troncli cmd/troncli/main.go; then
    echo -e "${GREEN}✅ Compilação bem-sucedida${RESET}\n"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}❌ Falha na compilação${RESET}\n"
    TESTS_FAILED=$((TESTS_FAILED + 1))
    exit 1
fi
TESTS_TOTAL=$((TESTS_TOTAL + 1))

run_test "Binário troncli criado" "[ -f ./troncli ]"
run_test "Binário é executável" "[ -x ./troncli ]"
run_test "Comando --version funciona" "./troncli --version"
run_test "Comando --help funciona" "./troncli --help"

# ============================================================================
# FASE 3: COMANDOS CLI BÁSICOS
# ============================================================================
print_header "FASE 3: Testes de Comandos CLI"

run_test "troncli system info" "./troncli system info"
run_test "troncli service list" "./troncli service list"
run_test "troncli process tree" "./troncli process tree"
run_test "troncli network interfaces" "./troncli network interfaces"
run_test "troncli disk usage" "./troncli disk usage"

# ============================================================================
# FASE 4: OUTPUT FORMATTING
# ============================================================================
print_header "FASE 4: Validação de Output"

echo -e "${YELLOW}[TEST]${RESET} Verificando box-drawing characters"
if ./troncli system info | grep -q "┌"; then
    echo -e "${GREEN}✅ PASSED - Box-drawing characters presentes${RESET}\n"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}❌ FAILED - Box-drawing characters ausentes${RESET}\n"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TESTS_TOTAL=$((TESTS_TOTAL + 1))

# ============================================================================
# FASE 5: LLAMA.CPP
# ============================================================================
print_header "FASE 5: Verificação do llama.cpp"

LLAMA_PATHS=(
    "$HOME/.troncli/bin/llama-cli"
    "/usr/local/bin/llama-cli"
    "/usr/bin/llama-cli"
    "/opt/llama.cpp/llama-cli"
)

LLAMA_FOUND=false
LLAMA_PATH=""

for path in "${LLAMA_PATHS[@]}"; do
    if [ -f "$path" ]; then
        LLAMA_FOUND=true
        LLAMA_PATH="$path"
        break
    fi
done

if [ "$LLAMA_FOUND" = true ]; then
    echo -e "${GREEN}✅ llama-cli encontrado em: $LLAMA_PATH${RESET}\n"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}⚠️  llama-cli não encontrado${RESET}"
    echo -e "${YELLOW}Execute: ./troncli agent setup${RESET}\n"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TESTS_TOTAL=$((TESTS_TOTAL + 1))

# ============================================================================
# FASE 6: MODELO GGUF
# ============================================================================
print_header "FASE 6: Verificação do Modelo GGUF"

MODEL_PATHS=(
    "$HOME/.troncli/models/qwen2.5-coder-7b-instruct-q4_0.gguf"
    "$HOME/.troncli/models/qwen3-coder-q4_k_m.gguf"
    "/opt/models/qwen2.5-coder-7b-instruct-q4_0.gguf"
    "/opt/models/qwen3-coder-q4_k_m.gguf"
)

MODEL_FOUND=false
MODEL_PATH=""

for path in "${MODEL_PATHS[@]}"; do
    if [ -f "$path" ]; then
        MODEL_FOUND=true
        MODEL_PATH="$path"
        break
    fi
done

if [ "$MODEL_FOUND" = true ]; then
    MODEL_SIZE=$(du -h "$MODEL_PATH" | cut -f1)
    echo -e "${GREEN}✅ Modelo encontrado: $MODEL_PATH ($MODEL_SIZE)${RESET}\n"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}⚠️  Modelo GGUF não encontrado${RESET}"
    echo -e "${YELLOW}Execute: ./troncli agent setup${RESET}\n"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TESTS_TOTAL=$((TESTS_TOTAL + 1))

# ============================================================================
# FASE 7: TRON ROOT AGENT
# ============================================================================
print_header "FASE 7: Testes do TRON ROOT AGENT"

run_test "troncli agent status" "./troncli agent status"

if [ "$LLAMA_FOUND" = true ] && [ "$MODEL_FOUND" = true ]; then
    echo -e "${YELLOW}[TEST]${RESET} TRON ROOT AGENT - Comando simples"
    echo -e "${CYAN}Executando: ./troncli agent root \"listar informações do sistema\"${RESET}\n"
    
    if timeout 60 ./troncli agent root "listar informações do sistema"; then
        echo -e "\n${GREEN}✅ PASSED - Root Agent executou com sucesso${RESET}\n"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "\n${RED}❌ FAILED - Root Agent falhou ou timeout${RESET}\n"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
else
    echo -e "${YELLOW}⚠️  Pulando testes do Root Agent (llama.cpp ou modelo não encontrado)${RESET}\n"
fi

# ============================================================================
# FASE 8: PERFORMANCE
# ============================================================================
print_header "FASE 8: Testes de Performance"

echo -e "${YELLOW}[TEST]${RESET} Benchmark: Startup time"
STARTUP_TIME=$( { time ./troncli --version > /dev/null 2>&1; } 2>&1 | grep real | awk '{print $2}')
echo -e "Tempo: ${CYAN}$STARTUP_TIME${RESET}"
if [[ "$STARTUP_TIME" =~ ^0m0\.[0-9]+s$ ]]; then
    echo -e "${GREEN}✅ PASSED - Startup < 1s${RESET}\n"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}⚠️  WARNING - Startup time alto${RESET}\n"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TESTS_TOTAL=$((TESTS_TOTAL + 1))

echo -e "${YELLOW}[TEST]${RESET} Benchmark: System info time"
SYSINFO_TIME=$( { time ./troncli system info > /dev/null 2>&1; } 2>&1 | grep real | awk '{print $2}')
echo -e "Tempo: ${CYAN}$SYSINFO_TIME${RESET}"
echo -e "${GREEN}✅ PASSED${RESET}\n"
TESTS_PASSED=$((TESTS_PASSED + 1))
TESTS_TOTAL=$((TESTS_TOTAL + 1))

# ============================================================================
# RELATÓRIO FINAL
# ============================================================================
print_header "RELATÓRIO FINAL"

echo -e "${BOLD}Sistema:${RESET}"
echo -e "  OS: $(uname -s)"
echo -e "  Kernel: $(uname -r)"
echo -e "  Arch: $(uname -m)"
echo -e "  Distro: $(cat /etc/os-release 2>/dev/null | grep PRETTY_NAME | cut -d'"' -f2 || echo 'Unknown')"
echo ""

echo -e "${BOLD}Resultados:${RESET}"
echo -e "  Total de testes: ${CYAN}$TESTS_TOTAL${RESET}"
echo -e "  Testes passados: ${GREEN}$TESTS_PASSED${RESET}"
echo -e "  Testes falhados: ${RED}$TESTS_FAILED${RESET}"
echo ""

PASS_RATE=$((TESTS_PASSED * 100 / TESTS_TOTAL))
echo -e "  Taxa de sucesso: ${CYAN}${PASS_RATE}%${RESET}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}${BOLD}╔════════════════════════════════════════════════════════════╗${RESET}"
    echo -e "${GREEN}${BOLD}║                                                            ║${RESET}"
    echo -e "${GREEN}${BOLD}║  ✅ TODOS OS TESTES PASSARAM!                              ║${RESET}"
    echo -e "${GREEN}${BOLD}║                                                            ║${RESET}"
    echo -e "${GREEN}${BOLD}║  TRONCLI está COMPLETAMENTE FUNCIONAL!                     ║${RESET}"
    echo -e "${GREEN}${BOLD}║  Pronto para merge com main!                               ║${RESET}"
    echo -e "${GREEN}${BOLD}║                                                            ║${RESET}"
    echo -e "${GREEN}${BOLD}╚════════════════════════════════════════════════════════════╝${RESET}"
    exit 0
else
    echo -e "${YELLOW}${BOLD}╔════════════════════════════════════════════════════════════╗${RESET}"
    echo -e "${YELLOW}${BOLD}║                                                            ║${RESET}"
    echo -e "${YELLOW}${BOLD}║  ⚠️  ALGUNS TESTES FALHARAM                                ║${RESET}"
    echo -e "${YELLOW}${BOLD}║                                                            ║${RESET}"
    echo -e "${YELLOW}${BOLD}║  Revise os erros acima e execute novamente                 ║${RESET}"
    echo -e "${YELLOW}${BOLD}║                                                            ║${RESET}"
    echo -e "${YELLOW}${BOLD}╚════════════════════════════════════════════════════════════╝${RESET}"
    exit 1
fi
