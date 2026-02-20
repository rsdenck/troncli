package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/ui/themes"
	"github.com/rivo/tview"
)

type SchedulerView struct {
	*tview.Flex
	cronTable  *tview.Table
	timerTable *tview.Table
	manager    ports.SchedulerManager
}

func NewSchedulerView(manager ports.SchedulerManager) *SchedulerView {
	v := &SchedulerView{
		Flex:       tview.NewFlex(),
		cronTable:  tview.NewTable(),
		timerTable: tview.NewTable(),
		manager:    manager,
	}

	v.setupUI()
	v.refreshData()
	return v
}

func (v *SchedulerView) setupUI() {
	v.SetDirection(tview.FlexRow)

	v.cronTable.SetBorders(true).SetTitle(" CRON JOBS ").SetBorderColor(themes.TronCyan)
	v.cronTable.SetSelectable(true, false)

	v.timerTable.SetBorders(true).SetTitle(" SYSTEMD TIMERS ").SetBorderColor(themes.TronBlue)
	v.timerTable.SetSelectable(true, false)

	v.AddItem(v.cronTable, 0, 1, true)
	v.AddItem(v.timerTable, 0, 1, false)
}

func (v *SchedulerView) refreshData() {
	// Cron
	v.cronTable.Clear()
	headers := []string{"Schedule", "Command"}
	for i, h := range headers {
		v.cronTable.SetCell(0, i, tview.NewTableCell(h).SetTextColor(themes.TronYellow).SetSelectable(false))
	}

	jobs, err := v.manager.ListCronJobs()
	if err == nil {
		for i, job := range jobs {
			row := i + 1
			v.cronTable.SetCell(row, 0, tview.NewTableCell(job.Schedule).SetTextColor(themes.TronCyan))
			v.cronTable.SetCell(row, 1, tview.NewTableCell(job.Command).SetTextColor(tcell.ColorWhite))
		}
	} else {
		v.cronTable.SetCell(1, 0, tview.NewTableCell(err.Error()).SetTextColor(tcell.ColorRed))
	}

	// Timers
	v.timerTable.Clear()
	headersTimer := []string{"Unit", "Next", "Left", "Last"}
	for i, h := range headersTimer {
		v.timerTable.SetCell(0, i, tview.NewTableCell(h).SetTextColor(themes.TronYellow).SetSelectable(false))
	}

	timers, err := v.manager.ListTimers(false)
	if err == nil {
		for i, t := range timers {
			row := i + 1
			v.timerTable.SetCell(row, 0, tview.NewTableCell(t.Unit).SetTextColor(themes.TronCyan))
			v.timerTable.SetCell(row, 1, tview.NewTableCell(t.Next).SetTextColor(tcell.ColorWhite))
			v.timerTable.SetCell(row, 2, tview.NewTableCell(t.Left).SetTextColor(tcell.ColorWhite))
			v.timerTable.SetCell(row, 3, tview.NewTableCell(t.Last).SetTextColor(tcell.ColorWhite))
		}
	} else {
		v.timerTable.SetCell(1, 0, tview.NewTableCell(err.Error()).SetTextColor(tcell.ColorRed))
	}
}
