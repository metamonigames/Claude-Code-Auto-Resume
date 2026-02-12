@echo off
set "WIN_PATH=%~dp0"
if "%WIN_PATH:~-1%"=="\" set "WIN_PATH=%WIN_PATH:~0,-1%"

wsl -d Ubuntu -e bash -c "cd \"$(wslpath '%WIN_PATH%')\" && chmod +x tmux_keep_alive.sh && bash tmux_keep_alive.sh"