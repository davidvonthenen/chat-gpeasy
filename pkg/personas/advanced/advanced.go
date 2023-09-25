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

func (c *Persona) Init(level interfaces.SkillType, model string) error {
	if len(model) == 0 {
		model = openai.GPT3Dot5Turbo
	}
	c.model = model
	c.level = level

	switch c.level {
	case interfaces.SkillTypeGeneric:
		c.conversation = make([]openai.ChatCompletionMessage, 0)
		c.conversation = append(c.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant. If you dont know the answer or some of what you might give is not factual, please say I dont know or omit that part of your reply.",
		})
		c.appendedResponse = true
	case interfaces.SkillTypeExpert:
		c.conversation = make([]openai.ChatCompletionMessage, 0)
		c.conversation = append(c.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are an expert in your field. If you dont know the answer or some of what you might give is not factual, please say I dont know or omit that part of your reply.",
		})
		c.appendedResponse = true
	case interfaces.SkillTypeDAN:
		c.conversation = make([]openai.ChatCompletionMessage, 0)
		c.conversation = append(c.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: DANPrompt,
		})
		c.appendedResponse = true
	case interfaces.SkillTypeSTAN:
		c.conversation = make([]openai.ChatCompletionMessage, 0)
		c.conversation = append(c.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: STANPrompt,
		})
		c.appendedResponse = true
	case interfaces.SkillTypeDUDE:
		c.conversation = make([]openai.ChatCompletionMessage, 0)
		c.conversation = append(c.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: DUDEPrompt,
		})
		c.appendedResponse = true
	case interfaces.SkillTypeJailBreak:
		c.conversation = make([]openai.ChatCompletionMessage, 0)
		c.conversation = append(c.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: JailBreakPrompt,
		})
		c.appendedResponse = true
	case interfaces.SkillTypeMongo:
		c.conversation = make([]openai.ChatCompletionMessage, 0)
		c.conversation = append(c.conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: MongoPrompt,
		})
		c.appendedResponse = true
	}

	return nil
}

func (c *Persona) GetConversation() ([]openai.ChatCompletionMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conversation == nil {
		return nil, interfaces.ErrInvalidInput
	}

	return c.conversation, nil
}

