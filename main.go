/*
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-10-18
@Description: The core code of the program
*/

package main

import (
	logging "NovelMaker/logging"
	Manager "NovelMaker/manager"
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var (
	//go:embed all:frontend/dist
	assets   embed.FS
	execPath string
	Args     string
	config   *Manager.Config
	cm       *Manager.ConfigManager
	logger   *logging.Log
)

func main() {
	execPath = getCurrentAbPath()
	cm = Manager.NewConfigManager(execPath)
	err := cm.LoadConfig()
	if err != nil {
		logger.Fatal(err.Error(), logging.RunFuncName())
	}
	config = cm.GetConfig()

	go Manager.StaticSevice(":"+config.StaticResource.Port, execPath)

	// Get the arguments from the command line
	// This is used to open the .no file while
	// the user directly open it
	ArgsLength := len(os.Args)
	if ArgsLength > 1 {
		Args = os.Args[1]
	} else {
		Args = ""
	}
	// Create an instance of the app structure
	app := NewApp()
	// AppMenu := Menu(app)

	windowSize := map[string]options.WindowStartState{
		"normal":     options.Normal,
		"fullscreen": options.Fullscreen,
		"maximised":  options.Maximised,
		"minimised":  options.Minimised,
	}

	GPUPolicy := map[string]linux.WebviewGpuPolicy{
		"never":  linux.WebviewGpuPolicyNever,
		"always": linux.WebviewGpuPolicyAlways,
		"auto":   linux.WebviewGpuPolicyOnDemand,
	}

	logLevel := map[int]logging.Level{
		1: logging.InfoLevel,
		2: logging.WarnLevel,
		3: logging.ErrorLevel,
		5: logging.DebugLevel,
		4: logging.FatalLevel,
	}

	logger = logging.NewLog(logLevel[config.Log.Level], true)
	// Create application with options
	err = wails.Run(&options.App{
		Title:     "NovelMaker",
		Width:     config.Appearance.Width,
		Height:    config.Appearance.Height,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// Menu:             AppMenu,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		WindowStartState: windowSize[config.Appearance.DefaultOpen],
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
			WebviewGpuIsDisabled: !config.Window.GPUAccelerate,
			WebviewUserDataPath:  config.Window.WebviewUserData,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "NovelMaker",
				Message: "© 2024 Hughie",
			},
		},
		Linux: &linux.Options{
			WindowIsTranslucent: config.Linux.WindowTransparent,
			WebviewGpuPolicy:    GPUPolicy[config.Linux.GPUStrategy],
			ProgramName:         "NovelMaker",
		},
	})
	if err != nil {
		logger.Fatal(err.Error(), logging.RunFuncName())
		println("Error:", err.Error())
	}
}
