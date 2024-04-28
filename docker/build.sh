#!/bin/bash

python3 ./buildscripts/buildpatchdb.py "../../roms" "../files/patches" 3
docker build . -t ironlauncher:dev