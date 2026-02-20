# Configuration

TronCLI looks for configuration files in the following locations:
1. `/etc/troncli/config.yaml`
2. `~/.config/troncli/config.yaml`

## Example Configuration

```yaml
general:
  timeout: 30
  verbose: false
  no_color: false

network:
  interface: eth0
  dns_servers:
    - 8.8.8.8
    - 1.1.1.1

remote:
  ssh_key_path: ~/.ssh/id_rsa
  default_user: root

themes:
  active: cyberpunk
```
