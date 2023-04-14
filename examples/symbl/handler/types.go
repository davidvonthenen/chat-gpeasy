// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package handler

import (
	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
)

/*
	Handler for messages
*/
type HandlerOptions struct {
	Simple *interfaces.SimpleChat
}

type Handler struct {
	// properties
	conversationID string

	// housekeeping
	simple *interfaces.SimpleChat
}
