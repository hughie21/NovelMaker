/*
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-10-18
@Description: The entry of the program
*/

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
