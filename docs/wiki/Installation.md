# Installation

## Automatic Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/rsdenck/troncli/main/setup-tron.sh | bash
```

## Manual Install

Download the latest release from the [Releases Page](https://github.com/rsdenck/troncli/releases).

```bash
# Example for Linux AMD64
wget https://github.com/rsdenck/troncli/releases/download/v0.2.18/troncli_0.2.18_linux_amd64.tar.gz
tar -xzf troncli_0.2.18_linux_amd64.tar.gz
sudo mv troncli /usr/local/bin/
```

## Packages (.deb / .rpm)

You can also download `.deb` or `.rpm` packages from the releases page.

```bash
# Debian/Ubuntu
sudo dpkg -i troncli_0.2.18_linux_amd64.deb

# RHEL/CentOS/Fedora
sudo rpm -i troncli_0.2.18_linux_amd64.rpm
```
