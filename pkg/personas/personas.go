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

const (
	openaiAPIURLv1 string = "https://api.openai.com/v1"
)

func DefaultConfig(OpenAiUrl, OpenAiApiKey string) (*PersonaOptions, error) {
	var openAiApiKey string
	if len(OpenAiApiKey) > 0 {
		openAiApiKey = OpenAiApiKey
	} else {
		if v := os.Getenv("OPENAI_API_KEY"); v != "" {
			openAiApiKey = v
		} else {
			klog.V(1).Infof("OPENAI_API_KEY not found\n")
			return nil, ErrNotFoundOpenAiKey
		}
	}

	config := openai.DefaultConfig(openAiApiKey)
	config.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if len(OpenAiUrl) > 0 {
		config.BaseURL = OpenAiUrl
	}

	opt := &PersonaOptions{
		&config,
		true,
	}

	return opt, nil
}

func newClient() (*openai.Client, error) {
	options, err := DefaultConfig("https://127.0.0.1/v1", "")
	if err != nil {
		return nil, err
	}

	client := openai.NewClientWithConfig(*options.ClientConfig)

	return client, nil
}

func newWithOptions(options *PersonaOptions) (*openai.Client, error) {
	// if options.BindPort == 0 {
	// 	options.BindPort = DefaultPort
	// }
	if options.DisableHostVerify {
		options.ClientConfig.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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

func NewCumulativeChatWithOptions(opt *PersonaOptions) (*interfaces.CumulativeChat, error) {
	client, err := newWithOptions(opt)
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
