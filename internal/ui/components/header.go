package components

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/ui/themes"
	"github.com/rivo/tview"
)

// Header displays system information at the top
type Header struct {
	*tview.TextView
	app *tview.Application
}

// NewHeader creates a new Header component
func NewHeader(app *tview.Application) *Header {
	h := &Header{
		TextView: tview.NewTextView(),
		app:      app,
	}

	h.SetDynamicColors(true)
	h.SetTextAlign(tview.AlignCenter)
	h.SetTextColor(tcell.ColorWhite)
	h.SetBackgroundColor(tcell.ColorBlack)
	h.SetBorder(true)
	h.SetBorderColor(tcell.ColorAqua) // Use neon cyan for border
	h.SetTitle(" TRONCLI: SYSTEM STATUS ")
	h.SetTitleColor(tcell.ColorAqua)

	// Start update loop in a goroutine
	go h.updateLoop()

	return h
}

func (h *Header) updateLoop() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		h.app.QueueUpdateDraw(func() {
			h.Update()
		})
	}
}

// Update refreshes the header content
func (h *Header) Update() {
	hostname, _ := os.Hostname()
	uptime := time.Since(startTime).Round(time.Second).String()

	text := fmt.Sprintf("%sHOST:%s %s  %sOS:%s %s  %sARCH:%s %s  %sUPTIME:%s %s  %sTIME:%s %s",
		themes.ColorNeonCyan, themes.ColorWhite, hostname,
		themes.ColorNeonCyan, themes.ColorWhite, runtime.GOOS,
		themes.ColorNeonCyan, themes.ColorWhite, runtime.GOARCH,
		themes.ColorNeonCyan, themes.ColorWhite, uptime,
		themes.ColorNeonCyan, themes.ColorWhite, time.Now().Format("15:04:05"))

	h.SetText(text)
}

var startTime = time.Now()
