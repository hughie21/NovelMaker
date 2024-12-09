package epub

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hughie21/NovelMaker/lib/html"
	"github.com/hughie21/NovelMaker/lib/utils"

	"github.com/beevik/etree"
)

var (
	Image = []string{"image/jpeg", "image/png", "image/gif", "image/svg", "image/tiff", "image/bmp", "image/webp", "image/x-icon"}
	Audio = []string{"audio/mpeg", "audio/ogg", "audio/wav", "audio/flac", "audio/aac", "audio/x-ms-wma", "audio/x-ms-wmv", "audio/x-ms-wav", "audio/x-ms-mp3", "audio/x-ms-mp4", "audio/x-ms-flac", "audio/x-ms-aac", "audio/x-ms-aiff", "audio/x-ms-aif", "audio/x-ms-aifc", "audio/x-ms-m4a", "audio/x-ms-m4b", "audio/x-ms-m4p", "audio/x-ms-m4r"}
)

/*
*
* File structure
*
 */

type Reader struct {
	JsonData
	Package      *etree.Document
	Images       []File
	tempDir      string
	targetPath   string
	resourcePath string
}

func NewReader(targetPath string, tempDir string) (*Reader, error) {
	if !utils.PathExists(targetPath) {
		return nil, errors.New("File does not exist")
	}
	_, name := filepath.Split(targetPath)
	err := Dezip(targetPath, tempDir)
	if err != nil {
		return nil, err
	}
	r := &Reader{
		tempDir:      filepath.Join(tempDir, fmt.Sprintf("~%s", name)),
		targetPath:   targetPath,
		resourcePath: filepath.Join(filepath.Join(tempDir, ".."), "resources", utils.GenerateHash([]byte(targetPath))),
	}
	return r, nil
}

func (r *Reader) checkEpub() error {
	if !utils.PathExists(filepath.Join(r.tempDir, "META-INF")) {
		return errors.New("META-INF not found")
	}
	if !utils.PathExists(filepath.Join(r.tempDir, "mimetype")) {
		return errors.New("mimetype not found")
	}
	fs, _ := os.Open(filepath.Join(r.tempDir, "mimetype"))
	defer fs.Close()
	mimetype := make([]byte, 20)
	fs.Read(mimetype)
	if string(mimetype) != "application/epub+zip" {
		return errors.New("mimetype not match")
	}
	return nil
}

func (r *Reader) Read() error {
	err := r.checkEpub()
	if err != nil {
		return err
	}
	// read container.xml
	container, err := os.ReadFile(filepath.Join(r.tempDir, "META-INF", "container.xml"))
	if err != nil {
		return err
	}
	// Pharse container.xml
	containerDoc := etree.NewDocument()
	err = containerDoc.ReadFromBytes(container)
	if err != nil {
		return err
	}
	optPath := containerDoc.FindElement("//rootfile[@media-type='application/oebps-package+xml']").SelectAttrValue("full-path", "")
	// read content.opt
	opt, err := os.ReadFile(filepath.Join(r.tempDir, optPath))
	if err != nil {
		return err
	}
	packageDoc := etree.NewDocument()
	err = packageDoc.ReadFromBytes(opt)
	if err != nil {
		return err
	}
	r.tempDir = filepath.Join(r.tempDir, filepath.Dir(optPath))
	r.Package = packageDoc
	return nil
}

func (r *Reader) moveImage() error {
	if !utils.PathExists(r.resourcePath) {
		err := os.Mkdir(r.resourcePath, os.ModePerm)
		if err != nil {
			return errors.New("create resource path error: " + err.Error())
		}
	}
	for _, image := range r.Images {
		ts, err := os.Create(filepath.Join(r.resourcePath, image.Name))
		if err != nil {
			return errors.New("create file error by: " + filepath.Join(r.resourcePath, image.Name))
		}
		defer ts.Close()
		ts.Write(image.Data)
	}
	return nil
}

