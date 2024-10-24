package manager

import (
	"testing"
)

func TestConfigManager_LoadConfig(t *testing.T) {
	cf := NewConfigManager("D:/NovelMaker")
	cf.LoadConfig()

	t.Log(cf.GetConfigByKey("Appearance", "Width"))
}
