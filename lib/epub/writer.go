// Description: The epub writing program
// Author: Hughie21
// Date: 2024-11-29
// license that can be found in the LICENSE file.

package epub

import (
	"archive/zip"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/hughie21/NovelMaker/lib/utils"
)

type XMLData struct {
	Container []byte
	Package   []byte
	Nav       []byte
	Text      []byte
}

type Writer struct {
	JsonData   JsonData
	tempDir    string
	XMLData    XMLData
	Media      []File
	targetPath string
}

func convertNav(point J_nav, depth int) NavLi {
	navLi := NavLi{
		A: NavA{
			Text:  point.Label,
			Class: fmt.Sprintf("sgc-toc-level-%d", depth),
			Href:  "Text/text.xhtml#" + point.Href,
		},
	}
	if len(point.Child) == 0 {
		return navLi
	}
	depth++
	ol := NavOl{
		Li: []NavLi{},
	}
	for _, child := range point.Child {
		ol.Li = append(ol.Li, convertNav(child, depth))
	}
	navLi.Ol = &ol
	return navLi
}

func WriteToFile(content []byte, FilePath string) {
	fs, _ := os.OpenFile(FilePath, os.O_RDWR|os.O_CREATE, 0766)
	fs.Write(content)
	fs.Close()
}

func NewWriter(targetPath string, tempDir string, jsonData *JsonData) *Writer {
	_, name := filepath.Split(targetPath)
	r := &Writer{
		JsonData:   *jsonData,
		tempDir:    filepath.Join(tempDir, fmt.Sprintf("~$%s", name)),
		targetPath: targetPath,
	}
	return r
}

func (w *Writer) formContainer() error {
	container := Container{
		Xmls:    "urn:oasis:names:tc:opendocument:xmlns:container",
		Version: "1.0",
		RootFiles: RootFiles{
			RootFile: []RootFile{
				{
					FullPath:  "OEBPS/content.opf",
					MediaType: "application/oebps-package+xml",
				},
			},
		},
	}
	b, _ := xml.MarshalIndent(container, "", "  ")
	w.XMLData.Container = b
	return nil
}

func (w *Writer) formPackage() error {
	root := "<?xml version='1.0' encoding='utf-8'?>\n"
	content := PackageNode{
		Xmlns:      "http://www.idpf.org/2007/opf",
		DC:         "http://purl.org/dc/elements/1.1/",
		DCTerm:     "http://purl.org/dc/terms/",
		Identifier: "BookId",
		Version:    "3.0",
		Metadata: MetadataNode{
			Title:    w.JsonData.MetaData.Title,
			Creators: []DCCreator{},
			Identifier: Identifier{
				Id:    "BookId",
				Value: w.JsonData.MetaData.Identifier,
			},
			Language:     w.JsonData.MetaData.Language,
			Contributors: []DCContributor{},
			Description:  w.JsonData.MetaData.Description,
			Publisher:    w.JsonData.MetaData.Publisher,
			Subject:      []DCSubject{},
			Metas:        []MetaNode{},
		},
		Manifest: ManifestNode{
			Items: []ItemNode{
				{
					Id:        "style.css",
					Href:      "Styles/style.css",
					MediaType: "text/css",
				},
				{
					Id:         "nav.xhtml",
					Href:       "nav.xhtml",
					MediaType:  "application/xhtml+xml",
					Properties: "nav",
				},
				{
					Id:        "text.xhtml",
					Href:      "Text/text.xhtml",
					MediaType: "application/xhtml+xml",
				},
			},
		},
		Spine: SpineNode{
			Items: []SpineItemNode{
				{
					Idref: "text.xhtml",
				},
			},
		},
	}
	content.Metadata.Metas = append(content.Metadata.Metas, MetaNode{
		Property: "dcterms:modified",
		Value:    time.Now().Format("2006-01-02T15:04:05Z"),
	})
	content.Metadata.Metas = append(content.Metadata.Metas, MetaNode{
		Property: "rendition:flow",
		Value:    "auto",
	})
	content.Metadata.Metas = append(content.Metadata.Metas, MetaNode{
		Property: "rendition:orientation",
		Value:    "auto",
	})
	for _, creator := range w.JsonData.MetaData.Creator {
		content.Metadata.Creators = append(content.Metadata.Creators, DCCreator{
			Id:    "id-" + creator + utils.RandomString(5),
			Value: creator,
		})
	}
	for _, contributor := range w.JsonData.MetaData.Contributers {
		content.Metadata.Contributors = append(content.Metadata.Contributors, DCContributor{
			Id:    "id-" + contributor + utils.RandomString(5),
			Value: contributor,
		})
	}
	for _, subject := range w.JsonData.MetaData.Subject {
		content.Metadata.Subject = append(content.Metadata.Subject, DCSubject{
			Id:    "id-" + subject + utils.RandomString(5),
			Value: subject,
		})
	}
	for _, item := range w.JsonData.Resources {
		if item.Id == "cover" {
			content.Manifest.Items = append(content.Manifest.Items, ItemNode{
				Id:         "image_" + item.Id,
				Href:       "Images/" + item.Name + utils.FileSuffix[item.Type],
				MediaType:  item.Type,
				Properties: "cover-image",
			})
			continue
		}
		content.Manifest.Items = append(content.Manifest.Items, ItemNode{
			Id:        "item_" + item.Id,
			Href:      "Images/" + item.Name + utils.FileSuffix[item.Type],
			MediaType: item.Type,
		})
	}
	b, _ := xml.MarshalIndent(content, "", "  ")
	b = append([]byte(root), b...)
	w.XMLData.Package = b
	return nil
}

