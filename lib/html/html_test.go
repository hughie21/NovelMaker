package html

import (
	"encoding/json"
	"testing"
)

func TestAst(t *testing.T) {
	//{"type":"doc","content":[{"type":"codeBlock","attrs":{"language":"go"},"content":[{"type":"text","text":"var test string"}]}]}
	//
	testContent := `
	<body>
		<p>a <strong>bold</strong> <em>italic</em> <s>strike</s></p>
	</body>
	`
	ast, err := LoadHTML([]byte(testContent))
	if err != nil {
		t.Error(err)
		return
	}
	t.Error("test")
	jsonNodes := ConvertIntoProseMirrorScheme(ast, map[string]TagParser{
		"p":     &TextParser{},
		"span":  &TextParser{},
		"table": &TableParser{},
		"h1":    &HeaderParser{Level: 1},
		"ol": &ListParser{
			Type: "orderedList",
		},
		"ul": &ListParser{
			Type: "bulletList",
		},
		"code": &CodeBlockParser{},
	})
	jsonData, _ := json.Marshal(jsonNodes)
	t.Log(string(jsonData))
}
