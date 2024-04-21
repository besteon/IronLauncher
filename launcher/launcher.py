from PySide6.QtWidgets import QVBoxLayout, QApplication, QListWidget, QListWidgetItem, QWidget, QFileDialog, QStackedWidget, QMainWindow
from PySide6.QtCore import QFile, Qt, QSize
from PySide6.QtGui import QPalette
from PySide6.QtUiTools import QUiLoader

import os
import sys
import hashlib
import requests
import subprocess

SUPPORTED_HASHES = {
    'dd5945db9b930750cb39d00c84da8571feebf417': 'Pokemon Fire Red v1.1'
}

class InstallWindow(QMainWindow):
    def __init__(self):
        super(InstallWindow, self).__init__()

        self.setWindowTitle("Iron Launcher Installer")

        layout = QVBoxLayout()

        loader = QUiLoader()
        file = QFile("installer.ui")
        file.open(QFile.ReadOnly)
        self.myWidget = loader.load(file, self)
        file.close()

        layout.addWidget(self.myWidget)

        widget = QWidget()
        widget.setLayout(layout)
        self.setCentralWidget(widget)

    def install(self):
        pass

    def installWindows(self):
        pass

    def installLinux(self):
        pass

class MainWindow(QMainWindow):
    def __init__(self):
        super(MainWindow, self).__init__()

        self.setWindowTitle("Iron Launcher")

        layout = QVBoxLayout()

        loader = QUiLoader()
        file = QFile("mainwindow.ui")
        file.open(QFile.ReadOnly)
        self.myWidget = loader.load(file, self)
        file.close()

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


        return

app = QApplication(sys.argv)
mainwindow = MainWindow()
mainwindow.setFixedWidth(600)
mainwindow.setFixedHeight(480)
#mainwindow.setWindowFlags(Qt.FramelessWindowHint)
mainwindow.show()
sys.exit(app.exec())