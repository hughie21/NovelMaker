// Description: Project's public functions and constants
// Author: Hughie21
// Date: 2024-12-20
// license that can be found in the LICENSE file.
package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hughie21/NovelMaker/lib/html"
	"golang.org/x/exp/rand"
)

func GetFileData(path string) []byte {
	fs, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	rawData, err := io.ReadAll(fs)
	if err != nil {
		panic(err)
	}

	return rawData
}

func GenerateHash(data []byte) (string, error) {
	h := sha1.New()
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil))[:20], nil
}

func CombineMap(m1, m2 map[string]html.TagParser) map[string]html.TagParser {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset)-1)]
	}
	return string(b)
}

func Contains(arr []string, target string) bool {
	str := strings.Join(arr, " ")
	return strings.Contains(str, target)
}

func NormalizeHTML(html string) string {
	// change the tag name to lowercase
	html = regexp.MustCompile(`<\/?([A-Z][A-Z0-9]*)\b[^>]*>`).ReplaceAllStringFunc(html, func(match string) string {
		return strings.ToLower(match)
	})

	// makesure the self-closing tags are closed properly
	selfClosingTags := map[string]bool{
		"area": true, "base": true, "br": true, "col": true, "embed": true, "hr": true, "img": true, "input": true, "link": true, "meta": true, "param": true, "source": true, "track": true, "wbr": true,
	}
	html = regexp.MustCompile(`<([a-z]+)([^>]*)\/?>`).ReplaceAllStringFunc(html, func(match string) string {
		re := regexp.MustCompile(`<([a-z]+)([^>]*)\/?>`)
		submatches := re.FindStringSubmatch(match)
		tagName := submatches[1]
		attributes := submatches[2]
		if selfClosingTags[tagName] {
			return `<` + tagName + attributes + `></` + tagName + `>`
		} else {
			return match
		}
	})

	return html
}

var FileSuffix = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/gif":       ".gif",
	"image/svg":       ".svg",
	"image/tiff":      ".tiff",
	"image/bmp":       ".bmp",
	"image/webp":      ".webp",
	"image/x-icon":    ".ico",
	"audio/mpeg":      ".mp3",
	"audio/ogg":       ".ogg",
	"audio/wav":       ".wav",
	"audio/flac":      ".flac",
	"audio/aac":       ".aac",
	"audio/x-ms-wma":  ".wma",
	"audio/x-ms-wmv":  ".wmv",
	"audio/x-ms-wav":  ".wav",
	"audio/x-ms-mp3":  ".mp3",
	"audio/x-ms-mp4":  ".mp4",
	"audio/x-ms-flac": ".flac",
	"audio/x-ms-aac":  ".aac",
	"audio/x-ms-aiff": ".aiff",
	"audio/x-ms-aif":  ".aif",
	"audio/x-ms-aifc": ".aifc",
	"audio/x-ms-m4a":  ".m4a",
	"audio/x-ms-m4b":  ".m4b",
	"audio/x-ms-m4p":  ".m4p",
	"audio/x-ms-m4r":  ".m4r",
}

func GetKeyByValue(value string) (string, bool) {
	for k, v := range FileSuffix {
		if v == value {
			return k, true
		}
	}
	return "", false
}

type marshal struct {
}

// modify the json marshal method
// when dump the struct to json
// it will turn the empty field into []
// instead of null
func Marshal(obj interface{}) string {
	v := reflect.ValueOf(obj)
	return new(marshal).parse(v)
}

// deal with the normal field
func (that *marshal) parse(v reflect.Value) string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.String:
		return `"` + v.String() + `"`
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Struct:
		return that.parseStruct(v, v.Type())
	case reflect.Map:
		return that.parseMap(v, v.Type())
	case reflect.Slice:
		return that.parseSlice(v, v.Type())
	case reflect.Interface:
		return that.parse(reflect.ValueOf(v.Interface()))
	}
	return "null"
}

// deal with the struct
func (that *marshal) parseStruct(v reflect.Value, t reflect.Type) string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	var str = make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fieldName := field.Name
		tag := field.Tag.Get("json")
		if tag != "" {
			f := strings.Split(tag, ",")
			fieldName = f[0]
		}
		if fieldName == "-" {
			continue
		}
		str[i] = `"` + fieldName + `":` + that.parse(value)
	}
	return "{" + strings.Join(str, ",") + "}"
}

// deal with the slice
func (that *marshal) parseSlice(v reflect.Value, _ reflect.Type) string {
	var str = make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		str[i] = that.parse(v.Index(i))
	}
	return `[` + strings.Join(str, ",") + "]"
}

// deal with the map
func (that *marshal) parseMap(v reflect.Value, _ reflect.Type) string {
	var str = make([]string, v.Len())
	m := v.MapRange()
	var i int
	for m.Next() {
		str[i] = `"` + m.Key().String() + `":` + that.parse(m.Value())
		i++
	}
	return "{" + strings.Join(str, ",") + "}"
}
