# Usage

TronCLI can be used in two modes: **Interactive TUI** and **Command Line Interface**.

## Interactive TUI

Simply run `troncli` without arguments to start the dashboard.

```bash
troncli
```

Use arrow keys to navigate, `Enter` to select, and `Esc` or `q` to go back/exit.

## Command Line Interface

### System Information
```bash
troncli system info
troncli system profile
```

### Network
```bash
troncli network scan --target google.com
troncli network trace google.com
troncli network sockets
```

### Services
```bash
troncli service list
troncli service status sshd
troncli service logs sshd
```

### Disk
```bash
troncli disk health
troncli disk usage
```
