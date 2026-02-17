package views

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/rivo/tview"
)

type DashboardView struct {
	*tview.Pages
	grid *tview.Grid

	cpuGauge  *tview.TextView
	memGauge  *tview.TextView
	swapGauge *tview.TextView
	loadText  *tview.TextView
	netText   *tview.TextView
	diskText  *tview.TextView
	procTable *tview.Table
	monitor   ports.SystemMonitor
	procMgr   ports.ProcessManager
	app       *tview.Application
	stop      chan struct{}
}

func NewDashboardView(app *tview.Application, monitor ports.SystemMonitor, procMgr ports.ProcessManager) *DashboardView {
	grid := tview.NewGrid()
	pages := tview.NewPages()
	pages.AddPage("main", grid, true, true)

	v := &DashboardView{
		Pages:     pages,
		grid:      grid,
		cpuGauge:  tview.NewTextView(),
		memGauge:  tview.NewTextView(),
		swapGauge: tview.NewTextView(),
		loadText:  tview.NewTextView(),
		netText:   tview.NewTextView(),
		diskText:  tview.NewTextView(),
		procTable: tview.NewTable(),
		monitor:   monitor,
		procMgr:   procMgr,
		app:       app,
		stop:      make(chan struct{}),
	}

	v.setupUI()
	v.startUpdateLoop()

	return v
}

func (v *DashboardView) setupUI() {
	v.grid.SetRows(10, 0)
	v.grid.SetColumns(0, 0, 0, 0) // 4 equal columns
	v.grid.SetBorders(false)

	// Widgets styling
	styleWidget := func(w *tview.TextView, title string) {
		w.SetBorder(true)
		w.SetTitle(title)
		w.SetTitleColor(tcell.ColorAqua)
		w.SetBorderColor(tcell.ColorAqua)
		w.SetTextAlign(tview.AlignCenter)
		w.SetDynamicColors(true)
	}

	styleWidget(v.cpuGauge, " CPU USAGE ")
	styleWidget(v.memGauge, " MEMORY ")
	styleWidget(v.swapGauge, " SWAP ")
	styleWidget(v.loadText, " LOAD AVG ")
	styleWidget(v.netText, " NETWORK I/O ")
	styleWidget(v.diskText, " DISK I/O ")

	v.procTable.SetBorders(true)
	v.procTable.SetTitle(" TOP PROCESSES (k: kill, r: renice) ")
	v.procTable.SetTitleColor(tcell.ColorAqua)
	v.procTable.SetBorderColor(tcell.ColorAqua)
	v.procTable.SetSelectable(true, false)

	v.procTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'k':
				v.killSelectedProcess()
				return nil
			case 'r':
				v.reniceSelectedProcess()
				return nil
			}
		}
		return event
	})

	// Layout
	// Row 0: Gauges/Stats
	v.grid.AddItem(v.cpuGauge, 0, 0, 1, 1, 0, 0, false)
	v.grid.AddItem(v.memGauge, 0, 1, 1, 1, 0, 0, false)
	v.grid.AddItem(v.netText, 0, 2, 1, 1, 0, 0, false)
	v.grid.AddItem(v.diskText, 0, 3, 1, 1, 0, 0, false)

	// Row 1: Processes (spanning all cols)
	v.grid.AddItem(v.procTable, 1, 0, 1, 4, 0, 0, true) // Focused by default
}

func (v *DashboardView) startUpdateLoop() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				v.app.QueueUpdateDraw(func() {
					v.update()
				})
			case <-v.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (v *DashboardView) Stop() {
	close(v.stop)
}

