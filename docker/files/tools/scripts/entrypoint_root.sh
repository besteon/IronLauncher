#!/bin/sh

chown -R $USER:$USER /dev/snd
exec runuser -u $USER $SCRIPTS/entrypoint_user.sh