// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package main

// streaming
import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	microphone "github.com/dvonthenen/symbl-go-sdk/pkg/audio/microphone"
	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"

	handler "github.com/dvonthenen/chat-gpeasy/examples/symbl/handler"
	initialize "github.com/dvonthenen/chat-gpeasy/pkg/initialize"
	personas "github.com/dvonthenen/chat-gpeasy/pkg/personas"
	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
)

func main() {
	// init the library
	initialize.Init(initialize.ChatGPTProxyInit{
		LogLevel: initialize.LogLevelStandard, // LogLevelStandard / LogLevelTrace
	})

	// context
	ctx := context.Background()

	// create the chatgpt client
	fmt.Printf("Connecting to Generative AI...\n")
	personaConfig, err := personas.DefaultConfig("", "")
	if err != nil {
		fmt.Printf("personas.DefaultConfig error: %v\n", err)
		os.Exit(1)
	}

	persona, err := personas.NewSimpleChatWithOptions(personaConfig)
	if err != nil {
		fmt.Printf("personas.NewSimpleChatWithOptions error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Connection Succeeded\n")

	// persona, err := personas.NewSimpleChat()
	// if err != nil {
	// 	fmt.Printf("personas.NewSimple error: %v\n", err)
	// 	os.Exit(1)
	// }

	(*persona).Init(interfaces.SkillTypeGeneric, "")

	// init library
	microphone.Initialize()

	// init the handler
	msgHandler := handler.NewHandler(handler.HandlerOptions{
		Simple: persona,
	})

	// create a new client
	symblConfig := symbl.GetDefaultConfig()
	symblConfig.Speaker.Name = "John Doe"
	symblConfig.Speaker.UserID = "john.doe@mymail.com"
	symblConfig.Config.DetectEntities = true
	symblConfig.Config.Sentiment = true

	options := symbl.StreamingOptions{
		SymblConfig: symblConfig,
		Callback:    msgHandler,
	}

	client, err := symbl.NewStreamClient(ctx, options)
	if err == nil {
		fmt.Println("Login Succeeded!")
	} else {
		fmt.Printf("New failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("ConversationID: %s\n", client.GetConversationId())

	err = client.Start()
	if err == nil {
		fmt.Printf("Streaming Session Started!\n")
	} else {
		fmt.Printf("client.Start failed. Err: %v\n", err)
		os.Exit(1)
	}

	// delay...
	time.Sleep(time.Second * 1)

	// mic stuf
	mic, err := microphone.New(microphone.AudioConfig{
		InputChannels: 1,
		SamplingRate:  16000,
	})
	if err != nil {
		fmt.Printf("Initialize failed. Err: %v\n", err)
		os.Exit(1)
	}

	// start the mic
	err = mic.Start()
	if err != nil {
		fmt.Printf("mic.Start failed. Err: %v\n", err)
		os.Exit(1)
	}

	go func() {
		// this is a blocking call
		mic.Stream(client)
	}()

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// close stream
	err = mic.Stop()
	if err != nil {
		fmt.Printf("mic.Stop failed. Err: %v\n", err)
		os.Exit(1)
	}

	// teardown library
	microphone.Teardown()

	// close client
	client.Stop()

	fmt.Printf("Succeeded!\n\n")
}
