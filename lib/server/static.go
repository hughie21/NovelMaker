// Description: Local services for storing media resources
// Author: Hughie21
// Date: 2024-12-20
// license that can be found in the LICENSE file.
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/hughie21/NovelMaker/lib/logging"
)

var logger *logging.Log

type StaticResource struct {
	Code     int
	FileList []string
}

var extensionToContentType = map[string]string{
	".html": "text/html; charset=utf-8",
	".css":  "text/css; charset=utf-8",
	".js":   "application/javascript",
	".xml":  "text/xml; charset=utf-8",
	".jpg":  "image/jpeg",
}

// Determine if an element exists in a slice
func InSlice(itmes []string, item string) bool {
	return sort.SearchStrings(itmes, item) < len(itmes)
}

// Exception handling
func ErrorHandler(e error, w http.ResponseWriter) {
	if e != nil {
		logger.Fatal(e.Error(), logging.RunFuncName())
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

// Print Catalog Information
func DirList(w http.ResponseWriter, r *http.Request, f http.File, parentDir string) {
	dirs, err := f.Readdir(-1)
	if err != nil {
		logger.Fatal(err.Error(), logging.RunFuncName())
		fmt.Println(w, http.StatusInternalServerError)
		return
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	staticResource := new(StaticResource)
	staticResource.Code = 0

	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		mdUrl := url.URL{Path: filepath.Join(parentDir, name)}
		mdPath, _ := url.PathUnescape(mdUrl.String())
		staticResource.FileList = append(staticResource.FileList, mdPath)
	}
	rawJson, _ := json.Marshal(staticResource)
	io.WriteString(w, string(rawJson))
}

// Static resource handler
func ResourceHandler(execPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ResourcePath := filepath.Join(execPath, "resources")
		path := ResourcePath + r.URL.Path
		f, err := os.Open(path)
		if err != nil {
			logger.Error(err.Error(), logging.RunFuncName())
			return
		}
		defer f.Close()

		d, err := f.Stat()
		if err != nil {
			logger.Error(err.Error(), logging.RunFuncName())
			return
		}

		if d.IsDir() {
			DirList(w, r, f, r.URL.Path)
			return
		}

		data, err := io.ReadAll(f)
		if err != nil {
			logger.Error(err.Error(), logging.RunFuncName())
			return
		}

		ext := filepath.Ext(path)
		if contentType := extensionToContentType[ext]; contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.Header().Set("Content-Length", strconv.FormatInt(d.Size(), 10))
		w.Write(data)
	}
}

// Start the static resource service
func StaticSevice(port string, execPath string, log *logging.Log) {
	logger = log
	logger.Info("Static resource service started", logging.RunFuncName())
	http.HandleFunc("/", ResourceHandler(execPath))
	http.ListenAndServe(port, nil)
}
