// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package proxy

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	klog "k8s.io/klog/v2"
)

func (p *ChatGPTProxy) getModels(c *gin.Context) {
	klog.V(6).Infof("getModels ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	modelList, err := p.chatgptClient.ListModels(ctx)
	if err != nil {
		klog.V(1).Infof("client.ListModels failed. Err: %v\n", err)
		klog.V(6).Infof("getModels LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "list models failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("ListModels Callback...\n")
		err = (*p.callback).ListModels(modelList)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] ListModels failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("getModels Succeeded\n")
	klog.V(6).Infof("getModels LEAVE\n")
	c.IndentedJSON(http.StatusOK, modelList)
}

// func (p *ChatGPTProxy) getModels(c *gin.Context) {
// 	klog.V(6).Infof("getModels ENTER\n")

// 	modelID := c.Param("model")

// 	klog.V(5).Infof("modelID: %s\n", modelID)
// 	for key, value := range c.Request.Header {
// 		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
// 	}

// 	ctx := context.Background()

// 	// TODO: hook here

// 	klog.V(4).Infof("getModels Succeeded\n")
// 	klog.V(6).Infof("getModels LEAVE\n")
// 	c.IndentedJSON(http.StatusOK, TODO)
// }
