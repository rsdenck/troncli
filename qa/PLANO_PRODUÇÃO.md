# PLANO PARA PRODUÇÃO - NUX

**Data:** 30 de Abril de 2026  
**Meta:** Score 9.4/10 (hoje: 8.5/10)  
**Repositório:** https://github.com/rsdenck/nux

---

## STATUS ATUAL (HOJE)

| Métrica | Valor | Meta | Status |
|---------|-------|------|--------|
| **Score Geral** | 8.5/10 | 9.4/10 | 🔄 Quase lá |
| **Cobertura Testes** | 35% | 70% | 🔄 35% restante |
| **Comandos Funcionais** | 15+ | 15+ | ✅ Completo |
| **Comandos Refatorados** | 2/15 | 15/15 | 🔄 13 restantes |
| **Segurança Aplicada** | 50% | 100% | 🔄 Aplicar em todos |
| **Testes E2E** | 0% | 70% | 🔄 Criar com mocks |
| **CI/CD** | 0% | 100% | 🔄 Criar pipeline |

---

## 🔥 PRÓXIMOS 5 DIAS (P0 Crítico)

### DIA 1-2: Refatorar TODOS os Comandos (P0)
**Objetivo:** Migrar TODOS os 15+ comandos para `core.Executor`

**Ja feito (2/15):**
- ✅ `disk.go` → `diskExecutor`
- ✅ `network.go` → `networkExecutor`

**Pendente (13/15):**
1. 🔄 `service.go` → `serviceExecutor`
2. 🔄 `pkg.go` → `pkgExecutor`
3. 🔄 `system.go` → `systemExecutor`
4. 🔄 `users.go` → `usersExecutor`
5. 🔄 `firewall.go` → `firewallExecutor`
6. 🔄 `container.go` → `containerExecutor`
7. 🔄 `bash.go` → `bashExecutor`
8. 🔄 `remote.go` → `remoteExecutor`
9. 🔄 `geoip.go` → `geoipExecutor`
10. 🔄 `process.go` → `processExecutor`
11. 🔄 `vault.go` → `vaultExecutor`
12. 🔄 `ask.go` → `askExecutor`
13. 🔄 `onboard.go` → `onboardExecutor`

**Template de Refatoração:**
```go
// ANTES (perigoso):
cmd := exec.Command("ls", "-la")
out, err := cmd.CombinedOutput()

// DEPOIS (seguro e testável):
var cmdExecutor core.Executor = &core.RealExecutor{}
out, err := cmdExecutor.CombinedOutput("ls", "-la")
```

---

### DIA 3: Aplicar Segurança em TODOS (P0)
**Objetivo:** Chamar `SanitizeInput()` e `ValidatePath()` em TODOS os comandos

```go
// Em CADA comando:
func Run(cmd *cobra.Command, args []string) {
    // Sanitize TODAS as entradas
    for i, arg := range args {
        args[i] = core.SanitizeInput(arg)
    }
    
    // Validate paths
    if len(args) > 0 && strings.Contains(args[0], "/") {
        if !core.ValidatePath(args[0]) {
            output.NewError("invalid path", "INVALID_PATH").Print()
            return
        }
    }
    
    // Use executor
    executor := &core.RealExecutor{}
    out, err := executor.CombinedOutput("command", args...)
}
```

---

### DIA 4-5: Testes E2E com MockExecutor (P0)
**Objetivo:** Atingir 70% cobertura

```go
func TestDiskListE2E(t *testing.T) {
    mock := &core.MockExecutor{
        Output: `{"blockdevices": [{"name":"sda","size":"100G"}]}`,
        Err:    nil,
    }
    
    // Testar fluxo completo
    out, err := mock.CombinedOutput("lsblk", "-J")
    if err != nil {
        t.Errorf("E2E failed: %v", err)
    }
    
    // Verificar resultado
    var result map[string]interface{}
    json.Unmarshal([]byte(out), &result)
    // Assertions...
}
```

