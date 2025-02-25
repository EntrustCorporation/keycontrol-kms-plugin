package utils

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// CallerHook is a Logrus hook that adds function name and line number
type CallerHook struct{}

// Fire is called for every log entry and adds caller info
func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	pc, _, line, ok := runtime.Caller(8) // Adjust depth based on log call stack
	if !ok {
		entry.Data["caller"] = "unknown"
		return nil
	}
	fn := runtime.FuncForPC(pc)
	entry.Data["caller"] = fmt.Sprintf("%s:%d", fn.Name(), line)
	return nil
}

// Levels returns the log levels the hook should be applied to
func (hook *CallerHook) Levels() []logrus.Level {
	return logrus.AllLevels // Apply to all log levels
}

func init() {
	//Create a new instance of logrus logger
	// f, err := os.OpenFile("/var/log/kmsplugin.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	// if err != nil {
	// 	fmt.Printf("error opening file: %v", err)
	// }
	Logger = logrus.New()

	Logger.AddHook(&CallerHook{})

	// Configure logger outputs
	Logger.SetOutput(os.Stdout) // Log to stdout

	//Logger.SetOutput(f)
	// Optionally set a formatter
	Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false, // Disable colors in log output
		FullTimestamp: true,  // Show full timestamp
	})
}
