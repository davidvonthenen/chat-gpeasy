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

func (p *ChatGPTProxy) postTranscription(c *gin.Context) {
	klog.V(6).Infof("postTranscription ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var audioRequest openai.AudioRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&audioRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postTranscription LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.CreateTranscription(ctx, audioRequest)
	if err != nil {
		klog.V(6).Infof("client.CreateTranscription failed. Err: %v\n", err)
		klog.V(6).Infof("postTranscription LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "transcription failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("CreateTranscription Callback...\n")
		err = (*p.callback).CreateTranscription(audioRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] CreateTranscription failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postTranscription Succeeded\n")
	klog.V(6).Infof("postTranscription LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}

func (p *ChatGPTProxy) postTranslation(c *gin.Context) {
	klog.V(6).Infof("postTranslation ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var audioRequest openai.AudioRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&audioRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postTranslation LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.CreateTranslation(ctx, audioRequest)
	if err != nil {
		klog.V(6).Infof("client.CreateTranslation failed. Err: %v\n", err)
		klog.V(6).Infof("postTranslation LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "transcription failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("CreateTranslation Callback...\n")
		err = (*p.callback).CreateTranslation(audioRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] CreateTranslation failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postTranslation Succeeded\n")
	klog.V(6).Infof("postTranslation LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}
