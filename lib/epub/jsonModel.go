// Description: Data structure definition for front-end and back-end data exchange
// Author: Hughie21
// Date: 2024-11-29
// license that can be found in the LICENSE file.

package epub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"os"

	"github.com/hughie21/NovelMaker/lib/utils"
)

// memory file
type File struct {
	Name string
	Data []byte
	Type string
}

// The cover image of the book
type J_cover struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// The metadata of the book
type J_MetaData struct {
	Title        string   `json:"title"`
	Creator      []string `json:"creator"`
	Identifier   string   `json:"identifier"`
	Language     string   `json:"language"`
	Contributers []string `json:"contributors"`
	Description  string   `json:"description"`
	Publisher    string   `json:"publisher"`
	Subject      []string `json:"subject"`
	Cover        J_cover  `json:"cover"`
	Meta         J_meta   `json:"meta"`
}

// The book display setting
type J_meta struct {
	TextDir     string `json:"textDirection"`
	Layout      string `json:"layout"`
	Flow        string `json:"flow"`
	Orientation string `json:"orientation"`
	Spread      string `json:"spread"`
	Proportions string `json:"proportions"`
}

// The navigation of the book
type J_nav struct {
	Id    string  `json:"id"`
	Label string  `json:"label"`
	Href  string  `json:"href"`
	Child []J_nav `json:"children"`
}

// the item in the book
type Resource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
}

// The book data structure for front-end and back-end data exchange
type JsonData struct {
	MetaData  J_MetaData `json:"metadata"`
	Content   string     `json:"content"`
	Resources []Resource `json:"resources"`
	Nav       []J_nav    `json:"toc"`
}

// return the json string of the data
func Dump(data *JsonData) string {
	jsonData := utils.Marshal(data)
	return string(jsonData)
}

// save the data to the file
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

// load the data from the file
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

// load the data from the json string
func LoadJson(RawData []byte, Mapping *JsonData) {
	json.Unmarshal(RawData, Mapping)
}
