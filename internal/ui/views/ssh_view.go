package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/ui/themes"
	"github.com/rivo/tview"
)

type SSHView struct {
	*tview.Flex
	list    *tview.List
	details *tview.TextView
	client  ports.SSHClient
	app     *tview.Application
}

func NewSSHView(app *tview.Application, client ports.SSHClient) *SSHView {
	v := &SSHView{
		Flex:    tview.NewFlex(),
		list:    tview.NewList(),
		details: tview.NewTextView(),
		client:  client,
		app:     app,
	}

	v.setupUI()
	v.loadData()

	return v
}

func (v *SSHView) setupUI() {
	v.SetDirection(tview.FlexRow)

	v.list.SetBorder(true).SetTitle(" SSH PROFILES ").SetBorderColor(tcell.ColorAqua).SetTitleColor(tcell.ColorAqua)
	v.list.ShowSecondaryText(false)
	v.list.SetSelectedBackgroundColor(tcell.ColorAqua)
	v.list.SetSelectedTextColor(tcell.ColorBlack)
	v.list.SetMainTextColor(tcell.ColorWhite)

	v.details.SetBorder(true).SetTitle(" CONNECTION INFO ").SetBorderColor(tcell.ColorAqua).SetTitleColor(tcell.ColorAqua)
	v.details.SetDynamicColors(true)

	v.AddItem(v.list, 0, 3, true)
	v.AddItem(v.details, 0, 1, false)

	v.list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// Connect logic
		v.connect(mainText)
	})

	v.list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		v.details.SetText(fmt.Sprintf("%sSelected Profile: %s%s\n\n%sPress Enter to Connect via rsd-sshm",
			themes.ColorNeonCyan, themes.ColorWhite, mainText, themes.ColorWarning))
	})
}

func (v *SSHView) loadData() {
	profiles, err := v.client.ListProfiles()
	if err != nil {
		v.details.SetText(fmt.Sprintf("[red]Error loading profiles: %v", err))
		// Add some dummy profiles if none found (mock behavior for demo)
		v.list.AddItem("demo-server-01", "", 0, nil)
		v.list.AddItem("demo-db-01", "", 0, nil)
		return
	}

	for _, p := range profiles {
		v.list.AddItem(p, "", 0, nil)
	}
}

func (v *SSHView) connect(profile string) {
	v.app.Suspend(func() {
		// Clear screen
		fmt.Print("\033[H\033[2J")
		fmt.Printf("Connecting to %s via rsd-sshm...\n", profile)

		err := v.client.Connect(profile)
		if err != nil {
			fmt.Printf("\nConnection failed: %v\nPress Enter to return...", err)
			fmt.Scanln()
		}
	})
}
