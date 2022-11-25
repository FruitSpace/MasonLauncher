package main

import (
	"embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"os"
	"strings"
)

//go:embed assets
var assets embed.FS

var LockFile = Lock{}

var (
	SrvId = "000S"
	Version = "0.4"
)


func main()	{
	if len(os.Args)>1 && strings.HasPrefix(os.Args[1],"-repatch:") {
		SrvId = strings.Split(os.Args[1],":")[1]
	}
	goApp := app.New()
	goApp.SetIcon(fyne.NewStaticResource("icon.png", GetIcon()))
	win := goApp.NewWindow("GhostLauncher")
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(800, 480))

	basePath:=CreateWorkdir()
	pwd,_:=os.Getwd()


	if lock:=CheckGDIntegrity(); lock!="" {
		// GD is installed
		err:=LockFile.ReadLock(lock)
		if err!=nil {
			dialog.ShowConfirm("Невозможно прочитать fruit.lock", err.Error(), func(b bool) {os.Exit(1)}, win)
			win.ShowAndRun()
		}
		win.SetContent(NewMainPage(win, basePath, pwd))
	} else {
		// GD is not installed
		GDPS,err:=LoadServerInfo(SrvId)
		if err!=nil {
			dialog.ShowConfirm("Ошибка", err.Error(), func(b bool) {os.Exit(1)}, win)
			win.ShowAndRun()
		}
		// Server is found
		LockFile.SrvId=GDPS.SrvId
		LockFile.Title=GDPS.Name
		win.SetTitle("GhostLauncher - Установка "+GDPS.Name)
		win.SetContent(NewInstallPage(win, basePath, pwd, GDPS))
	}


	win.ShowAndRun()

}

func GetIcon() []byte {
	s,_:=assets.ReadFile("assets/geometrydash.png")
	return s
}