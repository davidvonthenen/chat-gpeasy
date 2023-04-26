// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package proxy

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	klog "k8s.io/klog/v2"
)

func New(options ProxyOptions) (*ChatGPTProxy, error) {
	if options.BindPort == 0 {
		options.BindPort = DefaultPort
	}

	var openAiApiKey string
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		klog.V(4).Info("OPENAI_API_KEY found")
		openAiApiKey = v
	} else {
		klog.Errorf("OPENAI_API_KEY not found\n")
		return nil, ErrInvalidInput
	}

	proxy := &ChatGPTProxy{
		options:      &options,
		openAiApiKey: openAiApiKey,
		callback:     options.Callback,
	}
	return proxy, nil
}

func (p *ChatGPTProxy) Init() error {
	klog.V(6).Infof("ChatGPTProxy.Init ENTER\n")

	client := openai.NewClient(p.openAiApiKey)
	if client == nil {
		klog.V(1).Infof("openai.NewClient is nil\n")
		klog.V(6).Infof("ChatGPTProxy.Init LEAVE\n")
		return ErrInvalidOpenAiClient
	}

	ctx := context.Background()
	modelList, err := client.ListModels(ctx)
	if err != nil {
		klog.V(6).Infof("client.ListModels failed. Err: %v\n", err)
		klog.V(6).Infof("ChatGPTProxy.Init LEAVE\n")
		return err
	}

	klog.V(6).Infof("Model List\n")
	klog.V(6).Infof("-------------------------------------\n")
	for _, model := range modelList.Models {
		klog.V(6).Infof("Model: %s\n", model.ID)
	}
	klog.V(6).Infof("\n")

	// housekeeping
	p.chatgptClient = client

	klog.V(4).Infof("ChatGPTProxy.Init Succeeded\n")
	klog.V(6).Infof("ChatGPTProxy.Init LEAVE\n")

	return nil
}

func (p *ChatGPTProxy) Start() error {
	klog.V(6).Infof("ChatGPTProxy.Start ENTER\n")

	// redirect
	router := gin.Default()
	router.GET("/v1/models", p.getModels)
	// router.GET("/v1/models/:model", p.getModel) // NOT IMPLEMENTED IN GO SDK
	router.POST("/v1/completions", p.postCompletion)
	router.POST("/v1/chat/completions", p.postChatCompletion)
	router.POST("/v1/edits", p.postEdits)
	router.POST("/v1/images/generations", p.postCreateImage)
	router.POST("/v1/images/edits", p.postEditImage)
	router.POST("/v1/images/variations", p.postVariationImage)
	router.POST("/v1/embeddings", p.postEmbedding)
	router.POST("/v1/audio/transcriptions", p.postTranscription)
	router.POST("/v1/audio/translations", p.postTranslation)
	router.GET("/v1/files", p.getFiles)
	router.POST("/v1/files", p.postCreateFile)
	router.DELETE("/v1/files/:file_id", p.deleteFile)
	router.GET("/v1/files/:file_id", p.getFile)
	// router.GET("/v1/files/:file_id/content", p.getFileContent) // NOT IMPLEMENTED IN GO SDK
	router.POST("/v1/fine-tunes", p.postCreateFineTune)
	router.GET("/v1/fine-tunes", p.getFineTunes)
	router.GET("/v1/fine-tunes/:fine_tune_id", p.getFineTune)
	router.POST("/v1/fine-tunes/:fine_tune_id", p.postCancelFineTune)
	router.GET("/v1/fine-tunes/:fine_tune_id/events", p.getFineTuneEvent)
	router.DELETE("/v1/fine-tunes/:fine_tune_id", p.deleteFineTune)
	router.POST("/v1/moderations", p.postModeration)

	// server
	p.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", p.options.BindPort),
		Handler: router,
	}

	// start the main entry endpoint to direct traffic
	go func() {
		// this is a blocking call
		err := p.server.ListenAndServeTLS(p.options.CrtFile, p.options.KeyFile)
		if err != nil {
			klog.V(6).Infof("ListenAndServeTLS server stopped. Err: %v\n", err)
		}
	}()

	// TODO: start metrics and tracing

	klog.V(4).Infof("ChatGPTProxy.Start Succeeded\n")
	klog.V(6).Infof("ChatGPTProxy.Start LEAVE\n")

	return nil
}

func (p *ChatGPTProxy) Stop() error {
	klog.V(6).Infof("ChatGPTProxy.Stop ENTER\n")

	// TODO: stop metrics and tracing

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.server.Shutdown(ctx); err != nil {
		klog.V(1).Infof("Server Shutdown Failed. Err: %v\n", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		klog.V(1).Infof("timeout of 5 seconds.")
	}

	klog.V(4).Infof("ChatGPTProxy.Stop Succeeded\n")
	klog.V(6).Infof("ChatGPTProxy.Stop LEAVE\n")

	return nil
}

func (p *ChatGPTProxy) Teardown() error {
	klog.V(6).Infof("ChatGPTProxy.Teardown ENTER\n")

	// TODO: stop metrics and tracing

	// release client
	p.chatgptClient = nil

	klog.V(4).Infof("ChatGPTProxy.Teardown Succeeded\n")
	klog.V(6).Infof("ChatGPTProxy.Teardown LEAVE\n")

	return nil
}
