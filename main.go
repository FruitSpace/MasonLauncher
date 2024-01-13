//go:generate /Users/m41den/go/bin/goversioninfo.exe -icon=geometrydash.ico -company=Fruitspace -64
package main

import (
	"embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"os"
	"strings"
)

//go:embed assets
var assets embed.FS

var LockFile = Lock{}

var (
	SrvId   = "000S"
	Version = "0.8"
	Beta    = false
)

func main() {
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-repatch:") {
		arg := strings.Split(os.Args[1], ":")
		if len(arg) > 1 {
			SrvId = arg[1]
		}
		if len(arg) > 2 {
			Beta = arg[2] == "beta"
		}
	}

	go UploadMachineStatistics()

	goApp := app.New()
	goApp.SetIcon(fyne.NewStaticResource("icon.png", GetIcon()))
	goApp.Settings().SetTheme(theme.DarkTheme())
	win := goApp.NewWindow("GhostLauncher")
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(800, 480))

	//drv := fyne.CurrentApp().Driver() // Create splash (borderless) window
	//if drv, ok := drv.(desktop.Driver); ok {
	//	w := drv.CreateSplashWindow()
	//}

	basePath := CreateWorkdir()
	pwd, _ := os.Getwd()

	if lock := CheckGDIntegrity(); lock != "" {
		// GD is installed
		err := LockFile.ReadLock(lock)
		if err != nil {
			dialog.ShowConfirm("Невозможно прочитать fruit.lock", err.Error(), func(b bool) { os.Exit(1) }, win)
			win.ShowAndRun()
		}
		win.SetContent(NewMainPage(win, basePath, pwd))
	} else {
		// GD is not installed
		GDPS, err := LoadServerInfo(SrvId)
		if err != nil {
			dialog.ShowConfirm("Ошибка", err.Error(), func(b bool) {
				GDPS.Name = "Ошибка"
			}, win)
			win.ShowAndRun()
		}
		// Server is found
		LockFile.SrvId = GDPS.SrvId
		LockFile.Title = GDPS.Name
		suffix := ""
		if Beta {
			GDPS.Region = "no"
			suffix = " (Beta)"
		}
		win.SetTitle("GhostLauncher - Установка " + GDPS.Name + suffix)
		win.SetContent(NewInstallPage(win, basePath, pwd, GDPS))
	}

	win.ShowAndRun()

}

func GetIcon() []byte {
	s, _ := assets.ReadFile("assets/geometrydash.png")
	return s
}
