FROM ubuntu:22.04

#Variables for installation
ENV DEBIAN_FRONTEND=noninteractive
ENV DEBCONF_NONINTERACTIVE_SEEN=true
ENV XKB_DEFAULT_RULES=base

ENV USER=launcher
ARG HOME=/home/$USER

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
        xdelta3

RUN useradd -ms /bin/bash $USER && \
  wget https://github.com/btimofeev/lazy_ips/archive/refs/heads/master.zip && \
  unzip master.zip && cd lazy_ips-master && python3 setup.py install

USER $USER
WORKDIR /home/$USER
RUN mkdir /home/$USER/BizHawk && python3 -m venv venv

# Install Bizhawk & Ironmon-Tracker & Randomizer & Extensions
COPY BizHawk_devbuild_45e2307a.zip $HOME/BizHawk/
RUN unzip $HOME/BizHawk/BizHawk_devbuild_45e2307a.zip && \
  tar -xvf BizHawk_devbuild_45e2307a.tar -C BizHawk && \
  chmod +x BizHawk/EmuHawkMono.sh && \
  wget https://github.com/besteon/Ironmon-Tracker/releases/download/v8.5.1/Ironmon-Tracker.zip && \
  unzip Ironmon-Tracker.zip -d $HOME/BizHawk/Lua/ && \
  wget https://github.com/Ajarmar/universal-pokemon-randomizer-zx/releases/download/v4.6.0/PokeRandoZX-v4_6_0.zip && \
  unzip PokeRandoZX-v4_6_0.zip -d $HOME/ && \
  wget https://github.com/UTDZac/ImAttached-IronmonExtension/releases/download/v1.3/ImAttachedExtension.lua -O $HOME/BizHawk/Lua/Ironmon-Tracker/extensions/ImAttachedExtension.lua && \
  wget https://github.com/jtigues/AutoPedometerTrackerExtension/releases/download/Pedometer1.1/AutoPedometer.lua -O $HOME/BizHawk/Lua/Ironmon-Tracker/extensions/AutoPedometer.lua && \
  wget https://github.com/UTDZac/CrozwordsTourney-IronmonExtension/releases/download/v3.3/CrozwordsTourneyExtension.lua -O $HOME/BizHawk/Lua/Ironmon-Tracker/extensions/CrozwordsTourneyExtension.lua && \
  wget https://github.com/UTDZac/CalcAtk-IronmonExtension/releases/download/v1.1/CalcAtk.lua -O $HOME/BizHawk/Lua/Ironmon-Tracker/extensions/CalcAtk.lua && \
  wget https://github.com/UTDZac/DeathQuotes-IronmonExtension/releases/download/v1.1/DeathQuotes.lua -O $HOME/BizHawk/Lua/Ironmon-Tracker/extensions/DeathQuotes.lua && \
  wget https://github.com/mdmurphy2/Ironmon-fastforward/releases/download/v1.2/FastForwardBattleIntroExtension.lua -O $HOME/BizHawk/Lua/Ironmon-Tracker/extensions/DeathQuotes.lua && \
  wget https://github.com/UTDZac/SpriteIsMe-IronmonExtension/releases/download/v1.4/SpriteIsMe.zip && \
  unzip SpriteIsMe.zip -d $HOME/BizHawk/Lua/Ironmon-Tracker/extensions/ && \
  wget https://github.com/CyanSMP64/NatDexExtension/releases/download/v1.1.0/NatDexExtension_v1.1.0.zip && \
  unzip NatDexExtension_v1.1.0.zip -d $HOME/BizHawk/Lua/Ironmon-Tracker/extensions/ && \
  wget https://github.com/WaffleSmacker/SmackerTracker-IronmonExtension/releases/download/v1.2/SmackerTracker.zip && \
  unzip SmackerTracker.zip -d $HOME/BizHawk/Lua/Ironmon-Tracker/extensions/ && \
  wget https://github.com/Fellshadow/Ironmon-Tracker-AutoPokemonThemes/releases/download/v1.2/AutoPokemonThemes_v1.2.zip && \
  unzip AutoPokemonThemes_v1.2.zip -d $HOME/BizHawk/Lua/Ironmon-Tracker/extensions/

COPY entrypoint_user.sh $HOME/entrypoint_user.sh
COPY launcherapi.py $HOME/launcherapi.py
  
# java -Xmx4608M -jar PokeRandoZX.jar please-use-the-launcher

RUN mkdir -p $HOME/.local/share/fonts
COPY ["ironmon-tracker/Franklin Gothic Medium Regular.ttf", "/home/launcher/.local/share/fonts/"]
RUN fc-cache -f

COPY bizhawk/config.ini $HOME/BizHawk/
COPY ironmon-tracker/Settings.ini $HOME/BizHawk/Lua/Ironmon-Tracker/
COPY roms $HOME/roms


USER root
RUN chmod 777 $HOME/entrypoint_user.sh && $HOME/venv/bin/pip install flask
COPY entrypoint_root.sh /entrypoint_root.sh

EXPOSE 5000
ENTRYPOINT [ "/bin/sh", "/entrypoint_root.sh" ]
