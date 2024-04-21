#!/bin/sh

~/venv/bin/python $SCRIPTS/launcherapi.py &

xdotool search --sync --name "Lua Console" windowunmap &
~/BizHawk/EmuHawkMono.sh