func (w *Writer) formNav() error {
	root := `<?xml version="1.0" encoding="utf-8"?>
	<!DOCTYPE html>
	`
	nav := w.JsonData.Nav
	NAV := NavHTML{
		Xmlns: "http://www.w3.org/1999/xhtml",
		XEpub: "http://www.idpf.org/2007/ops",
		Head: HtmlHead{
			Title: w.JsonData.MetaData.Title,
			Link: []Link{
				{
					Rel:  "stylesheet",
					Href: "Styles/style.css",
					Type: "text/css",
				},
			},
		},
		Body: NavBody{
			Nav: Nav{
				Type: "toc",
				Ol: NavOl{
					Li: []NavLi{},
				},
			},
		},
	}

	depth := 1

	for _, point := range nav {
		NAV.Body.Nav.Ol.Li = append(NAV.Body.Nav.Ol.Li, convertNav(point, depth))
	}
	b, _ := xml.MarshalIndent(NAV, "", "  ")
	b = append([]byte(root), b...)
	w.XMLData.Nav = b
	return nil
}

func (w *Writer) formText() error {
	imagePathRegex := regexp.MustCompile(`http(s)?://127.0.0.1:(\d+)/[0-9a-z]+/`)
	text := XhtmlHTML{
		Xmlns: "http://www.w3.org/1999/xhtml",
		Lang:  "en",
		Header: XhtmlHead{
			Title: w.JsonData.MetaData.Title,
			Link: []Link{
				{
					Href: "../Styles/style.css",
					Rel:  "stylesheet",
					Type: "text/css",
				},
			},
		},
		Body: XhtmlBody{
			Section: imagePathRegex.ReplaceAllString(w.JsonData.Content, "../Images/"),
		},
	}
	b, _ := xml.MarshalIndent(text, "", "  ")
	w.XMLData.Text = b
	return nil
}

func (w *Writer) loadMedia() error {
	for _, item := range w.JsonData.Resources {
		if utils.Contains(Image, item.Type) {
			decodeData, err := base64.StdEncoding.DecodeString(item.Data)
			if err != nil {
				return err
			}
			w.Media = append(w.Media, File{Name: item.Name, Data: decodeData, Type: item.Type})
		}
	}
	return nil
}

