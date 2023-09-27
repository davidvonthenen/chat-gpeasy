// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package personas

import (
	"crypto/tls"
	"net/http"
	"os"

	openai "github.com/sashabaranov/go-openai"
	klog "k8s.io/klog/v2"
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
