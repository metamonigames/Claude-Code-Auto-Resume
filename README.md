# wakeclaude for Windows

A Windows port of [wakeclaude](https://github.com/rittikbasu/wakeclaude) — a TUI to schedule Claude Code prompts so your sessions keep running even when you hit rate limits or go to sleep.

All features of the original macOS version have been ported to Windows, replacing macOS-specific system APIs with native Windows equivalents.

## why it exists

When you hit the 5-hour session rate limit on your Claude plan, work stops mid-flow. This tool lets you schedule a prompt to auto-resume right when the limit resets — or queue up work to run while you sleep.

## what it does

- Pick a project, pick a session (or start a new one)
- Write the prompt
- Choose a model + permission mode
- Schedule it (one-time, daily, or weekly)
- Wakes your PC from sleep only when needed and runs the prompt
- Keeps logs + shows a simple run history
- Sends a native Windows toast notification on success/error

## how it works (Windows)

| Feature | macOS (original) | Windows (this port) |
|---|---|---|
| **Scheduling daemon** | `launchd` (LaunchDaemons) | Windows Task Scheduler (`schtasks`) |
| **Wake from sleep** | `pmset schedule wakeorpoweron` | Task Scheduler `WakeToRun` setting |
| **Token storage** | macOS Keychain (`security`) | DPAPI-encrypted XML file (`%APPDATA%\WakeClaude\token.xml`) |
| **Notifications** | AppleScript (`osascript`) | WinRT Toast via PowerShell |
| **Data directory** | `~/Library/Application Support/WakeClaude` | `%APPDATA%\WakeClaude` |
| **Privilege model** | Runs as root, drops to user via `launchctl asuser` | Task Scheduler runs directly as the logged-in user |

> **Note on wake from sleep:** Windows Task Scheduler can wake the PC if "Allow wake timers" is enabled in your active Power Plan (`Control Panel → Power Options → Change plan settings → Change advanced power settings → Sleep → Allow wake timers → Enable`).

## quickstart

1. Build:

   ```powershell
   go build -o wakeclaude.exe ./cmd/wakeclaude
   ```

2. Run:

   ```powershell
   .\wakeclaude.exe
   ```

## setup token (required)

wakeclaude uses a long-lived Claude Code token so scheduled prompts keep working even after you close the terminal.

Generate one in a separate terminal:

```bash
claude setup-token
```

Paste it into wakeclaude when prompted. It is stored in `%APPDATA%\WakeClaude\token.xml`, encrypted with Windows DPAPI (only readable by the same Windows user on the same machine).

## usage (TUI)

You'll see a simple menu:

- **Schedule a prompt** (project → session → prompt → model → permission → time)
- **Manage scheduled prompts** (edit/delete)
- **View run logs**

Controls:

- Arrow keys to move, `Enter` to select
- Type to search (projects, sessions, schedules, logs)
- `Esc` to go back, `q` to quit
- Prompt entry: `Ctrl+D` to continue

## models + permission modes

Models:
- `opus`
- `sonnet`
- `haiku`

Permission modes:
- `acceptEdits` – auto-accept file edits + filesystem access
- `plan` – read-only, no commands or file changes
- `bypassPermissions` – skips permission checks (use with care)

## logs + notifications

Data lives here:

- `%APPDATA%\WakeClaude\schedules.json`
- `%APPDATA%\WakeClaude\logs.jsonl`
- `%APPDATA%\WakeClaude\logs\*.log`

Run logs are retained (last 50) and shown in the TUI. Each run also sends a native Windows toast notification via PowerShell/WinRT.

## flags

- `--projects-root <path>`: override default `~/.claude/projects`
- `--run <id>`: internal (used by Task Scheduler)

## requirements

- Windows 10 or later (WinRT Toast notifications require Windows 10+)
- PowerShell 5.1+ (included with Windows 10)
- [Claude Code](https://claude.ai/download) installed and accessible in PATH
- Claude Code sessions under `~/.claude/projects`

## assumptions

- Claude Code sessions live under `~/.claude/projects`
- Root-level `.jsonl` files are sessions
- Project display names are derived from the session `cwd` when available
- `claude` (or `claude.exe` / `claude.cmd`) must be in your PATH at scheduling time

## contributing

Open a PR if you want — bug fixes, UX polish, or smarter scheduling ideas are welcome.

Original macOS version: https://github.com/rittikbasu/wakeclaude
