// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package utils

import (
	openai "github.com/sashabaranov/go-openai"

	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
)

// CompletionMessage: Chat-GPeasy to OpenAI
func ConvertCompletionMessage(message interfaces.CompletionMessage) *openai.ChatCompletionMessage {
	converted := openai.ChatCompletionMessage{
		Role:    message.Role,
		Content: message.Content,
		Name:    message.Name,
	}
	return &converted
}

func ConvertCompletionMessages(messages []interfaces.CompletionMessage) *[]openai.ChatCompletionMessage {
	converted := make([]openai.ChatCompletionMessage, 0)
	for _, message := range messages {
		converted = append(converted, *ConvertCompletionMessage(message))
	}
	return &converted
}

// CompletionMessage: OpenAI to Chat-GPeasy
func ConvertChatCompletionMessage(message openai.ChatCompletionMessage) *interfaces.CompletionMessage {
	cm := interfaces.CompletionMessage{
		Role:    message.Role,
		Content: message.Content,
		Name:    message.Name,
	}
	return &cm
}

func ConvertChatCompletionMessages(messages []openai.ChatCompletionMessage) *[]interfaces.CompletionMessage {
	converted := make([]interfaces.CompletionMessage, 0)
	for _, message := range messages {
		converted = append(converted, *ConvertChatCompletionMessage(message))
	}
	return &converted
}

// CompletionChoice: Chat-GPeasy to OpenAI
func ConvertCompletionChoice(choice interfaces.CompletionChoice) *openai.ChatCompletionChoice {
	converted := openai.ChatCompletionChoice{
		Index:        choice.Index,
		Message:      *ConvertCompletionMessage(*choice.Message),
		FinishReason: choice.FinishReason,
	}
	return &converted
}

func ConvertCompletionChoices(choices []interfaces.CompletionChoice) *[]openai.ChatCompletionChoice {
	converted := make([]openai.ChatCompletionChoice, 0)
	for _, choice := range choices {
		converted = append(converted, *ConvertCompletionChoice(choice))
	}
	return &converted
}

// CompletionChoice: OpenAI to Chat-GPeasy
func ConvertChatCompletionChoice(choice openai.ChatCompletionChoice) *interfaces.CompletionChoice {
	cc := interfaces.CompletionChoice{
		Index:        choice.Index,
		Message:      ConvertChatCompletionMessage(choice.Message),
		FinishReason: choice.FinishReason,
	}
	return &cc
}

func ConvertChatCompletionChoices(choices []openai.ChatCompletionChoice) *[]interfaces.CompletionChoice {
	converted := make([]interfaces.CompletionChoice, 0)
	for _, choice := range choices {
		converted = append(converted, *ConvertChatCompletionChoice(choice))
	}
	return &converted
}
