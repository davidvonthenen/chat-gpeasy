// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package proxy

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	klog "k8s.io/klog/v2"
)

func (p *ChatGPTProxy) postCompletion(c *gin.Context) {
	klog.V(6).Infof("postCompletion ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var completionRequest openai.CompletionRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&completionRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postCompletion LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.CreateCompletion(ctx, completionRequest)
	if err != nil {
		klog.V(6).Infof("client.CreateCompletion failed. Err: %v\n", err)
		klog.V(6).Infof("postCompletion LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "completion failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("CreateCompletion Callback...\n")
		err = (*p.callback).CreateCompletion(completionRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] CreateCompletion failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postCompletion Succeeded\n")
	klog.V(6).Infof("postCompletion LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}

func (p *ChatGPTProxy) postChatCompletion(c *gin.Context) {
	klog.V(6).Infof("postChatCompletion ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var completionRequest openai.ChatCompletionRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&completionRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postChatCompletion LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.CreateChatCompletion(ctx, completionRequest)
	if err != nil {
		klog.V(6).Infof("client.CreateChatCompletion failed. Err: %v\n", err)
		klog.V(6).Infof("postChatCompletion LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "chat completion failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("CreateChatCompletion Callback...\n")
		err = (*p.callback).CreateChatCompletion(completionRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] CreateChatCompletion failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postChatCompletion Succeeded\n")
	klog.V(6).Infof("postChatCompletion LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}

func (p *ChatGPTProxy) postEdits(c *gin.Context) {
	klog.V(6).Infof("postEdits ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var editsRequest openai.EditsRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&editsRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postEdits LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.Edits(ctx, editsRequest)
	if err != nil {
		klog.V(6).Infof("client.Edits failed. Err: %v\n", err)
		klog.V(6).Infof("postEdits LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "edits failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("Edits Callback...\n")
		err = (*p.callback).Edits(editsRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] Edits failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postEdits Succeeded\n")
	klog.V(6).Infof("postEdits LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}
