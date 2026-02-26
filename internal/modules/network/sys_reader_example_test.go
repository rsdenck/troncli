package network_test

import (
	"fmt"
	"log"

	"github.com/mascli/troncli/internal/modules/network"
)

// ExampleSysReader_ReadInterfaces demonstrates how to use SysReader to read network interfaces
func ExampleSysReader_ReadInterfaces() {
	// Create a new SysReader instance
	reader := network.NewSysReader()

	// Read all network interfaces from /sys/class/net
	interfaces, err := reader.ReadInterfaces()
	if err != nil {
		log.Fatalf("Failed to read interfaces: %v", err)
	}

	// Display interface information
	for _, iface := range interfaces {
		fmt.Printf("Interface: %s\n", iface.Name)
		fmt.Printf("  MAC Address: %s\n", iface.HardwareAddr)
		fmt.Printf("  MTU: %d\n", iface.MTU)
		fmt.Printf("  State: %s\n", iface.State)
		fmt.Printf("  Index: %d\n", iface.Index)
		fmt.Printf("  IP Addresses: %v\n", iface.IPAddresses)
		fmt.Println()
	}
}
