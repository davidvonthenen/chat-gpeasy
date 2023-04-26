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

func (p *ChatGPTProxy) postCreateImage(c *gin.Context) {
	klog.V(6).Infof("postCreateImage ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var imageRequest openai.ImageRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&imageRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postCreateImage LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.CreateImage(ctx, imageRequest)
	if err != nil {
		klog.V(6).Infof("client.CreateImage failed. Err: %v\n", err)
		klog.V(6).Infof("postCreateImage LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "image creation failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("CreateImage Callback...\n")
		err = (*p.callback).CreateImage(imageRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] CreateImage failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postCreateImage Succeeded\n")
	klog.V(6).Infof("postCreateImage LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}

func (p *ChatGPTProxy) postEditImage(c *gin.Context) {
	klog.V(6).Infof("postEditImage ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var imageRequest openai.ImageEditRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&imageRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postEditImage LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.CreateEditImage(ctx, imageRequest)
	if err != nil {
		klog.V(6).Infof("client.CreateEditImage failed. Err: %v\n", err)
		klog.V(6).Infof("postEditImage LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "image creation failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("CreateEditImage Callback...\n")
		err = (*p.callback).CreateEditImage(imageRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] CreateEditImage failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postEditImage Succeeded\n")
	klog.V(6).Infof("postEditImage LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}

func (p *ChatGPTProxy) postVariationImage(c *gin.Context) {
	klog.V(6).Infof("postVariationImage ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var imageRequest openai.ImageVariRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&imageRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postVariationImage LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.CreateVariImage(ctx, imageRequest)
	if err != nil {
		klog.V(6).Infof("client.CreateVariImage failed. Err: %v\n", err)
		klog.V(6).Infof("postVariationImage LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "image variation failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("CreateVariImage Callback...\n")
		err = (*p.callback).CreateVariImage(imageRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] CreateVariImage failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postVariationImage Succeeded\n")
	klog.V(6).Infof("postVariationImage LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}
