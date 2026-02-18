package scheduler

// On Windows, wake-from-sleep scheduling is handled by the Task Scheduler
// WakeToRun setting embedded in the task XML (see launchd.go).
// These functions are no-ops maintained for API compatibility with main.go.

func ScheduleWake(_ ScheduleEntry, _ string) error {
	return nil
}

func CancelWake(_ ScheduleEntry) error {
	return nil
}
