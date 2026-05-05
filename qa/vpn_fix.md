# Proton VPN Login Fix#

## Problem#

When using passwords with special characters like !, $, &, bash interprets them:
./nux proton vpn login --password cyX&KPjD3DYcybs6LNStMoQ5t$P48f&h$e!NqqX$
bash: !NqqX: event not found#

## Solution Applied#

### 1. Use --password-stdin#
Reads password from stdin, bypassing bash interpretation:
echo "password_with_special_chars" | ./nux proton vpn login --username user@protonmail.com --password-stdin#

### 2. Correct protonvpn command#
Changed from "signin" to "login" (correct protonvpn-cli subcommand)#

### 3. Vault integration#
After successful login, username is saved to vault:
- APIKeys["proton_username"] = username
- Config["proton_logged_in"] = true#

## Test#

# Login with special chars via stdin:
echo "pass!@#$%" | ./nux proton vpn login --username test@protonmail.com --password-stdin#

# Check vault:
./nux vault status#

# Connect after login:
./nux proton vpn connect#
