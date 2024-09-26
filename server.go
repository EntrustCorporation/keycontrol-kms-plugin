package main

import (
	"flag"
	"kms-plugin/grpcservice"
	"kms-plugin/kmsservice"
	"kms-plugin/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	utils.Logger.Level = logrus.DebugLevel
	sockFile := flag.String("sockFile", "", "unix Domain socket that gRpc server listens to")
	confFile := flag.String("confFile", "", "config File location for HyTrust KMS plugin")
	flag.Parse()

	if len(*sockFile) == 0 {
		utils.Logger.Fatal("No sockFile specified")
		//*sockFile = "/var/tmp/run.sock"
	}

	if len(*confFile) == 0 {
		utils.Logger.Fatal("No config file specified")
		//*confFile = "./confFile"
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	config := utils.LoadConfig(*confFile)
	keyControlKmsService, err := kmsservice.NewKeyControlKmsService(&config)
	if err != nil {
		utils.Logger.Fatal("failed to initialize kms service", err)
	}
	grpcService := grpcservice.NewGRPCService(*sockFile, 5*time.Minute, keyControlKmsService)

	go func() {
		if err := grpcService.ListenAndServe(); err != nil {
			utils.Logger.Fatal("failed to  run plugin", err)
		}
	}()
	<-sigChan
	grpcService.Shutdown()
}
