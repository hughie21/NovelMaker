// Description: This is the Go function that frontend can call for.
// Author: Hughie21
// Date: 2024-09-16
// license that can be found in the LICENSE file.
package core

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hughie21/NovelMaker/lib/logging"
	"github.com/hughie21/NovelMaker/lib/server"
	"github.com/hughie21/NovelMaker/lib/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"context"
)

// App struct
type App struct {
	ctx          context.Context
	core         *Core
	resourceFold string
}

// NewApp creates a new App application struct
func NewApp() *App {
	hash, err := utils.GenerateHash([]byte(time.Now().String()))
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		logger.LogOutPut(core.execPath)
		utils.ShowMessage("Error when generating hash: ", err.Error(), "error")
		panic(err)
	}
	return &App{
		resourceFold: hash,
	}
}

// the message formula that communiacte with the frontend
type Message struct {
	Code int
	Msg  string
	Data string
}

type ImageFIle struct {
	Code int
	Name string
	Id   string
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.core = NewCore()
	a.ctx = ctx
	if core.Args != "" {
		hash, err := utils.GenerateHash([]byte(core.Args))
		if err != nil {
			logger.Error(err.Error(), logging.RunFuncName())
			logger.LogOutPut(core.execPath)
			utils.ShowMessage("Error when generating hash: ", err.Error(), "error")
			panic(err)
		}
		a.resourceFold = hash
		return
	}
	defaultPath := filepath.Join(core.execPath, "resources", a.resourceFold)
	if !utils.PathExists(defaultPath) {
		err := os.Mkdir(defaultPath, os.ModePerm)
		if err != nil {
			logger.Fatal(err.Error(), logging.RunFuncName())
			logger.LogOutPut(a.core.execPath)
			utils.ShowMessage("Error when creating resource folder: ", err.Error(), "error")
			panic(err)
		}
	}
}

func (a *App) shutdown(ctx context.Context) {
	logger.Info("App shutdown", logging.RunFuncName())
	err := a.core.cm.SaveConfig()
	if err != nil {
		logger.LogOutPut(a.core.execPath)
		utils.ShowMessage("Error when writing config: ", err.Error(), "error")
		panic(err)
	}
	logger.Info("Config saved", logging.RunFuncName())
	os.RemoveAll(filepath.Join(a.core.execPath, "resources", a.resourceFold))
	logger.Info("Temporary resource folder removed", logging.RunFuncName())
	core.agt.Close()
	logger.LogOutPut(a.core.execPath)
}

// return the raw data of file
func (a *App) Fr(path string) string {
	fp, err := os.Open(path)
	if err != nil {
		// LogOutPut(err.Error(), runFuncName())
		logger.Error(err.Error(), logging.RunFuncName())
		return err.Error()
	}
	defer fp.Close()
	content, _ := io.ReadAll(fp)
	return string(content)
}

func (a *App) Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func (a *App) Base64Decode(str string) string {
	decodeBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes)
}

func (a *App) FileOpen() Message {
	var Message Message
	res := FileOpenDialog(a, "Epub File", "*.epub")
	if res == "" {
		Message.Code = -1
		Message.Msg = "cancel"
		return Message
	}
	os.RemoveAll(filepath.Join(a.core.execPath, "resources", a.resourceFold))
	logger.Info("Reading Epub file", logging.RunFuncName())
	futue, err := a.core.agt.Exec("reader", res)
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		Message.Code = 1
		Message.Msg = err.Error()
		return Message
	}
	ctxGet, cancelGet := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelGet()
	returnData, err := futue.Get(ctxGet)
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		Message.Code = 1
		Message.Msg = err.Error()
		return Message
	}
	if returnData.Err() != nil {
		logger.Error(returnData.Err().Error(), logging.RunFuncName())
		Message.Code = 1
		Message.Msg = returnData.Err().Error()
		return Message
	}
	hash, err := utils.GenerateHash([]byte(res))
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		Message.Code = 1
		Message.Msg = err.Error()
		return Message
	}
	a.resourceFold = hash
	Message.Code = 0
	Message.Msg = res
	Message.Data = returnData.Data().(string)
	return Message
}

