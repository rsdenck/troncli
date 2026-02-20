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
	details *tview.TextView
}

func NewDiskView(manager ports.DiskManager) *DiskView {
	v := &DiskView{
		Flex:    tview.NewFlex(),
		manager: manager,
		table:   tview.NewTable(),
		details: tview.NewTextView(),
	}

	v.setupUI()
	v.refreshData()
	return v
}

func (v *DiskView) setupUI() {
	v.SetDirection(tview.FlexRow)
	
	v.table.SetBorders(true).SetTitle(" DISK & STORAGE ").SetBorderColor(themes.TronCyan).SetTitleColor(themes.TronCyan)
	v.table.SetSelectable(true, false)

	v.details.SetBorder(true).SetTitle(" PARTITION DETAILS ").SetBorderColor(themes.TronBlue).SetTitleColor(themes.TronBlue)
	v.details.SetDynamicColors(true)

	v.AddItem(v.table, 0, 3, true)
	v.AddItem(v.details, 0, 1, false)

	v.table.SetSelectionChangedFunc(func(row, column int) {
		if row > 0 {
			cell := v.table.GetCell(row, 0)
			if cell != nil {
				// Just showing simple details for now, could be expanded
				name := cell.Text
				v.details.SetText(fmt.Sprintf("%sSelected Device: %s%s", themes.ColorNeonCyan, themes.ColorWhite, name))
			}
		}
	})
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
		usageColor := tcell.ColorWhite
		
		if d.MountPoint != "" {
			usage, err := v.manager.GetFilesystemUsage(d.MountPoint)
			if err == nil {
				percent := 0.0
				if usage.Total > 0 {
					percent = (float64(usage.Used) / float64(usage.Total)) * 100
				}
				usageStr = fmt.Sprintf("%.1f%%", percent)
				
				if percent > 90 {
					usageColor = tcell.ColorRed
				} else if percent > 75 {
					usageColor = tcell.ColorYellow
				} else {
					usageColor = tcell.ColorGreen
				}
			}
		}

		v.table.SetCell(row, 0, tview.NewTableCell(name).SetTextColor(themes.TronCyan))
		v.table.SetCell(row, 1, tview.NewTableCell(d.Size).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 2, tview.NewTableCell(d.Type).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 3, tview.NewTableCell(d.MountPoint).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 4, tview.NewTableCell(usageStr).SetTextColor(usageColor))

		row++
		for _, child := range d.Children {
			addDevice(child, prefix+"  └─ ")
		}
	}

	for _, d := range devices {
		addDevice(d, "")
	}
}
