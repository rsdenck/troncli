# ProtonVPN CLI

- **Repo:** https://github.com/rsdenck/skillnux/vpn/protonvpn-cli_install.go
- **Description:** ProtonVPN CLI for Linux - official command-line tool
- **Category:** VPN/Security

## Commands

### protonvpn
Main CLI tool for ProtonVPN.

### protonvpn signin
Sign in to your Proton account.

### protonvpn connect
Connect to VPN server (fastest by default).

### protonvpn disconnect
Disconnect from VPN.

### protonvpn status
Show connection status.

### protonvpn servers
List available servers.

## Install

```bash
# The NUX CLI will automatically detect your distribution and install the correct version.
# Or run manually:
nux skill install protonvpn-cli
```

## Supported Distributions

- Debian/Ubuntu: Installs via official Proton repository
- Fedora: Installs via DNF and official repository  
- Arch: Installs via AUR (protonvpn-cli)
- Other: Attempts generic installation

## Usage

After installation, use the official ProtonVPN CLI commands:

```bash
# Sign in
protonvpn signin your_username

# Connect to fastest server
protonvpn connect --fastest

# Check status
protonvpn status

# Disconnect
protonvpn disconnect
```

## Notes

- Requires gnome-keyring or KDE KWallet for credential storage
- Does not work on headless setups
- Official documentation: https://protonvpn.com/support/linux-cli
