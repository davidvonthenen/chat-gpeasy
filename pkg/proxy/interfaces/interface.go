// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

import (
	openai "github.com/sashabaranov/go-openai"
)

type ChatGPTCallback interface {
	CreateTranscription(openai.AudioRequest, openai.AudioResponse) error
	CreateTranslation(openai.AudioRequest, openai.AudioResponse) error

	CreateCompletion(openai.CompletionRequest, openai.CompletionResponse) error
	CreateChatCompletion(openai.ChatCompletionRequest, openai.ChatCompletionResponse) error
	Edits(openai.EditsRequest, openai.EditsResponse) error

	CreateEmbeddings(openai.EmbeddingRequest, openai.EmbeddingResponse) error

	ListFiles(openai.FilesList) error
	CreateFile(openai.FileRequest, openai.File) error
	DeleteFile(string) error
	GetFile(string, openai.File) error
	// GetFileContent: Not implemented in Go SDK

	CreateFineTune(openai.FineTuneRequest, openai.FineTune) error
	ListFineTunes(openai.FineTuneList) error
	GetFineTune(string, openai.FineTune) error
	CancelFineTune(string, openai.FineTune) error
	ListFineTuneEvents(string, openai.FineTuneEventList) error
	DeleteFineTune(string) error

	CreateImage(openai.ImageRequest, openai.ImageResponse) error
	CreateEditImage(openai.ImageEditRequest, openai.ImageResponse) error
	CreateVariImage(openai.ImageVariRequest, openai.ImageResponse) error

	ListModels(openai.ModelsList) error

	Moderations(openai.ModerationRequest, openai.ModerationResponse) error
}
