#!/bin/sh

chown -R $USER:$USER /dev/snd
exec runuser -u $USER /home/$USER/entrypoint_user.sh
