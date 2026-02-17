# TRONCLI

> **IDENTIDADE DO SISTEMA**
> 
> HOST: GITHUB
> OS: LINUX
> KERNEL: 6.8.0-RC1
> UPTIME: 99.99%

## VISÃO GERAL

TRONCLI é uma TUI (Interface de Usuário em Texto) de Administração de Sistemas Linux de nível de produção, escrita 100% em Go. Ela fornece monitoramento de sistema em tempo real, gerenciamento de LVM, análise de rede e auditoria de segurança com uma identidade visual inspirada em sistemas de grade de alta tecnologia.

**SEM MOCKS. SEM SIMULAÇÃO. CONTROLE REAL DO KERNEL.**

![Licença](https://img.shields.io/badge/LICENSE-MIT-00d9ff?style=for-the-badge&labelColor=000000)
![Versão Go](https://img.shields.io/badge/GO-1.22+-00d9ff?style=for-the-badge&labelColor=000000&logo=go)
![Plataforma](https://img.shields.io/badge/PLATFORM-LINUX-00d9ff?style=for-the-badge&labelColor=000000&logo=linux)
![Build](https://img.shields.io/badge/BUILD-PASSING-00d9ff?style=for-the-badge&labelColor=000000)

## MÓDULOS PRINCIPAIS

### [01] DASHBOARD DO SISTEMA
Métricas em tempo real de `/proc` e `/sys`.
- Uso de CPU (Usuário/Sistema/Ocioso)
- Memória (Ram/Swap)
- Carga Média (1/5/15)
- E/S de Disco e Taxa de Transferência de Rede

### [02] GERENCIADOR LVM
Interface direta para o Logical Volume Manager.
- Volumes Físicos (PV)
- Grupos de Volumes (VG)
- Volumes Lógicos (LV)
- Operações de Extensão/Redução

### [03] MATRIZ DE REDE
Análise avançada da pilha de rede.
- Estatísticas de interfaces
- Taxas de RX/TX em tempo real
- Configuração de DNS
- Estados de soquetes

### [04] AUDITORIA DE SEGURANÇA
Endurecimento do sistema e gerenciamento de usuários.
- Enumeração de Usuários/Grupos
- Gerenciamento de sessões SSH
- Matriz de permissões de arquivos
- Logs de auditoria

## INSTALAÇÃO

### PRÉ-REQUISITOS
- Kernel Linux 5.4+
- Privilégios de Root (para LVM/Auditoria)

### COMPILAR DO CÓDIGO FONTE

```bash
git clone https://github.com/rsdenck/troncli.git
cd troncli
go build -ldflags="-s -w" -o troncli cmd/troncli/main.go
./troncli
```

## ARQUITETURA

O sistema segue os princípios da Clean Architecture com estrita separação de responsabilidades.

```text
cmd/
  troncli/       # Ponto de entrada
internal/
  core/          # Lógica de domínio e Portas
  modules/       # Implementações (Específicas Linux)
  ui/            # Camada TUI (tview/tcell)
```

## CONTRIBUIÇÃO

1. Faça um Fork do repositório
2. Crie sua branch de feature (`git checkout -b feature/RecursoIncrivel`)
3. Commit suas alterações (`git commit -m 'feat: Adiciona RecursoIncrivel'`)
4. Push para a branch (`git push origin feature/RecursoIncrivel`)
5. Abra um Pull Request

## SEGURANÇA

Por favor, reporte vulnerabilidades para `security@troncli.local`.
Veja [SECURITY.md](SECURITY.md) para detalhes.

---
TRONCLI | FIM DE LINHA DO SISTEMA
