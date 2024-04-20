#!/bin/sh

~/venv/bin/python ~/launcherapi.py &

xdotool search --sync --name "Lua Console" windowunmap &
/home/$USER/BizHawk/EmuHawkMono.sh