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

func (p *Persona) InitWithProvided(model string, previous []openai.ChatCompletionMessage) error {
	if len(model) == 0 {
		model = openai.GPT3Dot5Turbo
	}
	p.model = model
	p.level = interfaces.SkillTypeCustom

	p.conversation = make([]openai.ChatCompletionMessage, 0)
	copy(p.conversation, previous)

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

func (p *Persona) EditConversation(index int, statement string) (*interfaces.StreamingCompletion, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	klog.V(6).Infof("advanced.EditConversation ENTER\n")
	klog.V(5).Infof("index: %d\n", index)
	klog.V(5).Infof("statement: %s\n", statement)

	numOfConvo := len(p.conversation)

	if index < 0 || index >= numOfConvo {
		klog.V(1).Infof("invalid index (%d) must be between 0 and %d\n", numOfConvo, (numOfConvo - 1))
		return nil, interfaces.ErrInvalidInput
	}
	if len(statement) == 0 {
		klog.V(1).Infof("statement is empty\n")
		return nil, interfaces.ErrInvalidInput
	}

	convo := make([]openai.ChatCompletionMessage, 0)

	for pos, msg := range p.conversation {
		if pos == (numOfConvo-1) && msg.Role == openai.ChatMessageRoleAssistant {
			klog.V(3).Infof("skip last chatgpt response for a redo/rebuild\n")
			break
		}
		if pos == index {
			if msg.Role == openai.ChatMessageRoleAssistant {
				klog.V(1).Infof("unable to edit the response to queries/prompts\n")
				klog.V(6).Infof("advanced.EditConversation LEAVE\n")
				return nil, interfaces.ErrInvalidInput
			}
			klog.V(3).Infof("placing:\n%s\n\nwith:\n%s\n", msg.Content, statement)
			convo = append(convo, openai.ChatCompletionMessage{
				Role:    msg.Role,
				Content: statement,
			})
		} else {
			convo = append(convo, msg)
		}
	}

	for _, msg := range convo {
		klog.V(5).Infof("Updated convo (type: %s): %s\n", msg.Role, msg.Content)
	}

	// requery
	request := openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: convo,
	}

	ctx := context.Background()
	stream, err := p.client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		klog.V(1).Infof("CreateChatCompletionStream error: %v\n", err)
		klog.V(6).Infof("advanced.EditConversation LEAVE\n")
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

	klog.V(4).Infof("advanced.Query Succeeded\n")
	klog.V(6).Infof("advanced.Query LEAVE\n")

	return &streamingcompletion, nil
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
		klog.V(1).Infof("CreateChatCompletionStream failed. Err: %v\n", err)
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
