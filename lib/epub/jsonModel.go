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

	utils "NovelMaker/lib/utils"
)

type File struct {
	Name string
	Data []byte
	Type string
}

type J_cover struct {
	Name string `json:"id"`
	Data string `json:"data"`
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
	Cover        J_cover  `json:"cover"`
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
	jsonData := utils.Marshal(data)
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
