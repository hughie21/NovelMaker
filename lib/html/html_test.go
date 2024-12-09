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
		<h1>
                序章
        </h1>
        <p>
                在寂静的店内响起的门铃声，让他抬起了头。
        </p>
	</body>
	`
	ast, err := LoadHTML([]byte(testContent))
	if err != nil {
		t.Error(err)
		return
	}
	t.Error("test")
	jsonNodes := ConvertIntoProseMirrorScheme(ast, map[string]TagParser{
		"text": &TextParser{},
		// "span":  &TextParser{},
		"table": &TableParser{},
		"h1":    &HeaderParser{Level: 1},
		"h2":    &HeaderParser{Level: 1},
		"h3":    &HeaderParser{Level: 1},
		"img":   &ImageParser{},
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
