package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mascli/troncli/internal/collectors/system"
	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/audit"
	"github.com/mascli/troncli/internal/modules/container"
	"github.com/mascli/troncli/internal/modules/disk"
	"github.com/mascli/troncli/internal/modules/lvm"
	"github.com/mascli/troncli/internal/modules/network"
	"github.com/mascli/troncli/internal/modules/process"
	"github.com/mascli/troncli/internal/modules/scheduler"
	"github.com/mascli/troncli/internal/modules/security"
	"github.com/mascli/troncli/internal/modules/service"
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
func NewApp() (*App, error) {
	themes.ApplyTheme()

	tviewApp := tview.NewApplication()

	// Initialize Services/Infra
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, err
	}

	sshClient, err := ssh.NewNativeSSHClient()
	if err != nil {
		return nil, err
	}
	lvmManager := lvm.NewLinuxLVMManager(true) // Sudo enabled by default
	auditManager := audit.NewUniversalAuditManager(executor, profile)
	systemMonitor := system.NewSystemMonitor()
	diskManager := disk.NewUniversalDiskManager(executor, profile)
	networkManager := network.NewUniversalNetworkManager(executor, profile)
	userManager := users.NewLinuxUserManager()
	processManager := process.NewUniversalProcessManager(executor, profile)
	containerManager := container.NewDockerManager()
	serviceManager := service.NewSystemdManager()
	schedulerManager := scheduler.NewLinuxSchedulerManager()
	securityManager := security.NewLinuxSecurityManager()

	header := components.NewHeader(tviewApp)
	statusBar := components.NewStatusBar()
	content := tview.NewPages()

	// Views
	dashboardView := views.NewDashboardView(tviewApp, systemMonitor, processManager)

	sshView := views.NewSSHView(tviewApp, sshClient)
	lvmView := views.NewLVMView(lvmManager)
	auditView := views.NewAuditView(auditManager)
	diskView := views.NewDiskView(diskManager)
	networkView := views.NewNetworkView(networkManager)
	usersView := views.NewUsersView(userManager)
	containerView := views.NewContainerView(containerManager)
	serviceView := views.NewServiceView(serviceManager)
	schedulerView := views.NewSchedulerView(schedulerManager)
	securityView := views.NewSecurityView(securityManager)

	content.AddPage("Dashboard", dashboardView, true, true)
	content.AddPage("SSH Manager", sshView, true, false)
	content.AddPage("LVM Manager", lvmView, true, false)
	content.AddPage("Disk & Storage", diskView, true, false)
	content.AddPage("Network", networkView, true, false)
	content.AddPage("Containers", containerView, true, false)
	content.AddPage("Services", serviceView, true, false)
	content.AddPage("Cron/Timers", schedulerView, true, false)
	content.AddPage("Users & Groups", usersView, true, false) // Sidebar key matches?
	content.AddPage("Audit", auditView, true, false)
	content.AddPage("Security", securityView, true, false)

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
