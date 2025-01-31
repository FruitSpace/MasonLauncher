package main

import (
	"embed"
	"flag"
	"github.com/wailsapp/wails/v2/pkg/options/mac"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var GDPS_ID string

func init() {
	PrepareLauncher()
	flag.StringVar(&GDPS_ID, "gdps", "", "GDPS ID to launch or install")
}

func main() {
	app := NewApp()
	err := wails.Run(&options.App{
		Title:  "Mason Launcher",
		Width:  900,
		Height: 576,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 22, G: 30, B: 43, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       false,
			},
			OnUrlOpen: func(url string) {

			},
		},
		//DisableResize: true,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
