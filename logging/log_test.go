package logging

import (
	"testing"
)

func TestLogging(t *testing.T) {
	var logger = NewLog(FatalLevel, true)
	logger.Error("Error", "ErrorTest")
	logger.Info("Info", "ErrorTest")
	logger.Debug("Debug", "DebugTest")
	logger.Fatal("Fatal", "FatalTest")

	var logger2 = NewLog(DebugLevel, true)
	logger2.Warning("Warning", "WarningTest")

	logger.LogOutPut("./")
}
