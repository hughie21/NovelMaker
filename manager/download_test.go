package manager

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImageDownloader_Download(t *testing.T) {
	// Create a test server with a handler that returns a sample image
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sample image data"))
	}))
	t.Log("Test server started on", ts.URL)
	defer ts.Close()

	downloader := NewImageDownloader("", 20)
	err := ProcessDownload(downloader, ts.URL)
	if err != nil {
		t.Error(err)
	}
}
