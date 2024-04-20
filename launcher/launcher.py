from PySide6.QtWidgets import QVBoxLayout, QApplication, QListWidget, QListWidgetItem, QWidget, QFileDialog, QStackedWidget, QMainWindow
from PySide6.QtCore import QFile, Qt
from PySide6.QtGui import QPalette
from PySide6.QtUiTools import QUiLoader

import sys

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

    def browseFiles(self):
        folder = QFileDialog.getExistingDirectory(self, 'Select ROMs folder')
        self.myWidget.romsLineEdit.setText(folder)
        

    def updateRomsText(self, text):
        self.myWidget.romsLineEdit.setText(text)

    def play(self):
        # TODO Make an API call to the container with all required settings
        return

app = QApplication(sys.argv)
mainwindow = MainWindow()
mainwindow.setFixedWidth(600)
mainwindow.setFixedHeight(480)
#mainwindow.setWindowFlags(Qt.FramelessWindowHint)
mainwindow.show()
sys.exit(app.exec())