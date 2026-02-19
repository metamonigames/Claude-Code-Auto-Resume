@echo off
set "TASK_NAME=ClaudeAutoPingTask"
set "WIN_PATH=%~dp0"
if "%WIN_PATH:~-1%"=="\" set "WIN_PATH=%WIN_PATH:~0,-1%"

:: Register the task to execute via WSL
schtasks /create /tn "%TASK_NAME%" /tr "wsl -d Ubuntu -e bash -c \"cd \\\"$(wslpath '%WIN_PATH%')\\\" && chmod +x claude_auto_ping.sh && ./claude_auto_ping.sh\"" /sc minute /mo 60 /f

echo [ON] Claude Auto-Ping task scheduled (Every 60 mins).
pause