func (w *Writer) toTemp() error {
	FoldName := w.tempDir
	var err error
	if err = os.Mkdir(FoldName, os.ModePerm); err != nil {
		return err
	}
	if err = os.MkdirAll(FoldName+"/META-INF", os.ModePerm); err != nil {
		return err
	}
	if err = os.MkdirAll(FoldName+"/OEBPS", os.ModePerm); err != nil {
		return err
	}
	if err = os.MkdirAll(FoldName+"/OEBPS/Text", os.ModePerm); err != nil {
		return err
	}
	if err = os.MkdirAll(FoldName+"/OEBPS/Styles", os.ModePerm); err != nil {
		return err
	}
	if err = os.MkdirAll(FoldName+"/OEBPS/Images", os.ModePerm); err != nil {
		return err
	}
	mimetype := `application/epub+zip`
	WriteToFile([]byte(mimetype), FoldName+"/mimetype")
	WriteToFile(w.XMLData.Container, filepath.Join(FoldName, "META-INF/container.xml"))
	WriteToFile(w.XMLData.Package, filepath.Join(FoldName, "OEBPS/content.opf"))
	WriteToFile(w.XMLData.Nav, filepath.Join(FoldName, "OEBPS/nav.xhtml"))
	WriteToFile(w.XMLData.Text, filepath.Join(FoldName, "OEBPS/Text/text.xhtml"))
	fs, err := os.Open(filepath.Join(filepath.Join(w.tempDir, "../../"), "style", "style.css"))
	if err != nil {
		return err
	}
	of, _ := os.Create(FoldName + "/OEBPS/Styles/style.css")
	io.Copy(of, fs)
	defer fs.Close()
	defer of.Close()
	for _, medium := range w.Media {
		err := os.WriteFile(FoldName+"/OEBPS/Images/"+medium.Name+utils.FileSuffix[medium.Type], medium.Data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) Write() error {
	err := w.formPackage()
	if err != nil {
		return err
	}
	err = w.formContainer()
	if err != nil {
		return err
	}
	err = w.formNav()
	if err != nil {
		return err
	}
	err = w.formText()
	if err != nil {
		return err
	}
	err = w.loadMedia()
	if err != nil {
		return err
	}
	err = w.toTemp()
	if err != nil {
		return err
	}
	err = WriteEpub(w.tempDir, w.targetPath)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) Close() {
	os.RemoveAll(w.tempDir)
}

// Generate an EPUB file from the temp fold
// @Quoted from https://github.com/gonejack/html-to-epub/blob/main/go-epub/write.go
func WriteEpub(tempDir string, destFilePath string) error {
	f, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	z := zip.NewWriter(f)
	defer z.Close()

	skipMimetypeFile := false

	var addFileToZip = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the path of the file relative to the folder we're zipping
		relativePath, err := filepath.Rel(tempDir, path)
		relativePath = filepath.ToSlash(relativePath)
		if err != nil {
			// tempDir and path are both internal, so we shouldn't get here
			return err
		}

		// Only include regular files, not directories
		if !info.Mode().IsRegular() {
			return nil
		}

		var w io.Writer
		if path == filepath.Join(tempDir, "mimetype") {
			// Skip the mimetype file if it's already been written
			if skipMimetypeFile == true {
				return nil
			}
			// The mimetype file must be uncompressed according to the EPUB spec
			w, err = z.CreateHeader(&zip.FileHeader{
				Name:   relativePath,
				Method: zip.Store,
			})
		} else {
			w, err = z.Create(relativePath)
		}
		if err != nil {
			return err
		}

		r, err := os.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()

		_, err = io.Copy(w, r)
		if err != nil {
			return err
		}

		return nil
	}

	// Add the mimetype file first
	mimetypeFilePath := filepath.Join(tempDir, "mimetype")
	mimetypeInfo, err := os.Lstat(mimetypeFilePath)
	if err != nil {
		return err
	}
	err = addFileToZip(mimetypeFilePath, mimetypeInfo, nil)
	if err != nil {
		return err
	}

	skipMimetypeFile = true

	err = filepath.Walk(tempDir, addFileToZip)
	if err != nil {
		return err
	}

	return nil
}
