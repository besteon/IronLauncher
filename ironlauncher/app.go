package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/ini.v1"
)

type SuppportedMode struct {
	Name               string
	RomHash            string
	Randomizer         string
	RandomizerSettings string
}

type SupportedGame struct {
	Name  string
	Modes []SuppportedMode
}

var SUPPORTED_HASHES = map[string]string{
	"dd5945db9b930750cb39d00c84da8571feebf417": "Pokemon Fire Red v1.1",
	"f3ae088181bf583e55daf962a92bb46f4f1d07b7": "Pokemon Emerald",
	"007d061e1abc8d9b56c6378c82fcfb3fc990adf3": "Pokemon Heart Gold",
	"d6bf0cc22ab1619335e5ae1ca0586180054147a9": "Pokemon Platinum",
}

var settings_file string = "ironlauncher.ini"
var datapath string = ""
var savespath string = ""
var gbaSettingsPath string = ""
var ndsSettingsPath string = ""

var CONTAINER string = "docker.io/besteon/ironlauncher:latest"

//var CONTAINER string = "besteon/ironlauncher:dev"

type Settings struct {
	RomsFolder  string `json:"romsFolder"`
	DefaultRom  string `json:"defaultRom"`
	DefaultMode string `json:"defaultMode"`
	QolPatches  string `json:"qolPatches"`
}

var settings Settings = Settings{}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func Which(cmd string) bool {
	_, err := exec.Command(cmd, "-h").Output()
	return err == nil
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.WindowSetMinSize(a.ctx, 640, 480)
	runtime.WindowSetMaxSize(a.ctx, 640, 480)

	a.InitFolderStructure()

	cfg, err := ini.Load(settings_file)
	if err != nil {
		fmt.Println(err.Error())
		settings.RomsFolder = ""
		settings.DefaultRom = ""
		settings.DefaultMode = ""
	} else {
		fmt.Println("loading settings")
		settings.RomsFolder = cfg.Section("settings").Key("romsFolder").String()
		settings.DefaultRom = cfg.Section("settings").Key("defaultRom").String()
		settings.DefaultMode = cfg.Section("settings").Key("defaultMode").String()
		settings.QolPatches = cfg.Section("settings").Key("qolPatches").String()
	}

}

func (a *App) InitFolderStructure() {
	hostinfo, _ := host.Info()

	if hostinfo.OS == "windows" {
		appdata := os.Getenv("APPDATA")
		appdata += "\\ironlauncher"

		err := os.MkdirAll(appdata, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}

		datapath = appdata + "\\data"
		err = os.MkdirAll(datapath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}

		savespath = datapath + "\\saves"
		err = os.MkdirAll(savespath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}

		gbaSettingsPath = datapath + "\\gba"
		err = os.MkdirAll(gbaSettingsPath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}

		ndsSettingsPath = datapath + "\\nds"
		err = os.MkdirAll(ndsSettingsPath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}

		file, _ := os.OpenFile(datapath+"\\config.ini", os.O_RDONLY|os.O_CREATE, 0644)
		file.Close()
		file, _ = os.OpenFile(gbaSettingsPath+"\\Settings.ini", os.O_RDONLY|os.O_CREATE, 0644)
		file.Close()
		file, _ = os.OpenFile(ndsSettingsPath+"\\Settings.ini", os.O_RDONLY|os.O_CREATE, 0644)
		file.Close()

		settings_file = datapath + "\\ironlauncher.ini"

	} else if hostinfo.OS == "linux" {
		appdata, _ := os.UserHomeDir()
		appdata += "/.local/share/ironlauncher"

		err := os.MkdirAll(appdata, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}

		datapath = appdata + "/data"
		err = os.MkdirAll(datapath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}

		savespath = datapath + "/saves"
		err = os.MkdirAll(savespath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}

		gbaSettingsPath = datapath + "/gba"
		err = os.MkdirAll(datapath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}

		ndsSettingsPath = datapath + "/nds"
		err = os.MkdirAll(datapath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}

		file, _ := os.OpenFile(datapath+"/config.ini", os.O_RDONLY|os.O_CREATE, 0644)
		file.Close()
		file, _ = os.OpenFile(gbaSettingsPath+"/Settings.ini", os.O_RDONLY|os.O_CREATE, 0644)
		file.Close()
		file, _ = os.OpenFile(ndsSettingsPath+"/Settings.ini", os.O_RDONLY|os.O_CREATE, 0644)
		file.Close()

		settings_file = datapath + "/ironlauncher.ini"
	}

}

func (a *App) AreDepsInstalled() bool {
	hostinfo, _ := host.Info()

	result := false
	if hostinfo.OS == "windows" {
		result = Which("podman")
	} else if hostinfo.OS == "linux" {
		result = Which("podman") && Which("pulseaudio") && os.Getenv("DISPLAY") != ""
	}
	fmt.Printf("AreDepsInstalled: %s", strconv.FormatBool(result))
	return result
}

