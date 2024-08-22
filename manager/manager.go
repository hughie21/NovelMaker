/*
@Author: Hughie
@CreateTime: 2024-7-20
@LastEditors: Hughie
@LastEditTime: 2024-08-1
@Description: This is the static resource manager of the program, use to display the static resource on the frontend
*/

package manager

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
)

type StaticResource struct {
	Code     int
	FileList []string
}

var AllowExt = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

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
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func DirList(w http.ResponseWriter, r *http.Request, f http.File) {
	dirs, err := f.Readdir(-1)
	if err != nil {
		fmt.Println(w, http.StatusInternalServerError)
		return
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	// fmt.Fprintf(w, "<h1> This is the static resource service of the NovelMaker </h1>\n")
	// fmt.Fprintf(w, "<pre>\n")
	staticResource := new(StaticResource)
	staticResource.Code = 0
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		url := url.URL{Path: name}
		// fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), name)
		staticResource.FileList = append(staticResource.FileList, url.String())
	}
	rawJson, _ := json.Marshal(staticResource)
	io.WriteString(w, string(rawJson))
}

func ResourceHandler(w http.ResponseWriter, r *http.Request) {
	ResourcePath := "./resources/"
	path := ResourcePath + r.URL.Path
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	if d.IsDir() {
		DirList(w, r, f)
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	ext := filepath.Ext(path)
	if contentType := extensionToContentType[ext]; contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.Header().Set("Content-Length", strconv.FormatInt(d.Size(), 10))
	w.Write(data)
}

func StaticSevice() {
	http.HandleFunc("/", ResourceHandler)

	http.ListenAndServe(":7288", nil)
}
