package main

import (
	"GhostPatcher/particles"
	"GhostPatcher/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/terminal"
	"image/color"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type TerminalLayout struct {
	Size fyne.Size
}

var installBtn *widget.Button

func NewTerminalLayout(size fyne.Size) *TerminalLayout {
	return &TerminalLayout{
		Size: size,
	}
}

func (l *TerminalLayout) Layout(objects []fyne.CanvasObject, s fyne.Size) {
	for _, o := range objects {
		o.Resize(s)
		o.Move(fyne.NewPos(0, 0))
	}
}

func (l *TerminalLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return l.Size
}

func helpScreenDialog(win fyne.Window) {
	subClr := color.Gray16{Y: 0xaaaf}
	msvcUrl, _ := url.Parse("https://files.catbox.moe/bdwtei.zip")
	discordUrl, _ := url.Parse("https://discord.gg/fruitspace")

	//open appdata/local in explorer
	openAppData := func() {
		pwd, _ := os.UserCacheDir()
		err := exec.Command("explorer.exe", pwd).Start()
		if err != nil {
			dialog.ShowError(err, win)
		}
	}

	infoAccordion := widget.NewAccordion(
		widget.NewAccordionItem(
			"Лаунчер выдает при запуске ошибку с длинной ссылкой",
			container.NewVBox(
				canvas.NewText("Используйте системный VPN во время установки (например iTop VPN, не реклама).", subClr),
				canvas.NewText("Не забудьте выключить его после установки", subClr),
			),
		),

		widget.NewAccordionItem(
			"Geometry Dash выдает ошибку 0x00000142 или похожую",
			container.NewVBox(
				canvas.NewText("Путь к приватному серверу содержит русские символы. Вы можете перенести папку", subClr),
				canvas.NewText("в C:/Games (обязательно диск C:), чтобы в пути были только английские символы", subClr),
			),
		),

		widget.NewAccordionItem(
			"Geometry Dash выдает ошибку \"MSVCP140.dll not found\" или похожую",
			container.NewVBox(
				canvas.NewText("Вам требуется установить MSVC (Microsoft Visual C++ Redistributable)", subClr),
				canvas.NewText("Мы собрали для вас архив со всем необходимым: установите оттуда все 4 файла", subClr),
				widget.NewHyperlink("https://files.catbox.moe/bdwtei.zip", msvcUrl),
			),
		),

		widget.NewAccordionItem(
			"Geometry Dash не запускается вообще (и не отображает ошибку)",
			container.NewVBox(
				canvas.NewText("1) Убедитесь, что игра находится на диске C:", subClr),
				canvas.NewText("2) Если не помогло, то попробуйте удалить сохранения: нажмите на кнопку ниже и удалите папку с названием", subClr),
				canvas.NewText("    вашего GDPS. Это баг самого Geometry Dash и альтернативы нет :(", subClr),
				widget.NewButton("Открыть папку", openAppData),
			),
		),

		widget.NewAccordionItem(
			"Я установил MegaHack v7, но вижу ошибку Themida - An error has occurred... wrong DLL present.",
			container.NewVBox(
				canvas.NewText("Приватные сервера не поддерживают мегахак, но некоторый пользователи сообщали о рабочем способе:", subClr),
				canvas.NewText("1) Перенесите файл приватного сервера (ИмяGDPS.exe) в папку офиц. GD из Steam", subClr),
				canvas.NewText("2) Если в той папке нет файла \"steam_appid.txt\", то создайте его и поместите в него 322170", subClr),
			),
		),

		widget.NewAccordionItem(
			"Ничего не помогло",
			container.NewVBox(
				canvas.NewText("И такое тоже бывает. Присоединятесь к нашему Discord-серверу и открывайте там тикет - мы обязательно поможем!", subClr),
				widget.NewHyperlink("https://discord.gg/fruitspace", discordUrl),
			),
		),
	)

	dialog.ShowCustom("Помощь", "Закрыть", container.NewVBox(
		canvas.NewText("Если у вас возникают проблемы при запуске Geometry Dash, то попробуйте следующие решения:", color.White),
		infoAccordion,
	), win)
}

func NewInstallPage(win fyne.Window, basePath string, pwd string, GDPS utils.Server) fyne.CanvasObject {

	r, w, _ := os.Pipe()

	os.Stdout = w
	os.Stderr = w

	t := terminal.New()
	go func() {
		t.RunWithConnection(os.Stdin, r)
		//t.RunLocalShell()
	}()

	// Make Card --- V(H(img, V(text, text)), btn)

	currentYear := strconv.Itoa(time.Now().Year())
	copyright := canvas.NewText("GhostLauncher v"+Version+" © "+currentYear+" Fruitspace", color.Gray16{Y: 0xaaaf})
	copyright.TextSize = 10
	copyright.Alignment = fyne.TextAlignCenter

	title := canvas.NewText(GDPS.Name, color.White)
	title.TextSize = 20

	texturesWarning := canvas.NewText("Используется Particle Engine", color.Gray16{Y: 0xaaaf})
	texturesWarning.TextSize = 12

	DataCard := container.NewVBox(
		title,
		canvas.NewText(
			fmt.Sprintf("Игроков: %d,   Уровней: %d", GDPS.Players, GDPS.Levels),
			color.White),
		texturesWarning,
	)

	logo := &canvas.Image{}
	if utils.CacheIcon(GDPS.Icon, "") != nil {
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

	installBtn = widget.NewButtonWithIcon("Установить", theme.DownloadIcon(), func() {
		installBtn.Disable()
		go InstallGD(GDPS, pwd, win)
	})
	helpBtn := widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
		go helpScreenDialog(win)
	})

	ActionBox := container.NewBorder(nil, nil, nil,

		helpBtn,
		installBtn,
	)
	Pane := container.NewCenter(container.NewVBox(
		Card,
		ActionBox,
	))
	tlay := container.New(NewTerminalLayout(fyne.NewSize(800, 160)),
		container.NewStack(canvas.NewRectangle(color.RGBA{R: 0x19, G: 0x19, B: 0x22, A: 255}), t))
	return container.NewBorder(copyright, tlay, nil, nil, Pane)
}