func (a *App) StartUp() bool {
	fmt.Println("Updating container...")
	a.UpdateContainer()
	fmt.Println("Container updated.")

	hostinfo, _ := host.Info()
	if hostinfo.OS == "windows" {
		fmt.Println("Initializing Windows Podman")
		InitWindowsPodman()

		if strings.Contains(hostinfo.Platform, "Windows 10") {
			fmt.Println("Starting vcxsrv and pulseaudio")
			appdata := os.Getenv("APPDATA")
			appdata += "\\ironlauncher"

			cmdStr := strings.Fields(fmt.Sprintf("%s\\vcxsrv\\vcxsrv.exe -ac -multiwindow", appdata))
			cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
			cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
			cmd.Start()

			cmdStr = strings.Fields(fmt.Sprintf("%s\\pulseaudio-1.1\\bin\\pulseaudio.exe", appdata))
			cmd = exec.Command(cmdStr[0], cmdStr[1:]...)
			cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
			cmd.Start()
		}
	}

	return true
}

func InitWindowsPodman() {
	cmdStr := strings.Fields("podman machine init")
	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	cmd.Run()

	cmdStr = strings.Fields("podman machine start")
	cmd = exec.Command(cmdStr[0], cmdStr[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	cmd.Run()
}

func (a *App) GetRomsFolder() string {
	result, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "ROMs Folder",
	})
	return result
}

func (a *App) GetRoms(path string) []string {

	var roms []string

	items, _ := os.ReadDir(path)
	for _, item := range items {
		if !item.IsDir() {
			if strings.HasSuffix(item.Name(), ".gba") ||
				strings.HasSuffix(item.Name(), ".gb") ||
				strings.HasSuffix(item.Name(), ".gbc") ||
				strings.HasSuffix(item.Name(), ".nds") {

				f, err := os.Open(path + string(os.PathSeparator) + item.Name())
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				h := sha1.New()
				if _, err := io.Copy(h, f); err != nil {
					log.Fatal(err)
				}

				hash := hex.EncodeToString(h.Sum(nil))
				_, exists := SUPPORTED_HASHES[hash]

				if exists {
					roms = append(roms, item.Name())
				}
			}
		}
	}

	return roms
}

func (a *App) InstallDependencies() bool {
	hostinfo, _ := host.Info()
	if hostinfo.OS == "windows" {
		out, _ := os.CreateTemp("", "*.exe")

		fmt.Println(out.Name())

		resp, err := http.Get("https://github.com/containers/podman/releases/download/v5.0.2/podman-5.0.2-setup.exe")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer resp.Body.Close()

		io.Copy(out, resp.Body)
		out.Close()

		cmdStr := strings.Fields(fmt.Sprintf("%s /install /passive /quiet", out.Name()))
		fmt.Println(cmdStr)
		cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
		err = cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
		}
		os.Remove(out.Name())

		if strings.Contains(hostinfo.Platform, "Windows 10") {
			fmt.Println("Downloading vcxsrv and pulseaudio")
			appdata := os.Getenv("APPDATA")
			appdata += "\\ironlauncher"

			cmdStr := strings.Fields(fmt.Sprintf("curl https://raw.githubusercontent.com/besteon/IronLauncher/master/win10/vcxsrv.zip -o %s\\vcxsrv.zip", appdata))
			cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
			cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
			err := cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
			}

			cmdStr = strings.Fields(fmt.Sprintf("curl https://raw.githubusercontent.com/besteon/IronLauncher/master/win10/pulseaudio-1.1.zip -o %s\\pulseaudio-1.1.zip", appdata))
			cmd = exec.Command(cmdStr[0], cmdStr[1:]...)
			cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
			err = cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println("Extracting vcxsrv and pulseaudio")

			cmdStr = strings.Fields(fmt.Sprintf("tar -xf %s\\vcxsrv.zip -C %s", appdata, appdata))
			cmd = exec.Command(cmdStr[0], cmdStr[1:]...)
			cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
			err = cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
			}

			cmdStr = strings.Fields(fmt.Sprintf("tar -xf %s\\pulseaudio-1.1.zip -C %s", appdata, appdata))
			cmd = exec.Command(cmdStr[0], cmdStr[1:]...)
			cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
			err = cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		return Which("podman")
	} else if hostinfo.OS == "linux" {
		cmdStr := strings.Fields("sudo -A apt install -y podman pulseaudio xwayland && Xwayland")
		cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "SUDO_ASKPASS=/usr/bin/ssh-askpass")
		cmd.Start()
		return true
	}
	return false
}

func (a *App) Play(romsFolder string, rom string) {
	a.StartContainer(romsFolder)
}

