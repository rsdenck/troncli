# Install#

Instala pacotes (universal)#

## Usage#

```
nux install <packages>
```

## Description#

Installs packages using the detected package manager (apt, dnf, yum, pacman, zypper, apk). Supports dry-run mode.

## Examples#

Install nginx:
nux install nginx

Dry run:
nux install --dry-run vim
