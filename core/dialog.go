// Description: This file encapsulated the system dialog provided by wails runtime, which is used to open and save files.
// Author: Hughie21
// Date: 2024-10-12
// license that can be found in the LICENSE file.

package core

import (
	logging "github.com/hughie21/NovelMaker/lib/logging"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func FileOpenDialog(app *App, displayName string, pattern string) string {
	selected, err := runtime.OpenFileDialog(app.ctx, runtime.OpenDialogOptions{
		Title: "Open File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: displayName,
				Pattern:     pattern,
			},
		},
		CanCreateDirectories:       true,
		ResolvesAliases:            true,
		TreatPackagesAsDirectories: true,
	})
	if err != nil {
		logger.Fatal(err.Error(), logging.RunFuncName())
		return err.Error()
	}
	return selected
}

func FileSaveDialog(app *App, filename string, pattern string) string {
	saved, err := runtime.SaveFileDialog(app.ctx, runtime.SaveDialogOptions{
		DefaultFilename: filename,
		Title:           "Save File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "NovelMaker Files",
				Pattern:     pattern,
			},
		},
	})
	if err != nil {
		logger.Fatal(err.Error(), logging.RunFuncName())
		return err.Error()
	}
	return saved
}
