from PySide6.QtWidgets import QVBoxLayout, QApplication, QListWidget, QListWidgetItem, QWidget, QFileDialog, QStackedWidget, QMainWindow
from PySide6.QtCore import QFile, Qt, QSize
from PySide6.QtGui import QPalette, QPixmap
from PySide6.QtUiTools import QUiLoader

import os
import sys
import hashlib
import socket
from shutil import which
import subprocess
import psutil
import time

BASEDIR = os.path.dirname(__file__)

APIKEY = "asdf"

SUPPORTED_HASHES = {
    'dd5945db9b930750cb39d00c84da8571feebf417': 'Pokemon Fire Red v1.1'
}

class InstallWindow(QMainWindow):
    def __init__(self):
        super(InstallWindow, self).__init__()

        self.setWindowTitle("Iron Launcher Installer")

        layout = QVBoxLayout()

        loader = QUiLoader()
        file = QFile(os.path.join(BASEDIR, "installer.ui"))
        file.open(QFile.ReadOnly)
        self.myWidget = loader.load(file, self)
        file.close()

        self.myWidget.installButton.clicked.connect(self.install)
        layout.addWidget(self.myWidget)

        widget = QWidget()
        widget.setLayout(layout)
        self.setCentralWidget(widget)

    def install(self):
        if os.name == 'nt':
            os.system('powershell -c "wget https://github.com/containers/podman/releases/download/v5.0.2/podman-5.0.2-setup.exe; podman-5.0.2-setup.exe /install /passive /quiet; rm podman-5.0.2-setup.exe"')
        elif os.name == 'posix':
            os.system('SUDO_ASKPASS=/usr/bin/ssh-askpass sudo -A apt install -y podman pulseaudio xwayland && Xwayland')

        MainWindow().show()
        self.destroy()

class MainWindow(QMainWindow):
    def __init__(self):
        super(MainWindow, self).__init__()

        self.setWindowTitle("Iron Launcher")

        layout = QVBoxLayout()

        loader = QUiLoader()
        file = QFile(os.path.join(BASEDIR, "mainwindow.ui"))
        file.open(QFile.ReadOnly)
        self.myWidget = loader.load(file, self)
        file.close()

        ironmonImg = QPixmap(os.path.join(BASEDIR, "ironmon.png"))
        self.myWidget.ironmonImgLabel.setPixmap(ironmonImg)

        self.myWidget.romsButton.clicked.connect(self.browseFiles)
        self.myWidget.romsLineEdit.textChanged.connect(self.updateRomsText)
        self.myWidget.playButton.clicked.connect(self.play)
        layout.addWidget(self.myWidget)

        widget = QWidget()
        widget.setLayout(layout)
        self.setCentralWidget(widget)

        palette = QPalette()
        self.setPalette(palette)

        self.launcher = { }

    def browseFiles(self):
        self.launcher['roms'] = { }

        folder = QFileDialog.getExistingDirectory(self, 'Select ROMs folder')
        self.myWidget.romsLineEdit.setText(folder)
        
        self.myWidget.romsList.clear()
        for file in os.listdir(self.myWidget.romsLineEdit.text()):
            if file.endswith('.gba') or file.endswith('.gb') or file.endswith('.gbc') or file.endswith('.nds'):
                with open(self.myWidget.romsLineEdit.text() + '/' + file, mode='rb') as f:
                    sha1 = hashlib.sha1(f.read()).hexdigest()
                    if sha1 in SUPPORTED_HASHES:
                        item = QListWidgetItem(SUPPORTED_HASHES[sha1])
                        size = QSize()
                        size.setHeight(25)
                        size.setWidth(150)
                        item.setSizeHint(size)
                        self.myWidget.romsList.addItem(item)

                        self.launcher['roms'][SUPPORTED_HASHES[sha1]] = file


    def updateRomsText(self, text):
        self.myWidget.romsLineEdit.setText(text)

    def play(self):
        # TODO Make an API call to the container with all required settings
        if 'roms' in self.launcher and len(self.launcher['roms']) > 0:
            if self.myWidget.romsList.currentItem().text():
                romFile = self.launcher['roms'][self.myWidget.romsList.currentItem().text()]
                print(romFile)
                
                self.startContainer()
                self.close()

                healthcmd = 'podman ps'
                if os.name == 'nt':
                    healthcmd = 'wsl ' + healthcmd
                
                # Wait for the container to start up
                while 'ironlauncher' not in os.popen(healthcmd).read():
                    time.sleep(1.0)
                
                # Wait for the container to stop
                while 'ironlauncher' in os.popen(healthcmd).read():
                    time.sleep(1.0)
                self.show()

    def startContainer(self):
        if os.name == 'nt':
            if 'vcxsrv' not in (p.name() for p in psutil.process_iter()):
                subprocess.Popen(f"{BASEDIR}/vcxsrv/xlaunch.exe -run {BASEDIR}/vcxsrv/config.xlaunch".split())

            path = self.myWidget.romsLineEdit.text()
            roms_path = os.popen(f'wsl wslpath -a "{path}"').read()
            ip = socket.gethostbyname(socket.gethostname())
            cmd = f'''
                wsl podman run \
                -e APIKEY={APIKEY} \
                -e DISPLAY={ip}:0 \
                -e PULSE_SERVER=/mnt/wslg/PulseServer \
                -v /mnt/wslg/:/mnt/wslg/ \
                -v {roms_path}:/roms \
                --net=host
                docker.io/besteon/ironlauncher:latest
            '''
            proc = subprocess.Popen(cmd.split())
        elif os.name == 'posix':
            roms_path = self.myWidget.romsLineEdit.text()
            display = os.environ['DISPLAY']
            xdg = os.environ['XDG_RUNTIME_DIR']
            home = os.path.expanduser('~')
            os.system('xhost +')
            cmd = f'''
                podman run \
                -e APIKEY={APIKEY} \
                -e DISPLAY={display} \
                -e PULSE_SERVER=unix:{xdg}/pulse/native \
                -v {xdg}/pulse/native:{xdg}/pulse/native \
                -v {home}/.config/pulse/cookie:/root/.config/pulse/cookie \
                -v /tmp/.X11-unix:/tmp/.X11-unix:ro \
                -v {roms_path}:/roms \
                --net=host \
                docker.io/besteon/ironlauncher:latest
            '''
            proc = subprocess.Popen(cmd.split())

def depsInstalled():        
    if os.name == 'nt': # Windows
        if which('podman') is not None:
            return True
        else:
            return False
    elif os.name == 'posix': # Linux
        if which('podman') is not None and which('pulseaudio') is not None and os.environ['DISPLAY'] is not None:
            return True
        else:
            return False

def updateContainer():
    cmd = f'podman pull docker.io/besteon/ironlauncher:latest'
    if os.name == 'nt':
        cmd = 'wsl ' + cmd
    os.system(cmd)


app = QApplication(sys.argv)
window = MainWindow() if depsInstalled() else InstallWindow()
window.setFixedWidth(600)
window.setFixedHeight(480)
window.show()
sys.exit(app.exec())