# TronCLI

Advanced Linux System Management CLI tool built with Go.

## Installation

### Automatic Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/mascli/troncli/main/install.sh | bash
```

### Manual Install

Download the latest release from the [Releases Page](https://github.com/mascli/troncli/releases).

```bash
# Example for Linux AMD64
wget https://github.com/mascli/troncli/releases/download/v0.2.18/troncli_0.2.18_linux_amd64.tar.gz
tar -xzf troncli_0.2.18_linux_amd64.tar.gz
sudo mv troncli /usr/local/bin/
```

### Packages (.deb / .rpm)

You can also download `.deb` or `.rpm` packages from the releases page.

```bash
# Debian/Ubuntu
sudo dpkg -i troncli_0.2.18_linux_amd64.deb

# RHEL/CentOS/Fedora
sudo rpm -i troncli_0.2.18_linux_amd64.rpm
```

## Development

### Requirements

- Go 1.24+
- Make

### Build Locally

```bash
make build
```

### Cross-Compile Script

```bash
./scripts/build.sh
```

### Release Process

Releases are automated via GitHub Actions when a new tag is pushed.

1. Create a new tag:
   ```bash
   git tag -a v0.2.18 -m "Release v0.2.18"
   ```
2. Push the tag:
   ```bash
   git push origin v0.2.18
   ```

The CI pipeline will:
- Run tests
- Build binaries for Linux (amd64, arm64, armv7)
- Generate packages (tar.gz, deb, rpm)
- Create a GitHub Release
- Upload artifacts

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
