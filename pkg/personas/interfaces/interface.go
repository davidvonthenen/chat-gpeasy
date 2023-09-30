// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

import (
	"context"
	"io"
)

// shared
type CompletionMessage struct {
	Role    string
	Content string
	Name    string
}
type CompletionChoice struct {
	Index   int
	Message *CompletionMessage
	// FinishReason string
}

// rest interfaces
type SimpleChat interface {
	Init(level SkillType, model string) error
	Query(ctx context.Context, statement string) (string, error)
}

type StandardChat interface {
	Init(level SkillType, model string) error
	Query(ctx context.Context, statement string) ([]CompletionChoice, error)
	AddDirective(directives string) error
	CommitResponse(index int) error
}

type AdvancedChat interface {
	Init(level SkillType, model string) error
	InitWithProvided(model string, previous []CompletionMessage) error
	DynamicInit(model string, previous []CompletionMessage) error
	GetConversation() ([]CompletionMessage, error)
	EditConversation(index int, statement string) ([]CompletionChoice, error)
	Query(ctx context.Context, role, statement string) ([]CompletionChoice, error)
	AddDirective(directives string) error
	AddUserContext(text string) error
	CommitResponse(index int) error
}

// streaming interfaces
type StreamingCompletion interface {
	Stream(w io.Writer) error
	Close() error
}

type StandardChatStream interface {
	Init(level SkillType, model string) error
	GetConversation() ([]CompletionMessage, error)
	Query(ctx context.Context, statement string) (*StreamingCompletion, error)
	AddDirective(directives string) error
}

type AdvancedChatStream interface {
	Init(level SkillType, model string) error
	InitWithProvided(model string, previous []CompletionMessage) error
	DynamicInit(model string, previous []CompletionMessage) error
	GetConversation() ([]CompletionMessage, error)
	EditConversation(index int, statement string) (*StreamingCompletion, error)
	Query(ctx context.Context, statement string) (*StreamingCompletion, error)
	AddDirective(directives string) error
	AddUserContext(text string) error
}
