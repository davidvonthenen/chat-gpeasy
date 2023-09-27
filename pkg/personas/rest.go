// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package personas

import (
	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
	advanced "github.com/dvonthenen/chat-gpeasy/pkg/personas/rest/advanced"
	simple "github.com/dvonthenen/chat-gpeasy/pkg/personas/rest/simple"
	standard "github.com/dvonthenen/chat-gpeasy/pkg/personas/rest/standard"
)

func NewSimpleChat() (*interfaces.SimpleChat, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}

	simpleChat, err := simple.New(client)
	if err != nil {
		return nil, err
	}

	var simple interfaces.SimpleChat
	simple = simpleChat

	return &simple, nil
}

func NewSimpleChatWithOptions(opt *PersonaOptions) (*interfaces.SimpleChat, error) {
	client, err := newWithOptions(opt)
	if err != nil {
		return nil, err
	}

	simpleChat, err := simple.New(client)
	if err != nil {
		return nil, err
	}

	var simple interfaces.SimpleChat
	simple = simpleChat

	return &simple, nil
}

func NewStandardChat() (*interfaces.StandardChat, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}

	standardChat, err := standard.New(client)
	if err != nil {
		return nil, err
	}

	var standard interfaces.StandardChat
	standard = standardChat

	return &standard, nil
}

func NewStandardChatWithOptions(opt *PersonaOptions) (*interfaces.StandardChat, error) {
	client, err := newWithOptions(opt)
	if err != nil {
		return nil, err
	}

	standardChat, err := standard.New(client)
	if err != nil {
		return nil, err
	}

	var standard interfaces.StandardChat
	standard = standardChat

	return &standard, nil
}

func NewAdvancedChat() (*interfaces.AdvancedChat, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}

	advancedChat, err := advanced.New(client)
	if err != nil {
		return nil, err
	}

	var advanced interfaces.AdvancedChat
	advanced = advancedChat

	return &advanced, nil
}

func NewAdvancedChatWithOptions(opt *PersonaOptions) (*interfaces.AdvancedChat, error) {
	client, err := newWithOptions(opt)
	if err != nil {
		return nil, err
	}

	advancedChat, err := advanced.New(client)
	if err != nil {
		return nil, err
	}

	var advanced interfaces.AdvancedChat
	advanced = advancedChat

	return &advanced, nil
}
