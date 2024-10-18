package sys

import (
	"testing"
)

func TestWindowSysMessageBox(t *testing.T) {
	ShowMessage("test", "test", "info")
}
