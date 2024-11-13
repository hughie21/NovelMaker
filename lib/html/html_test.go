package html

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestAstElement(t *testing.T) {
	fs, _ := os.Open("text.xhtml")
	defer fs.Close()
	rawData, _ := io.ReadAll(fs)
	ast := LoadHTML(rawData)
	// json_data, _ := json.Marshal(ast)
	// t.Logf("Ast JSON: %s", json_data)
	doc := ConvertIntoProseMirrorScheme(ast, map[string]TagParser{})
	json_data, _ := json.Marshal(doc)
	t.Logf("%s", json_data)
}
