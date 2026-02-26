package agent

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

// NotificationLevel defines notification urgency
type NotificationLevel int

const (
	NotificationInfo NotificationLevel = iota
	NotificationWarning
	NotificationCritical
)

// Notification represents a system notification
type Notification struct {
	Level     NotificationLevel
	Title     string
	Message   string
	Timestamp time.Time
	Source    string
}

// NotificationManager handles system notifications
type NotificationManager struct {
	Notifications []Notification
	LogFile       string
	Enabled       bool
}

func NewNotificationManager(logFile string) *NotificationManager {
	return &NotificationManager{
		Notifications: make([]Notification, 0),
		LogFile:       logFile,
		Enabled:       true,
	}
}

// Notify sends a notification to sysadmin
func (nm *NotificationManager) Notify(level NotificationLevel, title, message, source string) {
	if !nm.Enabled {
		return
	}

	notification := Notification{
		Level:     level,
		Title:     title,
		Message:   message,
		Timestamp: time.Now(),
		Source:    source,
	}

	nm.Notifications = append(nm.Notifications, notification)
	nm.displayNotification(notification)
	nm.logNotification(notification)
}

// displayNotification shows notification to console
func (nm *NotificationManager) displayNotification(n Notification) {
	var icon, color string

	switch n.Level {
	case NotificationInfo:
		icon = "ℹ️"
		color = "\033[36m" // Cyan
	case NotificationWarning:
		icon = "⚠️"
		color = "\033[33m" // Yellow
	case NotificationCritical:
		icon = "🚨"
		color = "\033[31m" // Red
	}

	fmt.Printf("\n%s%s %s\033[0m\n", color, icon, n.Title)
	fmt.Printf("%s%s\033[0m\n", color, n.Message)
	fmt.Printf("📅 %s | 🔧 %s\n\n", n.Timestamp.Format("2006-01-02 15:04:05"), n.Source)
}

// logNotification logs notification to file
func (nm *NotificationManager) logNotification(n Notification) {
	if nm.LogFile == "" {
		return
	}

	levelStr := []string{"INFO", "WARNING", "CRITICAL"}[n.Level]
	logEntry := fmt.Sprintf("[%s] %s - %s | Source: %s\n",
		n.Timestamp.Format("2006-01-02 15:04:05"),
		levelStr,
		n.Message,
		n.Source)

	file, err := os.OpenFile(nm.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("Failed to open notification log file", "error", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(logEntry); err != nil {
		slog.Error("Failed to write notification log", "error", err)
	}
}

// GetRecentNotifications returns recent notifications
func (nm *NotificationManager) GetRecentNotifications(since time.Duration) []Notification {
	var recent []Notification
	cutoff := time.Now().Add(-since)

	for _, notif := range nm.Notifications {
		if notif.Timestamp.After(cutoff) {
			recent = append(recent, notif)
		}
	}

	return recent
}

// Clear clears old notifications
func (nm *NotificationManager) Clear(olderThan time.Duration) {
	var filtered []Notification
	cutoff := time.Now().Add(-olderThan)

	for _, notif := range nm.Notifications {
		if notif.Timestamp.After(cutoff) {
			filtered = append(filtered, notif)
		}
	}

	nm.Notifications = filtered
}

// PreExecutionWarning warns about upcoming execution
func (nm *NotificationManager) PreExecutionWarning(intent, command string) {
	nm.Notify(
		NotificationWarning,
		"AI Agent - Pre-Execution Warning",
		fmt.Sprintf("Intent: %s\nCommand: %s\n\nThe AI agent is about to execute this command.", intent, command),
		"TRONCLI-Agent",
	)
}

// PostExecutionNotification notifies about completed execution
func (nm *NotificationManager) PostExecutionNotification(intent, command, result string, err error) {
	if err != nil {
		nm.Notify(
			NotificationCritical,
			"AI Agent - Execution Failed",
			fmt.Sprintf("Intent: %s\nCommand: %s\nError: %s", intent, command, err.Error()),
			"TRONCLI-Agent",
		)
	} else {
		nm.Notify(
			NotificationInfo,
			"AI Agent - Execution Success",
			fmt.Sprintf("Intent: %s\nCommand: %s\nResult: %s", intent, command, result),
			"TRONCLI-Agent",
		)
	}
}

// SystemStatusNotification sends system status updates
func (nm *NotificationManager) SystemStatusNotification(status string) {
	nm.Notify(
		NotificationInfo,
		"TRONCLI - System Status",
		status,
		"TRONCLI-Core",
	)
}

// SecurityAlert sends security-related notifications
func (nm *NotificationManager) SecurityAlert(message string) {
	nm.Notify(
		NotificationCritical,
		"🔒 Security Alert",
		message,
		"TRONCLI-Security",
	)
}

// PolicyViolationNotification sends policy violation alerts
func (nm *NotificationManager) PolicyViolationNotification(command, reason string) {
	nm.Notify(
		NotificationCritical,
		"🚫 Policy Violation",
		fmt.Sprintf("Command: %s\nReason: %s\n\nThis action was blocked by the policy engine.", command, reason),
		"TRONCLI-Policy",
	)
}
