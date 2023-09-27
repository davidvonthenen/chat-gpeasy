// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package advanced

import (
	"strings"
	"sync"

	openai "github.com/sashabaranov/go-openai"

	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
)

type CommitCallback interface {
	CommitResponse(response string) error
}

type StreamingChatCompletion struct {
	sb       strings.Builder
	stream   *openai.ChatCompletionStream
	stopChan chan struct{}
	callback *CommitCallback
}

type Persona struct {
	// openai
	client *openai.Client

	// options
	model string
	level interfaces.SkillType

	// last query
	appendedResponse bool
	conversation     []openai.ChatCompletionMessage
	request          *openai.ChatCompletionRequest
	response         *openai.ChatCompletionResponse

	// housekeeping
	mu sync.Mutex
}
