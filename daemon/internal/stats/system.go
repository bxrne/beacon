package stats

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/log"
)

func SendNotification(title, message string, logger *log.Logger) error {
	logger.Info("sending notification",
		"title", title,
		"message", message,
		"os", runtime.GOOS)

	switch runtime.GOOS {
	case "darwin":
		script := fmt.Sprintf(`display dialog "%s" with title "%s" buttons {"OK"} default button "OK" with icon caution`,
			message, title)
		cmd := exec.Command("osascript", "-e", script)
		return cmd.Run()

	case "linux":
		cmd := exec.Command("zenity", "--warning",
			"--title", title,
			"--text", message,
			"--width", "300")
		return cmd.Run()

	case "windows":
		script := fmt.Sprintf(`Add-Type -AssemblyName PresentationFramework;[System.Windows.MessageBox]::Show('%s','%s','OK','Warning')`,
			message, title)
		cmd := exec.Command("powershell", "-Command", script)
		return cmd.Run()

	default:
		return fmt.Errorf("notifications not supported on %s", runtime.GOOS)
	}
}
