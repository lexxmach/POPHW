package main

import (
	"fmt"
	"msg/internal"
	"msg/internal/configs"
	"msg/internal/sockets"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	config := getConfig()
	storage := internal.Must(config.DB.GetStorage())

	socket := sockets.NewServerSocket(config.Address, storage, nil, func(msg *sockets.ServerMessage) {
		if msg.Error != nil {
			logger.Errorf("Encounterd error: %q", msg.Error)
		} else {
			logger.Infof("User %s: %s", msg.Msg.Message.User, msg.Msg.Message.Message)
		}
	})

	go func() {
		<-irqSig
		err := socket.Stop()
		if err != nil {
			fmt.Printf("Failed to shutdown server: %q", err)
		}
	}()

	logger.Infof("Starting server on host: %q", config.Address)

	if err := socket.Start(); err != nil {
		fmt.Printf("Server stopped with response: %q", err)
		return
	}
}

func getConfig() *configs.ServerConfig {
	configPath := pflag.StringP("config", "c", "", "Path to config")
	pflag.Parse()

	return internal.Must(configs.GetServerConfig(*configPath))
}
