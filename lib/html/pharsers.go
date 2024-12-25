// Description: The parser corresponding to each tag
// Author: Hughie21
// Date: 2024-12-20
// license that can be found in the LICENSE file.
package html

import (
	"fmt"
	"html"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type (
	TagParser interface {
		Parse(node *AstElement) *PMNode
	}

	HeaderParser struct {
		Level int
	}
	TextParser  struct{}
	ImageParser struct {
		FoldName string
	}
	ListParser struct {
		Type string
	}
	TableParser struct {
		FoldName string // compatible to the image
	}
	BrParser        struct{}
	CodeBlockParser struct{}
)

// parse the style string to a map
func splitStyle(style string) map[string]string {
	result := make(map[string]string)
	for _, item := range strings.Split(style, ";") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		kv := strings.SplitN(item, ":", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			if (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) || (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) {
				value = value[1 : len(value)-1]
			}
			result[key] = value
		}
	}
	return result
}

// backward traversal of the ast tree
func backward(node *AstElement, textNode *PMNode) {
	s := NewStack()
	s.Push(node.Parent)
	for {
		currentNode := s.Pop().(*AstElement)
		if currentNode == nil {
			break
		}
		if currentNode.Tag == "p" {
			break
		}
		s.Push(currentNode.Parent)
		switch currentNode.Tag {
		case "span":
			handleSpanTag(currentNode, textNode)
		case "strong":
			handleStrongTag(textNode)
		case "em":
			handleEmTag(textNode)
		case "s":
			handleSTag(textNode)
		case "a":
			handleATag(currentNode, textNode)
		default:
			continue
		}
	}
}

// This method is used to parse the text in the node
// it will iterates through all the child nodes under the given node
//
// Whenever the type of a node is checked to be 3, the node is forward traversed
// to its parent node, looking for other possible format-related nodes
func FindTextTill(node *AstElement, parent *PMNode) {
	s := NewStack()
	s.Push(node)
	for {
		if s.Len() == 0 {
			break
		}
		node := s.Pop().(*AstElement)
		if node.Type == 3 {
			textNode := PMNode{
				Type:    "text",
				Mark:    []*PMNode{},
				Content: []*PMNode{},
				Text:    "",
			}
			textNode.Text = html.UnescapeString(node.Text)
			backward(node, &textNode)
			parent.Content = append([]*PMNode{&textNode}, parent.Content...)
			continue
		}
		for _, child := range node.Children {
			s.Push(child)
		}
	}
}

// dealing with the tag <a>
func handleATag(node *AstElement, parent *PMNode) {
	href, ok := node.Attrs["href"]
	if !ok {
		return
	}
	target, ok := node.Attrs["target"]
	if !ok {
		target = "_blank"
	}
	rel, ok := node.Attrs["rel"]
	if !ok {
		rel = "noopener noreferrer nofollow"
	}
	linkNode := &PMNode{
		Type: "link",
		Attrs: map[string]interface{}{
			"href":   href,
			"target": target,
			"rel":    rel,
			"class":  "",
		},
	}
	parent.Mark = append([]*PMNode{linkNode}, parent.Mark...)
}

// dealing with the tag <span>
func handleSpanTag(node *AstElement, parent *PMNode) {
	textStyleString, ok := node.Attrs["style"]
	if !ok {
		return
	}
	textStyle := splitStyle(html.UnescapeString(textStyleString))
	fontSize, ok := textStyle["font-size"]
	if !ok {
		fontSize = "1rem"
	}
	fontColor, ok := textStyle["color"]
	if !ok {
		fontColor = "#000000"
	}
	backgroundColor, ok := textStyle["background"]
	if !ok {
		backgroundColor = "transparent"
	}
	fontFamily, ok := textStyle["font-family"]
	if !ok {
		fontFamily = "sans-serif"
	}
	markNode := &PMNode{
		Type: "textStyle",
		Attrs: map[string]interface{}{
			"fontSize":        fontSize,
			"fontColor":       fontColor,
			"backgroundColor": backgroundColor,
			"fontFamily":      fontFamily,
		},
	}
	parent.Mark = append([]*PMNode{markNode}, parent.Mark...)
}

// dealing with the tag <strong>
func handleStrongTag(parent *PMNode) {
	bold := &PMNode{
		Type: "bold",
	}
	parent.Mark = append([]*PMNode{bold}, parent.Mark...)
}

// dealing with the tag <em>
func handleEmTag(parent *PMNode) {
	italic := &PMNode{
		Type: "italic",
	}
	parent.Mark = append([]*PMNode{italic}, parent.Mark...)
}

// dealing with the tag <s>
func handleSTag(parent *PMNode) {
	strike := &PMNode{
		Type: "strike",
	}
	parent.Mark = append([]*PMNode{strike}, parent.Mark...)
}

// Parse the tag <h1> to <h6>
func (p *HeaderParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	// level := map[string]int{"h1": 1, "h2": 2, "h3": 3, "h4": 4, "h5": 5, "h6": 6}[node.Tag]
	heading := &PMNode{
		Type:    "custom-heading",
		Attrs:   map[string]interface{}{"level": p.Level},
		Content: []*PMNode{},
	}
	FindTextTill(node, heading)
	return heading
}

