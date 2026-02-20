package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/rivo/tview"
)

type LVMView struct {
	*tview.Flex
	pvTable *tview.Table
	vgTable *tview.Table
	lvTable *tview.Table
	manager ports.LVMManager
}

func NewLVMView(manager ports.LVMManager) *LVMView {
	v := &LVMView{
		Flex:    tview.NewFlex(),
		pvTable: tview.NewTable(),
		vgTable: tview.NewTable(),
		lvTable: tview.NewTable(),
		manager: manager,
	}

	v.setupUI()
	v.loadData()

	return v
}

func (v *LVMView) setupUI() {
	v.SetDirection(tview.FlexRow)

	// PV Table
	v.pvTable.SetBorders(true).SetTitle(" PHYSICAL VOLUMES ").SetTitleColor(tcell.ColorAqua).SetBorderColor(tcell.ColorAqua)
	v.pvTable.SetSelectable(true, false)

	// VG Table
	v.vgTable.SetBorders(true).SetTitle(" VOLUME GROUPS ").SetTitleColor(tcell.ColorAqua).SetBorderColor(tcell.ColorAqua)
	v.vgTable.SetSelectable(true, false)

	// LV Table
	v.lvTable.SetBorders(true).SetTitle(" LOGICAL VOLUMES ").SetTitleColor(tcell.ColorAqua).SetBorderColor(tcell.ColorAqua)
	v.lvTable.SetSelectable(true, false)

	// Layout: Top row (PV + VG), Bottom row (LV)
	topRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	topRow.AddItem(v.pvTable, 0, 1, false)
	topRow.AddItem(v.vgTable, 0, 1, false)

	v.AddItem(topRow, 0, 1, false)
	v.AddItem(v.lvTable, 0, 1, true)
}

func (v *LVMView) loadData() {
	// PVs
	v.pvTable.SetCell(0, 0, tview.NewTableCell("PV NAME").SetTextColor(tcell.ColorYellow))
	v.pvTable.SetCell(0, 1, tview.NewTableCell("VG NAME").SetTextColor(tcell.ColorYellow))
	v.pvTable.SetCell(0, 2, tview.NewTableCell("SIZE").SetTextColor(tcell.ColorYellow))
	v.pvTable.SetCell(0, 3, tview.NewTableCell("FREE").SetTextColor(tcell.ColorYellow))

	pvs, err := v.manager.ListPhysicalVolumes()
	if err == nil {
		for i, pv := range pvs {
			row := i + 1
			v.pvTable.SetCell(row, 0, tview.NewTableCell(pv.Name).SetTextColor(tcell.ColorWhite))
			v.pvTable.SetCell(row, 1, tview.NewTableCell(pv.VGName).SetTextColor(tcell.ColorWhite))
			v.pvTable.SetCell(row, 2, tview.NewTableCell(pv.Size).SetTextColor(tcell.ColorWhite))
			v.pvTable.SetCell(row, 3, tview.NewTableCell(pv.Free).SetTextColor(tcell.ColorWhite))
		}
	}

	// VGs
	v.vgTable.SetCell(0, 0, tview.NewTableCell("VG NAME").SetTextColor(tcell.ColorYellow))
	v.vgTable.SetCell(0, 1, tview.NewTableCell("SIZE").SetTextColor(tcell.ColorYellow))
	v.vgTable.SetCell(0, 2, tview.NewTableCell("FREE").SetTextColor(tcell.ColorYellow))
	v.vgTable.SetCell(0, 3, tview.NewTableCell("PVs").SetTextColor(tcell.ColorYellow))
	v.vgTable.SetCell(0, 4, tview.NewTableCell("LVs").SetTextColor(tcell.ColorYellow))

	vgs, err := v.manager.ListVolumeGroups()
	if err == nil {
		for i, vg := range vgs {
			row := i + 1
			v.vgTable.SetCell(row, 0, tview.NewTableCell(vg.Name).SetTextColor(tcell.ColorWhite))
			v.vgTable.SetCell(row, 1, tview.NewTableCell(vg.Size).SetTextColor(tcell.ColorWhite))
			v.vgTable.SetCell(row, 2, tview.NewTableCell(vg.Free).SetTextColor(tcell.ColorWhite))
			v.vgTable.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%d", vg.PVCount)).SetTextColor(tcell.ColorWhite))
			v.vgTable.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%d", vg.LVCount)).SetTextColor(tcell.ColorWhite))
		}
	}

	// LVs
	v.lvTable.SetCell(0, 0, tview.NewTableCell("LV NAME").SetTextColor(tcell.ColorYellow))
	v.lvTable.SetCell(0, 1, tview.NewTableCell("VG NAME").SetTextColor(tcell.ColorYellow))
	v.lvTable.SetCell(0, 2, tview.NewTableCell("PATH").SetTextColor(tcell.ColorYellow))
	v.lvTable.SetCell(0, 3, tview.NewTableCell("SIZE").SetTextColor(tcell.ColorYellow))
	v.lvTable.SetCell(0, 4, tview.NewTableCell("STATUS").SetTextColor(tcell.ColorYellow))

	lvs, err := v.manager.ListLogicalVolumes()
	if err == nil {
		for i, lv := range lvs {
			row := i + 1
			color := tcell.ColorWhite
			if lv.Status != "" {
				// Simple status check logic
				if lv.Status == "available" { // This depends on what 'lvs' command returns, usually complex attributes
					color = tcell.ColorGreen
				}
			}
			v.lvTable.SetCell(row, 0, tview.NewTableCell(lv.Name).SetTextColor(color))
			v.lvTable.SetCell(row, 1, tview.NewTableCell(lv.VGName).SetTextColor(color))
			v.lvTable.SetCell(row, 2, tview.NewTableCell(lv.Path).SetTextColor(color))
			v.lvTable.SetCell(row, 3, tview.NewTableCell(lv.Size).SetTextColor(color))
			v.lvTable.SetCell(row, 4, tview.NewTableCell(lv.Status).SetTextColor(color))
		}
	}
}
