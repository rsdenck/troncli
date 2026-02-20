package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/ui/themes"
	"github.com/rivo/tview"
)

type ServiceView struct {
	*tview.Flex
	table   *tview.Table
	details *tview.TextView
	manager ports.ServiceManager
}

func NewServiceView(manager ports.ServiceManager) *ServiceView {
	v := &ServiceView{
		Flex:    tview.NewFlex(),
		table:   tview.NewTable(),
		details: tview.NewTextView(),
		manager: manager,
	}

	v.setupUI()
	v.refreshData()
	return v
}

func (v *ServiceView) setupUI() {
	v.SetDirection(tview.FlexRow)

	v.table.SetBorders(true).SetTitle(" SYSTEM SERVICES ").SetBorderColor(themes.TronCyan)
	v.table.SetSelectable(true, false)

	v.details.SetBorder(true).SetTitle(" SERVICE STATUS ").SetBorderColor(themes.TronBlue)

	v.AddItem(v.table, 0, 2, true)
	v.AddItem(v.details, 0, 1, false)

	v.table.SetSelectionChangedFunc(func(row, column int) {
		if row > 0 {
			cell := v.table.GetCell(row, 0) // Name
			if cell != nil {
				name := cell.Text
				status, _ := v.manager.GetServiceStatus(name)
				v.details.SetText(status)
			}
		}
	})
}

func (v *ServiceView) refreshData() {
	v.table.Clear()
	headers := []string{"Unit", "Load", "Active", "Sub", "Description"}
	for i, h := range headers {
		v.table.SetCell(0, i, tview.NewTableCell(h).SetTextColor(themes.TronYellow).SetSelectable(false))
	}

	services, err := v.manager.ListServices()
	if err != nil {
		v.details.SetText(fmt.Sprintf("Error listing services: %v", err))
		return
	}

	for i, s := range services {
		row := i + 1
		v.table.SetCell(row, 0, tview.NewTableCell(s.Name).SetTextColor(themes.TronCyan))
		v.table.SetCell(row, 1, tview.NewTableCell(s.LoadState).SetTextColor(tcell.ColorWhite))

		activeColor := tcell.ColorRed
		if s.ActiveState == "active" {
			activeColor = themes.TronGreen
		}
		v.table.SetCell(row, 2, tview.NewTableCell(s.ActiveState).SetTextColor(activeColor))
		v.table.SetCell(row, 3, tview.NewTableCell(s.SubState).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 4, tview.NewTableCell(s.Description).SetTextColor(tcell.ColorWhite))
	}
}
