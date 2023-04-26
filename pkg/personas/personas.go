// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package personas

import (
	"crypto/tls"
	"net/http"
	"os"

	openai "github.com/sashabaranov/go-openai"
	klog "k8s.io/klog/v2"

	advanced "github.com/dvonthenen/chat-gpeasy/pkg/personas/advanced"
	cumulative "github.com/dvonthenen/chat-gpeasy/pkg/personas/cumulative"
	interfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"
	simple "github.com/dvonthenen/chat-gpeasy/pkg/personas/simple"
)

func DefaultConfig() (*PersonaOptions, error) {
	var openAiApiKey string
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		openAiApiKey = v
	} else {
		klog.V(1).Infof("OPENAI_API_KEY not found\n")
		return nil, ErrNotFoundOpenAiKey
	}

	config := openai.DefaultConfig(openAiApiKey)
	config.BaseURL = "https://127.0.0.1/v1"
	config.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	opt := &PersonaOptions{
		&config,
		true,
	}

	return opt, nil
}

func newClient() (*openai.Client, error) {
	options, err := DefaultConfig()
	if err != nil {
		return nil, err
	}

	client := openai.NewClientWithConfig(*options.ClientConfig)

	return client, nil
}

func newWithOptions(options PersonaOptions, openAiApiKey string) (*openai.Client, error) {
	// if options.BindPort == 0 {
	// 	options.BindPort = DefaultPort
	// }
	if options.DisableHostVerify {
		options.ClientConfig.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	if len(openAiApiKey) == 0 {
		if v := os.Getenv("OPENAI_API_KEY"); v != "" {
			klog.V(4).Info("OPENAI_API_KEY found")
			openAiApiKey = v
		} else {
			klog.V(1).Infof("OPENAI_API_KEY not found\n")
			return nil, ErrNotFoundOpenAiKey
		}
	}

	client := openai.NewClientWithConfig(*options.ClientConfig)
	if client == nil {
		klog.V(1).Infof("NewClientWithConfig is nil\n")
		return nil, ErrInvalidOpenAiClient
	}

	return client, nil
}

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

func NewSimpleChatWithOptions(opt PersonaOptions, openAiApiKey string) (*interfaces.SimpleChat, error) {
	client, err := newWithOptions(opt, openAiApiKey)
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

func NewCumulativeChat() (*interfaces.CumulativeChat, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}

	cumulativeChat, err := cumulative.New(client)
	if err != nil {
		return nil, err
	}

	var cumulative interfaces.CumulativeChat
	cumulative = cumulativeChat

	return &cumulative, nil
}

func NewCumulativeChatWithOptions(opt PersonaOptions, openAiApiKey string) (*interfaces.CumulativeChat, error) {
	client, err := newWithOptions(opt, openAiApiKey)
	if err != nil {
		return nil, err
	}

	cumulativeChat, err := cumulative.New(client)
	if err != nil {
		return nil, err
	}

	var cumulative interfaces.CumulativeChat
	cumulative = cumulativeChat

	return &cumulative, nil
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

func NewAdvancedChatWithOptions(opt PersonaOptions, openAiApiKey string) (*interfaces.AdvancedChat, error) {
	client, err := newWithOptions(opt, openAiApiKey)
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
