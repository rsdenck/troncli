package commands

import (
	"fmt"
	"net/netip"
	"os"

	"github.com/oschwald/geoip2-golang/v2"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var geoipCmd = &cobra.Command{
	Use:   "geoip",
	Short: "IP geolocation and analysis",
	Long:  `Geolocate IPs, analyze logs, investigate connections, and integrate with firewall.`,
}

var geoipLookupCmd = &cobra.Command{
	Use:   "lookup <ip>",
	Short: "Lookup geolocation for an IP address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ipStr := args[0]

		ip, err := netip.ParseAddr(ipStr)
		if err != nil {
			output.NewError(fmt.Sprintf("invalid IP address: %s", ipStr), "GEOIP_INVALID_IP").Print()
			return
		}

		dbPath := "/opt/nux/geoip/GeoLite2-City.mmdb"
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			output.NewError("GeoIP database not found. Please install GeoLite2-City.mmdb at /opt/nux/geoip/", "GEOIP_DB_NOT_FOUND").Print()
			output.NewInfo("You can download it from: https://www.maxmind.com/en/geolite2/signup").Print()
			return
		}

		db, err := geoip2.Open(dbPath)
		if err != nil {
			output.NewError(fmt.Sprintf("failed to open GeoIP database: %s", err.Error()), "GEOIP_DB_ERROR").Print()
			return
		}
		defer db.Close()

		record, err := db.City(ip)
		if err != nil {
			output.NewError(fmt.Sprintf("failed to lookup IP: %s", err.Error()), "GEOIP_LOOKUP_ERROR").Print()
			return
		}

		result := map[string]interface{}{
			"ip":          ipStr,
			"country":     record.Country.ISOCode,
			"latitude":    record.Location.Latitude,
			"longitude":   record.Location.Longitude,
			"timezone":    record.Location.TimeZone,
			"postal_code": record.Postal.Code,
		}

		output.NewSuccess(result).WithMessage(fmt.Sprintf("Geolocation for %s", ipStr)).Print()
	},
}

var geoipWhoisCmd = &cobra.Command{
	Use:   "whois <ip>",
	Short: "Whois lookup for an IP address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ipStr := args[0]

		ip, err := netip.ParseAddr(ipStr)
		if err != nil {
			output.NewError(fmt.Sprintf("invalid IP address: %s", ipStr), "GEOIP_INVALID_IP").Print()
			return
		}

		asnDbPath := "/opt/nux/geoip/GeoLite2-ASN.mmdb"
		if _, err := os.Stat(asnDbPath); os.IsNotExist(err) {
			output.NewError("ASN database not found. Please install GeoLite2-ASN.mmdb at /opt/nux/geoip/", "GEOIP_ASN_NOT_FOUND").Print()
			return
		}

		db, err := geoip2.Open(asnDbPath)
		if err != nil {
			output.NewError(fmt.Sprintf("failed to open ASN database: %s", err.Error()), "GEOIP_DB_ERROR").Print()
			return
		}
		defer db.Close()

		record, err := db.ASN(ip)
		if err != nil {
			output.NewError(fmt.Sprintf("failed to lookup ASN: %s", err.Error()), "GEOIP_WHOIS_ERROR").Print()
			return
		}

		result := map[string]interface{}{
			"ip":           ipStr,
			"asn":          record.AutonomousSystemNumber,
			"organization": record.AutonomousSystemOrganization,
			"network":      record.Network.String(),
		}

		output.NewSuccess(result).WithMessage(fmt.Sprintf("Whois for %s", ipStr)).Print()
	},
}

var geoipSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup GeoIP databases",
	Run: func(cmd *cobra.Command, args []string) {
		geoipDir := "/opt/nux/geoip"

		if err := os.MkdirAll(geoipDir, 0755); err != nil {
			output.NewError(fmt.Sprintf("failed to create directory: %s", err.Error()), "GEOIP_SETUP_ERROR").Print()
			return
		}

		output.NewInfo("GeoIP setup").Print()
		output.NewInfo(fmt.Sprintf("Directory created: %s", geoipDir)).Print()
		output.NewInfo("Please download the following databases from MaxMind:").Print()
		output.NewInfo("1. GeoLite2-City.mmdb").Print()
		output.NewInfo("2. GeoLite2-ASN.mmdb").Print()
		output.NewInfo("3. GeoLite2-Country.mmdb").Print()
		output.NewInfo("\nDownload URL: https://www.maxmind.com/en/geolite2/signup").Print()
		output.NewInfo(fmt.Sprintf("\nAfter download, place the .mmdb files in: %s", geoipDir)).Print()

		output.NewSuccess(map[string]interface{}{
			"status":       "setup_complete",
			"database_dir": geoipDir,
			"next_step":    "Download and place .mmdb files",
		}).Print()
	},
}

func init() {
	geoipCmd.AddCommand(geoipLookupCmd)
	geoipCmd.AddCommand(geoipWhoisCmd)
	geoipCmd.AddCommand(geoipSetupCmd)
	rootCmd.AddCommand(geoipCmd)
}
