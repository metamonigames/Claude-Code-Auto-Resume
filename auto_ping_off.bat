@echo off
set "TASK_NAME=ClaudeAutoPingTask"

:: Delete the scheduled task
schtasks /delete /tn "%TASK_NAME%" /f

echo [OFF] Claude Auto-Ping task removed.
pause