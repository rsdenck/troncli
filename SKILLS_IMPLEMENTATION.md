# Implementação do Conceito NUX CLI Master

## Conceito
NUX é uma CLI Master/Manager que pode incorporar funcionalidades de outras CLIs através do sistema de Skills.

## Fluxo de uma Skill
1. **Skill Definition** (skills/nome.md)
   - Define repositório, comandos, tipo
   
2. **Skill Installation** (`nux skill install discord`)
   - Baixa/instala a CLI externa
   - Registra no vault (.skills/.nux.json)
   
3. **Skill Enable** (`nux skill enable discord`)
   - Ativa a skill para uso
   - Pode solicitar credenciais (API keys) armazenadas no vault
   
4. **Skill Usage** (`nux discord --help`)
   - NUX atua como proxy para a CLI externa
   - Ou executa comandos específicos da skill

## Implementação Técnica

### Estrutura do Vault (.skills/.nux.json)
```json
{
  "version": "1.0.0",
  "installed_skills": ["discord", "docker", "kubectl"],
  "enabled_skills": ["docker", "kubectl"],
  "api_keys": {
    "discord": "encrypted_or_token",
    "github": "ghp_xxxxx"
  },
  "ollama": {
    "host": "http://localhost:11434",
    "model": "qwen3-coder",
    "enabled": true
  },
  "vault_mode": true
}
```

### Comandos de Skill
- `nux skill install <skill>` - Instala skill
- `nux skill enable <skill>` - Habilita skill  
- `nux skill list` - Lista skills disponíveis
- `nux skill info <skill>` - Info da skill
- `nux skill search <termo>` - Busca skills
- `nux skill upgrade <skill>` - Atualiza skill
- `nux skill sync` - Sincroniza com repositórios

### Exemplo de Uso
```bash
# Instalar skill Discord
nux skill install discord

# Habilitar e configurar
nux skill enable discord
# Solicita credenciais e salva no vault

# Usar através do NUX
nux discord --help
# Ou executa comando específico
nux ask "send discord message to channel X"
```

## Próximos Passos
1. Implementar proxy de comandos para skills habilitadas
2. Criptografar API keys no vault
3. Adicionar autocomplete para skills
4. Implementar atualização automática de skills
