/*
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-07-12
*/

package main

import (
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
				DisplayName: "Epmb Files",
				Pattern:     pattern,
			},
		},
	})
	if err != nil {
		return err.Error()
	}
	return saved
}
