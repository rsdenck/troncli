# Status

Show service status (alias for 'service status')

## Usage

```
nux status <service>
```

## Description

Displays the current status of a systemd service with colored output indicating whether the service is running (green), stopped (orange), or has problems (red). Shows key service metadata in a compact table.

## Examples

Check status of httpd service:
nux status httpd

Check status of sshd service:
nux status sshd
