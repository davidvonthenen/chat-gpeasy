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
			Content: "You are a helpful assistant. If you dont know the answer or some of what you might give is not factual, please say I dont know or omit that part of your reply.",
		})
	case interfaces.SkillTypeExpert:
		p.conversation = make([]openai.ChatCompletionMessage, 0)
		p.conversation = append(p.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are an expert in your field. If you dont know the answer or some of what you might give is not factual, please say I dont know or omit that part of your reply.",
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

func (p *Persona) EditConversation(index int, statement string) ([]openai.ChatCompletionChoice, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	numOfConvo := len(p.conversation)

	if index < 0 || index >= numOfConvo {
		klog.V(1).Infof("invalid index (%d) must be between 0 and %d\n", numOfConvo, (numOfConvo - 1))
		return nil, interfaces.ErrInvalidInput
	}
	if len(statement) == 0 {
		klog.V(1).Infof("statement is empty\n")
		return nil, interfaces.ErrInvalidInput
	}

	klog.V(6).Infof("advanced.EditConversation ENTER\n")
	klog.V(5).Infof("index: %d\n", index)
	klog.V(5).Infof("statement: %s\n", statement)

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
	response, err := p.client.CreateChatCompletion(ctx, request)
	if err != nil {
		klog.V(1).Infof("CreateChatCompletion error: %v\n", err)
		klog.V(6).Infof("advanced.EditConversation LEAVE\n")
		return nil, err
	}

	// housekeeping
	p.appendedResponse = false
	p.conversation = convo
	p.request = &request
	p.response = &response

	if len(response.Choices) == 1 {
		p.appendedResponse = true
		p.conversation = append(p.conversation, response.Choices[0].Message)
	}

	for _, msg := range p.conversation {
		klog.V(5).Infof("Output message (type: %s): %s\n", msg.Role, msg.Content)
	}

	klog.V(4).Infof("advanced.Query Succeeded\n")
	klog.V(6).Infof("advanced.Query LEAVE\n")

	return response.Choices, nil
}

func (p *Persona) Query(ctx context.Context, role, statement string) ([]openai.ChatCompletionChoice, error) {
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

	switch role {
	case openai.ChatMessageRoleUser:
	case openai.ChatMessageRoleSystem:
	case openai.ChatMessageRoleAssistant:
	default:
		klog.V(1).Infof("Invalid role type = %s\n", role)
		return nil, interfaces.ErrInvalidInput
	}

	klog.V(6).Infof("advanced.Query ENTER\n")
	klog.V(5).Infof("role: %s\n", role)
	klog.V(5).Infof("statement: %s\n", statement)

	convo := make([]openai.ChatCompletionMessage, 0)
	for _, msg := range p.conversation {
		klog.V(5).Infof("Input message (type: %s): %s\n", msg.Role, msg.Content)
		convo = append(convo, msg)
	}
	klog.V(5).Infof("Input NEW message (type: %s): %s\n", role, statement)
	convo = append(convo, openai.ChatCompletionMessage{
		Role:    role,
		Content: statement,
	})

	request := openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: convo,
	}

	response, err := p.client.CreateChatCompletion(ctx, request)
	if err != nil {
		klog.V(1).Infof("CreateChatCompletion error: %v\n", err)
		klog.V(6).Infof("advanced.Query LEAVE\n")
		return nil, err
	}

	// housekeeping
	p.appendedResponse = false
	p.conversation = convo
	p.request = &request
	p.response = &response

	if len(response.Choices) == 1 {
		p.appendedResponse = true
		p.conversation = append(p.conversation, response.Choices[0].Message)
	}

	for _, msg := range p.conversation {
		klog.V(5).Infof("Output message (type: %s): %s\n", msg.Role, msg.Content)
	}

	klog.V(4).Infof("advanced.Query Succeeded\n")
	klog.V(6).Infof("advanced.Query LEAVE\n")

	return response.Choices, nil
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

func (p *Persona) CommitResponse(index int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	klog.V(6).Infof("advanced.CommitResponse ENTER\n")
	klog.V(5).Infof("index: %d\n", index)

	if p.appendedResponse {
		klog.V(1).Infof("already appended response\n")
		klog.V(6).Infof("advanced.CommitResponse LEAVE\n")
		return interfaces.ErrInvalidInput
	}
	if p.response == nil {
		klog.V(1).Infof("response is empty\n")
		klog.V(6).Infof("advanced.CommitResponse LEAVE\n")
		return interfaces.ErrInvalidInput
	}
	if len(p.response.Choices) == 0 {
		klog.V(1).Infof("not choices generated\n")
		klog.V(6).Infof("advanced.CommitResponse LEAVE\n")
		return interfaces.ErrInvalidInput
	}

	for _, msg := range p.response.Choices {
		if msg.Index == index {
			p.appendedResponse = true
			p.conversation = append(p.conversation, msg.Message)
			break
		}
	}

	if p.appendedResponse {
		klog.V(1).Infof("response %d not found\n", index)
		klog.V(6).Infof("advanced.CommitResponse LEAVE\n")
		return interfaces.ErrInvalidInput
	}

	for _, msg := range p.conversation {
		klog.V(5).Infof("Message (type: %s): %s\n", msg.Role, msg.Content)
	}

	klog.V(4).Infof("response %d found\n", index)
	klog.V(6).Infof("advanced.CommitResponse LEAVE\n")
	return nil
}
