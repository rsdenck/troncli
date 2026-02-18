package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/ui/themes"
	"github.com/rivo/tview"
)

type DiskView struct {
	*tview.Flex
	manager ports.DiskManager
	table   *tview.Table
}

func NewDiskView(manager ports.DiskManager) *DiskView {
	v := &DiskView{
		Flex:    tview.NewFlex(),
		manager: manager,
		table:   tview.NewTable(),
	}

	v.setupUI()
	v.refreshData()
	return v
}

func (v *DiskView) setupUI() {
	v.SetDirection(tview.FlexRow)
	v.SetBorder(true).SetTitle(" Disk & Storage ").SetBorderColor(themes.TronCyan)

	v.table.SetBorders(true).SetBorderColor(themes.TronBlue)
	v.table.SetSelectable(true, false)

	v.AddItem(v.table, 0, 1, true)
}

func (v *DiskView) refreshData() {
	v.table.Clear()

	// Headers
	headers := []string{"Name", "Size", "Type", "MountPoint", "Usage"}
	for i, h := range headers {
		v.table.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(themes.TronYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}

	devices, err := v.manager.ListBlockDevices()
	if err != nil {
		v.table.SetCell(1, 0, tview.NewTableCell("Error: "+err.Error()).SetTextColor(tcell.ColorRed))
		return
	}

	row := 1
	var addDevice func(d ports.BlockDevice, prefix string)
	addDevice = func(d ports.BlockDevice, prefix string) {
		name := prefix + d.Name

		usageStr := "-"
		if d.MountPoint != "" {
			usage, err := v.manager.GetFilesystemUsage(d.MountPoint)
			if err == nil {
				percent := 0.0
				if usage.Total > 0 {
					percent = (float64(usage.Used) / float64(usage.Total)) * 100
				}
				usageStr = fmt.Sprintf("%.1f%%", percent)
			}
		}

		v.table.SetCell(row, 0, tview.NewTableCell(name).SetTextColor(themes.TronCyan))
		v.table.SetCell(row, 1, tview.NewTableCell(d.Size).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 2, tview.NewTableCell(d.Type).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 3, tview.NewTableCell(d.MountPoint).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 4, tview.NewTableCell(usageStr).SetTextColor(themes.TronGreen))

		row++
		for _, child := range d.Children {
			addDevice(child, prefix+"  └─ ")
		}
	}

	for _, d := range devices {
		addDevice(d, "")
	}
}