func (a *App) NewFile() Message {
	var msg Message
	err := os.RemoveAll(filepath.Join(a.core.execPath, "resources", a.resourceFold))
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	hash, err := utils.GenerateHash([]byte(time.Now().String()))
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	a.resourceFold = hash
	os.Mkdir(filepath.Join(a.core.execPath, "resources", a.resourceFold), os.ModePerm)
	msg.Code = 0
	msg.Msg = "success"
	return msg
}

func (a *App) DirectLoading() Message {
	var msg Message
	if a.core.Args == "" {
		msg.Code = -1
		msg.Msg = "no file"
		return msg
	}

	future, err := a.core.agt.Exec("reader", a.core.Args)
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		msg.Code = 1
		msg.Msg = err.Error()
		return msg
	}
	ctxGet, cancelGet := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelGet()
	returnData, err := future.Get(ctxGet)
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		msg.Code = 1
		msg.Msg = err.Error()
		return msg
	}
	if returnData.Err() != nil {
		logger.Error(fmt.Sprintf("%s: %s", returnData.Err().Error(), a.core.Args), logging.RunFuncName())
		msg.Code = 1
		msg.Msg = returnData.Err().Error()
		return msg
	}
	msg.Code = 0
	msg.Msg = a.core.Args
	msg.Data = returnData.Data().(string)
	return msg
}

// func (a *App) FileImport() Message {
// 	var Message Message
// 	res := FileOpenDialog(a, "Text File", "*.md;*.txt")
// 	if res == "" {
// 		Message.Code = -1
// 		Message.Msg = "cancel"
// 		return Message
// 	}
// 	returnData := a.core.agt.Exec("markdown", res)
// 	if returnData.Err() != nil {
// 		logger.Error(returnData.Err().Error(), logging.RunFuncName())
// 		Message.Code = 1
// 		Message.Msg = returnData.Err().Error()
// 		return Message
// 	}
// 	Message.Code = 0
// 	Message.Msg = "success"
// 	Message.Data = returnData.Data().(string)
// 	return Message
// }

// corresponding to the "Save" button on the frontend
func (a *App) FileSave(name string, rawJson string, skip bool) Message {
	var msg Message
	res := name
	if !skip {
		res = FileSaveDialog(a, name, "*.epub")
	}
	if res == "" {
		msg.Code = -1
		msg.Msg = "cancel"
		return msg
	}
	logger.Info("Saving Epub file", logging.RunFuncName())
	future, err := a.core.agt.Exec("writer", res, rawJson)
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		msg.Code = 1
		msg.Msg = err.Error()
		return msg
	}
	ctxGet, cancelGet := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelGet()
	returnData, err := future.Get(ctxGet)
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		msg.Code = 1
		msg.Msg = err.Error()
		return msg
	}
	if returnData.Err() != nil {
		msg.Code = 1
		msg.Msg = returnData.Err().Error()
		logger.Error(returnData.Err().Error(), logging.RunFuncName())
		return msg
	}
	hash, err := utils.GenerateHash([]byte(res))
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	a.resourceFold = hash
	msg.Code = 0
	msg.Data = res
	msg.Msg = "success"
	return msg
}

// get the static resources list of the loacl server
func (a *App) GetStaticResources() Message {
	var msg Message
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("http://127.0.0.1:7288/" + a.resourceFold)
	if err != nil {
		logger.Fatal(err.Error(), logging.RunFuncName())
		msg.Code = 1
		msg.Msg = err.Error()
		return msg
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	msg.Code = 0
	msg.Msg = "success"
	msg.Data = string(body)
	return msg
}

// corresponding to the "Delete" button on the "insert picture"
func (a *App) FileDelete(name string) Message {
	path := filepath.Join(a.core.execPath, "resources", a.resourceFold, name)
	msg := new(Message)
	err := os.Remove(path)
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return *msg
	}
	msg.Code = 0
	msg.Msg = "success"
	msg.Data = ""
	return *msg
}

