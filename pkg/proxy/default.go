// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package proxy

import (
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	openai "github.com/sashabaranov/go-openai"
	klog "k8s.io/klog/v2"
)

type DefaultChatGPTCallback struct {
	AllDisable bool
}

func NewDefaultChatGPTCallback() *DefaultChatGPTCallback {
	var allDisableStr string
	if v := os.Getenv("SYMBL_ALL_DISABLE"); v != "" {
		klog.V(4).Info("SYMBL_ALL_DISABLE found")
		allDisableStr = v
	}

	allDisable := strings.EqualFold(strings.ToLower(allDisableStr), "true")

	return &DefaultChatGPTCallback{
		AllDisable: allDisable,
	}
}

func (dcc *DefaultChatGPTCallback) CreateTranscription(request openai.AudioRequest, response openai.AudioResponse) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("CreateTranscription:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) CreateTranslation(request openai.AudioRequest, response openai.AudioResponse) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("CreateTranslation:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) CreateCompletion(request openai.CompletionRequest, response openai.CompletionResponse) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("CreateCompletion:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) CreateChatCompletion(request openai.ChatCompletionRequest, response openai.ChatCompletionResponse) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("CreateChatCompletion:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) Edits(request openai.EditsRequest, response openai.EditsResponse) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("Edits:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) CreateEmbeddings(request openai.EmbeddingRequest, response openai.EmbeddingResponse) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("CreateEmbeddings:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) ListFiles(list openai.FilesList) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("ListFiles:\n\n")
	klog.Infof("Response:\n%s\n", spew.Sdump(list))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) CreateFile(request openai.FileRequest, response openai.File) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("CreateFile:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) DeleteFile(ID string) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("DeleteFile:\n\n")
	klog.Infof("Request:\nID = %s\n", ID)
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) GetFile(ID string, response openai.File) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("GetFile:\n\n")
	klog.Infof("Request:\nID = %s\n", ID)
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

// GetFileContent: Not implemented in Go SDK

func (dcc *DefaultChatGPTCallback) CreateFineTune(request openai.FineTuneRequest, response openai.FineTune) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("CreateFineTune:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) ListFineTunes(list openai.FineTuneList) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("ListFineTunes:\n\n")
	klog.Infof("Response:\n%s\n", spew.Sdump(list))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) GetFineTune(ID string, response openai.FineTune) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("GetFineTune:\n\n")
	klog.Infof("Request:\nID = %s\n", ID)
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) CancelFineTune(ID string, response openai.FineTune) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("CancelFineTune:\n\n")
	klog.Infof("Request:\nID = %s\n", ID)
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) ListFineTuneEvents(ID string, response openai.FineTuneEventList) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("ListFineTuneEvents:\n\n")
	klog.Infof("Request:\nID = %s\n", ID)
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) DeleteFineTune(ID string) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("DeleteFineTune:\n\n")
	klog.Infof("Request:\nID = %s\n", ID)
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) CreateImage(request openai.ImageRequest, response openai.ImageResponse) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("CreateImage:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) CreateEditImage(request openai.ImageEditRequest, response openai.ImageResponse) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("CreateEditImage:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) CreateVariImage(request openai.ImageVariRequest, response openai.ImageResponse) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("CreateVariImage:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) ListModels(list openai.ModelsList) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("ListModels:\n\n")
	klog.Infof("Response:\n%s\n", spew.Sdump(list))
	klog.Infof("-------------------------------\n\n")
	return nil
}

func (dcc *DefaultChatGPTCallback) Moderations(request openai.ModerationRequest, response openai.ModerationResponse) error {
	klog.Infof("\n\n-------------------------------\n")
	klog.Infof("Moderations:\n\n")
	klog.Infof("Request:\n%s\n\n", spew.Sdump(request))
	klog.Infof("Response:\n%s\n", spew.Sdump(response))
	klog.Infof("-------------------------------\n\n")
	return nil
}
