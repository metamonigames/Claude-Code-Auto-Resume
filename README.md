🤖 Claude Code Auto-Resume Guide
- This guide will help you automatically resume Claude Code whenever you see the "You've hit your limit" message. It checks your status every hour and automatically presses Enter to keep your work moving.

---
### 1. File Setup
Place all 6 files you received into a single folder of your choice.

```
tmux_keep_alive.sh
tmux_keep_alive_run.bat
tmux_keep_alive_on.bat
tmux_keep_alive_off.bat
tmux_launch-team.sh
tmux_launch-team.bat
```

### 2. No Manual Setup Required
Usually, complex "permission settings" are required for these tools. However, I have automated this process within the files, so you don't need to type any commands.

### 3. How to Start (One-time Setup)
1. In your folder, find the file named `tmux_keep_alive_on.bat`.
2. Right-click it and select "Run as administrator".
3. Once the black window shows `[ON] Done.`, you can press any key to close it.
 - Now, your computer will automatically check for limits every hour on the hour (e.g., 9:00, 10:00).

### 4. Running Claude Code
 - From now on, simply double-click `tmux_launch-team.bat` to start your Claude Code project.

---

### 💡 Quick Tips
 - What if I hit a limit? Just leave it alone. The system will detect the limit and "Double Enter" for you automatically at the top of the next hour.
 - How to stop the automation? If you no longer want the hourly checks, right-click `tmux_keep_alive_off.bat` and "Run as administrator".
 - Keep the files together: Do not move the files to different folders; they must stay together to work correctly.