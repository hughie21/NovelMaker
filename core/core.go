// Description: the core package is the core of the program, it is used to manage the configuration, the agent, and the logger.
// Author: Hughie21
// Date: 2024-11-21
// license that can be found in the LICENSE file.
package core

import (
	"embed"
	"os"
	"sync"
	"time"

	"github.com/hughie21/NovelMaker/lib/config"
	"github.com/hughie21/NovelMaker/lib/logging"
	"github.com/hughie21/NovelMaker/lib/server"
	"github.com/hughie21/NovelMaker/lib/utils"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var (
	core   *Core
	once   *sync.Once
	logger *logging.Log
)

type Core struct {
	execPath string
	Args     string
	config   *config.Config
	cm       *config.ConfigManager
	agt      *Agent
	logger   *logging.Log
}

// NewCore create a new core instance, it is a singleton.
func NewCore() *Core {
	if once == nil {
		once = &sync.Once{}
	}
	once.Do(func() {
		core = &Core{}
	})
	return core
}

func (c *Core) Init(assets embed.FS, app *App) *options.App {
	c.execPath = utils.GetCurrentAbPath()
	c.cm = config.NewConfigManager(c.execPath)
	err := c.cm.LoadConfig()
	if err != nil {
		utils.ShowMessage("Load Configuration Failed", err.Error(), "error")
	}
	c.config = c.cm.GetConfig()
	c.agt = NewAgent(c.config.Core.MaxTask, time.Duration(c.config.Core.Timeout))

	// The args is the path of the epub file which is used to open the file directly.
	ArgsLength := len(os.Args)
	if ArgsLength > 1 {
		c.Args = os.Args[1]
	} else {
		c.Args = ""
	}
	logLevel := map[int]logging.Level{
		1: logging.InfoLevel,
		2: logging.WarnLevel,
		3: logging.ErrorLevel,
		5: logging.DebugLevel,
		4: logging.FatalLevel,
	}

	c.logger = logging.NewLog(logLevel[c.config.Log.Level], true)

	logger = c.logger

	// Here is the static server that provide the static resources for the frontend.
	go server.StaticSevice(":"+c.config.StaticResource.Port, c.execPath)

	c.logger.Info("Core Initialized", logging.RunFuncName())
	return c.formOptions(assets, app)
}

// This method is to form the app options based on the configuration.
func (c *Core) formOptions(assets embed.FS, app *App) *options.App {
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

	return &options.App{
		Title:     "NovelMaker",
		Width:     c.config.Appearance.Width,
		Height:    c.config.Appearance.Height,
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
		WindowStartState: windowSize[c.config.Appearance.DefaultOpen],
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
			WebviewGpuIsDisabled: !c.config.Window.GPUAccelerate,
			WebviewUserDataPath:  c.config.Window.WebviewUserData,
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
			WindowIsTranslucent: c.config.Linux.WindowTransparent,
			WebviewGpuPolicy:    GPUPolicy[c.config.Linux.GPUStrategy],
			ProgramName:         "NovelMaker",
		},
	}
}
