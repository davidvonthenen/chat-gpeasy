// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package handler

import (
	"context"
	"fmt"

	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
)

func NewHandler(options HandlerOptions) *Handler {
	handler := Handler{
		simple: options.Simple,
	}
	return &handler
}

func (h *Handler) InitializedConversation(im *sdkinterfaces.InitializationMessage) error {
	h.conversationID = im.Message.Data.ConversationID
	fmt.Printf("conversationID: %s\n", h.conversationID)
	return nil
}

func (h *Handler) RecognitionResultMessage(rr *sdkinterfaces.RecognitionResult) error {
	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) MessageResponseMessage(mr *sdkinterfaces.MessageResponse) error {
	for _, msg := range mr.Messages {
		fmt.Printf("\n\nMessage [%s]: %s\n\n", msg.From.Name, msg.Payload.Content)
	}
	return nil
}

func (h *Handler) InsightResponseMessage(ir *sdkinterfaces.InsightResponse) error {
	for _, insight := range ir.Insights {
		switch insight.Type {
		case sdkinterfaces.InsightTypeQuestion:
			err := h.HandleQuestion(&insight, ir.SequenceNumber)
			if err != nil {
				fmt.Printf("HandleQuestion failed. Err: %v\n", err)
				return err
			}
		case sdkinterfaces.InsightTypeFollowUp:
			err := h.HandleFollowUp(&insight, ir.SequenceNumber)
			if err != nil {
				fmt.Printf("HandleFollowUp failed. Err: %v\n", err)
				return err
			}
		case sdkinterfaces.InsightTypeActionItem:
			err := h.HandleActionItem(&insight, ir.SequenceNumber)
			if err != nil {
				fmt.Printf("HandleActionItem failed. Err: %v\n", err)
				return err
			}
		default:
			fmt.Printf("\n\n-------------------------------\n")
			fmt.Printf("Unknown InsightResponseMessage: %s\n\n", insight.Type)
			fmt.Printf("-------------------------------\n\n")
			return nil
		}
	}

	return nil
}

func (h *Handler) TopicResponseMessage(tr *sdkinterfaces.TopicResponse) error {
	for _, curTopic := range tr.Topics {
		prompt := fmt.Sprintf("The topic of \"%s\" came up in this conversation I am having. Come up with 2 general facts and 2 obscure facts about this topic.\n", curTopic.Phrases)

		ctx := context.Background()
		choices, err := (*h.simple).Query(ctx, prompt)
		if err != nil {
			fmt.Printf("persona.Query error: %v\n", err)
			return err
		}
		fmt.Printf("\n\n-------------------------------\n")
		fmt.Printf("TOPIC:\n%s\n", prompt)
		fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
		fmt.Printf("-------------------------------\n\n")

	}

	return nil
}

func (h *Handler) TrackerResponseMessage(tr *sdkinterfaces.TrackerResponse) error {
	for _, curTracker := range tr.Trackers {
		for _, match := range curTracker.Matches {
			prompt := fmt.Sprintf("The topic of \"%s\" came up in this conversation I am having. Come up with 2 general facts and 2 obscure facts about this topic.\n", match.Value)

			ctx := context.Background()
			choices, err := (*h.simple).Query(ctx, prompt)
			if err != nil {
				fmt.Printf("persona.Query error: %v\n", err)
				return err
			}
			fmt.Printf("\n\n-------------------------------\n")
			fmt.Printf("TRACKER:\n%s\n", prompt)
			fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
			fmt.Printf("-------------------------------\n\n")
		}
	}

	return nil
}

func (h *Handler) EntityResponseMessage(er *sdkinterfaces.EntityResponse) error {
	for _, entity := range er.Entities {
		for _, match := range entity.Matches {
			prompt := fmt.Sprintf("Someone mentioned \"%s\" in this conversation I am having. Come up with 2 general facts and 2 obscure facts about \"%s\".\n", match.DetectedValue, match.DetectedValue)

			ctx := context.Background()
			choices, err := (*h.simple).Query(ctx, prompt)
			if err != nil {
				fmt.Printf("persona.Query error: %v\n", err)
				return err
			}
			fmt.Printf("\n\n-------------------------------\n")
			fmt.Printf("ENTITY:\n%s\n", prompt)
			fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
			fmt.Printf("-------------------------------\n\n")
		}
	}

	return nil
}

func (h *Handler) TeardownConversation(tm *sdkinterfaces.TeardownMessage) error {
	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) UserDefinedMessage(data []byte) error {
	// This is only needed on the client side and not on the plugin side.
	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) UnhandledMessage(byMsg []byte) error {
	fmt.Printf("\n\n-------------------------------\n")
	fmt.Printf("UnhandledMessage:\n%v\n", string(byMsg))
	fmt.Printf("-------------------------------\n\n")
	return ErrUnhandledMessage
}

func (h *Handler) HandleQuestion(insight *sdkinterfaces.Insight, number int) error {
	prompt := fmt.Sprintf("Someone is asking the following question below:\n\n\"%s\"\n\nCome up with 3 different concise answers that someone with more than a passing familiarity in the subject might reply with. Provide a small detail as proof of you being a familiar in that subject.\n", insight.Payload.Content)

	ctx := context.Background()
	choices, err := (*h.simple).Query(ctx, prompt)
	if err != nil {
		fmt.Printf("persona.Query error: %v\n", err)
		return err
	}
	fmt.Printf("\n\n-------------------------------\n")
	fmt.Printf("QUESTION:\n%s\n", prompt)
	fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
	fmt.Printf("-------------------------------\n\n")

	return nil
}

func (h *Handler) HandleActionItem(insight *sdkinterfaces.Insight, number int) error {
	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) HandleFollowUp(insight *sdkinterfaces.Insight, number int) error {
	// No implementation required. Return Succeess!
	return nil
}
