# Política de Segurança

## Versões Suportadas

Apenas a versão estável mais recente do `troncli` recebe atualizações de segurança.

| Versão | Suportado          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |
| < 1.0   | :x:                |

## Reportando uma Vulnerabilidade

**NÃO** relate vulnerabilidades de segurança através de issues públicas no GitHub.

Se você acredita ter encontrado uma vulnerabilidade de segurança no `troncli`, por favor, utilize a funcionalidade de **Private Vulnerability Reporting** do GitHub neste repositório, ou entre em contato diretamente com os mantenedores.

### Processo

1.  **Relato**: Você reporta a vulnerabilidade de forma privada.
2.  **Confirmação**: Confirmamos o recebimento em até 48 horas.
3.  **Investigação**: Investigamos o problema e determinamos o impacto.
4.  **Correção**: Desenvolvemos uma correção e testamos.
5.  **Release**: Lançamos uma versão corrigida.
6.  **Divulgação**: Divulgamos publicamente a vulnerabilidade após a correção estar disponível.

## Segurança da Cadeia de Suprimentos

Levamos a segurança da cadeia de suprimentos a sério.
- Todos os releases são assinados com Cosign.
- SBOMs são gerados para cada release (CycloneDX).
- Dependências são escaneadas diariamente.
- Pipelines de CI/CD são pinados por SHA.