func (c *Persona) EditConversation(index int, statement string) ([]openai.ChatCompletionChoice, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	numOfConvo := len(c.conversation)

	if index < 0 || index >= numOfConvo {
		klog.V(1).Infof("invalid index (%d) must be between 0 and %d\n", numOfConvo, (numOfConvo - 1))
		return nil, interfaces.ErrInvalidInput
	}
	if len(statement) == 0 {
		klog.V(1).Infof("statement is empty\n")
		return nil, interfaces.ErrInvalidInput
	}

	klog.V(6).Infof("cumulative.EditConversation ENTER\n")
	klog.V(5).Infof("index: %d\n", index)
	klog.V(5).Infof("statement: %s\n", statement)

	convo := make([]openai.ChatCompletionMessage, 0)

	for pos, msg := range c.conversation {
		if pos == (numOfConvo-1) && msg.Role == openai.ChatMessageRoleAssistant {
			klog.V(3).Infof("skip last chatgpt response for a redo/rebuild\n")
			break
		}
		if pos == index {
			if msg.Role == openai.ChatMessageRoleAssistant {
				klog.V(1).Infof("unable to edit the response to queries/prompts\n")
				klog.V(6).Infof("cumulative.EditConversation LEAVE\n")
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
		Model:    c.model,
		Messages: convo,
	}

	ctx := context.Background()
	response, err := c.client.CreateChatCompletion(ctx, request)
	if err != nil {
		klog.V(1).Infof("CreateChatCompletion error: %v\n", err)
		klog.V(6).Infof("cumulative.EditConversation LEAVE\n")
		return nil, err
	}

	// housekeeping
	c.appendedResponse = false
	c.conversation = convo
	c.request = &request
	c.response = &response

	if len(response.Choices) == 1 {
		c.appendedResponse = true
		c.conversation = append(c.conversation, response.Choices[0].Message)
	}

	for _, msg := range c.conversation {
		klog.V(5).Infof("Output message (type: %s): %s\n", msg.Role, msg.Content)
	}

	klog.V(4).Infof("cumulative.Query Succeeded\n")
	klog.V(6).Infof("cumulative.Query LEAVE\n")

	return response.Choices, nil
}

func (c *Persona) Query(ctx context.Context, role, statement string) ([]openai.ChatCompletionChoice, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.appendedResponse {
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

	klog.V(6).Infof("cumulative.Query ENTER\n")
	klog.V(5).Infof("role: %s\n", role)
	klog.V(5).Infof("statement: %s\n", statement)

	convo := make([]openai.ChatCompletionMessage, 0)
	for _, msg := range c.conversation {
		klog.V(5).Infof("Input message (type: %s): %s\n", msg.Role, msg.Content)
		convo = append(convo, msg)
	}
	klog.V(5).Infof("Input NEW message (type: %s): %s\n", role, statement)
	convo = append(convo, openai.ChatCompletionMessage{
		Role:    role,
		Content: statement,
	})

	request := openai.ChatCompletionRequest{
		Model:    c.model,
		Messages: convo,
	}

	response, err := c.client.CreateChatCompletion(ctx, request)
	if err != nil {
		klog.V(1).Infof("CreateChatCompletion error: %v\n", err)
		klog.V(6).Infof("cumulative.Query LEAVE\n")
		return nil, err
	}

	// housekeeping
	c.appendedResponse = false
	c.conversation = convo
	c.request = &request
	c.response = &response

	if len(response.Choices) == 1 {
		c.appendedResponse = true
		c.conversation = append(c.conversation, response.Choices[0].Message)
	}

	for _, msg := range c.conversation {
		klog.V(5).Infof("Output message (type: %s): %s\n", msg.Role, msg.Content)
	}

	klog.V(4).Infof("cumulative.Query Succeeded\n")
	klog.V(6).Infof("cumulative.Query LEAVE\n")

	return response.Choices, nil
}

func (c *Persona) AddDirective(directives string) error {
	klog.V(5).Infof("AddDirectives: %s\n", directives)

	if len(directives) == 0 {
		klog.V(1).Infof("directives is empty\n")
		return interfaces.ErrInvalidInput
	}

	c.conversation = append(c.conversation, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: directives,
	})

	return nil
}

func (c *Persona) CommitResponse(index int) error {
	klog.V(6).Infof("cumulative.CommitResponse ENTER\n")
	klog.V(5).Infof("index: %d\n", index)

	if c.appendedResponse {
		klog.V(1).Infof("already appended response\n")
		klog.V(6).Infof("cumulative.CommitResponse LEAVE\n")
		return interfaces.ErrInvalidInput
	}
	if c.response == nil {
		klog.V(1).Infof("response is empty\n")
		klog.V(6).Infof("cumulative.CommitResponse LEAVE\n")
		return interfaces.ErrInvalidInput
	}
	if len(c.response.Choices) == 0 {
		klog.V(1).Infof("not choices generated\n")
		klog.V(6).Infof("cumulative.CommitResponse LEAVE\n")
		return interfaces.ErrInvalidInput
	}

	for _, msg := range c.response.Choices {
		if msg.Index == index {
			c.appendedResponse = true
			c.conversation = append(c.conversation, msg.Message)
			break
		}
	}

	if c.appendedResponse {
		klog.V(1).Infof("response %d not found\n", index)
		klog.V(6).Infof("cumulative.CommitResponse LEAVE\n")
		return interfaces.ErrInvalidInput
	}

	for _, msg := range c.conversation {
		klog.V(5).Infof("Message (type: %s): %s\n", msg.Role, msg.Content)
	}

	klog.V(4).Infof("response %d found\n", index)
	klog.V(6).Infof("cumulative.CommitResponse LEAVE\n")
	return nil
}
