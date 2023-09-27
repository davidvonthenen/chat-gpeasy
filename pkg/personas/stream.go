// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package personas

import (
	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
	advanced "github.com/dvonthenen/chat-gpeasy/pkg/personas/stream/advanced"
)

func NewAdvancedChatStream() (*interfaces.AdvancedChatStream, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}

	advancedChat, err := advanced.New(client)
	if err != nil {
		return nil, err
	}

	var advanced interfaces.AdvancedChatStream
	advanced = advancedChat

	return &advanced, nil
}

func NewAdvancedChatStreamWithOptions(opt *PersonaOptions) (*interfaces.AdvancedChatStream, error) {
	client, err := newWithOptions(opt)
	if err != nil {
		return nil, err
	}

	advancedChat, err := advanced.New(client)
	if err != nil {
		return nil, err
	}

	var advanced interfaces.AdvancedChatStream
	advanced = advancedChat

	return &advanced, nil
}
