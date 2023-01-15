package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewInstallPage(win fyne.Window, basePath string, pwd string, GDPS Server) fyne.CanvasObject {
	downBar := widget.NewProgressBar()
	downBar.Max = 100

	// Make Card --- V(H(img, V(text, text)), btn)

	currentYear := strconv.Itoa(time.Now().Year())
	copyright := canvas.NewText("GhostLauncher v"+Version+" © "+currentYear+" Fruitspace", color.Gray16{Y: 0xaaaf})
	copyright.TextSize = 10
	copyright.Alignment = fyne.TextAlignCenter

	title := canvas.NewText(GDPS.Name, color.White)
	title.TextSize = 20

	texturesWarning := canvas.NewText("", color.Gray16{Y: 0xaaaf})
	if GDPS.TexturePack != "gdps_textures.zip" {
		texturesWarning = canvas.NewText("Данный GDPS использует текстурпак", color.Gray16{Y: 0xcccf})
		texturesWarning.TextSize = 12
	}

	DataCard := container.NewVBox(
		title,
		canvas.NewText(
			fmt.Sprintf("Игроков: %d,   Уровней: %d", GDPS.Players, GDPS.Levels),
			color.White),
		texturesWarning,
	)

	logo := &canvas.Image{}
	if CacheIcon(GDPS.Icon, "") != nil {
		i, _ := assets.Open("assets/geometrydash.png")
		logo = canvas.NewImageFromReader(i, "geometrydash.png")
	} else {
		logo = canvas.NewImageFromFile(basePath + "/cache/" + GDPS.Icon[strings.LastIndex(GDPS.Icon, "/")+1:])
	}

	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(100, 100))
	Card := container.NewHBox(
		logo,
		DataCard,
	)

	installBtn := widget.NewButtonWithIcon("Установить", theme.DownloadIcon(), func() {
		go InstallGD(GDPS, downBar, basePath, pwd, win)
	})
	Pane := container.NewCenter(container.NewVBox(Card, installBtn))
	return container.NewBorder(copyright, downBar, nil, nil, Pane)
}

