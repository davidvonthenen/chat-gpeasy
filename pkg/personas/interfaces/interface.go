// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type SimpleChat interface {
	Init(level SkillType, model string) error
	Query(ctx context.Context, statement string) ([]openai.ChatCompletionChoice, error)
	CommitResponse(index int) error
}

type CumulativeChat interface {
	Init(level SkillType, model string) error
	Query(ctx context.Context, statement string) ([]openai.ChatCompletionChoice, error)
	AddDirective(directives string) error
	CommitResponse(index int) error
}

type AdvancedChat interface {
	Init(level SkillType, model string) error
	GetConversation() ([]openai.ChatCompletionMessage, error)
	EditConversation(index int, statement string) ([]openai.ChatCompletionChoice, error)
	Query(ctx context.Context, role, statement string) ([]openai.ChatCompletionChoice, error)
	AddDirective(directives string) error
	CommitResponse(index int) error
}
