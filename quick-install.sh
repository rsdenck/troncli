#!/bin/bash
# TRONCLI - Instalação Rápida e Completa
# Instala TRONCLI + llama.cpp + Modelo GGUF

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

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
║              INSTALAÇÃO RÁPIDA                             ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
EOF
echo -e "${RESET}\n"

# Verificar se está rodando como root
if [ "$EUID" -eq 0 ]; then
    echo -e "${RED}❌ Não execute este script como root!${RESET}"
    echo -e "${YELLOW}Execute como usuário normal: ./quick-install.sh${RESET}"
    exit 1
fi

# Detectar distribuição
if [ -f /etc/os-release ]; then
    . /etc/os-release
    DISTRO=$ID
else
    DISTRO="unknown"
fi

echo -e "${CYAN}📋 Sistema detectado: $PRETTY_NAME${RESET}\n"

# ============================================================================
# PASSO 1: Instalar dependências
# ============================================================================
echo -e "${CYAN}${BOLD}[1/5] Instalando dependências...${RESET}\n"

case $DISTRO in
    ubuntu|debian)
        echo -e "${YELLOW}Instalando: git, build-essential, wget, golang${RESET}"
        sudo apt update
        sudo apt install -y git build-essential wget golang
        ;;
    fedora|rhel|centos)
        echo -e "${YELLOW}Instalando: git, gcc, make, wget, golang${RESET}"
        sudo dnf install -y git gcc gcc-c++ make wget golang
        ;;
    arch|manjaro)
        echo -e "${YELLOW}Instalando: git, base-devel, wget, go${RESET}"
        sudo pacman -S --noconfirm git base-devel wget go
        ;;
    alpine)
        echo -e "${YELLOW}Instalando: git, build-base, wget, go${RESET}"
        sudo apk add git build-base wget go
        ;;
    gentoo)
        echo -e "${YELLOW}Instalando: git, gcc, make, wget, go${RESET}"
        sudo emerge --ask dev-vcs/git sys-devel/gcc sys-devel/make net-misc/wget dev-lang/go
        ;;
    void)
        echo -e "${YELLOW}Instalando: git, gcc, make, wget, go${RESET}"
        sudo xbps-install -S git gcc make wget go
        ;;
    *)
        echo -e "${YELLOW}⚠️  Distribuição não reconhecida. Certifique-se de ter:${RESET}"
        echo -e "   - git, gcc, make, wget, go"
        read -p "Pressione ENTER para continuar..."
        ;;
esac

echo -e "${GREEN}✅ Dependências instaladas${RESET}\n"

# ============================================================================
# PASSO 2: Clonar e compilar TRONCLI
# ============================================================================
echo -e "${CYAN}${BOLD}[2/5] Clonando e compilando TRONCLI...${RESET}\n"

if [ -d "troncli" ]; then
    echo -e "${YELLOW}⚠️  Diretório troncli já existe. Atualizando...${RESET}"
    cd troncli
    git fetch origin
    git checkout dev
    git pull origin dev
else
    git clone https://github.com/rsdenck/troncli.git
    cd troncli
    git checkout dev
fi

echo -e "${YELLOW}Compilando...${RESET}"
go build -o troncli cmd/troncli/main.go

if [ -f "troncli" ]; then
    echo -e "${GREEN}✅ TRONCLI compilada com sucesso${RESET}\n"
else
    echo -e "${RED}❌ Falha na compilação${RESET}"
    exit 1
fi

# ============================================================================
# PASSO 3: Instalar llama.cpp
# ============================================================================
echo -e "${CYAN}${BOLD}[3/5] Instalando llama.cpp...${RESET}\n"

mkdir -p ~/.troncli/{bin,models}

if [ -f ~/.troncli/bin/llama-cli ]; then
    echo -e "${GREEN}✅ llama-cli já instalado${RESET}\n"
