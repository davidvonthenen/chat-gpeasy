// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package advanced

import (
	"sync"

	openai "github.com/sashabaranov/go-openai"

	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
)

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
