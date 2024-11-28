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

var logger = logging.NewLog(logging.FatalLevel, true)

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

func InSlice(itmes []string, item string) bool {
	return sort.SearchStrings(itmes, item) < len(itmes)
}

func ErrorHandler(e error, w http.ResponseWriter) {
	if e != nil {
		logger.Fatal(e.Error(), logging.RunFuncName())
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func DirList(w http.ResponseWriter, r *http.Request, f http.File) {
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
		url := url.URL{Path: name}
		staticResource.FileList = append(staticResource.FileList, url.String())
	}
	rawJson, _ := json.Marshal(staticResource)
	io.WriteString(w, string(rawJson))
}

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
			DirList(w, r, f)
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

func StaticSevice(port string, execPath string) {
	logger.Info("Static resource service started", logging.RunFuncName())
	http.HandleFunc("/", ResourceHandler(execPath))
	http.ListenAndServe(port, nil)
}
