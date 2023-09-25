// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package advanced

import (
	"io"

	klog "k8s.io/klog/v2"
)

func (scc StreamingChatCompletion) Stream(w io.Writer) error {
	scc.sb.Reset()

	for {
		exit := false
		select {
		case <-scc.stopChan:
			return nil
		default:
			response, err := scc.stream.Recv()
			if err == io.EOF {
				klog.V(7).Infof("sc.stream.Recv() finished successfully\n")
				exit = true
				break
			}
			if err != nil {
				klog.V(1).Infof("sc.stream.Recv() failed. Err: %v\n", err)
				return err
			}

			if len(response.Choices) == 0 {
				continue
			}

			sentence := response.Choices[0].Delta.Content
			klog.V(7).Infof("sentence to pass to w.Write: %s\n", sentence)
			scc.sb.WriteString(sentence)

			byteCount, err := w.Write([]byte(sentence))
			if err != nil {
				klog.V(1).Infof("w.Write failed. Err: %v\n", err)
				return err
			}
			klog.V(7).Infof("io.Writer succeeded. Bytes written: %d\n", byteCount)
		}

		if exit {
			break
		}
	}

	klog.V(5).Infof("Stream result: %s\n", scc.sb.String())
	(*scc.callback).CommitResponse(scc.sb.String())
	close(scc.stopChan) // stop the stream explicit so close() can return

	return nil
}

func (scc StreamingChatCompletion) Close() error {
	scc.stream.Close()
	<-scc.stopChan
	return nil
}