func (r *Reader) Pharse(extension map[string]html.TagParser) error {
	// Pharse package document
	packageDoc := r.Package

	idElem := packageDoc.FindElement("//dc:identifier")
	if idElem != nil {
		r.MetaData.Identifier = idElem.Text()
	} else {
		r.MetaData.Identifier = ""
	}

	titleElem := packageDoc.FindElement("//dc:title")
	if titleElem != nil {
		r.MetaData.Title = titleElem.Text()
	} else {
		r.MetaData.Title = ""
	}

	langElem := packageDoc.FindElement("//dc:language")
	if langElem != nil {
		r.MetaData.Language = langElem.Text()
	} else {
		r.MetaData.Language = "en"
	}

	pubElem := packageDoc.FindElement("//dc:publisher")
	if pubElem != nil {
		r.MetaData.Publisher = packageDoc.FindElement("//dc:publisher").Text()
	} else {
		r.MetaData.Publisher = ""
	}

	CreaterElements := packageDoc.FindElements("//dc:creator")
	for _, CreaterElement := range CreaterElements {
		r.MetaData.Creator = append(r.MetaData.Creator, CreaterElement.Text())
	}
	ContributorElements := packageDoc.FindElements("//dc:contributor")
	for _, ContributorElement := range ContributorElements {
		r.MetaData.Contributers = append(r.MetaData.Contributers, ContributorElement.Text())
	}

	metaPath, err := etree.CompilePath("//meta[@name='cover']")
	if err != nil {
		return err
	}
	itemProPath, err := etree.CompilePath("//item[@properties='cover-image']")
	if err != nil {
		return err
	}
	if cover := packageDoc.FindElementPath(metaPath); cover != nil {
		coverImageId := cover.SelectAttrValue("content", "")
		r.MetaData.Cover.Name = "cover"
		r.MetaData.Cover.Data = coverImageId
	} else if cover := packageDoc.FindElementPath(itemProPath); cover != nil {
		coverImagePath := cover.SelectAttrValue("id", "")
		r.MetaData.Cover.Name = "cover"
		r.MetaData.Cover.Data = coverImagePath
	} else {
		r.MetaData.Cover.Name = ""
		r.MetaData.Cover.Data = ""
	}

	manifestTextId := make(map[string]string)

	for _, item := range packageDoc.FindElements("//manifest/item") {
		itemId := item.SelectAttrValue("id", "")
		itemHref := item.SelectAttrValue("href", "")
		itemType := item.SelectAttrValue("media-type", "")
		_, name := filepath.Split(itemHref)
		if utils.Contains(Image, itemType) {
			if itemId == r.MetaData.Cover.Data {
				r.MetaData.Cover.Name = name
				imageRawData := utils.GetFileData(filepath.Join(r.tempDir, itemHref))
				r.MetaData.Cover.Data = string(base64.StdEncoding.EncodeToString(imageRawData))
				continue
			}
			r.Images = append(r.Images, File{
				Name: name,
				Data: utils.GetFileData(filepath.Join(r.tempDir, itemHref)),
				Type: itemType,
			})
			r.Resources = append(r.Resources, Resource{
				Id:   itemId,
				Name: name,
				Type: itemType,
				Data: "",
			})
		}
		if itemType == "application/xhtml+xml" {
			manifestTextId[itemId] = itemHref
		}
	}

	textNode := html.PMNode{
		Type:    "doc",
		Content: []*html.PMNode{},
	}
	for _, itemref := range packageDoc.FindElements("//spine/itemref") {
		ref := itemref.SelectAttrValue("idref", "")
		if href, ok := manifestTextId[ref]; ok {
			fs, _ := os.Open(filepath.Join(r.tempDir, href))
			defer fs.Close()
			rawData, _ := io.ReadAll(fs)
			ast, err := html.LoadHTML(rawData)
			if err != nil {
				return err
			}
			currentNode := html.ConvertIntoProseMirrorScheme(ast, utils.CombineMap(extension, map[string]html.TagParser{
				"img": &html.ImageParser{
					FoldName: utils.GenerateHash([]byte(r.targetPath)),
				},
				"h1": &html.HeaderParser{
					Level: 1,
				},
				"h2": &html.HeaderParser{
					Level: 2,
				},
				"h3": &html.HeaderParser{
					Level: 3,
				},
				"h4": &html.HeaderParser{
					Level: 4,
				},
				"h5": &html.HeaderParser{
					Level: 5,
				},
				"h6": &html.HeaderParser{
					Level: 6,
				},
				"p":     &html.TextParser{},
				"span":  &html.TextParser{},
				"table": &html.TableParser{},
				"image": &html.SVGParser{
					FoldName: utils.GenerateHash([]byte(r.targetPath)),
				},
				"br": &html.BrParser{},
				"ol": &html.ListParser{
					Type: "orderedList",
				},
				"ul": &html.ListParser{
					Type: "bulletList",
				},
				"code": &html.CodeBlockParser{},
			}))
			textNode.Content = append(textNode.Content, currentNode.Content...)
		}
	}
	data, err := json.Marshal(textNode)
	if err != nil {
		return err
	}
	r.Content = base64.StdEncoding.EncodeToString(data)
	err = r.moveImage()
	if err != nil {
		return err
	}
	return nil
}

func (r *Reader) Close() error {
	currentPath := filepath.Join(r.tempDir, "..")
	return os.RemoveAll(currentPath)
}

func Dezip(path string, tempDir string) error {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer archive.Close()

	_, name := filepath.Split(path)

	tempDir = filepath.Join(tempDir, fmt.Sprintf("~%s", name))
	for _, f := range archive.File {
		filePath := filepath.Join(tempDir, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(tempDir)+string(os.PathSeparator)) {
			return errors.New("zip: illegal file path")
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}
		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}
