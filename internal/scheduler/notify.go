package scheduler

import (
	"os"
	"os/exec"
	"strings"
)

// notifyScript sends a Windows 10/11 toast notification using WinRT via PowerShell.
// Values are passed via environment variables to prevent injection issues,
// and XML special characters are escaped using SecurityElement.Escape.
const notifyScript = `$ErrorActionPreference = 'SilentlyContinue'
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
$t = [System.Security.SecurityElement]::Escape($env:NOTIFY_TITLE)
$b = [System.Security.SecurityElement]::Escape($env:NOTIFY_BODY)
$xml = [Windows.Data.Xml.Dom.XmlDocument]::new()
$xml.LoadXml("<toast><visual><binding template='ToastGeneric'><text>$t</text><text>$b</text></binding></visual></toast>")
$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('WakeClaude').Show($toast)`

func NotifyRun(_ ScheduleEntry, logEntry LogEntry) {
	title, body := buildNotificationContent(logEntry)

	cmd := exec.Command("powershell",
		"-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", notifyScript)
	cmd.Env = append(os.Environ(),
		"NOTIFY_TITLE="+title,
		"NOTIFY_BODY="+body,
	)
	_ = cmd.Run()
}

func buildNotificationContent(logEntry LogEntry) (title, body string) {
	title = "WakeClaude"
	subtitle := "Run complete"
	message := logEntry.PromptPreview

	if logEntry.Status != "success" {
		subtitle = "Run failed"
		if isMeaningfulError(logEntry.Error) {
			message = logEntry.Error
		}
	}

	if strings.TrimSpace(message) == "" {
		if logEntry.Status == "success" {
			message = "Run finished."
		} else {
			message = "Run failed."
		}
	}

	message = truncateNotification(message, 140)
	body = subtitle + ": " + message
	return title, body
}

func isMeaningfulError(err string) bool {
	err = strings.TrimSpace(err)
	if err == "" {
		return false
	}
	lower := strings.ToLower(err)
	if strings.HasPrefix(lower, "exit status") {
		return false
	}
	return true
}

func truncateNotification(text string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(strings.TrimSpace(text))
	if len(runes) <= max {
		return string(runes)
	}
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}
