# Padrão de Saída NUX (seguindo GCX CLI)

## Princípios
1. **JSON por padrão** para parsing automático
2. **Saída limpa** sem decorações desnecessárias
3. **Consistência** em todos os comandos

## Formato de Saída

### Sucesso
```json
{
  "status": "success",
  "data": { ... },
  "message": "optional message"
}
```

### Erro
```json
{
  "status": "error",
  "error": "error message",
  "code": "ERR_CODE"
}
```

### Listagem
```json
{
  "status": "success",
  "total": 10,
  "items": [ ... ]
}
```

## Flags Globais
- `--json`: Força saída JSON
- `--yaml`: Força saída YAML
- `--quiet`: Suprime saída (apenas código de retorno)
