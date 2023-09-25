// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package simple

import (
	"context"

	openai "github.com/sashabaranov/go-openai"

	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
	advanced "github.com/dvonthenen/chat-gpeasy/pkg/personas/rest/advanced"
)

func New(client *openai.Client) (*Persona, error) {
	if client == nil {
		return nil, interfaces.ErrInvalidInput
	}
	advanced, err := advanced.New(client)
	if err != nil {
		return nil, err
	}
	return &Persona{
		persona: advanced,
	}, nil
}

func (c *Persona) Init(level interfaces.SkillType, model string) error {
	return c.persona.Init(interfaces.SkillType(level), model)
}

func (c *Persona) Query(ctx context.Context, statement string) (string, error) {
	choices, err := c.persona.Query(ctx, openai.ChatMessageRoleUser, statement)
	if err != nil {
		return "", err
	}
	if len(choices) == 0 {
		return "", interfaces.ErrEmptyChoices
	}

	if len(choices) > 1 {
		err = c.persona.CommitResponse(0)
		if err != nil {
			return "", err
		}
	}

	return choices[0].Message.Content, nil
}
