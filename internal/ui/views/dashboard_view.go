package views

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/rivo/tview"
)

type DashboardView struct {
	*tview.Grid
	cpuGauge  *tview.TextView
	memGauge  *tview.TextView
	swapGauge *tview.TextView
	loadText  *tview.TextView
	netText   *tview.TextView
	diskText  *tview.TextView
	procTable *tview.Table
	monitor   ports.SystemMonitor
	app       *tview.Application
	stop      chan struct{}
}

func NewDashboardView(app *tview.Application, monitor ports.SystemMonitor) *DashboardView {
	v := &DashboardView{
		Grid:      tview.NewGrid(),
		cpuGauge:  tview.NewTextView(),
		memGauge:  tview.NewTextView(),
		swapGauge: tview.NewTextView(),
		loadText:  tview.NewTextView(),
		netText:   tview.NewTextView(),
		diskText:  tview.NewTextView(),
		procTable: tview.NewTable(),
		monitor:   monitor,
		app:       app,
		stop:      make(chan struct{}),
	}

	v.setupUI()
	v.startUpdateLoop()

	return v
}

func (v *DashboardView) setupUI() {
	v.SetRows(10, 0)
	v.SetColumns(0, 0, 0, 0) // 4 equal columns
	v.SetBorders(false)

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
	v.procTable.SetTitle(" TOP PROCESSES ")
	v.procTable.SetTitleColor(tcell.ColorAqua)
	v.procTable.SetBorderColor(tcell.ColorAqua)

	// Layout
	// Row 0: Gauges/Stats
	v.AddItem(v.cpuGauge, 0, 0, 1, 1, 0, 0, false)
	v.AddItem(v.memGauge, 0, 1, 1, 1, 0, 0, false)
	v.AddItem(v.netText, 0, 2, 1, 1, 0, 0, false)
	v.AddItem(v.diskText, 0, 3, 1, 1, 0, 0, false)

	// Row 1: Processes (spanning all cols)
	v.AddItem(v.procTable, 1, 0, 1, 4, 0, 0, false)
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
	v.procTable.Clear()
	headers := []string{"PID", "USER", "CPU", "MEM", "COMMAND"}
	for i, h := range headers {
		v.procTable.SetCell(0, i, tview.NewTableCell(h).SetTextColor(tcell.ColorYellow).SetSelectable(false))
	}

	for i, p := range metrics.TopProcesses {
		row := i + 1
		v.procTable.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", p.PID)).SetTextColor(tcell.ColorWhite))
		v.procTable.SetCell(row, 1, tview.NewTableCell(p.User).SetTextColor(tcell.ColorWhite))
		v.procTable.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%.1f", p.CPU)).SetTextColor(tcell.ColorWhite))
		v.procTable.SetCell(row, 3, tview.NewTableCell(formatBytes(p.Memory)).SetTextColor(tcell.ColorWhite))
		v.procTable.SetCell(row, 4, tview.NewTableCell(p.Name).SetTextColor(tcell.ColorWhite))
	}
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