func InstallGD(GDPS Server, progressBar *widget.ProgressBar, basePath, pwd string, w fyne.Window) {
	progressBar.TextFormatter = func() string {
		return "Готовим файлы..."
	}
	os.MkdirAll(pwd+"/"+GDPS.Name, 0777)

	// Get Sizes and Etags
	dllSize, dllEtag, err := GetWebFileInfo("https://cdn.fruitspace.one/assets/gdps_dlls.zip")
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось получить информацию о библиотеках", func(b bool) { os.Exit(1) }, w)
		return
	}
	LockFile.DllEtag = dllEtag

	texturesSize, texturesEtag, err := GetWebFileInfo("https://cdn.fruitspace.one/assets/gdps_textures.zip")
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось получить информацию о базовых текстурах", func(b bool) { os.Exit(1) }, w)
		return
	}
	LockFile.TextureEtag = texturesEtag

	_, etag, err := GetWebFileInfo("https://cdn.fruitspace.one/assets/GhostLauncher.exe")
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось получить информацию о лаунчере", func(b bool) { os.Exit(1) }, w)
		return
	}
	LockFile.LauncherEtag = etag

	_, iconEtag, err := GetWebFileInfo(GDPS.Icon)
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось получить информацию об иконке", func(b bool) { os.Exit(1) }, w)
		return
	}
	LockFile.IconEtag = iconEtag

	// Download DLLs
	progressBar.TextFormatter = func() string {
		return fmt.Sprintf("Загружаем необходимые библиотеки... (%.2f%%)", progressBar.Value)
	}
	go UpdateProgress(progressBar, dllSize, basePath+"/gdps_dlls.zip")
	err = DownloadDefaultDLLs()
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось загрузить библиотеки", func(b bool) { os.Exit(1) }, w)
		return
	}
	time.Sleep(time.Millisecond * 200)
	progressBar.SetValue(0)
	DownloadingNow = make(chan int64)
	fmt.Sprintln("Set 0")

	// Verify Dlls size
	dllSizeLocal, err := GetFileSize(basePath + "/gdps_dlls.zip")
	if err != nil || int(dllSizeLocal) != dllSize {
		os.Remove(basePath + "/gdps_dlls.zip")
		dialog.ShowConfirm("Ошибка", "Не удалось загрузить библиотеки. Попробуйте еще раз", func(b bool) { os.Exit(1) }, w)
		return
	}

	// Download Textures
	progressBar.TextFormatter = func() string {
		return fmt.Sprintf("Загружаем базовый текстурпак... (%.2f%%)", progressBar.Value)
	}
	go UpdateProgress(progressBar, texturesSize, basePath+"/gdps_textures.zip")
	err = DownloadDefaultTextures()
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось загрузить базовые текстуры", func(b bool) { os.Exit(1) }, w)
		return
	}
	time.Sleep(time.Millisecond * 200)
	progressBar.SetValue(0)
	DownloadingNow = make(chan int64)

	// Verify Textures size
	texturesSizeLocal, err := GetFileSize(basePath + "/gdps_textures.zip")
	if err != nil || int(texturesSizeLocal) != texturesSize {
		os.Remove(basePath + "/gdps_textures.zip")
		dialog.ShowConfirm("Ошибка", "Не удалось загрузить базовые текстуры. Попробуйте еще раз", func(b bool) { os.Exit(1) }, w)
		return
	}

	if GDPS.TexturePack != "gdps_textures.zip" {

		texturesOverlaySize, texturesOverlayEtag, err := GetWebFileInfo("https://cdn.fruitspace.one/gdps_textures/" + GDPS.TexturePack)
		if err != nil {
			dialog.ShowConfirm("Ошибка", "Не удалось получить информацию о GDPS текстурах", func(b bool) { os.Exit(1) }, w)
			return
		}

		DialogLock := true
		dialog.ShowConfirm("Загрузка текстур FruitPack",
			fmt.Sprintf("%s использует кастомный текстурпак размером %.2fMb. Загрузить?", GDPS.Name, float64(texturesOverlaySize)/1024/1024),
			func(b bool) {
				fmt.Println(b)
				defer func() { DialogLock = false }()
				if !b {
					GDPS.TexturePack = "gdps_textures.zip"
					return
				}

				//Accept
				LockFile.TextureOverlayEtag = texturesOverlayEtag
				progressBar.TextFormatter = func() string {
					return fmt.Sprintf("Загружаем текстурпак для GDPS... (%.2f%%)", progressBar.Value)
				}
				go UpdateProgress(progressBar, texturesSize, basePath+"/"+GDPS.SrvId+"_textures.zip")
				err = DownloadCustomTextures(GDPS.SrvId, GDPS.TexturePack)
				if err != nil {
					dialog.ShowConfirm("Ошибка", "Не удалось загрузить текстуры для GDPS", func(b bool) { os.Exit(1) }, w)
					return
				}
				time.Sleep(time.Millisecond * 200)
				progressBar.SetValue(0)
				DownloadingNow = make(chan int64)

				// Verify Textures size
				texturesOverlaySizeLocal, err := GetFileSize(basePath + "/" + GDPS.SrvId + "_textures.zip")
				if err != nil || int(texturesOverlaySizeLocal) != texturesOverlaySize {
					os.Remove(basePath + "/" + GDPS.SrvId + "_textures.zip")
					dialog.ShowConfirm("Ошибка", "Не удалось загрузить текстуры для GDPS. Попробуйте еще раз", func(b bool) { os.Exit(1) }, w)
					return
				}
			}, w)

		for DialogLock {
			time.Sleep(time.Millisecond * 100)
		}
	}

	// Unpack all
	progressBar.TextFormatter = func() string { return "Распаковываем файлы..." }
	progressBar.SetValue(0)
	err = UnzipFile(basePath+"/gdps_dlls.zip", pwd+"/"+GDPS.Name)
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось распаковать библиотеки", func(b bool) { os.Exit(1) }, w)
		return
	}
	err = UnzipFile(basePath+"/gdps_textures.zip", pwd+"/"+GDPS.Name)
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось распаковать базовые текстуры", func(b bool) { os.Exit(1) }, w)
		return
	}
	if GDPS.TexturePack != "gdps_textures.zip" {
		err = UnzipFile(basePath+"/"+GDPS.SrvId+"_textures.zip", pwd+"/"+GDPS.Name+"/Resources")
		if err != nil {
			dialog.ShowConfirm("Ошибка", "Не удалось распаковать текстуры GDPS", func(b bool) { os.Exit(1) }, w)
			return
		}
	}

	// Download GD
	progressBar.TextFormatter = func() string { return "Загружаем " + GDPS.Name + "..." }
	progressBar.SetValue(0)
	repatch := RePatcher{}
	gd, err := repatch.DownloadPureGD()
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось загрузить "+GDPS.Name, func(b bool) { os.Exit(1) }, w)
		return
	}
	gd = repatch.PatchPureGD(GDPS.GetUrl(), gd)
	err = WriteBytes(pwd+"/"+GDPS.Name+"/"+GDPS.Name+".exe", gd)
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось записать "+GDPS.Name, func(b bool) { os.Exit(1) }, w)
		return
	}
	Update(pwd + "/" + GDPS.Name + "/GhostLauncher.exe")
	LockFile.WriteLock(pwd + "/" + GDPS.Name)
	progressBar.TextFormatter = func() string { return "Готово!" }
	progressBar.SetValue(0)
	dialog.ShowConfirm("Установка завершена", GDPS.Name+" успешно установлен. Хотите запустить?", func(b bool) {
		if b {
			StartBinaryDetached(pwd + "/" + GDPS.Name + "/" + GDPS.Name + ".exe")
			os.Exit(0)
		}
	}, w)
}

