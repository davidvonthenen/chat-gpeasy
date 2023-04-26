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

func (p *ChatGPTProxy) postCreateFineTune(c *gin.Context) {
	klog.V(6).Infof("postCreateFineTune ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var finetuneRequest openai.FineTuneRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&finetuneRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postCreateFineTune LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.CreateFineTune(ctx, finetuneRequest)
	if err != nil {
		klog.V(6).Infof("client.CreateFineTune failed. Err: %v\n", err)
		klog.V(6).Infof("postCreateFineTune LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "create file failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("CreateFineTune Callback...\n")
		err = (*p.callback).CreateFineTune(finetuneRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] CreateFineTune failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postCreateFineTune Succeeded\n")
	klog.V(6).Infof("postCreateFineTune LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}

func (p *ChatGPTProxy) getFineTunes(c *gin.Context) {
	klog.V(6).Infof("getFineTunes ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	fileList, err := p.chatgptClient.ListFineTunes(ctx)
	if err != nil {
		klog.V(1).Infof("client.ListFineTunes failed. Err: %v\n", err)
		klog.V(6).Infof("getFineTunes LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "list fine-tunes failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("ListFineTunes Callback...\n")
		err = (*p.callback).ListFineTunes(fileList)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] ListFineTunes failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("getFineTunes Succeeded\n")
	klog.V(6).Infof("getFineTunes LEAVE\n")
	c.IndentedJSON(http.StatusOK, fileList)
}

func (p *ChatGPTProxy) getFineTune(c *gin.Context) {
	klog.V(6).Infof("getFineTune ENTER\n")

	finetuneID := c.Param("fine_tune_id")

	klog.V(5).Infof("finetuneID: %s\n", finetuneID)
	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	resp, err := p.chatgptClient.GetFineTune(ctx, finetuneID)
	if err != nil {
		klog.V(6).Infof("client.GetFineTune failed. Err: %v\n", err)
		klog.V(6).Infof("getFineTune LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "get fine-tune failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("GetFineTune Callback...\n")
		err = (*p.callback).GetFineTune(finetuneID, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] GetFineTune failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("getFineTune Succeeded\n")
	klog.V(6).Infof("getFineTune LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}

func (p *ChatGPTProxy) postCancelFineTune(c *gin.Context) {
	klog.V(6).Infof("postCancelFineTune ENTER\n")

	finetuneID := c.Param("fine_tune_id")

	klog.V(5).Infof("finetuneID: %s\n", finetuneID)
	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	resp, err := p.chatgptClient.CancelFineTune(ctx, finetuneID)
	if err != nil {
		klog.V(6).Infof("client.CancelFineTune failed. Err: %v\n", err)
		klog.V(6).Infof("postCancelFineTune LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "cancel fine-tune failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("CancelFineTune Callback...\n")
		err = (*p.callback).CancelFineTune(finetuneID, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] CancelFineTune failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postCancelFineTune Succeeded\n")
	klog.V(6).Infof("postCancelFineTune LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}

func (p *ChatGPTProxy) getFineTuneEvent(c *gin.Context) {
	klog.V(6).Infof("getFineTuneEvent ENTER\n")

	finetuneID := c.Param("fine_tune_id")

	klog.V(5).Infof("finetuneID: %s\n", finetuneID)
	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	resp, err := p.chatgptClient.ListFineTuneEvents(ctx, finetuneID)
	if err != nil {
		klog.V(6).Infof("client.ListFineTuneEvents failed. Err: %v\n", err)
		klog.V(6).Infof("getFineTuneEvent LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "get fine-tune events failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("ListFineTuneEvents Callback...\n")
		err = (*p.callback).ListFineTuneEvents(finetuneID, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] ListFineTuneEvents failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("getFineTuneEvent Succeeded\n")
	klog.V(6).Infof("getFineTuneEvent LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}

func (p *ChatGPTProxy) deleteFineTune(c *gin.Context) {
	klog.V(6).Infof("deleteFineTune ENTER\n")

	finetuneID := c.Param("fine_tune_id")

	klog.V(5).Infof("finetuneID: %s\n", finetuneID)
	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	resp, err := p.chatgptClient.DeleteFineTune(ctx, finetuneID)
	if err != nil {
		klog.V(6).Infof("client.DeleteFineTune failed. Err: %v\n", err)
		klog.V(6).Infof("deleteFineTune LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "delete file failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("DeleteFineTune Callback...\n")
		err = (*p.callback).DeleteFineTune(finetuneID)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] DeleteFineTune failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("deleteFineTune Succeeded\n")
	klog.V(6).Infof("deleteFineTune LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}
