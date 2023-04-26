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

func (p *ChatGPTProxy) getFiles(c *gin.Context) {
	klog.V(6).Infof("getFiles ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	fileList, err := p.chatgptClient.ListFiles(ctx)
	if err != nil {
		klog.V(1).Infof("client.ListFiles failed. Err: %v\n", err)
		klog.V(6).Infof("getFiles LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "list files failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("ListFiles Callback...\n")
		err = (*p.callback).ListFiles(fileList)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] ListFiles failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("getFiles Succeeded\n")
	klog.V(6).Infof("getFiles LEAVE\n")
	c.IndentedJSON(http.StatusOK, fileList)
}

func (p *ChatGPTProxy) postCreateFile(c *gin.Context) {
	klog.V(6).Infof("postCreateFile ENTER\n")

	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	var fileRequest openai.FileRequest

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&fileRequest); err != nil {
		klog.V(1).Infof("BindJSON failed. Err: %v\n", err)
		klog.V(6).Infof("postCreateFile LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	resp, err := p.chatgptClient.CreateFile(ctx, fileRequest)
	if err != nil {
		klog.V(6).Infof("client.CreateFile failed. Err: %v\n", err)
		klog.V(6).Infof("postCreateFile LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "create file failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("CreateFile Callback...\n")
		err = (*p.callback).CreateFile(fileRequest, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] CreateFile failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("postCreateFile Succeeded\n")
	klog.V(6).Infof("postCreateFile LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}

func (p *ChatGPTProxy) deleteFile(c *gin.Context) {
	klog.V(6).Infof("deleteFile ENTER\n")

	fileID := c.Param("file_id")

	klog.V(5).Infof("fileID: %s\n", fileID)
	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	err := p.chatgptClient.DeleteFile(ctx, fileID)
	if err != nil {
		klog.V(6).Infof("client.DeleteFile failed. Err: %v\n", err)
		klog.V(6).Infof("deleteFile LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "delete file failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("DeleteFile Callback...\n")
		err = (*p.callback).DeleteFile(fileID)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] DeleteFile failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("deleteFile Succeeded\n")
	klog.V(6).Infof("deleteFile LEAVE\n")
	c.IndentedJSON(http.StatusOK, gin.H{"message": "delete file succeeded"})
}

func (p *ChatGPTProxy) getFile(c *gin.Context) {
	klog.V(6).Infof("getFile ENTER\n")

	fileID := c.Param("file_id")

	klog.V(5).Infof("fileID: %s\n", fileID)
	for key, value := range c.Request.Header {
		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
	}

	ctx := context.Background()

	resp, err := p.chatgptClient.GetFile(ctx, fileID)
	if err != nil {
		klog.V(6).Infof("client.GetFile failed. Err: %v\n", err)
		klog.V(6).Infof("getFile LEAVE\n")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "get file failed"})
		return
	}

	if p.callback != nil {
		klog.V(6).Infof("GetFile Callback...\n")
		err = (*p.callback).GetFile(fileID, resp)
		if err != nil {
			klog.V(1).Infof("[CALLBACK] GetFile failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("getFile Succeeded\n")
	klog.V(6).Infof("getFile LEAVE\n")
	c.IndentedJSON(http.StatusOK, resp)
}

// func (p *ChatGPTProxy) getFileContent(c *gin.Context) {
// 	klog.V(6).Infof("getFileContent ENTER\n")

// 	fileID := c.Param("file_id")

// 	klog.V(5).Infof("fileID: %s\n", fileID)
// 	for key, value := range c.Request.Header {
// 		klog.V(5).Infof("HTTP Header: %s = %v\n", key, value)
// 	}

// 	ctx := context.Background()

// 	resp, err := p.chatgptClient.GetFile(ctx, fileID)
// 	if err != nil {
// 		klog.V(6).Infof("client.GetFile failed. Err: %v\n", err)
// 		klog.V(6).Infof("getFileContent LEAVE\n")
// 		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "get file content failed"})
// 		return
// 	}

// 	// TODO: hook here

// 	klog.V(4).Infof("getFileContent Succeeded\n")
// 	klog.V(6).Infof("getFileContent LEAVE\n")
// 	c.IndentedJSON(http.StatusOK, resp)
// }
