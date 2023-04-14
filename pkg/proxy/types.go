// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package proxy

import (
	"net/http"

	openai "github.com/sashabaranov/go-openai"

	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/proxy/interfaces"
)

// ProxyOptions for the main HTTP endpoint
type ProxyOptions struct {
	Callback *interfaces.ChatGPTCallback
	CrtFile  string
	KeyFile  string
	BindPort int
}

type ChatGPTProxy struct {
	options *ProxyOptions

	// callback
	callback *interfaces.ChatGPTCallback

	// server
	server *http.Server

	// openai
	openAiApiKey  string
	chatgptClient *openai.Client
}