else
    echo -e "${YELLOW}Clonando llama.cpp...${RESET}"
    git clone https://github.com/ggerganov/llama.cpp ~/.troncli/llama.cpp
    
    cd ~/.troncli/llama.cpp
    
    # Verificar AVX2
    if grep -q avx2 /proc/cpuinfo; then
        echo -e "${GREEN}✅ AVX2 detectado - compilando otimizado${RESET}"
        make LLAMA_NATIVE=1
    else
        echo -e "${YELLOW}⚠️  AVX2 não detectado - compilando padrão${RESET}"
        make
    fi
    
    # Copiar binário
    if [ -f "llama-cli" ]; then
        cp llama-cli ~/.troncli/bin/
    elif [ -f "main" ]; then
        cp main ~/.troncli/bin/llama-cli
    else
        echo -e "${RED}❌ Binário llama-cli não encontrado${RESET}"
        exit 1
    fi
    
    chmod +x ~/.troncli/bin/llama-cli
    
    cd -
    
    echo -e "${GREEN}✅ llama.cpp instalado${RESET}\n"
fi

# ============================================================================
# PASSO 4: Baixar modelo GGUF
# ============================================================================
echo -e "${CYAN}${BOLD}[4/5] Baixando modelo Qwen2.5-Coder-7B (~4GB)...${RESET}\n"

MODEL_PATH=~/.troncli/models/qwen2.5-coder-7b-instruct-q4_0.gguf

if [ -f "$MODEL_PATH" ]; then
    echo -e "${GREEN}✅ Modelo já baixado${RESET}\n"
else
    echo -e "${YELLOW}Baixando... (isso pode levar alguns minutos)${RESET}"
    
    cd ~/.troncli/models/
    
    if command -v wget &> /dev/null; then
        wget -O qwen2.5-coder-7b-instruct-q4_0.gguf \
            https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_0.gguf
    elif command -v curl &> /dev/null; then
        curl -L -o qwen2.5-coder-7b-instruct-q4_0.gguf \
            https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_0.gguf
    else
        echo -e "${RED}❌ wget ou curl não encontrado${RESET}"
        exit 1
    fi
    
    cd -
    
    if [ -f "$MODEL_PATH" ]; then
        MODEL_SIZE=$(du -h "$MODEL_PATH" | cut -f1)
        echo -e "${GREEN}✅ Modelo baixado ($MODEL_SIZE)${RESET}\n"
    else
        echo -e "${RED}❌ Falha no download do modelo${RESET}"
        exit 1
    fi
fi

# ============================================================================
# PASSO 5: Teste final
# ============================================================================
echo -e "${CYAN}${BOLD}[5/5] Testando instalação...${RESET}\n"

echo -e "${YELLOW}Teste 1: TRONCLI version${RESET}"
./troncli --version
echo -e "${GREEN}✅ OK${RESET}\n"

echo -e "${YELLOW}Teste 2: TRONCLI system info${RESET}"
./troncli system info
echo -e "${GREEN}✅ OK${RESET}\n"

echo -e "${YELLOW}Teste 3: llama-cli version${RESET}"
~/.troncli/bin/llama-cli --version
echo -e "${GREEN}✅ OK${RESET}\n"

# ============================================================================
# INSTALAÇÃO COMPLETA
# ============================================================================
echo -e "${GREEN}${BOLD}"
cat << "EOF"
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║  ✅ INSTALAÇÃO COMPLETA!                                   ║
║                                                            ║
║  TRONCLI está pronta para uso!                             ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
EOF
echo -e "${RESET}\n"

echo -e "${CYAN}${BOLD}📝 Próximos passos:${RESET}\n"
echo -e "1. Adicionar ao PATH (opcional):"
echo -e "   ${YELLOW}export PATH=\"\$PWD:\$HOME/.troncli/bin:\$PATH\"${RESET}\n"
echo -e "2. Testar comandos básicos:"
echo -e "   ${YELLOW}./troncli system info${RESET}"
echo -e "   ${YELLOW}./troncli service list${RESET}"
echo -e "   ${YELLOW}./troncli process tree${RESET}\n"
echo -e "3. Testar TRON ROOT AGENT:"
echo -e "   ${YELLOW}./troncli agent root \"verificar saúde do sistema\"${RESET}"
echo -e "   ${YELLOW}./troncli agent root \"listar serviços ativos\"${RESET}\n"
echo -e "4. Executar suite de testes:"
echo -e "   ${YELLOW}chmod +x test-troncli.sh${RESET}"
echo -e "   ${YELLOW}./test-troncli.sh${RESET}\n"

echo -e "${CYAN}📚 Documentação completa: TEST_VALIDATION.md${RESET}\n"
