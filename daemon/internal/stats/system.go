package stats

import (
	"github.com/charmbracelet/log"
)

func SendNotification(title, message string, logger *log.Logger) error {
	logger.Info("notification sent",
		"title", title,
		"message", message)
	return nil
}

func SetScreenBrightness(percent int, logger *log.Logger) error {
	logger.Info("brightness changed",
		"percent", percent)
	return nil
}
