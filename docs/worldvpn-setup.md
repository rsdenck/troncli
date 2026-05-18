# WorldVPN Setup Guide for NUX CLI

This guide covers how to create a WorldVPN account, obtain credentials, configure NUX CLI with tunnel profiles, and connect using the nux wg command.

---

## 1. Creating a WorldVPN Account

### Option A: Free Trial (24 hours)

1. Open your browser and go to https://worldvpn.net/free-trial
2. Fill in the registration form with your name and a valid email address
3. Complete the captcha verification
4. Click Submit
5. Wait up to 5 minutes for the trial credentials to arrive in your email inbox
6. The email will contain your username and password for the VPN

### Option B: Paid Plan (Shared or Dedicated)

1. Go to https://worldvpn.net/pricing
2. Choose one of the three plans:
   - Shared Plan ($1.12/month): 200+ servers, 1 simultaneous connection
   - Dedicated Plan ($9/month): 1 dedicated server, 10 simultaneous connections
   - Reseller Plan ($1/month): 200+ servers, 3 simultaneous connections
3. Click "Get started" on your chosen plan
4. Complete the payment process
5. After payment confirmation, you will receive login credentials via email

---

## 2. Accessing Your Credentials

### From Email

The welcome email from WorldVPN contains:
- Username (format: xtrialXXXXX or similar)
- Password (numeric string)

### From Client Panel

1. Go to https://worldvpn.net/client/login
2. Log in with your credentials
3. Navigate to the dashboard to view or reset your VPN password
4. Your OpenVPN username and password are displayed there

---

## 3. Downloading Server Configurations

1. Go to https://worldvpn.net/servers
2. Browse the server list by country and location
3. For each server, click the OpenVPN Profile link (ZIP file)
4. Each ZIP file contains two configuration files:
   - `{location}_tcp.ovpn` (TCP protocol, port 80)
   - `{location}_udp.ovpn` (UDP protocol, port 80)

Example: Brazil S1 provides:
- https://worldvpn.net/ovpn/brazil%20s1.zip
- Contains: `brazil s1_tcp.ovpn` and `brazil s1_udp.ovpn`

---

## 4. Installing NUX (if not already installed)

```bash
curl -fsSL https://raw.githubusercontent.com/rsdenck/nux/main/src/install_nux.sh | bash
```

Or build from source:

```bash
git clone https://github.com/rsdenck/nux.git
cd nux
make build
sudo cp nux /usr/local/bin/
```

---

## 5. Configuring Tunnels in NUX

### Automatic Setup (downloads all 182 WorldVPN configs)

NUX includes a setup script that downloads all server configurations automatically:

```bash
bash scripts/setup-worldvpn-tunnels.sh
```

This script:
- Downloads all 182 ZIP files from WorldVPN (one per server location)
- Extracts the 364 OVPN configuration files (TCP and UDP for each server)
- Stores them in `~/.nux/tunnels/`
- Generates an index file at `~/.nux/tunnels/.index`
- Cleans up temporary files from `/tmp/`

### Manual Setup (add a single server)

```bash
# Create the tunnels directory
mkdir -p ~/.nux/tunnels

# Download a specific server configuration
curl -sL https://worldvpn.net/ovpn/brazil%20s1.zip -o /tmp/brazil_s1.zip

# Extract the .ovpn files
unzip -o /tmp/brazil_s1.zip -d ~/.nux/tunnels/

# Clean up
rm /tmp/brazil_s1.zip
```

### Verifying the configuration

```bash
nux wg list
```

Expected output shows all available tunnels:

```
Available tunnel configs (364):
TUNNEL                   TYPE     FILE
brazil_s1_tcp            OpenVPN  brazil_s1_tcp.ovpn
brazil_s1_udp            OpenVPN  brazil_s1_udp.ovpn
germany_s1_tcp           OpenVPN  germany_s1_tcp.ovpn
...
```

---

## 6. Connecting to a Tunnel

NUX uses OpenVPN to connect to WorldVPN tunnels. OpenVPN must be installed on the system.

### Install OpenVPN (if not installed)

Debian/Ubuntu:
```bash
sudo apt install openvpn
```

RHEL/Rocky/CentOS/Fedora:
```bash
sudo dnf install openvpn
```

Arch Linux:
```bash
sudo pacman -S openvpn
```

### Connect using NUX

Create an authentication file with your WorldVPN credentials:

```bash
echo "your_username" > /tmp/vpn-auth.txt
echo "your_password" >> /tmp/vpn-auth.txt
chmod 600 /tmp/vpn-auth.txt
```

Connect to a tunnel:

```bash
sudo openvpn --config ~/.nux/tunnels/brazil_s1_tcp.ovpn \
  --auth-user-pass /tmp/vpn-auth.txt \
  --daemon \
  --log /tmp/vpn-connect.log
```

