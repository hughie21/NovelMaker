// Description: Chnage the html string to the ast tree
// Author: Hughie21
// Date: 2024-11-29
// license that can be found in the LICENSE file.

package html

import (
	"errors"
	"regexp"
	"strings"

	"github.com/beevik/etree"
)

// the ast tree nodes
type AstElement struct {
	Tag      string
	Type     int
	Children []*AstElement
	Parent   *AstElement
	Attrs    map[string]string
	Text     string
}

var (
	ncname            = `[a-zA-Z_][\-\.0-9_a-zA-Z]*`
	qnameCapture      = `((?:` + ncname + `\:)?` + ncname + `)`
	startTagOpen      = regexp.MustCompile(`^<` + qnameCapture)
	endTag            = regexp.MustCompile(`^<\/` + qnameCapture + `[^>]*>`)
	attribute         = regexp.MustCompile(`^\s*([^\s"'<>\/=]+)(?:\s*(=)\s*(?:"([^"]*)"+|'([^']*)'+|([^\s"'=<>` + "`" + `]+)))?`)
	startTagClose     = regexp.MustCompile(`^\s*(\/?)>`)
	uselessCharacters = regexp.MustCompile(`>(\s*\n\s*|\s*\t\s*)<`)
)

// HTMLParser is a struct that contains the root element of the AST tree
// and a stack to keep track of the current element being processed.
type HTMLParser struct {
	// root element of the AST tree
	root *AstElement
	// stack to keep track of the current element being processed
	stack []*AstElement
}

// NewHTMLParser creates a new HTMLParser object and returns a pointer to it.
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

// createAstElement creates a new AstElement object with the given tag name, attributes, and parent element.
func (p *HTMLParser) createAstElement(tagName string, attrs map[string]string, parent *AstElement) *AstElement {
	return &AstElement{
		Tag:      tagName,
		Type:     1,
		Children: []*AstElement{},
		Parent:   parent,
		Attrs:    attrs,
		Text:     "",
	}
}

// Iterate through the characters to find the label that matches the start flag
func (p *HTMLParser) start(tagName string, attributes map[string]string, unary bool) {
	var parent *AstElement
	if len(p.stack) > 0 {
		parent = p.stack[len(p.stack)-1]
	}
	element := p.createAstElement(tagName, attributes, parent)
	if p.root == nil {
		p.root = element
	}
	if parent != nil {
		parent.Children = append(parent.Children, element)
	}
	if !unary {
		p.stack = append(p.stack, element)
	}
}

// Match End tag
func (p *HTMLParser) end(tagName string) {
	last := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	if last.Tag != tagName {
		panic("Error: tag not match")
	}
}

// Handling text between tags
func (p *HTMLParser) chars(text string) {
	// text = strings.TrimSpace(text)
	// text = strings.ReplaceAll(text, "\n", "")
	// text = uselessCharacters.ReplaceAllString(text, "")
	re := regexp.MustCompile(`&nbsp;|&ensp;|&emsp;|&thinsp;|&zwnj;|&zwj;|&lrm;|&rlm;`)
	text = re.ReplaceAllString(text, "")
	var parent *AstElement
	if len(p.stack) > 0 {
		parent = p.stack[len(p.stack)-1]
	}
	if text != "" {
		parent.Children = append(parent.Children, &AstElement{
			Type:   3,
			Text:   text,
			Tag:    parent.Tag,
			Attrs:  parent.Attrs,
			Parent: parent,
		})
	}
}

// parserHTML parses the given HTML string and returns an AstElement object.
//
// The main flow of this method is as follows:
//
//  1. Define an inner function, advance, that advances the parsing of the HTML string.
//  2. Define an inner function, parseStartTag, that parses the start tag and returns the tag name, attributes, and whether it is a single tag.
//  3. In the main loop, check to see if the length of the HTML string is greater than zero.
//  4. If the HTML string starts with “<”, try to parse the start tag.
//     - If parsing is successful, call the p.start method to process the start tag and continue to the next loop.
//     - If parsing fails, try to parse the end tag.
//     - If parsing succeeds, call p.end to process the end tag and advance the parsing process.
//  5. If the HTML string does not begin with “<”, extract the text and call the p.chars method to process it.
//  6. At the end of the loop, return the parsed root element, p.root.
func (p *HTMLParser) parserHTML(html string) *AstElement {
	advance := func(length int) {
		html = html[length:]
	}

	parseStartTag := func() (map[string]interface{}, bool) {
		start := startTagOpen.FindStringSubmatch(html)
		if start != nil {
			match := map[string]interface{}{
				"tagName": start[1],
				"attrs":   map[string]string{},
				"unary":   false,
			}
			advance(len(start[0]))
			var end []string
			var attr []string
			for end = startTagClose.FindStringSubmatch(html); end == nil; end = startTagClose.FindStringSubmatch(html) {
				attr = attribute.FindStringSubmatch(html)
				if attr == nil {
					break
				}
				match["attrs"].(map[string]string)[attr[1]] = attr[3]
				advance(len(attr[0]))
			}
			if end != nil {
				match["unary"] = end[1] == "/"
				advance(len(end[0]))
			}
			return match, true
		}
		return nil, false
	}

	for len(html) > 0 {
		textEnd := strings.Index(html, "<")
		if textEnd == 0 {
			startTagMatch, ok := parseStartTag()
			if ok {
				p.start(startTagMatch["tagName"].(string), startTagMatch["attrs"].(map[string]string), startTagMatch["unary"].(bool))
				continue
			}
			endTagMatch := endTag.FindStringSubmatch(html)
			if endTagMatch != nil {
				p.end(endTagMatch[1])
				advance(len(endTagMatch[0]))
				continue
			}
		}
		var text string
		if textEnd > 0 {
			text = html[:textEnd]
		}
		if text != "" {
			p.chars(text)
			advance(len(text))
		}
	}

	return p.root
}

// LoadHTML loads the given HTML byte slice and returns an AstElement object.
func LoadHTML(html []byte) (*AstElement, error) {
	textDoc := etree.NewDocument()
	if err := textDoc.ReadFromBytes(html); err != nil {
		return nil, err
	}
	body := textDoc.FindElement("//body")
	if body == nil {
		return nil, errors.New("no body element found")
	}
	bodyDoc := etree.NewDocument()
	bodyDoc.SetRoot(body.Copy())
	rawText, _ := bodyDoc.WriteToBytes()
	parser := NewHTMLParser()
	rawText = uselessCharacters.ReplaceAll(rawText, []byte("><"))
	ast := parser.parserHTML(string(rawText))
	return ast, nil
}
