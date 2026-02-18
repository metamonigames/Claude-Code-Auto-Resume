package scheduler

// EnsureSudo is a no-op on Windows.
// Windows Task Scheduler creates user-level tasks without requiring
// administrator privileges, so no elevation check is performed here.
func EnsureSudo() error {
	return nil
}
