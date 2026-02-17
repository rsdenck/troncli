package views

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/ui/themes"
	"github.com/rivo/tview"
)

type AuditView struct {
	*tview.Flex
	table   *tview.Table
	details *tview.TextView
	manager ports.AuditManager
}

func NewAuditView(manager ports.AuditManager) *AuditView {
	v := &AuditView{
		Flex:    tview.NewFlex(),
		table:   tview.NewTable(),
		details: tview.NewTextView(),
		manager: manager,
	}

	v.setupUI()
	v.loadData()

	return v
}

func (v *AuditView) setupUI() {
	v.SetDirection(tview.FlexRow)

	// Table setup
	v.table.SetBorders(true)
	v.table.SetBorderColor(tcell.ColorAqua)
	v.table.SetTitle(" AUDIT LOGS ")
	v.table.SetTitleColor(tcell.ColorAqua)
	v.table.SetSelectable(true, false)
	v.table.SetFixed(1, 1)

	// Details setup
	v.details.SetBorder(true)
	v.details.SetTitle(" DETAILS ")
	v.details.SetBorderColor(tcell.ColorAqua)
	v.details.SetDynamicColors(true)

	v.AddItem(v.table, 0, 2, true)
	v.AddItem(v.details, 0, 1, false)

	v.table.SetSelectedFunc(func(row, column int) {
		// Show details
		ref := v.table.GetCell(row, 0).GetReference()
		if entry, ok := ref.(ports.AuditEntry); ok {
			v.details.SetText(fmt.Sprintf("%sTime:%s %s\n%sUser:%s %s\n%sService:%s %s\n%sResult:%s %s\n\n%sMessage:%s\n%s",
				themes.ColorNeonCyan, themes.ColorWhite, entry.Timestamp.Format(time.RFC3339),
				themes.ColorNeonCyan, themes.ColorWhite, entry.User,
				themes.ColorNeonCyan, themes.ColorWhite, entry.Service,
				themes.ColorNeonCyan, themes.ColorWhite, entry.Result,
				themes.ColorNeonCyan, themes.ColorWhite, entry.Message,
			))
		}
	})
}

func (v *AuditView) loadData() {
	// Header
	headers := []string{"TIME", "SEVERITY", "USER", "SERVICE", "RESULT"}
	for i, h := range headers {
		v.table.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(tcell.ColorBlack).
			SetBackgroundColor(tcell.ColorAqua).
			SetSelectable(false))
	}

	// Load Auth logs
	entries, err := v.manager.GetAuthLogs(50)
	if err != nil {
		v.details.SetText(fmt.Sprintf("[red]Error loading audit logs: %v", err))
		return
	}

	for i, entry := range entries {
		row := i + 1
		color := tcell.ColorWhite
		if entry.Severity == "High" {
			color = tcell.ColorRed
		} else if entry.Result == "Fail" {
			color = tcell.ColorYellow
		}

		v.table.SetCell(row, 0, tview.NewTableCell(entry.Timestamp.Format("15:04:05")).SetTextColor(color).SetReference(entry))
		v.table.SetCell(row, 1, tview.NewTableCell(entry.Severity).SetTextColor(color))
		v.table.SetCell(row, 2, tview.NewTableCell(entry.User).SetTextColor(color))
		v.table.SetCell(row, 3, tview.NewTableCell(entry.Service).SetTextColor(color))
		v.table.SetCell(row, 4, tview.NewTableCell(entry.Result).SetTextColor(color))
	}
}
