// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	initialize "github.com/dvonthenen/chat-gpeasy/pkg/initialize"
	personas "github.com/dvonthenen/chat-gpeasy/pkg/personas"
	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
)

func main() {
	initialize.Init(initialize.ChatGPTProxyInit{
		LogLevel: initialize.LogLevelStandard, // LogLevelStandard / LogLevelTrace
	})

	// create the chatgpt client
	fmt.Printf("Connecting to Generative AI...\n")
	personaConfig, err := personas.DefaultConfig("", "")
	if err != nil {
		fmt.Printf("personas.DefaultConfig error: %v\n", err)
		os.Exit(1)
	}

	persona, err := personas.NewAdvancedChatWithOptions(personaConfig)
	if err != nil {
		fmt.Printf("personas.NewAdvancedChatWithOptions error: %v\n", err)
		os.Exit(1)
	}
	// OR
	// persona, err := personas.NewAdvancedChat()
	// if err != nil {
	// 	fmt.Printf("personas.NewCumulativeChat error: %v\n", err)
	// 	os.Exit(1)
	// }

	(*persona).Init(interfaces.SkillTypeGeneric, "")

	ctx := context.Background()

	// prompt 1
	prompt := "Hello! How are you doing?"
	choices, err := (*persona).Query(ctx, openai.ChatMessageRoleUser, prompt)
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
	choices, err = (*persona).Query(ctx, openai.ChatMessageRoleUser, prompt)
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

	// edit convo
	fmt.Printf("Oooops... I goofed. I need to edit this...\n\n\n")
	conversation, err := (*persona).GetConversation()
	if err != nil {
		fmt.Printf("persona.GetConversation error: %v\n", err)
		os.Exit(1)
	}

	for pos, msg := range conversation {
		if strings.Contains(msg.Content, "Long Beach, CA") {
			prompt = "Tell me about Laguna Beach, CA."
			choices, err := (*persona).EditConversation(pos, prompt)
			if err != nil {
				fmt.Printf("persona.EditConversation error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Me:\n%s\n", prompt)
			fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
		}
	}

	// divider
	fmt.Printf("\n\n\n")
	fmt.Printf("-------------------------------------------")
	fmt.Printf("\n\n\n")

	// prompt 3
	prompt = "I want more factual type data"
	fmt.Printf("Adding clarifying directives to AI:\n%s\n", prompt)
	err = (*persona).AddDirective(prompt)
	if err != nil {
		fmt.Printf("persona.AddDirective error: %v\n", err)
		os.Exit(1)
	}

	// divider
	fmt.Printf("\n\n\n")
	fmt.Printf("-------------------------------------------")
	fmt.Printf("\n\n\n")

	prompt = "Now... tell me about Laguna Beach, CA."
	choices, err = (*persona).Query(ctx, openai.ChatMessageRoleUser, prompt)
	if err != nil {
		fmt.Printf("persona.Query error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Me:\n%s\n", prompt)
	fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
}
