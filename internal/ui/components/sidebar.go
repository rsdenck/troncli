package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Sidebar displays the navigation menu
type Sidebar struct {
	*tview.List
}

// NewSidebar creates a new Sidebar component
func NewSidebar(onSelected func(index int, mainText string, secondaryText string, shortcut rune)) *Sidebar {
	s := &Sidebar{
		List: tview.NewList(),
	}

	s.SetBorder(true)
	s.SetBorderColor(tcell.ColorAqua)
	s.SetTitle(" MODULES ")
	s.SetTitleColor(tcell.ColorAqua)
	s.SetBackgroundColor(tcell.ColorBlack)
	s.SetMainTextColor(tcell.ColorWhite)
	s.SetSelectedTextColor(tcell.ColorBlack)
	s.SetSelectedBackgroundColor(tcell.ColorAqua)
	s.ShowSecondaryText(false)

	menuItems := []struct {
		Text     string
		Shortcut rune
	}{
		{"Dashboard", 'd'},
		{"LVM Manager", 'l'},
		{"Disk & Storage", 's'},
		{"Network", 'n'},
		{"SSH Manager", 'h'},
		{"Logs", 'g'},
		{"Audit", 'a'},
		{"Users & Perms", 'u'},
		{"Security", 'e'},
	}

	for _, item := range menuItems {
		s.AddItem(item.Text, "", item.Shortcut, nil)
	}

	s.SetSelectedFunc(onSelected)

	return s
}
