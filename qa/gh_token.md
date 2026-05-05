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
- Usuário: rsdenck
- Repositório: nux

## Uso#

```bash
# Configurar git com token:
git config --global credential.helper store
echo "https://ghp_N369EtWTrL3saR15tLrx7UYf2Fjp4F3VKz@github.com" > ~/.git-credentials

# Verificar token (requer gh CLI instalado):
gh auth token
```

## Configuração Atual#

```
Remote: https://github.com/rsdenck/nux.git
Credential Helper: store
```
