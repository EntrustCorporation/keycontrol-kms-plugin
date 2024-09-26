package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	//Create a new instance of logrus logger
	// f, err := os.OpenFile("/var/log/kmsplugin.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	// if err != nil {
	// 	fmt.Printf("error opening file: %v", err)
	// }
	Logger = logrus.New()

	// Configure logger outputs
	Logger.SetOutput(os.Stdout) // Log to stdout

	//Logger.SetOutput(f)
	// Optionally set a formatter
	Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false, // Disable colors in log output
		FullTimestamp: true,  // Show full timestamp
	})
}
