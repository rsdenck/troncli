package themes

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TronTheme defines the color palette inspired by TRON: Legacy
var TronTheme = tview.Theme{
	PrimitiveBackgroundColor:    tcell.ColorBlack,
	ContrastBackgroundColor:     tcell.ColorDarkBlue,
	MoreContrastBackgroundColor: tcell.Color16,   // Darker black/grey
	BorderColor:                 tcell.ColorAqua, // Neon Cyan for borders
	TitleColor:                  tcell.ColorAqua, // Neon Cyan for titles
	GraphicsColor:               tcell.ColorBlue,
	PrimaryTextColor:            tcell.ColorWhite,
	SecondaryTextColor:          tcell.ColorAqua,
	TertiaryTextColor:           tcell.ColorTeal,
	InverseTextColor:            tcell.ColorBlack,
	ContrastSecondaryTextColor:  tcell.ColorTeal,
}

// Color Markup Constants
const (
	ColorNeonBlue    = "[#0000FF]"
	ColorNeonCyan    = "[#00FFFF]"
	ColorWhite       = "[#FFFFFF]"
	ColorBlack       = "[#000000]"
	ColorAlertRed    = "[#FF0000]"
	ColorStatusGreen = "[#00FF00]"
	ColorWarning     = "[#FFFF00]"
)

// Tcell Color Constants for direct usage
const (
	TronCyan   = tcell.ColorAqua
	TronBlue   = tcell.ColorBlue
	TronYellow = tcell.ColorYellow
	TronGreen  = tcell.ColorGreen
	TronRed    = tcell.ColorRed
	TronWhite  = tcell.ColorWhite
	TronBlack  = tcell.ColorBlack
)

// ApplyTheme applies the Tron theme to the application
func ApplyTheme() {
	tview.Styles = TronTheme
}
