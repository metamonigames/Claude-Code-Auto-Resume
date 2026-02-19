# Claude WSL Auto-Ping

A lightweight automation suite to keep your **Claude Code** sessions alive in WSL and automate task progression during specific hours.

## Features

* **Hourly Keep-Alive**: Prevents session timeouts by sending a "ping" every 60 minutes.
* **Time-Aware Model Switching**: Automatically switches between models based on the time of day.
* **Day Shift (08:00 - 01:59)**: Uses the cost-effective **Haiku** model for simple pings.
* **Night Shift (02:00 - 07:59)**: Uses **Sonnet** to proceed with tasks and auto-commit changes.


* **Windows Integration**: Uses the native **Windows Task Scheduler** (`schtasks`) to ensure pings run in the background without needing an open terminal.
* 
**WSL Native**: Designed to work with Claude Code installed inside a WSL (Ubuntu) environment.



## Components

* 
`claude_auto_ping.sh`: The core logic script running inside WSL that detects the time and executes the appropriate Claude command.


* `auto_ping_on.bat`: Windows batch file to register the hourly task in the Windows Task Scheduler.
* `auto_ping_off.bat`: Windows batch file to stop and remove the scheduled task.

## Setup

1. **Place Files**: Keep all three files (`.sh`, `auto_ping_on.bat`, `auto_ping_off.bat`) in the same directory within your Windows filesystem.
2. 
**Permissions**: Ensure your WSL user has execution permissions for the shell script (the `on` batch file attempts to handle this automatically).


3. **Activation**:
* Run `auto_ping_on.bat` to start the automation.
* The task will now run every 60 minutes in the background.



## Usage Logic

| Time Range | Model (via Environment Var) | Prompt | Goal |
| --- | --- | --- | --- |
| **08:00 - 01:59** | `claude-3-haiku` | "ping" | Maintain Session |
| **02:00 - 07:59** | `claude-3-7-sonnet` | "Proceed with the next task..." | Autonomous Progress |

## Requirements

* 
**WSL (Ubuntu)** with Claude Code installed and authenticated.


* **Windows 10/11** for Task Scheduler support.
