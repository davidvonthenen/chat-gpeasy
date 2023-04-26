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
	c := openai.NewClientWithConfig(config)

	modelList, err := c.ListModels(ctx)
	if err != nil {
		fmt.Printf("c.ListModels failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, model := range modelList.Models {
		fmt.Printf("Model: %s\n", model.ID)
	}
}
