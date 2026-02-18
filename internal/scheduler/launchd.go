package scheduler

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

const taskFolder = "WakeClaude"

// TaskName returns the Windows Task Scheduler task name for a given schedule ID.
func TaskName(id string) string {
	return taskFolder + `\` + id
}

// LaunchdPath returns the Task Scheduler task name (kept for API compatibility).
func LaunchdPath(id string) string {
	return TaskName(id)
}

// EnsureLaunchd creates or updates a Windows Task Scheduler task for the given schedule.
func EnsureLaunchd(entry ScheduleEntry) error {
	nextRun, err := NextRun(entry, time.Now())
	if err != nil {
		return err
	}
	return createScheduledTask(entry, nextRun)
}

// RemoveLaunchd deletes the Windows Task Scheduler task for the given schedule.
func RemoveLaunchd(entry ScheduleEntry) error {
	return removeScheduledTask(entry.ID)
}

// RemoveLaunchdIfRoot removes the task unconditionally on Windows.
// (No root/admin check needed since Task Scheduler runs tasks as the user.)
func RemoveLaunchdIfRoot(entry ScheduleEntry) {
	_ = RemoveLaunchd(entry)
}

// buildCreateScript returns a PowerShell script that registers a Task Scheduler task.
// Line-continuation backticks are avoided to stay compatible with Go raw string literals.
func buildCreateScript() string {
	return `$ErrorActionPreference = 'Stop'
$id        = $env:WAKECLAUDE_TASK_ID
$binary    = $env:WAKECLAUDE_BINARY
$schedType = $env:WAKECLAUDE_SCHED_TYPE
$startTime = [DateTime]$env:WAKECLAUDE_START_TIME
$weekday   = $env:WAKECLAUDE_WEEKDAY
$folder    = $env:WAKECLAUDE_FOLDER

switch ($schedType) {
    'once'   { $trigger = New-ScheduledTaskTrigger -Once   -At $startTime }
    'daily'  { $trigger = New-ScheduledTaskTrigger -Daily  -At $startTime.ToString('HH:mm') }
    'weekly' { $trigger = New-ScheduledTaskTrigger -Weekly -DaysOfWeek $weekday -At $startTime.ToString('HH:mm') }
    default  { throw "Unknown schedule type: $schedType" }
}

$action   = New-ScheduledTaskAction -Execute $binary -Argument "--run $id"
$settings = New-ScheduledTaskSettingsSet -WakeToRun -MultipleInstances IgnoreNew -ExecutionTimeLimit ([TimeSpan]::Zero) -DisallowStartIfOnBatteries $false -StopIfGoingOnBatteries $false
$principal = New-ScheduledTaskPrincipal -UserId ([System.Security.Principal.WindowsIdentity]::GetCurrent().Name) -LogonType Interactive -RunLevel Limited

Register-ScheduledTask -TaskName "$folder\$id" -Trigger $trigger -Action $action -Settings $settings -Principal $principal -Force | Out-Null
`
}

const removeScript = `$ErrorActionPreference = 'SilentlyContinue'
Unregister-ScheduledTask -TaskName "$env:WAKECLAUDE_FOLDER\$env:WAKECLAUDE_TASK_ID" -Confirm:$false`

// createScheduledTask registers a new Windows Task Scheduler task using PowerShell.
// The task wakes the computer from sleep (WakeToRun) and runs as the current user.
func createScheduledTask(entry ScheduleEntry, nextRun time.Time) error {
	startTime := nextRun.Format("2006-01-02T15:04:05")
	weekday := strings.Title(strings.ToLower(entry.Schedule.Weekday))

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", buildCreateScript())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"WAKECLAUDE_TASK_ID="+entry.ID,
		"WAKECLAUDE_BINARY="+entry.BinaryPath,
		"WAKECLAUDE_SCHED_TYPE="+entry.Schedule.Type,
		"WAKECLAUDE_START_TIME="+startTime,
		"WAKECLAUDE_WEEKDAY="+weekday,
		"WAKECLAUDE_FOLDER="+taskFolder,
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("create scheduled task: %w", err)
	}
	return nil
}

// removeScheduledTask deletes a Task Scheduler task by ID.
func removeScheduledTask(id string) error {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", removeScript)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	cmd.Env = append(os.Environ(),
		"WAKECLAUDE_TASK_ID="+id,
		"WAKECLAUDE_FOLDER="+taskFolder,
	)
	return cmd.Run()
}
