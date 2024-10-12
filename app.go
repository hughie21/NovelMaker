/*
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-09-16
@Description: This is the Go function that frontend can call for.
*/

package main

import (
	epubMaker "NovelMaker/epub"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
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

var logFileName string

// App struct
type App struct {
	ctx context.Context
}

func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func LogOutPut(msg string, funcName string) {
	if logFileName == "" {
		todaystr := time.Now().Format("2006-01-02")
		todayint, _ := time.ParseInLocation("2006-01-02", todaystr, time.Local)
		logFileName = strconv.FormatInt(todayint.Unix(), 10) + ".log"
	}
	fileLogger := logger.NewFileLogger(filepath.Join("./log", logFileName))
	fileLogger.Error(fmt.Sprintf("%s -> %s", funcName, msg))
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// return the raw data of file
func (a *App) Fr(path string) string {
	fp, err := os.Open(path)
	if err != nil {
		LogOutPut(err.Error(), runFuncName())
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
		LogOutPut(err.Error(), runFuncName())
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
	res := FileOpenDialog(a, "Empb File", "*.no")
	dataStruct, err := epubMaker.Load(res)
	if res == "" {
		Message.Code = -1
		Message.Msg = "cancel"
		return Message
	}
	if err != nil {
		LogOutPut(err.Error(), runFuncName())
		Message.Code = 1
		Message.Msg = err.Error()
		return Message
	}
	jsonData := epubMaker.Dump(&dataStruct)
	Message.Code = 0
	Message.Msg = res
	Message.Data = jsonData
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
		LogOutPut(err.Error(), runFuncName())
		return msg
	}
	msg.Code = 0
	msg.Msg = "success"
	msg.Data = res
	return msg
}

// get the static resources list of the loacl server
func (a *App) GetStaticResources() string {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("http://127.0.0.1:7288/")
	if err != nil {
		LogOutPut(err.Error(), runFuncName())
		return "{code: 1, msg: " + err.Error() + "}"
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

// corresponding to the "Delete" button on the "insert picture"
func (a *App) FileDelete(name string) Message {
	path := filepath.Join("./resources", name)
	msg := new(Message)
	err := os.Remove(path)
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		LogOutPut(err.Error(), runFuncName())
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
	path := FileOpenDialog(a, "Image File", "*jpg;*png;*jpeg;*bmp")
	if path == "" {
		img.Code = -1
		return img
	}
	_, name := filepath.Split(path)
	imgFp, err := os.Open(path)
	if err != nil {
		img.Code = 1
		LogOutPut(err.Error(), runFuncName())
		return img
	}
	imgData, _ := io.ReadAll(imgFp)
	defer imgFp.Close()
	h := md5.New()
	h.Write(imgData)
	id := hex.EncodeToString(h.Sum(nil))[8:24]

	fp, err := os.Create(filepath.Join("./resources", id+".jpg"))
	if err != nil {
		img.Code = 1
		LogOutPut(err.Error(), runFuncName())
		return img
	}
	defer fp.Close()
	fp.Write(imgData)
	img.Code = 0
	img.Name = name
	img.Id = id
	return img
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
		LogOutPut(err.Error(), runFuncName())
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
	tmpPath := epubMaker.FormXML(JsonStruct)
	e := epubMaker.WriteEpub(tmpPath, path)
	if e != nil {
		e = os.RemoveAll(tmpPath)
		msg.Code = 1
		msg.Msg = e.Error()
		LogOutPut(e.Error(), runFuncName())
		return msg
	}
	e = os.RemoveAll(tmpPath)
	if e != nil {
		LogOutPut(e.Error(), runFuncName())
	}
	msg.Code = 0
	msg.Msg = "success"
	return msg
}

// get the base64 string of the image
func (a *App) GetImageData(filename string) Message {
	imagePath := filepath.Join("./resources", filename)
	var msg Message
	fs, err := os.Open(imagePath)
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		LogOutPut(err.Error(), runFuncName())
		return msg
	}
	defer fs.Close()
	imgData, _ := io.ReadAll(fs)
	msg.Code = 0
	msg.Msg = "success"
	msg.Data = base64.StdEncoding.EncodeToString(imgData)
	return msg
}
