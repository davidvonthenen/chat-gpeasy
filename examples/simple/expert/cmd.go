// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package main

import (
	"context"
	"fmt"
	"os"

	initialize "github.com/dvonthenen/chat-gpeasy/pkg/initialize"
	personas "github.com/dvonthenen/chat-gpeasy/pkg/personas"
	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
)

func main() {
	initialize.Init(initialize.ChatGPTProxyInit{
		LogLevel: initialize.LogLevelStandard, // LogLevelStandard / LogLevelTrace
	})

	persona, err := personas.NewSimpleChat()
	if err != nil {
		fmt.Printf("personas.NewSimpleChat error: %v\n", err)
		os.Exit(1)
	}

	(*persona).Init(interfaces.SkillTypeExpert, "")

	ctx := context.Background()

	// prompt 1
	prompt := "Hello! How are you doing?"
	choices, err := (*persona).Query(ctx, prompt)
	if err != nil {
		fmt.Printf("persona.Query error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Me:\n%s\n", prompt)
	fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)

	// divider
	fmt.Printf("\n\n\n")
	fmt.Printf("-------------------------------------------")
	fmt.Printf("\n\n\n")

	// prompt 2
	prompt = "Tell me about Long Beach, CA."
	choices, err = (*persona).Query(ctx, prompt)
	if err != nil {
		fmt.Printf("persona.Query error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Me:\n%s\n", prompt)
	fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
}
