# GitHub Token (gh CLI)

O token de acesso ao GitHub está sendo gerenciado pelo helper de credenciais do Git.

## Informações#

- O comando `gh` não está instalado no sistema
- Os pushes para o repositório funcionam via credential helper
- O token está armazenado no git credential helper

## Como obter o token manualmente#

```bash
# Se o gh estivesse instalado:
gh auth token

# Verificar credential helper:
git config --global credential.helper
```
