/*
@Author: Hughie
@CreateTime: 2024-10-31
@LastEditors: Hughie
@LastEditTime: 2024-11-1
@Description: Download middleware
*/

package manager

import (
	"NovelMaker/logging"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Middleware interface {
	Download(url string)
	Write()
	GetError() error
}

type ImageDownloader struct {
	Url     string
	timeout time.Duration
	Body    io.ReadCloser
	Path    string
	Name    string
	Err     error
}

func ProcessDownload(m Middleware, url string) error {
	m.Download(url)
	m.Write()
	if err := m.GetError(); err != nil {
		return err
	}
	return nil
}

func NewImageDownloader(path string, delay int) *ImageDownloader {
	return &ImageDownloader{
		timeout: time.Duration(delay) * time.Second,
		Url:     "",
		Body:    nil,
		Path:    path,
		Name:    "",
		Err:     nil,
	}
}

func (downloader *ImageDownloader) Download(url string) {
	client := http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: downloader.timeout,
			}).Dial,
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		downloader.Err = err
	}
	downloader.Body = resp.Body
	downloader.generateName()
}

func (downloader *ImageDownloader) generateName() {
	h := md5.New()
	h.Write([]byte(downloader.Url))
	id := hex.EncodeToString(h.Sum(nil))[8:24]
	downloader.Name = id + ".jpg"
}

func (downloader *ImageDownloader) Write() {
	if downloader.Err != nil {
		return
	}
	file, err := os.Create(filepath.Join(downloader.Path, downloader.Name))
	if err != nil {
		downloader.Err = err
		logger.Error(err.Error(), logging.RunFuncName())
	}
	defer file.Close()
	_, err = io.Copy(file, downloader.Body)
	downloader.Body.Close()
	if err != nil {
		downloader.Err = err
		logger.Error(err.Error(), logging.RunFuncName())
	}
}

func (downloader *ImageDownloader) GetError() error {
	return downloader.Err
}