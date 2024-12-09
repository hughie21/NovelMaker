package html

import (
	"sync"
)

type (
	Stack struct {
		top    *node
		length int
		lock   *sync.RWMutex
	}
	node struct {
		value interface{}
		prev  *node
	}
	PMNode struct {
		Type    string                 `json:"type"`
		Attrs   map[string]interface{} `json:"attrs,omitempty"`
		Content []*PMNode              `json:"content,omitempty"`
		Text    string                 `json:"text,omitempty"`
		Mark    []*PMNode              `json:"marks,omitempty"`
	}
	ParserContext struct {
		parsers map[string]TagParser
	}
)

// Create a new stack
func NewStack() *Stack {
	return &Stack{nil, 0, &sync.RWMutex{}}
}

// Return the number of items in the stack
func (s *Stack) Len() int {
	return s.length
}

// View the top item on the stack
func (s *Stack) Peek() interface{} {
	if s.length == 0 {
		return nil
	}
	return s.top.value
}

// Pop the top item of the stack and return it
func (s *Stack) Pop() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.length == 0 {
		return nil
	}
	n := s.top
	s.top = n.prev
	s.length--
	return n.value
}

// Push a value onto the top of the stack
func (s *Stack) Push(value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	n := &node{value, s.top}
	s.top = n
	s.length++
}

func NewPMNode() *PMNode {
	return &PMNode{
		Type:    "",
		Attrs:   make(map[string]interface{}),
		Content: []*PMNode{},
	}
}

func NewParserContext() *ParserContext {
	return &ParserContext{
		parsers: make(map[string]TagParser),
	}
}

func (c *ParserContext) Register(tag string, parser TagParser) {
	c.parsers[tag] = parser
}

func (c *ParserContext) Parse(node *AstElement) *PMNode {
	if parser, ok := c.parsers[node.Tag]; ok {
		return parser.Parse(node)
	}
	if node.Type == 3 {
		if node.Tag == "span" || node.Tag == "p" {
			return c.parsers["text"].Parse(node)
		}
	}
	return nil
}

func ConvertIntoProseMirrorScheme(root *AstElement, Parsers map[string]TagParser) *PMNode {
	doc := &PMNode{
		Type:    "doc",
		Content: []*PMNode{},
	}

	if root == nil {
		return nil
	}
	context := NewParserContext()
	for tag, parser := range Parsers {
		context.Register(tag, parser)
	}

	s := NewStack()
	s.Push(root)
	for {
		if s.Len() == 0 {
			break
		}
		node := s.Pop().(*AstElement)
		res := context.Parse(node)
		if res != nil {
			doc.Content = append([]*PMNode{res}, doc.Content...)
		}
		for _, child := range node.Children {
			if _, ok := Parsers[node.Tag]; ok {
				break
			}
			s.Push(child)
		}
	}
	return doc
}