func InstallGD(GDPS utils.Server, pwd string, w fyne.Window) {

	stopper := false

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	p := particles.NewParticle()

	// region Prepare Build Env
	where := filepath.Join(pwd, GDPS.Name)
	os.MkdirAll(filepath.Join(where, ".build"), 0777)

	logfile, err := os.Create(filepath.Join(where, ".build", "log.txt"))
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	defer logfile.Close()
	golog := func(s error) {
		if s != nil {
			fmt.Fprintln(logfile, s)
			fmt.Println(s)
			dialog.ShowError(s, w)
			stopper = true
			installBtn.Enable()
		}
	}

	p.InitFolder(filepath.Join(where, ".build"))
	manifest := p.GenerateMainfestFor(GDPS.SrvId, GDPS.Version)
	if len(GDPS.Recipe) > 0 {
		manifest = GDPS.Recipe
	}
	golog(os.WriteFile(filepath.Join(where, ".build", "particle.json"), []byte(manifest), 0755))
	if stopper {
		return
	}
	// endregion

	// Etags
	_, etag, err := utils.GetWebFileInfo("https://cdn.fruitspace.one/assets/GhostLauncher.exe")
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось получить информацию о лаунчере", func(b bool) { os.Exit(1) }, w)
		return
	}
	LockFile.LauncherEtag = etag

	_, iconEtag, err := utils.GetWebFileInfo(GDPS.Icon)
	if err != nil {
		dialog.ShowConfirm("Ошибка", "Не удалось получить информацию об иконке", func(b bool) { os.Exit(1) }, w)
		return
	}
	LockFile.IconEtag = iconEtag

	// region Build

	e := p.Prepare(filepath.Join(where, ".build"))
	if e != nil && strings.HasPrefix(e.Error(), "symlink") {
		helpPermissionsPage(w)
		return
	}
	golog(e)
	if stopper {
		return
	}
	logfile.Close()
	// endregion

	golog(p.MoveBuild(where))
	if stopper {
		return
	}
	utils.Update(pwd + "/" + GDPS.Name + "/GhostLauncher.exe")
	LockFile.WriteLock(pwd + "/" + GDPS.Name)
	dialog.ShowConfirm("Установка завершена", GDPS.Name+" успешно установлен. Хотите запустить?", func(b bool) {
		if b {
			utils.StartBinaryDetached(pwd + "/" + GDPS.Name + "/" + GDPS.Name + ".exe")
			os.Exit(0)
		}
	}, w)
	installBtn.Enable()
}