// Parse the tag <img> or <image>
func (p *ImageParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	imageAttr := make(map[string]interface{})
	if v, ok := node.Attrs["xlink:href"]; ok { // compatible to the svg
		_, imageName := filepath.Split(v)
		imageAttr["src"] = fmt.Sprintf("http://127.0.0.1:7288/%s/%s", p.FoldName, imageName)
	}
	if v, ok := node.Attrs["src"]; ok {
		_, imageName := filepath.Split(v)
		imageAttr["src"] = fmt.Sprintf("http://127.0.0.1:7288/%s/%s", p.FoldName, imageName)
	}
	if v, ok := node.Attrs["alt"]; ok {
		imageAttr["alt"] = v
		imageAttr["title"] = v
	} else {
		imageAttr["alt"] = ""
		imageAttr["title"] = ""
	}

	widthReg := regexp.MustCompile(`width:\s?(\d+)px`)
	style, ok := node.Attrs["style"]
	if ok {
		matches := widthReg.FindStringSubmatch(style)
		if len(matches) > 1 {
			width := matches[1]
			if intV, _ := strconv.Atoi(width); intV > 500 {
				width = "500"
			}
			imageAttr["style"] = widthReg.ReplaceAllString(style, fmt.Sprintf("width: %spx", width))
		} else {
			imageAttr["style"] = "width: 300px; height: auto; cursor: pointer;"
		}
	} else {
		imageAttr["style"] = "width: 300px; height: auto; cursor: pointer;"
	}

	image := &PMNode{
		Type:  "image",
		Attrs: imageAttr,
	}
	return image
}

// Parse the tag <p>
func (p *TextParser) Parse(node *AstElement) *PMNode {
	paragraph := &PMNode{
		Type:    "paragraph",
		Content: []*PMNode{},
	}

	FindTextTill(node, paragraph)

	return paragraph
}

// Parse the tag <ul> or <ol>
func (p *ListParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	bulletList := &PMNode{
		Type:    p.Type,
		Content: []*PMNode{},
	}
	s := NewStack()
	s.Push(node)
	for {
		if s.Len() == 0 {
			break
		}
		node := s.Pop().(*AstElement)
		if node.Tag == "li" {
			listItem := &PMNode{
				Type: "listItem",
				Content: []*PMNode{
					{
						Type:    "paragraph",
						Content: []*PMNode{},
					},
				},
			}
			FindTextTill(node, listItem.Content[0])
			bulletList.Content = append(bulletList.Content, listItem)
		} else {
			for _, child := range node.Children {
				s.Push(child)
			}
		}
	}
	return bulletList
}

// Parse the tag <table>
func (p *TableParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	table := &PMNode{
		Type:    "table",
		Content: []*PMNode{},
	}

	// parse the table cells
	sf := func(node *AstElement, parent *PMNode) {
		for _, child := range node.Children {
			if child.Tag == "td" || child.Tag == "th" {
				cellType := "tableCell"
				if child.Tag == "th" {
					cellType = "tableHeader"
				}
				colwidth, ok := child.Attrs["colwidth"]
				if !ok {
					colwidth = ""
				}
				colspan, err := strconv.Atoi(child.Attrs["colspan"])
				if err != nil {
					colspan = 1
				}
				rowspan, err := strconv.Atoi(child.Attrs["rowspan"])
				if err != nil {
					rowspan = 1
				}
				tableCell := &PMNode{
					Type: cellType,
					Attrs: map[string]interface{}{
						"colspan":  colspan,
						"rowspan":  rowspan,
						"colwidth": colwidth,
					},
					Content: []*PMNode{},
				}
				for _, grandChild := range child.Children {
					if grandChild.Tag == "p" {
						paragraph := &PMNode{
							Type:    "paragraph",
							Content: []*PMNode{},
						}
						FindTextTill(grandChild, paragraph)
						tableCell.Content = append(tableCell.Content, paragraph)
					} else if grandChild.Tag == "br" {
						br := &PMNode{
							Type: "paragraph",
						}
						tableCell.Content = append(tableCell.Content, br)
					} else if grandChild.Children[0].Tag == "img" { // if the cell contains an image
						imageParser := &ImageParser{
							FoldName: p.FoldName,
						}
						image := imageParser.Parse(grandChild.Children[0])
						if image != nil {
							tableCell.Content = append(tableCell.Content, image)
						}
					}
				}
				parent.Content = append(parent.Content, tableCell)
			}
		}
	}

	// parse the table rows
	parseRows := func(node *AstElement, parent *PMNode) {
		for _, child := range node.Children {
			if child.Tag == "tr" {
				tableRow := &PMNode{
					Type:    "tableRow",
					Content: []*PMNode{},
				}
				sf(child, tableRow)
				parent.Content = append(parent.Content, tableRow)
			}
		}
	}

	for _, child := range node.Children {
		if child.Tag == "tr" {
			tableRow := &PMNode{
				Type:    "tableRow",
				Content: []*PMNode{},
			}
			sf(child, tableRow)
			table.Content = append(table.Content, tableRow)
		} else if child.Tag == "thead" || child.Tag == "tbody" {
			parseRows(child, table)
		}
	}
	return table
}

// Parse the tag <br>
func (p *BrParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	br := &PMNode{
		Type: "paragraph",
	}
	return br
}

// Parse the tag <pre>
func (p *CodeBlockParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	codeBlock := &PMNode{
		Type: "codeBlock",
		Attrs: map[string]interface{}{
			"class": node.Attrs["class"],
		},
		Content: []*PMNode{},
	}
	FindTextTill(node, codeBlock)
	return codeBlock
}
