@echo off
set "TASK_NAME=ClaudeAutoPingTask"
set "SCRIPT_DIR=%cd%"
if "%SCRIPT_DIR:~-1%"=="\" set "SCRIPT_DIR=%SCRIPT_DIR:~0,-1%"

:: Convert Windows path to WSL path (e.g. C:\foo\bar -> /mnt/c/foo/bar)
set "WSL_PATH=%SCRIPT_DIR:\=/%"
set "DRIVE_LETTER=%WSL_PATH:~0,1%"
for %%a in ("A=a" "B=b" "C=c" "D=d" "E=e" "F=f" "G=g" "H=h" "I=i" "J=j" "K=k" "L=l" "M=m" "N=n" "O=o" "P=p" "Q=q" "R=r" "S=s" "T=t" "U=u" "V=v" "W=w" "X=x" "Y=y" "Z=z") do call set "DRIVE_LETTER=%%DRIVE_LETTER:%%~a%%"
set "WSL_PATH=/mnt/%DRIVE_LETTER%%WSL_PATH:~2%"

:: Register the task to execute via WSL
schtasks /create /tn "%TASK_NAME%" /tr "wsl -d Ubuntu -e bash -c 'cd %WSL_PATH% && chmod +x claude_auto_ping.sh && ./claude_auto_ping.sh'" /sc hourly /mo 1 /st 00:00 /f

echo [ON] Claude Auto-Ping task scheduled (Every hour at :00).
pause