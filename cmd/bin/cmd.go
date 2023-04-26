// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package main

import (
	"bufio"
	"fmt"
	"os"

	initproxy "github.com/dvonthenen/chat-gpeasy/pkg/initialize"
	chatgptproxy "github.com/dvonthenen/chat-gpeasy/pkg/proxy"
	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/proxy/interfaces"
)

func main() {
	initproxy.Init(initproxy.ChatGPTProxyInit{
		LogLevel: initproxy.LogLevelStandard,
	})

	var callback interfaces.ChatGPTCallback
	callback = chatgptproxy.NewDefaultChatGPTCallback()

	proxyServer, err := chatgptproxy.New(chatgptproxy.ProxyOptions{
		Callback: &callback,
		CrtFile:  "localhost.crt",
		KeyFile:  "localhost.key",
	})
	if err != nil {
		fmt.Printf("server.New failed. Err: %v\n", err)
		os.Exit(1)
	}

	// init
	err = proxyServer.Init()
	if err != nil {
		fmt.Printf("proxyServer.Init() failed. Err: %v\n", err)
		os.Exit(1)
	}

	// start
	fmt.Printf("Starting server...\n")
	err = proxyServer.Start()
	if err != nil {
		fmt.Printf("proxyServer.Start() failed. Err: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// stop
	err = proxyServer.Stop()
	if err != nil {
		fmt.Printf("proxyServer.Stop() failed. Err: %v\n", err)
	}

	// teardown
	err = proxyServer.Teardown()
	if err != nil {
		fmt.Printf("proxyServer.Stop() failed. Err: %v\n", err)
	}

	fmt.Printf("Succeeded!\n\n")
}
