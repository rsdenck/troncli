#!/bin/bash
# setup-worldvpn-tunnels.sh
# Downloads all WorldVPN OpenVPN configs and configures them as standard tunnels for NUX
set -e

NUX_TUNNELS_DIR="${HOME}/.nux/tunnels"
TMP_DIR="/tmp/worldvpn_download"
BASE_URL="https://worldvpn.net/ovpn"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
log_info()  { echo -e "${GREEN}[INFO]${NC}  $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $1"; }
log_err()   { echo -e "${RED}[ERROR]${NC} $1"; }
log_step()  { echo -e "${CYAN}  →${NC} $1"; }

# ASCII header
echo -e "\033[38;5;208m"
echo " ███╗   ██╗██╗   ██╗██╗  ██╗"
echo " ████╗  ██║██║   ██║╚██╗██╔╝"
echo " ██╔██╗ ██║██║   ██║ ╚███╔╝ "
echo " ██║╚██╗██║██║   ██║ ██╔██╗ "
echo " ██║ ╚████║╚██████╔╝██╔╝ ██╗"
echo " ╚═╝  ╚═══╝ ╚═════╝ ╚═╝  ╚═╝"
echo -e "\033[0m"
echo -e "\033[38;5;208mNUX — WorldVPN Tunnel Setup\033[0m"
echo "================================================"

mkdir -p "$NUX_TUNNELS_DIR" "$TMP_DIR"

# All VPN zip files from WorldVPN (lowercase, with spaces)
ZIPS=(
  "germany s1.zip" "germany s2.zip" "germany s3.zip" "germany s4.zip"
  "germany s5.zip" "germany s6.zip" "germany s7.zip" "germany s8.zip"
  "germany s9.zip" "germany s10.zip" "germany s11.zip" "germany s12.zip"
  "germany s13.zip" "germany s14.zip" "germany s15.zip" "germany s16.zip"
  "germany s17.zip" "germany s18.zip" "germany s19.zip" "germany s20.zip"
  "germany s21.zip" "germany s22.zip" "germany s23.zip" "germany s24.zip"
  "germany s25.zip" "germany s26.zip" "germany s27.zip" "germany s28.zip"
  "germany s29.zip" "germany s30.zip" "germany s31.zip" "germany s32.zip"
  "germany s33.zip" "germany s34.zip" "germany s35.zip" "germany s36.zip"
  "germany s37.zip" "germany s38.zip" "germany s39.zip" "germany s40.zip"
  "germany s41.zip" "germany s42.zip" "germany s43.zip" "germany s44.zip"
  "germany s45.zip" "germany s46.zip" "germany s47.zip" "germany s48.zip"
  "germany s49.zip" "germany s50.zip" "germany s51.zip" "germany s52.zip"
  "germany s53.zip" "germany s54.zip" "germany s55.zip" "germany s56.zip"
  "germany s57.zip" "germany s58.zip"
  "austria s1.zip" "austria s2.zip"
  "australia s1.zip" "australia s2.zip" "australia s3.zip" "australia s4.zip"
  "belgium.zip"
  "bulgaria s1.zip" "bulgaria s2.zip" "bulgaria s3.zip" "bulgaria s4.zip" "bulgaria s5.zip"
  "brazil s1.zip"
  "canada s1.zip" "canada s2.zip"
  "swiss s1.zip" "swiss s2.zip" "swiss s3.zip" "swiss s4.zip" "swiss s5.zip"
  "chile s1.zip"
  "czech republic s1.zip"
  "spain s1.zip"
  "finland s1.zip" "finland s2.zip" "finland s3.zip" "finland s4.zip"
  "france s1.zip" "france s2.zip" "france s3.zip" "france s4.zip" "france s5.zip" "france s6.zip"
  "greece s1.zip"
  "hong kong s1.zip" "hong kong s2.zip" "hong kong s3.zip" "hong kong s4.zip"
  "ireland s1.zip"
  "iceland s1.zip"
  "italy s1.zip" "italy s2.zip" "italy s3.zip" "italy s4.zip" "italy s5.zip"
  "japan s1.zip"
  "south korea.zip"
  "liechtenstein s1.zip"
  "luxembourg s1.zip"
  "latvia s1.zip" "latvia s2.zip" "latvia s3.zip" "latvia s4.zip" "latvia s5.zip"
  "malaysia s1.zip"
  "netherlands s1.zip" "netherlands s2.zip" "netherlands s3.zip" "netherlands s4.zip"
  "netherlands s5.zip" "netherlands s6.zip" "netherlands s7.zip" "netherlands s8.zip"
  "netherlands s9.zip" "netherlands s10.zip"
  "norway s1.zip"
  "new zealand s1.zip"
  "philippines s1.zip"
  "poland s1.zip" "poland s2.zip" "poland s3.zip" "poland s4.zip" "poland s5.zip" "poland s6.zip"
  "portugal.zip"
  "romania.zip"
  "russia s1.zip" "russia s2.zip" "russia s3.zip"
  "sweden s1.zip" "sweden s2.zip"
  "singapore s1.zip" "singapore s2.zip" "singapore s3.zip"
  "turkey.zip"
  "ukraine s1.zip" "ukraine s2.zip" "ukraine s3.zip" "ukraine s4.zip" "ukraine s5.zip"
  "united kingdom s1.zip" "united kingdom s2.zip" "united kingdom s3.zip"
  "united kingdom s4.zip" "united kingdom s5.zip" "united kingdom s6.zip"
  "united kingdom s7.zip" "united kingdom s8.zip"
  "california s1.zip" "california s2.zip" "california s3.zip" "california s4.zip" "california s5.zip"
  "dallas s1.zip" "dallas s2.zip" "dallas s3.zip" "dallas s4.zip" "dallas s5.zip"
  "florida s1.zip" "florida s2.zip"
  "illinois s1.zip" "illinois s2.zip" "illinois s3.zip" "illinois s4.zip" "illinois s5.zip"
  "missouri s1.zip"
  "new york s1.zip" "new york s2.zip" "new york s3.zip" "new york s4.zip" "new york s5.zip"
  "seattle s1.zip" "seattle s2.zip"
  "silicon valley s1.zip"
)

