@echo off
chcp 949 > nul

set "TASK_NAME=TmuxKeepAliveTask"
set "BATCH_PATH=%~dp0tmux_keep_alive_run.bat"

schtasks /create /tn "%TASK_NAME%" /tr "\"%BATCH_PATH%\"" /sc minute /mo 60 /st 00:00 /f

echo [ON] Done.
pause