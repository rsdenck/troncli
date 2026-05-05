# GitHub Token (gh CLI)

Token de acesso ao GitHub (usado pelo gh CLI e git credentials).

## Token#

```
ghp_N369EtWTrL3saR15tLrx7UYf2Fjp4F3VKz
```

## Informações#

- Tipo: Fine-grained personal access token
- Permissões: Repo (read/write)
- Expiração: Configurada no GitHub

## Uso#

```bash
# Configurar git com token:
git config --global credential.helper store
echo "https://ghp_N369EtWTrL3saR15tLrx7UYf2Fjp4F3VKz@github.com" > ~/.git-credentials

# Verificar token:
gh auth token  # Requer gh CLI instalado
```
