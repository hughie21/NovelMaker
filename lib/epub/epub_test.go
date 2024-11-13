package epub

import (
	"io"
	"os"
	"testing"
)

func TestReader(t *testing.T) {
	t.Error("Test Reader")
	reader, err := NewReader("D:/test.epub", "./")
	if err != nil {
		t.Errorf("NewReader() failed at %s: %v", t.Name(), err)
	}
	err = reader.Read()
	if err != nil {
		t.Fatalf("Read() failed at %s: %v", t.Name(), err)
	}
	err = reader.Pharse()
	if err != nil {
		t.Fatalf("Pharse() failed at %s: %v", t.Name(), err)
	}
	t.Log(reader.Content)
	err = reader.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriter(t *testing.T) {
	var jsonData JsonData
	fs, err := os.Open("raw.json")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()
	rawData, err := io.ReadAll(fs)
	if err != nil {
		t.Fatal(err)
	}
	LoadJson(rawData, &jsonData)
	w := NewWriter("test.epub", "./", &jsonData)
	err = w.Write()
	if err != nil {
		t.Error(err)
	}
	w.Close()
}