func (a *App) SaveDefaults(romsFolder string, game string, mode string, qolPatches string) {
	cfg := ini.Empty()
	cfg.NewSection("settings")
	cfg.Section("settings").NewKey("romsFolder", romsFolder)
	cfg.Section("settings").NewKey("defaultRom", game)
	cfg.Section("settings").NewKey("defaultMode", mode)
	cfg.Section("settings").NewKey("qolPatches", qolPatches)
	err := cfg.SaveTo(settings_file)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func (a *App) GetSettings() Settings {
	return settings
}

func (a *App) UpdateContainer() {
	cmdStrBuilder := "podman pull docker.io/besteon/ironlauncher:latest"

	cmdStr := strings.Fields(cmdStrBuilder)
	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	err := cmd.Run()

	if err != nil {
		fmt.Println(err.Error())
	}
}

func (a *App) StartContainer(path string) {
	hostinfo, _ := host.Info()
	if hostinfo.OS == "windows" {

		if strings.Contains(hostinfo.Platform, "Windows 10") {
			ip := GetOutboundIP()
			cmdStr := strings.Fields(strings.ReplaceAll(fmt.Sprintf(`wsl --distribution podman-machine-default podman run 
			-e 'DISPLAY=%s:0' 
			-e 'PULSE_SERVER=tcp:%s' 
			-v '%s:/roms' 
			-v '%s:/data' 
			-v '%s:/data/saves'
			-v '%s\config.ini:/home/launcher/BizHawk/config.ini'
			-v '%s\Settings.ini:/home/launcher/BizHawk/Lua/gba/Ironmon-Tracker/Settings.ini'
			-v '%s\Settings.ini:/home/launcher/BizHawk/Lua/nds/Ironmon-Tracker/Settings.ini'
			-v '%s:/home/launcher/ironlauncher.ini'
			--net=host 
			%s`, ip, ip, path, datapath, savespath, datapath, gbaSettingsPath, ndsSettingsPath, settings_file, CONTAINER), `\`, `\\`))
			cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
			cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
			err := cmd.Start()

			if err != nil {
				fmt.Println(err.Error())
			}
		} else if strings.Contains(hostinfo.Platform, "Windows 11") {
			cmdStr := strings.Fields(strings.ReplaceAll(fmt.Sprintf(`wsl --distribution podman-machine-default podman run 
			-e 'DISPLAY=:0' 
			-e 'PULSE_SERVER=/mnt/wslg/PulseServer' 
			-e 'WAYLAND_DISPLAY=wayland-0' 
			-v '/mnt/wslg/:/mnt/wslg/' 
			-v '/mnt/wslg/.X11-unix:/tmp/.X11-unix' 
			-v '%s:/roms'
			-v '%s:/data' 
			-v '%s:/data/saves'
			-v '%s\config.ini:/home/launcher/BizHawk/config.ini'
			-v '%s\Settings.ini:/home/launcher/BizHawk/Lua/gba/Ironmon-Tracker/Settings.ini'
			-v '%s\Settings.ini:/home/launcher/BizHawk/Lua/nds/Ironmon-Tracker/Settings.ini'
			-v '%s:/home/launcher/ironlauncher.ini'
			--net=host 
			%s`, path, datapath, savespath, datapath, gbaSettingsPath, ndsSettingsPath, settings_file, CONTAINER), `\`, `\\`))
			fmt.Println(cmdStr)
			cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
			cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
			err := cmd.Start()

			if err != nil {
				fmt.Println(err.Error())
			}
		}

	} else if hostinfo.OS == "linux" {
		display := os.Getenv("DISPLAY")
		xdg := os.Getenv("XDG_RUNTIME_DIR")
		home, _ := os.UserHomeDir()
		xhost := exec.Command("xhost", "+")
		xhost.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
		xhost.Start()

		cmdStr := strings.Fields(fmt.Sprintf(`podman run
			-e DISPLAY=%s
			-e PULSE_SERVER=unix:%s/pulse/native
			-v %s/pulse/native:%s/pulse/native
			-v %s/.config/pulse/cookie:/root/.config/pulse/cookie
			-v /tmp/.X11-unix:/tmp/.X11-unix:ro
			-v '%s:/roms'
			-v '%s:/data' 
			-v '%s:/data/saves'
			-v '%s/config.ini:/home/launcher/BizHawk/config.ini'
			-v '%s/Settings.ini:/home/launcher/BizHawk/Lua/gba/Ironmon-Tracker/Settings.ini'
			-v '%s/Settings.ini:/home/launcher/BizHawk/Lua/nds/Ironmon-Tracker/Settings.ini'
			-v '%s:/home/launcher/ironlauncher.ini'
			--net=host
			%s`, display, xdg, xdg, xdg, home, path, datapath, savespath, datapath, gbaSettingsPath, ndsSettingsPath, settings_file, CONTAINER))
		cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
		cmd.Start()
	}

	fmt.Println("StartContainer complete")
	go a.PollEmulator()
}

func (a *App) PollEmulator() {
	containerRunning := false
	for !containerRunning {
		cmdStr := strings.Fields("podman ps")
		cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
		out, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
		}
		if strings.Contains(string(out), "ironlauncher") {
			runtime.EventsEmit(a.ctx, "EMULATOR_OPEN")
			containerRunning = true
		}
		time.Sleep(1 * time.Second)
	}
	for {
		cmdStr := strings.Fields("podman ps")
		cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		if !strings.Contains(string(out), "ironlauncher") {
			runtime.EventsEmit(a.ctx, "EMULATOR_CLOSED")
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
