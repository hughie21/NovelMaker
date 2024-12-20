// Description: This main entrance of the program.
// Author: Hughie21
// Date: 2024-12-16
// license that can be found in the LICENSE file.

package main

import (
	"embed"

	"github.com/hughie21/NovelMaker/core"
	"github.com/hughie21/NovelMaker/lib/utils"
	"github.com/wailsapp/wails/v2"
)

var (
	//go:embed all:frontend/dist
	assets embed.FS
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			utils.ShowMessage("Core Dump", err.(error).Error(), "error")
		}
	}()

	// Create an instance of the app structure
	app := core.NewApp()

	core := core.NewCore()

	// Create application with options
	err := wails.Run(core.Init(assets, app))
	if err != nil {
		utils.ShowMessage("Core Dump", err.Error(), "error")
		panic(err)
	}
}
