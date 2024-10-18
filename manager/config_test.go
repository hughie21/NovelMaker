package manager

import (
	"testing"
)

func TestConfigManager_LoadConfig(t *testing.T) {
	cf := NewConfigManager("D:/NovelMaker")
	cf.LoadConfig()

	cf.SetConfig("Appearance", "Theme", "light")
}
