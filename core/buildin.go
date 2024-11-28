// Description: This is the two major inner plugins of NovelMaker, which are used to process epub files.
// Author: Hughie21
// Date: 2024-11-21
// license that can be found in the LICENSE file.

package core

import (
	"context"
	"errors"

	"github.com/hughie21/NovelMaker/lib/epub"
	"github.com/hughie21/NovelMaker/lib/html"
)

type (
	EpubWriter struct {
		PluginInfo
		tempDir string
	}
	EpubReader struct {
		PluginInfo
		tempDir    string
		extionsion map[string]html.TagParser
	}
)

var (
	epubDefaultInfo = PluginInfo{
		Name:     "EpubProcessor",
		Version:  "0.0.1",
		Type:     "epub",
		Priority: 999,
		Author:   "Hughie",
		Email:    "",
		Source:   "github.com/hughie21/NovelMaker/lib/epub",
	}
)

func NewEpubWriter(tempDir string) *EpubWriter {
	return &EpubWriter{
		PluginInfo: epubDefaultInfo,
		tempDir:    tempDir,
	}
}

func NewEpubReader(tempDir string) *EpubReader {
	return &EpubReader{
		PluginInfo: epubDefaultInfo,
		tempDir:    tempDir,
		extionsion: map[string]html.TagParser{},
	}
}

func (e *EpubWriter) Run(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("insufficient arguments")
	}
	targetPath, ok := args[0].(string)
	if !ok {
		return nil, errors.New("first argument must be a string")
	}

	rawData, ok := args[1].(string)
	if !ok {
		return nil, errors.New("second argument must be a string")
	}
	var jsonData epub.JsonData
	epub.LoadJson([]byte(rawData), &jsonData)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		writer := epub.NewWriter(targetPath, e.tempDir, &jsonData)
		defer writer.Close()
		if err := writer.Write(); err != nil {
			return 0, err
		}
		return 1, nil
	}
}

func (e *EpubReader) Run(ctx context.Context, args ...interface{}) (interface{}, error) {
	targetPath, ok := args[0].(string)
	if !ok {
		return nil, errors.New("first argument must be a string")
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		reader, err := epub.NewReader(targetPath, e.tempDir)
		defer reader.Close()
		if err != nil {
			return nil, err
		}
		err = reader.Read()
		if err != nil {
			return nil, err
		}
		err = reader.Pharse(e.extionsion)
		if err != nil {
			return nil, err
		}
		result := epub.Dump(&reader.JsonData)

		return result, nil
	}
}
