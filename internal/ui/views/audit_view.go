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
	v.table.SetBorderColor(themes.TronCyan)
	v.table.SetTitle(" AUDIT LOGS ")
	v.table.SetTitleColor(themes.TronCyan)
	v.table.SetSelectable(true, false)
	v.table.SetFixed(1, 1)

	// Details setup
	v.details.SetBorder(true)
	v.details.SetTitle(" EVENT DETAILS ")
	v.details.SetTitleColor(themes.TronBlue)
	v.details.SetBorderColor(themes.TronBlue)
	v.details.SetDynamicColors(true)

	// Layout: Table takes 70%, Details takes 30%
	v.AddItem(v.table, 0, 7, true)
	v.AddItem(v.details, 0, 3, false)

	v.table.SetSelectionChangedFunc(func(row, column int) {
		if row > 0 {
			cell := v.table.GetCell(row, 0)
			if cell != nil {
				ref := cell.GetReference()
				if entry, ok := ref.(ports.AuditEvent); ok {
					v.updateDetails(entry)
				}
			}
		}
	})
}

func (v *AuditView) updateDetails(entry ports.AuditEvent) {
	severityColor := themes.ColorWhite
	switch entry.Severity {
	case "CRITICAL":
		severityColor = themes.ColorAlertRed
	case "WARNING":
		severityColor = themes.ColorWarn
	case "INFO":
		severityColor = themes.ColorNeonCyan
	}

	text := fmt.Sprintf("%sEvent Analysis\n", themes.ColorNeonBlue)
	text += fmt.Sprintf("%s────────────────────────────────────────\n", "[grey]")
	text += fmt.Sprintf("%sTimestamp: %s%s\n", themes.ColorNeonCyan, themes.ColorWhite, entry.Timestamp.Format(time.RFC1123))
	text += fmt.Sprintf("%sSeverity : %s%s\n", themes.ColorNeonCyan, severityColor, entry.Severity)
	text += fmt.Sprintf("%sUser     : %s%s\n", themes.ColorNeonCyan, themes.ColorWhite, entry.User)
	text += fmt.Sprintf("%sService  : %s%s\n", themes.ColorNeonCyan, themes.ColorWhite, entry.Type) // Using Type as Service/Source
	text += fmt.Sprintf("%s────────────────────────────────────────\n", "[grey]")
	text += fmt.Sprintf("%sMessage  :\n%s%s\n", themes.ColorNeonCyan, themes.ColorWhite, entry.Message)

	// Add some hypothetical context or raw data if available
	text += fmt.Sprintf("\n%sRaw Data :\n%s%v", "[grey]", "[darkgrey]", entry)

	v.details.SetText(text)
}

func (v *AuditView) loadData() {
	// Header
	headers := []string{"TIME", "SEVERITY", "USER", "SERVICE", "MESSAGE"}
	for i, h := range headers {
		v.table.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(themes.TronYellow).
			SetBackgroundColor(tcell.ColorBlack).
			SetSelectable(false).
			SetAlign(tview.AlignCenter))
	}

	// Load Auth logs
	entries, err := v.manager.AnalyzeLogins(24 * time.Hour)
	if err != nil {
		v.details.SetText(fmt.Sprintf("[red]Error loading audit logs: %v", err))
		return
	}

	for i, entry := range entries {
		row := i + 1
		color := tcell.ColorWhite
		switch entry.Severity {
		case "CRITICAL":
			color = tcell.ColorRed
		case "WARNING":
			color = tcell.ColorYellow
		}

		v.table.SetCell(row, 0, tview.NewTableCell(entry.Timestamp.Format("15:04:05")).SetTextColor(themes.TronCyan).SetReference(entry))
		v.table.SetCell(row, 1, tview.NewTableCell(entry.Severity).SetTextColor(color))
		v.table.SetCell(row, 2, tview.NewTableCell(entry.User).SetTextColor(tcell.ColorWhite))
		v.table.SetCell(row, 3, tview.NewTableCell(entry.Type).SetTextColor(themes.TronBlue))
		v.table.SetCell(row, 4, tview.NewTableCell(entry.Message).SetTextColor(tcell.ColorWhite))
	}
}
