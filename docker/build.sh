#!/bin/bash

python3 ./buildscripts/buildpatchdb.py "../../roms" "../files/patches" 2
#docker build . -t ironlauncher:dev