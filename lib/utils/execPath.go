/*
@Author: Hughie
@CreateTime: 2024-10-14
@LastEditors: Hughie
@LastEditTime: 2024-11-1
@Description: This is the program that get the right execution path of the program when go build or run
*/

package utils

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

func getCurrentAbPathByExecutable() string {
	exePath, err := os.Getwd()
	if err != nil {
		ShowMessage("Error", "Failed to get current directory: "+err.Error(), "error")
	}
	return exePath
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
