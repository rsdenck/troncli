package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// StatusBar displays status messages at the bottom
type StatusBar struct {
	*tview.TextView
}

// NewStatusBar creates a new StatusBar component
func NewStatusBar() *StatusBar {
	s := &StatusBar{
		TextView: tview.NewTextView(),
	}

	s.SetDynamicColors(true)
	s.SetTextAlign(tview.AlignLeft)
	s.SetTextColor(tcell.ColorWhite)
	s.SetBackgroundColor(tcell.ColorBlack)
	s.SetBorder(true)
	s.SetBorderColor(tcell.ColorAqua)

	s.SetText("Ready. Press '?' for help.")

	return s
}

// SetStatus updates the status message
func (s *StatusBar) SetStatus(msg string) {
	s.SetText(msg)
}