TOTAL=${#ZIPS[@]}
COUNT=0

log_info "Found $TOTAL tunnel configs to download"

# Download in parallel with xargs (8 at a time)
download_one() {
  local zipname="$1"
  local url_safe=$(echo "$zipname" | sed 's/ /%20/g' | tr '[:upper:]' '[:lower:]')
  local url="$BASE_URL/$url_safe"
  local dest="$TMP_DIR/$zipname"
  if curl -sL -o "$dest" "$url" --connect-timeout 10 --max-time 60 2>/dev/null && [ -s "$dest" ]; then
    echo "OK:$zipname"
  else
    echo "FAIL:$zipname"
  fi
}
export BASE_URL TMP_DIR
export -f download_one

printf "%s\n" "${ZIPS[@]}" | xargs -P 8 -I {} bash -c 'download_one "$@"' _ {} 2>&1 | \
  awk -F: '{ if($1=="OK"){ok++; printf "."} else {fail++; printf "x"} } END { print "\nOK="ok,"FAIL="fail }'

# Extract all zips
echo ""
log_step "Extracting all zip files..."
for zipf in "$TMP_DIR"/*.zip; do
  [ -f "$zipf" ] || continue
  unzip -o -q -d "$TMP_DIR" "$zipf" 2>/dev/null || true
done

# Count extracted ovpn files
OVPN_COUNT=$(ls "$TMP_DIR"/*.ovpn 2>/dev/null | wc -l)
log_step "Found $OVPN_COUNT .ovpn config files"

# Move to NUX tunnels directory
log_step "Copying tunnel configs to $NUX_TUNNELS_DIR..."
for f in "$TMP_DIR"/*.ovpn; do
  [ -f "$f" ] || continue
  base=$(basename "$f")
  safe_name=$(echo "$base" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')
  cp "$f" "$NUX_TUNNELS_DIR/$safe_name"
done

# Generate index
log_step "Generating tunnel index..."
INDEX_FILE="$NUX_TUNNELS_DIR/.index"
echo "# NUX Tunnel Index - WorldVPN" > "$INDEX_FILE"
echo "# Generated: $(date -u)" >> "$INDEX_FILE"
echo "# Total OVPN configs: $(ls "$NUX_TUNNELS_DIR"/*.ovpn 2>/dev/null | wc -l)" >> "$INDEX_FILE"
echo "" >> "$INDEX_FILE"

for f in "$NUX_TUNNELS_DIR"/*.ovpn; do
  [ -f "$f" ] || continue
  base=$(basename "$f")
  name="${base%.*}"
  remote_ip=$(grep -m1 "^remote " "$f" | awk '{print $2}' || echo "unknown")
  remote_port=$(grep -m1 "^remote " "$f" | awk '{print $3}' || echo "1194")
  proto=$(grep -m1 "^proto " "$f" | awk '{print $2}' || echo "udp")
  echo "$name|$remote_ip|$remote_port|$proto|$base" >> "$INDEX_FILE"
done

# Quick stats
FINAL_COUNT=$(ls "$NUX_TUNNELS_DIR"/*.ovpn 2>/dev/null | wc -l)
echo ""
log_info "Download complete!"
echo "  Tunnels stored: $NUX_TUNNELS_DIR"
echo "  OVPN configs:   $FINAL_COUNT"
echo ""

# Cleanup /tmp/
log_step "Cleaning temporary files..."
rm -rf "$TMP_DIR"
find /tmp -name "*.ovpn" -maxdepth 1 -type f -delete 2>/dev/null || true
find /tmp -name "*worldvpn*" -maxdepth 1 -exec rm -rf {} + 2>/dev/null || true

echo ""
echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}  WorldVPN tunnels configured for NUX!${NC}"
echo -e "${GREEN}  $FINAL_COUNT tunnel configs ready${NC}"
echo -e "${GREEN}================================================${NC}"
