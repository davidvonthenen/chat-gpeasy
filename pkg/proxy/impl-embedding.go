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

func (p *ChatGPTProxy) postEmbedding(c *gin.Context) {
	klog.V(6).Infof("postEmbedding ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var embeddingRequest openai.EmbeddingRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&embeddingRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postEmbedding LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.CreateEmbeddings(ctx, embeddingRequest)
	if err != nil {
		klog.V(6).Infof("client.CreateEmbeddings failed. Err: %v\n", err)
		klog.V(6).Infof("postEmbedding LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "embedding failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("CreateEmbeddings Callback...\n")
		err = (*p.callback).CreateEmbeddings(embeddingRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] CreateEmbeddings failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postEmbedding Succeeded\n")
	klog.V(6).Infof("postEmbedding LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}
