// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package personas

import (
	openai "github.com/sashabaranov/go-openai"
)

// PersonaOptions for the main HTTP endpoint
type PersonaOptions struct {
	*openai.ClientConfig

	DisableHostVerify bool
}
