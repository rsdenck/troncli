# Contribuindo para o troncli

Obrigado pelo interesse em contribuir para o `troncli`! Este projeto segue padrões rigorosos de qualidade e segurança para garantir que seja adequado para uso em produção.

## Código de Conduta

Ao participar deste projeto, você concorda em seguir nosso Código de Conduta (implícito: respeito mútuo, profissionalismo e colaboração construtiva).

## Git Commit Messages
*   **IMPORTANT**: We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.
    *   `feat: add new feature`
    *   `fix: resolve bug`
    *   `docs: update documentation`
    *   `style: formatting, missing semi colons, etc; no code change`
    *   `refactor: refactoring production code`
    *   `test: adding missing tests, refactoring tests`
    *   `chore: updating build tasks, package manager configs, etc`

## Como Contribuir

1.  **Fork o repositório**.
2.  **Crie uma branch** para sua feature ou correção (`git checkout -b feature/nova-funcionalidade`).
3.  **Implemente suas mudanças**.
    *   Siga a Clean Architecture.
    *   **NUNCA** use dados falsos (mock). Use apenas fontes reais do sistema Linux (`/proc`, `/sys`, syscalls).
    *   Escreva testes unitários para cobrir suas mudanças.
    *   Execute `golangci-lint` localmente para garantir conformidade.
4.  **Assine seus commits** (DCO). Use `git commit -s -m "mensagem descritiva"`.
5.  **Envie para o GitHub** (`git push origin feature/nova-funcionalidade`).
6.  **Abra um Pull Request**.

## Padrões de Qualidade

*   **Go Idiomático**: Siga as melhores práticas da linguagem Go.
*   **Testes**: Cobertura mínima de 85%.
*   **Segurança**: Sem vulnerabilidades conhecidas, sem segredos no código.
*   **Performance**: O código deve ser eficiente e não bloquear a UI.
*   **Documentação**: Comentários claros e documentação atualizada.

## Estrutura do Projeto

*   `cmd/troncli`: Ponto de entrada da aplicação.
*   `internal/core`: Lógica de negócios e interfaces (Clean Architecture).
*   `internal/modules`: Implementações concretas dos módulos.
*   `internal/ui`: Interface do usuário (TUI).
*   `internal/collectors`: Coletores de dados do sistema.

## Reportando Bugs

Use o Issue Tracker do GitHub. Forneça detalhes sobre o ambiente (OS, Kernel, Versão do troncli) e passos para reproduzir o problema.

## Licença

Este projeto é licenciado sob a MIT License.
