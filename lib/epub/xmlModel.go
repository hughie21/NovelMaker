// Description: Structure mapping for each xml file of epub file
// Author: Hughie21
// Date: 2024-11-29
// license that can be found in the LICENSE file.

package epub

import "encoding/xml"

/*
*
* container.xml
*
 */
type Container struct {
	XMLName   xml.Name  `xml:"container"`
	Xmls      string    `xml:"xmlns,attr"`
	Version   string    `xml:"version,attr"`
	RootFiles RootFiles `xml:"rootfiles"`
}

type RootFiles struct {
	XMLName  xml.Name   `xml:"rootfiles"`
	RootFile []RootFile `xml:"rootfile"`
}

type RootFile struct {
	XMLName   xml.Name `xml:"rootfile"`
	FullPath  string   `xml:"full-path,attr"`
	MediaType string   `xml:"media-type,attr"`
}

/*
*
* content.opt
*
 */

type MetaNode struct {
	XMLName  xml.Name `xml:"meta"`
	Name     string   `xml:"name,attr,omitempty"`
	Content  string   `xml:"content,attr,omitempty"`
	Equiv    string   `xml:"http-equiv,attr,omitempty"`
	Property string   `xml:"property,attr,omitempty"`
	Value    string   `xml:",chardata"`
}

type DCCreator struct {
	XMLName xml.Name `xml:"dc:creator"`
	Id      string   `xml:"id,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type DCContributor struct {
	XMLName xml.Name `xml:"dc:contributor"`
	Id      string   `xml:"id,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type DCSubject struct {
	XMLName xml.Name `xml:"dc:subject"`
	Id      string   `xml:"id,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type Identifier struct {
	XMLName xml.Name `xml:"dc:identifier"`
	Id      string   `xml:"id,attr"`
	Value   string   `xml:",chardata"`
}

type MetadataNode struct {
	XMLName      xml.Name        `xml:"metadata"`
	Title        string          `xml:"dc:title"`
	Creators     []DCCreator     `xml:"dc:creator"`
	Identifier   Identifier      `xml:"dc:identifier"`
	Language     string          `xml:"dc:language"`
	Contributors []DCContributor `xml:"dc:contributor,omitempty"`
	Description  string          `xml:"dc:description,omitempty"`
	Publisher    string          `xml:"dc:publisher"`
	Subject      []DCSubject     `xml:"dc:subject,omitempty"`
	Date         string          `xml:"dc:date,omitempty"`
	Metas        []MetaNode      `xml:"meta"`
}

type ItemNode struct {
	XMLName    xml.Name `xml:"item"`
	Id         string   `xml:"id,attr"`
	Href       string   `xml:"href,attr"`
	MediaType  string   `xml:"media-type,attr"`
	Properties string   `xml:"properties,attr,omitempty"`
}

type ManifestNode struct {
	XMLName xml.Name   `xml:"manifest"`
	Items   []ItemNode `xml:"item"`
}

type SpineItemNode struct {
	XMLName xml.Name `xml:"itemref"`
	Idref   string   `xml:"idref,attr"`
}

type SpineNode struct {
	XMLName xml.Name        `xml:"spine"`
	Items   []SpineItemNode `xml:"itemref"`
}

type PackageNode struct {
	XMLName    xml.Name     `xml:"package"`
	Xmlns      string       `xml:"xmlns,attr"`
	Dir        string       `xml:"dir,attr,omitempty"`
	DC         string       `xml:"xmlns:dc,attr,omitempty"`      // dc: http://purl.org/dc/elements/1.1/
	DCTerm     string       `xml:"xmlns:dcterms,attr,omitempty"` // dcterms: http://purl.org/dc/terms/
	Identifier string       `xml:"unique-identifier,attr"`
	Version    string       `xml:"version,attr"`
	Metadata   MetadataNode `xml:"metadata"`
	Manifest   ManifestNode `xml:"manifest"`
	Spine      SpineNode    `xml:"spine"`
}

/*
*
* toc.ncx
*
 */

// This is abondoned in Epub v3

// type HeadNode struct {
// 	XMLName xml.Name   `xml:"head"`
// 	Meta    []MetaNode `xml:"meta"`
// }

// type TextNode struct {
// 	Text string `xml:"text"`
// }

// type Content struct {
// 	Src string `xml:"src,attr"`
// }

// type NavPoint struct {
// 	Navlable  TextNode   `xml:"navLabel"`
// 	Content   Content    `xml:"content"`
// 	Id        string     `xml:"id,attr"`
// 	PlayOrder int        `xml:"playOrder,attr"`
// 	NavPoints []NavPoint `xml:"navPoint"`
// }

// type NavMap struct {
// 	NavPoints []NavPoint `xml:"navPoint"`
// }

// type NcxNode struct {
// 	XMLName xml.Name   `xml:"ncx"`
// 	Xmlns   string     `xml:"xmlns,attr"`
// 	Version string     `xml:"version,attr"`
// 	Header  HeadNode   `xml:"head"`
// 	Title   TextNode   `xml:"docTitle"`
// 	Author  []TextNode `xml:"docAuthor"`
// 	NavMap  NavMap     `xml:"navMap"`
// }

/*
*
* nav.xhtml
*
 */
type Link struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
}

type Script struct {
	XMLName xml.Name `xml:"script"`
	Src     string   `xml:"src,attr"`
	Type    string   `xml:"type,attr"`
	Onload  string   `xml:"onload,attr,omitempty"`
	Defer   bool     `xml:"defer,attr,omitempty"`
}

type HtmlHead struct {
	Title string `xml:"title"`
	Link  []Link `xml:"link"`
}

type NavA struct {
	Text  string `xml:",chardata"`
	Href  string `xml:"href,attr"`
	Class string `xml:"class,attr,omitempty"`
}

type NavLi struct {
	XMLName xml.Name `xml:"li"`
	A       NavA     `xml:"a"`
	Ol      *NavOl   `xml:"ol,omitempty"`
}

type NavOl struct {
	XMLName xml.Name `xml:"ol"`
	Li      []NavLi  `xml:"li,omitempty"`
}

type Nav struct {
	XMLName xml.Name `xml:"nav"`
	Type    string   `xml:"epub:type,attr"`
	Ol      NavOl    `xml:"ol"`
}

type NavBody struct {
	XMLName xml.Name `xml:"body"`
	Nav     Nav      `xml:"nav"`
}

type NavHTML struct {
	XMLName xml.Name `xml:"html"`
	Xmlns   string   `xml:"xmlns,attr"`      // http://www.w3.org/1999/xhtml
	XEpub   string   `xml:"xmlns:epub,attr"` // http://www.idpf.org/2007/ops
	Head    HtmlHead `xml:"head"`
	Body    NavBody  `xml:"body"`
}

/*
*
* text.xhtml
*
 */
type XhtmlBody struct {
	XMLName xml.Name `xml:"body"`
	Section string   `xml:",innerxml"`
}

type XhtmlHead struct {
	XMLName xml.Name   `xml:"head"`
	Title   string     `xml:"title"`
	Meta    []MetaNode `xml:"meta"`
	Link    []Link     `xml:"link"`
	Scripts []Script   `xml:"script"`
}

type XhtmlHTML struct {
	XMLName xml.Name  `xml:"html"`
	Xmlns   string    `xml:"xmlns,attr"`
	Lang    string    `xml:"xml:lang,attr"`
	Header  XhtmlHead `xml:"head"`
	Body    XhtmlBody `xml:"body"`
}
