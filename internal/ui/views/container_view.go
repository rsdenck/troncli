package views

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/ui/themes"
	"github.com/rivo/tview"
)

type ContainerView struct {
	*tview.Flex
	table   *tview.Table
	details *tview.TextView
	manager ports.ContainerManager
}

func NewContainerView(manager ports.ContainerManager) *ContainerView {
	v := &ContainerView{
		Flex:    tview.NewFlex(),
		table:   tview.NewTable(),
		details: tview.NewTextView(),
		manager: manager,
	}

	v.setupUI()
	v.refreshData()
	return v
}

func (v *ContainerView) setupUI() {
	v.SetDirection(tview.FlexRow)

	v.table.SetBorders(true).SetTitle(" CONTAINERS ").SetBorderColor(themes.TronCyan)
	v.table.SetSelectable(true, false)
	
	v.details.SetBorder(true).SetTitle(" LOGS / DETAILS ").SetBorderColor(themes.TronBlue)

	v.AddItem(v.table, 0, 2, true)
	v.AddItem(v.details, 0, 1, false)

	v.table.SetSelectedFunc(func(row, column int) {
		// Actions: Stop/Start/Restart
		// TODO: Add modal for actions
	})
	
	v.table.SetSelectionChangedFunc(func(row, column int) {
		if row > 0 {
			cell := v.table.GetCell(row, 0) // ID
			if cell != nil {
				id := cell.Text
				logs, _ := v.manager.GetContainerLogs(id, 20)
				v.details.SetText(logs)
			}
		}
	})
}

func (v *ContainerView) refreshData() {
	v.table.Clear()
	headers := []string{"ID", "Names", "Image", "State", "Status"}
	for i, h := range headers {
		v.table.SetCell(0, i, tview.NewTableCell(h).SetTextColor(themes.TronYellow).SetSelectable(false))
	}

	containers, err := v.manager.ListContainers(true)
	if err != nil {
		v.details.SetText(fmt.Sprintf("Error listing containers: %v", err))
		return
	}

	for i, c := range containers {
		row := i + 1
		v.table.SetCell(row, 0, tview.NewTableCell(c.ID).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 1, tview.NewTableCell(strings.Join(c.Names, ",")).SetTextColor(themes.TronCyan))
		v.table.SetCell(row, 2, tview.NewTableCell(c.Image).SetTextColor(tcell.ColorWhite))
		
		stateColor := tcell.ColorRed
		if strings.Contains(strings.ToLower(c.State), "running") {
			stateColor = themes.TronGreen
		}
		v.table.SetCell(row, 3, tview.NewTableCell(c.State).SetTextColor(stateColor))
		v.table.SetCell(row, 4, tview.NewTableCell(c.Status).SetTextColor(tcell.ColorWhite))
	}
}
