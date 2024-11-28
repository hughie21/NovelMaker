package html

import (
	"encoding/json"
	"testing"
)

func TestAst(t *testing.T) {
	testContent := `<body><p>dawda</p><br></br><p>dawda</p><br></br><br></br><br></br><p>dwada</p><br></br><br></br><span style='display:flex;justify-content:left;'><img src='../Images/b882997817cfad8d.jpg' alt='b882997817cfad8d.jpg' title='b882997817cfad8d.jpg' style='width:50%;'/></span><br></br><br></br><span style='display:flex;justify-content:left;'><img src='../Images/0164b793d2f9252d.jpg' alt='0164b793d2f9252d.jpg' title='0164b793d2f9252d.jpg' style='width:50%;'/></span><br></br><br></br></body>`
	ast := LoadHTML([]byte(testContent))
	t.Error("test")
	jsonNodes := ConvertIntoProseMirrorScheme(ast, map[string]TagParser{})
	jsonData, _ := json.Marshal(jsonNodes)
	t.Log(string(jsonData))
}
