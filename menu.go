/*
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-08-16
@Description: The configuration of program menu
*/

package main

import (
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func Menu(app *App) *menu.Menu {
	const (
		CmdOrCtrlKey   keys.Modifier = "cmdorctrl"
		OptionOrAltKey keys.Modifier = "optionoralt"
		ShiftKey       keys.Modifier = "shift"
		ControlKey     keys.Modifier = "ctrl"
	)

	var (
		FileMenu *menu.MenuItem
		EditMenu *menu.MenuItem
		HelpMenu *menu.MenuItem
	)

	FileMenu = &menu.MenuItem{
		Label: "File",
		Type:  "Text",
		SubMenu: menu.NewMenuFromItems(
			&menu.MenuItem{
				Label:       "New",
				Type:        "Text",
				Accelerator: keys.CmdOrCtrl("N"),
				Click: func(cd *menu.CallbackData) {
					runtime.WindowExecJS(app.ctx, "document.getElementById('btn-new').click();")
				},
			},
			&menu.MenuItem{
				Label:       "Open",
				Type:        "Text",
				Accelerator: keys.CmdOrCtrl("O"),
				Click: func(cd *menu.CallbackData) {
					runtime.WindowExecJS(app.ctx, "document.getElementById('btn-open').click();")
				},
			},
			&menu.MenuItem{
				Label:       "Save",
				Type:        "Text",
				Accelerator: keys.CmdOrCtrl("S"),
				Click: func(cd *menu.CallbackData) {
					runtime.WindowExecJS(app.ctx, "document.getElementById('btn-save').click();")
				},
			},
			&menu.MenuItem{
				Label:       "Save As",
				Type:        "Text",
				Accelerator: keys.Combo("s", ShiftKey, CmdOrCtrlKey),
				Click: func(cd *menu.CallbackData) {
					runtime.WindowExecJS(app.ctx, "document.getElementById('btn-save-as').click();")
				},
			},
			&menu.MenuItem{
				Type: "Separator",
			},
			&menu.MenuItem{
				Label:       "Exit",
				Type:        "Text",
				Accelerator: keys.CmdOrCtrl("Q"),
				Click: func(cd *menu.CallbackData) {
					runtime.Quit(app.ctx)
				},
			},
		),
	}

	EditMenu = &menu.MenuItem{
		Label: "Edit",
		Type:  "Text",
		SubMenu: menu.NewMenuFromItems(
			&menu.MenuItem{
				Label:       "Undo",
				Type:        "Text",
				Accelerator: keys.CmdOrCtrl("Z"),
				Click: func(cd *menu.CallbackData) {

				},
			},
			&menu.MenuItem{
				Label:       "Redo",
				Type:        "Text",
				Accelerator: keys.CmdOrCtrl("Y"),
				Click: func(cd *menu.CallbackData) {

				},
			},
			&menu.MenuItem{
				Type: "Separator",
			},
			&menu.MenuItem{
				Label:       "Cut",
				Type:        "Text",
				Accelerator: keys.CmdOrCtrl("X"),
				Click: func(cd *menu.CallbackData) {

				},
			},
			&menu.MenuItem{
				Label:       "Copy",
				Type:        "Text",
				Accelerator: keys.CmdOrCtrl("C"),
				Click: func(cd *menu.CallbackData) {

				},
			},
			&menu.MenuItem{
				Label:       "Paste",
				Type:        "Text",
				Accelerator: keys.CmdOrCtrl("V"),
				Click: func(cd *menu.CallbackData) {

				},
			},
			&menu.MenuItem{
				Type: "Separator",
			},
			&menu.MenuItem{
				Label:       "Selet All",
				Type:        "Text",
				Accelerator: keys.CmdOrCtrl("A"),
				Click: func(cd *menu.CallbackData) {

				},
			},
		),
	}
	HelpMenu = &menu.MenuItem{
		Label: "Help",
		Type:  "Text",
		Click: func(cd *menu.CallbackData) {

		},
	}

	AppMenu := menu.NewMenuFromItems(FileMenu, EditMenu, HelpMenu)
	return AppMenu
}
