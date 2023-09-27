// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package advanced

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
)

func New(client *openai.Client) (*Persona, error) {
	if client == nil {
		return nil, interfaces.ErrInvalidInput
	}
	return &Persona{
		client:       client,
		conversation: make([]openai.ChatCompletionMessage, 0),
	}, nil
}

func (p *Persona) Init(level interfaces.SkillType, model string) error {
	if p.appendedResponse {
		klog.V(1).Infof("Init has already been called\n")
		return interfaces.ErrInitAlready
	}

	if len(model) == 0 {
		model = openai.GPT3Dot5Turbo
	}
	p.model = model
	p.level = level

	switch p.level {
	case interfaces.SkillTypeGeneric:
		p.conversation = make([]openai.ChatCompletionMessage, 0)
		p.conversation = append(p.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant.",
		})
	case interfaces.SkillTypeExpert:
		p.conversation = make([]openai.ChatCompletionMessage, 0)
		p.conversation = append(p.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are an expert in your field.",
		})
	case interfaces.SkillTypeDAN:
		p.conversation = make([]openai.ChatCompletionMessage, 0)
		p.conversation = append(p.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: interfaces.DANPrompt,
		})
	case interfaces.SkillTypeSTAN:
		p.conversation = make([]openai.ChatCompletionMessage, 0)
		p.conversation = append(p.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: interfaces.STANPrompt,
		})
	case interfaces.SkillTypeDUDE:
		p.conversation = make([]openai.ChatCompletionMessage, 0)
		p.conversation = append(p.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: interfaces.DUDEPrompt,
		})
	case interfaces.SkillTypeJailBreak:
		p.conversation = make([]openai.ChatCompletionMessage, 0)
		p.conversation = append(p.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: interfaces.JailBreakPrompt,
		})
	case interfaces.SkillTypeMongo:
		p.conversation = make([]openai.ChatCompletionMessage, 0)
		p.conversation = append(p.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: interfaces.MongoPrompt,
		})
	}

	p.appendedResponse = true

	return nil
}

func (p *Persona) InitWithProvided(previous []openai.ChatCompletionMessage) error {
	p.conversation = make([]openai.ChatCompletionMessage, 0)
	p.conversation = append(p.conversation, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: interfaces.MongoPrompt,
	})
	p.appendedResponse = true

	return nil
}

func (p *Persona) GetConversation() ([]openai.ChatCompletionMessage, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conversation == nil {
		return nil, interfaces.ErrInvalidInput
	}

	return p.conversation, nil
}

func (p *Persona) Query(ctx context.Context, statement string) (*interfaces.StreamingCompletion, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.appendedResponse {
		klog.V(1).Infof("the response hasn't been appended yet\n")
		return nil, interfaces.ErrInvalidInput
	}

	if len(statement) == 0 {
		klog.V(1).Infof("statement is empty\n")
		return nil, interfaces.ErrInvalidInput
	}

	klog.V(6).Infof("streaming.Query ENTER\n")
	klog.V(5).Infof("statement: %s\n", statement)

	convo := make([]openai.ChatCompletionMessage, 0)
	for _, msg := range p.conversation {
		klog.V(5).Infof("Input message (type: %s): %s\n", msg.Role, msg.Content)
		convo = append(convo, msg)
	}
	klog.V(5).Infof("Input NEW message: %s\n", statement)
	convo = append(convo, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: statement,
	})

	stream, err := p.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: convo,
		Stream:   true,
	})
	if err != nil {
		klog.V(1).Infof("error creating chat completion stream", err)
		return nil, err
	}

	// housekeeping
	p.appendedResponse = false
	p.conversation = convo

	var cb CommitCallback
	cb = p

	scc := StreamingChatCompletion{
		stopChan: make(chan struct{}),
		stream:   stream,
		callback: &cb,
	}

	var streamingcompletion interfaces.StreamingCompletion
	streamingcompletion = scc

	for _, msg := range p.conversation {
		klog.V(5).Infof("Output message (type: %s): %s\n", msg.Role, msg.Content)
	}

	klog.V(4).Infof("streaming.Query Succeeded\n")
	klog.V(6).Infof("streaming.Query LEAVE\n")

	return &streamingcompletion, nil
}

func (p *Persona) AddDirective(directives string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	klog.V(5).Infof("AddDirectives: %s\n", directives)

	if len(directives) == 0 {
		klog.V(1).Infof("directives is empty\n")
		return interfaces.ErrInvalidInput
	}

	p.conversation = append(p.conversation, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: directives,
	})

	return nil
}

func (p *Persona) CommitResponse(response string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	klog.V(6).Infof("advanced.CommitResponse ENTER\n")

	if p.appendedResponse {
		klog.V(1).Infof("already appended response\n")
		klog.V(6).Infof("advanced.CommitResponse LEAVE\n")
		return interfaces.ErrInvalidInput
	}

	p.appendedResponse = true
	p.conversation = append(p.conversation, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: response,
	})

	for _, msg := range p.conversation {
		klog.V(5).Infof("Message (type: %s): %s\n", msg.Role, msg.Content)
	}

	klog.V(6).Infof("advanced.CommitResponse LEAVE\n")
	return nil
}
