@echo off
chcp 949 > nul

set "TASK_NAME=TmuxKeepAliveTask"

schtasks /delete /tn "%TASK_NAME%" /f

echo [OFF] Done.
pause