/*
@Author: Hughie
@CreateTime: 2024-8-2
@LastEditors: Hughie
@LastEditTime: 2024-08-7
@Description: Transform the Json format to the xml format and write to the epub file
*/

package epub

import (
	"archive/zip"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

var logFileName string

func LogOutPut(msg string, funcName string) {
	if logFileName == "" {
		todaystr := time.Now().Format("2006-01-02")
		todayint, _ := time.ParseInLocation("2006-01-02", todaystr, time.Local)
		logFileName = strconv.FormatInt(todayint.Unix(), 10) + ".log"
	}
	fileLogger := logger.NewFileLogger(filepath.Join("./log", logFileName))
	fileLogger.Error(fmt.Sprintf("%s -> %s", funcName, msg))
}

func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

/*
*
* File structure
*
 */

type ImagesFile struct {
	Name string
	Data []byte
}

type EpubFile struct {
	OPT    []byte
	NCX    []byte
	NAV    []byte
	Texts  []byte
	Images []ImagesFile
}

var Files EpubFile

func LoadImage(jsonData *JsonData) {
	for _, item := range jsonData.Resources {
		if item.Type == "image/jpeg" || item.Type == "image/png" {
			decodeData, e := base64.StdEncoding.DecodeString(item.Data)
			if e != nil {
				fmt.Println("Error decoding image:", e)
				continue
			}
			Files.Images = append(Files.Images, ImagesFile{Name: item.Name, Data: decodeData})
		}
	}
}

func convertNavMap(jNav J_nav, count int) NavPoint {
	navPoint := NavPoint{
		Navlable: TextNode{
			Text: jNav.Label,
		},
		Content: Content{
			Src: "Text/text.xhtml#" + jNav.Href,
		},
		Id:        jNav.Id,
		PlayOrder: count,
	}
	for _, child := range jNav.Child {
		count++
		navPoint.NavPoints = append(navPoint.NavPoints, convertNavMap(child, count))
	}
	return navPoint
}

func convertNavDiv(jNav J_nav, depth int) NavDiv {
	navDiv := NavDiv{
		Class: "sgc-toc-level-" + strconv.Itoa(depth),
		Navlable: NaVA{
			Text: jNav.Label,
			Href: "Text/text.xhtml#" + jNav.Href,
		},
		Child: []NavDiv{},
	}
	depth++
	for _, child := range jNav.Child {
		navDiv.Child = append(navDiv.Child, convertNavDiv(child, depth))
	}
	return navDiv
}

func FormNcx(jsonData *JsonData) []byte {
	root := `<?xml version="1.0" encoding="UTF-8" standalone="no" ?>
<!DOCTYPE ncx PUBLIC "-//NISO//DTD ncx 2005-1//EN"
"http://www.daisy.org/z3986/2005/ncx-2005-1.dtd">
`
	nav := jsonData.Nav
	NCX := NcxNode{
		Xmlns:   "http://www.daisy.org/z3986/2005/ncx/",
		Version: "2005-1",
		Header: HeadNode{
			Meta: []MetaNode{
				{
					Name:    "dtb:uid",
					Content: jsonData.MetaData.Identifier,
				},
				{
					Name:    "dtb:depth",
					Content: "0",
				},
				{
					Name:    "dtb:totalPageCount",
					Content: "0",
				},
				{
					Name:    "dtb:maxPageNumber",
					Content: "0",
				},
			},
		},
		Title: TextNode{
			Text: jsonData.MetaData.Title,
		},
		Author: []TextNode{},
		NavMap: NavMap{
			NavPoints: []NavPoint{},
		},
	}
	for _, author := range jsonData.MetaData.Creator {
		NCX.Author = append(NCX.Author, TextNode{
			Text: author,
		})
	}
	count := 1
	for _, point := range nav {
		NCX.NavMap.NavPoints = append(NCX.NavMap.NavPoints, convertNavMap(point, count))
		count++
	}
	b, _ := xml.MarshalIndent(NCX, "", "  ")
	b = append([]byte(root), b...)
	return b
}

func FormNav(jsonData *JsonData) []byte {
	root := `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
`
	nav := jsonData.Nav
	NAV := NavHTML{
		Xmlns: "http://www.w3.org/1999/xhtml",
		Header: HtmlHead{
			Title: jsonData.MetaData.Title,
			Link: []Link{
				{
					Href: "Styles/style.css",
					Rel:  "stylesheet",
					Type: "text/css",
				},
			},
		},
		Body: NavBody{
			Nav: []NavDiv{},
		},
	}
	depth := 1
	for _, point := range nav {
		NAV.Body.Nav = append(NAV.Body.Nav, convertNavDiv(point, depth))
		depth++
	}
	b, _ := xml.MarshalIndent(NAV, "", "  ")
	b = append([]byte(root), b...)
	return b
}

func FormText(jsonData *JsonData) []byte {
	text := XhtmlHTML{
		Xmlns: "http://www.w3.org/1999/xhtml",
		Lang:  "en",
		Header: XhtmlHead{
			Title: jsonData.MetaData.Title,
			Meta: []MetaNode{
				{
					Content: "text/html; charset=utf-8",
				},
			},
			Link: []Link{
				{
					Href: "../Styles/style.css",
					Rel:  "stylesheet",
					Type: "text/css",
				},
			},
		},
		Body: XhtmlBody{
			Section: jsonData.Content,
		},
	}
	b, _ := xml.MarshalIndent(text, "", "  ")
	return b
}

func FormPack(pack *JsonData) []byte {
	root := "<?xml version='1.0' encoding='utf-8'?>\n"
	content := PackageNode{
		Xmlns:      "http://www.idpf.org/2007/opf",
		Identifier: "BookId",
		Version:    "2.0",
		Metadata: MetadataNode{
			Xmlns:    "http://purl.org/dc/elements/1.1/",
			Opt:      "http://www.idpf.org/2007/opf",
			Title:    pack.MetaData.Title,
			Creators: []DCCreator{},
			Identifier: Identifier{
				Id:    "BookId",
				Value: pack.MetaData.Identifier,
			},
			Language:     pack.MetaData.Language,
			Contributors: []DCContributor{},
			Description:  pack.MetaData.Description,
			Publisher:    pack.MetaData.Publisher,
			Subject:      []DCSubject{},
			Date:         pack.MetaData.Date,
			Metas:        []MetaNode{},
		},
		Manifest: ManifestNode{
			Items: []ItemNode{
				{
					Id:        "ncx",
					Href:      "toc.ncx",
					MediaType: "application/x-dtbncx+xml",
				},
				{
					Id:        "style.css",
					Href:      "Styles/style.css",
					MediaType: "text/css",
				},
				{
					Id:        "nav.xhtml",
					Href:      "nav.xhtml",
					MediaType: "application/xhtml+xml",
				},
				{
					Id:        "text.xhtml",
					Href:      "Text/text.xhtml",
					MediaType: "application/xhtml+xml",
				},
			},
		},
		Spine: SpineNode{
			Toc: "ncx",
			Items: []SpineItemNode{
				{
					Idref: "text.xhtml",
				},
			},
		},
	}
	for _, creator := range pack.MetaData.Creator {
		content.Metadata.Creators = append(content.Metadata.Creators, DCCreator{
			Id:    "id-" + creator,
			Value: creator,
		})
	}
	for _, contributor := range pack.MetaData.Contributers {
		content.Metadata.Contributors = append(content.Metadata.Contributors, DCContributor{
			Id:    "id-" + contributor,
			Value: contributor,
		})
	}
	for _, meta := range pack.MetaData.Meta {
		content.Metadata.Metas = append(content.Metadata.Metas, MetaNode{
			Name:    meta.Name,
			Content: meta.Content,
		})
	}
	for _, subject := range pack.MetaData.Subject {
		content.Metadata.Subject = append(content.Metadata.Subject, DCSubject{
			Id:    "id-" + subject,
			Value: subject,
		})
	}
	for _, item := range pack.Resources {
		content.Manifest.Items = append(content.Manifest.Items, ItemNode{
			Id:        item.Id,
			Href:      "Images/" + item.Name,
			MediaType: item.Type,
		})
	}
	b, _ := xml.MarshalIndent(content, "", "  ")
	b = append([]byte(root), b...)
	return b
}

func FormXML(jsonData JsonData) string {
	Files.OPT = FormPack(&jsonData)
	Files.NAV = FormNav(&jsonData)
	Files.NCX = FormNcx(&jsonData)
	Files.Texts = FormText(&jsonData)
	LoadImage(&jsonData)
	path := filepath.Join("tmp", jsonData.MetaData.Title)
	toFile(path, &Files)
	return path
}

func WirteToFile(content []byte, FilePath string) {
	fs, _ := os.OpenFile(FilePath, os.O_RDWR|os.O_CREATE, 0766)
	fs.Write(content)
	fs.Close()
}

func toFile(FoldName string, Files *EpubFile) {
	os.Mkdir(FoldName, os.ModePerm)
	os.MkdirAll(FoldName+"/META-INF", os.ModePerm)
	os.MkdirAll(FoldName+"/OEBPS", os.ModePerm)
	os.MkdirAll(FoldName+"/OEBPS/Text", os.ModePerm)
	os.MkdirAll(FoldName+"/OEBPS/Styles", os.ModePerm)
	os.MkdirAll(FoldName+"/OEBPS/Images", os.ModePerm)
	container := `<?xml version="1.0" encoding="UTF-8"?>
	<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
		<rootfiles>
			<rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
	</rootfiles>
	</container>
	`
	mimetype := `application/epub+zip`
	WirteToFile([]byte(mimetype), FoldName+"/mimetype")
	WirteToFile([]byte(container), FoldName+"/META-INF/container.xml")
	WirteToFile(Files.OPT, FoldName+"/OEBPS/content.opf")
	WirteToFile(Files.NCX, FoldName+"/OEBPS/toc.ncx")
	WirteToFile(Files.NAV, FoldName+"/OEBPS/nav.xhtml")
	fs, err := os.Open("./style.css")
	if err != nil {
		LogOutPut(err.Error(), runFuncName())
	}
	of, _ := os.Create(FoldName + "/OEBPS/Styles/style.css")
	io.Copy(of, fs)
	// for _, files := range Files.Texts {
	// 	WirteToFile(files.Data, FoldName+"/OEBPS/Text/"+files.Name)
	// }
	WirteToFile(Files.Texts, FoldName+"/OEBPS/Text/text.xhtml")
	for _, images := range Files.Images {
		err := os.WriteFile(FoldName+"/OEBPS/Images/"+images.Name, images.Data, 0644)
		if err != nil {
			LogOutPut(err.Error(), runFuncName())
		}
	}
}

// Generate an EPUB file from the temp fold
// @Quoted from https://github.com/gonejack/html-to-epub/blob/main/go-epub/write.go
func WriteEpub(tempDir string, destFilePath string) error {
	f, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			LogOutPut(err.Error(), runFuncName())
		}
	}()

	z := zip.NewWriter(f)
	defer func() {
		if err := z.Close(); err != nil {
			LogOutPut(err.Error(), runFuncName())
		}
	}()

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
			LogOutPut(err.Error(), runFuncName())
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
			LogOutPut(err.Error(), runFuncName())
		}

		r, err := os.Open(path)
		if err != nil {
			LogOutPut(err.Error(), runFuncName())
		}
		defer func() {
			if err := r.Close(); err != nil {
				LogOutPut(err.Error(), runFuncName())
			}
		}()

		_, err = io.Copy(w, r)
		if err != nil {
			LogOutPut(err.Error(), runFuncName())
		}

		return nil
	}

	// Add the mimetype file first
	mimetypeFilePath := filepath.Join(tempDir, "mimetype")
	mimetypeInfo, err := os.Lstat(mimetypeFilePath)
	if err != nil {
		LogOutPut(err.Error(), runFuncName())
	}
	err = addFileToZip(mimetypeFilePath, mimetypeInfo, nil)
	if err != nil {
		LogOutPut(err.Error(), runFuncName())
	}

	skipMimetypeFile = true

	err = filepath.Walk(tempDir, addFileToZip)
	if err != nil {
		LogOutPut(err.Error(), runFuncName())
	}

	return nil
}
