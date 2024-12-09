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
	TableParser struct{}
	SVGParser   struct {
		FoldName string
	}
	BrParser        struct{}
	CodeBlockParser struct{}
)

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

// Find the text node till the end of the tree
func FindTextTill(node *AstElement, parent *PMNode) {
	s := NewStack()
	s.Push(node)
	for {
		if s.Len() == 0 {
			break
		}
		node := s.Pop().(*AstElement)
		if node.Type == 3 {
			if node.Tag == "span" {
				textStyleString, ok := node.Attrs["style"]
				if !ok {
					textNode := &PMNode{
						Type: "text",
						Text: strings.TrimSpace(html.UnescapeString(node.Text)),
					}
					parent.Content = append([]*PMNode{textNode}, parent.Content...)
					continue
				}
				textStyle := splitStyle(html.UnescapeString(textStyleString))
				var checks []bool
				fontSize, ok := textStyle["font-size"]
				checks = append(checks, ok)
				fontColor, ok := textStyle["color"]
				checks = append(checks, ok)
				backgroundColor, ok := textStyle["background"]
				checks = append(checks, ok)
				fontFamily, ok := textStyle["font-family"]
				checks = append(checks, ok)
				if checks[0] && checks[1] && checks[2] && checks[3] {
					textNode := &PMNode{
						Type: "text",
						Text: strings.TrimSpace(html.UnescapeString(node.Text)),
						Mark: []*PMNode{
							{
								Type: "textStyle",
								Attrs: map[string]interface{}{
									"fontSize":        fontSize,
									"fontColor":       fontColor,
									"backgroundColor": backgroundColor,
									"fontFamily":      fontFamily,
								},
							},
						},
					}
					parent.Content = append([]*PMNode{textNode}, parent.Content...)
					continue
				}
			} else if node.Tag == "strong" {
				bold := &PMNode{
					Type: "text",
					Mark: []*PMNode{
						{
							Type: "bold",
						},
					},
					Text: strings.TrimSpace(html.UnescapeString(node.Text)),
				}
				parent.Content = append([]*PMNode{bold}, parent.Content...)
				continue
			} else if node.Tag == "em" {
				italic := &PMNode{
					Type: "text",
					Mark: []*PMNode{
						{
							Type: "italic",
						},
					},
					Text: strings.TrimSpace(html.UnescapeString(node.Text)),
				}
				parent.Content = append([]*PMNode{italic}, parent.Content...)
				continue
			} else if node.Tag == "s" {
				strike := &PMNode{
					Type: "text",
					Mark: []*PMNode{
						{
							Type: "strike",
						},
					},
					Text: strings.TrimSpace(html.UnescapeString(node.Text)),
				}
				parent.Content = append([]*PMNode{strike}, parent.Content...)
				continue
			}

			textNode := &PMNode{
				Type: "text",
				Text: strings.TrimSpace(html.UnescapeString(node.Text)),
			}
			parent.Content = append([]*PMNode{textNode}, parent.Content...)
			continue
		}
		for _, child := range node.Children {
			s.Push(child)
		}
	}
}

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

func (p *ImageParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	imageAttr := make(map[string]interface{})
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
			imageAttr["style"] = fmt.Sprintf("width: %spx; height: auto; cursor: pointer;", width)
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

func (p *TextParser) Parse(node *AstElement) *PMNode {
	paragraph := &PMNode{
		Type:    "paragraph",
		Content: []*PMNode{},
	}

	FindTextTill(node, paragraph)

	return paragraph
}

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
				tableCell := &PMNode{
					Type: cellType,
					Attrs: map[string]interface{}{
						"colspan":  1,
						"rowspan":  1,
						"colwidth": "",
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

func (p *SVGParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	imageAttr := make(map[string]interface{})
	if v, ok := node.Attrs["xlink:href"]; ok { // compatible to the svg
		_, imageName := filepath.Split(v)
		imageAttr["src"] = fmt.Sprintf("http://127.0.0.1:7288/%s/%s", p.FoldName, imageName)
	}
	if v, ok := node.Attrs["href"]; ok {
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
			imageAttr["style"] = fmt.Sprintf("width: %spx; height: auto; cursor: pointer;", width)
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

func (p *BrParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	br := &PMNode{
		Type: "paragraph",
	}
	return br
}

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
