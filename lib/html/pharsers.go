package html

import (
	"fmt"
	"path/filepath"
)

type (
	TagParser interface {
		Parse(node *AstElement) *PMNode
	}

	HeaderParser struct{}
	TextParser   struct{}
	ImageParser  struct{}
	ListParser   struct{}
	TableParser  struct{}
	SVGParser    struct{}
)

func (p *HeaderParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	level := map[string]int{"h1": 1, "h2": 2, "h3": 3, "h4": 4, "h5": 5, "h6": 6}[node.Tag]
	heading := &PMNode{
		Type:    "custom-heading",
		Attrs:   map[string]string{"level": fmt.Sprintf("%d", level)},
		Content: []*PMNode{},
	}
	findTextTill(node, heading)
	return heading
}

func (p *ImageParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	imageAttr := make(map[string]string)
	if v, ok := node.Attrs["src"]; ok {
		_, imageName := filepath.Split(v)
		imageAttr["src"] = fmt.Sprintf("http://127.0.0.1:7288/%s", imageName)
	}
	if v, ok := node.Attrs["alt"]; ok {
		imageAttr["alt"] = v
		imageAttr["title"] = v
	} else {
		imageAttr["alt"] = ""
		imageAttr["title"] = ""
	}

	imageAttr["zoom"] = "50"
	imageAttr["pos"] = "left"

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
	textNode := &PMNode{
		Type: "text",
		Text: node.Tag,
	}
	paragraph.Content = append(paragraph.Content, textNode)
	return paragraph
}

func (p *ListParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	bulletList := &PMNode{
		Type:    "bulletList",
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
				Type:    "listItem",
				Content: []*PMNode{},
			}
			findTextTill(node, listItem)
			bulletList.Content = append(bulletList.Content, listItem)
		}
		for _, child := range node.Children {
			s.Push(child)
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

	//
	// this function is to find the table cell when reach the tr tag
	//
	// the reason why we need to do this is because the text node is not the direct child of the tr tag
	//
	// so we need an extra funtion to find the table cell node
	//
	sf := func(node *AstElement, parent *PMNode) {
		s := NewStack()
		s.Push(node)
		for {
			if s.Len() == 0 {
				break
			}
			node := s.Pop().(*AstElement)
			if node.Tag == "td" {
				tableCell := &PMNode{
					Type: "tableCell",
					Attrs: map[string]string{
						"colspan":  "1",
						"rowspan":  "1",
						"colwidth": "",
					},
					Content: []*PMNode{},
				}
				findTextTill(node, tableCell)
				parent.Content = append(parent.Content, tableCell)
			}
			for _, child := range node.Children {
				s.Push(child)
			}
		}
	}

	s := NewStack()
	s.Push(node)
	for {
		if s.Len() == 0 {
			break
		}
		node := s.Pop().(*AstElement)
		if node.Tag == "tr" {
			tableRow := &PMNode{
				Type:    "tableRow",
				Content: []*PMNode{},
			}
			sf(node, tableRow)
			table.Content = append(table.Content, tableRow)
		}
		for _, child := range node.Children {
			s.Push(child)
		}
	}
	return table
}

func (p *SVGParser) Parse(node *AstElement) *PMNode {
	if node == nil {
		return nil
	}
	imageAttr := make(map[string]string)
	if v, ok := node.Attrs["xlink:href"]; ok { // compatible to the svg
		_, imageName := filepath.Split(v)
		imageAttr["src"] = fmt.Sprintf("http://127.0.0.1:7288/%s", imageName)
	}
	if v, ok := node.Attrs["href"]; ok {
		_, imageName := filepath.Split(v)
		imageAttr["src"] = fmt.Sprintf("http://127.0.0.1:7288/%s", imageName)
	}
	if v, ok := node.Attrs["alt"]; ok {
		imageAttr["alt"] = v
		imageAttr["title"] = v
	} else {
		imageAttr["alt"] = ""
		imageAttr["title"] = ""
	}

	imageAttr["zoom"] = "50"
	imageAttr["pos"] = "left"

	image := &PMNode{
		Type:  "image",
		Attrs: imageAttr,
	}
	return image
}
