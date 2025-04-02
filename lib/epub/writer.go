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

// The struct that store the content of each file
type XMLData struct {
	Container []byte
	Package   []byte
	Nav       []byte
	Text      []byte
}

// The epub writer struct
type Writer struct {
	// The JsonData struct that store the content of the epub
	JsonData JsonData

	// The temp fold that store the content of the epub
	tempDir string
	// the byte content of the epub
	XMLData XMLData
	// the media content of the epub
	Media []File
	// the output path of the epub
	targetPath string
}

// Convert the J_nav struct to NavLi struct
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

// Write the byte content to the file
func WriteToFile(content []byte, FilePath string) {
	fs, _ := os.OpenFile(FilePath, os.O_RDWR|os.O_CREATE, 0766)
	fs.Write(content)
	fs.Close()
}

// constructor of the Writer struct
func NewWriter(targetPath string, tempDir string, jsonData *JsonData) *Writer {
	_, name := filepath.Split(targetPath)
	r := &Writer{
		JsonData:   *jsonData,
		tempDir:    filepath.Join(tempDir, fmt.Sprintf("~$%s", name)),
		targetPath: targetPath,
	}
	return r
}

// form the container.xml file
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
		Dir:        w.JsonData.MetaData.Meta.TextDir,
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
					Id:        "katex.css",
					Href:      "Styles/katex.css",
					MediaType: "text/css",
				},
				{
					Id:        "katex.ams.regular",
					Href:      "Styles/fonts/KaTeX_AMS-Regular.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.caligraphic.bold",
					Href:      "Styles/fonts/KaTeX_Caligraphic-Bold.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.caligraphic.regular",
					Href:      "Styles/fonts/KaTeX_Caligraphic-Regular.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.fraktur.bold",
					Href:      "Styles/fonts/KaTeX_Fraktur-Bold.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.fraktur.regular",
					Href:      "Styles/fonts/KaTeX_Fraktur-Regular.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.main.bold",
					Href:      "Styles/fonts/KaTeX_Main-Bold.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.main.bolditalic",
					Href:      "Styles/fonts/KaTeX_Main-BoldItalic.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.main.italic",
					Href:      "Styles/fonts/KaTeX_Main-Italic.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.main.regular",
					Href:      "Styles/fonts/KaTeX_Main-Regular.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.math.bolditalic",
					Href:      "Styles/fonts/KaTeX_Math-BoldItalic.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.math.italic",
					Href:      "Styles/fonts/KaTeX_Math-Italic.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.sanserif.bold",
					Href:      "Styles/fonts/KaTeX_SansSerif-Bold.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.sanserif.italic",
					Href:      "Styles/fonts/KaTeX_SansSerif-Italic.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.sanserif.regular",
					Href:      "Styles/fonts/KaTeX_SansSerif-Regular.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.script.regular",
					Href:      "Styles/fonts/KaTeX_Script-Regular.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.size1.regular",
					Href:      "Styles/fonts/KaTeX_Size1-Regular.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.size2.regular",
					Href:      "Styles/fonts/KaTeX_Size2-Regular.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.size3.regular",
					Href:      "Styles/fonts/KaTeX_Size3-Regular.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.size4.regular",
					Href:      "Styles/fonts/KaTeX_Size4-Regular.ttf",
					MediaType: "application/vnd.ms-opentype",
				},
				{
					Id:        "katex.typewriter.regular",
					Href:      "Styles/fonts/KaTeX_Typewriter-Regular.ttf",
					MediaType: "application/vnd.ms-opentype",
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
				{
					Id:        "katex.js",
					Href:      "Scripts/katex.min.js",
					MediaType: "application/javascript",
				},
				{
					Id:        "auto-render.js",
					Href:      "Scripts/auto-render.js",
					MediaType: "application/javascript",
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
	// if the layout is reflowable, add the flow property
	// if the layout is pre-paginated, add the viewport property
	// but the pre-paginated and the flow property can't exist at the same time
	layout := w.JsonData.MetaData.Meta.Layout
	if layout == "reflowable" {
		content.Metadata.Metas = append(content.Metadata.Metas, MetaNode{
			Property: "rendition:flow",
			Value:    w.JsonData.MetaData.Meta.Flow,
		})
	}
	content.Metadata.Metas = append(content.Metadata.Metas, MetaNode{
		Property: "dcterms:modified",
		Value:    time.Now().Format("2006-01-02T15:04:05Z"),
	})
	content.Metadata.Metas = append(content.Metadata.Metas, MetaNode{
		Property: "rendition:layout",
		Value:    layout,
	})
	content.Metadata.Metas = append(content.Metadata.Metas, MetaNode{
		Property: "rendition:spread",
		Value:    w.JsonData.MetaData.Meta.Spread,
	})
	content.Metadata.Metas = append(content.Metadata.Metas, MetaNode{
		Property: "rendition:orientation",
		Value:    w.JsonData.MetaData.Meta.Orientation,
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

// form the nav.xhtml file
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

// form the text.xhtml file
func (w *Writer) formText() error {
	// replace the image path to the local path
	imagePathRegex := regexp.MustCompile(`http(s)?://127.0.0.1:(\d+)/(\\[0-9a-z]+\\|[0-9a-z]+/)`)
	// remove illegal characters
	illegalCharRegex := regexp.MustCompile(`&nbsp;|&ensp;|&emsp;|&thinsp;|&zwnj;|&zwj;|&lrm;|&rlm;`)
	w.JsonData.Content = illegalCharRegex.ReplaceAllString(w.JsonData.Content, " ")
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
				{
					Href: "../Styles/katex.css",
					Rel:  "stylesheet",
					Type: "text/css",
				},
			},
			Scripts: []Script{
				{
					Src:   "../Scripts/katex.min.js",
					Type:  "text/javascript",
					Defer: true,
				},
				{
					Src:    "../Scripts/auto-render.js",
					Type:   "text/javascript",
					Defer:  true,
					Onload: "renderMathInElement(document.body);",
				},
			},
		},
		Body: XhtmlBody{
			Section: imagePathRegex.ReplaceAllString(w.JsonData.Content, "../Images/"),
		},
	}
	if w.JsonData.MetaData.Meta.Layout == "pre-paginated" {
		viewport := w.JsonData.MetaData.Meta.Proportions
		if viewport == "auto" {
			viewport = "width=device-width, height=device-height"
		}
		text.Header.Meta = append(text.Header.Meta, MetaNode{
			Name:    "viewport",
			Content: viewport,
		})
	}
	b, _ := xml.MarshalIndent(text, "", "  ")
	w.XMLData.Text = b
	return nil
}

// move the images of the imagehost back to the epub
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

func copyFile(file io.ReadCloser, dst string) error {
	defer file.Close()
	of, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer of.Close()
	io.Copy(of, file)
	return nil
}

func loadResources(path, dst string) error {
	r, err := zip.OpenReader(filepath.Join(path, "res.pak"))
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		if filepath.Ext(f.Name) == ".css" {
			// Load CSS files
			rc, err := f.Open()
			if err != nil {
				return err
			}
			copyFile(rc, filepath.Join(dst, "OEBPS/Styles", f.Name))
		}
		if filepath.Ext(f.Name) == ".js" {
			// Load JS files
			rc, err := f.Open()
			if err != nil {
				return err
			}
			copyFile(rc, filepath.Join(dst, "OEBPS/Scripts", f.Name))
		}
		if filepath.Ext(f.Name) == ".ttf" {
			// Load font files
			rc, err := f.Open()
			if err != nil {
				return err
			}
			copyFile(rc, filepath.Join(dst, "OEBPS/Styles/fonts", f.Name))
		}
	}
	return nil
}

// form all the required files to the temp fold
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
	if err = os.MkdirAll(FoldName+"/OEBPS/Scripts", os.ModePerm); err != nil {
		return err
	}
	if err = os.MkdirAll(FoldName+"/OEBPS/Styles/fonts", os.ModePerm); err != nil {
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
	err = loadResources(filepath.Join(w.tempDir, "../../"), FoldName)
	if err != nil {
		return err
	}
	for _, medium := range w.Media {
		err := os.WriteFile(FoldName+"/OEBPS/Images/"+medium.Name+utils.FileSuffix[medium.Type], medium.Data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// write the content
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
