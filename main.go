/*
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-08-16
@Description: The core code of the program
*/

package main

import (
	Manager "NovelMaker/manager"
	"embed"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

var execPath string

var Args string

func main() {
	go Manager.StaticSevice()
	execPath, _ = os.Executable()
	execPath = filepath.Dir(execPath)
	ArgsLength := len(os.Args)
	if ArgsLength > 1 {
		Args = os.Args[1]
	} else {
		Args = ""
	}
	// Create an instance of the app structureu
	app := NewApp()
	AppMenu := Menu(app)
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "NovelMaker",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             AppMenu,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		WindowStartState: options.Maximised,
		Windows: &windows.Options{
			Theme: windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