func NewMainPage(win fyne.Window, basePath string, pwd string) fyne.CanvasObject {

	// Check internet connection
	_, _, inetErr := utils.GetWebFileInfo("https://google.com")
	GDPS := utils.Server{
		SrvId: LockFile.SrvId,
		Name:  LockFile.Title,
	}
	desc := "Добро пожаловать!"
	stat := "Офлайн режим"
	manager := utils.SaveManager{}
	xerr := manager.Open(GDPS.Name)
	if xerr == nil {
		uname := manager.GetUname()
		desc = "Добро пожаловать, " + uname + "!"
	}
	if inetErr == nil {
		GDPS, _ = utils.LoadServerInfo(LockFile.SrvId)
		stat = fmt.Sprintf("Игроков: %d,   Уровней: %d", GDPS.Players, GDPS.Levels)
		_, etag, err := utils.GetWebFileInfo("https://cdn.fruitspace.one/assets/GhostLauncher.exe")
		if err == nil && LockFile.LauncherEtag != etag {
			utils.SelfUpdate()
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
		canvas.NewText(desc, color.White),
		xstat,
	)

	logo := &canvas.Image{}
	if utils.CacheIcon(GDPS.Icon, LockFile.IconEtag) != nil {
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

	runBtn := widget.NewButtonWithIcon("Запустить", theme.MediaPlayIcon(), func() {
		utils.StartBinaryDetached(pwd + "/" + GDPS.Name + ".exe")
		os.Exit(0)
	})
	installBtn = widget.NewButtonWithIcon("Переустановить", theme.UploadIcon(), func() {
		// Get pwd parent folder
		pwd = filepath.Dir(pwd)
		installBtn.Disable()
		go InstallGD(GDPS, pwd, win)
	})
	Pane := container.NewCenter(container.NewVBox(Card, runBtn, installBtn))
	return container.NewBorder(copyright, nil, nil, nil, Pane)

	//_, iconEtag, err:= GetWebFileInfo(GDPS.Icon)
	//if err != nil {
	//	dialog.ShowConfirm("Ошибка", "Не удалось получить информацию об иконке", func(b bool) {os.Exit(1)}, w)
	//	return
	//}
	//LockFile.IconEtag=iconEtag
}

func helpPermissionsPage(win fyne.Window) {
	logo1 := &canvas.Image{}
	i1, _ := assets.Open("assets/enable_dev.png")
	logo1 = canvas.NewImageFromReader(i1, "dev.png")
	logo1.FillMode = canvas.ImageFillContain
	logo1.SetMinSize(fyne.NewSize(40*16, 360))

	logo2 := &canvas.Image{}
	i2, _ := assets.Open("assets/enable_adm.png")
	logo2 = canvas.NewImageFromReader(i2, "adm.png")
	logo2.FillMode = canvas.ImageFillContain
	logo2.SetMinSize(fyne.NewSize(504, 360))

	page := container.NewScroll(container.NewVBox(
		canvas.NewText("Для Windows 10 и выше, включите режим разработчика в настройках", color.White),
		canvas.NewText("Где найти: Настройки -> Система -> Для разработчиков", color.White),
		logo1,
		canvas.NewText("Для Windows 7 и выше, в настройках совместимости файла отметьте", color.White),
		canvas.NewText("опцию \"Запуск от имени администратора\"", color.White),
		logo2,
		canvas.NewText("Это ограничение Windows для символьных ссылок, а копировать много", color.White),
		canvas.NewText("раз файлы и засорять диск мы не хотим. Спасибо за понимание", color.White),
	))
	page.SetMinSize(fyne.NewSize(40*16, 360))
	dialog.ShowCustom("Ошибка разрешений", "понятно", page, win)
}
