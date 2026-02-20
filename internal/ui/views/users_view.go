package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/ui/themes"
	"github.com/rivo/tview"
)

type UsersView struct {
	*tview.Flex
	table   *tview.Table
	manager ports.UserManager
}

func NewUsersView(manager ports.UserManager) *UsersView {
	v := &UsersView{
		Flex:    tview.NewFlex(),
		table:   tview.NewTable(),
		manager: manager,
	}

	v.setupUI()
	v.refreshData()
	return v
}

func (v *UsersView) setupUI() {
	v.SetDirection(tview.FlexRow)
	v.table.SetBorders(true).SetTitle(" USERS ").SetBorderColor(themes.TronCyan)
	v.table.SetSelectable(true, false)
	v.AddItem(v.table, 0, 1, true)
}

func (v *UsersView) refreshData() {
	v.table.Clear()
	headers := []string{"Username", "UID", "GID", "Shell", "Home"}
	for i, h := range headers {
		v.table.SetCell(0, i, tview.NewTableCell(h).SetTextColor(themes.TronYellow).SetSelectable(false))
	}

	users, err := v.manager.ListUsers()
	if err != nil {
		// Show error in table
		v.table.SetCell(1, 0, tview.NewTableCell(fmt.Sprintf("Error: %v", err)).SetTextColor(tcell.ColorRed))
		return
	}

	for i, u := range users {
		row := i + 1
		v.table.SetCell(row, 0, tview.NewTableCell(u.Username).SetTextColor(themes.TronCyan))
		v.table.SetCell(row, 1, tview.NewTableCell(u.UID).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 2, tview.NewTableCell(u.GID).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 3, tview.NewTableCell(u.Shell).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 4, tview.NewTableCell(u.HomeDir).SetTextColor(tcell.ColorWhite))
	}
}
