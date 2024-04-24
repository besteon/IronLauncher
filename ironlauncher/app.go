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
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var SUPPORTED_HASHES = map[string]string{
	"dd5945db9b930750cb39d00c84da8571feebf417": "Pokemon Fire Red v1.1",
}

var APIKEY string = "asdf"

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
}

func (a *App) AreDepsInstalled() bool {
	hostinfo, _ := host.Info()
	if hostinfo.OS == "windows" {
		return Which("podman")
	} else if hostinfo.OS == "linux" {
		return Which("podman") && Which("pulseaudio") && os.Getenv("DISPLAY") != ""
	}

	return false
}

func (a *App) StartUp() bool {
	hostinfo, _ := host.Info()
	if hostinfo.OS == "windows" {
		InitWindowsPodman()
	}

	return true
}

func InitWindowsPodman() {
	cmdStr := strings.Fields("podman machine init")
	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	cmd.Run()

	cmdStr = strings.Fields("podman machine start")
	cmd = exec.Command(cmdStr[0], cmdStr[1:]...)
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
				name, exists := SUPPORTED_HASHES[hash]

				if exists {
					roms = append(roms, name)
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
		err = cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
		}
		os.Remove(out.Name())

		return Which("podman")
	} else if hostinfo.OS == "linux" {
		cmdStr := strings.Fields("sudo -A apt install -y podman pulseaudio xwayland && Xwayland")
		cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
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

func (a *App) UpdateContainer() {
	hostinfo, _ := host.Info()
	cmdStrBuilder := "podman pull docker.io/besteon/ironlauncher:latest"
	if hostinfo.OS == "windows" {
		cmdStrBuilder = "wsl " + cmdStrBuilder
	}

	cmdStr := strings.Fields(cmdStrBuilder)
	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	cmd.Run()
}

func (a *App) StartContainer(path string) {
	hostinfo, _ := host.Info()
	if hostinfo.OS == "windows" {
		cmdStr := []string{}
		if strings.Contains(hostinfo.Platform, "Windows 10") {
			ip := GetOutboundIP()
			cmdStr = strings.Fields(fmt.Sprintf("wsl podman run -e 'APIKEY=%s' -e 'DISPLAY=%s:0' -e 'PULSE_SERVER=/mnt/wslg/PulseServer' -v '/mnt/wslg/:/mnt/wslg/' -v '%s:/roms' --net=host docker.io/besteon/ironlauncher:latest", APIKEY, ip, path))
			cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
			cmd.Start()
		} else if strings.Contains(hostinfo.Platform, "Windows 11") {
			cmdStr = strings.Fields(fmt.Sprintf("wsl podman run -e 'APIKEY=%s' -e 'DISPLAY=:0' -e 'PULSE_SERVER=/mnt/wslg/PulseServer' -e 'WAYLAND_DISPLAY=wayland-0' -v '/mnt/wslg/:/mnt/wslg/' -v '/mnt/wslg/.X11-unix:/tmp/.X11-unix' -v '%s:/roms' --net=host docker.io/besteon/ironlauncher:latest", APIKEY, path))
			cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
			cmd.Start()
		}

	} else if hostinfo.OS == "linux" {
		display := os.Getenv("DISPLAY")
		xdg := os.Getenv("XDG_RUNTIME_DIR")
		home, _ := os.UserHomeDir()
		xhost := exec.Command("xhost", "+")
		xhost.Start()

		cmdStr := strings.Fields(fmt.Sprintf(`podman run
			-e APIKEY=%s
			-e DISPLAY=%s
			-e PULSE_SERVER=unix:%s/pulse/native
			-v %s/pulse/native:%s/pulse/native
			-v %s/.config/pulse/cookie:/root/.config/pulse/cookie
			-v /tmp/.X11-unix:/tmp/.X11-unix:ro
			-v %s:/roms
			--net=host
			docker.io/besteon/ironlauncher:latest`, APIKEY, display, xdg, xdg, xdg, home, path))
		cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
		cmd.Start()
	}

	fmt.Println("StartContainer complete")
	go a.PollEmulator()
}

func (a *App) PollEmulator() {
	time.Sleep(5 * time.Second)
	for {
		cmdStr := strings.Fields("podman ps")
		out, err := exec.Command(cmdStr[0], cmdStr[1:]...).Output()
		fmt.Printf("%s", out)
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