// corresponding to the "Upload" button on the "insert picture"
func (a *App) ImageUpload() ImageFIle {
	var img ImageFIle
	var AllowExt = a.core.config.StaticResource.AllowExt
	conSuffix := strings.Join(AllowExt, ";")
	conSuffix = strings.ReplaceAll(conSuffix, ".", "*")
	path := FileOpenDialog(a, "Image File", conSuffix)
	if path == "" {
		logger.Info("Cancel upload", logging.RunFuncName())
		img.Code = -1
		return img
	}
	_, name := filepath.Split(path)
	suffix := filepath.Ext(name)
	set := make(map[string]struct{})
	for _, v := range AllowExt {
		set[v] = struct{}{}
	}
	if _, ok := set[suffix]; !ok {
		logger.Warning("The file type is not allowed", logging.RunFuncName())
		img.Code = -1
		return img
	}
	imgFp, err := os.Open(path)
	if err != nil {
		img.Code = 1
		logger.Error(err.Error(), logging.RunFuncName())
		return img
	}
	imgData, _ := io.ReadAll(imgFp)
	defer imgFp.Close()
	h := md5.New()
	h.Write(imgData)
	id := hex.EncodeToString(h.Sum(nil))[8:24]

	fp, err := os.Create(filepath.Join(a.core.execPath, "resources", a.resourceFold, id+".jpg"))
	if err != nil {
		img.Code = 1
		logger.Error(err.Error(), logging.RunFuncName())
		return img
	}
	defer fp.Close()
	fp.Write(imgData)
	img.Code = 0
	img.Name = name
	img.Id = id
	return img
}

func (a *App) LoadImage(data string) Message {
	var msg Message
	imgData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	h := md5.New()
	h.Write(imgData)
	id := hex.EncodeToString(h.Sum(nil))[8:24]
	imagePath := filepath.Join(a.core.execPath, "resources", id+".jpg")
	if _, err := os.Stat(imagePath); err == nil {
		msg.Code = -1
		msg.Msg = "resources folder already exists"
		return msg
	}
	fp, err := os.Create(imagePath)
	if err != nil {
		msg.Code = 1
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	defer fp.Close()
	fp.Write(imgData)
	msg.Code = 0
	msg.Msg = "success"
	msg.Data = id
	return msg
}

// corresponding to the "Cover Upload" on the "Book info"
func (a *App) OpenImage() Message {
	var msg Message
	path := FileOpenDialog(a, "Image File", "*jpg;*png;*jpeg;*bmp")
	if path == "" {
		msg.Code = 0
		msg.Msg = "cancel"
		return msg
	}
	img, err := os.Open(path)
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	defer img.Close()
	imgData, _ := io.ReadAll(img)
	msg.Code = 0
	msg.Msg = "success"
	msg.Data = base64.StdEncoding.EncodeToString(imgData)
	return msg
}

// get the base64 string of the image
func (a *App) GetImageData(filename string) Message {
	imagePath := filepath.Join(a.core.execPath, "resources", filename)
	var msg Message
	fs, err := os.Open(imagePath)
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	defer fs.Close()
	imgData, _ := io.ReadAll(fs)
	msg.Code = 0
	msg.Msg = "success"
	msg.Data = base64.StdEncoding.EncodeToString(imgData)
	return msg
}

// get the configure
func (a *App) GetConfig(sector string, key string) Message {
	var msg Message
	if sector == "" || key == "" {
		msg.Code = -1
		msg.Msg = "sector or key is empty"
		return msg
	}
	value, err := a.core.cm.GetConfigByKey(sector, key)
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	msg.Code = 0
	msg.Msg = "success"
	msg.Data = value
	return msg
}

func (a *App) SetConfig(sector string, key string, value string) Message {
	var msg Message
	if sector == "" || key == "" {
		msg.Code = -1
		msg.Msg = "sector or key is empty"
		return msg
	}
	err := a.core.cm.SetConfig(sector, key, value)
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	msg.Code = 0
	msg.Msg = "success"
	return msg
}

func (a *App) GetVersion() {
	info := a.core.cm.GetInfo()
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "NovelMaker",
		Message: fmt.Sprintf("Version: %-10s\nCommit: %-10s\nBuild Time: %-10s\nGolang: %-10s\nWails: %-10s\nNodejs: %-10s\nCopyright: %-10s", info["version"], info["commit"], info["buildTime"], "v1.20.3", "v2.8.1", "22.11.0", "©2024 Chenxi Guan"),
	})
}

func (a *App) ImageDownload(url string) Message {
	var msg Message
	downloader := server.NewImageDownloader(filepath.Join(a.core.execPath, "resources", a.resourceFold), a.core.config.Dowload.Timeout)
	err := server.ProcessDownload(downloader, url)
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	msg.Code = 0
	msg.Msg = "success"
	return msg
}

func (a *App) Trace(source, stack string) {
	logger.Trace(source, stack)
}