func (v *DashboardView) update() {
	metrics, err := v.monitor.GetMetrics()
	if err != nil {
		v.cpuGauge.SetText(fmt.Sprintf("[red]Error: %v", err))
		return
	}

	// CPU
	v.cpuGauge.SetText(fmt.Sprintf("\n\n[neon-cyan]%.1f%%", metrics.CPUUsage))

	// Memory
	memPercent := float64(metrics.MemUsed) / float64(metrics.MemTotal) * 100
	v.memGauge.SetText(fmt.Sprintf("\n\n[neon-cyan]%.1f%%\n[white]%s / %s",
		memPercent,
		formatBytes(metrics.MemUsed),
		formatBytes(metrics.MemTotal)))

	// Network
	v.netText.SetText(fmt.Sprintf("\nRX: [green]%s/s\n[white]TX: [yellow]%s/s\n\n[grey]Total RX: %s\nTotal TX: %s",
		formatBytes(metrics.NetworkIO.RxRate),
		formatBytes(metrics.NetworkIO.TxRate),
		formatBytes(metrics.NetworkIO.RxBytes),
		formatBytes(metrics.NetworkIO.TxBytes)))

	// Disk
	v.diskText.SetText(fmt.Sprintf("\n\nRead: [green]%s\n[white]Write: [yellow]%s\n[white]IOPS: [cyan]%d",
		formatBytes(metrics.DiskIO.ReadBytes),
		formatBytes(metrics.DiskIO.WriteBytes),
		metrics.DiskIO.IOPS))

	// Update Process Table
	// Store current selection
	row, _ := v.procTable.GetSelection()

	v.procTable.Clear()
	headers := []string{"PID", "USER", "CPU", "MEM", "COMMAND"}
	for i, h := range headers {
		v.procTable.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetBackgroundColor(tcell.ColorDarkBlue))
	}

	for i, p := range metrics.TopProcesses {
		r := i + 1
		v.procTable.SetCell(r, 0, tview.NewTableCell(fmt.Sprintf("%d", p.PID)).SetTextColor(tcell.ColorWhite).SetReference(p.PID))
		v.procTable.SetCell(r, 1, tview.NewTableCell(p.User).SetTextColor(tcell.ColorWhite))
		v.procTable.SetCell(r, 2, tview.NewTableCell(fmt.Sprintf("%.1f", p.CPU)).SetTextColor(tcell.ColorWhite))
		v.procTable.SetCell(r, 3, tview.NewTableCell(formatBytes(p.Memory)).SetTextColor(tcell.ColorWhite))
		v.procTable.SetCell(r, 4, tview.NewTableCell(p.Name).SetTextColor(tcell.ColorWhite))
	}

	// Restore selection if possible
	if row > 0 && row <= len(metrics.TopProcesses) {
		v.procTable.Select(row, 0)
	}
}

func (v *DashboardView) killSelectedProcess() {
	row, _ := v.procTable.GetSelection()
	if row <= 0 {
		return
	}

	ref := v.procTable.GetCell(row, 0).GetReference()
	if ref == nil {
		return
	}
	pid := ref.(int)

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to KILL process %d?", pid)).
		AddButtons([]string{"Cancel", "SIGTERM", "SIGKILL"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Cancel" {
				v.RemovePage("modal")
				v.app.SetFocus(v.procTable)
				return
			}

			err := v.procMgr.KillProcess(pid, buttonLabel)
			if err != nil {
				v.showError(fmt.Sprintf("Failed to kill process: %v", err))
			} else {
				v.RemovePage("modal")
				v.app.SetFocus(v.procTable)
			}
		})

	v.AddPage("modal", modal, true, true)
	v.app.SetFocus(modal)
}

func (v *DashboardView) reniceSelectedProcess() {
	row, _ := v.procTable.GetSelection()
	if row <= 0 {
		return
	}

	ref := v.procTable.GetCell(row, 0).GetReference()
	if ref == nil {
		return
	}
	pid := ref.(int)

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Renice process %d\nSelect new priority:", pid)).
		AddButtons([]string{"Cancel", "-10 (High)", "0 (Normal)", "10 (Low)", "19 (Idle)"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Cancel" {
				v.RemovePage("modal")
				v.app.SetFocus(v.procTable)
				return
			}

			prio := 0
			switch buttonLabel {
			case "-10 (High)":
				prio = -10
			case "0 (Normal)":
				prio = 0
			case "10 (Low)":
				prio = 10
			case "19 (Idle)":
				prio = 19
			}

			err := v.procMgr.ReniceProcess(pid, prio)
			if err != nil {
				v.showError(fmt.Sprintf("Failed to renice process: %v", err))
			} else {
				v.RemovePage("modal")
				v.app.SetFocus(v.procTable)
			}
		})

	v.AddPage("modal", modal, true, true)
	v.app.SetFocus(modal)
}

func (v *DashboardView) showError(msg string) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			v.RemovePage("error")
			v.app.SetFocus(v.procTable)
		})
	v.AddPage("error", modal, true, true)
	v.app.SetFocus(modal)
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
