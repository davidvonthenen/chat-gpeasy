// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	var openAiApiKey string
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		openAiApiKey = v
	} else {
		fmt.Printf("OPENAI_API_KEY not found\n")
		os.Exit(1)
	}

	config := openai.DefaultConfig(openAiApiKey)
	config.BaseURL = "https://127.0.0.1/v1"
	config.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	ctx := context.Background()
	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello! How are you doing?",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
