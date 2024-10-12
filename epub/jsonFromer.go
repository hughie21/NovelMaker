/*
@Author: Hughie
@CreateTime: 2024-8-2
@LastEditors: Hughie
@LastEditTime: 2024-08-7
@Description: The definition of the json Format and the json marshal method
*/

package epub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

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
func (this *marshal) parse(v reflect.Value) string {
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
		return this.parseStruct(v, v.Type())
	case reflect.Map:
		return this.parseMap(v, v.Type())
	case reflect.Slice:
		return this.parseSlice(v, v.Type())
	case reflect.Interface:
		return this.parse(reflect.ValueOf(v.Interface()))
	}
	return "null"
}

// deal with the struct
func (this *marshal) parseStruct(v reflect.Value, t reflect.Type) string {
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
		str[i] = `"` + fieldName + `":` + this.parse(value)
	}
	return "{" + strings.Join(str, ",") + "}"
}

// deal with the slice
func (this *marshal) parseSlice(v reflect.Value, t reflect.Type) string {
	var str = make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		str[i] = this.parse(v.Index(i))
	}
	return `[` + strings.Join(str, ",") + "]"
}

// deal with the map
func (this *marshal) parseMap(v reflect.Value, t reflect.Type) string {
	var str = make([]string, v.Len())
	m := v.MapRange()
	var i int
	for m.Next() {
		str[i] = `"` + m.Key().String() + `":` + this.parse(m.Value())
		i++
	}
	return "{" + strings.Join(str, ",") + "}"
}

type J_meta struct {
	Name    string `json:"id"`
	Content string `json:"content"`
}

type J_MetaData struct {
	Title        string   `json:"title"`
	Creator      []string `json:"creator"`
	Identifier   string   `json:"identifier"`
	Language     string   `json:"language"`
	Contributers []string `json:"contributors"`
	Description  string   `json:"description"`
	Publisher    string   `json:"publisher"`
	Subject      []string `json:"subject"`
	Date         string   `json:"date"`
	Meta         []J_meta `json:"meta"`
}

type J_nav struct {
	Id    string  `json:"id"`
	Label string  `json:"label"`
	Href  string  `json:"href"`
	Child []J_nav `json:"children"`
}

type Resource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
}

type JsonData struct {
	MetaData  J_MetaData `json:"metadata"`
	Content   string     `json:"content"`
	Resources []Resource `json:"resources"`
	Nav       []J_nav    `json:"toc"`
}

func Dump(data *JsonData) string {
	jsonData := Marshal(data)
	return string(jsonData)
}

func SaveToFile(jsonData *JsonData, filePath string) error {
	var SaveFile = func(path string, data []byte) error {
		return os.WriteFile(path, data, 0644)
	}
	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(jsonData)
	SaveFile(filePath, buf.Bytes())
	return err
}

func Load(filePath string) (JsonData, error) {
	File, err := os.Open(filePath)
	if err != nil {
		return JsonData{}, err
	}
	defer File.Close()
	rawData, _ := io.ReadAll(File)
	decoder := gob.NewDecoder(bytes.NewReader(rawData))
	var data JsonData
	decoder.Decode(&data)
	return data, nil
}

func LoadJson(RawData []byte, Mapping *JsonData) {
	json.Unmarshal(RawData, Mapping)
}