func UpdateProgress(progressBar *widget.ProgressBar, size int, target string) {
	stop := false
	for {
		fmt.Printf(".")
		select {
		case <-DownloadingNow:
			stop = true
		default:
			val := GetDownloadPercent(size, target)
			progressBar.SetValue(val)
		}
		if stop {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func NewMainPage(win fyne.Window, basePath string, pwd string) fyne.CanvasObject {

	// Check internet connection
	_, _, inetErr := GetWebFileInfo("https://google.com")
	GDPS := Server{
		SrvId: LockFile.SrvId,
		Name:  LockFile.Title,
	}
	desc := "Добро пожаловать!"
	stat := "Офлайн режим"
	manager := SaveManager{}
	xerr := manager.Open(GDPS.Name)
	if xerr == nil {
		uname := manager.GetUname()
		desc = "Добро пожаловать, " + uname + "!"
	}
	if inetErr == nil {
		GDPS, _ = LoadServerInfo(LockFile.SrvId)
		stat = fmt.Sprintf("Игроков: %d,   Уровней: %d", GDPS.Players, GDPS.Levels)
		_, etag, err := GetWebFileInfo("https://cdn.fruitspace.one/assets/GhostLauncher.exe")
		if err == nil && LockFile.LauncherEtag != etag {
			SelfUpdate()
			LockFile.LauncherEtag = etag
			LockFile.WriteLock(pwd)
			dialog.ShowInformation("Обновление", "Обновление завершено. Перезапустите лаунчер", win)
		}
	}
	win.SetTitle("GhostLauncher - " + GDPS.Name)

	currentYear := strconv.Itoa(time.Now().Year())
	copyright := canvas.NewText("GhostLauncher v"+Version+" © "+currentYear+" Fruitspace", color.Gray16{Y: 0xaaaf})
	copyright.TextSize = 10
	copyright.Alignment = fyne.TextAlignCenter

	title := canvas.NewText(GDPS.Name, color.White)
	title.TextSize = 20
	xstat := canvas.NewText(stat, color.White)
	xstat.TextSize = 10
	DataCard := container.NewVBox(
		title,
		canvas.NewText(
			desc,
			color.White),
		xstat,
	)

	logo := &canvas.Image{}
	if CacheIcon(GDPS.Icon, LockFile.IconEtag) != nil {
		i, _ := assets.Open("assets/geometrydash.png")
		logo = canvas.NewImageFromReader(i, "geometrydash.png")
	} else {
		logo = canvas.NewImageFromFile(basePath + "/cache/" + GDPS.Icon[strings.LastIndex(GDPS.Icon, "/")+1:])
	}

	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(100, 100))
	Card := container.NewHBox(
		logo,
		DataCard,
	)

	installBtn := widget.NewButtonWithIcon("Запустить", theme.MediaPlayIcon(), func() {
		StartBinaryDetached(pwd + "/" + GDPS.Name + ".exe")
		os.Exit(0)
	})
	Pane := container.NewCenter(container.NewVBox(Card, installBtn))
	return container.NewBorder(copyright, nil, nil, nil, Pane)

	//_, iconEtag, err:= GetWebFileInfo(GDPS.Icon)
	//if err != nil {
	//	dialog.ShowConfirm("Ошибка", "Не удалось получить информацию об иконке", func(b bool) {os.Exit(1)}, w)
	//	return
	//}
	//LockFile.IconEtag=iconEtag
}
