package scheduler

import (
	"fmt"
	"os"
	"path/filepath"
)


const (
	appName = "WakeClaude"
)

type Store struct {
	BaseDir      string
	SchedulesDir string
	LogsDir      string
	Schedules    string
	Logs         string
}

func DefaultStore() (*Store, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolve config directory: %w", err)
	}

	base := filepath.Join(configDir, appName)
	return &Store{
		BaseDir:      base,
		SchedulesDir: base,
		LogsDir:      filepath.Join(base, "logs"),
		Schedules:    filepath.Join(base, "schedules.json"),
		Logs:         filepath.Join(base, "logs.jsonl"),
	}, nil
}

func (s *Store) Ensure() error {
	if err := os.MkdirAll(s.SchedulesDir, 0o755); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}
	if err := os.MkdirAll(s.LogsDir, 0o755); err != nil {
		return fmt.Errorf("create logs directory: %w", err)
	}
	return nil
}
