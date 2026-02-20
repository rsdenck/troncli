package themes

import (
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Global Color Variables (Semantic)
var (
	ColorBorder     = tcell.ColorAqua
	ColorTitle      = tcell.ColorAqua
	ColorPrimary    = tcell.ColorBlue
	ColorSecondary  = tcell.ColorTeal
	ColorSuccess    = tcell.ColorGreen
	ColorWarning    = tcell.ColorYellow
	ColorError      = tcell.ColorRed
	ColorText       = tcell.ColorWhite
	ColorBackground = tcell.ColorBlack
)

// Legacy Aliases (mapped to semantic variables for backward compatibility if needed,
// but we will refactor to use semantic names)
var (
	TronCyan   = ColorBorder
	TronBlue   = ColorPrimary
	TronYellow = ColorWarning
	TronGreen  = ColorSuccess
	TronRed    = ColorError
	TronWhite  = ColorText
	TronBlack  = ColorBackground
)

// TronTheme defines the color palette inspired by TRON: Legacy
var TronTheme = tview.Theme{
	PrimitiveBackgroundColor:    tcell.ColorBlack,
	ContrastBackgroundColor:     tcell.ColorDarkBlue,
	MoreContrastBackgroundColor: tcell.Color16,
	BorderColor:                 tcell.ColorAqua,
	TitleColor:                  tcell.ColorAqua,
	GraphicsColor:               tcell.ColorBlue,
	PrimaryTextColor:            tcell.ColorWhite,
	SecondaryTextColor:          tcell.ColorAqua,
	TertiaryTextColor:           tcell.ColorTeal,
	InverseTextColor:            tcell.ColorBlack,
	ContrastSecondaryTextColor:  tcell.ColorTeal,
}

// CyberpunkTheme defines a vibrant, neon-heavy palette
var CyberpunkTheme = tview.Theme{
	PrimitiveBackgroundColor:    tcell.ColorBlack,
	ContrastBackgroundColor:     tcell.ColorDarkMagenta,
	MoreContrastBackgroundColor: tcell.Color17,
	BorderColor:                 tcell.ColorFuchsia, // Neon Pink
	TitleColor:                  tcell.ColorYellow,  // Neon Yellow
	GraphicsColor:               tcell.ColorAqua,    // Neon Cyan
	PrimaryTextColor:            tcell.ColorWhite,
	SecondaryTextColor:          tcell.ColorFuchsia,
	TertiaryTextColor:           tcell.ColorYellow,
	InverseTextColor:            tcell.ColorBlack,
	ContrastSecondaryTextColor:  tcell.ColorYellow,
}

// MatrixTheme defines a green-on-black palette
var MatrixTheme = tview.Theme{
	PrimitiveBackgroundColor:    tcell.ColorBlack,
	ContrastBackgroundColor:     tcell.ColorDarkGreen,
	MoreContrastBackgroundColor: tcell.Color22, // Dark Green
	BorderColor:                 tcell.ColorGreen,
	TitleColor:                  tcell.ColorGreen,
	GraphicsColor:               tcell.ColorDarkGreen,
	PrimaryTextColor:            tcell.ColorGreen,
	SecondaryTextColor:          tcell.ColorLime,
	TertiaryTextColor:           tcell.ColorDarkGreen,
	InverseTextColor:            tcell.ColorBlack,
	ContrastSecondaryTextColor:  tcell.ColorLime,
}

// Color Markup Constants (Strings for tview tags)
var (
	ColorNeonBlue    = "[#0000FF]"
	ColorNeonCyan    = "[#00FFFF]"
	ColorWhite       = "[#FFFFFF]"
	ColorBlack       = "[#000000]"
	ColorAlertRed    = "[#FF0000]"
	ColorStatusGreen = "[#00FF00]"
	ColorWarn        = "[#FFFF00]"
)

// ApplyTheme applies the theme based on TRONCLI_THEME env var
func ApplyTheme() {
	themeName := strings.ToLower(os.Getenv("TRONCLI_THEME"))

	switch themeName {
	case "cyberpunk":
		tview.Styles = CyberpunkTheme
		updateSemanticColors(CyberpunkTheme)
		// Update Markup
		ColorNeonBlue = "[#00FFFF]" // Cyan
		ColorNeonCyan = "[#FF00FF]" // Pink
		ColorWarn = "[#FFFF00]"     // Yellow
	case "matrix":
		tview.Styles = MatrixTheme
		updateSemanticColors(MatrixTheme)
		// Update Markup
		ColorNeonBlue = "[#006400]" // DarkGreen
		ColorNeonCyan = "[#00FF00]" // Green
		ColorWarn = "[#008000]"     // Green
		ColorWhite = "[#00FF00]"    // All white becomes Green
	default:
		tview.Styles = TronTheme
		updateSemanticColors(TronTheme)
	}
}

func updateSemanticColors(t tview.Theme) {
	ColorBorder = t.BorderColor
	ColorTitle = t.TitleColor
	ColorPrimary = t.GraphicsColor
	ColorSecondary = t.SecondaryTextColor
	ColorText = t.PrimaryTextColor
	ColorBackground = t.PrimitiveBackgroundColor

	// Aliases update
	TronCyan = ColorBorder
	TronBlue = ColorPrimary
	TronYellow = ColorWarning
	TronGreen = ColorSuccess
	TronRed = ColorError
	TronWhite = ColorText
	TronBlack = ColorBackground
}
