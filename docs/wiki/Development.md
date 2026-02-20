# Development

## Prerequisites

- Go 1.24+
- Make
- Docker (for testing packaging)

## Building

```bash
make build
```

## Running Tests

```bash
make test
```

## Architecture

The project follows Clean Architecture:

- **cmd/**: Entry points.
- **internal/core/ports**: Interfaces (Hexagonal Architecture).
- **internal/core/services**: Business logic.
- **internal/core/adapter**: Implementations (Repositories, Executors).
- **internal/modules**: Feature-specific implementations.
- **internal/ui**: Presentation layer (Console/TUI).
