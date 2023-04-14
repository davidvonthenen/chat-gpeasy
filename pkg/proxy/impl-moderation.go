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

func (p *ChatGPTProxy) postModeration(c *gin.Context) {
	klog.V(6).Infof("postModeration ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var moderationRequest openai.ModerationRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&moderationRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postModeration LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.Moderations(ctx, moderationRequest)
	if err != nil {
		klog.V(6).Infof("client.Moderations failed. Err: %v\n", err)
		klog.V(6).Infof("postModeration LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "post moderations failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("Moderations Callback...\n")
		err = (*p.callback).Moderations(moderationRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] LisModerationstModels failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postModeration Succeeded\n")
	klog.V(6).Infof("postModeration LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}
