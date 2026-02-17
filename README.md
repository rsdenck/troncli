<div align="center">

<img src="assets/banner.svg" width="100%" alt="TRONCLI" />

<table>
<tr>
<td align="left" width="55%">

### TRONCLI | System Administration TUI

<br>

**Production Grade Linux Tool**
Real-time Monitoring | LVM Management | Security Auditing

<br>

_"Building systems that do not wake people up at 3 AM."_

</td>
<td align="center" width="45%">
<h3>TRONCLI<br>INTERFACE</h3>
</td>
</tr>
</table>

</div>

### Status

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-000000?style=for-the-badge&logo=go&logoColor=00d9ff" />
  <img src="https://img.shields.io/badge/Platform-Linux-000000?style=for-the-badge&logo=linux&logoColor=00d9ff" />
  <img src="https://img.shields.io/badge/License-MIT-000000?style=for-the-badge&logoColor=00d9ff" />
  <img src="https://img.shields.io/badge/Build-Passing-000000?style=for-the-badge&logoColor=00d9ff" />
</div>

<br>

### Core Modules

<div align="center">

<table>
<tr>
<td valign="top" width="50%">

<h3>System Dashboard</h3>

<img src="https://img.shields.io/badge/CPU_Usage-000000?style=for-the-badge&logoColor=00d9ff" />
<img src="https://img.shields.io/badge/Memory_Stats-000000?style=for-the-badge&logoColor=00d9ff" />
<img src="https://img.shields.io/badge/Load_Average-000000?style=for-the-badge&logoColor=00d9ff" />
<img src="https://img.shields.io/badge/Network_IO-000000?style=for-the-badge&logoColor=00d9ff" />

<hr>

<h3>LVM Manager</h3>

<img src="https://img.shields.io/badge/Physical_Volumes-000000?style=for-the-badge&logoColor=00d9ff" />
<img src="https://img.shields.io/badge/Volume_Groups-000000?style=for-the-badge&logoColor=00d9ff" />
<img src="https://img.shields.io/badge/Logical_Volumes-000000?style=for-the-badge&logoColor=00d9ff" />

</td>

<td valign="top" width="50%">

<h3>Network Matrix</h3>

<img src="https://img.shields.io/badge/Interface_Stats-000000?style=for-the-badge&logoColor=00d9ff" />
<img src="https://img.shields.io/badge/RX_TX_Rates-000000?style=for-the-badge&logoColor=00d9ff" />
<img src="https://img.shields.io/badge/Socket_States-000000?style=for-the-badge&logoColor=00d9ff" />

<hr>

<h3>Security Audit</h3>

<img src="https://img.shields.io/badge/User_Enum-000000?style=for-the-badge&logoColor=00d9ff" />
<img src="https://img.shields.io/badge/SSH_Sessions-000000?style=for-the-badge&logoColor=00d9ff" />
<img src="https://img.shields.io/badge/Audit_Logs-000000?style=for-the-badge&logoColor=00d9ff" />

</td>
</tr>
</table>

</div>

### Installation

```bash
git clone https://github.com/rsdenck/troncli.git
cd troncli
go build -ldflags="-s -w" -o troncli cmd/troncli/main.go
./troncli
```

### Architecture

The system follows Clean Architecture principles with strict separation of concerns.

```text
cmd/
  troncli/       # Entry Point
internal/
  core/          # Domain Logic & Ports
  modules/       # Implementations (Linux Specific)
  ui/            # TUI Layer (tview/tcell)
```

### Security

Please report vulnerabilities to `ranlens.denck@protonmnail.com`.
See [SECURITY.md](SECURITY.md) for details.