**Testes prioritários:**
1. TestDiskListE2E
2. TestNetworkListE2E
3. TestServiceListE2E
4. TestPkgInstallE2E
5. TestUserAddE2E

---

## DIA 6: CI/CD com GitHub Actions (P2)

**Objetivo:** Automatizar testes e lint

```yaml
# .github/workflows/go.yml
name: Go
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests with coverage
        run: |
          go test ./... -coverprofile=coverage.out
          go tool cover -html=coverage.out -o coverage.html
      
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
      
      - name: Build
        run: go build -o nux ./cmd/nux
```

---

## DIA 7: Release Candidate v0.4.0-rc

**Objetivo:** Preparar lançamento público

1. ✅ TODOS os comandos refatorados
2. ✅ Segurança aplicada em todos
3. ✅ Testes E2E criados (70% coverage)
4. ✅ CI/CD funcionando
5. ✅ Documentação atualizada
6. ✅ Build para todas as distros (GoReleaser)

---

## 📊 MÉTRICAS ESPERADAS APÓS 7 DIAS

| Métrica | Atual | Após 7 dias | Melhoria |
|---------|-------|-------------|----------|
| **Score Geral** | 8.5/10 | **9.4/10** | +0.9 |
| **Cobertura Testes** | 35% | **70%** | +35% |
| **Comandos Refatorados** | 2/15 | **15/15** | +13 |
| **Segurança Aplicada** | 50% | **100%** | +50% |
| **Testes E2E** | 0% | **20%** | +20% |
| **CI/CD** | 0% | **100%** | +100% |

---

## 🎯 VEREDITO FINAL

### ✅ O que está BOM:
- Base técnica sólida (Go + Cobra)
- 172 skills documentadas
- 15+ comandos principais funcionais
- Vault profissional implementado
- Ask com múltiplos providers
- Onboard testado e funcionando
- Testes unitários passando (35% coverage)
- Build sempre limpo
- Executor interface criado
- Segurança inicial implementada

### ❌ O que AINDA ESTÁ RUIM:
- **13 comandos ainda usam `exec.Command` direto (P0)**
- **Segurança não aplicada em todos comandos (P0)**
- **Testes E2E inexistentes (P0)**
- Output não padronizado em todos comandos (P1)
- Sem CI/CD automatizado (P2)

---

## 🚀 RECOMENDAÇÃO FINAL

**O NUX TEM BASE EXCELENTE PARA SER UM PRODUTO PROFISSIONAL.**

**DEVE SER LANÇADO COMO BETA PÚBLICO APÓS:**
1. ✅ Refatorar TODOS os 15+ comandos para `core.Executor`
2. ✅ Aplicar segurança em TODOS os comandos
3. ✅ Criar testes E2E (70% cobertura)
4. ✅ Configurar CI/CD básico

**Score Estimado pós-correções:** **9.4/10**

**NÃO deve ser abandonado. Deve ser madurado.**

---

## 📋 COMANDOS DISPONÍVEIS (Atual)

```bash
nux                      # Root ✅
nux onboard              # Onboard wizard ✅
nux doctor               # System health ✅
nux vault                # Vault management ✅
  nux vault show
  nux vault set-key
  nux vault get-key
nux ask                  # AI multi-provider ✅
  nux ask query
  nux ask config
nux disk                 # Disk/LVM ✅ (refactored)
nux network              # Network management ✅ (refactored)
nux service              # Services 🔄 (next)
nux pkg                  # Packages 🔄
nux system               # System info 🔄
nux users                # User management 🔄
nux firewall             # Firewall 🔄
nux container            # Docker/Podman 🔄
nux bash                 # Bash execution 🔄
nux remote               # SSH 🔄
nux geoip                # Geolocation 🔄
```

---

**Relatório gerado em:** 30/04/2026 14:40 BRT  
**Próxima revisão:** Após refatoração completa + testes E2E  

---
*NUX - Linux CLI Master/Manager - Rumo à Produção!*
