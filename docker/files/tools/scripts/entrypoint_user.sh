#!/bin/bash

~/venv/bin/python $SCRIPTS/launcher.py

xdotool search --sync --name "Lua Console" windowunmap &
~/BizHawk/EmuHawkMono.sh