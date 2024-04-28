FROM ubuntu:22.04

#Variables for installation
ENV DEBIAN_FRONTEND=noninteractive
ENV DEBCONF_NONINTERACTIVE_SEEN=true
ENV XKB_DEFAULT_RULES=base

#Install dependencies
RUN apt-get update && \
        apt-get install -y \
        unzip \
        gnupg \
        apt-transport-https \
        dbus-x11 \
        software-properties-common \
        libxv1 \
        libglu1-mesa \
        xauth \
        x11-utils \
        xorg \
        x11-xkb-utils \
        software-properties-common \
        bzip2 \
        libgtk2.0-0 \
        libncursesw5 \
        libopenal1 \
        libsdl-image1.2 \
        libsdl-ttf2.0-0 \
        libsdl1.2debian \
        libsndfile1 \
        supervisor \
        ucspi-tcp \
        wget \
        build-essential \
        ccache \
        mono-complete \
        liblua5.4-dev \
        libopenal-dev \
        libcanberra-gtk-module \
        alsa-utils \
        default-jre \
        xdotool \
        python3 \
        python3-pip \
        python3.10-venv \
        libcairo2-dev \
        curl \
        xdelta3 \
        pulseaudio \
        chromium-browser \
        nano

ENV USER=launcher
ENV HOME=/home/$USER
ARG BIZHAWK_FILES=$HOME/files/bizhawk
ARG PATCHES=$HOME/files/patches
ARG TRACKER=$HOME/files/ironmon-tracker
ARG TOOLS=$HOME/files/tools
ENV BINARIES=$TOOLS/binaries
ENV SCRIPTS=$TOOLS/scripts

ADD files $HOME/files

RUN useradd -ms /bin/bash $USER && chown -R $USER $HOME && \
  mkdir /roms  && mkdir /data && mkdir /data/saves && \
  chown $USER /roms && chown -R $USER /data && chmod +x $BINARIES/flips && \
  wget https://github.com/btimofeev/lazy_ips/archive/refs/heads/master.zip && \
  unzip master.zip && cd lazy_ips-master && python3 setup.py install

USER $USER
WORKDIR $HOME

# Install Bizhawk & Ironmon-Tracker & Randomizer & Extensions
RUN mkdir -p $HOME/.local/share/fonts && mv $TRACKER/Franklin\ Gothic\ Medium\ Regular.ttf $HOME/.local/share/fonts && fc-cache -f && \
  python3 -m venv venv && \
  echo '[general]\ndrivers = pulse' >> ~/.alsoftrc && \
  mkdir /home/$USER/BizHawk && unzip $BIZHAWK_FILES/BizHawk_devbuild_45e2307a.zip && \
  tar -xvf BizHawk_devbuild_45e2307a.tar -C BizHawk && \
  chmod +x BizHawk/EmuHawkMono.sh && \
  wget https://github.com/besteon/Ironmon-Tracker/releases/download/v8.5.1/Ironmon-Tracker.zip && \
  wget https://github.com/Brian0255/NDS-Ironmon-Tracker/releases/download/6.2.4/NDS-Ironmon-Tracker-6.2.4.zip && \
  mkdir $HOME/BizHawk/Lua/gba && \
  mkdir $HOME/BizHawk/Lua/nds && \
  unzip Ironmon-Tracker.zip -d $HOME/BizHawk/Lua/gba && \
  unzip NDS-Ironmon-Tracker-6.2.4.zip -d $HOME/BizHawk/Lua/nds && \
  wget https://github.com/Ajarmar/universal-pokemon-randomizer-zx/releases/download/v4.6.0/PokeRandoZX-v4_6_0.zip && \
  unzip PokeRandoZX-v4_6_0.zip -d $HOME/ && \
  wget https://github.com/UTDZac/ImAttached-IronmonExtension/releases/download/v1.3/ImAttachedExtension.lua -O $HOME/BizHawk/Lua/gba/Ironmon-Tracker/extensions/ImAttachedExtension.lua && \
  wget https://github.com/jtigues/AutoPedometerTrackerExtension/releases/download/Pedometer1.1/AutoPedometer.lua -O $HOME/BizHawk/Lua/gba/Ironmon-Tracker/extensions/AutoPedometer.lua && \
  wget https://github.com/UTDZac/CrozwordsTourney-IronmonExtension/releases/download/v3.3/CrozwordsTourneyExtension.lua -O $HOME/BizHawk/Lua/gba/Ironmon-Tracker/extensions/CrozwordsTourneyExtension.lua && \
  wget https://github.com/UTDZac/CalcAtk-IronmonExtension/releases/download/v1.1/CalcAtk.lua -O $HOME/BizHawk/Lua/gba/Ironmon-Tracker/extensions/CalcAtk.lua && \
  wget https://github.com/UTDZac/DeathQuotes-IronmonExtension/releases/download/v1.1/DeathQuotes.lua -O $HOME/BizHawk/Lua/gba/Ironmon-Tracker/extensions/DeathQuotes.lua && \
  wget https://github.com/mdmurphy2/Ironmon-fastforward/releases/download/v1.2/FastForwardBattleIntroExtension.lua -O $HOME/BizHawk/Lua/gba/Ironmon-Tracker/extensions/DeathQuotes.lua && \
  wget https://github.com/UTDZac/SpriteIsMe-IronmonExtension/releases/download/v1.4/SpriteIsMe.zip && \
  unzip SpriteIsMe.zip -d $HOME/BizHawk/Lua/gba/Ironmon-Tracker/extensions/ && \
  wget https://github.com/CyanSMP64/NatDexExtension/releases/download/v1.1.2/NatDexExtension_v1.1.2.zip && \
  unzip NatDexExtension_v1.1.2.zip -d $HOME/BizHawk/Lua/gba/Ironmon-Tracker/extensions/ && \
  wget https://github.com/WaffleSmacker/SmackerTracker-IronmonExtension/releases/download/v1.2/SmackerTracker.zip && \
  unzip SmackerTracker.zip -d $HOME/BizHawk/Lua/gba/Ironmon-Tracker/extensions/ && \
  wget https://github.com/Fellshadow/Ironmon-Tracker-AutoPokemonThemes/releases/download/v1.2/AutoPokemonThemes_v1.2.zip && \
  unzip AutoPokemonThemes_v1.2.zip -d $HOME/BizHawk/Lua/gba/Ironmon-Tracker/extensions/ && \
  mv $BIZHAWK_FILES/config.ini $HOME/Bizhawk.ini && \
  mv $TRACKER/gba/Settings.ini $HOME/GbaSettings.ini && \
  mv $TRACKER/nds/Settings.ini $HOME/NdsSettings.ini

USER root
RUN chmod 777 $SCRIPTS/entrypoint_user.sh && $HOME/venv/bin/pip install networkx

EXPOSE 5000

ENTRYPOINT [ "/bin/sh", "/home/launcher/files/tools/scripts/entrypoint_root.sh" ]
