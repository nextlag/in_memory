package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/fatih/color"
	config "github.com/nextlag/in_memory/config/client"
	"github.com/nextlag/in_memory/internal"
	"github.com/nextlag/in_memory/internal/client/network"
	"github.com/nextlag/in_memory/pkg/parse"
)

func main() {
	cfg := config.Load()

	maxMessageSize, err := parse.Size(cfg.Network.MaxMessageSize)
	if err != nil {
		color.Red("failed to parse max message size %v", err)
	}

	var options []network.TCPClientOption
	options = append(options, network.WithClientIdleTimeout(cfg.Network.IdleTimeout))
	options = append(options, network.WithClientBufferSize(uint(maxMessageSize)))

	reader := bufio.NewReader(os.Stdin)
	client, err := network.NewTCPClient(cfg.Network.ServerAddress, options...)
	if err != nil {
		color.Red("failed to connect with server %v", err)
	}

	for {
		fmt.Print("> ")
		request, err := reader.ReadString('\n')
		if errors.Is(err, syscall.EPIPE) {
			color.Red("connection was closed %v", err)
		} else if err != nil {
			color.Red("failed to read query %v", err)
		}
		request = strings.TrimSpace(request)

		response, err := client.Send([]byte(request))
		if errors.Is(err, syscall.EPIPE) {
			color.Red("connection was closed %v", err)
		} else if err != nil {
			color.Red("failed to send query %v", err)
		}

		responseStr := string(response)

		if strings.Contains(responseStr, internal.ResponseOk) {
			color.Green(responseStr)
			continue
		}
		color.Red(responseStr)
	}
}
