package views

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/ui/themes"
	"github.com/rivo/tview"
)

type NetworkView struct {
	*tview.Flex
	manager ports.NetworkManager
	table   *tview.Table
}

func NewNetworkView(manager ports.NetworkManager) *NetworkView {
	v := &NetworkView{
		Flex:    tview.NewFlex(),
		manager: manager,
		table:   tview.NewTable(),
	}

	v.setupUI()
	v.refreshData()
	return v
}

func (v *NetworkView) setupUI() {
	v.SetDirection(tview.FlexRow)
	v.SetBorder(true).SetTitle(" Network Interfaces ").SetBorderColor(themes.TronCyan)

	v.table.SetBorders(true).SetBorderColor(themes.TronBlue)
	v.table.SetSelectable(true, false)

	v.AddItem(v.table, 0, 1, true)
}

func (v *NetworkView) refreshData() {
	v.table.Clear()

	// Headers
	headers := []string{"Name", "Status", "MAC Address", "IP Addresses"}
	for i, h := range headers {
		v.table.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(themes.TronYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}

	ifaces, err := v.manager.GetInterfaces()
	if err != nil {
		v.table.SetCell(1, 0, tview.NewTableCell("Error: "+err.Error()).SetTextColor(tcell.ColorRed))
		return
	}

	for i, iface := range ifaces {
		row := i + 1
		status := "DOWN"
		color := tcell.ColorRed
		if iface.Flags&1 != 0 { // net.FlagUp
			status = "UP"
			color = themes.TronGreen
		}

		v.table.SetCell(row, 0, tview.NewTableCell(iface.Name).SetTextColor(themes.TronCyan))
		v.table.SetCell(row, 1, tview.NewTableCell(status).SetTextColor(color))
		v.table.SetCell(row, 2, tview.NewTableCell(iface.HardwareAddr).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 3, tview.NewTableCell(strings.Join(iface.IPAddresses, ", ")).SetTextColor(tcell.ColorWhite))
	}
}