Wait a few seconds and verify the connection:

```bash
# Check the VPN interface
ip addr show tun0

# Check public IP (should show the VPN server IP)
curl -4 ifconfig.me

# List active interfaces
nux wg list
```

### Disconnect

```bash
sudo pkill -f "openvpn.*tun0"
```

Or kill all OpenVPN processes:

```bash
sudo killall openvpn
```

---

## 7. NUX WireGuard Commands Reference

| Command | Description |
|---------|-------------|
| `nux wg list` | List all active WireGuard interfaces and available tunnel configs |
| `nux wg status` | Show WireGuard interface status from kernel |
| `nux wg show` | Show raw WireGuard interface configuration |
| `nux wg quick-status` | Show wg-quick managed interface status |
| `nux wg genkey` | Generate a WireGuard keypair |
| `nux wg genpsk` | Generate a WireGuard pre-shared key |
| `nux wg connect [config]` | Connect a WireGuard interface via wg-quick |
| `nux wg disconnect` | Disconnect WireGuard interface wg0 |
| `nux wg install` | Install WireGuard tools (wg, wg-quick, wgcf) |
| `nux wg warp generate` | Generate Cloudflare Warp config with wgcf |
| `nux wg warp register` | Register with Cloudflare Warp via wgcf |
| `nux wg warp connect` | Connect Cloudflare Warp via wgcf |
| `nux wg warp disconnect` | Disconnect from Cloudflare Warp |
| `nux wg warp status` | Check Cloudflare Warp connection status |

---

## 8. Troubleshooting

### Authentication Failed

```
AUTH: Received control message: AUTH_FAILED
```

Cause: Invalid username or password.
Solution: Verify credentials in the WorldVPN client panel. Reset the password if necessary.
         Confirm the trial period has not expired (trial lasts 24 hours).

### Connection Timeout

```
TCP connection established with [AF_INET]x.x.x.x:80
TLS: Initial packet from [AF_INET]x.x.x.x:80
[no further progress]
```

Cause: Network firewall blocking the VPN protocol.
Solution: Try the UDP configuration file instead of TCP, or vice versa.

### No .conf files found

```
nux wg connect: No .conf files found in /etc/wireguard/
```

Cause: `nux wg connect` expects WireGuard `.conf` files in `/etc/wireguard/`.
Note: WorldVPN provides OpenVPN (`.ovpn`) files, not WireGuard (`.conf`) files.
      Use the `openvpn` command directly to connect to WorldVPN tunnels.
      The `nux wg list` command shows them as available OpenVPN tunnel configs.

### Tunnel configs not showing

```
nux wg list shows: WireGuard interfaces: none
(no tunnel section)
```

Cause: The `~/.nux/tunnels/` directory is empty or does not exist.
Solution: Run the setup script again: `bash scripts/setup-worldvpn-tunnels.sh`
         Or manually copy .ovpn files to `~/.nux/tunnels/`.

---

## 9. Server Locations

WorldVPN provides 182 server locations across 35 countries:

| Country | Servers |
|---------|---------|
| Germany | 58 (S1-S58) |
| Netherlands | 10 (S1-S10) |
| United Kingdom | 8 (S1-S8) |
| France | 6 (S1-S6) |
| United States | California (5), Dallas (5), New York (5), Illinois (5), Florida (2), Seattle (2), Silicon Valley (1), Missouri (1) |
| Italy | 5 (S1-S5) |
| Latvia | 5 (S1-S5) |
| Poland | 6 (S1-S6) |
| Switzerland | 5 (S1-S5) |
| Bulgaria | 5 (S1-S5) |
| Ukraine | 5 (S1-S5) |
| Russia | 3 (S1-S3) |
| Singapore | 3 (S1-S3) |
| Australia | 4 (S1-S4) |
| Canada | 2 (S1-S2) |
| Austria | 2 (S1-S2) |
| Hong Kong | 4 (S1-S4) |
| Finland | 4 (S1-S4) |
| Sweden | 2 (S1-S2) |
| Others | Belgium, Brazil, Chile, Czech Republic, Greece, Iceland, Ireland, Japan, South Korea, Liechtenstein, Luxembourg, Malaysia, Norway, New Zealand, Philippines, Portugal, Romania, Spain, Turkey |

Each server provides both TCP and UDP configuration files.

---

## 10. Security Notes

- WorldVPN enforces a no-logs policy
- Connection uses AES-256-GCM encryption (TLS 1.3)
- Store your credentials securely. Do not commit them to version control
- The `~/.nux/tunnels/` directory contains VPN configuration files only. Credentials are never stored there
- Authentication files (`vpn-auth.txt`) should be created per-session and deleted after disconnection
- Always verify your public IP after connecting to confirm the VPN tunnel is active
