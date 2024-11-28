package html

import (
	"regexp"
	"strings"

	"github.com/beevik/etree"
)

type AstElement struct {
	Tag      string
	Type     int
	Children []*AstElement
	Parent   *AstElement
	Attrs    map[string]string
}

var (
	ncname            = `[a-zA-Z_][\-\.0-9_a-zA-Z]*`
	qnameCapture      = `((?:` + ncname + `\:)?` + ncname + `)`
	startTagOpen      = regexp.MustCompile(`^<` + qnameCapture)
	endTag            = regexp.MustCompile(`^<\/` + qnameCapture + `[^>]*>`)
	attribute         = regexp.MustCompile(`^\s*([^\s"'<>\/=]+)(?:\s*(=)\s*(?:"([^"]*)"+|'([^']*)'+|([^\s"'=<>` + "`" + `]+)))?`)
	startTagClose     = regexp.MustCompile(`^\s*(\/?)>`)
	uselessCharacters = regexp.MustCompile(`\n+(\t+)?`)
)

type HTMLParser struct {
	root  *AstElement
	stack []*AstElement
}

func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

func (p *HTMLParser) createAstElement(tagName string, attrs map[string]string, parent *AstElement) *AstElement {
	return &AstElement{
		Tag:      tagName,
		Type:     1,
		Children: []*AstElement{},
		Parent:   parent,
		Attrs:    attrs,
	}
}

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

func (p *HTMLParser) end(tagName string) {
	last := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	if last.Tag != tagName {
		panic("Error: tag not match")
	}
}

func (p *HTMLParser) chars(text string) {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", "")
	text = uselessCharacters.ReplaceAllString(text, "")
	var parent *AstElement
	if len(p.stack) > 0 {
		parent = p.stack[len(p.stack)-1]
	}
	if text != "" {
		parent.Children = append(parent.Children, &AstElement{
			Type:   3,
			Tag:    text,
			Parent: parent,
		})
	}
}

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

func LoadHTML(html []byte) *AstElement {
	textDoc := etree.NewDocument()
	if err := textDoc.ReadFromBytes(html); err != nil {
		panic(err)
	}
	body := textDoc.FindElement("//body")
	if body == nil {
		panic("no body element found")
	}
	bodyDoc := etree.NewDocument()
	bodyDoc.SetRoot(body.Copy())
	rawText, _ := bodyDoc.WriteToBytes()
	parser := NewHTMLParser()
	ast := parser.parserHTML(string(rawText))
	return ast
}
