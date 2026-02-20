package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/ui/themes"
	"github.com/rivo/tview"
)

type SecurityView struct {
	*tview.Flex
	statusText *tview.TextView
	scanTable  *tview.Table
	manager    ports.SecurityManager
}

func NewSecurityView(manager ports.SecurityManager) *SecurityView {
	v := &SecurityView{
		Flex:       tview.NewFlex(),
		statusText: tview.NewTextView(),
		scanTable:  tview.NewTable(),
		manager:    manager,
	}

	v.setupUI()
	v.refreshData()
	return v
}

func (v *SecurityView) setupUI() {
	v.SetDirection(tview.FlexRow)

	v.statusText.SetBorder(true).SetTitle(" SECURITY STATUS ").SetBorderColor(themes.TronCyan)
	v.statusText.SetDynamicColors(true)

	v.scanTable.SetBorders(true).SetTitle(" VULNERABILITY SCAN ").SetBorderColor(themes.TronRed)
	
	v.AddItem(v.statusText, 3, 1, false)
	v.AddItem(v.scanTable, 0, 1, true)
}

func (v *SecurityView) refreshData() {
	installed := v.manager.IsToolInstalled()
	status := "NOT INSTALLED"
	color := "[red]"
	if installed {
		status = "INSTALLED"
		color = "[green]"
	}

	v.statusText.SetText(fmt.Sprintf("CVE-BIN-TOOL: %s%s[white]\n\nPress 's' to run scan (Not implemented yet)", color, status))

	// Placeholder for scan results
	v.scanTable.Clear()
	v.scanTable.SetCell(0, 0, tview.NewTableCell("Scan results will appear here").SetTextColor(tcell.ColorGray))
}
