package scheduler

import (
	"fmt"
	"strings"

	"wakeclaude/internal/app"
)

// loadOAuthToken retrieves the stored OAuth token.
// On Windows, Task Scheduler runs tasks directly as the user,
// so no privilege dropping (launchctl asuser) is needed.
func loadOAuthToken(_ ScheduleEntry) (string, error) {
	token, err := app.LoadOAuthToken()
	if err != nil {
		return "", fmt.Errorf("missing setup token; run %s", app.ClaudeSetupTokenCmd)
	}
	if strings.TrimSpace(token) == "" {
		return "", fmt.Errorf("missing setup token; run %s", app.ClaudeSetupTokenCmd)
	}
	return token, nil
}
