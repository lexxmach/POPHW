package main

import (
	"bufio"
	"context"
	"fmt"
	"msg/internal"
	"msg/internal/configs"
	"msg/internal/sockets"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func getLine(header string, reader bufio.Reader) (string, error) {
	if header != "" {
		fmt.Println(header)
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return line, nil
}

func getToken(header string) (string, error) {
	if header != "" {
		fmt.Println(header)
	}
	var token string
	_, err := fmt.Scanf("%s", &token)
	if err != nil {
		return "", err
	}

	return token, err
}

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	reader := bufio.NewReader(os.Stdin)

	config := getConfig()

	user, err := getToken("Input user name: ")
	if err != nil {
		panic(err)
	}

	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	socket, err := sockets.NewClientSocket(context.Background(), user, url.URL{Scheme: "ws", Host: config.ServerAddress, Path: "/"}, func(msg *sockets.ClientMessage) {
		if msg.Error != nil {
			logger.Errorf("Encounterd error: %q", msg.Error)
		} else {
			logger.Infof("%s: %s", msg.Msg.Message.User, msg.Msg.Message.Message)
		}
	})
	if err != nil {
		logger.Error(fmt.Errorf("host: %q failed to connect: %w", config.ServerAddress, err))
		return
	}

	go func() {
		<-irqSig
		err := socket.Stop()
		if err != nil {
			fmt.Printf("Failed to shutdown client: %q", err)
			panic(err)
		}
		os.Exit(0)
	}()

	for {
		msg, err := getLine("", *reader)
		if err != nil {
			fmt.Printf("Encountered error: %q", err)
			continue
		}

		err = socket.SendMessage(msg)
		if err != nil {
			fmt.Printf("Failed to deliver message: %q", err)
			break
		}
	}
}

func getConfig() *configs.ClientConfig {
	configPath := pflag.StringP("config", "c", "", "Path to config")
	pflag.Parse()

	return internal.Must(configs.GetClientConfig(*configPath))
}
