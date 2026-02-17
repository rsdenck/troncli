package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/collectors/system"
	"github.com/mascli/troncli/internal/modules/audit"
	"github.com/mascli/troncli/internal/modules/disk"
	"github.com/mascli/troncli/internal/modules/lvm"
	"github.com/mascli/troncli/internal/modules/network"
	"github.com/mascli/troncli/internal/modules/ssh"
	"github.com/mascli/troncli/internal/modules/users"
	"github.com/mascli/troncli/internal/ui/components"
	"github.com/mascli/troncli/internal/ui/themes"
	"github.com/mascli/troncli/internal/ui/views"
	"github.com/rivo/tview"
)

// App is the main application structure
type App struct {
	TviewApp  *tview.Application
	Header    *components.Header
	Sidebar   *components.Sidebar
	Content   *tview.Pages
	StatusBar *components.StatusBar
	Layout    *tview.Grid
}

// NewApp creates and initializes the application
func NewApp() *App {
	themes.ApplyTheme()

	tviewApp := tview.NewApplication()

	// Initialize Services/Infra
	sshClient := ssh.NewRSDSSHMClient()
	lvmManager := lvm.NewLinuxLVMManager(true) // Sudo enabled by default
	auditManager := audit.NewLinuxAuditManager()
	systemMonitor := system.NewSystemMonitor()
	diskManager := disk.NewLinuxDiskManager()
	networkManager := network.NewLinuxNetworkManager()
	userManager := users.NewLinuxUserManager()

	header := components.NewHeader(tviewApp)
	statusBar := components.NewStatusBar()
	content := tview.NewPages()

	// Views
	dashboardView := views.NewDashboardView(tviewApp, systemMonitor)

	sshView := views.NewSSHView(tviewApp, sshClient)
	lvmView := views.NewLVMView(lvmManager)
	auditView := views.NewAuditView(auditManager)
	// New Views (Placeholders until implemented properly)
	diskView := views.NewDiskView(diskManager)
	networkView := views.NewNetworkView(networkManager)
	usersView := views.NewUsersView(userManager)

	content.AddPage("Dashboard", dashboardView, true, true)
	content.AddPage("SSH Manager", sshView, true, false)
	content.AddPage("LVM Manager", lvmView, true, false)
	content.AddPage("Disk & Storage", diskView, true, false)
	content.AddPage("Network", networkView, true, false)
	content.AddPage("Users & Groups", usersView, true, false)
	content.AddPage("Audit", auditView, true, false)

	// Callback for sidebar selection
	sidebar := components.NewSidebar(func(index int, mainText string, secondaryText string, shortcut rune) {
		if content.HasPage(mainText) {
			content.SwitchToPage(mainText)
		} else {
			// Placeholder for missing pages
			placeholder := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Module " + mainText + " not implemented yet.")
			placeholder.SetBorder(true).SetTitle(" " + mainText + " ").SetBorderColor(tcell.ColorAqua)
			content.AddPage(mainText, placeholder, true, true)
		}
		statusBar.SetStatus("Selected: " + mainText)
	})

	// Layout
	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0).
		SetBorders(false).
		AddItem(header, 0, 0, 1, 2, 0, 0, false).
		AddItem(sidebar, 1, 0, 1, 1, 0, 0, true).
		AddItem(content, 1, 1, 1, 1, 0, 0, false).
		AddItem(statusBar, 2, 0, 1, 2, 0, 0, false)

	app := &App{
		TviewApp:  tviewApp,
		Header:    header,
		Sidebar:   sidebar,
		Content:   content,
		StatusBar: statusBar,
		Layout:    grid,
	}

	// Splash Screen
	splash := views.NewSplashView(tviewApp, func() {
		tviewApp.SetRoot(grid, true)
	})

	tviewApp.SetRoot(splash, true)

	return app, nil
}

// Run starts the application
func (a *App) Run() error {
	return a.TviewApp.EnableMouse(true).Run()
}
