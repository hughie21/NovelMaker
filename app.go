/*
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-09-16
@Description: This is the Go function that frontend can call for.
*/

package main

import (
	epubMaker "NovelMaker/lib/epub"
	logging "NovelMaker/logging"
	Manager "NovelMaker/manager"
	sys "NovelMaker/sys"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

	"bytes"
	"context"
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
)

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

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	logger.Info("App started", logging.RunFuncName())
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	logger.Info("App shutdown", logging.RunFuncName())
	err := logger.LogOutPut(execPath)
	if err != nil {
		sys.ShowMessage("Error when writing log: ", err.Error(), "error")
		panic(err)
	}
	err = cm.SaveConfig()
	if err != nil {
		sys.ShowMessage("Error when writing config: ", err.Error(), "error")
		panic(err)
	}
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

// when user open the epmb file directly, this function will be called
func (a *App) DirectLoading() Message {
	var msg Message
	fmt.Println(Args)
	if Args == "" {
		msg.Code = -1
		msg.Msg = "no args"
		return msg
	}
	dataStruct, err := epubMaker.Load(Args)
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	jsonData := epubMaker.Dump(&dataStruct)
	msg.Code = 0
	msg.Msg = "success"
	msg.Data = jsonData
	return msg
}

// corresponding to the "Open" button on the frontend
func (a *App) FileOpen() Message {
	var Message Message
	res := FileOpenDialog(a, "Epub File", "*.epub")
	// dataStruct, err := epubMaker(res)
	if res == "" {
		Message.Code = -1
		Message.Msg = "cancel"
		return Message
	}
	reader, err := epubMaker.NewReader(res, filepath.Join(execPath, "tmp"))
	if err != nil {
		Message.Code = 1
		Message.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return Message
	}
	err = reader.Read()
	if err != nil {
		Message.Code = 1
		Message.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return Message
	}
	err = reader.Pharse()
	if err != nil {
		Message.Code = 1
		Message.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return Message
	}
	reader.Close()
	jsonData := epubMaker.Dump(&reader.JsonData)
	Message.Code = 0
	Message.Msg = res
	Message.Data = jsonData
	return Message
}

func (a *App) FileImport() Message {
	var Message Message
	res := FileOpenDialog(a, "Text File", "*.md;*.txt")
	if res == "" {
		Message.Code = -1
		Message.Msg = "cancel"
		return Message
	}
	fp, err := os.Open(res)
	if err != nil {
		Message.Code = 1
		Message.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return Message
	}
	defer fp.Close()
	content, _ := io.ReadAll(fp)
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
	}
	Message.Code = 0
	Message.Msg = "success"
	Message.Data = buf.String()
	return Message
}

// corresponding to the "Save" button on the frontend
func (a *App) FileSave(name string, rawJson string, skip bool) Message {
	res := name
	if !skip {
		res = FileSaveDialog(a, name, "*.no")
	}

	var msg Message
	if res == "" {
		msg.Code = -1
		msg.Msg = "cancel"
		return msg
	}
	var JsonStruct epubMaker.JsonData
	epubMaker.LoadJson([]byte(rawJson), &JsonStruct)
	err := epubMaker.SaveToFile(&JsonStruct, res)
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		logger.Error(err.Error(), logging.RunFuncName())
		return msg
	}
	msg.Code = 0
	msg.Msg = "success"
	msg.Data = res
	return msg
}

// get the static resources list of the loacl server
func (a *App) GetStaticResources() Message {
	var msg Message
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("http://127.0.0.1:7288/")
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
	path := filepath.Join(execPath, "resources", name)
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
	var AllowExt = config.StaticResource.AllowExt
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

	fp, err := os.Create(filepath.Join(execPath, "resources", id+".jpg"))
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
	imagePath := filepath.Join(execPath, "resources", id+".jpg")
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

// corresponding to the "Export" button on the frontend
func (a *App) Publish(name string, rawJson string) Message {
	var JsonStruct epubMaker.JsonData
	var msg Message
	path := FileSaveDialog(a, name, "*.epub")
	if path == "" {
		msg.Code = -1
		return msg
	}
	epubMaker.LoadJson([]byte(rawJson), &JsonStruct)
	logger.Info("Start to export to EPUB", logging.RunFuncName())
	// tmpPath := "epubMaker.(JsonStruct)"
	// e := epubMaker.WriteEpub(tmpPath, path)
	// if e != nil {
	// 	e = os.RemoveAll(tmpPath)
	// 	msg.Code = 1
	// 	msg.Msg = e.Error()
	// 	logger.Error(e.Error(), logging.RunFuncName())
	// 	return msg
	// }
	// e = os.RemoveAll(tmpPath)
	// if e != nil {
	// 	logger.Error(e.Error(), logging.RunFuncName())
	// }
	// msg.Code = 0
	// msg.Msg = "success"
	// return msg
	writer := epubMaker.NewWriter(path, filepath.Join(execPath, "tmp"), &JsonStruct)
	err := writer.Write()
	defer writer.Close()
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

// get the base64 string of the image
func (a *App) GetImageData(filename string) Message {
	imagePath := filepath.Join(execPath, "resources", filename)
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
	value, err := cm.GetConfigByKey(sector, key)
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
	err := cm.SetConfig(sector, key, value)
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

func (a *App) ImageDownload(url string) Message {
	var msg Message
	downloader := Manager.NewImageDownloader(filepath.Join(execPath, "resources"), config.Dowload.Timeout)
	err := Manager.ProcessDownload(downloader, url)
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
