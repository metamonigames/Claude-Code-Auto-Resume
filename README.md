# Claude WSL Auto-Ping

A lightweight automation suite to keep your **Claude Code** sessions alive in WSL and automate task progression during off-hours.

## Features

* **Hourly Keep-Alive**: Prevents session timeouts by sending a "ping" every 60 minutes.
* **Time-Aware Model Switching**: Automatically switches between models based on the time of day.
  * **Day Shift (08:00 - 01:59)**: Uses cost-effective **Haiku** for a fire-and-forget ping (sends and exits immediately).
  * **Night Shift (02:00 - 07:59)**: Uses **Sonnet** to proceed with tasks and auto-commit changes.
* **Windows Integration**: Uses the native **Windows Task Scheduler** (`schtasks`) to run in the background without an open terminal.
* **WSL Native**: Designed for Claude Code installed inside a WSL (Ubuntu) environment.

## Components

| File | Description |
| --- | --- |
| `claude_auto_ping.sh` | Core logic — detects time of day and runs the appropriate Claude command. |
| `auto_ping_on.bat` | Registers the hourly scheduled task in Windows Task Scheduler. |
| `auto_ping_off.bat` | Removes the scheduled task. |
| `tmux_launch-team.sh` | Launches a Claude Code agent-team session inside tmux. |
| `tmux_launch-team.bat` | Windows wrapper to start the tmux team session via WSL. |

## Setup

1. **Place Files**: Keep all files in the same directory within your Windows filesystem.
2. **Permissions**: Ensure your WSL user has execution permissions for the shell scripts (the `on` batch file handles this automatically).
3. **Activation**: Run `auto_ping_on.bat` to start the automation. The task runs every 60 minutes in the background.
4. **Deactivation**: Run `auto_ping_off.bat` to remove the scheduled task.

## Usage Logic

| Time Range | Model | Prompt | Behavior |
| --- | --- | --- | --- |
| **08:00 - 01:59** | `claude-haiku-4-5-20251001` | "ping" | Fire-and-forget (exits immediately) |
| **02:00 - 07:59** | `claude-sonnet-4-6` | "Proceed with the next task..." | Waits for completion |

## Requirements

* **WSL (Ubuntu)** with Claude Code installed and authenticated.
* **Windows 10/11** for Task Scheduler support.
* **tmux** (only for the team launcher feature).
